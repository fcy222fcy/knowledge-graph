# AGENTS.md

> 每次修改代码前必须先阅读本文档，遵守其中的规则。

## 项目介绍

软件工程知识问答平台后端，面向课程演示场景。

主要功能模块：

- 知识问答：关键词匹配返回答案、引用和相关知识点
- 知识图谱：返回节点和关系数据
- 图谱构建：从课程资料中抽取知识点与关系，写入图谱
- 资源管理：保存和展示资料元数据
- 题库练习：提交答案并返回评分解析
- 学习分析：聚合问答、答题和薄弱知识点

Go 后端提供 REST API，前端 Vue 3 通过 Vite 代理调用后端接口。

## 技术栈

- 语言：Go 1.24.0
- HTTP 框架：Gin v1.10.1
- ORM：GORM v1.25.12
- 数据库：MySQL（库名 `software_qa_platform`，字符集 utf8mb4）
- 跨域：gin-contrib/cors v1.7.6
- 配置：`.env` 环境变量

端口分配：

| 服务 | 端口 |
|------|------|
| 前端 | 5173 |
| Go 后端 | 8080 |

## 目录结构

```
cmd/server/              -- 程序入口 (main.go)
internal/
  controller/            -- 控制器层（参数绑定、调用 service、返回响应）
  service/               -- 业务逻辑层
  repository/            -- 数据库访问层（GORM CRUD）
  model/                 -- GORM 数据模型
  dto/                   -- 请求/响应结构体（Request/Response 后缀）
  routes/                -- 路由注册和依赖组装
  database/              -- MySQL 连接和 AutoMigrate
  seed/                  -- 演示数据初始化
  middleware/            -- 中间件（CORS、权限等）
doc/                     -- 项目文档
  api.md                 -- 接口文档（需同步维护）
  business-rules.md      -- 业务规则
  backend-rules.md       -- 后端代码规范
prototype/               -- 前端 UI 原型
  se-platform-mvp.html   -- MVP 版本完整原型图（可直接在浏览器打开）
```

## 分层规则

### Controller 层

负责：

- 接收 HTTP 请求
- 参数绑定和校验（使用 `binding` tag）
- 调用 Service 方法
- 返回统一 JSON 响应

禁止：

- 不直接写 SQL
- 不写复杂业务逻辑
- 不直接操作数据库

### Service 层

负责：

- 业务逻辑处理
- 权限判断
- 状态管理
- 调用 Repository
- 开启和管理事务

### Repository 层

负责：

- 数据库查询（CRUD）
- 封装 GORM 调用

禁止：

- 不写复杂业务逻辑

### Model 层

- 对应 MySQL 数据库表，使用 GORM tag 映射字段

### DTO 层

- 请求结构体：以 `Request` 结尾（如 `CreateKnowledgePointRequest`）
- 响应结构体：以 `Response` 结尾（如 `KnowledgePointResponse`）
- JSON tag 使用 `snake_case`
- 使用 `binding:"required"` 标记必填字段

## 统一响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

错误响应：

```json
{
  "code": 400,
  "message": "错误描述",
  "data": null
}
```

## 常用命令

```bash
# 启动后端
go run ./cmd/server

# 整理依赖
go mod tidy

# 运行全部测试
go test ./...

# 测试指定模块
go test ./internal/service/...
go test ./internal/repository/...
```

## AI 修改代码要求

1. **先说明需要修改哪些文件**，列出完整文件路径和修改原因。
2. **遵守 controller / service / repository 分层**，各层职责不能混淆。
3. **不要随便修改无关模块**，只改与需求相关的文件。
4. **新增接口时需要同时更新 `docs/api.md`**，保持接口文档与代码同步。
5. **修改接口字段时需要提醒前端同步修改**，避免前后端不一致。
6. **写完后说明需要执行哪些测试命令**，包括 `go test` 和 curl 测试命令。

## 补充说明

- 项目当前无源代码，首次启动后 GORM 会自动建表并写入演示数据。
- 数据库配置在 `.env` 文件中管理，不要硬编码。
- 多表修改必须使用 GORM 事务。
- 业务错误使用 `errors.NewBizError("message")` 返回，不要直接 panic。
