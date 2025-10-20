# 用户API规范

## 概述

用户API提供用户信息的查询和管理功能，展示Ginx框架的标准接口实现模式。

## 核心功能

### Requirement: 用户信息查询接口
系统必须提供根据用户ID查询用户信息的接口。

#### Scenario: 根据有效用户ID查询用户信息
- **GIVEN** 系统中存在用户ID为123的用户
- **WHEN** 客户端发送 `GET /api/v1/users/123` 请求
- **THEN** 系统返回HTTP状态码200和用户详细信息
- **AND** 响应包含用户ID、用户名、邮箱、姓名和状态字段

#### Scenario: 查询不存在的用户
- **GIVEN** 系统中不存在用户ID为999的用户
- **WHEN** 客户端发送 `GET /api/v1/users/999` 请求
- **THEN** 系统返回HTTP状态码404和用户不存在的错误信息
- **AND** 错误信息支持中英文国际化

#### Scenario: 使用无效的用户ID查询
- **GIVEN** 客户端提供无效的用户ID（如负数或非数字）
- **WHEN** 客户端发送 `GET /api/v1/users/-1` 或 `GET /api/v1/users/abc` 请求
- **THEN** 系统返回HTTP状态码400和参数错误信息
- **AND** 错误信息明确指出用户ID格式不正确

### Requirement: 用户数据模型定义
系统必须定义标准的用户数据结构。

#### Scenario: 用户信息结构定义
- **GIVEN** 系统需要返回用户信息
- **WHEN** 定义用户数据模型
- **THEN** 必须包含以下字段：
  - `id`: 用户唯一标识符（int64类型）
  - `username`: 用户名（string类型）
  - `email`: 邮箱地址（string类型）
  - `name`: 真实姓名（string类型）
  - `status`: 用户状态（string类型）
- **AND** 所有字段必须包含JSON序列化标签

```go
type User struct {
    ID       int64  `json:"id"`       // 用户唯一标识符
    Username string `json:"username"` // 用户名
    Email    string `json:"email"`    // 邮箱地址
    Name     string `json:"name"`     // 真实姓名
    Status   string `json:"status"`   // 用户状态
}
```

### Requirement: 参数校验规则
系统必须对输入参数进行严格校验。

#### Scenario: 用户ID参数校验
- **GIVEN** 接口接收用户ID作为路径参数
- **WHEN** 定义参数校验规则
- **THEN** 用户ID必须满足以下条件：
  - 必须为正整数
  - 最小值为1
  - 使用 `validate:"required,min=1"` 标签
- **AND** 参数类型声明为 `in:"path"`

### Requirement: 错误码定义
系统必须定义用户查询相关的错误码。

#### Scenario: 用户不存在错误码
- **GIVEN** 查询的用户不存在
- **WHEN** 定义错误码
- **THEN** 错误码必须遵循 `HTTP状态码 * 1e8 + 递增ID` 格式
- **AND** 必须包含中英文错误描述注释
- **AND** 支持参数注入显示具体的用户ID

#### Scenario: 用户ID无效错误码
- **GIVEN** 提供的用户ID格式无效
- **WHEN** 定义错误码
- **THEN** 错误码必须基于400状态码
- **AND** 错误描述必须明确指出参数格式问题

```go
const (
    // @errZH 请求参数错误
    // @errEN bad request
    BadRequest StatusError = http.StatusBadRequest*1e8 + iota + 1
    // @errZH 用户ID无效，必须为正整数
    // @errEN invalid user ID, must be a positive integer
    InvalidUserID StatusError = http.StatusBadRequest*1e8 + iota + 1
)

const (
    // @errZH 资源未找到
    // @errEN not found
    NotFound StatusError = http.StatusNotFound*1e8 + iota + 1
    // @errZH 用户不存在，用户ID：{{.UserID}}
    // @errEN user not found, user ID: {{.UserID}}
    UserNotFound StatusError = http.StatusNotFound*1e8 + iota + 1
)
```

### Requirement: 接口实现规范
接口实现必须遵循Ginx框架规范。

#### Scenario: HandleOperator接口实现
- **GIVEN** 需要实现用户查询接口
- **WHEN** 定义接口结构体
- **THEN** 必须实现以下方法：
  - `Path()`: 返回 "/users/:id"
  - `Method()`: 继承自 `ginx.MethodGet`
  - `Output()`: 实现查询逻辑和响应处理
- **AND** 结构体必须包含用户ID字段和相应的标签

```go
type GetUserInfo struct {
    ginx.MethodGet
    ID string `in:"path" name:"id" validate:"required"` // 用户ID
}

func (g *GetUserInfo) Path() string {
    return "/users/:id"
}

func (g *GetUserInfo) Output(ctx *gin.Context) (interface{}, error) {
    // 实现查询逻辑
}
```

#### Scenario: 路由注册
- **GIVEN** 接口实现完成
- **WHEN** 注册路由
- **THEN** 必须在V1Router中注册该接口
- **AND** 接口路径必须符合RESTful规范

```go
func init() {
    V1Router.Register(&user.GetUserInfo{})
}
```

### Requirement: 模拟数据提供
系统必须提供模拟用户数据用于演示。

#### Scenario: 内存数据存储
- **GIVEN** 这是一个示例项目
- **WHEN** 实现数据查询逻辑
- **THEN** 使用内存中的模拟数据
- **AND** 至少包含3个不同的用户记录
- **AND** 数据结构必须与用户模型保持一致

#### Scenario: 数据查询逻辑
- **GIVEN** 接收到用户ID查询请求
- **WHEN** 执行查询逻辑
- **THEN** 在模拟数据中查找对应的用户
- **AND** 找到用户时返回完整用户信息
- **AND** 未找到用户时返回相应错误码

## 实现示例

### 完整接口实现
```go
package user

import (
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/shrewx/ginx"
    "ginx-example/constants/status_error"
    "ginx-example/models"
)

type GetUserInfo struct {
    ginx.MethodGet
    ID string `in:"path" name:"id" validate:"required"`
}

func (g *GetUserInfo) Path() string {
    return "/users/:id"
}

func (g *GetUserInfo) Output(ctx *gin.Context) (interface{}, error) {
    userID, err := strconv.ParseInt(g.ID, 10, 64)
    if err != nil || userID <= 0 {
        return nil, status_error.InvalidUserID
    }

    user, found := models.GetUserByID(userID)
    if !found {
        return nil, status_error.UserNotFound.WithParams(map[string]interface{}{
            "UserID": userID,
        })
    }

    return models.GetUserInfoResponse{User: *user}, nil
}
```

### 响应格式
```json
{
  "user": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "name": "管理员",
    "status": "active"
  }
}
```