import type { Page } from '@playwright/test'

/**
 * 设置已登录状态并 mock 所有 API 响应
 * 防止因后端未运行导致 401 清除 token
 */
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

  // Mock 所有 API 请求，防止 401 清除 token
  await page.route('**/api/**', (route) => {
    const url = route.request().url()

    // 首页统计
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

    // 趋势数据
    if (url.includes('/analytics/trends')) {
      return route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          code: 200,
          data: {
            daily_stats: [
              { date: '2024-01-01', learning_hours: 2, correct_rate: 0.75 },
              { date: '2024-01-02', learning_hours: 1.5, correct_rate: 0.8 },
              { date: '2024-01-03', learning_hours: 3, correct_rate: 0.85 }
            ]
          },
          message: 'success'
        })
      })
    }

    // 问答历史
    if (url.includes('/ask/history') || url.includes('/ask/sessions')) {
      return route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          code: 200,
          data: {
            list: [
              { conversation_id: 1, title: '测试会话', last_question: '什么是软件工程？', message_count: 5 }
            ],
            total: 1
          },
          message: 'success'
        })
      })
    }

    // 图谱数据
    if (url.includes('/graph')) {
      return route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          code: 200,
          data: {
            nodes: [
              { id: '1', name: '软件工程', category: '概念', description: '软件工程的定义' },
              { id: '2', name: '需求分析', category: '阶段', description: '需求分析阶段' }
            ],
            edges: [
              { source: '1', target: '2', type: 'CONTAINS' }
            ],
            summary: { node_count: 2, edge_count: 1 }
          },
          message: 'success'
        })
      })
    }

    // 文档列表
    if (url.includes('/document')) {
      return route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          code: 200,
          data: {
            list: [],
            total: 0
          },
          message: 'success'
        })
      })
    }

    // 知识掌握度
    if (url.includes('/analytics/knowledge-mastery')) {
      return route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          code: 200,
          data: [],
          message: 'success'
        })
      })
    }

    // 薄弱点
    if (url.includes('/analytics/weak-points')) {
      return route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          code: 200,
          data: [],
          message: 'success'
        })
      })
    }

    // 热点知识点
    if (url.includes('/analytics/hot-points')) {
      return route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          code: 200,
          data: [],
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
