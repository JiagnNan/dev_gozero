package common

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

// revokedKeyPrefix 撤销名单的键前缀，统一前缀便于按前缀运维清理与监控。
const revokedKeyPrefix = "jwt:revoked:"

// TokenBlacklist 基于 Redis 的 JWT 撤销名单。
//
// JWT 是无状态凭据，服务端无法「删除」已签发的 token，因此登出只能采用
// 「黑名单 + TTL」方案：写入时把 TTL 设为该 token 的剩余有效期，
// token 自然过期后键自动消失，Redis 数据量不会无限增长。
type TokenBlacklist struct {
	rds *redis.Redis
}

// NewTokenBlacklist 构造撤销名单，复用 ServiceContext 中的 Redis 连接。
func NewTokenBlacklist(rds *redis.Redis) *TokenBlacklist {
	return &TokenBlacklist{rds: rds}
}

// Revoke 把 token 写入黑名单，ttl 应传该 token 的剩余有效期。
// ttl <= 0 表示 token 已过期，无需写入。
func (b *TokenBlacklist) Revoke(ctx context.Context, token string, ttl time.Duration) error {
	if ttl <= 0 {
		return nil
	}

	seconds := int(ttl.Seconds())
	if seconds < 1 {
		// SETEX 不接受 0 或负数的过期时间，向上取整到 1 秒
		seconds = 1
	}

	return b.rds.SetexCtx(ctx, revokedKey(token), "1", seconds)
}

// IsRevoked 判断 token 是否已被撤销（已登出）。
func (b *TokenBlacklist) IsRevoked(ctx context.Context, token string) (bool, error) {
	return b.rds.ExistsCtx(ctx, revokedKey(token))
}

// revokedKey 用 SHA-256 摘要作为键，原因有二：
//  1. JWT 长度通常 200~500 字节，直接当键浪费内存；摘要固定 64 字符；
//  2. Redis 被 dump 时不会泄露可直接使用的明文凭据。
func revokedKey(token string) string {
	sum := sha256.Sum256([]byte(token))
	return revokedKeyPrefix + hex.EncodeToString(sum[:])
}
