import { test, expect } from '@playwright/test'
import { setupLoggedIn } from './helpers'

test.describe('首页测试', () => {
  test.beforeEach(async ({ page }) => {
    await setupLoggedIn(page)
    await page.goto('/home')
    await page.waitForLoadState('networkidle')
    // 初始化 Pinia store 中的 userInfo（store 不会自动从 localStorage 读取 userInfo）
    await page.evaluate(() => {
      const stored = localStorage.getItem('userInfo')
      if (stored) {
        const userInfo = JSON.parse(stored)
        // 查找 Pinia 实例并设置用户信息
        const app = document.querySelector('#app')?.__vue_app__
        if (app) {
          const pinia = app.config.globalProperties.$pinia
          if (pinia?.state?.value?.user) {
            pinia.state.value.user.userInfo = userInfo
          }
        }
      }
    })
  })

  test('页面加载正确', async ({ page }) => {
    await expect(page).toHaveTitle(/SE智图问答/)
    await expect(page.locator('.app-layout')).toBeVisible()
    await expect(page.locator('.sidebar')).toBeVisible()
    await expect(page.locator('.main')).toBeVisible()
  })

  test('侧边栏导航完整', async ({ page }) => {
    await expect(page.locator('.logo-text')).toHaveText('SE智图问答')

    const navItems = page.locator('.nav-item')
    await expect(navItems).toHaveCount(5)

    await expect(page.locator('.nav-item:has-text("首页")')).toBeVisible()
    await expect(page.locator('.nav-item:has-text("知识图谱")')).toBeVisible()
    await expect(page.locator('.nav-item:has-text("问答中心")')).toBeVisible()
    await expect(page.locator('.nav-item:has-text("资料管理")')).toBeVisible()
    await expect(page.locator('.nav-item:has-text("分析统计")')).toBeVisible()
  })

  test('用户信息显示', async ({ page }) => {
    await expect(page.locator('.user-avatar')).toBeVisible()
    await expect(page.locator('.user-name')).toHaveText('测试用户')
    await expect(page.locator('.user-role')).toHaveText('软件工程专业')
  })

  test('首页内容加载', async ({ page }) => {
    await expect(page.locator('.quick-actions')).toBeVisible()
  })

  test('导航到各个页面', async ({ page }) => {
    await page.locator('.nav-item:has-text("知识图谱")').click()
    await expect(page).toHaveURL('/knowledge-graph')

    await page.locator('.nav-item:has-text("问答中心")').click()
    await expect(page).toHaveURL('/qa')

    await page.locator('.nav-item:has-text("资料管理")').click()
    await expect(page).toHaveURL('/files')

    await page.locator('.nav-item:has-text("分析统计")').click()
    await expect(page).toHaveURL('/stats')

    await page.locator('.nav-item:has-text("首页")').click()
    await expect(page).toHaveURL('/home')
  })

  test('退出登录功能', async ({ page }) => {
    const userCard = page.locator('.user-card')
    await userCard.click()

    await expect(page.locator('.el-message-box')).toBeVisible()
    await expect(page.locator('text=确定要退出登录吗？')).toBeVisible()

    await page.locator('button:has-text("确定")').click()
    await expect(page).toHaveURL('/login')
  })

  test('页面响应式布局', async ({ page }) => {
    await page.setViewportSize({ width: 1200, height: 800 })
    await expect(page.locator('.sidebar')).toBeVisible()
    await expect(page.locator('.main')).toBeVisible()

    await page.setViewportSize({ width: 375, height: 667 })
    await expect(page.locator('.main')).toBeVisible()
  })
})
