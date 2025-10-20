---
inclusion: always
---

# Ginx 框架开发规则

## 接口开发约定

### 接口定义规则
- **必须实现 HandleOperator 接口**：所有 API 接口必须实现 `Path()`, `Method()`, `Output()` 三个方法
- **一个接口一个文件**：每个接口单独一个 Go 文件，文件名与结构体名称对应（如 `GetUserInfo` → `get_user_info.go`）
- **结构体命名**：使用动词+名词的格式，如 `GetUserInfo`, `CreateUser`, `UpdateUserProfile`
- **包组织**：按业务领域组织包，如 `apis/user/`, `apis/file/`, `apis/order/`
- **路径参数格式**：使用冒号语法，如 `/users/:id`，而不是 `/users/{id}`
- **接口路径**：返回相对于业务分组的路径，如 `/:id` 而不是 `/users/:id`

### 参数绑定规则
- **使用 tag 声明参数类型**：
  - `in:"query"` - URL 查询参数
  - `in:"path"` - URL 路径参数  
  - `in:"header"` - HTTP 头部参数
  - `in:"form"` - 表单数据
  - `in:"multipart"` - 文件上传
  - `in:"urlencoded"` - URL 编码数据
  - `in:"body"` - JSON 请求体

- **参数校验**：使用 `validate` tag 进行参数校验，字段都需要有中文注释
  ```go
  // 用户名称
  Username string `in:"query" validate:"required,min=3,max=20"`
  // 邮箱
  Email    string `in:"body" json:"email" validate:"required,email"`
  ```

### 响应处理规则
- **默认 JSON 响应**：直接返回结构体，框架自动序列化为 JSON
- **文件下载**：使用 `ginx.NewAttachment(filename, contentType)`
  ```go
  file := ginx.NewAttachment("data.txt", ginx.MineApplicationOctetStream)
  file.Write([]byte("hello world"))
  return file, nil
  ```
- **HTML 响应**：使用 `ginx.NewHTML()`
  ```go
  html := ginx.NewHTML()
  html.Write([]byte("<h1>Hello</h1>"))
  return html, nil
  ```
- **图片响应**：使用 `ginx.NewImagePNG()`, `ginx.NewImageJPEG()` 等
- **自定义响应**：实现 `MineDescriber` 接口或直接使用 `ctx.Data()`

## 错误处理约定

### 错误码定义规则
- **文件位置**：统一放在 `constants/status_error/error.go`
- **错误码格式**：`HTTP状态码 * 1e8 + 递增ID`
  ```go
  BadRequest StatusError = http.StatusBadRequest*1e8 + iota + 1  // 40000000001
  UserNotFound StatusError = http.StatusNotFound*1e8 + iota + 1  // 40400000001
  ```
- **错误码定义位置**：新的错误码必须定义在相同HTTP状态码的现有错误码后面，不存在则自己定义。
  ```go
  BadRequest StatusError = http.StatusBadRequest*1e8 + iota + 1  // 
  UserIDInvalid


  NotFound StatusError = http.StatusNotFound*1e8 + iota + 1  // 
  UserNotFound
  ```


### 错误注释规则
- **必须包含中英文注释**：
  ```go
  // @errZH 用户未找到，ID：{{.UserID}}
  // @errEN user not found, ID: {{.UserID}}
  UserNotFound StatusError = http.StatusNotFound*1e8 + iota + 1
  ```

### 错误使用规则
- **简单错误**：直接返回错误码
  ```go
  return nil, status_error.UserNotFound
  ```
- **带参数错误**：使用 `WithParams()` 注入参数
  ```go
  return nil, status_error.UserNotFound.WithParams(map[string]interface{}{
      "UserID": userID,
  })
  ```
- **数据库错误**：Controller中捕获的数据库错误统一使用 `DataOperationFailed`，并使用 `logx` 记录错误日志
  ```go
  if result.Error != nil {
      logx.Errorf("获取用户信息失败, userID: %d, error: %v", userID, result.Error)
      return nil, status_error.DataOperationFailed
  }
  ```

