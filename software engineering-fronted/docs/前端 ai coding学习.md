# 前端 AI Coding 指南

## 1. 文档目的

这份文档用于指导 AI 帮你编写前端代码。

适合场景：

- 你是后端开发，想让 AI 帮你写前端页面
- 你已经有后端接口，希望 AI 根据接口生成页面
- 你希望 AI 写出来的前端代码结构统一、可维护、能运行

核心思想：

> 不要直接让 AI “帮我写页面”，而是先告诉 AI：项目结构是什么、接口是什么、页面需要什么功能、代码要怎么写、写完怎么检查。

---

## 2. 整体使用流程

推荐按下面 5 个阶段来做：

1. 项目信息阶段：告诉 AI 前端项目是什么
2. 页面需求阶段：告诉 AI 页面要长什么样、有什么功能
3. 接口对接阶段：告诉 AI 后端接口怎么用
4. 代码规范阶段：告诉 AI 代码应该怎么组织
5. 检查验证阶段：让 AI 写完后根据命令检查代码

---

# 阶段一：项目信息阶段

## 目标

让 AI 先理解你的前端项目，而不是直接开始写代码。

## 需要生成的文档

建议在项目根目录创建：

```text
AGENTS.md
```

## AGENTS.md 应该写什么

这个文件是给 AI 看的“项目说明书”。

建议包含：

- 项目介绍
- 技术栈
- 目录结构
- 常用命令
- 开发规则
- 接口请求规则
- AI 修改代码时的注意事项

## AGENTS.md 示例

```md
# AGENTS.md

## 项目介绍

这是一个后台管理系统前端项目，用于管理用户、店铺、评价、权限、系统设置等功能。

后端由 Go + Gin 提供接口，前端通过 `/api/v1` 前缀访问后端接口。

## 技术栈

- vue
- Tailwind CSS / Ant Design / shadcn-ui

## 目录结构

- `app/`：页面路由
- `components/`：公共组件
- `services/`：接口请求封装
- `types/`：TypeScript 类型定义
- `hooks/`：自定义 Hook
- `lib/`：工具函数
- `docs/`：项目文档

## 开发规则

1. 页面代码放在 `app/` 目录。
2. 接口请求统一写在 `services/` 目录。
3. TypeScript 类型统一写在 `types/` 目录。
4. 不要在页面组件里直接写 axios 请求。
5. 页面需要有 loading、错误提示、空数据处理。
6. 表单提交需要有 loading 状态。
7. 删除操作需要二次确认。
8. 不要随便改动无关文件。

## 常用命令

```bash
npm install
npm run dev
npm run lint
npm run type-check
npm run build
```

## 后端接口约定

后端接口统一使用 `/api/v1` 前缀。

接口返回格式：

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

## AI 修改代码要求

AI 修改代码时需要：

1. 先说明计划修改哪些文件。
2. 按照现有项目结构生成代码。
3. 不要破坏已有功能。
4. 写完后说明需要执行哪些检查命令。
```

## 这个阶段需要对 AI 说什么

可以这样说：

```text
请先阅读 AGENTS.md，理解当前前端项目的技术栈、目录结构和开发规则。
后续生成代码时，必须遵守 AGENTS.md 中的规则。
```

## 这个阶段怎么用

每次让 AI 写前端代码之前，先让 AI 读取或参考 `AGENTS.md`。

例如：

```text
根据 AGENTS.md 的项目规范，帮我生成店铺管理页面。
```

---

# 阶段二：页面需求阶段

## 目标

告诉 AI 具体页面要做什么。

很多 AI 写出来的页面不好用，是因为你只说了：

```text
帮我写一个用户管理页面。
```

这太模糊了。

你应该告诉 AI：

- 页面路径是什么
- 页面有哪些区域
- 表格有哪些列
- 查询条件有哪些
- 按钮有哪些
- 弹窗有哪些
- 用户点击按钮后发生什么

## 需要生成的文档

建议创建：

```text
docs/pages.md
```

## pages.md 应该写什么

每个页面建议包含：

- 页面名称
- 页面路径
- 页面功能
- 页面结构
- 查询条件
- 表格字段
- 操作按钮
- 弹窗说明
- 交互说明

## pages.md 示例

