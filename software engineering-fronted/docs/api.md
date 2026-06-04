# API 文档

## 基础信息

- **Base URL**: `http://localhost:8080`
- **API 版本**: v1
- **Content-Type**: `application/json`
- **文件上传**: `multipart/form-data`

## 统一响应格式

所有 API 返回统一的响应格式：

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

- `code`: 业务状态码，`200` 表示成功
- `message`: 响应消息
- `data`: 响应数据

## 错误码

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |
| 1001 | 用户不存在 |
| 1002 | 用户已存在 |
| 1003 | 用户名或密码错误 |
| 1004 | 用户已被禁用 |
| 1005 | 无效的令牌 |
| 1006 | 令牌已过期 |
| 2001 | 文档不存在 |
| 2002 | 文档解析失败 |
| 3001 | 知识点不存在 |
| 3002 | 关系不存在 |
| 4001 | 题目不存在 |
| 4002 | 答案格式错误 |

## 认证方式

使用 JWT Bearer Token 认证：

```
Authorization: Bearer {token}
```

---

## API 接口列表

### 健康检查

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/health | 检查服务状态 |

### 认证模块 `/api/v1/auth`

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /auth/register | 用户注册 |
| POST | /auth/login | 用户登录 |
| POST | /auth/refresh | 刷新token |

### 用户管理 `/api/v1/users`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /users/profile | 获取当前用户信息 | 是 |
| PUT | /users/profile | 更新用户信息 | 是 |
| POST | /users/password | 修改密码 | 是 |
| GET | /users/list | 用户列表（分页） | 是 |

### 资源管理 `/api/v1/documents`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /documents | 上传文档 | 是 |
| GET | /documents | 文档列表 | 是 |
| GET | /documents/:id | 文档详情 | 是 |
| PUT | /documents/:id | 更新文档信息 | 是 |
| DELETE | /documents/:id | 删除文档 | 是 |
| GET | /documents/:id/content | 获取文档内容 | 是 |

### 知识点管理 `/api/v1/knowledge`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /knowledge/points | 知识点列表 | 是 |
| GET | /knowledge/points/:id | 知识点详情 | 是 |
| POST | /knowledge/points | 新增知识点 | 是 |
| PUT | /knowledge/points/:id | 更新知识点 | 是 |
| DELETE | /knowledge/points/:id | 删除知识点 | 是 |
| GET | /knowledge/relations | 关系列表 | 是 |
| POST | /knowledge/relations | 新增关系 | 是 |
| PUT | /knowledge/relations/:id | 更新关系 | 是 |
| DELETE | /knowledge/relations/:id | 删除关系 | 是 |

### 知识图谱 `/api/v1/graph`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /graph | 获取图谱数据 | 是 |
| POST | /graph/build | 从文档构建图谱 | 是 |
| GET | /graph/build/latest | 最近构建结果 | 是 |
| GET | /graph/build/history | 构建历史记录 | 是 |

### 题库 `/api/v1/questions`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /questions | 题目列表（分页） | 是 |
| GET | /questions/:id | 题目详情 | 是 |
| POST | /questions | 新增题目 | 是 |
| PUT | /questions/:id | 更新题目 | 是 |
| DELETE | /questions/:id | 删除题目 | 是 |

### 答题 `/api/v1/quizzes`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /quizzes/submit | 提交答题 | 是 |
| GET | /quizzes/history | 答题历史 | 是 |
| GET | /quizzes/:id | 答题详情 | 是 |

### 知识问答 `/api/v1/ask`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /ask/sessions | 新建问答会话 | 是 |
| GET | /ask/sessions | 会话列表 | 是 |
| GET | /ask/sessions/:id/messages | 会话消息列表 | 是 |
| POST | /ask | 提问 | 是 |
| GET | /ask/history | 问答历史 | 是 |

### 学习分析 `/api/v1/analytics`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /analytics/overview | 总览统计 | 是 |
| GET | /analytics/hot-knowledge-points | 热门知识点 | 是 |
| GET | /analytics/knowledge-mastery | 知识点掌握度 | 是 |
| GET | /analytics/weak-points | 薄弱知识点 | 是 |
| GET | /analytics/trends | 趋势数据 | 是 |

---

## 接口详细说明

### 健康检查

#### 检查服务状态

