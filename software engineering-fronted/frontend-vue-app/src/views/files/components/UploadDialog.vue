<template>
  <el-dialog v-model="visible" title="上传文档" width="500px" destroy-on-close>
    <el-form label-width="80px">
      <el-form-item label="选择文件">
        <el-upload
          ref="uploadRef"
          drag
          :auto-upload="false"
          :limit="1"
          accept=".pdf,.docx,.pptx,.md,.txt"
          :on-change="handleFileChange"
          :on-exceed="handleExceed"
        >
          <el-icon class="el-icon--upload" :size="40"><UploadFilled /></el-icon>
          <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
          <template #tip>
            <div class="el-upload__tip">支持 .pdf, .docx, .pptx, .md, .txt，单文件最大 10MB</div>
          </template>
        </el-upload>
      </el-form-item>
      <el-form-item label="文档标题">
        <el-input v-model="title" placeholder="不填则使用文件名" />
      </el-form-item>
      <el-form-item label="描述">
        <el-input v-model="description" type="textarea" :rows="2" placeholder="可选" />
      </el-form-item>
    </el-form>
    <el-progress v-if="uploading" :percentage="uploadProgress" :status="uploadStatus" />
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="uploading" @click="handleUpload">上传</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage, genFileId } from 'element-plus'
import type { UploadFile, UploadInstance, UploadRawFile } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'
import { uploadDocument } from '@/services/document'

const visible = defineModel<boolean>({ default: false })
const emit = defineEmits<{ success: [] }>()

const uploadRef = ref<UploadInstance>()
const selectedFile = ref<File | null>(null)
const title = ref('')
const description = ref('')
const uploading = ref(false)
const uploadProgress = ref(0)

const uploadStatus = computed(() => uploadProgress.value >= 100 ? 'success' : undefined)

const handleFileChange = (file: UploadFile) => {
  selectedFile.value = file.raw || null
}

const handleExceed = (files: File[]) => {
  uploadRef.value?.clearFiles()
  const file = files[0] as UploadRawFile
  file.uid = genFileId()
  uploadRef.value?.handleStart(file)
}

const handleUpload = async () => {
  if (!selectedFile.value) {
    ElMessage.warning('请先选择文件')
    return
  }

  const maxSize = 10 * 1024 * 1024
  if (selectedFile.value.size > maxSize) {
    ElMessage.error('文件大小不能超过 10MB')
    return
  }

  uploading.value = true
  uploadProgress.value = 0

  try {
    await uploadDocument(
      selectedFile.value,
      title.value || undefined,
      description.value || undefined,
      (percent) => { uploadProgress.value = percent }
    )

    uploadProgress.value = 100
    ElMessage.success('上传成功')
    visible.value = false
    emit('success')
  } catch (error) {
    ElMessage.error('上传失败')
  } finally {
    uploading.value = false
    uploadProgress.value = 0
  }
}
</script>