```md
# 页面说明文档

## 管理员管理页面

### 页面路径

`/admin/admin-users`

### 页面功能

用于管理后台管理员账号。

### 页面结构

页面包含以下区域：

1. 查询区域
2. 操作按钮区域
3. 数据表格区域
4. 分页区域
5. 新增/编辑弹窗
6. 删除确认弹窗

### 查询条件

- 关键词：根据管理员昵称搜索
- 状态：启用 / 禁用

### 操作按钮

- 查询
- 重置
- 新增管理员

### 表格字段

| 字段 | 说明 |
|---|---|
| id | 管理员 ID |
| nick_name | 管理员昵称 |
| status | 状态 |
| created_at | 创建时间 |
| actions | 操作列 |

### 操作列

每一行包含：

- 编辑
- 删除
- 分配角色

### 新增管理员弹窗

字段：

- 昵称
- 密码
- 角色
- 状态

### 编辑管理员弹窗

字段：

- 昵称
- 角色
- 状态

### 交互说明

1. 页面打开时自动加载管理员列表。
2. 点击查询，根据查询条件刷新列表。
3. 点击重置，清空查询条件并重新加载列表。
4. 点击新增，打开新增管理员弹窗。
5. 点击编辑，打开编辑弹窗并回显当前行数据。
6. 点击删除，弹出确认框，确认后调用删除接口。
7. 新增、编辑、删除成功后，需要刷新列表。
```

## 这个阶段需要对 AI 说什么

```text
请根据 docs/pages.md 中的“管理员管理页面”说明，生成对应的前端页面。
页面需要包含查询、表格、分页、新增、编辑、删除功能。
```

## 这个阶段怎么用

当你想生成新页面时，先在 `docs/pages.md` 中写清楚页面需求。

然后让 AI 根据这个文档生成页面。

---

# 阶段三：接口对接阶段

## 目标

让 AI 知道前端要调用哪些后端接口。

如果接口不清楚，AI 很容易乱写接口路径、乱猜字段名。

## 需要生成的文档

建议创建：

```text
docs/api.md
```

不过接口文档主要由后端维护，前端 AI 只负责使用。

## 前端需要关注什么

前端需要从接口文档里知道：

- 请求方法：GET / POST / PUT / DELETE
- 请求路径
- 请求参数
- 请求体
- 返回数据结构
- 字段类型
- 分页格式

## 这个阶段需要对 AI 说什么

```text
请根据 docs/api.md 中的接口说明，生成前端请求封装。
要求：
1. 接口请求写到 services 目录。
2. TypeScript 类型写到 types 目录。
3. 页面中不要直接写 axios。
```

## 推荐生成的前端文件

以管理员管理页面为例，AI 应该生成：

```text
services/adminUserService.ts
types/adminUser.ts
app/admin/admin-users/page.tsx
components/admin-users/AdminUserSearchForm.tsx
components/admin-users/AdminUserTable.tsx
components/admin-users/AdminUserModal.tsx
```

## services 示例

```ts
// services/adminUserService.ts

import { request } from '@/lib/request'
import type {
  AdminUserListParams,
  AdminUserListResponse,
  CreateAdminUserPayload,
  UpdateAdminUserPayload,
} from '@/types/adminUser'

export function getAdminUserList(params: AdminUserListParams) {
  return request.get<AdminUserListResponse>('/api/v1/admin/admin-users', { params })
}

export function createAdminUser(data: CreateAdminUserPayload) {
  return request.post('/api/v1/admin/admin-users', data)
}

export function updateAdminUser(id: number, data: UpdateAdminUserPayload) {
  return request.put(`/api/v1/admin/admin-users/${id}`, data)
}

export function deleteAdminUser(id: number) {
  return request.delete(`/api/v1/admin/admin-users/${id}`)
}
```

## types 示例

```ts
// types/adminUser.ts

export interface AdminUser {
  id: number
  nick_name: string
  status: number
  created_at: string
}

export interface AdminUserListParams {
  page: number
  page_size: number
  keyword?: string
  status?: number
}

export interface AdminUserListResponse {
  list: AdminUser[]
  total: number
}

export interface CreateAdminUserPayload {
  nick_name: string
  password: string
  role_ids: number[]
  status: number
}

export interface UpdateAdminUserPayload {
  nick_name: string
  role_ids: number[]
  status: number
}
```

