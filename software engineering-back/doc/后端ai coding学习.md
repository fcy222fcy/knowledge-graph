# 后端 AI 协作指南

## 1. 文档目的

这份文档用于指导后端开发者如何配合 AI Coding。

适合场景：

- 你主要写 Go / Gin / MySQL / Redis 后端
- 你希望 AI 帮你生成前端页面
- 你希望后端接口文档清楚，让 AI 不乱猜接口
- 你希望后端代码也能通过 AI 辅助生成、修改、测试

核心思想：

> 后端负责把接口、字段、业务规则说清楚，AI 才能稳定生成前端页面或后端代码。

---

## 2. 整体使用流程

推荐按下面 5 个阶段来做：

1. 项目信息阶段：告诉 AI 后端项目是什么
2. 接口文档阶段：告诉 AI 接口路径、参数、返回值
3. 业务规则阶段：告诉 AI 每个接口背后的业务逻辑
4. 代码规范阶段：告诉 AI 后端代码应该怎么分层
5. 测试验证阶段：让 AI 写完代码后生成测试和验证命令

---

# 阶段一：项目信息阶段

## 目标

让 AI 先理解你的后端项目结构。

## 需要生成的文档

建议在项目根目录创建：

```text
AGENTS.md
```

如果前后端是两个仓库，后端仓库也单独放一份 `AGENTS.md`。

## AGENTS.md 应该写什么

建议包含：

- 项目介绍
- 技术栈
- 目录结构
- 数据库说明
- 中间件说明
- 常用命令
- 分层规则
- AI 修改代码要求

## AGENTS.md 示例

```md
# AGENTS.md

## 项目介绍

这是一个 Go + Gin 后端项目，为后台管理系统和用户端提供接口。

项目包含用户、店铺、评价、权限、系统设置、排行榜等模块。

## 技术栈

- Go
- Gin
- GORM
- MySQL
- Redis
- JWT
- Zap Logger
- WebSocket

## 目录结构

- `cmd/`：程序入口
- `internal/router/`：路由注册
- `internal/controller/`：控制器层
- `internal/service/`：业务逻辑层
- `internal/repository/`：数据库访问层
- `internal/model/`：数据库模型
- `internal/dto/`：请求和响应结构
- `internal/middleware/`：中间件
- `pkg/response/`：统一响应
- `pkg/logger/`：日志工具
- `docs/`：项目文档

## 分层规则

1. controller 只负责参数绑定、调用 service、返回响应。
2. service 负责业务逻辑。
3. repository 负责数据库操作。
4. model 对应数据库表。
5. dto 用于请求参数和响应数据。
6. 不要在 controller 里直接写 SQL。
7. 不要在 repository 里写复杂业务判断。

## 常用命令

```bash
go run ./cmd/server
go test ./...
go test ./internal/service/...
go test ./internal/repository/...
```

## AI 修改代码要求

1. 先说明需要修改哪些文件。
2. 遵守 controller / service / repository 分层。
3. 不要随便修改无关模块。
4. 新增接口时需要同时更新 docs/api.md。
5. 修改接口字段时需要提醒前端同步修改。
6. 写完后说明需要执行哪些测试命令。
```

## 这个阶段需要对 AI 说什么

```text
请先阅读 AGENTS.md，理解当前 Go 后端项目的技术栈、目录结构和分层规则。
后续生成或修改代码时，必须遵守这些规则。
```

---

# 阶段二：接口文档阶段

## 目标

让 AI 清楚知道接口怎么调用。

这个阶段非常重要，因为前端 AI 主要依赖接口文档来生成页面。

## 需要生成的文档

建议创建：

```text
docs/api.md
```

## api.md 应该写什么

每个接口建议包含：

- 模块名称
- 接口名称
- 请求方法
- 请求路径
- 请求参数
- 请求体
- 返回示例
- 字段说明
- 错误情况

## api.md 示例

```md
# 接口文档

## 管理员模块

### 获取管理员列表

#### 请求方式

GET `/api/v1/admin/admin-users`

#### 请求参数

| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| page | number | 是 | 页码 |
| page_size | number | 是 | 每页数量 |
| keyword | string | 否 | 搜索关键词 |
| status | number | 否 | 状态：1 启用，0 禁用 |

#### 返回示例

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "nick_name": "admin",
        "status": 1,
        "created_at": "2026-06-04 10:00:00"
      }
    ],
    "total": 1
  }
}
```

#### 返回字段说明

| 字段 | 类型 | 说明 |
|---|---|---|
| id | number | 管理员 ID |
| nick_name | string | 管理员昵称 |
| status | number | 状态 |
| created_at | string | 创建时间 |

---

### 新增管理员

#### 请求方式

POST `/api/v1/admin/admin-users`

#### 请求体

```json
{
  "nick_name": "admin",
  "password": "123456",
  "role_ids": [1, 2],
  "status": 1
}
```