## 中间件开发约定

### 中间件定义规则
- **实现 TypeOperator 接口**：包含 `Output()` 和 `Type()` 方法
- **文件位置**：统一放在 `middleware/` 目录
- **命名规范**：使用功能描述命名，如 `JWTAuth`, `RateLimit`, `CORS`

### 中间件使用规则
- **全局中间件**：在路由组创建时添加
  ```go
  AuthRouter = ginx.NewRouter(ginx.Group("api/v1"), &middleware.JWTAuth{})
  ```
- **参数获取**：通过结构体字段和 tag 获取请求参数
- **上下文传递**：使用 `ctx.Set()` 传递数据给后续处理器

## 路由组织约定

### 路由组规则
- **按版本分组**：`api/v1`, `api/v2`
- **按权限分组**：`public`, `auth`, `admin`
- **按业务分组**：`user`, `order`, `payment`, `file`

### 路由注册架构
采用三层路由注册架构，确保代码组织清晰和模块化：

```
apis/
├── root.go              # 版本分组注册（第一层）
├── user/
│   ├── router.go        # 用户业务分组（第二层）
│   ├── get_user_info.go # 具体接口实现（第三层）
│   ├── create_user.go
│   └── list_users.go
├── file/
│   ├── router.go        # 文件业务分组（第二层）
│   ├── upload_file.go   # 具体接口实现（第三层）
│   └── download_file.go
└── order/
    ├── router.go        # 订单业务分组（第二层）
    └── ...
```

### 路由注册规则
- **三层路由架构**：采用分层路由注册结构确保代码组织清晰
  1. **第一层 - 版本分组**：`apis/root.go` 管理API版本和注册业务分组
  2. **第二层 - 业务分组**：`apis/{module}/router.go` 定义业务模块路由分组
  3. **第三层 - 接口实现**：`apis/{module}/{interface}.go` 具体接口实现

- **业务分组定义规则**：
  ```go
  // apis/user/router.go
  package user
  
  import "github.com/shrewx/ginx"
  
  // Router 用户模块路由分组
  var Router = ginx.NewRouter(ginx.Group("users"))
  
  func init() {
      // 注册用户相关接口
      Router.Register(&GetUserInfo{})  // GET /api/v1/users/:id
      Router.Register(&CreateUser{})   // POST /api/v1/users
      Router.Register(&ListUsers{})    // GET /api/v1/users
  }
  ```

- **版本分组注册规则**：
  ```go
  // apis/root.go
  var V1Router = ginx.NewRouter(ginx.Group("api/v1"))
  
  func init() {
      V1Router.Register(user.Router)   // 注册用户模块路由
      V1Router.Register(file.Router)   // 注册文件模块路由
  }
  ```

- **路径规则**：
  - 最终路径：`/api/v1/{group}/{interface_path}`
  - 接口路径：相对于业务分组的路径
  - 路径参数：使用冒号语法 `/users/:id`

- **命名规范**：
  - 业务分组：使用复数形式 `users`, `files`, `orders`
  - 路由变量：统一命名为 `Router`
  - 分组文件：统一命名为 `router.go`


## 代码生成约定

### go:generate 指令
- **错误码生成**：
  ```go
  //go:generate toolx gen error -p error_codes -c StatusError
  //go:generate toolx gen errorYaml -p error_codes -o ../i18n -c StatusError
  ```
- **i18n 生成**：
  ```go
  //go:generate toolx gen i18n prefix errors.references CommonField
  //go:generate toolx gen i18nYaml -p errors.references -o ../i18n -c CommonField
  ```
- **OpenAPI 生成**：
  ```go
  //go:generate toolx gen openapi
  ```

### 自动化规范
- **参数绑定**：通过 tag 自动解析和校验请求参数
- **响应处理**：支持多种 MIME 类型 (JSON、HTML、文件下载、图片等)
- **错误处理**：结构化错误码支持国际化和参数注入
- **文档生成**：自动生成 OpenAPI 规范、错误处理和 i18n 文件

