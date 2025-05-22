package config

import (
	"os"
	"strconv"
)

// AppConfig 应用配置结构体
type AppConfig struct {
	// 数据库配置
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	// Redis配置
	RedisHost string
	RedisPort int
	RedisDB   int
	RedisPass string

	// 应用配置
	JWTSecret     string
	APIPort       string
	FileStorePath string
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *AppConfig {
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))
	redisPort, _ := strconv.Atoi(getEnv("REDIS_PORT", "6379"))
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	return &AppConfig{
		// 数据库配置
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "chartsystem"),

		// Redis配置
		RedisHost: getEnv("REDIS_HOST", "localhost"),
		RedisPort: redisPort,
		RedisDB:   redisDB,
		RedisPass: getEnv("REDIS_PASSWORD", ""),

		// 应用配置
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
		APIPort:       getEnv("API_PORT", "8080"),
		FileStorePath: getEnv("FILE_STORE_PATH", "./uploads"),
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
