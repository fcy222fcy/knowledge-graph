# 基于知识图谱的软件工程课程问答平台

## 项目简介

本项目是一个基于知识图谱的软件工程课程智能问答平台，旨在为学生提供智能化的学习辅助服务。平台通过构建软件工程领域的知识图谱，结合大语言模型技术，实现课程知识的智能问答、文档解析、知识图谱可视化等功能。

## 主要功能

- **智能问答**：基于知识图谱和文档检索的智能问答系统
- **文档管理**：支持 PDF、Markdown、DOCX 格式的文档上传与解析
- **知识图谱**：自动构建和可视化软件工程知识图谱
- **作业管理**：支持作业发布、提交和自动评分
- **学习统计**：提供学习数据统计和分析功能
- **管理员后台**：提供用户管理、数据统计等管理功能

## 技术栈

### 后端
- Go 1.21+
- Gin Web 框架
- MySQL 数据库
- Neo4j 图数据库（知识图谱存储）
- 向量数据库（文档向量化）

### 前端
- Vue 3
- TypeScript
- Element Plus UI 组件库
- Vite 构建工具

## 项目结构

```
├── software engineering-back/    # 后端服务
│   ├── internal/                 # 核心业务逻辑
│   │   ├── api/                 # API 路由和控制器
│   │   ├── app/                 # 应用启动配置
│   │   ├── model/               # 数据模型定义
│   │   ├── repository/          # 数据访问层
│   │   └── service/             # 业务服务层
│   └── pkg/                     # 公共工具包
│       └── database/            # 数据库连接
├── software engineering-fronted/ # 前端应用
│   ├── admin-vue-app/           # 管理员前端
│   └── frontend-vue-app/        # 用户前端
└── 课程文档/                     # 项目文档（已排除）
```

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 16+
- MySQL 8.0+
- Neo4j 4.4+

### 后端启动
```bash
cd software engineering-back
go mod tidy
go run main.go
```

### 前端启动
```bash
cd software engineering-fronted/frontend-vue-app
npm install
npm run dev
```

### 管理员端启动
```bash
cd software engineering-fronted/admin-vue-app
npm install
npm run dev
```

## 开发说明

- 后端默认运行在 `http://localhost:8080`
- 前端开发服务器运行在 `http://localhost:5173`
- 管理员端运行在 `http://localhost:5174`

## 许可证

本项目仅用于学习和研究目的。
