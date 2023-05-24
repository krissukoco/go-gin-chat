<script setup lang="ts">
import { RegistrationRequest, User } from '~/types/user';

const mode: Ref<string> = ref('login');
const token: Ref<string|null> = useSessionToken();
const user: Ref<User|null> = useUser();
const config = useAppConfig();

async function onLogin(username: string, password: string) {
  const response = await fetch(`${config.apiBaseUrl}/auth/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ username, password })
  });

  const data = await response.json();
  console.log(data);
  token.value = data.token;
  user.value = data.user;
  // Redirect to chat page
  useRouter().push('/chat');
}

async function onRegister(request: RegistrationRequest) {
  const response = await fetch(`${config.apiBaseUrl}/auth/register`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(request)
  });

  const data = await response.json();
  console.log(data);
  token.value = data.token;
  user.value = data.user;
  // Redirect to chat page
  if (response.ok) {
    mode.value = 'login'
  }
}

</script>

<template>
  <div class="main">
    <h1>Welcome to Chat App</h1>
    <h3>Created By: <a href="https://github.com/krissukoco" target="_blank">Kris Sukoco</a></h3>
    <h3><a href="https://github.com/krissukoco/go-gin-chat" target="_blank">Github Repo</a></h3>
    <LoginForm 
      v-if="mode === 'login'" 
      @login="onLogin"
      @switch="mode = 'register'"
    />
    <RegisterForm 
      v-else
      @register="onRegister"
      @switch="mode = 'login'"
    />
  </div>
</template>

<style scoped>

.main {
  display: flex;
  flex-direction: column;
  width: 100%;
  justify-content: center;
  align-items: center;
  /* height: 50vh; */
  margin: 20% 0;
  background-color: #120323;
}
</style>