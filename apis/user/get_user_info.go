package user

import (
	"ginx-example/global"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shrewx/ginx"

	"ginx-example/constants/status_error"
	"ginx-example/controller"
	"ginx-example/models"
)

func init() {
	Router.Register(&GetUserInfo{})
}

// GetUserInfo 获取用户信息接口
type GetUserInfo struct {
	ginx.MethodGet
	ID string `in:"path" name:"id" validate:"required"` // 用户ID
}

// GetUserInfoResponse 获取用户信息响应结构体
type GetUserInfoResponse struct {
	User models.User `json:"user"` // 用户信息
}

// Path 返回接口路径
func (g *GetUserInfo) Path() string {
	return "/users/:id"
}

// Output 实现接口逻辑
func (g *GetUserInfo) Output(ctx *gin.Context) (interface{}, error) {
	// 将字符串ID转换为int64
	userID, err := strconv.ParseInt(g.ID, 10, 64)
	if err != nil || userID <= 0 {
		return nil, status_error.InvalidUserID
	}

	// 创建用户控制器
	userController := controller.NewUserController(global.DB)

	// 查询用户信息
	user, err := userController.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// 返回用户信息
	return GetUserInfoResponse{
		User: *user,
	}, nil
}
