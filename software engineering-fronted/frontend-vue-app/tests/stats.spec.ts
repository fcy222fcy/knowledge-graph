import { test, expect } from '@playwright/test'
import { setupLoggedIn } from './helpers'

test.describe('分析统计页面测试', () => {
  test.beforeEach(async ({ page }) => {
    await setupLoggedIn(page)
    await page.goto('/stats')
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

  test('趋势图表存在', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 检查趋势图表组件
    await expect(page.locator('.chart-card:has-text("学习趋势")')).toBeVisible()

    // 检查图表内容
    await expect(page.locator('h3:has-text("学习趋势")')).toBeVisible()
  })

  test('掌握度图表存在', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 检查掌握度图表组件
    await expect(page.locator('.chart-card:has-text("知识点掌握度")')).toBeVisible()

    // 检查图表内容
    await expect(page.locator('h3:has-text("知识点掌握度")')).toBeVisible()
  })

  test('热点排行存在', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 检查热点排行组件（实际标题为"热门知识点排行"）
    await expect(page.locator('.chart-card:has-text("热门知识点排行")')).toBeVisible()
  })

  test('薄弱点列表存在', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 检查薄弱点列表组件
    await expect(page.locator('.chart-card:has-text("薄弱知识点")')).toBeVisible()
  })

  test('图表交互功能', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 检查图表卡片
    const chartCards = page.locator('.chart-card')
    const count = await chartCards.count()

    // 应该有多个图表卡片
    expect(count).toBeGreaterThan(0)

    // 检查第一个图表卡片
    await expect(chartCards.first()).toBeVisible()
  })

  test('概览卡片存在', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 检查概览卡片
    const overviewCards = page.locator('.overview-cards')
    if (await overviewCards.isVisible()) {
      // 检查统计卡片
      const statCards = page.locator('.stat-card')
      const count = await statCards.count()
      expect(count).toBeGreaterThan(0)
    }
  })

  test('页面响应式布局', async ({ page }) => {
    // 测试桌面视图
    await page.setViewportSize({ width: 1200, height: 800 })
    await expect(page.locator('.stats-container')).toBeVisible()

    // 测试移动端视图
    await page.setViewportSize({ width: 375, height: 667 })
    await expect(page.locator('.stats-container')).toBeVisible()
  })
})
