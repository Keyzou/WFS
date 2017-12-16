<template>
  <div id="profile">
    <v-gravatar v-bind:email="user.email" class="avatar" :size="200"/>
    <div class="content">
      <h1>{{user.username}} <small>(It's you !)</small></h1>
      <div class="profile-info">
        <p v-if="user.followers === 0">No followers ! <font-awesome-icon icon="frown" /></p>
        <p v-else>Followers: {{user.followers}} <font-awesome-icon icon="smile" /></p>
      </div>
      <div class="followers">
        <h3 @click="toggleFollowersList()">
          Followers <small>({{followers.length}})</small>
          <span class="toggler" v-if="showFollowersList()">-</span>
          <span class="toggler" v-else>+</span>
        </h3>
        <div class="follow-list" v-show="showFollowersList()">
          <ul v-if="followers.length > 0" >
            <li v-bind:key="index" v-bind:post="u" v-for="(u, index) in followers"><el-button type="text" @click="redirect(u.id)" class="user">{{u.username | capitalize }}</el-button></li>
          </ul>
          <p v-else>
            None yet !
          </p>
        </div>
      </div>
      <div class="following">
        <h3 @click="toggleFollowingList()">
          Following <small>({{following.length}})</small>
          <span class="toggler" v-if="showFollowingsList()">-</span>
          <span class="toggler" v-else>+</span>
        </h3>
        <div class="follow-list" v-show="showFollowingsList()">
          <ul v-if="following.length > 0" >
            <li v-bind:key="index" v-bind:post="u" v-for="(u, index) in following"><el-button type="text" @click="redirect(u.id)" class="user">{{u.username | capitalize }}</el-button></li>
          </ul>
          <p v-else>
            None yet !
          </p>
        </div>
      </div>
      <div class="update">
        <h3 @click="toggleUpdate()">
          Update Informations
          <span class="toggler" v-if="showUpdate()">-</span>
          <span class="toggler" v-else>+</span>
        </h3>
        <el-form v-show="showUpdate()" :inline="true" ref="formUpdate" style="padding-top:1rem;"  :model="formUpdate" :rules="rules">
          <el-form-item label="Username" prop="username">
            <el-input v-model="formUpdate.username" @change="onChange('formUpdate')" placeholder="Username"></el-input>
          </el-form-item>
          <el-form-item label="Password" prop="password">
            <el-input v-model="formUpdate.password" @change="onChange('formUpdate')" type="password" placeholder="Password"></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :disabled="UpdateButton" @click="onSubmit('formUpdate')">Update</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>
<script>
import auth from '@/auth'
import router from '@/router'
export default {
  data () {
    return {
      showFollowerList: false,
      showFollowingList: false,
      showUpdateForm: false,
      user: undefined,
      formUpdate: {
        username: '',
        password: ''
      },
      followers: [],
      following: [],
      UpdateButton: true,
      rules: {
        password: [
          { min: 6, message: 'Password must be at least 6 characters long !', trigger: 'blur' }
        ],
        username: [
          { min: 3, message: 'Username must be at least 3 characters long !', trigger: 'blur' }]
      }
    }
  },
  methods: {
    onSubmit (formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.$http.post('http://localhost:1323/user/update', {username: this.formUpdate.username, password: this.formUpdate.password}).then((data) => {
            console.log(data.body)
            this.user.username = data.body.username
          })
        } else {
          return false
        }
      })
    },
    onChange (formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.UpdateButton = false
        } else {
          this.UpdateButton = true
        }
      })
    },
    showFollowingsList () {
      return this.showFollowingList
    },
    toggleFollowingList () {
      this.$nextTick(function () {
        this.showFollowingList = !this.showFollowingList
      })
    },
    showFollowersList () {
      return this.showFollowerList
    },
    toggleFollowersList () {
      this.$nextTick(function () {
        console.log(this.followers.length)
        this.showFollowerList = !this.showFollowerList
      })
    },
    showUpdate () {
      return this.showUpdateForm
    },
    toggleUpdate () {
      this.$nextTick(function () {
        this.showUpdateForm = !this.showUpdateForm
      })
    },
    redirect (id) {
      console.log(id)
      router.push('/user/' + id)
    }
  },
  beforeCreate: function () {
    this.$http.get('http://localhost:1323/user/' + auth.getUserClaims(localStorage.getItem('token')).id).then((data) => {
      console.log(data)
      this.user = data.body
    })
    this.$http.get('http://localhost:1323/user/followers').then((data) => {
      console.log(data)
      this.followers = data.body
    })
    this.$http.get('http://localhost:1323/user/following').then((data) => {
      console.log(data)
      this.following = data.body
    })
  }
}
</script>
<style lang="scss" scoped>#profile {
  margin-top: 2rem;
  display: flex;
  align-items: flex-start;
  .avatar {
    border: 10px solid white;
    box-shadow: 0px 0px 10px 0px rgba(0, 0, 0, 0.25);
  }
  .content {
    .follow-list ul {
      list-style: none;
      text-align: auto;
      column-count: 3;
      column-rule: 1px solid #e5e5ee;
      li{
        padding-bottom: 0.75rem;
        padding-top: 0.75rem;
        color: #bbb;
      }
      column-gap: 40px;
      background: white;
      margin: 0;
      padding-top: 2rem;
      padding-bottom: 2rem;
    }
    .follow-list .user{
      font-size: 1rem;
      margin: 0;
      padding:0;
    }
    .follow-list p{
      color: #e5e5ee;
      font-size: 1.5rem;
      font-style: italic;
    }
    .update {
      text-align: center;
      margin-bottom: 1rem;
    }
    .followers, .following, .update {
      h3 {
        background: #f8f9fa;
        .toggler {
          font-size: 1.5rem;
          right: 1.5rem;
          vertical-align: middle;
          position: absolute;
        }
        small{
          color: #d5d5dd;
        }
        -webkit-touch-callout: none;
        -webkit-user-select: none;
        -khtml-user-select: none;
        -moz-user-select: none;
        -ms-user-select: none;
        user-select: none;
        position: relative;
        border-top: 1px solid #e5e5ee;
        border-bottom: 1px solid #e5e5ee;
        height: 50px;
        line-height: 50px;
        margin:0;
        padding:0;
        min-height: 50px;
        color: #a5a5aa;
        &:hover {
          background: #fbfbfd;
          color: #409EFF;
          cursor: pointer;
        }
      }
      margin-bottom: 20px;
      text-align: center;
    }
    background: white;
    flex-grow: 1;
    margin-left: 1rem;
    border: 0.125rem dashed #e5e5ee;
    .profile-info {
      p {
        margin: 0;
        height: 30vh;
        line-height: 30vh;
        text-align: center;
        font-size: 3rem;
        font-family: "Nunito";
        color: #e5e5ee;
        cursor: default;
      }
      min-height: 30vh;
      height: 30vh;
    }
    h1 {
      small {
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

