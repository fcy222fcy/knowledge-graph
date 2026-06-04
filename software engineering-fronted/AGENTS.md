# AGENTS.md

## 项目介绍

**SE智图问答** 是一个基于知识图谱的软件工程课程问答平台。

系统支持上传课程相关资料（PDF、PPT、Word、Markdown等），自动抽取知识点及其关系，构建软件工程课程知识图谱。学生可通过自然语言提问，系统基于知识图谱进行智能检索和回答，同时提供可视化学习功能。

### 核心功能

- **用户模块**：注册、登录、JWT鉴权
- **资料管理**：文件上传、解析状态管理
- **知识图谱**：图谱可视化展示、节点搜索、关系筛选
- **智能问答**：自然语言提问、知识来源展示、历史会话
- **学习统计**：学习时长、问答次数、趋势图表

### 系统架构

```
┌─────────────────────────────────────┐
│         前端 (Vue 3 + Vite)         │
└─────────────────────────────────────┘
                 │
                 │ HTTP (REST API)
                 ↓
┌─────────────────────────────────────┐
│         后端 (Go + Gin)             │
│  - API网关 / 业务逻辑              │
│  - 调用 Python AI 服务             │
└─────────────────────────────────────┘
                 │
        ┌────────┴────────┐
        ↓                 ↓
┌──────────────┐  ┌──────────────┐
│ Neo4j 图数据库│  │ FAISS 向量库  │
│ (知识图谱)   │  │ (语义检索)    │
└──────────────┘  └──────────────┘
```

---

## 技术栈

### 前端

| 技术 | 版本 | 说明 |
|------|------|------|
| Vue | 3.x | 前端框架 |
| Vite | 5.x | 构建工具 |
| Tailwind CSS | 3.x | 样式框架 |
| Vue Router | 4.x | 路由管理 |
| Pinia | 2.x | 状态管理 |
| Axios | 1.x | HTTP 请求 |

### 后端

| 技术 | 说明 |
|------|------|
| Go 1.21+ | 后端语言 |
| Gin | Web 框架 |
| GORM | ORM 框架 |
| JWT | 认证鉴权 |

### 数据存储

| 技术 | 用途 |
|------|------|
| MySQL | 用户、资料、问答记录等基础数据 |
| Neo4j | 知识图谱节点和关系 |
| FAISS | 文本向量语义检索 |

---

## 目录结构

```
software engineering-fronted/
├── src/                        # 前端源码
│   ├── assets/                 # 静态资源（图片、字体等）
│   ├── components/             # 公共组件
│   │   └── common/             # 通用基础组件
│   ├── composables/            # 组合式函数 (hooks)
│   ├── layouts/                # 布局组件
│   ├── router/                 # 路由配置
│   ├── services/               # 接口请求封装
│   │   ├── auth.ts             # 认证相关接口
│   │   ├── file.ts             # 资料管理接口
│   │   ├── graph.ts            # 知识图谱接口
│   │   ├── qa.ts               # 问答中心接口
│   │   └── stats.ts            # 学习统计接口
│   ├── stores/                 # Pinia 状态管理
│   ├── types/                  # TypeScript 类型定义
│   ├── utils/                  # 工具函数
│   │   └── request.ts          # Axios 封装
│   ├── views/                  # 页面视图
│   │   ├── auth/               # 登录/注册页
│   │   ├── home/               # 首页
│   │   ├── graph/              # 知识图谱页
│   │   ├── qa/                 # 问答中心页
│   │   ├── files/              # 资料管理页
│   │   └── stats/              # 分析统计页
│   ├── App.vue                 # 根组件
│   └── main.ts                 # 入口文件
├── public/                     # 公共静态资源
├── docs/                       # 项目文档
├── prototype/                  # 原型图
├── index.html                  # HTML 入口
├── package.json                # 依赖配置
├── vite.config.ts              # Vite 配置
├── tailwind.config.js          # Tailwind 配置
├── tsconfig.json               # TypeScript 配置
└── AGENTS.md                   # 本文档
```

---

## 常用命令

```bash
# 安装依赖
npm install

# 启动开发服务器 (http://localhost:5173)
npm run dev

# 代码检查
npm run lint

# 类型检查
npm run type-check

# 构建生产版本
npm run build

# 预览生产版本
npm run preview
```

---

## 开发规则

### 代码组织

1. **页面入口** → `src/views/xxx/` 目录
2. **接口请求** → `src/services/` 目录
3. **类型定义** → `src/types/` 目录
4. **公共组件** → `src/components/common/` 目录
5. **业务组件** → `src/views/xxx/components/` 目录（与页面同级）
6. **组合式函数** → `src/composables/` 目录

### 组件拆分规则

每个后台管理页面建议拆分为：

```
views/xxx/
├── index.vue              # 页面主组件（状态管理、逻辑编排）
├── components/
│   ├── SearchForm.vue     # 搜索表单组件
│   ├── DataTable.vue      # 数据表格组件
│   └── FormModal.vue      # 新增/编辑弹窗组件
```

### 组件职责

