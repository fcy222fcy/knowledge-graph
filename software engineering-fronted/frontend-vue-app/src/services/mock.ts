// Mock 数据 - 用于前端独立开发
// 在 request.ts 中引入此文件即可启用 mock 模式

// ============================================================
// Mock 数据
// ============================================================

export const mockDocuments = [
  {
    id: 1,
    title: '软件工程导论.pdf',
    description: '软件工程基础知识教材',
    file_type: '.pdf',
    file_size: 2457600,
    status: 'completed',
    created_at: '2024-01-15T10:30:00Z',
    updated_at: '2024-01-15T10:30:00Z'
  },
  {
    id: 2,
    title: '需求分析文档模板.docx',
    description: '需求分析标准模板',
    file_type: '.docx',
    file_size: 512000,
    status: 'completed',
    created_at: '2024-01-16T14:20:00Z',
    updated_at: '2024-01-16T14:20:00Z'
  },
  {
    id: 3,
    title: 'UML建模指南.pptx',
    description: 'UML图使用指南',
    file_type: '.pptx',
    file_size: 3072000,
    status: 'processing',
    created_at: '2024-01-17T09:15:00Z',
    updated_at: '2024-01-17T09:15:00Z'
  },
  {
    id: 4,
    title: '敏捷开发实践.md',
    description: '敏捷开发方法论总结',
    file_type: '.md',
    file_size: 25600,
    status: 'completed',
    created_at: '2024-01-18T16:45:00Z',
    updated_at: '2024-01-18T16:45:00Z'
  },
  {
    id: 5,
    title: '测试用例规范.txt',
    description: '软件测试用例编写规范',
    file_type: '.txt',
    file_size: 12800,
    status: 'pending',
    created_at: '2024-01-19T11:00:00Z',
    updated_at: '2024-01-19T11:00:00Z'
  }
]

export const mockOverview = {
  total_learning_hours: 128.5,
  total_questions_asked: 256,
  total_quizzes_taken: 42,
  average_correct_rate: 0.78,
  knowledge_points_mastered: 45,
  knowledge_points_total: 80,
  today_learning_hours: 3.5,
  today_questions_asked: 12
}

export const mockMastery = [
  { knowledge_point_id: 1, knowledge_point_name: '需求分析', mastery_rate: 0.85, level: 'mastered', total_questions: 20, correct_answers: 17 },
  { knowledge_point_id: 2, knowledge_point_name: '系统设计', mastery_rate: 0.72, level: 'learning', total_questions: 15, correct_answers: 11 },
  { knowledge_point_id: 3, knowledge_point_name: '编码实现', mastery_rate: 0.90, level: 'mastered', total_questions: 25, correct_answers: 23 },
  { knowledge_point_id: 4, knowledge_point_name: '软件测试', mastery_rate: 0.45, level: 'weak', total_questions: 18, correct_answers: 8 },
  { knowledge_point_id: 5, knowledge_point_name: '项目管理', mastery_rate: 0.68, level: 'learning', total_questions: 12, correct_answers: 8 },
  { knowledge_point_id: 6, knowledge_point_name: '配置管理', mastery_rate: 0.55, level: 'weak', total_questions: 10, correct_answers: 6 },
  { knowledge_point_id: 7, knowledge_point_name: '质量保证', mastery_rate: 0.78, level: 'learning', total_questions: 14, correct_answers: 11 },
  { knowledge_point_id: 8, knowledge_point_name: '维护演化', mastery_rate: 0.62, level: 'learning', total_questions: 8, correct_answers: 5 }
]

export const mockHotPoints = [
  { knowledge_point_id: 1, knowledge_point_name: '需求分析', heat: 1250 },
  { knowledge_point_id: 2, knowledge_point_name: '系统设计', heat: 980 },
  { knowledge_point_id: 4, knowledge_point_name: '软件测试', heat: 875 },
  { knowledge_point_id: 3, knowledge_point_name: '编码实现', heat: 820 },
  { knowledge_point_id: 5, knowledge_point_name: '项目管理', heat: 650 }
]

