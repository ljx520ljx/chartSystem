package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ljx520ljx/chartSystem/internal/model"
	"github.com/ljx520ljx/chartSystem/internal/repository"
)

// AuthServiceImpl 认证服务实现
type AuthServiceImpl struct {
	repos     *repository.Repositories
	jwtSecret string
}

// JWTClaims JWT声明结构
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// NewAuthService 创建认证服务
func NewAuthService(repos *repository.Repositories) AuthService {
	return &AuthServiceImpl{
		repos:     repos,
		jwtSecret: getJWTSecret(),
	}
}

// 获取JWT密钥
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "default-very-secret-key-should-be-changed"
	}
	return secret
}

// Login 用户登录
func (s *AuthServiceImpl) Login(usernameOrEmail, password string) (string, *model.User, error) {
	// 尝试通过用户名查找用户
	user, err := s.repos.User.GetByUsername(usernameOrEmail)
	if err != nil {
		// 如果通过用户名找不到，尝试通过邮箱查找
		user, err = s.repos.User.GetByEmail(usernameOrEmail)
		if err != nil {
			return "", nil, errors.New("用户名或密码错误")
		}
	}

	// 验证密码
	if !user.VerifyPassword(password) {
		return "", nil, errors.New("用户名或密码错误")
	}

	// 生成JWT令牌
	token, err := s.generateToken(user)
	if err != nil {
		return "", nil, fmt.Errorf("生成令牌失败: %w", err)
	}

	return token, user, nil
}

// Register 用户注册
func (s *AuthServiceImpl) Register(username, email, password string) (*model.User, error) {
	// 检查用户名是否已存在
	_, err := s.repos.User.GetByUsername(username)
	if err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	_, err = s.repos.User.GetByEmail(email)
	if err == nil {
		return nil, errors.New("邮箱已被注册")
	}

	// 创建用户
	user := &model.User{
		Username: username,
		Email:    email,
		Password: password,
		RoleID:   2, // 默认为普通用户角色
	}

	if err := s.repos.User.Create(user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return user, nil
}

// ValidateToken 验证令牌
func (s *AuthServiceImpl) ValidateToken(tokenString string) (*model.User, error) {
	// 解析JWT令牌
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析令牌失败: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("无效的令牌")
	}

	// 获取声明
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("无效的令牌声明")
	}

	// 获取用户
	user, err := s.repos.User.GetByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("找不到用户: %w", err)
	}

	return user, nil
}

// 生成JWT令牌
func (s *AuthServiceImpl) generateToken(user *model.User) (string, error) {
	// 设置令牌过期时间为24小时
	expirationTime := time.Now().Add(24 * time.Hour)

	// 创建声明
	claims := &JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
} 