package service

import (
	"errors"
	"fmt"

	"github.com/ljx520ljx/chartSystem/internal/model"
	"github.com/ljx520ljx/chartSystem/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserServiceImpl 用户服务实现
type UserServiceImpl struct {
	repos *repository.Repositories
}

// NewUserService 创建用户服务
func NewUserService(repos *repository.Repositories) UserService {
	return &UserServiceImpl{repos: repos}
}

// GetByID 通过ID获取用户
func (s *UserServiceImpl) GetByID(id uint) (*model.User, error) {
	return s.repos.User.GetByID(id)
}

// GetByUsername 通过用户名获取用户
func (s *UserServiceImpl) GetByUsername(username string) (*model.User, error) {
	return s.repos.User.GetByUsername(username)
}

// UpdateUser 更新用户信息
func (s *UserServiceImpl) UpdateUser(user *model.User) error {
	// 检查用户名是否已被其他用户使用
	existingUser, err := s.repos.User.GetByUsername(user.Username)
	if err == nil && existingUser.ID != user.ID {
		return errors.New("用户名已被使用")
	}

	// 检查邮箱是否已被其他用户使用
	existingUser, err = s.repos.User.GetByEmail(user.Email)
	if err == nil && existingUser.ID != user.ID {
		return errors.New("邮箱已被使用")
	}

	// 更新用户
	return s.repos.User.Update(user)
}

// DeleteUser 删除用户
func (s *UserServiceImpl) DeleteUser(id uint) error {
	return s.repos.User.Delete(id)
}

// ListUsers 获取用户列表
func (s *UserServiceImpl) ListUsers(page, pageSize int) ([]*model.User, int64, error) {
	offset := (page - 1) * pageSize
	return s.repos.User.List(offset, pageSize)
}

// ChangePassword 修改用户密码
func (s *UserServiceImpl) ChangePassword(userID uint, oldPassword, newPassword string) error {
	// 获取用户
	user, err := s.repos.User.GetByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if !user.VerifyPassword(oldPassword) {
		return errors.New("旧密码不正确")
	}

	// 检查新密码长度
	if len(newPassword) < 6 {
		return errors.New("新密码长度不能少于6个字符")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	user.Password = string(hashedPassword)
	return s.repos.User.Update(user)
} 