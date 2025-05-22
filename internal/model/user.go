package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"size:50;not null;unique"`
	Email     string    `json:"email" gorm:"size:100;not null;unique"`
	Password  string    `json:"-" gorm:"size:100;not null"`
	RoleID    uint      `json:"role_id" gorm:"not null;default:2"`
	Role      *Role     `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Files     []File    `json:"files,omitempty" gorm:"foreignKey:UserID"`
}

// BeforeSave 保存前加密密码
func (u *User) BeforeSave(tx *gorm.DB) error {
	// 只有在密码被修改时才进行加密
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// VerifyPassword 验证密码
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// Role 角色模型
type Role struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	Name        string       `json:"name" gorm:"size:50;not null;unique"`
	Description string       `json:"description" gorm:"size:255"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions"`
}

// Permission 权限模型
type Permission struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:50;not null;unique"`
	Description string    `json:"description" gorm:"size:255"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
