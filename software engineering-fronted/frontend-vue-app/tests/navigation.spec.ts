import { test, expect } from '@playwright/test'
import { setupLoggedIn } from './helpers'

test.describe('导航功能测试', () => {
  test.beforeEach(async ({ page }) => {
    await setupLoggedIn(page)
    await page.goto('/home')
    await page.waitForLoadState('networkidle')

    // 导航测试需要 Pinia store 中的 userInfo，手动同步
    await page.evaluate(() => {
      const stored = localStorage.getItem('userInfo')
      if (stored) {
        const userInfo = JSON.parse(stored)
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

  test('侧边栏导航高亮', async ({ page }) => {
    // 检查首页导航高亮
    await expect(page.locator('.nav-item.active')).toHaveText('首页')

    // 导航到知识图谱
    await page.locator('.nav-item:has-text("知识图谱")').click()
    await expect(page).toHaveURL('/knowledge-graph')
    await expect(page.locator('.nav-item.active')).toHaveText('知识图谱')

    // 导航到问答中心
    await page.locator('.nav-item:has-text("问答中心")').click()
    await expect(page).toHaveURL('/qa')
    await expect(page.locator('.nav-item.active')).toHaveText('问答中心')

    // 导航到资料管理
    await page.locator('.nav-item:has-text("资料管理")').click()
    await expect(page).toHaveURL('/files')
    await expect(page.locator('.nav-item.active')).toHaveText('资料管理')

    // 导航到分析统计
    await page.locator('.nav-item:has-text("分析统计")').click()
    await expect(page).toHaveURL('/stats')
    await expect(page.locator('.nav-item.active')).toHaveText('分析统计')
  })

  test('路由守卫功能', async ({ page }) => {
    // 清除登录状态（同时移除 mock）
    await page.unroute('**/api/**')
    await page.evaluate(() => {
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
    })

    // 导航到登录页（确保在已知状态）
    await page.goto('/login')
    await page.waitForLoadState('networkidle')

    // 再次确保 token 已移除（addInitScript 可能在页面加载时重新设置）
    await page.evaluate(() => {
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
    })

    // 使用客户端导航触发路由守卫（避免 addInitScript 在 page.goto 时重新设置 token）
    await page.evaluate(() => {
      const app = document.querySelector('#app')?.__vue_app__
      if (app) {
        const router = app.config.globalProperties.$router
        router.push('/home')
      }
    })

    // 等待 URL 变化
    await page.waitForTimeout(2000)

    // 应该重定向到登录页面
    await expect(page).toHaveURL('/login')
  })

  test('页面标题更新', async ({ page }) => {
    // 检查首页标题
    await expect(page).toHaveTitle(/首页/)

    // 导航到知识图谱
    await page.locator('.nav-item:has-text("知识图谱")').click()
    await expect(page).toHaveTitle(/知识图谱/)

    // 导航到问答中心
    await page.locator('.nav-item:has-text("问答中心")').click()
    await expect(page).toHaveTitle(/问答中心/)

    // 导航到资料管理
    await page.locator('.nav-item:has-text("资料管理")').click()
    await expect(page).toHaveTitle(/资料管理/)

    // 导航到分析统计
    await page.locator('.nav-item:has-text("分析统计")').click()
    await expect(page).toHaveTitle(/分析统计/)
  })

  test('Logo区域显示正确', async ({ page }) => {
    // 检查Logo区域
    await expect(page.locator('.logo')).toBeVisible()
    await expect(page.locator('.logo-mark')).toHaveText('SE')
    await expect(page.locator('.logo-text')).toHaveText('SE智图问答')
  })

  test('浏览器前进后退', async ({ page }) => {
    // 导航到知识图谱
    await page.locator('.nav-item:has-text("知识图谱")').click()
    await expect(page).toHaveURL('/knowledge-graph')

    // 导航到问答中心
    await page.locator('.nav-item:has-text("问答中心")').click()
    await expect(page).toHaveURL('/qa')

    // 浏览器后退
    await page.goBack()
    await expect(page).toHaveURL('/knowledge-graph')

    // 浏览器前进
    await page.goForward()
    await expect(page).toHaveURL('/qa')
  })
})
