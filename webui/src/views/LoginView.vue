<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from '../services/axios.js'

const router = useRouter()

const User = ref({
  Username: '',
  UserId: ''
})

const doLogin = async () => {
  try {
    //  LOGIN SENZA Authorization
    const response = await axios.post('/login', {
      Username: User.value.Username
    })

    //  Salvo l’ID ricevuto
    User.value.UserId = response.data.Identifier

    localStorage.setItem('userId', User.value.UserId)
    localStorage.setItem('username', User.value.Username)

    //  Imposto Authorization SOLO DOPO il login
    axios.defaults.headers.common['Authorization'] = User.value.UserId

    router.push('/home')
  } catch (e) {
    console.error(e)
    alert('Errore durante il login: ' + e.message)
  }
}
</script>

<template>
  <div class="text-center">
    <div>Login</div>
    <input v-model="User.Username" placeholder="Type here" />
    <button class="btn btn-primary" @click="doLogin">Go</button>
  </div>
</template>

<style>
</style>
