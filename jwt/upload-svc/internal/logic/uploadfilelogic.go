// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"upload-svc/internal/svc"
	"upload-svc/internal/types"
)

// maxUploadSize 业务级大小上限。
// Config.MaxBytes 是框架级闸门（整个请求体），这里再收紧一次：
// 一个请求可能带多个文件，框架放行不代表单个文件合规。
const maxUploadSize = 10 << 20 // 10 MB

type UploadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UploadFile 把上传流写入本地目录，返回可访问的相对 URL。
func (l *UploadFileLogic) UploadFile(file multipart.File, header *multipart.FileHeader) (*types.UploadResp, error) {
	if header.Size > maxUploadSize {
		return nil, fmt.Errorf("文件超过 %d MB 限制", maxUploadSize>>20)
	}

	// 首次上传时目录可能还不存在，MkdirAll 对已存在目录是幂等的
	if err := os.MkdirAll(l.svcCtx.Config.UploadDir, 0o755); err != nil {
		return nil, fmt.Errorf("创建上传目录失败: %w", err)
	}

	safeName := buildSafeName(header.Filename)
	dstPath := filepath.Join(l.svcCtx.Config.UploadDir, safeName)

	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	// 用 io.Copy 的返回值作为真实大小，不直接信任 header.Size
	size, err := io.Copy(dst, file)
	if err != nil {
		// 写一半失败必须删掉残留文件，否则目录里会堆积 0 字节或半截文件
		_ = os.Remove(dstPath)
		return nil, fmt.Errorf("写入文件失败: %w", err)
	}

	l.Infof("file saved: %s (%d bytes)", dstPath, size)

	return &types.UploadResp{
		Filename: header.Filename,
		Size:     size,
		URL:      "/files/" + safeName,
	}, nil
}

// buildSafeName 构造落盘文件名。
//
// 直接用用户传来的 Filename 拼路径有两个风险：
//  1. ../../etc/passwd 这类路径穿越；
//  2. Windows 上 \ 与 / 混用，文件可能被写到目标目录之外。
//
// filepath.Base + Clean 只保留最后一段，再把残留分隔符替换掉即可堵住。
// 前置时间戳是为了同名文件不互相覆盖。
func buildSafeName(raw string) string {
	base := filepath.Base(filepath.Clean(raw))
	base = strings.NewReplacer("/", "_", "\\", "_").Replace(base)
	// 去掉开头的点，防止生成 .env / .. 这类特殊名字
	base = strings.TrimLeft(base, ".")
	if base == "" {
		base = "unnamed"
	}

	return fmt.Sprintf("%d_%s", time.Now().UnixNano(), base)
}