---

# 阶段四：代码规范阶段

## 目标

让 AI 写出来的前端代码结构统一，不要全部堆在一个页面文件里。

## 需要生成的文档

建议创建：

```text
docs/frontend-rules.md
```

也可以和本文件合并。

## 代码组织规则

推荐规则：

```text
页面入口：app/xxx/page.tsx
接口请求：services/xxxService.ts
类型定义：types/xxx.ts
页面组件：components/xxx/
公共组件：components/common/
工具函数：lib/
```

## 组件拆分规则

一个后台管理页面建议拆成：

```text
page.tsx
XxxSearchForm.tsx
XxxTable.tsx
XxxModal.tsx
```

例如：

```text
app/admin/admin-users/page.tsx
components/admin-users/AdminUserSearchForm.tsx
components/admin-users/AdminUserTable.tsx
components/admin-users/AdminUserModal.tsx
services/adminUserService.ts
types/adminUser.ts
```

## 组件职责

### page.tsx

负责：

- 管理页面状态
- 调用接口
- 处理查询、分页、新增、编辑、删除
- 组合子组件

### SearchForm

负责：

- 展示查询条件
- 点击查询时把条件传给父组件
- 点击重置时通知父组件

### Table

负责：

- 展示表格数据
- 显示 loading
- 触发编辑、删除等操作

### Modal

负责：

- 展示新增/编辑表单
- 表单校验
- 提交数据给父组件

## 这个阶段需要对 AI 说什么

```text
请按照以下规则生成代码：
1. 页面入口放到 app 目录。
2. 接口封装放到 services 目录。
3. 类型定义放到 types 目录。
4. 页面组件放到 components 对应业务目录。
5. 不要把所有代码写在一个 page.tsx 文件里。
6. 不要在子组件里重复请求接口。
```

---

# 阶段五：检查验证阶段

## 目标

让 AI 写完代码之后，不只是“看起来有代码”，还要能运行、能编译、能通过检查。

## 需要配置的命令

建议在 `package.json` 中有这些命令：

```json
{
  "scripts": {
    "dev": "next dev",
    "lint": "next lint",
    "type-check": "tsc --noEmit",
    "build": "next build"
  }
}
```

## AI 写完后需要检查什么

让 AI 写完后说明：

```text
请告诉我：
1. 修改了哪些文件？
2. 新增了哪些文件？
3. 需要执行哪些命令？
4. 如果命令报错，应该优先检查哪里？
```

## 推荐检查命令

```bash
npm run lint
npm run type-check
npm run build
```

## 这个阶段需要对 AI 说什么

```text
代码生成完成后，请检查：
1. TypeScript 类型是否完整。
2. 是否存在未使用变量。
3. 是否存在错误的接口字段。
4. 是否存在错误的导入路径。
5. 是否可以通过 npm run lint、npm run type-check、npm run build。
```

---

# 最终推荐 Prompt 模板

以后你可以直接复制下面这段给 AI：

```text
你是一个前端 AI Coding 助手。

请根据以下文档生成前端代码：

1. AGENTS.md
2. docs/pages.md
3. docs/api.md
4. docs/frontend-guide.md

开发要求：

1. 使用当前项目已有技术栈。
2. 页面入口放到 app 目录。
3. 接口请求封装到 services 目录。
4. TypeScript 类型定义放到 types 目录。
5. 页面组件放到 components 对应业务目录。
6. 不要在页面里直接写 axios。
7. 不要把所有代码都写在一个文件里。
8. 不要修改无关文件。
9. 写完后说明修改了哪些文件。
10. 写完后告诉我需要执行哪些命令验证。

现在请根据 docs/pages.md 中的页面说明和 docs/api.md 中的接口说明，生成对应页面。
```

---

# 你作为后端开发者应该怎么配合 AI 写前端

你主要负责三件事：

1. 写清楚接口文档
2. 写清楚页面功能
3. 让 AI 根据文档生成前端

最小可行文档：

```text
AGENTS.md
docs/api.md
docs/pages.md
```

先不用搞太复杂。

只要这三个文档写清楚，AI 写前端就会比你直接一句“帮我写页面”稳定很多。
