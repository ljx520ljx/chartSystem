package repository

import (
	"errors"

	"github.com/ljx520ljx/chartSystem/internal/model"
	"gorm.io/gorm"
)

// UserRepoImpl 用户存储库实现
type UserRepoImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建用户存储库
func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepoImpl{db: db}
}

// Create 创建用户
func (r *UserRepoImpl) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// GetByID 通过ID获取用户
func (r *UserRepoImpl) GetByID(id uint) (*model.User, error) {
	var user model.User
	result := r.db.Preload("Role").Preload("Role.Permissions").First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetByUsername 通过用户名获取用户
func (r *UserRepoImpl) GetByUsername(username string) (*model.User, error) {
	var user model.User
	result := r.db.Preload("Role").Preload("Role.Permissions").Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetByEmail 通过邮箱获取用户
func (r *UserRepoImpl) GetByEmail(email string) (*model.User, error) {
	var user model.User
	result := r.db.Preload("Role").Preload("Role.Permissions").Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, result.Error
	}
	return &user, nil
}

// Update 更新用户
func (r *UserRepoImpl) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户
func (r *UserRepoImpl) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// List 获取用户列表
func (r *UserRepoImpl) List(offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	// 获取总数
	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := r.db.Preload("Role").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
} 