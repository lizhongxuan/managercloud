<template>
  <el-form
    ref="formRef"
    :model="form"
    :rules="rules"
    label-width="120px"
  >
    <el-form-item label="名称" prop="name">
      <el-input v-model="form.name" />
    </el-form-item>
    <el-form-item label="类型" prop="type">
      <el-select v-model="form.type" placeholder="请选择中间件类型">
        <el-option label="MySQL" value="mysql" />
        <el-option label="Redis" value="redis" />
        <el-option label="PostgreSQL" value="postgresql" />
        <el-option label="ZooKeeper" value="zookeeper" />
      </el-select>
    </el-form-item>
    <el-form-item label="版本" prop="version">
      <el-input v-model="form.version" />
    </el-form-item>
    <el-form-item label="主机地址" prop="host">
      <el-input v-model="form.host" />
    </el-form-item>
    <el-form-item label="端口" prop="port">
      <el-input v-model="form.port" />
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="submitForm">确定</el-button>
      <el-button @click="$emit('cancel')">取消</el-button>
    </el-form-item>
  </el-form>
</template>

<script setup>
import { ref, defineEmits } from 'vue'

const emit = defineEmits(['submit', 'cancel'])

const form = ref({
  name: '',
  type: '',
  version: '',
  host: '',
  port: ''
})

const rules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  host: [{ required: true, message: '请输入主机地址', trigger: 'blur' }],
  port: [{ required: true, message: '请输入端口', trigger: 'blur' }]
}

const submitForm = () => {
  emit('submit', form.value)
}
</script> 