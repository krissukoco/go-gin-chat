<script setup lang="ts">

import { RegistrationRequest } from '~/types/user';

const username: Ref<string> = ref("");
const password: Ref<string> = ref("");
const confirmPassword: Ref<string> = ref("");
const name: Ref<string> = ref("");
const location: Ref<string> = ref("");

const emit = defineEmits(['register', 'switch']);

function onClickRegister() {
  emit('register', {
    username: username.value,
    password: password.value,
    confirmPassword: confirmPassword.value,
    name: name.value,
    location: location.value
  } satisfies RegistrationRequest);
}

</script>

<template>
  <div class="login-form">
    <input class="default-input" type="text" name="username" placeholder="Pick a username" v-model="username">
    <input class="default-input" type="text" name="name" placeholder="Your Display Name" v-model="password">
    <input class="default-input" type="password" name="password" placeholder="Password, must be at least 8 characters">
    <input class="default-input" type="password" name="password" placeholder="Confirm Password">
    <button class="default-btn" @click="() => onClickRegister()">Register</button>
    <h3>Already have an account? <a @click="() => emit('switch')">Login</a></h3>
  </div>
</template>