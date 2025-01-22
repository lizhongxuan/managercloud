<template>
  <div class="middleware-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>中间件列表</span>
          <el-button type="primary" @click="showCreateDialog">添加中间件</el-button>
        </div>
      </template>

      <el-tabs v-model="activeTab" @tab-click="handleTabClick">
        <el-tab-pane label="MySQL" name="mysql">
          <middleware-table :type="'mysql'" :data="mysqlList" />
        </el-tab-pane>
        <el-tab-pane label="Redis" name="redis">
          <middleware-table :type="'redis'" :data="redisList" />
        </el-tab-pane>
        <el-tab-pane label="PostgreSQL" name="postgresql">
          <middleware-table :type="'postgresql'" :data="pgList" />
        </el-tab-pane>
        <el-tab-pane label="ZooKeeper" name="zookeeper">
          <middleware-table :type="'zookeeper'" :data="zkList" />
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 创建中间件对话框 -->
    <el-dialog
      v-model="createDialogVisible"
      title="添加中间件"
      width="50%"
    >
      <middleware-form
        ref="createForm"
        :type="activeTab"
        @submit="handleCreate"
        @cancel="createDialogVisible = false"
      />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import MiddlewareTable from '@/components/MiddlewareTable.vue'
import MiddlewareForm from '@/components/MiddlewareForm.vue'
import { getMiddlewareList, createMiddleware } from '@/api/middleware'

const activeTab = ref('mysql')
const createDialogVisible = ref(false)
const mysqlList = ref([])
const redisList = ref([])
const pgList = ref([])
const zkList = ref([])

// 加载中间件列表
const loadMiddlewareList = async (type) => {
  try {
    const res = await getMiddlewareList({ type })
    switch(type) {
      case 'mysql':
        mysqlList.value = res.data
        break
      case 'redis':
        redisList.value = res.data
        break
      case 'postgresql':
        pgList.value = res.data
        break
      case 'zookeeper':
        zkList.value = res.data
        break
    }
  } catch (error) {
    ElMessage.error('获取中间件列表失败')
  }
}

// 处理标签页切换
const handleTabClick = () => {
  loadMiddlewareList(activeTab.value)
}

onMounted(() => {
  loadMiddlewareList(activeTab.value)
})
</script> 