export const mockWeakPoints = [
  {
    knowledge_point_id: 4,
    knowledge_point_name: '软件测试',
    correct_rate: 0.45,
    suggested_questions: [
      { id: 101, title: '什么是单元测试？' },
      { id: 102, title: '集成测试和系统测试的区别？' }
    ]
  },
  {
    knowledge_point_id: 6,
    knowledge_point_name: '配置管理',
    correct_rate: 0.55,
    suggested_questions: [
      { id: 103, title: '版本控制的基本概念？' },
      { id: 104, title: '分支管理策略有哪些？' }
    ]
  }
]

export const mockTrend = {
  daily_stats: [
    { date: '2024-01-13', learning_hours: 2.5, questions_asked: 8, correct_rate: 0.75 },
    { date: '2024-01-14', learning_hours: 3.2, questions_asked: 12, correct_rate: 0.80 },
    { date: '2024-01-15', learning_hours: 1.8, questions_asked: 6, correct_rate: 0.70 },
    { date: '2024-01-16', learning_hours: 4.0, questions_asked: 15, correct_rate: 0.85 },
    { date: '2024-01-17', learning_hours: 2.8, questions_asked: 10, correct_rate: 0.78 },
    { date: '2024-01-18', learning_hours: 3.5, questions_asked: 12, correct_rate: 0.82 },
    { date: '2024-01-19', learning_hours: 3.5, questions_asked: 12, correct_rate: 0.78 }
  ]
}

export const mockGraph = {
  nodes: [
    { id: 1, name: '需求分析', description: '软件需求的获取、分析和规格说明', document_id: 1, category: '核心概念' },
    { id: 2, name: '系统设计', description: '软件系统的架构和详细设计', document_id: 1, category: '核心概念' },
    { id: 3, name: '编码实现', description: '将设计转换为可执行代码', document_id: 2, category: '核心概念' },
    { id: 4, name: '软件测试', description: '验证软件质量和功能正确性', document_id: 1, category: '核心概念' },
    { id: 5, name: '项目管理', description: '软件项目的计划、组织和控制', document_id: 3, category: '管理' },
    { id: 6, name: 'UML图', description: '统一建模语言的可视化表示', document_id: 3, category: '工具' },
    { id: 7, name: '敏捷开发', description: '迭代式软件开发方法', document_id: 4, category: '方法论' },
    { id: 8, name: '瀑布模型', description: '线性顺序的开发模型', document_id: 1, category: '方法论' },
    { id: 9, name: '配置管理', description: '软件配置项的标识和控制', document_id: 2, category: '管理' },
    { id: 10, name: '质量保证', description: '确保软件质量的活动', document_id: 4, category: '管理' }
  ],
  edges: [
    { id: 1, source: 1, target: 2, relation_type: 'DEPENDS_ON', description: '需求分析是系统设计的前提' },
    { id: 2, source: 2, target: 3, relation_type: 'DEPENDS_ON', description: '系统设计指导编码实现' },
    { id: 3, source: 3, target: 4, relation_type: 'DEPENDS_ON', description: '编码完成后进行测试' },
    { id: 4, source: 1, target: 6, relation_type: 'RELATED', description: '需求分析使用UML图' },
    { id: 5, source: 2, target: 6, relation_type: 'RELATED', description: '系统设计使用UML图' },
    { id: 6, source: 5, target: 7, relation_type: 'PART_OF', description: '敏捷开发是项目管理方法' },
    { id: 7, source: 5, target: 8, relation_type: 'PART_OF', description: '瀑布模型是项目管理方法' },
    { id: 8, source: 7, target: 8, relation_type: 'RELATED', description: '敏捷和瀑布是对比方法' },
    { id: 9, source: 5, target: 9, relation_type: 'RELATED', description: '项目管理包含配置管理' },
    { id: 10, source: 5, target: 10, relation_type: 'RELATED', description: '项目管理包含质量保证' },
    { id: 11, source: 4, target: 10, relation_type: 'DEPENDS_ON', description: '测试支撑质量保证' }
  ],
  summary: {
    node_count: 10,
    edge_count: 11
  }
}