#### 请求字段说明

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| nick_name | string | 是 | 管理员昵称 |
| password | string | 是 | 密码 |
| role_ids | number[] | 是 | 角色 ID 列表 |
| status | number | 是 | 状态：1 启用，0 禁用 |

#### 返回示例

```json
{
  "code": 200,
  "message": "创建成功",
  "data": null
}
```

---

### 修改管理员

#### 请求方式

PUT `/api/v1/admin/admin-users/:id`

#### 请求体

```json
{
  "nick_name": "admin",
  "role_ids": [1, 2],
  "status": 1
}
```

---

### 删除管理员

#### 请求方式

DELETE `/api/v1/admin/admin-users/:id`
```

## 这个阶段需要对 AI 说什么

```text
请根据 docs/api.md 中的接口文档，生成前端 services 和 types。
接口路径、字段名、请求参数、返回结构必须严格按照文档，不要自己猜。
```

## 这个阶段怎么用

当你写完一个后端模块后，要同步更新 `docs/api.md`。

例如你新增了“店铺管理接口”，就要在 `docs/api.md` 加上：

```text
店铺列表接口
新增店铺接口
编辑店铺接口
删除店铺接口
启用/禁用店铺接口
```

然后再让 AI 写前端页面。

---

# 阶段三：业务规则阶段

## 目标

让 AI 理解接口背后的业务规则。

接口文档只说明“怎么调用”，业务规则说明“为什么这么做”。

## 需要生成的文档

建议创建：

```text
docs/business-rules.md
```

## business-rules.md 应该写什么

建议按模块写：

- 模块说明
- 业务规则
- 权限规则
- 状态流转
- 特殊限制
- 边界情况

## business-rules.md 示例

```md
# 业务规则文档

## 管理员模块

### 业务说明

管理员模块用于管理后台系统账号。

### 业务规则

1. 超级管理员不能被删除。
2. 普通管理员不能修改自己的角色。
3. 禁用管理员后，该管理员不能继续登录。
4. 修改管理员角色后，需要重新计算权限。
5. 如果系统支持强制下线，修改角色后应让该管理员重新登录。

### 权限规则

| 操作 | 权限标识 |
|---|---|
| 查看管理员 | admin_user:view |
| 新增管理员 | admin_user:create |
| 编辑管理员 | admin_user:update |
| 删除管理员 | admin_user:delete |
| 分配角色 | admin_user:assign_role |

### 状态说明

| 状态值 | 说明 |
|---|---|
| 1 | 启用 |
| 0 | 禁用 |

### 边界情况

1. 删除不存在的管理员，应返回错误。
2. 新增管理员时昵称不能重复。
3. 密码不能为空，且需要加密存储。
```

## 这个阶段需要对 AI 说什么

```text
请根据 docs/business-rules.md 中的业务规则实现后端逻辑。
不要只写 CRUD，需要处理权限、状态、重复数据、边界情况。
```

## 对前端 AI 有什么用

前端也可以根据业务规则生成更合理的交互。

例如：

- 超级管理员不显示删除按钮
- 禁用状态显示红色标签
- 权限不足时隐藏按钮
- 删除前显示确认框

你可以对前端 AI 说：

```text
请根据 docs/business-rules.md 中的业务规则完善页面交互。
例如超级管理员不能删除，状态需要用标签展示。
```

---

# 阶段四：代码规范阶段

## 目标

让 AI 写后端代码时遵守你的 Go 项目结构。

## 需要生成的文档

建议创建：

```text
docs/backend-rules.md
```

## backend-rules.md 应该写什么

建议包含：

- 分层规则
- 命名规则
- 错误处理规则
- 日志规则
- 事务规则
- Redis 使用规则
- 权限校验规则

## backend-rules.md 示例

```md
# 后端代码规范

## 分层规则

### controller 层

负责：

- 接收请求
- 参数绑定
- 参数校验
- 调用 service
- 返回统一响应

不能做：

- 不直接写 SQL
- 不写复杂业务逻辑
- 不直接操作 Redis

### service 层

负责：

- 业务逻辑
- 权限判断
- 状态判断
- 调用 repository
- 开启事务

### repository 层

负责：

- 数据库查询
- 数据库新增、修改、删除
- 不写复杂业务逻辑

## DTO 规则

请求参数使用 `Request` 结尾：

```go
type CreateAdminUserRequest struct {
    NickName string `json:"nick_name" binding:"required"`
    Password string `json:"password" binding:"required"`
}
```

响应结构使用 `Response` 结尾：

```go
type AdminUserResponse struct {
    ID       int    `json:"id"`
    NickName string `json:"nick_name"`
}
```

## 错误处理规则

业务错误返回业务错误码，不直接 panic。

示例：

```go
return errors.NewBizError("管理员不存在")
```

## 日志规则

关键失败操作需要记录日志：

```go
logger.Error("创建管理员失败", zap.Error(err))
```

## 事务规则

多个表同时修改时必须使用事务。

