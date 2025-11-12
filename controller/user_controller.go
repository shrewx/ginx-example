package controller

import (
	"context"
	"ginx-example/constants/status_error"
	"ginx-example/models"
	"github.com/pkg/errors"
	"github.com/shrewx/ginx"
	"github.com/shrewx/ginx/pkg/dbhelper"
	"github.com/shrewx/ginx/pkg/logx"
	"gorm.io/gorm"
)

// UserController 用户控制器
type UserController struct {
	db *gorm.DB
}

// NewUserController 创建用户控制器实例
func NewUserController(db *gorm.DB) *UserController {
	return &UserController{
		db: db,
	}
}

// GetUserByID 根据用户ID获取用户信息
func (uc *UserController) GetUserByID(userID int64) (*models.User, error) {
	var user models.User

	// 使用GORM查询用户
	result := uc.db.First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ginx.WithStack(status_error.UserNotFound.WithParams(map[string]interface{}{
				"UserID": userID,
			}))
		}
		// 记录数据库错误日志并返回统一错误
		logx.Errorf("获取用户信息失败, userID: %d, error: %v", userID, result.Error)
		return nil, ginx.WithStack(status_error.DataOperationFailed)
	}

	return &user, nil
}

// CreateUser 创建用户
func (uc *UserController) CreateUser(ctx context.Context, user *models.User) error {
	result := dbhelper.GetCtxDB(ctx, uc.db).Create(user)
	if result.Error != nil {
		// 记录数据库错误日志并返回统一错误
		logx.Errorf("创建用户失败, user: %+v, error: %v", user, result.Error)
		return ginx.WithStack(status_error.DataOperationFailed)
	}
	return nil
}

// UpdateUser 更新用户信息
func (uc *UserController) UpdateUser(ctx context.Context, user *models.User) error {
	result := dbhelper.GetCtxDB(ctx, uc.db).Save(user)
	if result.Error != nil {
		// 记录数据库错误日志并返回统一错误
		logx.Errorf("更新用户信息失败, user: %+v, error: %v", user, result.Error)
		return ginx.WithStack(status_error.DataOperationFailed)
	}
	return nil
}

// DeleteUser 删除用户（软删除）
func (uc *UserController) DeleteUser(ctx context.Context, userID int64) error {
	result := dbhelper.GetCtxDB(ctx, uc.db).Delete(&models.User{}, userID)
	if result.Error != nil {
		// 记录数据库错误日志并返回统一错误
		logx.Errorf("删除用户失败, userID: %d, error: %v", userID, result.Error)
		return ginx.WithStack(status_error.DataOperationFailed)
	}
	if result.RowsAffected == 0 {
		return ginx.WithStack(status_error.UserNotFound.WithParams(map[string]interface{}{
			"UserID": userID,
		}))
	}
	return nil
}

// ListUsers 获取用户列表
func (uc *UserController) ListUsers(offset, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// 获取总数
	if err := uc.db.Model(&models.User{}).Count(&total).Error; err != nil {
		// 记录数据库错误日志并返回统一错误
		logx.Errorf("获取用户总数失败, error: %v", err)
		return nil, 0, ginx.WithStack(status_error.DataOperationFailed)
	}

	// 分页查询
	if err := uc.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		// 记录数据库错误日志并返回统一错误
		logx.Errorf("获取用户列表失败, offset: %d, limit: %d, error: %v", offset, limit, err)
		return nil, 0, ginx.WithStack(status_error.DataOperationFailed)
	}

	return users, total, nil
}
