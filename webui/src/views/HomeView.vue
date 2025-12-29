<script setup>
import { ref, onMounted } from 'vue'
import axios from '../services/axios.js'

// Dichiara le variabili reattive
const errormsg = ref(null);
const loading = ref(false);
const some_data = ref(null);
const username = ref(localStorage.getItem('username')); 
const UserId = ref(localStorage.getItem('userId'));   // Rendiamo reattivo il valore

// Funzione per fare il refresh
const refresh = async () => {
  loading.value = true;
  errormsg.value = null;
  try {
    let response = await axios.get("/");  // Assicurati di avere accesso a $axios
    some_data.value = response.data;
  } catch (e) {
    errormsg.value = e.toString();
  }
  loading.value = false;
};

const getMyConversations = async () => {
  try {
    const response = await axios.get(`/users/${UserId.value}/conversations`, {
      
    })

    axios.defaults.headers.common['Authorization'] = UserId.value
  } catch (e) {
    alert(e)
  }
}
// Chiamata alla funzione refresh quando il componente viene montato
onMounted(() => {
  refresh();
  
});
</script>

<template>
  <div>
    <div
      class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">Welcome, {{ username }}!</h1>  <!-- Ora username è reattivo -->
      
      <div class="btn-toolbar mb-2 mb-md-0">
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
            Refresh
          </button>
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="exportList">
            Export
          </button>
        </div>
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-primary" @click="newItem">
            New
          </button>
        </div>
      </div>
      
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
  </div>
</template>

<style>
</style>
