import { test, expect } from '@playwright/test'

test.describe('登录页面测试', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login')
  })

  test('页面加载正确', async ({ page }) => {
    // 检查页面标题
    await expect(page).toHaveTitle(/SE智图问答/)

    // 检查登录表单存在
    await expect(page.locator('.login-card')).toBeVisible()
    await expect(page.locator('.login-header h1')).toHaveText('欢迎回来')
    await expect(page.locator('.login-header p')).toHaveText('登录 SE智图问答，继续你的学习之旅')
  })

  test('表单元素完整', async ({ page }) => {
    // 检查用户名输入框
    const usernameInput = page.locator('input[placeholder="请输入用户名"]')
    await expect(usernameInput).toBeVisible()

    // 检查密码输入框
    const passwordInput = page.locator('input[placeholder="请输入密码"]')
    await expect(passwordInput).toBeVisible()

    // 检查登录按钮
    const loginButton = page.locator('button:has-text("登录")')
    await expect(loginButton).toBeVisible()

    // 检查记住我选项
    await expect(page.locator('text=记住我')).toBeVisible()

    // 检查忘记密码链接
    await expect(page.locator('text=忘记密码？')).toBeVisible()

    // 检查注册链接
    await expect(page.locator('text=还没有账号？')).toBeVisible()
    await expect(page.locator('text=立即注册')).toBeVisible()
  })

  test('表单验证 - 空值提交', async ({ page }) => {
    // 点击登录按钮 without filling form
    const loginButton = page.locator('button:has-text("登录")')
    await loginButton.click()

    // 检查验证错误消息
    await expect(page.locator('text=请输入用户名')).toBeVisible()
    await expect(page.locator('text=请输入密码')).toBeVisible()
  })

  test('表单验证 - 仅填写用户名', async ({ page }) => {
    // 填写用户名
    const usernameInput = page.locator('input[placeholder="请输入用户名"]')
    await usernameInput.fill('testuser')

    // 点击登录按钮
    const loginButton = page.locator('button:has-text("登录")')
    await loginButton.click()

    // 检查密码验证错误
    await expect(page.locator('text=请输入密码')).toBeVisible()
  })

  test('表单验证 - 仅填写密码', async ({ page }) => {
    // 填写密码
    const passwordInput = page.locator('input[placeholder="请输入密码"]')
    await passwordInput.fill('password123')

    // 点击登录按钮
    const loginButton = page.locator('button:has-text("登录")')
    await loginButton.click()

    // 检查用户名验证错误
    await expect(page.locator('text=请输入用户名')).toBeVisible()
  })

  test('导航到注册页面', async ({ page }) => {
    // 点击注册链接
    const registerLink = page.locator('text=立即注册')
    await registerLink.click()

    // 验证导航到注册页面
    await expect(page).toHaveURL('/register')
  })

  test('密码显示/隐藏功能', async ({ page }) => {
    const passwordInput = page.locator('input[placeholder="请输入密码"]')

    // 初始状态应该是密码类型
    await expect(passwordInput).toHaveAttribute('type', 'password')

    // 点击显示密码按钮（如果存在）
    const showPasswordButton = page.locator('.el-input__suffix').first()
    if (await showPasswordButton.isVisible()) {
      await showPasswordButton.click()
      // 密码应该显示为文本
      await expect(passwordInput).toHaveAttribute('type', 'text')
    }
  })

  test('页面响应式布局', async ({ page }) => {
    // 测试桌面视图
    await page.setViewportSize({ width: 1200, height: 800 })
    await expect(page.locator('.login-card')).toBeVisible()

    // 测试移动端视图
    await page.setViewportSize({ width: 375, height: 667 })
    await expect(page.locator('.login-card')).toBeVisible()
  })
})
