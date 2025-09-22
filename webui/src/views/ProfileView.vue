<script>
export default {
    data() {
    return {
      username: ''
    }
  },
  methods: {

  },
  mounted() {
    this.username = localStorage.getItem('username')
  }
}
</script> 

<template>
    <div>
        <h1>Benvenuto, {{ this.username }}<</h1>
    </div>
    

</template>