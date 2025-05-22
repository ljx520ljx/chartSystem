package utils

import (
	"log"

	"github.com/ljx520ljx/chartSystem/internal/model"
	"gorm.io/gorm"
)

// InitDatabase 初始化数据库表和默认数据
func InitDatabase(db *gorm.DB) error {
	// 自动迁移表结构
	if err := migrateSchema(db); err != nil {
		return err
	}

	// 添加默认角色和权限
	if err := seedRolesAndPermissions(db); err != nil {
		return err
	}

	// 添加默认管理员用户
	if err := seedAdminUser(db); err != nil {
		return err
	}

	log.Println("数据库初始化完成")
	return nil
}

// 迁移数据库表结构
func migrateSchema(db *gorm.DB) error {
	// 创建数据表
	return db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.File{},
		&model.DataChannel{},
		&model.FileProcessing{},
		&model.Marker{},
		&model.Analysis{},
	)
}

// 添加默认角色和权限
func seedRolesAndPermissions(db *gorm.DB) error {
	// 定义权限
	permissions := []model.Permission{
		{Name: "user:read", Description: "查看用户信息"},
		{Name: "user:write", Description: "修改用户信息"},
		{Name: "user:delete", Description: "删除用户"},
		{Name: "file:read", Description: "查看文件"},
		{Name: "file:write", Description: "上传和修改文件"},
		{Name: "file:delete", Description: "删除文件"},
		{Name: "file:process", Description: "处理文件"},
		{Name: "analysis:read", Description: "查看分析结果"},
		{Name: "analysis:write", Description: "创建和修改分析"},
		{Name: "analysis:delete", Description: "删除分析结果"},
	}

	// 创建权限
	for _, perm := range permissions {
		var existingPerm model.Permission
		if err := db.Where("name = ?", perm.Name).First(&existingPerm).Error; err != nil {
			if err := db.Create(&perm).Error; err != nil {
				return err
			}
		}
	}

	// 获取所有权限
	var allPermissions []model.Permission
	if err := db.Find(&allPermissions).Error; err != nil {
		return err
	}

	// 定义角色和权限关系
	roles := []struct {
		Role        model.Role
		Permissions []string
	}{
		{
			Role: model.Role{
				ID:          1,
				Name:        "admin",
				Description: "系统管理员",
			},
			Permissions: []string{
				"user:read", "user:write", "user:delete",
				"file:read", "file:write", "file:delete", "file:process",
				"analysis:read", "analysis:write", "analysis:delete",
			},
		},
		{
			Role: model.Role{
				ID:          2,
				Name:        "user",
				Description: "普通用户",
			},
			Permissions: []string{
				"user:read",
				"file:read", "file:write", "file:process",
				"analysis:read", "analysis:write",
			},
		},
	}

	// 创建角色并分配权限
	for _, roleData := range roles {
		var existingRole model.Role
		if err := db.Where("name = ?", roleData.Role.Name).First(&existingRole).Error; err != nil {
			// 如果角色不存在，创建新角色
			if err := db.Create(&roleData.Role).Error; err != nil {
				return err
			}
			existingRole = roleData.Role
		}

		// 添加权限关联
		var rolePermissions []model.Permission
		for _, permName := range roleData.Permissions {
			for _, perm := range allPermissions {
				if perm.Name == permName {
					rolePermissions = append(rolePermissions, perm)
					break
				}
			}
		}

		// 更新角色权限
		if err := db.Model(&existingRole).Association("Permissions").Replace(rolePermissions); err != nil {
			return err
		}
	}

	return nil
}

// 添加默认管理员用户
func seedAdminUser(db *gorm.DB) error {
	var adminUser model.User
	if err := db.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		// 创建管理员用户
		adminUser = model.User{
			Username: "admin",
			Email:    "admin@example.com",
			Password: "admin123",
			RoleID:   1, // 管理员角色
		}
		return db.Create(&adminUser).Error
	}
	return nil
} 