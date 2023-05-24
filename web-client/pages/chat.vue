<script setup lang="ts">
const token = useSessionToken();
const user = useUser();
console.log("CHAT TOKEN: ", token.value);

</script>

<template>
  <div class="chat-main">
    <div class="list">
      <h3>Hi {{ user?.name  }}!</h3>
      <div class="default-input">
        <input type="text">
      </div>
    </div>
  </div>
</template>

<style scoped>
.chat-main {
  display: grid;
  grid-template-columns: 5fr 7fr;
  height: 100%;
}
.list {
  border: solid 0.5px green;
}
</style>