<script>
import { ref, onMounted } from 'vue';

export default {
  setup() {
    const username = ref('');

    onMounted(() => {
      username.value = localStorage.getItem('username');
    });

    return {
      username
    };
  }
}
</script>

<template>
    <div>
        <h1>Benvenuto, {{ this.username }}<</h1>
    </div>
    

</template>