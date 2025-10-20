package controller

import (
	"testing"
	"ginx-example/models"
	"ginx-example/global"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB 设置测试数据库
func setupTestDB() {
	// 使用内存SQLite数据库进行测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	
	global.DB = db
	
	// 迁移表结构
	if err := models.AutoMigrate(); err != nil {
		panic("failed to migrate database")
	}
}

func TestUserController_CreateAndGetUser(t *testing.T) {
	// 设置测试数据库
	setupTestDB()
	
	controller := NewUserController()
	
	// 创建测试用户
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "测试用户",
		Status:   "active",
	}
	
	// 测试创建用户
	err := controller.CreateUser(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	
	// 测试获取用户
	retrievedUser, err := controller.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, retrievedUser.Username)
	assert.Equal(t, user.Email, retrievedUser.Email)
	assert.Equal(t, user.Name, retrievedUser.Name)
}

func TestUserController_GetUserNotFound(t *testing.T) {
	// 设置测试数据库
	setupTestDB()
	
	controller := NewUserController()
	
	// 测试获取不存在的用户
	_, err := controller.GetUserByID(999)
	assert.Error(t, err)
}

func TestUserController_ListUsers(t *testing.T) {
	// 设置测试数据库
	setupTestDB()
	
	controller := NewUserController()
	
	// 创建测试用户
	users := []*models.User{
		{Username: "user1", Email: "user1@example.com", Name: "用户1", Status: "active"},
		{Username: "user2", Email: "user2@example.com", Name: "用户2", Status: "active"},
	}
	
	for _, user := range users {
		err := controller.CreateUser(user)
		assert.NoError(t, err)
	}
	
	// 测试获取用户列表
	userList, total, err := controller.ListUsers(0, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, userList, 2)
}