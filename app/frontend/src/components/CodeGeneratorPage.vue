<script setup lang="ts">
import {ref, reactive} from 'vue'

import {EditGeneratorOption, GetGeneratorOptions, GetProjectInfo, SetGeneratorOption} from "../../wailsjs/go/main/App";
import {generator} from "../../wailsjs/go/models";
import {EventsOn} from "../../wailsjs/runtime";

import DatabaseImporterModal from "./DatabaseImporterModal.vue";
import SqlImporterModal from "./SqlImporterModal.vue";
import GRPCCodeGenerateModal from "./GRPCCodeGenerateModal.vue";
import RESTCodeGenerateModal from "./RESTCodeGenerateModal.vue";
import FrontendCodeGenerateModal from "./FrontendCodeGenerateModal.vue";

const openDatabaseImporter = ref<boolean>(false);
const openSqlImporter = ref<boolean>(false);

const grpcCodeGenerateImporter = ref<boolean>(false);
const restCodeGenerateImporter = ref<boolean>(false);
const frontendCodeGenerateImporter = ref<boolean>(false);

// 快速选择服务
const quickSelectService = ref<string>('');

// 表格数据
const tableData = ref<Array<{ id: number; tableName: string; service: string; exclude: boolean }>>([])

// 服务选项
const serviceOptions = reactive<Array<{ label: string; value: string }>>([])

function handleGenerateGRPCCode() {
  grpcCodeGenerateImporter.value = true;
}

function handleGenerateRESTCode() {
  restCodeGenerateImporter.value = true;
}

function handleGenerateFrontendCode() {
  frontendCodeGenerateImporter.value = true;
}

function handleDatabaseImport() {
  openDatabaseImporter.value = true;
}

function handleSQLImport() {
  openSqlImporter.value = true;
}

// 一键选中所有表的某个服务
async function handleQuickSelectService(service: string) {
  tableData.value.forEach(row => {
    row.service = service;
  });

  const opts = await GetGeneratorOptions();
  for (let i = 0; i < opts.length; i++) {
    opts[i].service = service;
  }
  await SetGeneratorOption(opts);

  quickSelectService.value = '';
}

/**
 * 处理服务选择变更
 */
async function handleServiceChange(row: generator.Option) {
  console.log('服务已更改：', row);
  await EditGeneratorOption(row);

  const opts = await GetGeneratorOptions();
  console.log('切换服务：', opts);
}

/**
 * 处理排除状态变更
 */
async function handleExcludeChange(row: generator.Option) {
  console.log('排除状态已更改：', row);
  await EditGeneratorOption(row);

  const opts = await GetGeneratorOptions();
  console.log('切换排除状态：', opts);
}

/**
 * 刷新服务选项列表
 */
async function refreshServiceOptions() {
  const pi = await GetProjectInfo();
  if (pi && pi.Services) {
    serviceOptions.length = 0; // 清空现有选项
    pi.Services.forEach(service => {
      serviceOptions.push({label: service, value: service});
    });
  }
}

async function refreshTableData() {
  const opts = await GetGeneratorOptions();
  console.log('刷新表数据，当前选项：', opts);

  tableData.value = []; // 清空现有数据
  // if (!opts || !opts.Tables) {
  //   return;
  // }

  tableData.value = opts;

  console.log('刷新表数据完成：', tableData);
}

EventsOn('project-opened', () => {
  console.log("project-opened");
  refreshServiceOptions();
})
EventsOn('table-imported', () => {
  console.log("table-imported");
  refreshTableData();
})
</script>

<template>
  <div class="code-generator-container">
    <a-card title="代码生成" class="full-card">
      <template #extra>
        <a-space>
          <a-button type="primary" @click="handleDatabaseImport">数据库导入</a-button>
          <a-button type="primary" @click="handleSQLImport">SQL导入</a-button>
          <a-button type="primary" danger @click="handleGenerateGRPCCode">生成gRPC代码</a-button>
          <a-button type="primary" danger @click="handleGenerateRESTCode">生成REST代码</a-button>
          <a-button type="primary" danger @click="handleGenerateFrontendCode">生成前端代码</a-button>
        </a-space>
      </template>

      <vxe-table
          :data="tableData"
          :row-config="{ keyField: 'id' }"
          size="small"
          class="table-content"
      >
        <vxe-column field="tableName" title="表名" width="40%"/>
        <vxe-column field="service" title="服务" width="30%">
          <template #header>
            <div class="service-header">
              <span>服务</span>
              <a-select
                  v-model:value="quickSelectService"
                  :options="serviceOptions"
                  placeholder="一键全选"
                  style="width: 150px; margin-left: 8px"
                  @change="handleQuickSelectService"
                  allow-clear
              />
            </div>
          </template>
          <template #default="{ row }">
            <a-select
                v-model:value="row.service"
                :options="serviceOptions"
                placeholder="选择服务"
                style="width: 100%"
                @change="handleServiceChange(row)"
            />
          </template>
        </vxe-column>
        <vxe-column field="exclude" title="排除" width="30%" align="center">
          <template #default="{ row }">
            <a-switch
                v-model:checked="row.exclude"
                :style="{ backgroundColor: row.exclude ? '#ff4d4f' : undefined }"
                @change="handleExcludeChange(row)"
            />
          </template>
        </vxe-column>
      </vxe-table>
    </a-card>
  </div>
  <DatabaseImporterModal
      v-model:open="openDatabaseImporter"
  />
  <SqlImporterModal
      v-model:open="openSqlImporter"/>
  <GRPCCodeGenerateModal
      v-model:open="grpcCodeGenerateImporter"/>
  <RESTCodeGenerateModal
      v-model:open="restCodeGenerateImporter"/>
  <FrontendCodeGenerateModal
      v-model:open="frontendCodeGenerateImporter"/>
</template>

<style scoped>
.code-generator-container {
  width: 100%;
  height: 100%;
  padding: 0;
  margin: 0;
  box-sizing: border-box;
}

.full-card {
  width: 100%;
  height: 100%;
  box-sizing: border-box;
}

:deep(.ant-card) {
  height: 100%;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
}

:deep(.ant-card-head) {
  flex-shrink: 0;
}

:deep(.ant-card-body) {
  flex: 1;
  overflow: auto;
  padding: 16px;
  box-sizing: border-box;
}

/* Switch 排除状态样式 */
:deep(.ant-switch-checked) {
  background-color: #ff4d4f !important;
}

.table-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
  padding: 12px;
  background-color: #fafafa;
  border-radius: 4px;
}

.table-header span {
  font-weight: 500;
  color: #333;
}

.service-header {
  display: flex;
  align-items: center;
  gap: 4px;
  width: 100%;
}
</style>
