<script>
import { ref, onMounted } from 'vue'
import axios from '../services/axios.js'

const username = ref(localStorage.getItem('username')); 
const userid = ref(localStorage.getItem('userId')); 
const getUser = async () => {
  try {

    const response = await axios.get(`/users/${this.userId}`)
  } catch (e) {
    console.error(e)
    alert('Error: ' + e.message)
  }
}

onMounted(() => {
  refresh();
});
</script>

<template>
    
  
    

</template>