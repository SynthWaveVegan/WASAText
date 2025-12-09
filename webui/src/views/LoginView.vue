<script setup>
  import ErrorMsg from '../components/ErrorMsg.vue'
  import {useRouter} from 'vue-router'
  import axios from 'axios'

  const router = useRouter();
  
  const doLogin = async () => {

    try {
      const response = await axios.post('/login');
      router.push("/home")
    } catch (e) {
      alert(e);
    }
  };
  

</script>

<template>
  <div class="d-flex justify-content-center align-items-center">
    <img src="../images/textlogo.jpg"
         class="img-fluid rounded" style="max-width: 300px;" />
  </div>

  <div class="text-center">
    <div class="text-center">Login</div>
    <input v-model="username" placeholder="Type here" />
    <button type="button" class="btn btn-primary" @click="doLogin">Go</button>
  </div>
</template>



<style>
</style>
