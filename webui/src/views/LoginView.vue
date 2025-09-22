<script>
export default {
  data() {
    return {
      username: ''
    }
  },
  methods: {
    loadmenu() {
      if (this.username.trim() !== '') {
        this.$router.push({ name: 'home', query: { username: this.username } })
      } else {
        alert('Please insert your nickname.')
        this.username = localStorage.setItem('username')
      }
      
    }
  },
  mounted() {
    console.log('Componente montato!')
		
	}
}
</script>

<template>
  <div class="d-flex justify-content-center align-items-center">
    <img src="../images/textlogo.jpg"
         class="img-fluid rounded" style="max-width: 300px;" />
  </div>

  <div class="text-center">
    <div class="text-center">Login</div>
    <input v-model="username" placeholder="Type here" />
    <button type="button" class="btn btn-primary" @click="loadmenu">Go</button>
  </div>
</template>



<style>
</style>
