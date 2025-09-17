<script setup>
    import { ref } from 'vue'
    import { useRouter } from 'vue-router';

    const text = ref('')
    const router = useRouter()


    function loadmenu() {
        if (text.value.trim() !== '') {
            router.push({ name: 'home', query: { text: text.value} });
        } else {
        alert('Per favore inserisci un testo prima di continuare.');
      }
        
    }
</script>

<template>
    <div class="d-flex">
        <div class="flex-shrink-0">
            <img src="../images/textlogo.jpg">
        </div>
    </div>
    
    <div class="text-center">
        
        <div class="text-center">login</div>
        <input v-model="text" placeholder="Type here"></input>
        <button type="button" class="btn btn-primary" @click="loadmenu()">Go</button>
    </div>
    
</template>

<style>
</style>