```http
GET /api/v1/health
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "status": "ok",
    "service": "software-engineering-backend"
  }
}
```

---

### 认证模块

#### 用户注册

```http
POST /api/v1/auth/register
```

**请求参数**:
```json
{
  "username": "student001",
  "password": "123456",
  "email": "student001@example.com",
  "nickname": "张三"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名（3-50字符） |
| password | string | 是 | 密码（6-50字符） |
| email | string | 是 | 邮箱地址 |
| nickname | string | 否 | 昵称（最多50字符） |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": null
}
```

#### 用户登录

```http
POST /api/v1/auth/login
```

**请求参数**:
```json
{
  "username": "student001",
  "password": "123456"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "student001",
      "email": "student001@example.com",
      "nickname": "张三",
      "avatar": "",
      "status": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

#### 刷新Token

```http
POST /api/v1/auth/refresh
```

**请求参数**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

---

### 用户管理

以下接口需要认证，请在请求头中携带 Token。

#### 获取当前用户信息

```http
GET /api/v1/users/profile
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "id": 1,
    "username": "student001",
    "email": "student001@example.com",
    "nickname": "张三",
    "avatar": "",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 更新用户信息

```http
PUT /api/v1/users/profile
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "nickname": "新昵称",
  "avatar": "https://example.com/avatar.jpg"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| nickname | string | 否 | 昵称 |
| avatar | string | 否 | 头像 URL |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": null
}
```

#### 修改密码

```http
POST /api/v1/users/password
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "old_password": "123456",
  "new_password": "654321"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| old_password | string | 是 | 旧密码 |
| new_password | string | 是 | 新密码（6-50字符） |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": null
}
```

#### 获取用户列表（分页）

```http
GET /api/v1/users/list?page=1&size=10
Authorization: Bearer {token}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 是 | 页码，从 1 开始 |
| size | int | 是 | 每页大小，1-100 |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "student001",
        "email": "student001@example.com",
        "nickname": "张三",
        "avatar": "",
        "status": 1,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "size": 10,
    "total_page": 10
  }
}
```

---

### 资源管理

#### 上传文档

```http
POST /api/v1/documents
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | 文档文件（支持 .md, .txt, .pdf, .docx） |
| title | string | 否 | 文档标题（不填则使用文件名） |
| description | string | 否 | 文档描述 |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "id": 1,
    "title": "软件工程课堂笔记",
    "description": "第三章需求分析",
    "filename": "notes.md",
    "file_size": 1024,
    "file_type": ".md",
    "status": "pending",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 文档列表

```http
GET /api/v1/documents?page=1&size=10&keyword=软件
Authorization: Bearer {token}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 是 | 页码 |
| size | int | 是 | 每页大小 |
| keyword | string | 否 | 搜索关键词 |
| status | string | 否 | 状态筛选：pending/processing/completed/failed |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "软件工程课堂笔记",
        "description": "第三章需求分析",
        "filename": "notes.md",
        "file_size": 1024,
        "file_type": ".md",
        "status": "completed",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 25,
    "page": 1,
    "size": 10,
    "total_page": 3
  }
}
```

#### 文档详情

```http
GET /api/v1/documents/:id
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "id": 1,
    "title": "软件工程课堂笔记",
    "description": "第三章需求分析",
    "filename": "notes.md",
    "file_size": 1024,
    "file_type": ".md",
    "status": "completed",
    "content_preview": "需求分析是软件工程的第一阶段...",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 更新文档信息

```http
PUT /api/v1/documents/:id
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "title": "新标题",
  "description": "新描述"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": null
}
```

#### 删除文档

```http
DELETE /api/v1/documents/:id
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": null
}
```

#### 获取文档内容

```http
GET /api/v1/documents/:id/content
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "id": 1,
    "title": "软件工程课堂笔记",
    "content": "# 第三章 需求分析\n\n需求分析是软件工程..."
  }
}
```

---

### 知识点管理

#### 知识点列表

```http
GET /api/v1/knowledge/points?page=1&size=10&keyword=需求
Authorization: Bearer {token}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 是 | 页码 |
| size | int | 是 | 每页大小 |
| keyword | string | 否 | 搜索关键词 |
| document_id | int | 否 | 按文档筛选 |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "需求分析",
        "description": "识别和确认用户需求的过程",
        "document_id": 1,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "size": 10,
    "total_page": 5
  }
}
```

