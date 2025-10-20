# Ginx框架编码规则补充

## 新增编码规则

基于项目实践，我们补充了以下编码规则，以确保代码的一致性和可维护性。

### 1. API参数注释规范

#### 规则说明
所有API的请求参数和响应参数都需要有对应的注释说明。

#### 实施要求
- **请求参数**：每个参数字段后必须添加注释，说明参数的用途、格式要求等
- **响应参数**：每个响应字段后必须添加注释，说明字段的含义

#### 代码示例

**请求参数注释**
```go
type GetUserInfo struct {
    ginx.MethodGet
    ID string `in:"path" name:"id" validate:"required"` // 用户ID，必须为正整数
}

type CreateUserRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50"`        // 用户名，3-50字符，必填，唯一
    Email    string `json:"email" validate:"required,email"`                  // 邮箱地址，必填，唯一，格式验证
    Name     string `json:"name" validate:"required,min=1,max=100"`           // 真实姓名，1-100字符，必填
    Status   string `json:"status" validate:"omitempty,oneof=active inactive"` // 用户状态，可选值：active/inactive，默认active
}
```

**响应参数注释**
```go
type GetUserInfoResponse struct {
    User models.User `json:"user"` // 用户信息
}

type ListUsersResponse struct {
    Users    []models.User `json:"users"`     // 用户列表
    Total    int64         `json:"total"`     // 总用户数
    Page     int           `json:"page"`      // 当前页码
    PageSize int           `json:"page_size"` // 每页数量
}
```

### 2. 数据库表注册规范

#### 规则说明
新表的创建都需要在模型文件中添加表注册代码。

#### 实施要求
- 在模型文件中添加 `init()` 函数
- 使用 `dbhelper.RegisterTable()` 注册表模型

#### 代码示例
```go
package models

import (
    "github.com/shrewx/ginx/pkg/dbhelper"
    "gorm.io/gorm"
    "time"
)

func init() {
    dbhelper.RegisterTable(&User{})
}

type User struct {
    ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
    Username  string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
    // ... 其他字段
}
```

### 3. Controller错误处理规范

#### 规则说明
Controller中如果捕获的数据库错误，统一使用 `DataOperationFailed` 错误码，并使用 `logx` 打印对应的错误在返回之前。

#### 实施要求
- 使用 `errors.Is()` 判断特定错误类型
- 数据库操作错误统一返回 `DataOperationFailed`
- 使用 `logx.Errorf()` 记录详细错误信息
- 业务逻辑错误返回具体的业务错误码

#### 代码示例

**正确的错误处理**
```go
func (uc *UserController) GetUserByID(userID int64) (*models.User, error) {
    var user models.User
    
    result := uc.db.First(&user, userID)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            // 业务逻辑错误：用户不存在
            return nil, status_error.UserNotFound.WithParams(map[string]interface{}{
                "UserID": userID,
            })
        }
        // 数据库操作错误：记录日志并返回统一错误
        logx.Errorf("获取用户信息失败, userID: %d, error: %v", userID, result.Error)
        return nil, status_error.DataOperationFailed
    }
    
    return &user, nil
}
```

**错误码定义**
```go
const (
    // @errZH 数据操作失败，请稍后重试
    // @errEN data operation failed, please try again later
    DataOperationFailed StatusError = http.StatusInternalServerError*1e8 + iota + 1
)
```

## 规范更新位置

这些编码规则已经更新到以下文档中：

1. **项目规范** (`openspec/project.md`)
   - 添加了API参数注释要求
   - 添加了表注册规范
   - 添加了Controller错误处理规范

2. **Ginx框架规范** (`openspec/specs/ginx-framework/spec.md`)
   - 新增API参数注释规范章节
   - 新增数据库表注册规范章节
   - 新增Controller错误处理规范章节

3. **开发规则** (`.kiro/steering/ginx-development-rules.md`)
   - 补充了数据库模型约定
   - 更新了错误使用规则
   - 添加了参数注释要求

## 实施检查清单

在代码审查时，请确保以下项目都已满足：

### API接口检查
- [ ] 所有请求参数都有详细注释
- [ ] 所有响应参数都有详细注释
- [ ] 参数校验规则正确设置
- [ ] 错误处理符合规范

### 数据模型检查
- [ ] 新模型已添加表注册代码
- [ ] GORM标签设置正确
- [ ] 字段注释完整

### Controller检查
- [ ] 数据库错误统一使用 `DataOperationFailed`
- [ ] 错误日志记录完整
- [ ] 业务错误使用具体错误码
- [ ] 使用 `errors.Is()` 判断错误类型

### 文档检查
- [ ] API文档更新
- [ ] 错误码文档更新
- [ ] 使用示例完整

## 工具支持

建议使用以下工具来确保规范的执行：

1. **代码生成工具**
   ```bash
   go generate ./...
   ```

2. **代码检查工具**
   ```bash
   golangci-lint run
   ```

3. **测试工具**
   ```bash
   go test ./...
   ```

通过这些编码规则的实施，我们可以确保项目代码的一致性、可读性和可维护性。