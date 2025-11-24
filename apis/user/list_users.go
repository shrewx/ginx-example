package user

import (
	"ginx-example/global"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shrewx/ginx"

	"ginx-example/controller"
	"ginx-example/models"
)

func init() {
	Router.Register(&ListUsers{})
}

// ListUsers 获取用户列表接口
type ListUsers struct {
	ginx.MethodGet
	Page     string `in:"query" name:"page" validate:"omitempty,min=1"`              // 页码，默认1
	PageSize string `in:"query" name:"page_size" validate:"omitempty,min=1,max=100"` // 每页数量，默认10，最大100
}

// ListUsersResponse 用户列表响应结构体
type ListUsersResponse struct {
	Users    []models.User `json:"users"`     // 用户列表
	Total    int64         `json:"total"`     // 总用户数
	Page     int           `json:"page"`      // 当前页码
	PageSize int           `json:"page_size"` // 每页数量
}

// Path 返回接口路径
func (l *ListUsers) Path() string {
	return ""
}

// Output 实现接口逻辑
func (l *ListUsers) Output(ctx *gin.Context) (interface{}, error) {
	// 解析分页参数
	page := 1
	pageSize := 10

	if l.Page != "" {
		if p, err := strconv.Atoi(l.Page); err == nil && p > 0 {
			page = p
		}
	}

	if l.PageSize != "" {
		if ps, err := strconv.Atoi(l.PageSize); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 创建用户控制器
	userController := controller.NewUserController(global.DB)

	// 获取用户列表
	users, total, err := userController.ListUsers(offset, pageSize)
	if err != nil {
		return nil, err
	}

	// 返回用户列表
	return ListUsersResponse{
		Users:    users,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