#### 知识点详情

```http
GET /api/v1/knowledge/points/:id
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "id": 1,
    "name": "需求分析",
    "description": "识别和确认用户需求的过程",
    "document_id": 1,
    "document_title": "软件工程课堂笔记",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 新增知识点

```http
POST /api/v1/knowledge/points
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "name": "软件测试",
  "description": "验证软件是否满足需求的过程",
  "document_id": 1
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 知识点名称 |
| description | string | 否 | 知识点描述 |
| document_id | int | 否 | 关联文档ID |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "id": 2
  }
}
```

#### 更新知识点

```http
PUT /api/v1/knowledge/points/:id
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "name": "软件测试（更新）",
  "description": "验证软件是否满足需求的过程（更新）"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": null
}
```

#### 删除知识点

```http
DELETE /api/v1/knowledge/points/:id
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": null
}
```

#### 关系列表

```http
GET /api/v1/knowledge/relations?page=1&size=10&point_id=1
Authorization: Bearer {token}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 是 | 页码 |
| size | int | 是 | 每页大小 |
| point_id | int | 否 | 按知识点筛选（来源或目标） |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "list": [
      {
        "id": 1,
        "source_id": 1,
        "source_name": "需求分析",
        "target_id": 2,
        "target_name": "软件测试",
        "relation_type": "RELATED",
        "description": "需求分析是软件测试的前置环节",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 30,
    "page": 1,
    "size": 10,
    "total_page": 3
  }
}
```

#### 新增关系

```http
POST /api/v1/knowledge/relations
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "source_id": 1,
  "target_id": 2,
  "relation_type": "RELATED",
  "description": "需求分析是软件测试的前置环节"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| source_id | int | 是 | 来源知识点ID |
| target_id | int | 是 | 目标知识点ID |
| relation_type | string | 是 | 关系类型：RELATED/DEPENDS_ON/PART_OF |
| description | string | 否 | 关系描述 |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "id": 1
  }
}
```

#### 更新关系

```http
PUT /api/v1/knowledge/relations/:id
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "relation_type": "DEPENDS_ON",
  "description": "软件测试依赖于需求分析的结果"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": null
}
```

#### 删除关系

```http
DELETE /api/v1/knowledge/relations/:id
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": null
}
```

---

### 知识图谱

#### 获取图谱数据

```http
GET /api/v1/graph?document_id=1&keyword=测试&relation_type=DEPENDS_ON
Authorization: Bearer {token}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| document_id | int | 否 | 按文档筛选 |
| keyword | string | 否 | 按知识点名称模糊搜索 |
| relation_type | string | 否 | 按关系类型筛选：RELATED/DEPENDS_ON/PART_OF |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "nodes": [
      {
        "id": 1,
        "name": "需求分析",
        "description": "识别和确认用户需求的过程",
        "document_id": 1,
        "category": "需求相关"
      },
      {
        "id": 2,
        "name": "软件测试",
        "description": "验证软件是否满足需求的过程",
        "document_id": 1,
        "category": "测试相关"
      }
    ],
    "edges": [
      {
        "id": 1,
        "source": 1,
        "target": 2,
        "relation_type": "RELATED",
        "description": "需求分析是软件测试的前置环节"
      }
    ],
    "summary": {
      "node_count": 2,
      "edge_count": 1
    }
  }
}
```

#### 从文档构建图谱

```http
POST /api/v1/graph/build
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "document_ids": [1, 2]
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| document_ids | array | 是 | 文档ID列表 |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "build_id": 1,
    "created_points": 15,
    "created_relations": 23,
    "chunk_count": 8,
    "vector_count": 42,
    "status": "completed",
    "message": "知识图谱构建完成"
  }
}
```

#### 获取最近构建结果

```http
GET /api/v1/graph/build/latest
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "build_id": 1,
    "document_ids": [1, 2],
    "created_points": 15,
    "created_relations": 23,
    "chunk_count": 8,
    "vector_count": 42,
    "status": "completed",
    "message": "知识图谱构建完成",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 构建历史记录

```http
GET /api/v1/graph/build/history?page=1&size=10
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "list": [
      {
        "build_id": 1,
        "document_ids": [1, 2],
        "created_points": 15,
        "created_relations": 23,
        "status": "completed",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "size": 10,
    "total_page": 1
  }
}
```

