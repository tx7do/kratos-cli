<script setup lang="ts">
import {ref} from "vue";

import {OpenProject, SelectFolder} from "../../wailsjs/go/main/App";
import {detect} from "../../wailsjs/go/models";

const projectInfo = ref<detect.ProjectInfo>()

async function handleOpenProject() {
  try {
    const path = await SelectFolder();
    if (path) {
      console.log("选择的项目路径：", path);
      const pi = await OpenProject(path);
      console.log("项目信息：", pi);
      projectInfo.value = pi;
    }
  } catch (err) {
    console.error("选择文件夹出错：", err);
  }
}
</script>

<template>
  <div class="header">
    <!-- 左侧：按钮 -->
    <div class="header-left">
      <a-button type="primary" @click="handleOpenProject">打开项目</a-button>
    </div>

    <!-- 右侧：项目信息（横向排列） -->
    <div class="header-right" v-if="projectInfo">
      <div class="info-item">
        <span class="info-label">项目</span>
        <span class="info-value">{{ projectInfo.ModPath }}</span>
      </div>
      <div class="info-item">
        <span class="info-label">Go版本</span>
        <span class="info-value">{{ projectInfo.GoVersion }}</span>
      </div>
      <div class="info-item">
        <span class="info-label">服务数</span>
        <span class="info-value">{{ projectInfo.Services?.length ?? 0 }}</span>
      </div>
      <div class="info-item">
        <span class="info-label">有API</span>
        <span class="info-value">{{ projectInfo.HasApi ? '有' : '无' }}</span>
      </div>
    </div>

    <!-- 未打开项目时的提示 -->
    <div class="header-right header-prompt" v-else>
      <span class="prompt-text">您需要首先打开一个微服务项目</span>
    </div>
  </div>
</template>

<style scoped>
.header {
  height: 60px;
  width: 100%;
  max-width: 100vw;
  background: #fff;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  align-items: center;
  justify-content: space-between; /* 按钮和信息分两侧 */
  padding: 0 24px;
  box-sizing: border-box;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 24px; /* 信息项之间的间距 */
  font-size: 13px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.info-label {
  color: #8c8c8c;
  font-size: 12px;
}

.info-value {
  color: #262626;
  font-weight: 500;
}

.header-prompt {
  gap: 0;
}

.prompt-text {
  color: #8c8c8c;
  font-size: 13px;
  font-style: italic;
}
</style>