例如：

- 创建管理员并分配角色
- 修改管理员角色
- 删除角色并删除角色权限关系

## 权限规则

需要管理员权限的接口必须使用权限中间件。

新增接口时，需要同步补充权限标识。
```

## 这个阶段需要对 AI 说什么

```text
请按照 docs/backend-rules.md 中的后端代码规范实现接口。
要求：
1. controller 只做参数处理和响应。
2. service 写业务逻辑。
3. repository 写数据库操作。
4. 涉及多表修改必须使用事务。
5. 新增接口后同步更新 docs/api.md。
```

---

# 阶段五：测试验证阶段

## 目标

让 AI 写完后端代码后，生成测试或验证命令。

## 需要生成的文档

建议创建：

```text
docs/test-guide.md
```

## test-guide.md 应该写什么

建议包含：

- 单元测试命令
- 接口 curl 测试
- 数据库检查 SQL
- Redis 检查命令
- 常见错误排查

## test-guide.md 示例

```md
# 测试验证指南

## Go 单元测试

```bash
go test ./...
```

测试指定模块：

```bash
go test ./internal/service/...
go test ./internal/repository/...
```

## 接口测试

### 获取管理员列表

```bash
curl -X GET "http://localhost:9090/api/v1/admin/admin-users?page=1&page_size=10" \
  -H "Authorization: Bearer <token>"
```

### 新增管理员

```bash
curl -X POST "http://localhost:9090/api/v1/admin/admin-users" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "nick_name": "test_admin",
    "password": "123456",
    "role_ids": [1],
    "status": 1
  }'
```

## 数据库检查

```sql
SELECT * FROM admin_users ORDER BY id DESC LIMIT 10;
SELECT * FROM admin_user_roles WHERE admin_user_id = 1;
```

## Redis 检查

```bash
redis-cli keys "*admin*"
```

## 常见问题

1. 接口 401：检查 token 是否过期。
2. 接口 403：检查权限是否配置。
3. 接口 500：查看后端日志。
4. 数据没更新：检查事务是否提交。
```

## 这个阶段需要对 AI 说什么

```text
写完代码后，请生成：
1. go test 命令。
2. curl 测试命令。
3. 需要检查的数据库 SQL。
4. 可能的错误排查点。
```

---

# 后端生成接口时的推荐 Prompt

以后你可以复制下面这段给 AI：

```text
你是一个 Go 后端 AI Coding 助手。

请根据以下文档实现后端接口：

1. AGENTS.md
2. docs/api.md
3. docs/business-rules.md
4. docs/backend-rules.md

开发要求：

1. 使用 Go + Gin + GORM。
2. controller 只负责参数绑定、调用 service、返回响应。
3. service 负责业务逻辑。
4. repository 负责数据库操作。
5. DTO 放到 dto 目录。
6. model 放到 model 目录。
7. 多表修改必须使用事务。
8. 需要权限校验的接口必须接入权限中间件。
9. 新增或修改接口后，需要更新 docs/api.md。
10. 写完后说明修改了哪些文件。
11. 写完后给出 go test、curl、SQL 检查命令。

现在请根据文档实现指定模块。
```

---

# 后端配合 AI 写前端的推荐 Prompt

如果你的目标是让 AI 根据后端接口写前端，可以这样说：

```text
你是一个前端 AI Coding 助手。

我已经在 docs/api.md 中写好了后端接口文档。
请根据接口文档生成前端页面。

要求：

1. 不要自己猜接口路径。
2. 不要自己改字段名。
3. 请求参数、请求体、返回结构必须以 docs/api.md 为准。
4. 接口请求封装到 services 目录。
5. 类型定义放到 types 目录。
6. 页面组件按功能拆分。
7. 写完后告诉我修改了哪些文件。
8. 写完后告诉我需要执行哪些命令验证。
```

---

# 最小可行文档清单

如果你不想一开始写太多文档，至少准备这 3 个：

```text
AGENTS.md
docs/api.md
docs/pages.md
```

含义：

| 文档 | 作用 |
|---|---|
| AGENTS.md | 告诉 AI 项目结构和开发规则 |
| docs/api.md | 告诉 AI 后端接口 |
| docs/pages.md | 告诉 AI 前端页面需求 |

有了这三个，AI 就能比较稳定地帮你写前端页面。

---

# 推荐项目文档目录

```text
docs/
├── api.md
├── pages.md
├── business-rules.md
├── backend-rules.md
├── frontend-rules.md
└── test-guide.md
```

如果是前后端分离两个仓库：

后端仓库：

```text
AGENTS.md
docs/
├── api.md
├── business-rules.md
├── backend-rules.md
└── test-guide.md
```

前端仓库：

```text
AGENTS.md
docs/
├── pages.md
├── frontend-rules.md
└── api.md
```

其中前端仓库里的 `docs/api.md` 可以从后端仓库复制过来，或者只放前端需要用到的接口。