export const mockSessions = [
  {
    conversation_id: 1,
    title: '需求分析基础',
    last_question: '什么是需求分析？',
    message_count: 4,
    created_at: '2024-01-19T10:00:00Z',
    updated_at: '2024-01-19T10:30:00Z'
  },
  {
    conversation_id: 2,
    title: '软件测试方法',
    last_question: '单元测试和集成测试的区别？',
    message_count: 6,
    created_at: '2024-01-18T14:00:00Z',
    updated_at: '2024-01-18T15:00:00Z'
  },
  {
    conversation_id: 3,
    title: '敏捷开发实践',
    last_question: 'Scrum和Kanban有什么区别？',
    message_count: 8,
    created_at: '2024-01-17T09:00:00Z',
    updated_at: '2024-01-17T10:30:00Z'
  }
]

export const mockMessages = [
  { message_id: 1, role: 'user', content: '什么是需求分析？', created_at: '2024-01-19T10:00:00Z' },
  {
    message_id: 2,
    role: 'assistant',
    content: '需求分析是软件工程中的一个关键阶段，主要任务是理解和定义系统必须做什么。它包括以下几个方面：\n\n1. **需求获取**：通过访谈、问卷、观察等方式收集用户需求\n2. **需求分析**：对收集的需求进行整理、分类和优先级排序\n3. **需求规格说明**：将分析结果文档化，形成需求规格说明书\n4. **需求验证**：确保需求的正确性、完整性和一致性\n\n需求分析的质量直接影响后续的设计和开发工作。',
    created_at: '2024-01-19T10:00:05Z',
    sources: [
      { document_id: 1, document_title: '软件工程导论.pdf', content: '需求分析是软件生命周期中的第一个阶段，也是最重要的阶段之一...' }
    ],
    related_knowledge_points: [
      { id: 1, name: '需求工程' },
      { id: 2, name: '用例分析' }
    ]
  },
  { message_id: 3, role: 'user', content: '需求分析有哪些常用方法？', created_at: '2024-01-19T10:05:00Z' },
  {
    message_id: 4,
    role: 'assistant',
    content: '常用的需求分析方法包括：\n\n1. **结构化分析方法**\n   - 数据流图 (DFD)\n   - 实体关系图 (ERD)\n   - 状态转换图\n\n2. **面向对象分析方法**\n   - 用例图\n   - 类图\n   - 序列图\n\n3. **原型法**\n   - 快速原型\n   - 进化原型\n\n4. **需求获取技术**\n   - 访谈\n   - 问卷调查\n   - 联合应用开发 (JAD)\n   - 观察',
    created_at: '2024-01-19T10:05:05Z',
    sources: [
      { document_id: 1, document_title: '软件工程导论.pdf', content: '结构化分析方法是一种传统的需求分析方法...' }
    ],
    related_knowledge_points: [
      { id: 3, name: '结构化分析' },
      { id: 4, name: '面向对象分析' }
    ]
  }
]

// ============================================================
// Mock 延迟函数
// ============================================================

const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms))

const randomDelay = () => delay(300 + Math.random() * 500)

// ============================================================
// Mock API 函数
// ============================================================

export const mockAuth = {
  login: async (params: any) => {
    await randomDelay()
    if (params.username === 'admin' && params.password === '123456') {
      return {
        code: 200,
        data: {
          token: 'mock-token-xxx',
          user: {
            id: 1,
            username: 'admin',
            email: 'admin@example.com',
            avatar: '',
            role: 'student',
            major: '软件工程',
            created_at: '2024-01-01T00:00:00Z'
          }
        },
        message: '登录成功'
      }
    }
    return { code: 401, data: null, message: '用户名或密码错误' }
  },
  register: async (_params: any) => {
    await randomDelay()
    return { code: 200, data: null, message: '注册成功' }
  }
}

