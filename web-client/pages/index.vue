<script lang="ts" setup>
import { useSessionToken } from '../composables/session';

const sessionToken = useSessionToken();

onMounted(() => {
  if (!sessionToken.value) {
    useRouter().push('/login');
  }
});

console.log("Session Token: ", sessionToken.value)
</script>