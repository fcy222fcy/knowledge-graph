<template>
  <div class="admin-knowledge" v-loading="loading">
    <div class="page-header">
      <h2 class="page-title">知识点管理</h2>
      <div class="header-actions">
        <el-button type="primary" @click="fetchData">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 知识点表格 -->
    <el-card shadow="never" style="margin-bottom: 20px;">
      <template #header>
        <span>知识点列表</span>
      </template>
      <el-table :data="knowledgePoints" stripe border style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" min-width="150" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="category" label="分类" width="120" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-popconfirm
              title="确定要删除该知识点吗？"
              @confirm="handleDeletePoint(row)"
            >
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pointPage"
        v-model:page-size="pointPageSize"
        :page-sizes="[10, 20, 50]"
        :total="pointTotal"
        layout="total, sizes, prev, pager, next"
        style="margin-top: 16px; justify-content: flex-end;"
        @size-change="fetchPoints"
        @current-change="fetchPoints"
      />
    </el-card>

    <!-- 关系表格 -->
    <el-card shadow="never">
      <template #header>
        <span>知识点关系列表</span>
      </template>
      <el-table :data="relations" stripe border style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="source_name" label="源知识点" width="150" />
        <el-table-column prop="relation_type" label="关系类型" width="120" />
        <el-table-column prop="target_name" label="目标知识点" width="150" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-popconfirm
              title="确定要删除该关系吗？"
              @confirm="handleDeleteRelation(row)"
            >
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="relationPage"
        v-model:page-size="relationPageSize"
        :page-sizes="[10, 20, 50]"
        :total="relationTotal"
        layout="total, sizes, prev, pager, next"
        style="margin-top: 16px; justify-content: flex-end;"
        @size-change="fetchRelations"
        @current-change="fetchRelations"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import {
  getKnowledgePoints,
  deleteKnowledgePoint,
  getKnowledgeRelations,
  deleteKnowledgeRelation
} from '@/services/admin'
import type { KnowledgePointItem, KnowledgeRelationItem } from '@/services/admin'

const loading = ref(false)

// 知识点相关
const knowledgePoints = ref<KnowledgePointItem[]>([])
const pointTotal = ref(0)
const pointPage = ref(1)
const pointPageSize = ref(10)

// 关系相关
const relations = ref<KnowledgeRelationItem[]>([])
const relationTotal = ref(0)
const relationPage = ref(1)
const relationPageSize = ref(10)

const fetchPoints = async () => {
  try {
    const res = await getKnowledgePoints({ page: pointPage.value, size: pointPageSize.value })
    knowledgePoints.value = res.data.list
    pointTotal.value = res.data.total
  } catch (error) {
    ElMessage.error('获取知识点列表失败')
  }
}

const fetchRelations = async () => {
  try {
    const res = await getKnowledgeRelations({ page: relationPage.value, size: relationPageSize.value })
    relations.value = res.data.list
    relationTotal.value = res.data.total
  } catch (error) {
    ElMessage.error('获取关系列表失败')
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    await Promise.all([fetchPoints(), fetchRelations()])
  } finally {
    loading.value = false
  }
}

const handleDeletePoint = async (row: KnowledgePointItem) => {
  try {
    await deleteKnowledgePoint(row.id)
    ElMessage.success('删除成功')
    fetchPoints()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const handleDeleteRelation = async (row: KnowledgeRelationItem) => {
  try {
    await deleteKnowledgeRelation(row.id)
    ElMessage.success('删除成功')
    fetchRelations()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.admin-knowledge {
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