export const mockDocument = {
  getDocumentList: async (params?: any) => {
    await randomDelay()
    let list = [...mockDocuments]
    if (params?.keyword) {
      list = list.filter(d => d.title.includes(params.keyword))
    }
    if (params?.status) {
      list = list.filter(d => d.status === params.status)
    }
    return {
      code: 200,
      data: {
        list,
        total: list.length
      }
    }
  },
  getDocumentDetail: async (id: number) => {
    await randomDelay()
    const doc = mockDocuments.find(d => d.id === id)
    return { code: 200, data: doc || null }
  },
  getDocumentContent: async (id: number) => {
    await randomDelay()
    return {
      code: 200,
      data: {
        content: `这是文档 ${id} 的内容预览。\n\n软件工程是一门研究如何系统地、规范地、可量化地开发、运行和维护软件的学科。\n\n它涵盖了软件生命周期的各个阶段，包括需求分析、设计、编码、测试、部署和维护。`
      }
    }
  },
  uploadDocument: async (file: File, title?: string, description?: string) => {
    await delay(1000)
    return {
      code: 200,
      data: {
        id: mockDocuments.length + 1,
        title: title || file.name,
        description: description || '',
        file_type: '.' + file.name.split('.').pop(),
        file_size: file.size,
        status: 'pending',
        created_at: new Date().toISOString()
      }
    }
  },
  deleteDocument: async (_id: number) => {
    await randomDelay()
    return { code: 200, data: null }
  }
}

export const mockStats = {
  getOverview: async () => {
    await randomDelay()
    return { code: 200, data: mockOverview }
  },
  getKnowledgeMastery: async () => {
    await randomDelay()
    return { code: 200, data: mockMastery }
  },
  getHotKnowledgePoints: async (limit?: number) => {
    await randomDelay()
    return { code: 200, data: mockHotPoints.slice(0, limit || 5) }
  },
  getWeakPoints: async (limit?: number) => {
    await randomDelay()
    return { code: 200, data: mockWeakPoints.slice(0, limit || 5) }
  },
  getTrends: async (_days?: number) => {
    await randomDelay()
    return { code: 200, data: mockTrend }
  }
}

export const mockGraphApi = {
  getGraphData: async (params?: any) => {
    await randomDelay()
    let data = { ...mockGraph }
    if (params?.keyword) {
      data = {
        ...data,
        nodes: data.nodes.filter(n => n.name.includes(params.keyword))
      }
    }
    return { code: 200, data }
  },
  buildGraph: async (_documentIds: number[]) => {
    await delay(2000)
    return {
      code: 200,
      data: {
        build_id: 1,
        created_points: 10,
        created_relations: 11,
        chunk_count: 50,
        vector_count: 100,
        status: 'success',
        message: '图谱构建成功'
      }
    }
  }
}

export const mockQa = {
  getSessions: async (_params?: any) => {
    await randomDelay()
    return { code: 200, data: { list: mockSessions, total: mockSessions.length } }
  },
  getSessionMessages: async (_sessionId: number, _params?: any) => {
    await randomDelay()
    return { code: 200, data: { list: mockMessages, total: mockMessages.length } }
  },
  createSession: async (_title?: string) => {
    await randomDelay()
    const newSession = {
      conversation_id: mockSessions.length + 1,
      title: '新会话',
      last_question: '',
      message_count: 0,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    }
    return { code: 200, data: newSession }
  },
  askQuestion: async (params: any) => {
    await delay(1000 + Math.random() * 1000)
    const answers: Record<string, string> = {
      '什么是需求分析？': '需求分析是软件工程中的一个关键阶段，主要任务是理解和定义系统必须做什么。',
      '软件测试有哪些方法？': '软件测试方法包括单元测试、集成测试、系统测试、验收测试等。',
      '什么是软件生命周期？': '软件生命周期是指软件从产生到报废的整个过程。',
      '敏捷开发和瀑布模型的区别？': '敏捷开发是迭代式的，强调快速响应变化；瀑布模型是线性的，强调文档和流程。'
    }
    const answer = answers[params.question] || `关于"${params.question}"的回答：\n\n这是一个很好的问题。在软件工程中，我们需要综合考虑多个因素来给出全面的答案。`
    return {
      code: 200,
      data: {
        question_id: Date.now(),
        answer,
        sources: [{ document_id: 1, document_title: '软件工程导论.pdf', content: '相关知识内容...' }],
        related_knowledge_points: [{ id: 1, name: '软件工程基础' }],
        created_at: new Date().toISOString()
      }
    }
  },
  getAskHistory: async (params?: any) => {
    await randomDelay()
    const history = mockSessions.map(s => ({
      id: s.conversation_id,
      question: s.last_question,
      created_at: s.updated_at
    }))
    return { code: 200, data: { list: history.slice(0, params?.size || 5), total: history.length } }
  }
}
