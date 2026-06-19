import { test, expect } from '@playwright/test'
import { setupLoggedIn } from './helpers'

test.describe('问答中心测试', () => {
  test.beforeEach(async ({ page }) => {
    await setupLoggedIn(page)
    await page.goto('/qa')
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

  test('会话列表存在', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 检查会话列表存在
    await expect(page.locator('.session-list')).toBeVisible()

    // 检查会话列表头部
    await expect(page.locator('.session-header')).toBeVisible()
    await expect(page.locator('h3:has-text("会话列表")')).toBeVisible()

    // 检查新建会话按钮
    await expect(page.locator('button:has-text("新建")')).toBeVisible()
  })

  test('聊天面板存在', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 检查聊天面板存在
    await expect(page.locator('.chat-panel')).toBeVisible()

    // 检查消息输入区域
    await expect(page.locator('.input-area')).toBeVisible()

    // 检查消息输入框
    await expect(page.locator('textarea[placeholder*="输入你的问题"]')).toBeVisible()

    // 检查发送按钮
    await expect(page.locator('button:has-text("发送")')).toBeVisible()
  })

  test('空状态显示', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 检查消息区域
    const messagesArea = page.locator('.messages-area')
    await expect(messagesArea).toBeVisible()

    // 如果没有消息，应该显示空状态（实际标题为"开始提问吧"）
    const emptyChat = page.locator('.empty-chat')
    if (await emptyChat.isVisible()) {
      await expect(page.locator('h3:has-text("开始提问吧")')).toBeVisible()
    }
  })

  test('新建会话功能', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 点击新建会话按钮
    const newSessionButton = page.locator('button:has-text("新建")')
    await newSessionButton.click()

    // 等待可能的API请求完成
    await page.waitForTimeout(1000)
  })

  test('发送消息功能', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 找到消息输入框
    const messageInput = page.locator('textarea[placeholder*="输入你的问题"]')
    await expect(messageInput).toBeVisible()

    // 输入消息
    await messageInput.fill('什么是软件工程？')

    // 检查发送按钮是否启用
    const sendButton = page.locator('button:has-text("发送")')
    await expect(sendButton).toBeEnabled()

    // 点击发送按钮
    await sendButton.click()

    // 等待可能的API请求完成
    await page.waitForTimeout(1000)
  })

  test('会话切换功能', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 检查会话列表
    const sessionItems = page.locator('.session-item')
    const count = await sessionItems.count()

    if (count > 0) {
      // 点击第一个会话
      await sessionItems.first().click()

      // 检查聊天面板更新
      await expect(page.locator('.chat-panel')).toBeVisible()
    }
  })

  test('消息区域显示', async ({ page }) => {
    // 等待页面加载
    await page.waitForLoadState('networkidle')

    // 检查消息区域
    await expect(page.locator('.messages-area')).toBeVisible()
  })

  test('页面响应式布局', async ({ page }) => {
    // 测试桌面视图
    await page.setViewportSize({ width: 1200, height: 800 })
    await expect(page.locator('.qa-container')).toBeVisible()

    // 测试移动端视图
    await page.setViewportSize({ width: 375, height: 667 })
    await expect(page.locator('.qa-container')).toBeVisible()
  })
})
