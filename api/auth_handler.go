package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ljx520ljx/chartSystem/internal/service"
)

// LoginRequest 登录请求结构
type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// AuthResponse 认证响应结构
type AuthResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	} `json:"user"`
}

// HandleLogin 处理登录请求
func HandleLogin(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
			return
		}

		// 调用登录服务
		token, user, err := services.Auth.Login(req.UsernameOrEmail, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// 构造响应
		response := AuthResponse{
			Token: token,
		}
		response.User.ID = user.ID
		response.User.Username = user.Username
		response.User.Email = user.Email
		if user.Role != nil {
			response.User.Role = user.Role.Name
		}

		c.JSON(http.StatusOK, response)
	}
}

// HandleRegister 处理注册请求
func HandleRegister(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
			return
		}

		// 调用注册服务
		user, err := services.Auth.Register(req.Username, req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 登录新注册的用户
		token, _, err := services.Auth.Login(req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "注册成功但登录失败"})
			return
		}

		// 构造响应
		response := AuthResponse{
			Token: token,
		}
		response.User.ID = user.ID
		response.User.Username = user.Username
		response.User.Email = user.Email
		if user.Role != nil {
			response.User.Role = user.Role.Name
		}

		c.JSON(http.StatusCreated, response)
	}
} 