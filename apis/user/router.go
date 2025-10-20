package user

import "github.com/shrewx/ginx"

// Router 用户模块路由分组
var Router = ginx.NewRouter(ginx.Group("users"))

func init() {
	// 注册用户相关接口
	Router.Register(&GetUserInfo{})  // GET /users/:id
	Router.Register(&CreateUser{})   // POST /users
	Router.Register(&ListUsers{})    // GET /users
}