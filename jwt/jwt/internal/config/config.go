// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	RefreshSecret string
	RefreshExpire int64
	// Redis 用于存放 JWT 撤销名单（登出后的 token 指纹）。
	// 放在配置里而不是硬编码，便于区分开发/测试/生产实例。
	Redis redis.RedisConf
}
