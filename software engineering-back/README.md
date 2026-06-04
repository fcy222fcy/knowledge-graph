# 软件工程知识问答平台 MVP

一个面向课程演示的知识问答平台原型，采用 Vue 3 + Go Gin + GORM + MySQL。

## 后端架构

后端按 MVC 思路组织：

- `internal/model`：GORM 数据模型。
- `internal/controller`：HTTP 请求绑定、响应返回。
- `internal/service`：问答、答题、分析等业务逻辑。
- `internal/repository`：数据库访问，封装 GORM 查询。
- `internal/routes`：路由注册和依赖组装。
- `internal/database`：MySQL 连接和 AutoMigrate。
- `internal/seed`：演示数据初始化。

## 运行前准备

### 1. 创建 MySQL 数据库

```sql
CREATE DATABASE software_qa_platform DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. 配置 `.env`

复制 `.env.example` 为 `.env`，填写数据库账号密码：

```powershell
Copy-Item .env.example .env
```

需要重点修改：

```env
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=你的MySQL密码
DB_NAME=software_qa_platform
```

### 3. 启动后端

```powershell
go mod tidy
go run ./cmd/server
```

默认地址由 `.env` 的 `SERVER_PORT` 决定。`.env.example` 使用 `8080`；如果你改成 `8899`，后端地址就是 `http://localhost:8899`。

### 4. 启动前端

```powershell
cd frontend-vue-app
npm install
npm run dev
```

默认地址：`http://localhost:5173`

### 5. 访问系统

浏览器打开：

```text
http://localhost:5173
```

前端会通过 Vite 代理请求 `http://localhost:8080` 的后端接口。第一次启动后端时，GORM 会自动建表；如果知识点、题目和资料表为空，会写入演示数据。

## 前端风格约束

后续每次修改前端都必须遵守：

- 极简风格，优先清晰、留白、低装饰。
- 配色避免明显“AI 产品感”，不使用大面积蓝紫渐变、发光、霓虹、科技感背景。
- 主色以白、灰、墨色、中性边框为主，只用少量低饱和强调色表达状态。
- 页面应像课程工具/后台应用，不做营销式首页。

## 示例资料

用户提供的示例资料已整理到：

- `doc/examples/Redis(Remote dictionary server).md`
- `doc/examples/SQL语法.md`
- `doc/examples/名词解释.md`

这些文件可作为后续“资源管理、文档解析、知识抽取”的样例输入。

## MVP 功能

- 知识问答：关键词匹配返回答案、引用和相关知识点
- 知识图谱：返回 MySQL 中的节点和关系
- 图谱构建：从课程资料中按规则抽取知识点和关系，并写入图谱数据
- 资源管理：保存和展示资料元数据
- 题库练习：提交答案并返回评分解析
- 学习分析：按演示用户聚合问答、答题和薄弱知识点

## 知识图谱构建

前端“知识图谱”页面提供资料选择和“构建图谱”按钮。系统会基于资料说明和示例资料文件内容，使用确定性规则抽取知识点与关系，写入 `KnowledgePoint` 和 `KnowledgeRelation` 表。构建完成后，页面会显示新增知识点、新增关系、跳过重复知识点数量，并自动刷新图谱。

相关接口：

- `POST /api/v1/graph/build`
- `GET /api/v1/graph/build/latest`
- `GET /api/v1/graph`
