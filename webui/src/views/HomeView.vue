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

const Conversation = reactive({
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
    axios.defaults.headers.common['Authorization'] = UserId.value

    const response = await axios.post(`/users/${UserId.value}/conversations`, {
      Name: Conversation.ChatName
    })
    
    const newConversation = {
        conversationId: response.data.Identifier.Identifier,
        UserHosting: response.data.userHosting.Identifier,
        UserConnected: response.data.UserConnected.Identifier,
        ChatName: response.data.Name
      };

      console.log(response.data)
      // Aggiunta della nuova conversazione alla lista
      
      conversations.value.push(newConversation);

  } catch(e){
    alert(e.response ? e.response.data : e.message)
  }
}

const openConversation = (id) => {

  

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
          <input v-model="Conversation.ChatName" placeholder="Search User..." />
          <button class="btn btn-primary" @click="createConversation">Go</button></div>
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
      
    <div class="container mt-4">
      <div class="row">
        <div class="col-12">
          <div v-for="conv in conversations">
            <div class="card mb-3">
              <div class="card-body">
                <h5 class="card-title" @click="openConversation(conv.conversationId)">
                  {{ conv.ChatName }}
                </h5>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

    
  </div>
  
</template>

<style>
</style>
