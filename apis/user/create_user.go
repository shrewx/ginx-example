package user

import (
	"ginx-example/global"
	"github.com/gin-gonic/gin"
	"github.com/shrewx/ginx"

	"ginx-example/controller"
	"ginx-example/models"
)

func init() {
	Router.Register(&CreateUser{})
}

// CreateUser 创建用户接口
type CreateUser struct {
	ginx.MethodPost
	Body CreateUserRequest `in:"body" json:"body"` // 创建用户请求体
}

// CreateUserRequest 创建用户请求结构体
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`         // 用户名，3-50字符，必填，唯一
	Email    string `json:"email" validate:"required,email"`                   // 邮箱地址，必填，唯一，格式验证
	Name     string `json:"name" validate:"required,min=1,max=100"`            // 真实姓名，1-100字符，必填
	Status   string `json:"status" validate:"omitempty,oneof=active inactive"` // 用户状态，可选值：active/inactive，默认active
}

// CreateUserResponse 创建用户响应结构体
type CreateUserResponse struct {
	User models.User `json:"user"` // 创建的用户信息
}

// Path 返回接口路径
func (c *CreateUser) Path() string {
	return "/users"
}

// Output 实现接口逻辑
func (c *CreateUser) Output(ctx *gin.Context) (interface{}, error) {
	// 创建用户控制器
	userController := controller.NewUserController(global.DB)

	// 构建用户对象
	user := &models.User{
		Username: c.Body.Username,
		Email:    c.Body.Email,
		Name:     c.Body.Name,
		Status:   c.Body.Status,
	}

	// 如果状态为空，设置默认值
	if user.Status == "" {
		user.Status = "active"
	}

	// 创建用户
	if err := userController.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// 返回创建的用户信息
	return CreateUserResponse{
		User: *user,
	}, nil
}
