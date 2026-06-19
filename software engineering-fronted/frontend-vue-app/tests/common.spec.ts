import { test, expect } from '@playwright/test'
import { setupLoggedIn } from './helpers'

test.describe('通用功能测试', () => {
  test('页面加载性能', async ({ page }) => {
    // 测量页面加载时间
    const startTime = Date.now()
    await page.goto('/login')
    await page.waitForLoadState('networkidle')
    const loadTime = Date.now() - startTime

    // 页面加载时间应该在3秒以内
    expect(loadTime).toBeLessThan(3000)

    console.log(`页面加载时间: ${loadTime}ms`)
  })

  test('内存泄漏检测', async ({ page }) => {
    // 获取初始内存使用
    const initialMemory = await page.evaluate(() => {
      return (performance as any).memory?.usedJSHeapSize || 0
    })

    // 多次导航页面
    for (let i = 0; i < 5; i++) {
      await page.goto('/login')
      await page.goto('/home')
    }

    // 获取最终内存使用
    const finalMemory = await page.evaluate(() => {
      return (performance as any).memory?.usedJSHeapSize || 0
    })

    // 内存增长应该在合理范围内（10MB以内）
    if (initialMemory > 0 && finalMemory > 0) {
      const memoryGrowth = finalMemory - initialMemory
      expect(memoryGrowth).toBeLessThan(10 * 1024 * 1024)
      console.log(`内存增长: ${memoryGrowth / 1024 / 1024}MB`)
    }
  })

  test('控制台错误检测', async ({ page }) => {
    const errors: string[] = []

    // 监听控制台错误
    page.on('console', (msg) => {
      if (msg.type() === 'error') {
        errors.push(msg.text())
      }
    })

    // 监听页面错误
    page.on('pageerror', (error) => {
      errors.push(error.message)
    })

    // 访问各个页面
    const pages = ['/login', '/home', '/knowledge-graph', '/qa', '/files', '/stats']

    for (const pagePath of pages) {
      await page.goto(pagePath)
      await page.waitForLoadState('networkidle')
    }

    // 检查是否有严重错误
    const criticalErrors = errors.filter(error =>
      !error.includes('ResizeObserver') &&
      !error.includes('Non-Error promise rejection') &&
      !error.includes('Loading chunk')
    )

    if (criticalErrors.length > 0) {
      console.log('发现控制台错误:', criticalErrors)
    }

    // 不应该有严重错误
    expect(criticalErrors.length).toBe(0)
  })

  test('响应式断点测试', async ({ page }) => {
    await setupLoggedIn(page)

    // 测试不同断点
    const breakpoints = [
      { width: 375, height: 667, name: '移动端' },
      { width: 768, height: 1024, name: '平板' },
      { width: 1024, height: 768, name: '小桌面' },
      { width: 1200, height: 800, name: '桌面' },
      { width: 1920, height: 1080, name: '大桌面' }
    ]

    for (const breakpoint of breakpoints) {
      await page.setViewportSize({ width: breakpoint.width, height: breakpoint.height })
      await page.goto('/home')
      await page.waitForLoadState('networkidle')

      // 检查页面是否正常显示
      await expect(page.locator('.app-layout')).toBeVisible()
      console.log(`${breakpoint.name} (${breakpoint.width}x${breakpoint.height}) - 正常`)
    }
  })

  test('网络请求监控', async ({ page }) => {
    const requests: string[] = []
    const failedRequests: string[] = []

    // 监听网络请求
    page.on('request', (request) => {
      if (request.url().includes('/api/')) {
        requests.push(request.url())
      }
    })

    // 监听失败的请求
    page.on('requestfailed', (request) => {
      if (request.url().includes('/api/')) {
        failedRequests.push(request.url())
      }
    })

    // 访问首页
    await page.goto('/home')
    await page.waitForLoadState('networkidle')

    console.log(`API请求数量: ${requests.length}`)
    console.log(`失败请求数量: ${failedRequests.length}`)

    if (failedRequests.length > 0) {
      console.log('失败的请求:', failedRequests)
    }
  })

  test('本地存储功能', async ({ page }) => {
    // 访问登录页面
    await page.goto('/login')

    // 设置本地存储
    await page.evaluate(() => {
      localStorage.setItem('testKey', 'testValue')
    })

    // 验证本地存储
    const value = await page.evaluate(() => {
      return localStorage.getItem('testKey')
    })

    expect(value).toBe('testValue')

    // 清除本地存储
    await page.evaluate(() => {
      localStorage.removeItem('testKey')
    })

    // 验证清除
    const clearedValue = await page.evaluate(() => {
      return localStorage.getItem('testKey')
    })

    expect(clearedValue).toBeNull()
  })

  test('Cookie功能', async ({ page }) => {
    // 设置Cookie
    await page.goto('/login')
    await page.context().addCookies([{
      name: 'testCookie',
      value: 'testValue',
      domain: 'localhost',
      path: '/'
    }])

    // 验证Cookie
    const cookies = await page.context().cookies()
    const testCookie = cookies.find(cookie => cookie.name === 'testCookie')

    expect(testCookie).toBeDefined()
    expect(testCookie?.value).toBe('testValue')
  })
})
