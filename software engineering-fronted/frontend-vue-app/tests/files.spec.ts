import { test, expect } from '@playwright/test'
import { setupLoggedIn } from './helpers'

test.describe('资料管理页面测试', () => {
  test.beforeEach(async ({ page }) => {
    await setupLoggedIn(page)
    await page.goto('/files')
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

  test('文档表格存在', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 检查页面头部存在
    await expect(page.locator('.page-header')).toBeVisible()
    await expect(page.locator('h2:has-text("资料管理")')).toBeVisible()

    // 检查上传按钮
    await expect(page.locator('button:has-text("上传文档")')).toBeVisible()

    // 检查搜索栏
    await expect(page.locator('.search-bar')).toBeVisible()
    await expect(page.locator('input[placeholder="搜索文档..."]')).toBeVisible()
  })

  test('上传按钮存在', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 检查上传按钮
    await expect(page.locator('button:has-text("上传文档")')).toBeVisible()
  })

  test('上传对话框功能', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 点击上传按钮
    const uploadButton = page.locator('button:has-text("上传文档")')
    await uploadButton.click()

    // 等待对话框出现
    await page.waitForTimeout(500)

    // 检查上传对话框出现
    await expect(page.locator('.el-dialog')).toBeVisible()
  })

  test('文档详情功能', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 检查文档详情组件（如果存在文档的话）
    const viewButtons = page.locator('button:has-text("查看")')
    const count = await viewButtons.count()

    if (count > 0) {
      // 点击查看按钮
      await viewButtons.first().click()

      // 检查详情抽屉出现
      await expect(page.locator('.el-drawer')).toBeVisible()
    }
  })

  test('文档操作按钮', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 检查操作按钮（查看、删除）
    const viewButtons = page.locator('button:has-text("查看")')
    const deleteButtons = page.locator('button:has-text("删除")')

    // 如果有文档，应该有操作按钮
    const viewCount = await viewButtons.count()
    if (viewCount > 0) {
      await expect(viewButtons.first()).toBeVisible()
    }
  })

  test('分页功能', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 检查分页包装器
    await expect(page.locator('.pagination-wrapper')).toBeVisible()
  })

  test('搜索功能', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 检查搜索框
    const searchInput = page.locator('input[placeholder="搜索文档..."]')
    await expect(searchInput).toBeVisible()

    // 输入搜索内容
    await searchInput.fill('软件工程')

    // 检查状态筛选下拉框
    const statusSelect = page.locator('.search-bar .el-select')
    await expect(statusSelect).toBeVisible()
  })

  test('页面响应式布局', async ({ page }) => {
    // 测试桌面视图
    await page.setViewportSize({ width: 1200, height: 800 })
    await expect(page.locator('.files-container')).toBeVisible()

    // 测试移动端视图
    await page.setViewportSize({ width: 375, height: 667 })
    await expect(page.locator('.files-container')).toBeVisible()
  })
})
