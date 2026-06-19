# SE智图问答前端自动化测试

## 概述

本项目使用Playwright进行前端自动化测试，覆盖了所有主要页面和功能。

## 测试结构

```
tests/
├── login.spec.ts           # 登录页面测试
├── register.spec.ts        # 注册页面测试
├── home.spec.ts            # 首页测试
├── knowledge-graph.spec.ts # 知识图谱页面测试
├── qa.spec.ts              # 问答中心测试
├── files.spec.ts           # 资料管理页面测试
├── stats.spec.ts           # 分析统计页面测试
├── navigation.spec.ts      # 导航功能测试
├── common.spec.ts          # 通用功能测试
└── README.md               # 本文档
```

## 测试覆盖范围

### 1. 登录页面测试 (login.spec.ts)
- ✅ 页面加载正确
- ✅ 表单元素完整
- ✅ 表单验证 - 空值提交
- ✅ 表单验证 - 仅填写用户名
- ✅ 表单验证 - 仅填写密码
- ✅ 导航到注册页面
- ✅ 密码显示/隐藏功能
- ✅ 页面响应式布局

### 2. 注册页面测试 (register.spec.ts)
- ✅ 页面加载正确
- ✅ 表单元素完整
- ✅ 表单验证 - 空值提交
- ✅ 密码确认验证
- ✅ 导航到登录页面
- ✅ 页面响应式布局

### 3. 首页测试 (home.spec.ts)
- ✅ 页面加载正确
- ✅ 侧边栏导航完整
- ✅ 用户信息显示
- ✅ 首页内容加载
- ✅ 导航到各个页面
- ✅ 退出登录功能
- ✅ 页面响应式布局

### 4. 知识图谱页面测试 (knowledge-graph.spec.ts)
- ✅ 页面加载正确
- ✅ 图谱工具栏完整
- ✅ 图谱画布存在
- ✅ 节点详情面板
- ✅ 构建图谱对话框
- ✅ 图谱搜索功能
- ✅ 图谱交互功能
- ✅ 页面响应式布局

### 5. 问答中心测试 (qa.spec.ts)
- ✅ 页面加载正确
- ✅ 会话列表存在
- ✅ 聊天面板存在
- ✅ 空状态显示
- ✅ 新建会话功能
- ✅ 发送消息功能
- ✅ 会话切换功能
- ✅ 消息气泡显示
- ✅ 来源参考功能
- ✅ 页面响应式布局

### 6. 资料管理页面测试 (files.spec.ts)
- ✅ 页面加载正确
- ✅ 文档表格存在
- ✅ 上传按钮存在
- ✅ 上传对话框功能
- ✅ 文档详情功能
- ✅ 文档操作按钮
- ✅ 分页功能
- ✅ 搜索功能
- ✅ 页面响应式布局

### 7. 分析统计页面测试 (stats.spec.ts)
- ✅ 页面加载正确
- ✅ 趋势图表存在
- ✅ 掌握度图表存在
- ✅ 热点排行存在
- ✅ 薄弱点列表存在
- ✅ 图表交互功能
- ✅ 时间范围选择
- ✅ 数据刷新功能
- ✅ 页面响应式布局

### 8. 导航功能测试 (navigation.spec.ts)
- ✅ 侧边栏导航高亮
- ✅ 路由守卫功能
- ✅ 页面标题更新
- ✅ Logo点击返回首页
- ✅ 浏览器前进后退

### 9. 通用功能测试 (common.spec.ts)
- ✅ 页面加载性能
- ✅ 内存泄漏检测
- ✅ 控制台错误检测
- ✅ 响应式断点测试
- ✅ 网络请求监控
- ✅ 本地存储功能
- ✅ Cookie功能

## 运行测试

### 安装依赖

```bash
cd software engineering-fronted/frontend-vue-app
npm install
npx playwright install
```

### 运行所有测试

```bash
npx playwright test
```

### 运行特定测试文件

```bash
npx playwright test tests/login.spec.ts
npx playwright test tests/home.spec.ts
```

### 运行特定测试用例

```bash
npx playwright test -g "页面加载正确"
```

### 生成测试报告

```bash
npx playwright show-report
```

### 调试模式

```bash
npx playwright test --headed
npx playwright test --debug
```

## 测试配置

测试配置文件为 `playwright.config.ts`，包含以下配置：

- **测试目录**: `./tests`
- **基础URL**: `http://localhost:5173`
- **浏览器**: Chromium, Firefox, WebKit
- **重试次数**: CI环境2次，本地0次
- **截图**: 失败时自动截图
- **跟踪**: 首次重试时记录跟踪

## 测试数据

测试使用模拟数据，不需要真实的后端服务：

- **登录状态**: 通过localStorage模拟已登录状态
- **用户信息**: 使用固定的测试用户数据
- **API请求**: 使用mock数据或跳过API调用

## 注意事项

1. **测试环境**: 测试应该在独立的测试环境中运行，避免影响生产数据
2. **测试数据**: 每次测试前会自动清理测试数据
3. **并行测试**: 支持并行测试，但需要注意测试数据隔离
4. **CI/CD**: 测试可以在CI/CD管道中自动运行

## 故障排查

### 测试失败

1. 检查前端服务器是否正常运行
2. 检查测试环境配置
3. 查看测试报告和截图
4. 检查控制台错误信息

### 性能问题

1. 增加测试超时时间
2. 优化测试数据准备
3. 使用并行测试提高效率

## 贡献指南

1. 添加新测试时，请遵循现有的测试结构
2. 确保测试用例具有描述性的名称
3. 添加适当的注释说明测试目的
4. 运行所有测试确保没有破坏现有功能

## 相关资源

- [Playwright官方文档](https://playwright.dev/)
- [Vue Testing Handbook](https://vue-testing-handbook.dev/)
- [Element Plus测试指南](https://element-plus.org/zh-CN/guide/quick-start.html)
