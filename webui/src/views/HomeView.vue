<script setup>
import { ref, onMounted } from 'vue'
import axios from '../services/axios.js'
import { reactive } from 'vue'
import Chat from '../components/Chat.vue'


const errormsg = ref(null);
const loading = ref(false);
const some_data = ref(null);
const username = ref(localStorage.getItem('username')); 
const UserId = ref(localStorage.getItem('userId'));  

const Conversations = reactive({
  conversationId: '',
  UserHosting: '',
  UserConnected: '',
  ChatName: ''
})

const conversations = ref([]);

// Funzione per fare il refresh
const refresh = async () => {
  loading.value = true;
  errormsg.value = null;
  try {
    let response = await axios.get("/");  
    some_data.value = response.data;
  } catch (e) {
    errormsg.value = e.toString();
  }
  loading.value = false;
};

const createConversation = async () => {
  try {
    const response = await axios.post(`/users/${UserId}/conversations`, {
      ChatName: username
    })
    Conversations.conversationId = response.data.conversationId
    Conversations.UserHosting = response.data.UserHosting
    Conversations.UserConnected = response.data.UserConnected
    Conversations.ChatName = response.data.ChatName

  } catch(e){
    alert(e)
  }
}

const openConversation = (id) => {

  router.push(`/conversation/${id}`);
  localStorage.setItem('conversationId', id);

}


// Chiamata alla funzione refresh quando il componente viene montato
onMounted(() => {
  

  
});
</script>

<template>
  <div>
    <div
      class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">Welcome, {{ username }}!</h1>  <!-- Ora username è reattivo -->
      
      <div class="btn-toolbar mb-2 mb-md-0">
        <div>
          <input placeholder="Search User..." />
          <button class="btn btn-primary">Go</button></div>
        </div>
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
      

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

    
  </div>
  
</template>

<style>
</style>
