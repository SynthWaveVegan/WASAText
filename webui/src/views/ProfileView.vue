<script setup>
import { ref, onMounted } from 'vue'
import axios from '../services/axios.js'

// Dichiarazione delle variabili reattive

const UserPhoto = ref('')
const editing = ref(false)

const startEdit = () => {
  editing.value = true
}
const UserId = ref(localStorage.getItem('userId'));  
const username = ref(localStorage.getItem('username')); 


const setMyUsername = async (newUsername) => {
  console.log("ok 1")
  try {
    const response = await axios.put(`/users/${UserId.value}`, {
      Name: newUsername.value
    })
    console.log("ok 2")
    
    localStorage.setItem('username', newUsername.value)
    console.log("ok 3")

    axios.defaults.headers.common['Authorization'] = UserId
    console.log("ok 4")
    editing.value = false
    
    

    
  } catch (e) {
    alert(e)
  }

}
 



onMounted(() => {
    
});
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
      <input v-model="newUsername.value" class="form-control d-inline-block w-auto" />
      <button class="btn btn-sm btn-success ms-2" @click="setMyUsername">
        Save
      </button>
      <button class="btn btn-sm btn-secondary ms-1" @click="editing = false">
        Cancel
      </button>
  </div>
  
  
    

</template>