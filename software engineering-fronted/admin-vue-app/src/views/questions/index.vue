<template>
  <div class="questions-page">
    <div class="page-header">
      <h2>题目管理</h2>
      <el-button type="primary" :icon="Plus" @click="handleCreate">新增题目</el-button>
    </div>

    <!-- 筛选栏 -->
    <div class="filter-bar page-card">
      <el-input
        v-model="keyword"
        placeholder="搜索题目"
        :prefix-icon="Search"
        clearable
        style="width: 280px"
        @keyup.enter="fetchQuestions"
      />
      <el-select v-model="typeFilter" placeholder="题目类型" clearable style="width: 140px">
        <el-option label="单选题" value="single" />
        <el-option label="多选题" value="multiple" />
        <el-option label="判断题" value="judge" />
        <el-option label="填空题" value="fill" />
      </el-select>
      <el-button type="primary" @click="fetchQuestions">查询</el-button>
    </div>

    <!-- 题目列表 -->
    <div class="page-card">
      <el-table :data="questions" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="题目" min-width="300" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ getTypeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="answer" label="答案" width="120" show-overflow-tooltip />
        <el-table-column prop="knowledge_point_name" label="关联知识点" width="140">
          <template #default="{ row }">
            <span v-if="row.knowledge_point_name">{{ row.knowledge_point_name }}</span>
            <span v-else class="text-muted">--</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper" v-if="total > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="fetchQuestions"
          @current-change="fetchQuestions"
        />
      </div>
    </div>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑题目' : '新增题目'"
      width="650px"
    >
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="题目类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择题目类型">
            <el-option label="单选题" value="single" />
            <el-option label="多选题" value="multiple" />
            <el-option label="判断题" value="judge" />
            <el-option label="填空题" value="fill" />
          </el-select>
        </el-form-item>
        <el-form-item label="题目内容" prop="title">
          <el-input
            v-model="formData.title"
            type="textarea"
            :rows="3"
            placeholder="请输入题目内容"
          />
        </el-form-item>
        <el-form-item
          v-if="formData.type === 'single' || formData.type === 'multiple'"
          label="选项"
        >
          <div class="options-list">
            <div v-for="(opt, index) in formData.options" :key="index" class="option-row">
              <span class="option-label">{{ String.fromCharCode(65 + index) }}.</span>
              <el-input v-model="formData.options[index]" placeholder="请输入选项内容" />
              <el-button
                v-if="formData.options.length > 2"
                type="danger"
                :icon="Delete"
                circle
                size="small"
                @click="removeOption(index)"
              />
            </div>
            <el-button type="primary" link @click="addOption">+ 添加选项</el-button>
          </div>
        </el-form-item>
        <el-form-item label="正确答案" prop="answer">
          <el-input v-model="formData.answer" placeholder="请输入正确答案" />
        </el-form-item>
        <el-form-item label="题目解析">
          <el-input
            v-model="formData.analysis"
            type="textarea"
            :rows="3"
            placeholder="请输入题目解析（可选）"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Plus, Search, Delete } from '@element-plus/icons-vue'
import {
  getQuestions,
  createQuestion,
  updateQuestion,
  deleteQuestion,
} from '@/services/admin'

const loading = ref(false)
const questions = ref<Record<string, unknown>[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const keyword = ref('')
const typeFilter = ref('')

const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const submitLoading = ref(false)
const formRef = ref<FormInstance>()

const formData = reactive({
  type: 'single',
  title: '',
  options: ['', '', '', ''],
  answer: '',
  analysis: '',
})

const formRules = {
  type: [{ required: true, message: '请选择题目类型', trigger: 'change' }],
  title: [{ required: true, message: '请输入题目内容', trigger: 'blur' }],
  answer: [{ required: true, message: '请输入正确答案', trigger: 'blur' }],
}

function getTypeLabel(type: string) {
  const map: Record<string, string> = {
    single: '单选题',
    multiple: '多选题',
    judge: '判断题',
    fill: '填空题',
  }
  return map[type] || type
}

function addOption() {
  formData.options.push('')
}

function removeOption(index: number) {
  formData.options.splice(index, 1)
}

async function fetchQuestions() {
  loading.value = true
  try {
    const data = await getQuestions({
      page: currentPage.value,
      size: pageSize.value,
      keyword: keyword.value || undefined,
      type: typeFilter.value || undefined,
    }) as Record<string, unknown>
    questions.value = (data.records as Record<string, unknown>[]) || []
    total.value = (data.total as number) || 0
  } catch (error) {
    console.error('获取题目列表失败:', error)
  } finally {
    loading.value = false
  }
}

function handleCreate() {
  isEdit.value = false
  editId.value = null
  formData.type = 'single'
  formData.title = ''
  formData.options = ['', '', '', '']
  formData.answer = ''
  formData.analysis = ''
  dialogVisible.value = true
}

function handleEdit(row: Record<string, unknown>) {
  isEdit.value = true
  editId.value = row.id as number
  formData.type = (row.type as string) || 'single'
  formData.title = (row.title as string) || ''
  formData.options = (row.options as string[]) || ['', '', '', '']
  formData.answer = (row.answer as string) || ''
  formData.analysis = (row.analysis as string) || ''
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      const payload = {
        title: formData.title,
        type: formData.type,
        answer: formData.answer,
        analysis: formData.analysis,
        options: (formData.type === 'single' || formData.type === 'multiple')
          ? formData.options.filter(o => o.trim())
          : undefined,
      }

      if (isEdit.value && editId.value) {
        await updateQuestion(editId.value, payload)
        ElMessage.success('更新成功')
      } else {
        await createQuestion(payload)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      fetchQuestions()
    } catch (error) {
      console.error('操作失败:', error)
    } finally {
      submitLoading.value = false
    }
  })
}

async function handleDelete(row: Record<string, unknown>) {
  try {
    await ElMessageBox.confirm(
      `确定要删除该题目吗？此操作不可恢复。`,
      '确认删除',
      { type: 'warning' }
    )
    await deleteQuestion(row.id as number)
    ElMessage.success('删除成功')
    fetchQuestions()
  } catch {
    // 用户取消
  }
}

onMounted(() => {
  fetchQuestions()
})
</script>

<style scoped>
.questions-page {
  max-width: 1200px;
}

.filter-bar {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 20px;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.options-list {
  width: 100%;
}

.option-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.option-label {
  font-weight: 500;
  color: var(--text-secondary);
  width: 24px;
}

.text-muted {
  color: var(--text-muted);
}
</style>
