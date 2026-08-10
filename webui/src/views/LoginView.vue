<script setup>
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import axios from '../services/axios.js'
import vuelogo from '@/assets/dataimages/Vue.js_Logo_2.svg.png'


const router = useRouter()

const User = reactive({
  Username: '',
  UserId: '',
  UserPhoto: ''
})


const doLogin = async () => {
  
  try {
    const response = await axios.post('/login', {
      Name: User.Username
    })

    User.UserId = response.data.Identifier

    localStorage.setItem('userId', User.UserId)
    localStorage.setItem('username', User.Username)
    
    const userResponse = await axios.get(`/users/${User.UserId}`, {
      headers: { 'Authorization': User.UserId }
    });
    User.UserPhoto = userResponse.data.UserPhoto
    localStorage.setItem('userphoto', User.UserPhoto)

    axios.defaults.headers.common['Authorization'] = User.UserId
    router.push('/home')

  } catch (e) {
    alert(e)
  }
}
</script>

<template>
  <div class="text-center">
    <div class="container mt-4 text-center">
      <img :src="vuelogo" class="img-fluid mb-3" style="max-width: 180px;">
    </div>
    <div><strong>WASATEXT</strong></div>
    <div>Login</div>
    <input v-model="User.Username" placeholder="Type here" />
    <button class="btn btn-primary" @click="doLogin">Go</button>
  </div>
</template>

<style>
  
</style>