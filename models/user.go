// models/user.go
package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:50;unique;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Role      int       `gorm:"not null;default:1" json:"role"` // 1-普通用户，2-管理员
	RealName  string    `gorm:"size:50" json:"real_name"`
	Email     string    `gorm:"size:100" json:"email"`
	Phone     string    `gorm:"size:20" json:"phone"`
	LastLogin time.Time `json:"last_login"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
