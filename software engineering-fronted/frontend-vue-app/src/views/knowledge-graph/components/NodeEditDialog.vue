<template>
  <el-dialog
    v-model="visible"
    title="编辑知识点"
    width="500px"
    @closed="resetForm"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="80px"
      v-loading="loading"
    >
      <el-form-item label="名称" prop="name">
        <el-input v-model="form.name" placeholder="请输入知识点名称" />
      </el-form-item>
      <el-form-item label="分类" prop="category">
        <el-input v-model="form.category" placeholder="请输入分类" />
      </el-form-item>
      <el-form-item label="描述" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="4"
          placeholder="请输入描述"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="loading" @click="handleSubmit">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import type { GraphNode } from '@/types/graph'
import { updateKnowledgePoint } from '@/services/graph'

const visible = defineModel<boolean>({ default: false })

const props = defineProps<{
  node: GraphNode | null
}>()

const emit = defineEmits<{
  (e: 'success'): void
}>()

const formRef = ref<FormInstance>()
const loading = ref(false)

const form = ref({
  name: '',
  category: '',
  description: ''
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入知识点名称', trigger: 'blur' },
    { max: 100, message: '名称不能超过100个字符', trigger: 'blur' }
  ],
  category: [
    { max: 50, message: '分类不能超过50个字符', trigger: 'blur' }
  ],
  description: [
    { max: 500, message: '描述不能超过500个字符', trigger: 'blur' }
  ]
}

watch(() => props.node, (node) => {
  if (node) {
    form.value = {
      name: node.name || '',
      category: node.category || '',
      description: node.description || ''
    }
  }
})

const resetForm = () => {
  formRef.value?.resetFields()
}

const handleSubmit = async () => {
  if (!formRef.value || !props.node) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      await updateKnowledgePoint(props.node!.id, {
        name: form.value.name,
        category: form.value.category,
        description: form.value.description
      })
      ElMessage.success('保存成功')
      visible.value = false
      emit('success')
    } catch (error) {
      ElMessage.error('保存失败')
    } finally {
      loading.value = false
    }
  })
}
</script>
