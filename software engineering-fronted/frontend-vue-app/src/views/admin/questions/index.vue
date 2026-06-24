<template>
  <div class="admin-questions" v-loading="loading">
    <div class="page-header">
      <h2 class="page-title">题目管理</h2>
      <div class="header-actions">
        <el-button type="primary" @click="fetchQuestions">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 题目表格 -->
    <el-table :data="questions" stripe border style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="title" label="题目" min-width="300" show-overflow-tooltip />
      <el-table-column prop="type" label="类型" width="100">
        <template #default="{ row }">
          <el-tag>{{ row.type === 'single' ? '单选' : '多选' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="difficulty" label="难度" width="100">
        <template #default="{ row }">
          <el-tag :type="getDifficultyType(row.difficulty)">{{ getDifficultyLabel(row.difficulty) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="answer" label="答案" width="100" />
      <el-table-column prop="created_at" label="创建时间" width="180" />
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="handleView(row)">查看</el-button>
          <el-popconfirm
            title="确定要删除该题目吗？"
            @confirm="handleDelete(row)"
          >
            <template #reference>
              <el-button type="danger" link>删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <el-pagination
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :page-sizes="[10, 20, 50]"
      :total="total"
      layout="total, sizes, prev, pager, next, jumper"
      style="margin-top: 16px; justify-content: flex-end;"
      @size-change="fetchQuestions"
      @current-change="fetchQuestions"
    />

    <!-- 查看详情对话框 -->
    <el-dialog v-model="detailVisible" title="题目详情" width="600px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="题目">{{ detailQuestion?.title }}</el-descriptions-item>
        <el-descriptions-item label="类型">{{ detailQuestion?.type === 'single' ? '单选' : '多选' }}</el-descriptions-item>
        <el-descriptions-item label="难度">{{ getDifficultyLabel(detailQuestion?.difficulty || '') }}</el-descriptions-item>
        <el-descriptions-item label="选项">
          <div v-if="detailQuestion?.options">
            <div v-for="opt in detailQuestion.options" :key="opt.key">
              <strong>{{ opt.key }}.</strong> {{ opt.value }}
            </div>
          </div>
        </el-descriptions-item>
        <el-descriptions-item label="答案">{{ detailQuestion?.answer }}</el-descriptions-item>
        <el-descriptions-item label="解析">{{ detailQuestion?.explanation }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { getQuestions, deleteQuestion } from '@/services/admin'
import type { QuestionItem } from '@/services/admin'

const loading = ref(false)
const questions = ref<QuestionItem[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

// 详情相关
const detailVisible = ref(false)
const detailQuestion = ref<QuestionItem | null>(null)

const getDifficultyType = (difficulty: string) => {
  const types: Record<string, string> = {
    easy: 'success',
    medium: 'warning',
    hard: 'danger'
  }
  return types[difficulty] || 'info'
}

const getDifficultyLabel = (difficulty: string) => {
  const labels: Record<string, string> = {
    easy: '简单',
    medium: '中等',
    hard: '困难'
  }
  return labels[difficulty] || '未知'
}

const fetchQuestions = async () => {
  loading.value = true
  try {
    const res = await getQuestions({ page: page.value, size: pageSize.value })
    questions.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    ElMessage.error('获取题目列表失败')
  } finally {
    loading.value = false
  }
}

const handleView = (row: QuestionItem) => {
  detailQuestion.value = row
  detailVisible.value = true
}

const handleDelete = async (row: QuestionItem) => {
  try {
    await deleteQuestion(row.id)
    ElMessage.success('删除成功')
    fetchQuestions()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

onMounted(() => {
  fetchQuestions()
})
</script>

<style scoped>
.admin-questions {
  background: #fff;
  padding: 24px;
  border-radius: 8px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}
</style>
