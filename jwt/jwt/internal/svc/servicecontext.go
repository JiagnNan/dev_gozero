// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"context"
	"jwt/common"
	"jwt/internal/config"
	"jwt/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

type UserRepository interface {
	FindOneByPhone(ctx context.Context, phone string) (model.User, error)
	FindOne(ctx context.Context, id int64) (model.User, error)
}

type ServiceContext struct {
	Config    config.Config
	UserModel UserRepository
	// Redis 供业务侧直接使用（如登录失败计数、在线状态等）。
	Redis *redis.Redis
	// Blacklist 封装 JWT 撤销名单的读写，登出与鉴权中间件共用同一实例。
	Blacklist *common.TokenBlacklist
}

func NewServiceContext(c config.Config, v model.UserRepository) *ServiceContext {
	rds := redis.MustNewRedis(c.Redis)

	return &ServiceContext{
		Config:    c,
		UserModel: v,
		Redis:     rds,
		Blacklist: common.NewTokenBlacklist(rds),
	}
}
