<script lang="ts" setup>
import Index from "./components/index.vue";
</script>

<template>
  <a-config-provider
      :theme="{
      components: {
        Radio: {
          colorPrimary: '#00b96b',
        },
      },
    }"
  >
    <Index/>
  </a-config-provider>
</template>

<style>
</style>
