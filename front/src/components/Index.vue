<template>
<div id="index">
  <login-component></login-component>
  <br>
  <register-component></register-component>
</div>
</template>

<script>
import LoginComponent from '@/components/layout/LoginComponent'
import RegisterComponent from '@/components/layout/RegisterComponent'
export default {
  components: {
    LoginComponent,
    RegisterComponent
  },
  name: 'index',
  data () {
    return {
      token: localStorage.getItem('token'),
      msg: 'Welcome to Your Vue.js App'
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
}
</style>