### 生成文件管理
- **不要修改生成文件**：所有 `*__generated.go` 文件都是自动生成的，不要手动修改
- **版本控制**：生成的文件应该提交到版本控制系统
- **重新生成**：修改源码后需要重新运行 `go generate`

## 国际化约定

### i18n 字段定义
- **文件位置**：`constants/i18n/fields.go`
- **注释格式**：
  ```go
  // @i18nZH 用户名
  // @i18nEN username
  Username Field = "username"
  ```

### i18n 文件组织
- **目录结构**：`constants/i18n/`
- **文件命名**：`zh_*.yaml`, `en_*.yaml`
- **键值结构**：使用点号分隔的层级结构

## 数据库模型约定

### 模型定义规则
- **表注册**：新表的创建都需要在模型文件中添加表注册
  ```go
  func init() {
      dbhelper.RegisterTable(&User{})
  }
  ```
- **参数注释**：定义的表数据，字段参数需要有对应的字段中文注释
  ```go
  // 用户ID
  ID string `in:"path" name:"id" validate:"required"`
  // 用户信息
  User models.User `json:"user"` 
  ```

## 项目结构约定

```
ginx-example/
├── apis/                    # API 接口定义
│   ├── user/               # 用户相关接口
│   ├── file/               # 文件相关接口
│   └── root.go             # 路由注册
├── constants/              # 常量定义
│   ├── status_error/       # 错误码定义
│   └── i18n/              # 国际化文件
├── global/                 # 全局配置
├── middleware/             # 中间件
├── models/                 # 数据模型
├── controller/             # 数据控制层
├── internal/               # 内部方法定义
├── cmd/                   # 命令行程序
```

## 开发流程约定

### 新接口开发流程
1. **定义接口结构体**：实现 `HandleOperator` 接口
2. **定义请求/响应结构体**：包含必要的 tag 和校验规则
3. **实现业务逻辑**：在 `Output()` 方法中实现
4. **定义错误码**：如果需要新的错误类型
5. **注册路由**：在对应的路由组中注册
6. **生成代码**：运行 `go generate` 生成相关文件
7. **测试接口**：编写单元测试和集成测试

### 错误处理流程
1. **定义错误码**：在 `status_error/error.go` 中定义
2. **添加注释**：包含中英文错误描述
3. **生成文件**：运行 `go generate` 生成错误码实现和 i18n 文件
4. **使用错误码**：在接口中返回对应的错误码

### 代码审查要点
- **接口定义**：检查是否正确实现了 `HandleOperator` 接口
- **参数校验**：确保所有必要参数都有校验规则
- **错误处理**：检查错误码使用是否正确，错误信息是否清晰
- **代码生成**：确保生成的文件已更新并提交
- **测试覆盖**：检查是否有足够的测试覆盖

## 性能优化约定

### 参数绑定优化
- **避免过度校验**：只对必要字段进行校验
- **合理使用缓存**：对频繁访问的数据进行缓存
- **批量操作**：避免在循环中进行数据库操作

### 错误处理优化
- **错误码复用**：相同类型的错误使用同一个错误码
- **参数注入优化**：避免创建过多临时对象
- **国际化缓存**：缓存已翻译的错误消息

## 安全约定

### 参数安全
- **输入校验**：所有外部输入都必须进行校验
- **SQL 注入防护**：使用参数化查询
- **XSS 防护**：对输出内容进行转义

### 认证授权
- **JWT 验证**：使用标准的 JWT 中间件
- **权限检查**：在需要的接口上添加权限中间件
- **敏感信息**：不在日志中记录敏感信息

## 监控和日志约定

### 日志记录
- **结构化日志**：使用 JSON 格式的结构化日志
- **日志级别**：合理使用 DEBUG, INFO, WARN, ERROR 级别
- **关键信息**：记录请求 ID、用户 ID、操作类型等关键信息

### 性能监控
- **响应时间**：监控接口响应时间
- **错误率**：监控接口错误率
- **资源使用**：监控 CPU、内存、数据库连接等资源使用情况