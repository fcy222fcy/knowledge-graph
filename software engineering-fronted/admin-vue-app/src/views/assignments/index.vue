<template>
  <div class="assignments-page">
    <div class="page-header">
      <h2>作业管理</h2>
      <el-button type="primary" @click="openCreateDialog">＋ 新建作业</el-button>
    </div>

    <!-- 作业列表 -->
    <el-table :data="assignments" stripe v-loading="loading">
      <el-table-column prop="title" label="作业名称" min-width="200" />
      <el-table-column prop="chapter" label="章节" width="120" />
      <el-table-column label="题目数" width="80" align="center">
        <template #default="{ row }">{{ row.question_num }} 题</template>
      </el-table-column>
      <el-table-column label="提交" width="100" align="center">
        <template #default="{ row }">{{ row.submit_count }}</template>
      </el-table-column>
      <el-table-column label="总分" width="80" align="center">
        <template #default="{ row }">{{ row.total_score }} 分</template>
      </el-table-column>
      <el-table-column label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.status === 'draft'" type="warning">草稿</el-tag>
          <el-tag v-else-if="row.status === 'published'" type="success">已发布</el-tag>
          <el-tag v-else type="info">已截止</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="deadline" label="截止时间" width="170" />
      <el-table-column label="操作" width="280" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="viewDetail(row)">查看</el-button>
          <el-button v-if="row.status !== 'draft'" size="small" @click="viewSubmissions(row)">查看提交</el-button>
          <el-button v-if="row.status === 'draft'" size="small" type="primary" @click="openEditDialog(row)">编辑</el-button>
          <el-button v-if="row.status === 'draft'" size="small" type="success" @click="handlePublish(row)">发布</el-button>
          <el-button v-if="row.status === 'published'" size="small" type="warning" @click="handleClose(row)">关闭</el-button>
          <el-button v-if="row.status !== 'published'" size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="page"
      :page-size="10"
      :total="total"
      layout="total, prev, pager, next"
      @current-change="fetchAssignments"
      style="margin-top: 16px; justify-content: flex-end;"
    />

    <!-- 新建/编辑作业弹窗 -->
    <AssignmentDialog
      v-model="dialogVisible"
      :assignment="editingAssignment"
      @saved="onSaved"
    />

    <!-- 查看详情弹窗 -->
    <el-dialog v-model="detailVisible" title="作业详情" width="800px" destroy-on-close>
      <div v-if="detailData">
        <p><strong>作业名称：</strong>{{ detailData.title }}</p>
        <p><strong>章节：</strong>{{ detailData.chapter || '—' }}</p>
        <p><strong>截止时间：</strong>{{ detailData.deadline }}</p>
        <p><strong>总分：</strong>{{ detailData.total_score }} 分</p>
        <el-divider />
        <div v-for="(q, i) in detailData.questions" :key="q.id" class="question-card">
          <div class="q-header">
            <span class="q-num">第 {{ i + 1 }} 题</span>
            <span class="q-type">{{ q.type === 'single' ? '单选' : q.type === 'multiple' ? '多选' : '判断' }} · {{ q.score }}分</span>
          </div>
          <div class="q-title">{{ q.title }}</div>
          <div class="q-options">
            <div v-for="opt in q.options" :key="opt.key" class="q-option">
              <span class="q-key" :class="{ correct: q.answer.includes(opt.key) }">{{ opt.key }}</span>
              {{ opt.value }}
            </div>
          </div>
          <div v-if="q.explanation" class="q-explanation">解析：{{ q.explanation }}</div>
        </div>
      </div>
    </el-dialog>

    <!-- 提交列表弹窗 -->
    <el-dialog v-model="submissionVisible" title="提交列表" width="700px" destroy-on-close>
      <el-table :data="submissions" stripe>
        <el-table-column prop="username" label="学生" />
        <el-table-column label="得分" width="120" align="center">
          <template #default="{ row }">{{ row.score }} / {{ row.total_score }}</template>
        </el-table-column>
        <el-table-column label="正确率" width="100" align="center">
          <template #default="{ row }">
            <span :style="{ color: getCorrectRate(row) >= 60 ? '#4caf50' : '#f56c6c', fontWeight: 600 }">
              {{ getCorrectRate(row) }}%
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" />
        <el-table-column prop="submitted_at" label="提交时间" width="170" />
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getAssignments, deleteAssignment, publishAssignment, closeAssignment,
  getAssignment, getSubmissions
} from '@/services/assignment'
import AssignmentDialog from './components/AssignmentDialog.vue'

const assignments = ref<any[]>([])
const loading = ref(false)
const page = ref(1)
const total = ref(0)

const dialogVisible = ref(false)
const editingAssignment = ref<any>(null)

const detailVisible = ref(false)
const detailData = ref<any>(null)

const submissionVisible = ref(false)
const submissions = ref<any[]>([])

const fetchAssignments = async () => {
  loading.value = true
  try {
    const res = await getAssignments({ page: page.value, size: 10 })
    assignments.value = res.list || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

const openCreateDialog = () => {
  editingAssignment.value = null
  dialogVisible.value = true
}

const openEditDialog = (row: any) => {
  editingAssignment.value = { ...row }
  dialogVisible.value = true
}

const viewDetail = async (row: any) => {
  const res = await getAssignment(row.id)
  detailData.value = res
  detailVisible.value = true
}

const viewSubmissions = async (row: any) => {
  loading.value = true
  try {
    const res = await getSubmissions(row.id)
    submissions.value = res.list || []
    submissionVisible.value = true
  } finally {
    loading.value = false
  }
}

const getCorrectRate = (row: any) => {
  if (!row.total_score || row.total_score === 0) return 0
  return Math.round((row.score / row.total_score) * 100)
}

const handlePublish = async (row: any) => {
  await ElMessageBox.confirm(`确定发布作业「${row.title}」？`, '发布确认')
  await publishAssignment(row.id)
  ElMessage.success('已发布')
  fetchAssignments()
}

const handleClose = async (row: any) => {
  await ElMessageBox.confirm(`确定关闭作业「${row.title}」？关闭后学生将无法提交。`, '关闭确认')
  await closeAssignment(row.id)
  ElMessage.success('已关闭')
  fetchAssignments()
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm(`确定删除作业「${row.title}」？`, '删除确认', { type: 'warning' })
  await deleteAssignment(row.id)
  ElMessage.success('已删除')
  fetchAssignments()
}

const onSaved = () => {
  fetchAssignments()
}

onMounted(() => {
  fetchAssignments()
})
</script>

<style scoped>
.assignments-page { padding: 0; }
.page-header {
  display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px;
}
.page-header h2 { margin: 0; font-size: 20px; }
.question-card {
  border: 1px solid #e0e0e0; border-radius: 8px; padding: 14px; margin-bottom: 10px;
}
.q-header { display: flex; justify-content: space-between; margin-bottom: 6px; }
.q-num { font-weight: 600; color: #1a73e8; font-size: 13px; }
.q-type { font-size: 12px; color: #999; }
.q-title { font-size: 14px; margin-bottom: 6px; }
.q-options { padding-left: 8px; }
.q-option { font-size: 13px; padding: 2px 0; display: flex; align-items: center; gap: 6px; }
.q-key {
  width: 22px; height: 22px; border-radius: 50%; background: #eee;
  display: inline-flex; align-items: center; justify-content: center;
  font-size: 12px; font-weight: 600; color: #666; flex-shrink: 0;
}
.q-key.correct { background: #4caf50; color: #fff; }
.q-explanation { font-size: 12px; color: #999; margin-top: 6px; padding-top: 6px; border-top: 1px dashed #eee; }
</style>
