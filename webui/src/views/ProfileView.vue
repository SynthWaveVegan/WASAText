<script setup>
import { ref, onMounted } from 'vue'
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import axios from '../services/axios.js'
import propic from '@/assets/images/istockphoto-1495088043-612x612.jpg'

const editing = ref(false)

const startEdit = () => {
  editing.value = true
}

const newUsername = ref("")
const newPhoto = ref("")
const UserId = ref(localStorage.getItem('userId'));  
const username = ref(localStorage.getItem('username')); 


const SetMyUsername = async () => {
  
  try {
    axios.defaults.headers.common['Authorization'] = UserId.value

    const response = await axios.put(`/users/${UserId.value}`, {
      Name: newUsername.value
    })
    
    localStorage.setItem('username', newUsername.value)
    username.value = ref(newUsername)  
      
    editing.value = false
    
  } catch (e) {
    alert(e)
    }
}

const SetMyPhoto = async () => {

  try{
    axios.defaults.headers.common['Authorization'] = UserId.value

    const response = await axios.put(`/users/${UserId.value}/photo`, {
      Photo: newPhoto.value
    })
      
  }catch(e) {
    alert(e)
  }
}

</script>

<template>
  <div
      class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
     
    <h1 class="h2"><strong>{{ username }}'s Profile</strong></h1> 
    
    
  </div>
  <div v-if="!editing">
    <span class="fw-bold">{{ username }}</span>
    <button class="btn btn-sm btn-outline-primary ms-2"  @click="startEdit">Change Username</button>
  </div>
  <div v-else>
      <input v-model="newUsername" class="form-control d-inline-block w-auto" />
      <button class="btn btn-sm btn-success ms-2" @click="SetMyUsername">
        Save
      </button>
      <button class="btn btn-sm btn-secondary ms-1" @click="editing = false">
        Cancel
      </button>
  </div>
  <div>
    <img :src="propic" class="img-fluid mb-3" style="max-width: 180px;" alt="change profile pic">
  </div>
  
  
    

</template>