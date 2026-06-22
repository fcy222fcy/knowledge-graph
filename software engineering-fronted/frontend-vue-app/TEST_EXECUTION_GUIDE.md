# SE智图问答系统 - 测试执行指南

## 一、快速开始

### 1.1 环境要求

- **Node.js**: v18.0.0 或更高版本
- **npm**: v9.0.0 或更高版本
- **操作系统**: Windows 10/11, macOS 12+, Ubuntu 20.04+

### 1.2 安装步骤

```bash
# 1. 进入前端项目目录
cd software engineering-fronted/frontend-vue-app

# 2. 安装项目依赖
npm install

# 3. 安装 Playwright 浏览器
npx playwright install

# 4. 验证安装
npx playwright --version
```

---

## 二、运行测试

### 2.1 运行所有测试

```bash
# 运行所有测试（无头模式）
npm test

# 或者直接使用 playwright 命令
npx playwright test
```

### 2.2 运行特定测试

```bash
# 运行登录测试
npx playwright test tests/login.spec.ts

# 运行首页测试
npx playwright test tests/home.spec.ts

# 运行问答中心测试
npx playwright test tests/qa.spec.ts

# 运行知识图谱测试
npx playwright test tests/knowledge-graph.spec.ts

# 运行资料管理测试
npx playwright test tests/files.spec.ts

# 运行分析统计测试
npx playwright test tests/stats.spec.ts

# 运行导航测试
npx playwright test tests/navigation.spec.ts

# 运行通用功能测试
npx playwright test tests/common.spec.ts
```

### 2.3 调试模式

```bash
# 有头模式运行（可以看到浏览器操作）
npm run test:headed

# 调试模式运行（逐步执行）
npm run test:debug

# 调试特定测试
npx playwright test tests/login.spec.ts --debug
```

### 2.4 查看测试报告

```bash
# 生成并查看HTML报告
npm run test:report

# 或者直接打开报告
npx playwright show-report
```

---

## 三、测试配置

### 3.1 Playwright配置文件

配置文件位于 `playwright.config.ts`，主要配置项：

```typescript
import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  // 测试目录
  testDir: './tests',

  // 并行执行
  fullyParallel: true,

  // CI环境禁止.only
  forbidOnly: !!process.env.CI,

  // 重试次数
  retries: process.env.CI ? 2 : 0,

  // 并发数
  workers: process.env.CI ? 1 : undefined,

  // 报告器
  reporter: 'html',

  // 全局配置
  use: {
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },

  // 浏览器项目
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
    // 移动端测试
    {
      name: 'mobile-chrome',
      use: { ...devices['Pixel 5'] },
    },
    {
      name: 'mobile-safari',
      use: { ...devices['iPhone 12'] },
    },
  ],

  // Web服务器配置
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: !process.env.CI,
  },
})
```

---

## 四、测试数据管理

### 4.1 Mock数据设置

测试使用 `tests/helpers.ts` 中的 `setupLoggedIn` 函数来设置登录状态：

```typescript
import type { Page } from '@playwright/test'

export async function setupLoggedIn(page: Page) {
  // 在页面加载前设置 localStorage
  await page.addInitScript(() => {
    localStorage.setItem('token', 'test-token')
    localStorage.setItem('userInfo', JSON.stringify({
      id: 1,
      username: 'testuser',
      nickname: '测试用户',
      email: 'test@example.com',
      avatar: '',
      status: 1,
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z'
    }))
  })

  // Mock 所有 API 请求
  await page.route('**/api/**', (route) => {
    const url = route.request().url()

    // 根据URL返回不同的Mock数据
    if (url.includes('/analytics/overview')) {
      return route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          code: 200,
          data: {
            total_learning_hours: 12.5,
            total_questions_asked: 48,
            total_quizzes_taken: 5,
            average_correct_rate: 78.5,
            knowledge_points_mastered: 15,
            knowledge_points_total: 30
          },
          message: 'success'
        })
      })
    }

    // 默认：返回空成功响应
    return route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        code: 200,
        data: null,
        message: 'success'
      })
    })
  })
}
```

### 4.2 Mock数据API端点

| API端点 | 描述 | Mock数据示例 |
|---------|------|-------------|
| `/api/analytics/overview` | 学习统计概览 | 学习时长、问题数量、测验次数等 |
| `/api/analytics/trends` | 学习趋势数据 | 每日学习时长、正确率等 |
| `/api/ask/history` | 问答历史记录 | 会话列表、最后问题等 |
| `/api/graph` | 知识图谱数据 | 节点、边、摘要信息 |
| `/api/document` | 文档列表数据 | 文档列表、总数 |
| `/api/analytics/knowledge-mastery` | 知识点掌握度 | 掌握度数据 |
| `/api/analytics/weak-points` | 薄弱点 | 薄弱知识点列表 |
| `/api/analytics/hot-points` | 热门知识点 | 热门知识点列表 |