---

### 题库

#### 题目列表

```http
GET /api/v1/questions?page=1&size=10&keyword=需求
Authorization: Bearer {token}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 是 | 页码 |
| size | int | 是 | 每页大小 |
| keyword | string | 否 | 搜索关键词 |
| knowledge_point_id | int | 否 | 按知识点筛选 |
| difficulty | string | 否 | 难度：easy/medium/hard |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "以下哪个不是需求分析的活动？",
        "type": "single",
        "difficulty": "easy",
        "knowledge_point_id": 1,
        "knowledge_point_name": "需求分析",
        "options": [
          {"key": "A", "value": "需求获取"},
          {"key": "B", "value": "需求分析"},
          {"key": "C", "value": "代码编写"},
          {"key": "D", "value": "需求验证"}
        ],
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "size": 10,
    "total_page": 10
  }
}
```

#### 题目详情

```http
GET /api/v1/questions/:id
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "id": 1,
    "title": "以下哪个不是需求分析的活动？",
    "type": "single",
    "difficulty": "easy",
    "knowledge_point_id": 1,
    "knowledge_point_name": "需求分析",
    "options": [
      {"key": "A", "value": "需求获取"},
      {"key": "B", "value": "需求分析"},
      {"key": "C", "value": "代码编写"},
      {"key": "D", "value": "需求验证"}
    ],
    "answer": "C",
    "explanation": "代码编写属于编码阶段，不是需求分析的活动",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 新增题目

```http
POST /api/v1/questions
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "title": "以下哪个不是需求分析的活动？",
  "type": "single",
  "difficulty": "easy",
  "knowledge_point_id": 1,
  "options": [
    {"key": "A", "value": "需求获取"},
    {"key": "B", "value": "需求分析"},
    {"key": "C", "value": "代码编写"},
    {"key": "D", "value": "需求验证"}
  ],
  "answer": "C",
  "explanation": "代码编写属于编码阶段，不是需求分析的活动"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 题目标题 |
| type | string | 是 | 题目类型：single/multiple |
| difficulty | string | 是 | 难度：easy/medium/hard |
| knowledge_point_id | int | 是 | 关联知识点ID |
| options | array | 是 | 选项列表 |
| answer | string | 是 | 正确答案（多选用逗号分隔，如"A,B"） |
| explanation | string | 否 | 答案解析 |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "id": 1
  }
}
```

#### 更新题目

```http
PUT /api/v1/questions/:id
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "title": "更新后的题目",
  "answer": "D"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": null
}
```

#### 删除题目

```http
DELETE /api/v1/questions/:id
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": null
}
```

---

### 答题

#### 提交答题

```http
POST /api/v1/quizzes/submit
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "question_id": 1,
  "user_answer": "C"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| question_id | int | 是 | 题目ID |
| user_answer | string | 是 | 用户答案（多选用逗号分隔） |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "quiz_id": 1,
    "question_id": 1,
    "user_answer": "C",
    "correct_answer": "C",
    "is_correct": true,
    "explanation": "代码编写属于编码阶段，不是需求分析的活动",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 答题历史

```http
GET /api/v1/quizzes/history?page=1&size=10&knowledge_point_id=1
Authorization: Bearer {token}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 是 | 页码 |
| size | int | 是 | 每页大小 |
| knowledge_point_id | int | 否 | 按知识点筛选 |
| is_correct | bool | 否 | 按正确/错误筛选 |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "list": [
      {
        "quiz_id": 1,
        "question_id": 1,
        "question_title": "以下哪个不是需求分析的活动？",
        "user_answer": "C",
        "correct_answer": "C",
        "is_correct": true,
        "knowledge_point_id": 1,
        "knowledge_point_name": "需求分析",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "size": 10,
    "total_page": 5
  }
}
```

#### 答题详情

