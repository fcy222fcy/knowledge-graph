<template>
  <div class="admin-system" v-loading="loading">
    <div class="page-header">
      <h2 class="page-title">系统设置</h2>
    </div>

    <el-card shadow="never">
      <template #header>
        <span>系统配置信息</span>
      </template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="服务器端口">{{ config.server_port }}</el-descriptions-item>
        <el-descriptions-item label="MinIO 端点">{{ config.minio_endpoint }}</el-descriptions-item>
        <el-descriptions-item label="MinIO 桶名">{{ config.minio_bucket }}</el-descriptions-item>
        <el-descriptions-item label="Ollama 地址">{{ config.ollama_url }}</el-descriptions-item>
        <el-descriptions-item label="Ollama 模型">{{ config.ollama_model }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card shadow="never" style="margin-top: 20px;">
      <template #header>
        <span>系统状态</span>
      </template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="AI 服务">
          <el-tag type="success">正常运行</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="存储服务">
          <el-tag type="success">正常运行</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="数据库">
          <el-tag type="success">正常运行</el-tag>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getSystemConfig } from '@/services/admin'
import type { SystemConfig } from '@/services/admin'

const loading = ref(false)
const config = ref<SystemConfig>({
  ollama_url: '',
  ollama_model: '',
  minio_endpoint: '',
  minio_bucket: '',
  server_port: ''
})

const fetchConfig = async () => {
  loading.value = true
  try {
    const res = await getSystemConfig()
    config.value = res.data
  } catch (error) {
    ElMessage.error('获取系统配置失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchConfig()
})
</script>

<style scoped>
.admin-system {
  background: #fff;
  padding: 24px;
  border-radius: 8px;
}

.page-header {
  margin-bottom: 20px;
}

.page-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}
</style>
