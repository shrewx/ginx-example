package apis

import (
	"ginx-example/apis/file"
	"ginx-example/apis/user"
	"github.com/shrewx/ginx"
)

// V1Router API v1版本路由组
var V1Router = ginx.NewRouter(ginx.Group("api/v1"))

func init() {
	// 注册各个业务模块的路由分组
	V1Router.Register(user.Router) // 注册用户模块路由
	V1Router.Register(file.Router) // 注册文件模块路由
}
