import axios from 'axios'
import router from '@/router'
import jwtDecode from 'jwt-decode'
export default {
  user: {
    loggedIn: false
  },

  login (context, creds, redirect) {
    var current = this
    axios({
      method: 'post',
      url: 'http://localhost:1323/login',
      crossDomain: true,
      data: creds
    })
    .then(function (response) {
      console.log(response.data['jwt'])
      localStorage.setItem('token', response.data['jwt'])
      current.user.loggedIn = true
      if (redirect) {
        router.go(redirect)
      }
    })
    .catch(function () {
      context.showError = true
    })
  },

  setupAuth () {
    var token = localStorage.getItem('token')
    if (token) {
      this.user.loggedIn = true
    } else {
      this.user.loggedIn = false
    }
  },

  getAuthHeader () {
    return {
      'Authorization': 'Bearer ' + localStorage.getItem('token'),
      'Content-Type': 'application/json'
    }
  },

  getUserClaims (token) {
    return jwtDecode(token)
  },

  logout () {
    localStorage.removeItem('token')
    this.user.loggedIn = false
    router.push('/')
  }
}
