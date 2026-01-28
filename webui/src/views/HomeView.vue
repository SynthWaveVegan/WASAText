<script setup>
import { ref, onMounted } from 'vue'
import axios from '../services/axios.js'
import { reactive } from 'vue'


const errormsg = ref(null);
const loading = ref(false);
const some_data = ref(null);
const username = ref(localStorage.getItem('username')); 
const UserId = ref(localStorage.getItem('userId'));  

const Conversation = reactive({
  
  Id: '',

  Name: ''
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


const openConversation = (id) => {

  router.push(`/conversation/${id}`);
  localStorage.setItem('conversationId', id);

}


const createConversation = async () => {
  axios.defaults.headers.common['Authorization'] = UserId.value

  try {
    const response = await axios.post(
      `/users/${UserId.value}/conversations`,
      {
        name: Conversation.Name
      }
    )

    Conversation.Id = response.data.Identifier.Identifier

    conversations.value.push({
      Id: Conversation.Id,
      Name: Conversation.Name
    });

    localStorage.setItem('conversationId', Conversation.Id)

    Conversation.Name = '';

  } catch (e) {
    console.error(e)
    alert('Errore nella creazione della conversazione')
  }
}

const getConversation = async () => {
  try {
    const response = await axios.get(`/users/${UserId.value}/conversations/${ConversationId.value}`, {

    })
  } catch (e) {
    alert(e)
  }
}


const getMyConversations = async () => {
  loading.value = true;
  errormsg.value = null;

  try {
    // Effettua la chiamata GET per ottenere le conversazioni
    const response = await axios.get(`/users/${UserId.value}/conversations`);
    
    // Salva la lista delle conversazioni nella variabile reattiva
    conversations.value = response.data;  // Assegna la risposta alla lista delle conversazioni

    // Imposta l'Authorization per tutte le richieste future
    axios.defaults.headers.common['Authorization'] = UserId.value;

    console.log('Conversations loaded:', conversations.value); // Log per verificare la risposta
  } catch (e) {
    // Gestione degli errori
    console.error(e);
    errormsg.value = 'Errore nel recuperare le conversazioni';
  }

  loading.value = false;
};
// Chiamata alla funzione refresh quando il componente viene montato
onMounted(() => {
  
  getMyConversations();
  
});
</script>

<template>
  <div>
    <div
      class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">Welcome, {{ username }}!</h1>  <!-- Ora username è reattivo -->
      
      <div class="btn-toolbar mb-2 mb-md-0">
        <div>
          <input v-model="Conversation.Name" placeholder="Search User..." />
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
      

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

    <div class="row mt-4">
      <div 
        v-for="Conversation in conversations" 
        :key="Conversation.id"
        class="col-sm-12 col-md-4 col-lg-3 mb-4"
      >
        <div 
          class="card shadow-sm p-3 mb-5 bg-white rounded"
          style="cursor: pointer;"
          @click="openConversation(Conversation.id)"
        >
          <div class="card-body">
            <h5 class="card-title text-center">{{ Conversation.Name }}</h5>
          </div>
        </div>
      </div>
    </div>
  </div>
  
</template>

<style>
</style>
