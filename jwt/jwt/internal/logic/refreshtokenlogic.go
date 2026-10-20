// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
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
