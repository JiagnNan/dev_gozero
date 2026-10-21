// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"jwt/internal/svc"
	"jwt/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Logout 把 access token 与 refresh token 写入 Redis 撤销名单。
//
// 两个 token 必须同时撤销：只拉黑 access token 的话，攻击者仍可用 refresh token
// 换取新的 access token，登出等于没做。这也是 JWT 登出最容易被漏掉的点。
func (l *LogoutLogic) Logout(req *types.LogoutReq) (*types.MesageResp, error) {
	if err := l.revoke(req.AccessToken, l.svcCtx.Config.Auth.AccessSecret); err != nil {
		return nil, err
	}
	if err := l.revoke(req.RefreshToken, l.svcCtx.Config.RefreshSecret); err != nil {
		return nil, err
	}

	l.Infof("logout success, userId=%v", l.ctx.Value("user_id"))

	return &types.MesageResp{
		Code:   0,
		Mesage: "登出成功",
	}, nil
}

// revoke 解析 token 剩余有效期并写入黑名单。
//
// 对非法 token 与已过期 token 一律静默跳过，原因：
//   - 登出接口必须幂等，重复登出、拿过期 token 登出都不应报错；
//   - 签名校验失败的 token 不能写进黑名单，否则会污染 Redis（任意人可塞键）。
func (l *LogoutLogic) revoke(tokenStr, secret string) error {
	if tokenStr == "" {
		return nil
	}

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		// 限定签名算法，防 alg=none 与算法混淆攻击
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		l.Errorf("skip revoking invalid token: %v", err)
		return nil
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		// token 没有 exp（理论上不该出现），无法确定 TTL，退化为不拉黑
		l.Errorf("token has no exp claim, skip revoking")
		return nil
	}

	// TTL 必须精确等于 token 剩余有效期，否则键会早于 token 失效（登出失效）
	// 或晚于 token 失效（Redis 里堆积无用键）。
	ttl := time.Until(time.Unix(int64(exp), 0))
	if ttl <= 0 {
		return nil
	}

	return l.svcCtx.Blacklist.Revoke(l.ctx, tokenStr, ttl)
}
