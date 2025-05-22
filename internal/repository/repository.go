package repository

import (
	"github.com/go-redis/redis/v8"
	"github.com/ljx520ljx/chartSystem/internal/model"
	"gorm.io/gorm"
)

// Repositories 所有存储库的集合
type Repositories struct {
	User        UserRepository
	File        FileRepository
	DataChannel DataChannelRepository
	Role        RoleRepository
	Analysis    AnalysisRepository
	db          *gorm.DB
	rdb         *redis.Client
}

// NewRepositories 创建存储库集合
func NewRepositories(db *gorm.DB, rdb *redis.Client) *Repositories {
	return &Repositories{
		User:        NewUserRepository(db),
		File:        NewFileRepository(db),
		DataChannel: NewDataChannelRepository(db),
		Role:        NewRoleRepository(db),
		Analysis:    NewAnalysisRepository(db),
		db:          db,
		rdb:         rdb,
	}
}

// UserRepository 用户存储库接口
type UserRepository interface {
	Create(user *model.User) error
	GetByID(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
	List(offset, limit int) ([]*model.User, int64, error)
}

// FileRepository 文件存储库接口
type FileRepository interface {
	Create(file *model.File) error
	GetByID(id uint) (*model.File, error)
	Update(file *model.File) error
	Delete(id uint) error
	ListByUser(userID uint, offset, limit int) ([]*model.File, int64, error)
	List(offset, limit int) ([]*model.File, int64, error)
}

// DataChannelRepository 数据通道存储库接口
type DataChannelRepository interface {
	Create(channel *model.DataChannel) error
	GetByID(id uint) (*model.DataChannel, error)
	GetByFileID(fileID uint) ([]*model.DataChannel, error)
	Update(channel *model.DataChannel) error
	Delete(id uint) error
}

// RoleRepository 角色存储库接口
type RoleRepository interface {
	Create(role *model.Role) error
	GetByID(id uint) (*model.Role, error)
	GetByName(name string) (*model.Role, error)
	Update(role *model.Role) error
	Delete(id uint) error
	List() ([]*model.Role, error)
}

// AnalysisRepository 分析结果存储库接口
type AnalysisRepository interface {
	Create(analysis *model.Analysis) error
	GetByID(id uint) (*model.Analysis, error)
	GetByFileID(fileID uint) ([]*model.Analysis, error)
	Update(analysis *model.Analysis) error
	Delete(id uint) error
}
