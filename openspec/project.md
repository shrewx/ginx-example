# 项目上下文

## 项目目的
这是一个基于 ginx 框架构建的 Go Web API 示例项目。它展示了构建 REST API 的结构化方法，包含适当的错误处理、国际化和数据库集成。

## 技术栈
- **编程语言**: Go 1.24.3
- **Web 框架**: Ginx (基于 Gin 的封装框架)
  - 参考 httptransport 设计思想
  - 支持自动代码生成 (OpenAPI、错误码、i18n、客户端 SDK)
  - 规范化接口定义和参数处理
  - 三层路由架构：版本分组 → 业务分组 → 接口实现
  - 统一的 HandleOperator 接口约定
- **数据库**: GORM，支持 MySQL 和 PostgreSQL
- **命令行工具**: Cobra 命令行接口
- **国际化**: go-i18n/v2 多语言支持
- **配置管理**: cleanenv 基于环境的配置
- **日志记录**: logrus 配合 lumberjack 日志轮转
- **可观测性**: OpenTelemetry 配合 Jaeger 和 Zipkin 导出器
- **服务发现**: Consul 集成
- **代码生成**: toolx 工具集 (错误码、OpenAPI、i18n)

## 项目约定

### 代码风格
- 遵循标准 Go 约定 (gofmt, golint)
- 使用有意义的包名并按领域组织
- 错误常量使用 HTTP 状态码乘以 1e8 加上递增 ID
- 使用 `go:generate` 指令配合 toolx 生成代码
- 接口定义：一个接口一个文件，文件名与类名相同
- 所有路由必须实现 `HandleOperator` 接口
- 中间件必须实现 `TypeOperator` 接口
- 使用 tag 声明参数类型：`in:"query|path|header|form|multipart|urlencoded|body"`
- **路径参数格式**：使用冒号语法，如 `/users/:id`，而不是 `/users/{id}`
- **错误码定义位置**：新的错误码必须定义在相同HTTP状态码的现有错误码后面，不存在则自己定义。
- **API参数注释**：所有API的请求参数字段和响应参数字读都需要有对应的注释
- **表注册规范**：新表的创建都需要在模型文件中添加 `func init() { dbhelper.RegisterTable(&ModelName{}) }` 注册
- **Controller错误处理**：如果捕获的数据库错误，统一使用 `DataOperationFailed` 错误码，并使用 `logx` 打印对应的错误在返回之前
- 错误码注释使用 `@errZH` 和 `@errEN`
- 国际化注释使用 `@i18nZH` 和 `@i18nEN`

### 架构模式
- **分层架构**: APIs → Controller → Models → 数据库
- **路由模式**: 基于 Ginx 的三层路由组织架构
  - 第一层：版本分组 (`apis/root.go`)
  - 第二层：业务分组 (`apis/{module}/router.go`)
  - 第三层：接口实现 (`apis/{module}/{interface}.go`)
- **接口模式**: 实现 `HandleOperator` 接口的标准化接口定义
  - 必须实现 `Path()`, `Method()`, `Output()` 三个方法
  - 支持路径参数、查询参数、请求体等多种参数类型
- **中间件模式**: 实现 `TypeOperator` 接口的中间件系统
- **参数绑定**: 通过 tag 自动解析和校验请求参数
  - 支持 `in:"query|path|header|form|multipart|urlencoded|body"`
  - 集成 validator 库进行参数校验
- **响应处理**: 支持多种 MIME 类型 (JSON、HTML、文件下载、图片等)
- **配置模式**: 全局配置结构体支持 YAML
- **错误处理**: 结构化错误码支持国际化和参数注入
  - 错误码格式：`HTTP状态码 * 1e8 + 递增ID`
  - 支持错误参数注入和多语言消息
- **代码生成**: 使用 toolx 生成 OpenAPI 规范、错误处理和 i18n 文件
  - 通过注释驱动生成，保持代码和文档同步
  - 支持多种输出格式

### 测试策略
- 使用 Go 内置测试框架
- 专注于业务逻辑的单元测试
- 数据库操作的集成测试

### Git 工作流
- 使用约定式提交
- 新功能开发使用特性分支
- 提交前生成代码 (OpenAPI 规范、错误翻译)

## 领域上下文
这是一个基于 Ginx 框架的 Web API 示例项目，展示了：

### 核心特性
- **接口约定**: 通过 `HandleOperator` 接口统一API定义规范
- **三层路由**: 版本分组 → 业务分组 → 接口实现的清晰架构
- **自动代码生成**: OpenAPI 文档、错误码、i18n 文件、客户端 SDK
- **参数绑定**: 通过 tag 声明参数类型，支持自动解析和校验
- **多种响应**: JSON、HTML、文件下载、图片等多种响应类型
- **结构化错误**: 统一的错误码格式和国际化支持
- **中间件系统**: 基于 `TypeOperator` 接口的标准化中间件

### 开发体验
- **规范化开发**: 统一的接口定义和参数处理方式
- **自动化工具**: 减少手动维护文档和代码的工作量
- **类型安全**: 利用 Go 语言类型系统进行编译时检查
- **IDE 友好**: 良好的代码补全和错误提示支持
- **测试友好**: 清晰的模块划分便于单元测试和集成测试

### 技术集成
- **GORM 数据库**: 完整的 ORM 支持和数据库迁移
- **配置管理**: 基于 YAML 的配置文件和环境变量支持
- **日志系统**: 结构化日志记录和错误追踪
- **可观测性**: OpenTelemetry 集成，支持链路追踪和监控
- **国际化**: 完整的多语言支持 (中文/英文)

## 重要约束
- **接口约定**: 所有API接口必须实现 `HandleOperator` 接口
- **路由组织**: 采用三层路由架构，确保代码组织清晰
- **错误码规范**: 遵循 `HTTP状态码 * 1e8 + 递增ID` 模式
- **参数注释**: 所有API参数必须有详细的中文注释说明
- **表注册**: 新数据表必须使用 `dbhelper.RegisterTable()` 注册
- **错误处理**: Controller层数据库错误统一使用 `DataOperationFailed` 并记录日志
- **路径格式**: 使用冒号语法 `/users/:id` 而不是 `/users/{id}`
- **代码生成**: 通过 `go:generate` 指令自动生成相关文件
- **数据库迁移**: 在启动时自动运行
- **配置管理**: 通过 YAML 文件提供配置

## 外部依赖
- **数据库**: MySQL 或 PostgreSQL (可配置)
- **服务发现**: Consul (可选)
- **监控**: Jaeger 或 Zipkin 分布式追踪
- **日志**: 基于文件的日志轮转
