// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"encoding/json"
	"errors"

	"jwt/internal/svc"
	"jwt/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProfileLogic {
	return &GetProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProfileLogic) GetProfile(req *types.ProfileReq) (resp *types.ProfileResp, err error) {
	// todo: add your logic here and delete this line
	userId, _ := l.ctx.Value("userId").(json.Number).Int64()
	role, _ := l.ctx.Value("role").(string)

	if role != "admin" && userId != req.Id {
		return nil, errors.New("无权访问")
	}
	user, _ := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	return &types.ProfileResp{
		UserName:     "qy",
		Password:     user.Password,
		Age:          user.Age,
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
	}, nil
}
