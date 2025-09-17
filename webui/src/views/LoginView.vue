<script>
export default {
  data() {
    return {
      text: ''
    }
  },
  methods: {
    loadmenu() {
      if (this.text.trim() !== '') {
        this.$router.push({ name: 'home', query: { text: this.text } })
      } else {
        alert('Please insert your nickname.')
      }
    }
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
    <input v-model="text" placeholder="Type here" />
    <button type="button" class="btn btn-primary" @click="loadmenu">Go</button>
  </div>
</template>



<style>
</style>