```http
GET /api/v1/quizzes/:id
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "quiz_id": 1,
    "question": {
      "id": 1,
      "title": "以下哪个不是需求分析的活动？",
      "type": "single",
      "difficulty": "easy",
      "options": [
        {"key": "A", "value": "需求获取"},
        {"key": "B", "value": "需求分析"},
        {"key": "C", "value": "代码编写"},
        {"key": "D", "value": "需求验证"}
      ],
      "answer": "C",
      "explanation": "代码编写属于编码阶段，不是需求分析的活动"
    },
    "user_answer": "C",
    "is_correct": true,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

### 知识问答

#### 新建问答会话

```http
POST /api/v1/ask/sessions
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "title": "黑盒测试和白盒测试的区别"
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 否 | 会话标题；不传则由首个问题自动生成 |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "conversation_id": 12,
    "title": "黑盒测试和白盒测试的区别",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 会话列表

```http
GET /api/v1/ask/sessions?page=1&size=10
Authorization: Bearer {token}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 是 | 页码 |
| size | int | 是 | 每页大小 |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "list": [
      {
        "conversation_id": 12,
        "title": "黑盒测试和白盒测试的区别",
        "last_question": "黑盒测试和白盒测试的区别是什么？",
        "message_count": 4,
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 8,
    "page": 1,
    "size": 10,
    "total_page": 1
  }
}
```

#### 会话消息列表

```http
GET /api/v1/ask/sessions/12/messages?page=1&size=20
Authorization: Bearer {token}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 是 | 页码 |
| size | int | 是 | 每页大小 |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "conversation_id": 12,
    "title": "黑盒测试和白盒测试的区别",
    "list": [
      {
        "message_id": 101,
        "role": "user",
        "content": "黑盒测试和白盒测试的区别是什么？",
        "created_at": "2024-01-01T00:00:00Z"
      },
      {
        "message_id": 102,
        "role": "assistant",
        "content": "黑盒测试和白盒测试都是软件测试的方法，两者的主要区别如下：",
        "created_at": "2024-01-01T00:00:05Z"
      }
    ],
    "total": 2,
    "page": 1,
    "size": 20,
    "total_page": 1
  }
}
```

#### 提问

```http
POST /api/v1/ask
Authorization: Bearer {token}
```

**请求参数**:
```json
{
  "question": "什么是需求分析？",
  "conversation_id": 12
}
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| question | string | 是 | 用户问题 |
| conversation_id | int | 否 | 会话 ID；不传则自动创建新会话 |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "conversation_id": 12,
    "question_id": 35,
    "answer": "需求分析是软件工程的第一阶段，用于识别和确认用户需求的过程。主要包括需求获取、需求分析、需求规格说明和需求验证等活动。",
    "confidence": 0.85,
    "sources": [
      {
        "document_id": 1,
        "document_title": "软件工程课堂笔记",
        "content": "需求分析是软件工程的第一阶段..."
      }
    ],
    "related_knowledge_points": [
      {
        "id": 1,
        "name": "需求分析",
        "description": "识别和确认用户需求的过程"
      },
      {
        "id": 3,
        "name": "需求获取",
        "description": "从用户处收集需求"
      }
    ],
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 问答历史

```http
GET /api/v1/ask/history?page=1&size=10&conversation_id=12
Authorization: Bearer {token}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 是 | 页码 |
| size | int | 是 | 每页大小 |
| conversation_id | int | 否 | 按会话筛选 |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "list": [
      {
        "id": 1,
        "conversation_id": 12,
        "conversation_title": "黑盒测试和白盒测试的区别",
        "question": "什么是需求分析？",
        "answer": "需求分析是软件工程的第一阶段...",
        "confidence": 0.85,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 30,
    "page": 1,
    "size": 10,
    "total_page": 3
  }
}
```

---

### 学习分析

#### 总览统计

```http
GET /api/v1/analytics/overview
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "today_learning_hours": 1.5,
    "today_questions_asked": 12,
    "total_learning_hours": 48.5,
    "total_questions_asked": 45,
    "total_quizzes_taken": 120,
    "average_correct_rate": 0.785,
    "knowledge_points_mastered": 15,
    "knowledge_points_total": 25,
    "mastery_rate": 0.6
  }
}
```

#### 热门知识点

```http
GET /api/v1/analytics/hot-knowledge-points?limit=10
Authorization: Bearer {token}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| limit | int | 否 | 返回数量，默认10 |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": [
    {
      "knowledge_point_id": 1,
      "knowledge_point_name": "软件生命周期",
      "heat": 352,
      "question_count": 120,
      "quiz_count": 88
    },
    {
      "knowledge_point_id": 2,
      "knowledge_point_name": "需求分析",
      "heat": 298,
      "question_count": 96,
      "quiz_count": 74
    }
  ]
}
```

#### 知识点掌握度

```http
GET /api/v1/analytics/knowledge-mastery
Authorization: Bearer {token}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": [
    {
      "knowledge_point_id": 1,
      "knowledge_point_name": "需求分析",
      "total_questions": 10,
      "correct_answers": 8,
      "mastery_rate": 0.8,
      "level": "mastered"
    },
    {
      "knowledge_point_id": 2,
      "knowledge_point_name": "软件测试",
      "total_questions": 8,
      "correct_answers": 4,
      "mastery_rate": 0.5,
      "level": "learning"
    }
  ]
}
```

掌握度等级说明：
- `mastered`: 掌握（>=80%）
- `learning`: 学习中（50%-80%）
- `weak`: 薄弱（<50%）

#### 薄弱知识点

```http
GET /api/v1/analytics/weak-points?limit=10
Authorization: Bearer {token}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| limit | int | 否 | 返回数量，默认10 |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": [
    {
      "knowledge_point_id": 2,
      "knowledge_point_name": "软件测试",
      "correct_rate": 0.5,
      "suggested_questions": [
        {
          "id": 20,
          "title": "以下哪个是黑盒测试方法？"
        },
        {
          "id": 21,
          "title": "测试用例设计的原则是什么？"
        }
      ]
    }
  ]
}
```

#### 趋势数据

```http
GET /api/v1/analytics/trends?days=7
Authorization: Bearer {token}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| days | int | 否 | 统计天数，默认7 |

