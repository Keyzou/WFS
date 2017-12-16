<template>
<div class="comment" v-if="!deleted">
  <el-popover
  ref="popover4"
  placement="top"
  trigger="hover"
  v-on:show="getFollowing()">
    <el-button v-if="!followed && !(this.userClaims.id === this.comment.author.id)" type="primary" @click="follow(comment.author.id)"><font-awesome-icon icon="user-plus" /> Follow <strong>({{followers}})</strong></el-button>
    <el-button v-else-if="this.userClaims.id === this.comment.author.id" type="primary" @click="redirect(this.comment.author.id)">
      Profile
    </el-button>
    <el-button v-else type="danger" @click="unfollow(comment.author.id)"><font-awesome-icon icon="user-times" /> Unfollow <strong>({{followers}})</strong></el-button>
  </el-popover>
  <v-gravatar v-bind:email="comment.author.email" />
  <div class="content">
    <el-button v-if="this.userClaims.id === this.comment.author.id" type="text" @click="deleteComment()" class="deleteButton">&times;</el-button>
    <h4><el-button type="text" @click="redirect(comment.author.id)" v-popover:popover4>{{comment.author.username}}</el-button> <small>| {{comment.created}}</small></h4>
    <p v-html="comment.content"></p>
    <div class="links">
      <el-button type="text" v-if="!commentLiked" @click="likeComment(comment.id)">{{commentLikes}} Likes <font-awesome-icon icon="thumbs-up" /></el-button>
      <el-button type="text" v-else @click="unlikeComment(comment.id)">{{commentLikes}} Likes <font-awesome-icon icon="thumbs-up" flip="vertical"/></el-button>
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
      commentLiked: this.comment.liked,
      followers: this.comment.author.followers,
      followed: this.comment.isFollowing,
      deleted: false,
      commentLikes: this.comment.likes,
      userClaims: auth.getUserClaims(localStorage.getItem('token'))
    }
  },
  name: 'Comment',
  props: ['comment'],
  methods: {
    deleteComment () {
      this.$http.delete('http://localhost:1323/deleteComment/' + this.comment.id, {id: this.comment.id}).then((data) => {
        this.deleted = true
      })
    },
    follow (userId) {
      this.$http.post('http://localhost:1323/follow/' + userId).then((data) => {
        this.followers = data.body.followers
        this.followed = true
      })
    },
    unfollow (userId) {
      this.$http.post('http://localhost:1323/unfollow/' + userId).then((data) => {
        this.followers = data.body.followers
        this.followed = false
      })
    },
    likeComment (commentId) {
      this.$http.post('http://localhost:1323/likeComment', {id: this.comment.id}).then((data) => {
        this.commentLikes = data.body.likes
        this.commentLiked = data.body.liked
      })
    },
    unlikeComment (commentId) {
      this.$http.post('http://localhost:1323/unlikeComment', {id: this.comment.id}).then((data) => {
        this.commentLikes = data.body.likes
        this.commentLiked = data.body.liked
      })
    },
    getFollowing () {
      this.$http.get('http://localhost:1323/user/' + this.comment.author.id + '/following').then((data) => {
        this.followed = data.body
      })
      this.$http.get('http://localhost:1323/user/' + this.comment.author.id + '/followers').then((data) => {
        console.log(data.body)
        this.followers = data.body
      })
    },
    redirect (id) {
      router.push(this.userClaims.id === id ? '/profile' : '/user/' + id)
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

.comment{
  display: flex;
  &:nth-child(2n){
    background: #fafafa;
  }
  background: white;
  border-left: 1px solid #e5e5ee;
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

