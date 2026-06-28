<script setup>
import { ref, onMounted } from 'vue'
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import axios from '../services/axios.js'


const editing = ref(false)

const startEdit = () => {
  editing.value = true
}

const newUsername = ref("")
const UserId = ref(localStorage.getItem('userId'));  
const username = ref(localStorage.getItem('username')); 
const userphoto = ref(localStorage.getItem("userphoto") || "");


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

const fileInput = ref(null);

const openFilePicker = () => {
  fileInput.value.click();
};

const selectedFile = ref(null);



const photoPath = ref("");

const updatePhoto = async () => {
  try {
    axios.defaults.headers.common['Authorization'] = UserId.value

    await axios.put(`/users/${UserId.value}/photo`, {
      Photo: photoPath.value
    });

    userphoto.value = photoPath.value;
    localStorage.setItem('userphoto', photoPath.value);

    photoPath.value = "";

  } catch (e) {
    console.error(e);
    alert("something went wrong");
  }
};

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
    <span class="fw-bold">{{ username }}'s Photo</span>
    <img :src="userphoto" alt="User photo" />
    <input
      v-model="photoPath"
      type="text"
      placeholder="Insert path"
    />

    <button @click="updatePhoto">
      Save profile pic
    </button>
    <li>you have to save the image in the images folder</li>
  </div>
  
  
    

</template>