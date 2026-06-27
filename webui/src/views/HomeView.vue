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
      Messages: response.data.Messages,
      IsGroup: response.data.IsGroup
    };

    
    console.log(newConversation);
    conversations.value.push(newConversation);
    Conversation.ChatName = '';

  } catch (e) {
    alert(e.response ? e.response.data : e.message);
  }
};

const getConversation = async (id) => {
  try {
    axios.defaults.headers.common['Authorization'] = UserId.value;
    const response = await axios.get(`/users/${UserId.value}/conversations/${id}`);

    

    const thisConversation = {
      conversationId: response.data.Identifier.Identifier,  
      Users: response.data.Users.map(user => ({
        Name: user.Name,
        userId: user.Identifier,  
        UserPhoto: user.UserPhoto
      })),
      ChatName: response.data.Name,
      IsGroup: response.data.IsGroup,
      Messages: response.data.Messages || []  
    };

    console.log("Risposta di getConversation:", thisConversation);
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

      await markMessageRead();

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

const markMessageRead = async () => {
  try {
    axios.defaults.headers.common['Authorization'] = UserId.value;
    await axios.put(`/users/${UserId.value}/conversations/${selectedConversation.value.conversationId}/messages/read`)
  }catch(e){
    alert(e)
  }
}

const closeConversation = () => {
  selectedConversation.value = null;  // Chiudi la conversazione
  Messages.value = [];  // Pulisci i messaggi
};

const getMyConversations = async () => {
  try {
    axios.defaults.headers.common['Authorization'] = UserId.value;
    const response = await axios.get(`/users/${UserId.value}/conversations`);

    if (response.data && Array.isArray(response.data) && response.data.length > 0) {
      

      const allConversations = await Promise.all(
        response.data.map(async (conversation) => {
          const detailedConversation = await getConversation(conversation.Identifier);
          return detailedConversation;
        })
      );

      conversations.value = allConversations;
      console.log("Conversazioni recuperate:", allConversations);
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
  ConversationId: '',
  IsRead: '',
  IsForwarded: '',
  Comments: [],
})


const sendMessage = async () => {
  try {

    axios.defaults.headers.common['Authorization'] = UserId.value;

    let response;

    if (selectedFile.value) {

      const formData = new FormData();

      formData.append("image", selectedFile.value);

      response = await axios.post(
        `/users/${UserId.value}/conversations/${selectedConversation.value.conversationId}/messages`,
        formData,
        {
          headers: {
            "Content-Type": "multipart/form-data",
            "MediaType": "image"
          }
        }
      );

    } else {

      response = await axios.post(
        `/users/${UserId.value}/conversations/${selectedConversation.value.conversationId}/messages`,
        {
          MessageBody: Message.MessageBody
        },
        {
          headers: {
            MediaType: "text"
          }
        }
      );

    }

    const updatedMessage = {
      MessageBody: response.data.MessageBody,
      MessageId: response.data.messageId,
      UploaderId: response.data.UploaderId,
      Date: response.data.Date,
      MediaType: response.data.MediaType,
      ConversationId: response.data.ConversationId,
      IsRead: response.data.IsRead,
      IsForwarded: response.data.IsForwarded,
      Comments: response.data.Comments || [],
      User: {
        Name: response.data.User.Name,
        userId: response.data.User.userId.Identifier,
        UserPhoto: response.data.User.UserPhoto
      }
    };

    Messages.value.push(updatedMessage);

    Message.MessageBody = "";
    selectedFile.value = null;

  } catch (e) {
    console.error(e);
  }
};


const showConversationModal = ref(false);
const messageToForward = ref(null);

const OpenForwardModal = async (message) => {
  messageToForward.value = message.messageId.Identifier;
  console.log("messaggio da inoltrare:", messageToForward.value)
  await getMyConversations();
  showConversationModal.value = true;
};


const forwardMessage = async (convId) => {
  try{
    axios.defaults.headers.common['MediaType'] = "text"; 
    axios.defaults.headers.common['Authorization'] = UserId.value;
    const response = await axios.post(`/users/${UserId.value}/conversations/${selectedConversation.value.conversationId}/messages/forward/${messageToForward.value}`, 
      {Identifier : convId}
    )

    alert("Message Forwarded!");
    const updatedMessage = {
      MessageBody: response.data.MessageBody,
      MessageId: response.data.messageId,
      UploaderId: response.data.UploaderId,
      Date: response.data.Date,
      MediaType: response.data.MediaType,
      ConversationId: response.data.ConversationId,
      IsRead: response.data.IsRead,
      IsForwarded: response.data.IsForwarded,
      Comments: response.data.Comments || [] ,
      User: {  
        Name: response.data.User.Name,
        userId: response.data.User.userId.Identifier,
        UserPhoto: response.data.User.UserPhoto
      }
    };

    Messages.value = [...Messages.value, updatedMessage];
    
    Message.MessageBody = '';
    Message.MessageId = '';
    Message.UploaderId = '';
    Message.Date = '';
    Message.MediaType = '';
    Message.ConversationId = '';
    Message.User = '';
    Message.MessageBody = '';
    Message.IsRead = '';
    Message.IsForwarded = '';
    Message.Comments = [];

    showConversationModal.value = false;

    messageToForward.value = null;

    await openConversation(convId);

  }catch(e){
    alert(e)
  }
}

const GroupName = ref("");
const showInputForConversationId = ref(null);

const toggleGroupInput = (conversationId) => {
      showInputForConversationId.value = conversationId;
      GroupName.value = '';  
    };

const closeGroupInput = () => {
      showInputForConversationId.value = null; 
      GroupName.value = ''; 
    };

const CreateGroup = async () => {
  try {
    axios.defaults.headers.common['Authorization'] = UserId.value;

    const response = await axios.post(`/users/${UserId.value}/group`, {
      Name: GroupName.value,
      conversationId: showInputForConversationId.value  
    });

    console.log("Dati ottenuti da getConversation:", response.data);

    const newGroup = {
      conversationId: response.data.conversationId.Identifier, 
      Users: response.data.Users.map(user => ({
        Name: user.Name,
        userId: user.Identifier,  
        UserPhoto: user.UserPhoto
      })),
      ChatName: response.data.Name, 
      Messages: response.data.Messages,
      ChatPhoto: response.data.ChatPhoto,
    };

    
    console.log("group:", newGroup);
    conversations.value.push(newGroup);

    showInputForConversationId.value = null;
    GroupName.value = '';
  } catch (e) {
    alert(e);
  }
};

const GroupModalInput = ref("")
const AddingUser = ref("")

const AddGroupModal = async (id) => {
  GroupModalInput.value = id
}

const CloseAddGroupModal = async () => {
  GroupModalInput.value = null
  AddingUser.value = null
}

const LeaveGroupInput = ref("")
const showLeaveModal = ref(false)

const LeaveGroupModal = (id) => {
  LeaveGroupInput.value = id
  showLeaveModal.value = true
}

const closeLeaveModal = () => {
  showLeaveModal.value = false
  LeaveGroupInput.value = ""
}

const LeaveGroup = async () => {
  try {
    axios.defaults.headers.common['Authorization'] = UserId.value;

    await axios.delete(`/users/${UserId.value}/group/${LeaveGroupInput.value}`)

    alert("You left the group")

    showLeaveModal.value = false
    selectedConversation.value = null;
    await getMyConversations();

  } catch (e) {
    alert(e)
  }
}

const AddToGroup = async () => {
  try {

    await axios.put(
      `/users/${UserId.value}/group/${GroupModalInput.value}`,
      {
        Name: AddingUser.value,
      }
    )

    await openConversation(selectedConversation.value.conversationId)

    AddingUser.value = ""
    GroupModalInput.value = null

  } catch (e) {

    console.log(e.response)

    alert("Errore durante l'aggiunta utente")
  }
}






const EditGroupModalInput = ref("")


const EditGroupModal = (id) => {
  EditGroupModalInput.value = id

}

const CloseEditGroupModal = () => {
  EditGroupModalInput.value = null
}

const newGroupName = ref("")

const SetGroupName = async () => {
  try{
    const response = await axios.put(`/groups/${EditGroupModalInput.value}`, {
      Name: newGroupName.value
    })
    selectedConversation.value.ChatName = newGroupName.value
    
    CloseEditGroupModal()
    await getMyConversations()

  }catch(e){
    alert(e)
  }
}

const fileInput = ref(null);

const openFilePicker = () => {
  fileInput.value.click();
};

const selectedFile = ref(null);

const onFileSelected = (event) => {
  selectedFile.value = event.target.files[0];
};

















const messageToComment = ref(null);
const showCommentModal = ref(false);

const OpenCommentModal = async (message) => {
  console.log(message);
  console.log(message.messageId);
  
  messageToComment.value = message.messageId.Identifier;
  Comments.value = message.Comments || [];
  showCommentModal.value = true;
};

const Comment = reactive({
  CommentId: '',
  MessageId: '',
  CommentBody: '',
  CommentDate: '',
  UploaderId: '',
})

const reactions = ['👍', '❤️', '😂', '😮', '😢', '😡'];

const commentMessage = async (reaction) => {
  try{
    axios.defaults.headers.common['Authorization'] = UserId.value;
    const response = await axios.post(`/users/${UserId.value}/messages/${messageToComment.value}/comments`,
      {CommentBody: reaction}
    );

    const UpdatedComment = {
      CommentId: response.data.commentId,
      MessageId: response.data.messageId,
      CommentBody: response.data.CommentBody,
      CommentDate: response.data.Date,
      UploaderId: response.data.UploaderId,
      User: {  
        Name: response.data.User.Name,
        userId: response.data.User.userId.Identifier,
        UserPhoto: response.data.User.UserPhoto
      }
    }

    

    Comments.value.push(UpdatedComment);
    console.log("Commento inviato:", UpdatedComment);

    
    const msgIndex = Messages.value.findIndex(msg => msg.MessageId === messageToComment.value);
    if (msgIndex !== -1) {
      if (!Messages.value[msgIndex].Comments) {
        Messages.value[msgIndex].Comments = [];
      }
      Messages.value[msgIndex].Comments.push(UpdatedComment);
    }

    Comment.CommentId = '';
    Comment.MessageId = '';
    Comment.CommentBody = '';
    Comment.CommentDate = '';
    Comment.UploaderId = '';
    Comment.User = '';

  }catch(e){
    alert(e)
  }
}

const uncommentMessage = async (id) => {
  try {
    axios.defaults.headers.common['Authorization'] = UserId.value;
        
    await axios.delete(`/users/${UserId.value}/messages/${messageToComment.value}/comments/${id}`);    
    Comments.value = Comments.value.filter(comment => comment.commentId !== id);
        
    const msgIndex = Messages.value.findIndex(msg => msg.messageId === messageToComment.value);
    if (msgIndex !== -1 && Messages.value[msgIndex].Comments) {
      Messages.value[msgIndex].Comments = Messages.value[msgIndex].Comments.filter(c => c.commentId !== id);
    }
    alert("message canceled")
  } catch(e) { 
    alert(e);
  }
}
const closeCommentModal = () => {
  showCommentModal.value = false;
  messageToComment.value = null;
}

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
              <button class="btn btn-success btn-sm" v-if="conv.IsGroup == 'no'" @click="toggleGroupInput(conv.conversationId)">New Group</button>
            </h5>
            <div v-if="showInputForConversationId == conv.conversationId">
              <input v-model="GroupName" type="text" class="form-control mt-2" placeholder="Enter group name" />
              <div class="mt-2 d-flex justify-content-between">
                <button class="btn btn-primary btn-sm" @click="CreateGroup()">Submit</button>
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
    <div class="col-8 d-flex flex-column" style="height: 500px; width: 1480px; border: 1px solid #ccc;">
      <div class="p-2 border-bottom">        
        <h5>{{ selectedConversation.ChatName }} 
          <button class="btn btn-secondary btn-sm" v-if="selectedConversation.IsGroup == 'yes'" @click="EditGroupModal(selectedConversation.conversationId)">✏️</button>
          <button class="btn btn-danger btn-sm" v-if="selectedConversation.IsGroup == 'yes'" @click="LeaveGroupModal(selectedConversation.conversationId)">Leave</button>
          <button class="btn btn-danger btn-sm float-end" @click="closeConversation">X</button>
          <div
            v-if="showLeaveModal"
            class="modal fade show d-block"
            tabindex="-1"
            style="background: rgba(0,0,0,0.5);"
          >
            <div class="modal-dialog modal-dialog-centered">
              <div class="modal-content">

                <div class="modal-header">
                  <h5 class="modal-title">
                    Confirm
                  </h5>

                  <button
                    type="button"
                    class="btn-close"
                    @click="closeLeaveModal"
                  ></button>
                </div>

                <div class="modal-body">
                  <p>Are you sure you want to leave this group?</p>
                </div>

                <div class="modal-footer">
                  <button
                    type="button"
                    class="btn btn-secondary"
                    @click="closeLeaveModal"
                  >
                    No
                  </button>

                  <button
                    type="button"
                    class="btn btn-danger"
                    @click="LeaveGroup"
                  >
                    Yes, Leave
                  </button>
                </div>

              </div>
            </div>
          </div>
        </h5>
        <div class="user-list mt-2">
         <h6>Members:</h6>
          <div class="d-flex flex-wrap">
            <span v-for="user in selectedConversation.Users" :key="user.Id" class="badge bg-primary me-2 mb-2">
              {{ user.Name }}
            </span>
            <button class="btn btn-success btn-sm" v-if="selectedConversation.IsGroup == 'yes'" @click="AddGroupModal(selectedConversation.conversationId)">+</button>
          </div>
        </div>
        <div v-if="EditGroupModalInput == selectedConversation.conversationId">
          <input v-model="newGroupName" />
          <button @click="SetGroupName">Salva</button>
          <button @click="CloseEditGroupModal">Annulla</button>
        </div>
        <div v-if="GroupModalInput == selectedConversation.conversationId"
              class="card shadow-sm border-0 mt-3">
              <div class="card-body p-3">

                <label class="form-label fw-semibold text-muted">
                  Add user to group
                </label>

                <div class="input-group">
                  <span class="input-group-text bg-light">
                    👤
                  </span>

                  <input
                    v-model="AddingUser"
                    type="text"
                    class="form-control"
                    placeholder="Enter username"
                  />

                  <button
                    class="btn btn-primary"
                    @click="AddToGroup"
                  >
                    Add
                  </button>
                </div>

                <div class="text-end mt-2">
                  <button
                    class="btn btn-outline-secondary btn-sm"
                   @click="CloseAddGroupModal"
                  >
                    Cancel
                  </button>
                </div>

              </div>
            </div>
      </div>
      <div class="flex-grow-1 overflow-auto p-2">
        <div v-for="msg in Messages" :key="msg.MessageId" class="mb-2">
          
          <div :class="{
            'd-flex justify-content-end': msg.User.Name == username, 
            'd-flex justify-content-start': msg.User.Name != username
          }">
            
            <div
            class="alert d-inline-block w-auto"
            style="min-width: 400px; max-width: 70%;"
            :class="{
            'alert-primary': msg.User.Name == username,
            'alert-success': msg.User.Name != username
            }">
              <button class="btn btn-success btn-sm float-end ms-1"  @click="OpenForwardModal(msg)">⇉</button>
              <button class="btn btn-secondary btn-sm float-end ms-1"  @click="OpenCommentModal(msg)">💬</button>
              <div v-if="showConversationModal"
                  class="modal fade show d-block"
                  tabindex="-1"
                  style="background-color: rgba(0,0,0,0.5); backdrop-filter: blur(5px);">
                <div class="modal-dialog modal-dialog-centered modal-md">
                  <div class="modal-content">
                    <div class="modal-header">
                      <h5 class="modal-title">
                        Select Conversation
                      </h5>
                      <button type="button"
                              class="btn-close"
                              @click="showConversationModal = false">
                      </button>
                    </div>
                    <div class="modal-body">
                      <div v-for="conv in conversations"
                          :key="conv.conversationId"
                          class="list-group mb-2">
                        <button class="list-group-item list-group-item-action"
                                @click="forwardMessage(conv.conversationId)">
                          {{ conv.ChatName }}
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div v-if="showCommentModal">
                <div class="modal fade" tabindex="-1" :class="{ show: showCommentModal }" 
                    style="display: block;" v-if="showCommentModal">
                  <div class="modal-dialog">
                    <div class="modal-content">
                      <div class="modal-header">
                        <h5 class="modal-title">Comments</h5>
                        <button type="button" class="btn-close" @click="closeCommentModal"></button>
                      </div>
                      <div class="modal-body">
                        <div class="comments-window mb-3" style="max-height: 300px; overflow-y: auto; border: 1px solid #ddd; padding: 10px; border-radius: 5px;">
                          <ul class="list-group">
                            <li class="list-group-item" v-for="c in Comments" :key="c.CommentId">
                              <strong>{{ c.User.Name }}:</strong> {{ c.CommentBody }}
                              <button class="btn btn-danger btn-sm float-end" v-if="c.User.Name == username" @click="uncommentMessage(c.CommentId)">X</button>
                            </li>
                            <li v-if="Comments.length === 0" class="list-group-item text-muted">No comments yet</li>
                          </ul>
                        </div>

                        <h6>Reactions:</h6>
                        <div class="d-flex gap-2">
                          <button 
                            v-for="reaction in reactions" 
                            :key="reaction" 
                            class="btn btn-light fs-3" 
                            @click="commentMessage(reaction)">
                            {{ reaction }}
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <button class="btn btn-danger btn-sm float-end" v-if="msg.User.Name == username" @click="deleteMessage(msg.MessageId.Identifier)">X</button>
              <strong>{{ msg.User.Name }}:</strong> {{ msg.MessageBody }}
              <div class="small text-muted">{{ msg.Date }}</div>
              <span v-if="msg.IsForwarded == 'yes'" class="badge bg-primary me-1">↠↠</span>
              <div class="message-status" v-if="msg.User.Name === username">

                  <span v-if="msg.IsRead == 'yes'" class="badge bg-success me-1">•</span>
                  <span v-if="msg.IsRead == 'yes'" class="badge bg-success me-1">•</span>

                  <span v-if="msg.IsRead == 'no'" class="badge bg-secondary me-1">•</span>
                  
              </div>
            </div>
            
          </div>
        </div>
      </div>
      <div class="p-2 border-top d-flex">
        <input v-model="Message.MessageBody" class="form-control me-2" placeholder="Write..." />
        <button class="btn btn-primary" @click="sendMessage">Send</button>
        <button class="camera-button" @click="openFilePicker">
          📷
        </button>

        <input
          ref="fileInput"
          type="file"
          hidden
          accept="image/*"
          @change="onFileSelected"
        />
      </div>
    </div>
  </div>

      <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

    
</div>
  
</template>

<style>
</style>
