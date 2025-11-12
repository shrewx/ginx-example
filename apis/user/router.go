package user

import "github.com/shrewx/ginx"

// Router 用户模块路由分组
var Router = ginx.NewRouter(ginx.Group("users"))
