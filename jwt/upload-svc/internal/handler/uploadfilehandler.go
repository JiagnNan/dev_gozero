// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"upload-svc/internal/logic"
	"upload-svc/internal/svc"
)

// allowedTypes 允许上传的 MIME 白名单。
// 只看 MIME、不看扩展名：扩展名是用户随便改的，伪造成本为零。
var allowedTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/gif":       true,
	"image/webp":      true,
	"application/pdf": true,
	"text/plain":      true,
}

func UploadFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ParseMultipartForm 的参数是「内存阈值」而非大小上限：
		// 32MB 以内的部分留在内存，超出部分自动溢写到系统临时目录。
		// 真正的上传大小上限由 Config.MaxBytes 控制。
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		// 必须清理溢写产生的临时文件，否则每次大文件上传都会在 %TEMP% 留一份副本。
		defer func() {
			if r.MultipartForm != nil {
				_ = r.MultipartForm.RemoveAll()
			}
		}()

		file, header, err := r.FormFile("file")
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		defer file.Close()

		// 读前 512 字节嗅探真实 MIME（http.DetectContentType 的固定要求）
		buf := make([]byte, 512)
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		mimeType := http.DetectContentType(buf[:n])
		if !allowedTypes[mimeType] {
			// 注意：go-zero v1.10.x 的签名是
			//   httpx.Error(w, err, fns ...func(http.ResponseWriter, error))
			// 第三个参数是可变函数而不是状态码，官方文档里的
			//   httpx.Error(w, err, http.StatusUnsupportedMediaType)
			// 会编译失败。需要指定状态码时直接用标准库的 http.Error。
			http.Error(w, fmt.Sprintf("不支持的文件类型: %s", mimeType), http.StatusUnsupportedMediaType)
			return
		}

		// 嗅探读掉了前 512 字节，必须把游标拨回起点，
		// 否则落盘的文件会缺失头部（图片直接损坏，且从返回的 size 看不出来）。
		if _, err = file.Seek(0, io.SeekStart); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUploadFileLogic(r.Context(), svcCtx)
		resp, err := l.UploadFile(file, header)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
