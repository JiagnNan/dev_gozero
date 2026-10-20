package model

import (
	"context"
	"gorm.io/gorm"
)

type user_model struct {
	db *gorm.DB
}

type User struct {
	Id           int64  `gorm:"primary_key" json:"id"`
	Phone        string `gorm:"not null;unique_index:idx_phone" json:"phone"`
	Password     string `gorm:"not null;unique_index:idx_phone" json:"password"`
	Age          int    `gorm:"not null;unique_index:idx_age" json:"age"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

type UserRepository interface {
	FindOneByPhone(ctx context.Context, phone string) (User, error)
	FindOne(ctx context.Context, id int64) (User, error)
}
