// Package middleware 存放 go-zero REST 中间件。
package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"

	"jwt/internal/svc"
)

// CheckBlacklist 校验当前请求携带的 token 是否已被登出撤销。
//
// 装配位置很关键：必须挂在 rest.WithJwt 之后。go-zero 的路由装配顺序是
//
//	原生中间件 -> JWT 校验 -> server.Use 中间件 -> route.Handler（本中间件）
//
// 所以到这里时 token 的签名与有效期已被验证过，本中间件只回答一个问题：
// 「这个 token 是否已被显式登出」。这样未认证的请求不会打到 Redis。
func CheckBlacklist(svcCtx *svc.ServiceContext) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractBearerToken(r)
			if tokenStr == "" {
				next(w, r)
				return
			}

			revoked, err := svcCtx.Blacklist.IsRevoked(r.Context(), tokenStr)
			if err != nil {
				// fail-closed：Redis 故障时拒绝请求。
				// 若选择放行，已登出的 token 会重新生效，安全代价高于短暂不可用。
				logx.WithContext(r.Context()).Errorf("check token blacklist failed: %v", err)
				writeJSON(w, http.StatusServiceUnavailable, 503, "服务暂时不可用，请稍后重试")
				return
			}
			if revoked {
				writeJSON(w, http.StatusUnauthorized, 401, "登录状态已失效，请重新登录")
				return
			}

			next(w, r)
		}
	}
}

// extractBearerToken 从 Authorization: Bearer <token> 中取出 token。
func extractBearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return ""
	}

	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return strings.TrimSpace(parts[1])
}

func writeJSON(w http.ResponseWriter, statusCode, code int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, _ = fmt.Fprintf(w, `{"code":%d,"mesage":%q}`, code, message)
}
