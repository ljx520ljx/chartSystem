package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ljx520ljx/chartSystem/internal/model"
	"github.com/ljx520ljx/chartSystem/internal/service"
)

// UpdateUserRequest 更新用户请求结构
type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// ChangePasswordRequest 修改密码请求结构
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// HandleGetCurrentUser 获取当前用户信息
func HandleGetCurrentUser(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// HandleGetUser 获取用户信息
func HandleGetUser(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取URL参数中的用户ID
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}

		// 通过服务获取用户信息
		user, err := services.User.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// HandleUpdateUser 更新用户信息
func HandleUpdateUser(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取URL参数中的用户ID
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}

		// 从上下文中获取当前用户
		currentUser, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
			return
		}
		authUser := currentUser.(*model.User)

		// 检查权限（只能修改自己的信息，除非是管理员）
		if authUser.ID != uint(id) && authUser.Role.Name != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限修改此用户信息"})
			return
		}

		// 获取请求数据
		var req UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
			return
		}

		// 获取要更新的用户
		user, err := services.User.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			return
		}

		// 更新用户信息
		if req.Username != "" {
			user.Username = req.Username
		}
		if req.Email != "" {
			user.Email = req.Email
		}

		// 保存更新
		if err := services.User.UpdateUser(user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户失败"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// HandleDeleteUser 删除用户
func HandleDeleteUser(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取URL参数中的用户ID
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}

		// 从上下文中获取当前用户
		currentUser, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
			return
		}
		authUser := currentUser.(*model.User)

		// 检查权限（只能删除自己，除非是管理员）
		if authUser.ID != uint(id) && authUser.Role.Name != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除此用户"})
			return
		}

		// 执行删除
		if err := services.User.DeleteUser(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "用户已删除"})
	}
}

// HandleListUsers 获取用户列表
func HandleListUsers(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取当前用户
		currentUser, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
			return
		}
		authUser := currentUser.(*model.User)

		// 检查权限（只有管理员可以查看所有用户）
		if authUser.Role.Name != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限查看用户列表"})
			return
		}

		// 获取分页参数
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

		// 获取用户列表
		users, total, err := services.User.ListUsers(page, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"users": users,
			"total": total,
			"page":  page,
			"size":  pageSize,
		})
	}
}

// HandleChangePassword 修改密码
func HandleChangePassword(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取URL参数中的用户ID
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
			return
		}

		// 从上下文中获取当前用户
		currentUser, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
			return
		}
		authUser := currentUser.(*model.User)

		// 检查权限（只能修改自己的密码）
		if authUser.ID != uint(id) {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限修改此用户密码"})
			return
		}

		// 获取请求数据
		var req ChangePasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
			return
		}

		// 修改密码
		if err := services.User.ChangePassword(uint(id), req.OldPassword, req.NewPassword); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "密码已修改"})
	}
} 