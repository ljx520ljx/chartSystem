package service

import (
	"github.com/ljx520ljx/chartSystem/internal/model"
	"github.com/ljx520ljx/chartSystem/internal/repository"
)

// Services 服务集合
type Services struct {
	Auth     AuthService
	User     UserService
	File     FileService
	Analysis AnalysisService
}

// NewServices 创建服务集合
func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Auth:     NewAuthService(repos),
		User:     NewUserService(repos),
		File:     NewFileService(repos),
		Analysis: NewAnalysisService(repos),
	}
}

// AuthService 认证服务接口
type AuthService interface {
	Login(usernameOrEmail, password string) (string, *model.User, error)
	Register(username, email, password string) (*model.User, error)
	ValidateToken(token string) (*model.User, error)
}

// UserService 用户服务接口
type UserService interface {
	GetByID(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id uint) error
	ListUsers(page, pageSize int) ([]*model.User, int64, error)
	ChangePassword(userID uint, oldPassword, newPassword string) error
}

// FileService 文件服务接口
type FileService interface {
	Upload(file *model.File) error
	GetByID(id uint) (*model.File, error)
	UpdateFile(file *model.File) error
	DeleteFile(id uint) error
	ListByUser(userID uint, page, pageSize int) ([]*model.File, int64, error)
	ListAll(page, pageSize int) ([]*model.File, int64, error)
	ProcessFile(fileID uint) error
	GetDataByChannel(channelID uint, startTime, endTime float64, maxPoints int) ([]float64, []float64, error)
	AddMarker(marker *model.Marker) error
	GetMarkers(fileID uint) ([]*model.Marker, error)
}

// AnalysisService 分析服务接口
type AnalysisService interface {
	CreateAnalysis(analysis *model.Analysis) error
	GetByID(id uint) (*model.Analysis, error)
	GetByFileID(fileID uint) ([]*model.Analysis, error)
	UpdateAnalysis(analysis *model.Analysis) error
	DeleteAnalysis(id uint) error
	RunAnalysis(analysisID uint) error
}
