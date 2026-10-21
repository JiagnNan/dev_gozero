// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"

	"jwt/common"
	"jwt/internal/svc"
	"jwt/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshReq) (resp *types.LoginResp, err error) {
	// 登出后 refresh token 已被写入黑名单，必须在这里拦住。
	// 原因：/user/refresh 是公开路由（此时 access token 通常已过期，过不了 JWT 中间件），
	// 所以黑名单校验只能放在 logic 层。漏掉这步，登出后仍能无限换出新 token。
	revoked, err := l.svcCtx.Blacklist.IsRevoked(l.ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if revoked {
		return nil, errors.New("登录状态已失效，请重新登录")
	}

	// todo: add your logic here and delete this line
	user, _ := l.svcCtx.UserModel.FindOne(l.ctx, req.Id)

	accessToken, err := common.GenerateAccessToken(l.svcCtx.Config.RefreshSecret, user.Id, "admin")
	refreshToken, _ := common.GenerateRefreshToken(l.svcCtx.Config.RefreshSecret, req.RefreshToken, req.Id)

	return &types.LoginResp{
		Id:           user.Id,
		UserName:     "QY",
		Age:          user.Age,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    l.svcCtx.Config.Auth.AccessExpire,
	}, err
}
