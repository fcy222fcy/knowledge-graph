<template>
  <el-dialog v-model="visible" title="提交资料" width="520px" @closed="resetForm">
    <el-form :model="form" label-width="80px">
      <el-form-item label="选择文件">
        <el-upload
          ref="uploadRef"
          :auto-upload="false"
          :limit="1"
          accept=".pdf,.doc,.docx,.txt,.md"
          :on-change="handleFileChange"
          :on-exceed="handleExceed"
          drag
        >
          <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
          <div class="el-upload__text">
            将文件拖到此处，或<em>点击上传</em>
          </div>
          <template #tip>
            <div class="el-upload__tip">支持 PDF、Word、TXT、Markdown 格式</div>
          </template>
        </el-upload>
      </el-form-item>
      <el-form-item label="标题">
        <el-input v-model="form.title" placeholder="文档标题（可选）" />
      </el-form-item>
      <el-form-item label="描述">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="3"
          placeholder="文档描述（可选）"
        />
      </el-form-item>
    </el-form>
    <el-progress
      v-if="uploading"
      :percentage="uploadProgress"
      :status="uploadProgress === 100 ? 'success' : undefined"
      style="margin-top: 8px"
    />
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="uploading" :disabled="!selectedFile" @click="handleUpload">
        上传
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage, type UploadInstance, type UploadFile } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'
import { uploadDocument } from '@/services/document'

const visible = defineModel<boolean>({ default: false })
const emit = defineEmits<{ success: [] }>()

const uploadRef = ref<UploadInstance>()
const selectedFile = ref<File | null>(null)
const uploading = ref(false)
const uploadProgress = ref(0)

const form = reactive({
  title: '',
  description: ''
})

const handleFileChange = (file: UploadFile) => {
  selectedFile.value = file.raw ?? null
}

const handleExceed = () => {
  ElMessage.warning('只能上传一个文件，请先移除已选文件')
}

const handleUpload = async () => {
  if (!selectedFile.value) return
  uploading.value = true
  uploadProgress.value = 0
  try {
    await uploadDocument(
      selectedFile.value,
      form.title || undefined,
      form.description || undefined,
      (percent) => { uploadProgress.value = percent }
    )
    ElMessage.success('上传成功')
    emit('success')
    visible.value = false
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error('上传失败，请重试')
  } finally {
    uploading.value = false
  }
}

const resetForm = () => {
  selectedFile.value = null
  uploadProgress.value = 0
  form.title = ''
  form.description = ''
  uploadRef.value?.clearFiles()
}
</script>

<style scoped>
.el-upload__tip {
  color: var(--text-secondary);
  font-size: 12px;
}
</style>
