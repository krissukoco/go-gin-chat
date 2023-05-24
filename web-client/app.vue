<script setup lang="ts">
const config = useAppConfig();
console.log(config.apiBaseUrl);
</script>

<template>
  <div class="_main">
    <div class="app">
      <NuxtPage />
    </div>
  </div>
</template>

<style>

</style>