| 组件 | 职责 |
|------|------|
| `index.vue` | 管理页面状态、调用接口、处理查询/分页/增删改、组合子组件 |
| `SearchForm.vue` | 展示查询条件、触发查询/重置事件 |
| `DataTable.vue` | 展示表格数据、显示 loading、触发编辑/删除操作 |
| `FormModal.vue` | 展示新增/编辑表单、表单校验、提交数据 |

### 代码规范

1. **不要在页面组件里直接写 axios 请求**，必须通过 `services/` 封装调用
2. **TypeScript 类型统一写在 `types/` 目录**
3. **页面需要有 loading 状态、错误提示、空数据处理**
4. **表单提交需要有 loading 状态，防止重复提交**
5. **删除操作需要二次确认**
6. **不要随便改动无关文件**
7. **组件命名使用 PascalCase**（如 `SearchForm.vue`）
8. **文件/目录命名使用 kebab-case**（如 `auth-login`）

---

## 接口请求规则

### 基础配置

- **API 前缀**：`/api/v1`
- **开发环境代理**：`http://localhost:8080`
- **超时时间**：10 秒

### 请求封装

所有接口请求统一使用 `src/utils/request.ts` 中的 Axios 实例：

```typescript
// src/utils/request.ts
import axios from 'axios'
import { useUserStore } from '@/stores/user'

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
})

// 请求拦截器：自动携带 Token
request.interceptors.request.use((config) => {
  const userStore = useUserStore()
  if (userStore.token) {
    config.headers.Authorization = `Bearer ${userStore.token}`
  }
  return config
})

// 响应拦截器：统一错误处理
request.interceptors.response.use(
  (response) => response.data,
  (error) => {
    // 401 跳转登录
    // 其他错误提示
    return Promise.reject(error)
  }
)

export default request
```

### 接口返回格式

后端接口统一返回格式：

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

### 接口编写示例

```typescript
// src/services/auth.ts
import request from '@/utils/request'
import type { LoginParams, LoginResponse } from '@/types/auth'

// 用户登录
export function login(data: LoginParams) {
  return request.post<LoginResponse>('/auth/login', data)
}

// 用户注册
export function register(data: RegisterParams) {
  return request.post('/auth/register', data)
}

// 获取用户信息
export function getUserInfo() {
  return request.get('/user/info')
}
```

### 类型定义示例

```typescript
// src/types/auth.ts
export interface LoginParams {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: {
    id: number
    username: string
    email: string
  }
}
```

---

## AI 修改代码注意事项

### 修改前

1. **先说明计划修改哪些文件**，以及每个文件的修改内容
2. **阅读现有代码**，理解当前实现逻辑
3. **确认不会破坏已有功能**

### 修改中

1. **按照现有项目结构生成代码**
2. **遵循 AGENTS.md 中的开发规则**
3. **保持代码风格一致**（注释、命名、缩进等）
4. **不要直接修改无关文件**

### 修改后

1. **说明修改了哪些文件**
2. **说明新增了哪些文件**
3. **列出需要执行的检查命令**：

```bash
npm run lint        # 代码检查
npm run type-check  # 类型检查
npm run build       # 构建验证
```

4. **如果命令报错，优先检查**：
   - TypeScript 类型是否完整
   - 是否存在未使用变量
   - 是否存在错误的接口字段
   - 是否存在错误的导入路径

---

## 页面清单

| 页面 | 路径 | 说明 |
|------|------|------|
| 登录页 | `/login` | 用户登录 |
| 注册页 | `/register` | 用户注册 |
| 首页 | `/` | Dashboard，展示今日统计、功能入口、最近提问 |
| 知识图谱 | `/graph` | 图谱可视化、节点搜索、关系筛选 |
| 问答中心 | `/qa` | 自然语言问答、历史会话 |
| 资料管理 | `/files` | 文件上传、列表展示、删除 |
| 分析统计 | `/stats` | 学习时长、问答次数、趋势图表 |

---

## API 接口清单

### 认证模块

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/auth/register` | POST | 用户注册 |
| `/api/v1/auth/login` | POST | 用户登录 |

### 用户模块

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/user/info` | GET | 获取用户信息 |
| `/api/v1/user/update` | PUT | 更新用户信息 |

### 资料管理模块

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/files/upload` | POST | 上传文件 |
| `/api/v1/files` | GET | 获取文件列表 |
| `/api/v1/files/:id` | DELETE | 删除文件 |

### 知识图谱模块

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/graph` | GET | 获取图谱数据 |
| `/api/v1/graph/search` | GET | 搜索知识点 |
| `/api/v1/graph/node/:id` | GET | 获取节点详情 |
| `/api/v1/graph/build` | POST | 触发图谱构建 |

### 问答模块

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/ask` | POST | 提问 |
| `/api/v1/qa/history` | GET | 获取历史会话 |
| `/api/v1/qa/session/:id` | GET | 获取会话详情 |

### 统计模块

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/v1/stats/overview` | GET | 获取统计概览 |
| `/api/v1/stats/trend` | GET | 获取学习趋势 |
