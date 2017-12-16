// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import VueResource from 'vue-resource'
import App from './App'
import router from './router'
import fontawesome from '@fortawesome/fontawesome'
import FontAwesomeIcon from '@fortawesome/vue-fontawesome'
import ElementUI from 'element-ui'
import { faUserPlus, faUser, faSignInAlt, faHome, faGlobe, faSignOutAlt, faThumbsUp, faPaperPlane, faCommentAlt, faUserTimes, faSearch, faFrown, faSmile } from '@fortawesome/fontawesome-free-solid'
import Gravatar from 'vue-gravatar'
import auth from './auth'

fontawesome.library.add(faUserPlus, faUser, faSignInAlt, faHome, faGlobe, faSignOutAlt, faThumbsUp, faPaperPlane, faCommentAlt, faUserTimes, faSearch, faFrown, faSmile)

Vue.component('v-gravatar', Gravatar)
Vue.component(FontAwesomeIcon.name, FontAwesomeIcon)
Vue.use(ElementUI)
Vue.config.productionTip = false
Vue.use(VueResource)

Vue.filter('capitalize', function (value) {
  if (!value) return ''
  value = value.toString()
  return value.charAt(0).toUpperCase() + value.slice(1)
})

auth.setupAuth()
if (auth.user.loggedIn) {
  console.log(localStorage.getItem('token'))
  Vue.http.headers.common['Authorization'] = 'Bearer ' + localStorage.getItem('token')
}

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  template: '<App/>',
  components: { App },
  render: h => h(App)
})
