import { test, expect } from '@playwright/test'

test.describe('注册页面测试', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/register')
    await page.waitForLoadState('networkidle')
  })

  test('页面加载正确', async ({ page }) => {
    await expect(page).toHaveTitle(/SE智图问答/)
    await expect(page.locator('.register-card')).toBeVisible()
    await expect(page.locator('h1')).toHaveText('创建账号')
  })

  test('表单元素完整', async ({ page }) => {
    await expect(page.locator('input[placeholder="请输入用户名"]')).toBeVisible()
    await expect(page.locator('input[placeholder="请输入邮箱"]')).toBeVisible()
    await expect(page.locator('input[placeholder="请输入密码"]')).toBeVisible()
    await expect(page.locator('input[placeholder="请再次输入密码"]')).toBeVisible()
    await expect(page.locator('button:has-text("注册")')).toBeVisible()
    await expect(page.locator('text=已有账号？')).toBeVisible()
    await expect(page.locator('text=立即登录')).toBeVisible()
  })

  test('表单验证 - 空值提交', async ({ page }) => {
    const registerButton = page.locator('button:has-text("注册")')
    await registerButton.click()

    await expect(page.locator('text=请输入用户名')).toBeVisible()
    await expect(page.locator('text=请输入邮箱')).toBeVisible()
    await expect(page.locator('text=请输入密码')).toBeVisible()
    await expect(page.locator('text=请再次输入密码')).toBeVisible()
  })

  test('密码确认验证', async ({ page }) => {
    await page.locator('input[placeholder="请输入用户名"]').fill('testuser')
    await page.locator('input[placeholder="请输入邮箱"]').fill('test@example.com')
    await page.locator('input[placeholder="请输入密码"]').fill('password123')
    await page.locator('input[placeholder="请再次输入密码"]').fill('password456')

    const registerButton = page.locator('button:has-text("注册")')
    await registerButton.click()

    await expect(page.locator('text=两次输入的密码不一致')).toBeVisible()
  })

  test('导航到登录页面', async ({ page }) => {
    const loginLink = page.locator('text=立即登录')
    await loginLink.click()
    await expect(page).toHaveURL('/login')
  })

  test('页面响应式布局', async ({ page }) => {
    await page.setViewportSize({ width: 1200, height: 800 })
    await expect(page.locator('.register-card')).toBeVisible()

    await page.setViewportSize({ width: 375, height: 667 })
    await expect(page.locator('.register-card')).toBeVisible()
  })
})
