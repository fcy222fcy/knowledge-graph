import { test, expect } from '@playwright/test'
import { setupLoggedIn } from './helpers'

test.describe('知识图谱页面测试', () => {
  test.beforeEach(async ({ page }) => {
    await setupLoggedIn(page)
    await page.goto('/knowledge-graph')
    await page.waitForLoadState('networkidle')
  })

  test('页面加载正确', async ({ page }) => {
    // 检查页面标题
    await expect(page).toHaveTitle(/SE智图问答/)

    // 检查主布局存在
    await expect(page.locator('.app-layout')).toBeVisible()
    await expect(page.locator('.sidebar')).toBeVisible()
    await expect(page.locator('.main')).toBeVisible()
  })

  test('图谱工具栏完整', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 检查工具栏存在
    await expect(page.locator('.graph-toolbar')).toBeVisible()

    // 检查构建图谱按钮
    await expect(page.locator('button:has-text("构建图谱")')).toBeVisible()

    // 检查刷新按钮
    await expect(page.locator('button:has-text("刷新")')).toBeVisible()

    // 检查搜索框
    await expect(page.locator('input[placeholder="搜索知识点..."]')).toBeVisible()

    // 检查关系类型筛选
    await expect(page.locator('.graph-toolbar .el-select')).toBeVisible()
  })

  test('图谱画布存在', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 检查图谱画布存在
    await expect(page.locator('.graph-canvas')).toBeVisible()

    // 检查SVG元素（D3.js图谱）
    await expect(page.locator('.graph-canvas svg')).toBeVisible()
  })

  test('图谱摘要信息', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 检查图谱摘要信息（如果有的话）
    const graphSummary = page.locator('.graph-summary')
    if (await graphSummary.isVisible()) {
      // 检查节点和关系统计
      await expect(page.locator('text=节点:')).toBeVisible()
      await expect(page.locator('text=关系:')).toBeVisible()
    }
  })

  test('构建图谱对话框', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 点击构建图谱按钮
    const buildButton = page.locator('button:has-text("构建图谱")')
    await buildButton.click()

    // 检查对话框出现
    await expect(page.locator('.el-dialog')).toBeVisible()
  })

  test('图谱搜索功能', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 找到搜索框
    const searchInput = page.locator('input[placeholder="搜索知识点..."]')
    await expect(searchInput).toBeVisible()

    // 输入搜索内容
    await searchInput.fill('软件工程')

    // 按回车触发搜索
    await searchInput.press('Enter')
  })

  test('图谱刷新功能', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 点击刷新按钮
    const refreshButton = page.locator('button:has-text("刷新")')
    await expect(refreshButton).toBeVisible()
    await refreshButton.click()
  })

  test('页面响应式布局', async ({ page }) => {
    // 测试桌面视图
    await page.setViewportSize({ width: 1200, height: 800 })
    await expect(page.locator('.graph-container')).toBeVisible()

    // 测试移动端视图
    await page.setViewportSize({ width: 375, height: 667 })
    await expect(page.locator('.graph-container')).toBeVisible()
  })
})
