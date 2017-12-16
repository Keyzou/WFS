<template>
  <el-menu
      v-bind:default-active="currentIndex"
      class="el-menu-vertical"
      active-text-color="#409EFF"
      router>
      <el-menu-item v-if="user.loggedIn" index="/feed">
        <font-awesome-icon icon="globe" />
        <span slot="title">Feed</span>
      </el-menu-item>
      <el-menu-item v-if="user.loggedIn" index="/profile">
        <font-awesome-icon icon="user" />
        <span slot="title">Profile</span>
      </el-menu-item>

      <el-menu-item v-if="!user.loggedIn" index="/login">
        <font-awesome-icon icon="sign-in-alt" />
        <span slot="title">Login</span>
      </el-menu-item>
      <el-menu-item v-if="!user.loggedIn" index="/register">
        <font-awesome-icon icon="user-plus" />
        <span slot="title">Register</span>
      </el-menu-item>
    </el-menu>
</template>

<script>
  import auth from '@/auth'
  import FontAwesomeIcon from '@fortawesome/vue-fontawesome'
  export default {
    data () {
      return {
        isCollapse: true,
        menuHandler: 'el-icon-caret-right',
        user: auth.user,
        currentIndex: '/feed'
      }
    },
    methods: {
      handleOpen (key, keyPath) {
        console.log(key, keyPath)
      },
      handleClose (key, keyPath) {
        console.log(key, keyPath)
      },
      toggleMenu () {
        this.isCollapse = !this.isCollapse
        this.menuHandler = this.menuHandler === 'el-icon-caret-right' ? 'el-icon-caret-left' : 'el-icon-caret-right'
      }
    },
    components: {
      FontAwesomeIcon
    },
    beforeCreate: function () {
      this.$nextTick(function () {
        this.currentIndex = this.$route.path
      })
    }
  }
</script>

<style lang="scss">

.el-aside{
  background: #f1f3f5;
  border-right: 0.125rem solid #f1f3f5;
  overflow: hidden;
}

.el-menu-vertical{
  border: 0;
  margin-top: 50px;
  background: #f8f9fa;
  height:100%;
}

.btn-collapse{
    margin-top: 10px;
    margin-bottom: 10px;
}  

.el-menu-vertical .el-menu-item{
    color: #697285;
    background: #f8f9fa;
    font-family: "Nunito";
    font-size: 1rem;
    font-weight: bold;
    &:not(:first-child){
      border-top: 0.125rem solid #f1f3f5;
    }

    &:last-child{
      border-bottom: 0.125rem solid #f1f3f5;
    }

    &:hover{
      background: white;
      color: #409EFF;
      svg {
        transition-property: fill;
        transition-duration: 0.2s;
        fill: #409EFF;
      }

    }    

    svg{
      transition-property: fill;
      transition-duration: 0.2s;
      vertical-align: middle;
      padding-left: 5px;
      margin-right: 5px;
      fill: mix(#878d99, #409EFF, 80%);
    }
}

.el-menu-item .fa-icon {
    vertical-align: middle;
    padding-left: 5px;
    margin-right: 10px;
    color: mix(#878d99, #409EFF, 80%);
}

.el-menu-item.is-active .fa-icon {
    color: #409EFF;
}

.el-menu-vertical:not(.el-menu--collapse) {
    width: 300px;
    
  }

</style>

