<template>
  <div class="monitor">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>监控面板</span>
        </div>
      </template>
      <el-tabs v-model="activeTab">
        <el-tab-pane label="性能监控" name="performance">
          性能监控内容
        </el-tab-pane>
        <el-tab-pane label="告警规则" name="alerts">
          告警规则内容
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const activeTab = ref('performance')
</script> 