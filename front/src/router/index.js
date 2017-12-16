import Vue from 'vue'
import Router from 'vue-router'
import Index from '@/components/Index'
import Feed from '@/components/Feed'
import Profile from '@/components/Profile'
import User from '@/components/User'
import LoginComponent from '@/components/layout/LoginComponent'
import RegisterComponent from '@/components/layout/RegisterComponent'
import auth from '@/auth'
import 'element-ui/lib/theme-chalk/index.css'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Index',
      component: Index,
      beforeEnter: (to, from, next) => {
        if (!auth.user.loggedIn) {
          next()
        } else {
          next('/profile')
        }
      }
    },
    {
      path: '/login',
      name: 'Login',
      component: LoginComponent
    },
    {
      path: '/register',
      name: 'Register',
      component: RegisterComponent
    },
    {
      path: '/feed',
      name: 'Feed',
      component: Feed,
      beforeEnter: (to, from, next) => {
        if (auth.user.loggedIn) {
          next()
        } else {
          next(false)
        }
      }
    },
    {
      path: '/profile',
      name: 'Profile',
      component: Profile,
      beforeEnter: (to, from, next) => {
        if (auth.user.loggedIn) {
          next()
        } else {
          next(false)
        }
      }
    },
    {
      path: '/user/:id',
      name: 'User',
      props: true,
      component: User,
      beforeEnter: (to, from, next) => {
        if (auth.user.loggedIn) {
          next()
        } else {
          next(false)
        }
      }
    }
  ]
})
