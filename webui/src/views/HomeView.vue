<script setup>
import { ref, onMounted } from 'vue'
import axios from '../services/axios.js'
import { reactive } from 'vue'



const errormsg = ref(null);
const loading = ref(false);
const some_data = ref(null);
const username = ref(localStorage.getItem('username')); 
const UserId = ref(localStorage.getItem('userId'));  
const selectedConversation = ref(null)
const Messages = ref([]);
const Comments = ref([]);



const Conversation = reactive({
  conversationId: '',
  Users: [],
  ChatName: ''
})

const conversations = ref([]);

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
    axios.defaults.headers.common['Authorization'] = UserId.value;

    const response = await axios.post(`/users/${UserId.value}/conversations`, {
      Name: Conversation.ChatName
    });

    
    const newConversation = {
      conversationId: response.data.Identifier.Identifier, 
      Users: response.data.Users.map(user => ({
        Name: user.Name,
        userId: user.Identifier,  // Accesso corretto
        UserPhoto: user.UserPhoto
      })),
      ChatName: response.data.Name, 
      Messages: response.data.Messages
    };

    console.log(response.data);
    conversations.value.push(newConversation);

  } catch (e) {
    alert(e.response ? e.response.data : e.message);
  }
};

const getConversation = async (id) => {
  try {
    axios.defaults.headers.common['Authorization'] = UserId.value;
    const response = await axios.get(`/users/${UserId.value}/conversations/${id}`);

    console.log("Risposta di getConversation:", response.data);

    const thisConversation = {
      conversationId: response.data.Identifier.Identifier,  
      Users: response.data.Users.map(user => ({
        Name: user.Name,
        userId: user.Identifier,  
        UserPhoto: user.UserPhoto
      })),
      ChatName: response.data.Name,
      Messages: response.data.Messages || []  
    };

    return thisConversation;

  } catch (e) {
    console.error("Errore nella chiamata a getConversation:", e);
    alert(e);
  }
};



const openConversation = async (id) => {
  loading.value = true;
  selectedConversation.value = null;

  try {
    const data = await getConversation(id);
    console.log("Dati ottenuti da getConversation:", data); 

    
    if (data && data.conversationId) {
      selectedConversation.value = data;
      Messages.value = data.Messages || [];  

      
    } else {
      console.error("La risposta non contiene un conversationId valido.");
      alert("Errore: la conversazione non è stata trovata.");
    }
  } catch (e) {
    alert(e);
  } finally {
    loading.value = false;
  }
};

const closeConversation = () => {
  selectedConversation.value = null;  // Chiudi la conversazione
  Messages.value = [];  // Pulisci i messaggi
};

const getMyConversations = async () => {
  try {
    axios.defaults.headers.common['Authorization'] = UserId.value;
    const response = await axios.get(`/users/${UserId.value}/conversations`);

    if (response.data && Array.isArray(response.data) && response.data.length > 0) {
      console.log("Conversazioni recuperate:", response.data);

      const allConversations = await Promise.all(
        response.data.map(async (conversation) => {
          const detailedConversation = await getConversation(conversation.Identifier);
          return detailedConversation;
        })
      );

      conversations.value = allConversations;

    } else {
      console.log("Nessuna conversazione trovata.");
      conversations.value = []; 
    }
  } catch (e) {
    console.error("Errore durante il recupero delle conversazioni:", e);
    alert("Si è verificato un errore durante il recupero delle conversazioni.");
  }
};
const Message = reactive({
  MessageBody: '',
  MessageId: '',
  UploaderId: '',
  Date: '',
  MediaType: '',
  ConversationId: ''
})


const sendMessage = async () => {
  try {
    axios.defaults.headers.common['Authorization'] = UserId.value;
    axios.defaults.headers.common['MediaType'] = "Text";  

    const message = {
      MessageBody: Message.MessageBody  
    };

    const response = await axios.post(`/users/${UserId.value}/conversations/${selectedConversation.value.conversationId}/messages`, message);

    console.log("Messaggio inviato con successo:", response.data);
    
    const updatedMessage = {
      MessageBody: response.data.MessageBody,
      MessageId: response.data.MessageId,
      UploaderId: response.data.UploaderId,
      Date: response.data.Date,
      MediaType: response.data.MediaType,
      ConversationId: response.data.ConversationId,
      IsRead: response.data.IsRead,
      User: {  
        Name: response.data.User.Name,
        userId: response.data.User.userId.Identifier,
        UserPhoto: response.data.User.UserPhoto
      }
    };

    
    Messages.value = [...Messages.value, updatedMessage];
    

    
    message.MessageBody = '';
    message.MessageId = '';
    message.UploaderId = '';
    message.Date = '';
    message.MediaType = '';
    message.ConversationId = '';
    message.User = '';
    Message.MessageBody = '';
    message.IsRead = '';

  } catch (e) {
    console.error(e);
    alert("Si è verificato un errore durante l'invio del messaggio.");
  }
};

