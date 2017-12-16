<template>
  <el-menu class="el-menu" mode="horizontal" active-text-color="#409EFF">
    <h1 class="navbar-brand">Twittos</h1>
    <el-menu-item index="1" v-if="user.loggedIn" style="line-height:75px;height:75px;" @click="handleLogout()"><font-awesome-icon icon="sign-out-alt" /> Logout</el-menu-item>
    <el-menu-item index="2"  v-if="user.loggedIn" style="float: right;line-height:75px;height:75px;margin-right: 3rem;">
        <el-autocomplete
          class="inline-input"
          v-model="formSearch.search"
          :fetch-suggestions="querySearch"
          placeholder="Search user or email..."
          :trigger-on-focus="false"
          @select="handleSelect"
        >
          <i slot="prefix" class="el-input__icon el-icon-search"></i>
        </el-autocomplete>
    </el-menu-item>
  </el-menu>
</template>

<script>
import auth from '@/auth'
import router from '@/router'
import FontAwesomeIcon from '@fortawesome/vue-fontawesome'
export default {
  name: 'header',
  data () {
    return {
      user: auth.user,
      formSearch: {
        search: ''
      },
      timeout: null,
      queried: []
    }
  },
  methods: {
    handleLogout () {
      auth.logout()
    },
    querySearch (queryString, cb) {
      this.queryBack(queryString)
      // call callback function to return suggestions
      clearTimeout(this.timeout)
      this.timeout = setTimeout(() => {
        var array = []
        this.queried.forEach((item) => {
          array.push({value: item.username, id: item.id})
        })
        cb(array)
      }, 3000 * Math.random())
    },
    createFilter (queryString) {
      return (link) => {
        return (link.username.toLowerCase().includes(queryString.toLowerCase()))
      }
    },
    queryBack (str) {
      this.$http.get('http://localhost:1323/queryUsers/' + str).then((data) => {
        this.queried = data.body
      })
    },
    handleSelect (item) {
      router.push('/user/' + item.id)
    }
  },
  components: {
    FontAwesomeIcon
  }
}
</script>

<style lang="scss" scoped>
h1.navbar-brand{
  line-height: 75px;
  margin-top: 0;
  margin-bottom: 0;
  margin-left: 2rem;
  margin-right: 2rem;
  float: left;
  color: #409EFF;
  font-family: "Nunito";
  font-size: 1.7rem;
  font-weight: bold;
  cursor: default;
}

.el-menu-item{
  font-weight: bold;
  font-family: "Nunito";
  font-size: 1rem;
  &:hover{
    color: #409EFF;
  }
}

</style>
