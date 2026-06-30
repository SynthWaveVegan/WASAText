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
const userphoto = ref(localStorage.getItem("userphoto"));


const setMyUserName = async () => {
  
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

const photoPath = ref("");

const setMyPhoto = async () => {
  try {
    axios.defaults.headers.common['Authorization'] = UserId.value

    

    const fullPath = photoPath.value;
    const normalizedPath = fullPath.replace(/\\/g, '/')


    const index = normalizedPath.indexOf("/images");

    if (index !== -1) {
      const publicPath = fullPath.substring(index);

      const baseUrl = typeof window !== 'undefined' ? window.location.origin : '';
      const fullUrl = baseUrl + publicPath;
      
      localStorage.setItem("userphoto", publicPath);

      await axios.put(`/users/${UserId.value}/photo`, {
        Photo: fullUrl
      });

      userphoto.value = publicPath;
    }

  } catch (e) {
    console.error(e);
    alert("something went wrong");
  }
};

</script>

<template>
  <div class="container py-4">
  <h1 class="mb-4">{{ username }}'s Profile</h1>
  <div class="card shadow-sm">
    <div class="card-body">
      <div class="text-center mb-4">
        <img
          :src="userphoto"
          class="rounded-circle border"
          style="width: 150px; height: 150px; object-fit: cover;"
          alt="Profile photo"
        >
      </div>
      <div class="mb-4">
        <label class="form-label fw-bold">Username</label>
        <div v-if="!editing" class="d-flex align-items-center gap-2">
          <span class="fs-5">{{ username }}</span>
          <button
            class="btn btn-outline-primary btn-sm"
            @click="startEdit">
            Change
          </button>
        </div>
        <div v-else class="input-group" style="max-width:400px;">
          <input
            v-model="newUsername"
            class="form-control"
          >
          <button
            class="btn btn-success"
            @click="setMyUserName">
            Save
          </button>
          <button
            class="btn btn-secondary"
            @click="editing = false">
            Cancel
          </button>
        </div>
      </div>
      <div>
        <label class="form-label fw-bold">
          Profile photo
        </label>
        <div class="input-group mb-2">
          <input
            v-model="photoPath"
            type="text"
            class="form-control"
            placeholder="/images/avatar.jpg"
          >
          <button
            class="btn btn-primary"
            @click="setMyPhoto">
            Save
          </button>
        </div>
        <small class="text-muted">
          Save the image inside <strong>public/images</strong> and insert the
          path like <code>/images/avatar.jpg</code>.
        </small>
      </div>
    </div>
  </div>
</div>
  
  
    

</template>