const markMessageRead = async (id) => {
  try {
    axios.defaults.headers.common['Authorization'] = UserId.value;
    await axios.put(`/users/${UserId.value}/conversations/${selectedConversation.value.conversationId}/messages/${id}/read`)
  }catch(e){
    alert(e)
  }
}
const GroupName = ref("");
const showInputForConversationId = ref(null);

const toggleGroupInput = (conversationId) => {
      showInputForConversationId.value = conversationId;
      GroupName.value = '';  // Resetta il nome del gruppo ogni volta che si apre il campo di input
    };

const closeGroupInput = () => {
      showInputForConversationId.value = null; // Chiudi la barra di input
      GroupName.value = '';  // Resetta il valore del gruppo
    };

const addToGroup = async () => {
  try {
    

    axios.defaults.headers.common['Authorization'] = UserId.value;

    const response = await axios.put(`/users/${UserId.value}/group`, {
      Name: GroupName.value,
      conversationId: showInputForConversationId.value  
    });

    console.log("Dati ottenuti da getConversation:", response.data);

    const newGroup = {
        conversationId: response.data.Identifier,
        Users: response.data.Users,
        ChatName: response.data.Name,
        Messages: response.data.Messages || []  
      };

    console.log("group:", newGroup);
    conversations.value.push(newGroup);

    showInputForConversationId.value = null;
    GroupName.value = '';
  } catch (e) {
    alert(e);
  }
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
      <h1 class="h2">Welcome, {{ username }}!</h1>  
      
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
  </div>
  <div class="container mt-4">
  <div class="row">
    <div class="col-12">
      <div v-for="conv in conversations" :key="conv.conversationId">
        <div class="card mb-3">
          <div class="card-body">
            <h5 class="card-title d-flex justify-content-between align-items-center">
              {{ conv.ChatName }}
              <button class="btn btn-success btn-sm" @click="toggleGroupInput(conv.conversationId)">Add to Group</button>
            </h5>
            <div v-if="showInputForConversationId == conv.conversationId">
              <input v-model="GroupName" type="text" class="form-control mt-2" placeholder="Enter group name" />
              <div class="mt-2 d-flex justify-content-between">
                <button class="btn btn-primary btn-sm" @click="addToGroup()">Submit</button>
                <button class="btn btn-danger btn-sm ml-2" @click="closeGroupInput">Close</button>
              </div>
            </div>
            <div v-if="showInputForConversationId != conv.conversationId">
              <button class="btn btn-primary" @click="openConversation(conv.conversationId)">Open</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<div class="col-8">
  <div v-if="loading">Loading...</div>

  <div v-else-if="selectedConversation">
    <div class="col-8 d-flex flex-column" style="height: 500px; border: 1px solid #ccc;">
      <div class="p-2 border-bottom">
        <h5>{{ selectedConversation.ChatName }}
          <button class="btn btn-danger btn-sm float-end" @click="closeConversation">X</button>
        </h5>
        <div class="user-list mt-2">
         <h6>Members:</h6>
          <div class="d-flex flex-wrap">
            <span v-for="user in selectedConversation.Users" :key="user.Id" class="badge bg-primary me-2 mb-2">
              {{ user.Name }}
            </span>
          </div>
        </div>
      </div>
      <div class="flex-grow-1 overflow-auto p-2">
        <div v-for="msg in Messages" :key="msg.MessageId" class="mb-2">
        
          <div :class="{
            'd-flex justify-content-end': msg.User.Name == username, 
            'd-flex justify-content-start': msg.User.Name != username
          }">
            
            <div class="alert" :class="{
              'alert-primary': msg.User.Name == username, 
              'alert-success': msg.User.Name != username
            }">
              <strong>{{ msg.User.Name }}:</strong> {{ msg.MessageBody }}
              <div>{{ msg.Date }}</div>
              <div class="message-status" v-if="msg.User.Name === username">

                  <span v-if="msg.IsRead == 'Yes'" class="badge bg-success me-1">•</span>
                  <span v-if="msg.IsRead == 'Yes'" class="badge bg-success me-1">•</span>

                  <span v-if="msg.IsRead == 'No'" class="badge bg-secondary me-1">•</span>
              </div>
            </div>
            
          </div>
        </div>
      </div>
      <div class="p-2 border-top d-flex">
        <input v-model="Message.MessageBody" class="form-control me-2" placeholder="Write..." />
        <button class="btn btn-primary" @click="sendMessage">Send</button>
      </div>
    </div>
  </div>

      <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

    
</div>
  
</template>

<style>
</style>
