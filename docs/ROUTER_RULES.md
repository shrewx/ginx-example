# 路由注册规则详解

## 概述

Ginx框架采用三层路由注册架构，确保代码组织清晰、模块化程度高，便于维护和扩展。

## 三层路由架构

### 第一层：版本分组注册
**文件位置**: `apis/root.go`
**职责**: 注册各个业务模块的路由分组，管理API版本

```go
package apis

import (
    "ginx-example/apis/user"
    "ginx-example/apis/file"
    "ginx-example/apis/order"
    "github.com/shrewx/ginx"
)

// V1Router API v1版本路由组
var V1Router = ginx.NewRouter(ginx.Group("api/v1"))

func init() {
    // 注册各个业务模块的路由分组
    V1Router.Register(user.Router)  // 注册用户模块路由
    V1Router.Register(file.Router)  // 注册文件模块路由
    V1Router.Register(order.Router) // 注册订单模块路由
}
```

### 第二层：业务分组定义
**文件位置**: `apis/{module}/router.go`
**职责**: 定义业务模块的路由分组，注册该模块下的所有接口

```go
// apis/user/router.go
package user

import "github.com/shrewx/ginx"

// Router 用户模块路由分组
var Router = ginx.NewRouter(ginx.Group("users"))

func init() {
    // 注册用户相关接口
    Router.Register(&GetUserInfo{})  // GET /users/:id
    Router.Register(&CreateUser{})   // POST /users
    Router.Register(&ListUsers{})    // GET /users
    Router.Register(&UpdateUser{})   // PUT /users/:id
    Router.Register(&DeleteUser{})   // DELETE /users/:id
}
```

### 第三层：接口实现
**文件位置**: `apis/{module}/{interface}.go`
**职责**: 实现具体的API接口逻辑

```go
// apis/user/get_user_info.go
package user

import (
    "github.com/gin-gonic/gin"
    "github.com/shrewx/ginx"
    // ...
)

// GetUserInfo 获取用户信息接口
type GetUserInfo struct {
    ginx.MethodGet
    ID string `in:"path" name:"id" validate:"required"` // 用户ID，必须为正整数
}

func (g *GetUserInfo) Path() string {
    return "/:id"  // 相对于业务分组的路径
}

func (g *GetUserInfo) Output(ctx *gin.Context) (interface{}, error) {
    // 接口实现逻辑
}
```

## 路由注册规则

### 1. 业务分组命名规则
- **分组名称**: 使用复数形式的英文名词，如 `users`, `files`, `orders`
- **变量命名**: 使用 `Router` 作为分组变量名
- **文件命名**: 统一使用 `router.go` 作为分组定义文件名

### 2. 接口路径规则
- **绝对路径**: 最终路径为 `/api/v1/{group}/{interface_path}`
- **相对路径**: 接口的 `Path()` 方法返回相对于业务分组的路径
- **RESTful规范**: 遵循RESTful API设计规范

### 3. 注册时机规则
- **自动注册**: 通过 `init()` 函数自动注册，无需手动调用
- **注册顺序**: 先注册接口到业务分组，再注册业务分组到版本分组
- **延迟加载**: 支持按需加载业务模块

## 实际路径映射

### 用户模块示例
```
业务分组: "users"
版本分组: "api/v1"

接口定义                    相对路径        最终路径
GetUserInfo.Path() = "/:id"     →  /api/v1/users/:id
CreateUser.Path() = ""          →  /api/v1/users
ListUsers.Path() = ""           →  /api/v1/users
UpdateUser.Path() = "/:id"      →  /api/v1/users/:id
DeleteUser.Path() = "/:id"      →  /api/v1/users/:id
```

### 文件模块示例
```
业务分组: "files"
版本分组: "api/v1"

接口定义                        相对路径            最终路径
UploadFile.Path() = ""              →  /api/v1/files
DownloadFile.Path() = "/:id"        →  /api/v1/files/:id
DeleteFile.Path() = "/:id"          →  /api/v1/files/:id
GetFileList.Path() = ""             →  /api/v1/files
```

## 权限分组扩展

对于需要权限控制的场景，可以创建多个路由分组：

```go
// apis/user/router.go
package user

import "github.com/shrewx/ginx"

// PublicRouter 公开访问的用户接口
var PublicRouter = ginx.NewRouter(ginx.Group("users"))

// AuthRouter 需要认证的用户接口
var AuthRouter = ginx.NewRouter(ginx.Group("users"))

// AdminRouter 需要管理员权限的用户接口
var AdminRouter = ginx.NewRouter(ginx.Group("users"))

func init() {
    // 公开接口
    PublicRouter.Register(&CreateUser{})    // 用户注册
    
    // 认证接口
    AuthRouter.Register(&GetUserInfo{})     // 获取用户信息
    AuthRouter.Register(&UpdateUser{})      // 更新用户信息
    
    // 管理员接口
    AdminRouter.Register(&ListUsers{})      // 用户列表
    AdminRouter.Register(&DeleteUser{})     // 删除用户
}
```

然后在 `apis/root.go` 中分别注册：

```go
func init() {
    // 公开路由
    V1Router.Register(user.PublicRouter)
    
    // 认证路由（需要添加认证中间件）
    AuthV1Router := ginx.NewRouter(ginx.Group("api/v1"), &middleware.JWTAuth{})
    AuthV1Router.Register(user.AuthRouter)
    
    // 管理员路由（需要添加管理员权限中间件）
    AdminV1Router := ginx.NewRouter(ginx.Group("api/v1"), &middleware.AdminAuth{})
    AdminV1Router.Register(user.AdminRouter)
}
```

## 最佳实践

### 1. 模块化开发
- 每个业务模块独立开发和维护
- 接口变更只影响对应的业务分组
- 便于团队协作和代码审查

### 2. 版本管理
- 通过版本分组管理API版本
- 支持多版本并存
- 便于API升级和兼容性管理

### 3. 权限控制
- 通过不同的路由分组实现权限控制
- 中间件可以精确控制到业务模块级别
- 支持细粒度的权限管理

### 4. 测试友好
- 每个业务分组可以独立测试
- 便于编写集成测试
- 支持模块级别的性能测试

## 注意事项

1. **循环依赖**: 避免业务模块之间的循环依赖
2. **命名冲突**: 确保不同模块的接口名称不冲突
3. **路径冲突**: 注意相同HTTP方法和路径的冲突
4. **中间件顺序**: 注意中间件的执行顺序
5. **性能考虑**: 避免在 `init()` 函数中执行耗时操作

通过遵循这些路由注册规则，可以构建出结构清晰、易于维护的API服务。