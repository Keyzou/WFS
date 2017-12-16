<template>
<div class="post" v-if="!deleted">
  <el-popover
  ref="popover4"
  placement="top"
  trigger="hover" v-on:show="getFollowing()">
    <el-button v-if="!followed && !(this.userClaims.id === this.post.author.id)" type="primary" @click="follow(post.author.id)"><font-awesome-icon icon="user-plus" /> Follow <strong>({{followers}})</strong></el-button>
    <el-button v-else-if="this.userClaims.id === this.post.author.id" type="primary" @click="goProfile()">
      Profile
    </el-button>
    <el-button v-else type="danger" @click="unfollow(post.author.id)"><font-awesome-icon icon="user-times" /> Unfollow <strong>({{followers}})</strong></el-button>
  </el-popover>
  <v-gravatar v-bind:email="post.author.email" />
  <div class="content">
    <el-button v-if="this.userClaims.id === this.post.author.id" type="text" @click="deletePost()" class="deleteButton">&times;</el-button>
    <h4><el-button type="text" @click="redirect(post.author.id)" v-popover:popover4>{{post.author.username}}</el-button> <small>| {{post.created}}</small></h4>
    <p v-html="post.content"></p>
    <div class="links">
      <el-button type="text" v-if="!postLiked" @click="likePost(post.id)">{{postLikes}} Likes <font-awesome-icon icon="thumbs-up" /></el-button>
      <el-button type="text" v-else @click="unlikePost(post.id)">{{postLikes}} Likes <font-awesome-icon icon="thumbs-up" flip="vertical"/></el-button>
    </div>
    <comment-list v-bind:comments="post.comments"></comment-list>
    <el-form :model="addCommentForm" :rules="rules" ref="addCommentForm">
      <el-form-item style="margin-bottom: 0; margin-top: 10px;" prop="content">
        <el-input type="textarea" placeholder="Your post!" v-model="addCommentForm.content"></el-input>
      </el-form-item>
      <el-input type="hidden" v-model="addCommentForm.postId" v-bind:value="post.id"></el-input>
      <el-form-item>
        <el-button type="primary" @click="submitForm('addCommentForm')" size="mini"><font-awesome-icon icon="comment-alt" /> Comment</el-button>
      </el-form-item>
    </el-form>
  </div>
</div>
</template>
<script>
import auth from '@/auth'
import router from '@/router'
import CommentList from '@/components/layout/CommentList'
export default{
  props: ['post'],
  name: 'Post',
  components: {
    CommentList
  },
  data () {
    return {
      postLiked: this.post.liked,
      postLikes: this.post.likes,
      followers: this.post.author.followers,
      followed: this.post.isFollowing,
      router: router,
      deleted: false,
      addCommentForm: {
        content: '',
        postId: this.post.id
      },
      userClaims: auth.getUserClaims(localStorage.getItem('token')),
      rules: {
        content: [
          { required: true, message: 'Content cannot be empty!' }
        ]
      }
    }
  },
  methods: {
    deletePost () {
      console.log(this.post.id)
      this.$http.delete('http://localhost:1323/deletePost/' + this.post.id, {id: this.post.id}).then((data) => {
        this.deleted = true
      })
    },
    follow (userId) {
      this.$http.post('http://localhost:1323/follow/' + userId).then((data) => {
        console.log(data)
        this.followers = data.body.followers
        this.followed = true
      })
    },
    unfollow (userId) {
      this.$http.post('http://localhost:1323/unfollow/' + userId).then((data) => {
        console.log(data)
        this.followers = data.body.followers
        this.followed = false
      })
    },
    likePost (postId) {
      this.$http.post('http://localhost:1323/likePost', {id: this.post.id}).then((data) => {
        console.log(data)
        this.postLikes = data.body.likes
        this.postLiked = data.body.liked
      })
    },
    unlikePost (postId) {
      this.$http.post('http://localhost:1323/unlikePost', {id: this.post.id}).then((data) => {
        console.log(data)
        this.postLikes = data.body.likes
        this.postLiked = data.body.liked
      })
    },
    submitForm (formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          console.log(this.addCommentForm.postId)
          this.$http.post('http://localhost:1323/create-comment', {content: this.addCommentForm.content, postId: this.addCommentForm.postId}).then((data) => {
            this.post.comments.unshift(data.body)
          })
        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
    redirect (id) {
      router.push(this.userClaims.id === id ? '/profile' : '/user/' + id)
    },
    getFollowing () {
      this.$http.get('http://localhost:1323/user/' + this.post.author.id + '/following').then((data) => {
        this.followed = data.body
      })
      this.$http.get('http://localhost:1323/user/' + this.post.author.id + '/followers').then((data) => {
        this.followers = data.body
      })
    },
    resetForm (formName) {
      this.addPostForm.resetFields()
    }
  }

}
</script>

<style lang="scss" scoped>
.el-popover{
  button{
    width: 100%;
  }
}

.deleteButton {
  float: right;
  margin: 0;
  font-size: 2rem;
  color: #bbb;
}


.post{
  display: flex;
  background: white;
  &:not(:last-child){
    border-bottom: 1px dashed #e5e5ee;
  }
  padding: 1rem;
  align-items: flex-start;

  .content{
    margin-left: 1rem;
    width: 100%;
    h4{
      button{
        color: mix(#777, #409EFF, 20%);
        font-weight: bold;
        margin: 0;
        padding: 0;
      }
      margin: 0;
      color: mix(#777, #409EFF, 20%);
      small{
        color: #bbb;
        text-transform: uppercase;
      }
      font-size: 0.9rem;
    }

    p{
      border-left: 1px solid #e5e5ee;
      padding-left: 10px;
      margin: 5px 0;

      white-space: pre-wrap; 
      word-wrap: break-word;

    }

    .links{
      margin-bottom: 1rem;
      a{
        text-decoration: none;
        color: #409EFF;
        font-weight: bold;
        font-size: 0.8rem;
        margin-right: 10px;
      }
    }
  }

}
</style>

