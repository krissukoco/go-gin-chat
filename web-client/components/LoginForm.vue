<script setup lang="ts">

const username: Ref<string> = ref("");
const password: Ref<string> = ref("");

const emit = defineEmits(['login', 'switch'])

</script>

<template>
  <div class="login-form">
    <input class="default-input" type="text" placeholder="Username" v-model="username">
    <input class="default-input" type="password" placeholder="Password" v-model="password">
    <button class="default-btn" @click="() => emit('login', username, password)">Login</button>
    <h3>Don't have an account? <a @click="() => emit('switch')">Register</a></h3>
  </div>
</template>

<style scoped>
h3 {
  text-align: center;
}
</style>