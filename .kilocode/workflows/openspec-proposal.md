<!-- OPENSPEC:START -->
**Guardrails**
- Favor straightforward, minimal implementations first and add complexity only when it is requested or clearly required.
- Keep changes tightly scoped to the requested outcome.
- Refer to `openspec/AGENTS.md` (located inside the `openspec/` directory—run `ls openspec` or `openspec update` if you don't see it) if you need additional OpenSpec conventions or clarifications.
- Identify any vague or ambiguous details and ask the necessary follow-up questions before editing files.

**Steps**
1. Review `openspec/project.md`, run `openspec list` and `openspec list --specs`, and inspect related code or docs (e.g., via `rg`/`ls`) to ground the proposal in current behaviour; note any gaps that require clarification.
2. Choose a unique verb-led `change-id` and scaffold `proposal.md`, `tasks.md`, and `design.md` (when needed) under `openspec/changes/<id>/`.
3. Map the change into concrete capabilities or requirements, breaking multi-scope efforts into distinct spec deltas with clear relationships and sequencing.
4. Capture architectural reasoning in `design.md` when the solution spans multiple systems, introduces new patterns, or demands trade-off discussion before committing to specs.
5. Draft spec deltas in `changes/<id>/specs/<capability>/spec.md` (one folder per capability) using `## ADDED|MODIFIED|REMOVED Requirements` with at least one `#### Scenario:` per requirement and cross-reference related capabilities when relevant.
6. Draft `tasks.md` as an ordered list of small, verifiable work items that deliver user-visible progress, include validation (tests, tooling), and highlight dependencies or parallelizable work.
7. Validate with `openspec validate <id> --strict` and resolve every issue before sharing the proposal.

**Reference**
- Use `openspec show <id> --json --deltas-only` or `openspec show <spec> --type spec` to inspect details when validation fails.
- Search existing requirements with `rg -n "Requirement:|Scenario:" openspec/specs` before writing new ones.
- Explore the codebase with `rg <keyword>`, `ls`, or direct file reads so proposals align with current implementation realities.
<!-- OPENSPEC:END -->

# 用户查询接口提案

## 提案概述
已创建完整的用户查询接口提案，包含以下文件：

### 📋 提案文档
- **`openspec/changes/add-user-query-api/proposal.md`** - 变更概述和动机
- **`openspec/changes/add-user-query-api/tasks.md`** - 详细实现任务列表
- **`openspec/changes/add-user-query-api/design.md`** - 架构设计和技术决策

### 📝 规范文档
- **`openspec/changes/add-user-query-api/specs/user-query-api/spec.md`** - 完整的接口规范

## 提案要点

### 🎯 核心功能
- 实现 `GET /api/v1/users/{id}` 用户查询接口
- 遵循 Ginx 框架的 `HandleOperator` 接口规范
- 支持参数校验和错误处理
- 提供中英文国际化错误消息

### 🏗️ 技术实现
- **接口定义**: 实现标准的 RESTful API
- **参数校验**: 使用 `validate` tag 进行用户ID校验
- **错误处理**: 定义用户相关的结构化错误码
- **数据模型**: 创建标准的用户信息结构体
- **模拟数据**: 使用内存数据演示接口功能

### 📦 文件结构
```
apis/user/get_user_info.go     # 接口实现
models/user.go                 # 用户数据模型
constants/status_error/error.go # 错误码定义（扩展）
apis/root.go                   # 路由注册（修改）
```

### ✅ 验收标准
- [ ] 接口正确实现 `HandleOperator` 接口
- [ ] 参数校验和错误处理正常工作
- [ ] 错误消息支持中英文国际化
- [ ] 代码生成工具正常运行
- [ ] 接口可通过 HTTP 请求正常访问

## 下一步行动
1. 运行 `openspec validate add-user-query-api --strict` 验证提案
2. 根据验证结果调整规范
3. 开始实现任务列表中的具体工作项