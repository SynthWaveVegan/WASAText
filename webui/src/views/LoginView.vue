<script setup>
  import ErrorMsg from '../components/ErrorMsg.vue'
  import {useRouter} from 'vue-router'
  import axios from '../services/axios.js'
  import {ref} from 'vue' 

  const router = useRouter();

  const User = ref({
    Username: '',
    UserId: ''
  })

  const Authorization = ref(null)
  
  
  const doLogin = async () => {

    try {
      
      const response = await axios.post('/login', {
        Username: User.Username
      });
      const config = {
        headers: {
          'Content-Type': 'application/json'
        }
      }
      
      User.UserId = response.data.Identifier

      localStorage.setItem('userId', User.UserId);
      localStorage.setItem('username', User.Username);
      localStorage.setItem('Authorization', `Bearer ${User.UserId}`);

      router.push("/home")

      axios.defaults.headers['Authorization'] = `Bearer ${User.UserId}`;

    } catch (e) {
      console.error(e)
      alert('Errore durante il login: ' + e.message)
    }
  };
  //const doLogin = () => {
    
    //router.push("/home")
  //}

  
  

</script>

<template>
  <div class="text-center">
    <div class="text-center">Login</div>
    <input v-model="User.Username" placeholder="Type here" />
    <button type="button" class="btn btn-primary" @click="doLogin">Go</button>
  </div>
</template>



<style>
</style>
