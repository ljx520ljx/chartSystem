package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ljx520ljx/chartSystem/api"
	"github.com/ljx520ljx/chartSystem/config"
	"github.com/ljx520ljx/chartSystem/internal/repository"
	"github.com/ljx520ljx/chartSystem/internal/service"
	"github.com/ljx520ljx/chartSystem/internal/utils"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("未找到.env文件，使用默认配置")
	}

	// 初始化配置
	cfg := config.LoadConfig()

	// 初始化数据库连接
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 初始化数据库表和默认数据
	if err := utils.InitDatabase(db); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 初始化Redis连接
	rdb := config.InitRedis(cfg)

	// 初始化存储库
	repos := repository.NewRepositories(db, rdb)

	// 初始化服务
	services := service.NewServices(repos)

	// 设置Gin模式
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin路由器
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 创建上传目录
	uploadPath := cfg.FileStorePath
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadPath, 0755); err != nil {
			log.Fatalf("创建上传目录失败: %v", err)
		}
	}

	// 注册API路由
	api.RegisterRoutes(r, services)

	// 启动服务器
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("服务器启动在 http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