---

## 五、测试用例说明

### 5.1 测试文件结构

```
tests/
├── helpers.ts                    # 测试辅助函数
├── login.spec.ts                 # 登录页面测试
├── register.spec.ts              # 注册页面测试
├── home.spec.ts                  # 首页测试
├── knowledge-graph.spec.ts       # 知识图谱测试
├── qa.spec.ts                    # 问答中心测试
├── files.spec.ts                 # 资料管理测试
├── stats.spec.ts                 # 分析统计测试
├── navigation.spec.ts            # 导航功能测试
└── common.spec.ts                # 通用功能测试
```

### 5.2 测试用例统计

| 测试文件 | 用例数量 | 覆盖模块 |
|----------|----------|----------|
| login.spec.ts | 8 | 用户认证 |
| register.spec.ts | 6 | 用户认证 |
| home.spec.ts | 7 | 首页 |
| knowledge-graph.spec.ts | 8 | 知识图谱 |
| qa.spec.ts | 9 | 问答中心 |
| files.spec.ts | 9 | 资料管理 |
| stats.spec.ts | 8 | 分析统计 |
| navigation.spec.ts | 5 | 导航功能 |
| common.spec.ts | 7 | 通用功能 |
| **总计** | **67** | - |

---

## 六、常见问题

### 6.1 测试失败排查

```bash
# 1. 查看详细错误信息
npx playwright test --reporter=list

# 2. 查看测试截图
npx playwright show-trace test-results/*/trace.zip

# 3. 重试失败的测试
npx playwright test --retries=3
```

### 6.2 浏览器问题

```bash
# 重新安装浏览器
npx playwright install --force

# 安装系统依赖（Linux）
npx playwright install-deps
```

### 6.3 端口冲突

如果默认端口 5173 被占用：

```bash
# 修改 playwright.config.ts 中的端口
webServer: {
  command: 'npm run dev -- --port 3000',
  url: 'http://localhost:3000',
  reuseExistingServer: !process.env.CI,
}
```

---

## 七、持续集成

### 7.1 GitHub Actions

在 `.github/workflows/playwright.yml` 中配置：

```yaml
name: Playwright Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: 18
      - name: Install dependencies
        run: npm ci
      - name: Install Playwright Browsers
        run: npx playwright install --with-deps
      - name: Run Playwright tests
        run: npx playwright test
      - name: Upload test results
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: playwright-report
          path: playwright-report/
          retention-days: 30
```

### 7.2 本地CI运行

```bash
# 模拟CI环境运行
CI=true npm test
```

---

## 八、测试最佳实践

### 8.1 编写测试用例

1. **使用描述性测试名称**: 清晰描述测试目的
2. **保持测试独立性**: 每个测试应该独立运行
3. **使用beforeEach**: 设置测试前置条件
4. **合理使用Mock**: 避免依赖外部服务
5. **添加适当等待**: 使用waitForLoadState等待页面加载

### 8.2 测试维护

1. **定期更新Mock数据**: 保持与实际API一致
2. **清理测试数据**: 避免测试间相互影响
3. **监控测试覆盖率**: 确保关键功能被覆盖
4. **及时修复失败测试**: 保持测试套件健康

### 8.3 性能优化

1. **并行执行**: 利用Playwright的并行能力
2. **减少等待时间**: 使用智能等待而非固定等待
3. **复用浏览器上下文**: 减少浏览器启动开销
4. **优化Mock数据**: 保持Mock数据简洁

---

## 九、测试报告解读

### 9.1 HTML报告

```bash
# 打开HTML报告
npx playwright show-report
```

报告包含：
- **测试结果概览**: 通过/失败/跳过统计
- **详细测试结果**: 每个测试用例的执行结果
- **错误信息**: 失败测试的详细错误堆栈
- **截图和追踪**: 失败测试的截图和执行追踪

### 9.2 测试趋势

定期查看测试报告，关注：
- **测试通过率**: 保持在95%以上
- **测试执行时间**: 控制在合理范围内
- **失败测试分析**: 及时修复问题

---

## 十、联系与支持

如有问题，请查看：
- **Playwright官方文档**: https://playwright.dev
- **项目README**: 查看项目根目录的README文件
- **测试用例文档**: 查看 TEST_CASES.md 文件
