<template>
  <div id="profile">
    <v-gravatar v-bind:email="user.email" class="avatar" :size="200"/>
    <div class="content">
      <h1>{{user.username}}
      <el-button type="primary" v-if="!isFollowing" @click="follow(user.id)"><font-awesome-icon icon="user-plus" /> Follow</el-button>
      <el-button type="danger" v-else @click="unfollow(user.id)"><font-awesome-icon icon="user-times" /> Unfollow</el-button>
      </h1>
      <div class="profile-info">
        <p v-if="user.followers === 0">No followers ! <font-awesome-icon icon="frown" /></p>
        <p v-else>Followers: {{user.followers}} <font-awesome-icon icon="smile" /></p>
      </div>
    </div>
  </div>
</template>
<script>
import auth from '@/auth'
import router from '@/router'
export default {
  props: ['id'],
  data () {
    return {
      user: undefined,
      isFollowing: false
    }
  },
  mounted: function () {
    if (auth.getUserClaims(localStorage.getItem('token')).id === parseInt(this.id)) {
      router.push('/profile')
    }
    this.$http.get('http://localhost:1323/user/' + this.id).then((data) => {
      console.log(data)
      this.user = data.body
    })
    this.$http.get('http://localhost:1323/user/' + this.id + '/following').then((data) => {
      console.log(data)
      this.isFollowing = data.body
    })
  },
  methods: {
    follow (userId) {
      this.$http.post('http://localhost:1323/follow/' + userId).then((data) => {
        console.log(data)
        this.user.followers++
        this.isFollowing = true
      })
    },
    unfollow (userId) {
      this.$http.post('http://localhost:1323/unfollow/' + userId).then((data) => {
        console.log(data)
        this.user.followers--
        this.isFollowing = false
      })
    }
  }
}
</script>
<style lang="scss" scoped>
#profile{
  margin-top: 2rem;
  display: flex;
  align-items: flex-start;

  .el-button{
    vertical-align: middle;
  }

  .avatar{
    border: 10px solid white;
    box-shadow: 0px 0px 10px 0px rgba(0,0,0,0.25);
  }

  .content{
    background: white;
    flex-grow: 1;
    margin-left: 1rem;
    border: 0.125rem dashed #e5e5ee;
    .profile-info{
      p{
        margin: 0;
        height: 20vh;
        line-height: 20vh;
        text-align: center;
        font-size: 3rem;
        font-family: "Nunito";
        color: #e5e5ee;
        cursor: default;
      }
      min-height: 20vh;
      height: 20vh;
    }

    h1{
      small{
        color: #bbb;
        font-size: 2rem;
        font-variant: all-small-caps;
      }
      text-align: center;
      padding-left: 2rem;
      border-bottom: 0.125rem dashed #e5e5ee;
      font-size: 3rem;
      text-transform: capitalize;
      color: #409EFF;
      font-family: "Nunito";
      margin: 0;
    }
  }
}
</style>

