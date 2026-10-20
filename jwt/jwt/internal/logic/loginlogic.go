// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"
	"jwt/commom"

	"jwt/internal/svc"
	"jwt/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// todo: add your logic here and delete this line
	user, err := l.svcCtx.UserModel.FindOneByPhone(l.ctx, req.Phone)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}
	if !commom.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}
	accessToken, _ := commom.GenerateAccessToken(l.svcCtx.Config.RefreshSecret, user.Id, "admin")
	refreshToken, _ := commom.GenerateRefreshToken(l.svcCtx.Config.RefreshSecret, user.RefreshToken, user.Id)

	return &types.LoginResp{
		Id:           user.Id,
		UserName:     "qy",
		Password:     user.Password,
		Age:          user.Age,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    l.svcCtx.Config.Auth.AccessExpire,
	}, nil
}
