// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"context"
	"jwt/internal/config"
	"jwt/model"
)

type UserRepository interface {
	FindOneByPhone(ctx context.Context, phone string) (model.User, error)
	FindOne(ctx context.Context, id int64) (model.User, error)
}

type ServiceContext struct {
	Config    config.Config
	UserModel UserRepository
}

func NewServiceContext(c config.Config, v model.UserRepository) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		UserModel: v,
	}
}