**响应示例**:
```json
{
  "code": 200,
  "message": "成功",
  "data": {
    "daily_stats": [
      {
        "date": "2024-01-01",
        "questions_asked": 5,
        "learning_hours": 1.2,
        "correct_rate": 0.8
      },
      {
        "date": "2024-01-02",
        "questions_asked": 3,
        "learning_hours": 0.8,
        "correct_rate": 0.75
      }
    ],
    "weekly_trend": [
      {
        "week": "2024-W01",
        "avg_correct_rate": 0.78,
        "total_learning_hours": 8.6,
        "total_questions_asked": 31
      }
    ]
  }
}
```

---

## 使用示例

### cURL

```bash
# 用户注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"student001","password":"123456","email":"student001@example.com"}'

# 用户登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"student001","password":"123456"}'

# 上传文档
curl -X POST http://localhost:8080/api/v1/documents \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@notes.md" \
  -F "title=软件工程笔记"

# 提问
curl -X POST http://localhost:8080/api/v1/ask \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"question":"什么是需求分析？"}'

# 新建问答会话并在会话中提问
curl -X POST http://localhost:8080/api/v1/ask/sessions \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"黑盒测试和白盒测试的区别"}'

curl -X POST http://localhost:8080/api/v1/ask \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"question":"黑盒测试和白盒测试的区别是什么？","conversation_id":12}'

# 构建知识图谱
curl -X POST http://localhost:8080/api/v1/graph/build \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"document_ids":[1,2]}'

# 按关键词和关系类型筛选图谱
curl "http://localhost:8080/api/v1/graph?keyword=测试&relation_type=DEPENDS_ON" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 获取学习分析
curl http://localhost:8080/api/v1/analytics/overview \
  -H "Authorization: Bearer YOUR_TOKEN"

# 获取热门知识点
curl "http://localhost:8080/api/v1/analytics/hot-knowledge-points?limit=5" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### JavaScript (Fetch)

```javascript
// 用户登录
const loginRes = await fetch('http://localhost:8080/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ username: 'student001', password: '123456' })
});
const { data: { token } } = await loginRes.json();

// 提问
const askRes = await fetch('http://localhost:8080/api/v1/ask', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({ question: '什么是需求分析？' })
});
const { data: { answer, confidence } } = await askRes.json();
console.log(`回答：${answer}（置信度：${confidence}）`);
```

---

## 注意事项

1. **Token 有效期**: Token 默认有效期为 24 小时
2. **密码安全**: 密码使用 bcrypt 加密存储，服务端无法查看明文密码
3. **文件上传限制**: 单文件最大 10MB
4. **分页参数**: page 从 1 开始，size 范围 1-100
5. **时间格式**: 所有时间使用 UTC 时区，格式为 ISO 8601

