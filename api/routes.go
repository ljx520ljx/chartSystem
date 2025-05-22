package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ljx520ljx/chartSystem/internal/middleware"
	"github.com/ljx520ljx/chartSystem/internal/service"
)

// RegisterRoutes 注册所有API路由
func RegisterRoutes(r *gin.Engine, services *service.Services) {
	// 认证中间件
	authMiddleware := middleware.AuthMiddleware(services.Auth)

	// 公开路由组
	public := r.Group("/api")
	{
		// 状态检查
		public.GET("/health", HealthCheck)

		// 认证相关路由
		auth := public.Group("/auth")
		{
			auth.POST("/login", HandleLogin(services))
			auth.POST("/register", HandleRegister(services))
		}
	}

	// 受保护的路由组
	protected := r.Group("/api")
	protected.Use(authMiddleware)
	{
		// 用户相关路由
		users := protected.Group("/users")
		{
			users.GET("/me", HandleGetCurrentUser(services))
			users.GET("/:id", HandleGetUser(services))
			users.PUT("/:id", HandleUpdateUser(services))
			users.DELETE("/:id", HandleDeleteUser(services))
			users.GET("", HandleListUsers(services))
			users.PUT("/:id/password", HandleChangePassword(services))
		}

		// 文件相关路由
		files := protected.Group("/files")
		{
			files.POST("", HandleUploadFile(services))
			files.GET("/:id", HandleGetFile(services))
			files.PUT("/:id", HandleUpdateFile(services))
			files.DELETE("/:id", HandleDeleteFile(services))
			files.GET("", HandleListFiles(services))
			files.POST("/:id/process", HandleProcessFile(services))
			files.GET("/:id/channels", HandleGetFileChannels(services))
			files.GET("/channels/:id/data", HandleGetChannelData(services))
			files.POST("/:id/markers", HandleAddMarker(services))
			files.GET("/:id/markers", HandleGetMarkers(services))
		}

		// 分析相关路由
		analyses := protected.Group("/analyses")
		{
			analyses.POST("", HandleCreateAnalysis(services))
			analyses.GET("/:id", HandleGetAnalysis(services))
			analyses.PUT("/:id", HandleUpdateAnalysis(services))
			analyses.DELETE("/:id", HandleDeleteAnalysis(services))
			analyses.GET("/file/:fileId", HandleGetFileAnalyses(services))
			analyses.POST("/:id/run", HandleRunAnalysis(services))
		}
	}
}

// HealthCheck 健康检查处理器
func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
} 