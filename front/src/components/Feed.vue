<template>
<div id="index">
  <h1><font-awesome-icon icon="comment-alt" size="lg" /> My Personal Space <font-awesome-icon icon="comment-alt" size="lg" flip="horizontal" /></h1>
  <el-form :model="addPostForm" :rules="rules" ref="addPostForm">
    <el-form-item label="Add a post" prop="content">
      <el-input type="textarea" placeholder="Your post!" v-model="addPostForm.content"></el-input>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="submitForm('addPostForm')"><font-awesome-icon icon="paper-plane" /> Post</el-button>
    </el-form-item>
  </el-form>
  <post-list v-bind:posts="posts"></post-list>
</div>
</template>

<script>
import FontAwesomeIcon from '@fortawesome/vue-fontawesome'
import PostList from '@/components/layout/PostList'

export default {
  components: {
    FontAwesomeIcon,
    PostList
  },
  name: 'feed',
  data () {
    return {
      posts: [],
      addPostForm: {
        content: ''
      },
      rules: {
        content: [
          { required: true, message: 'Content cannot be empty!' }
        ]
      }
    }
  },
  mounted: function () {
    this.$http.get('http://localhost:1323/feed').then((data) => {
      console.log(data.body)
      this.posts = data.body
    })
  },
  methods: {
    submitForm (formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.$http.post('http://localhost:1323/create-post', {content: this.addPostForm.content}).then((data) => {
            this.posts = this.updatePostList()
          })
        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
    resetForm (formName) {
      this.addPostForm.resetFields()
    },
    updatePostList () {
      console.log('triggered')
      this.$http.get('http://localhost:1323/feed').then((data) => {
        console.log(data.body)
        this.posts = data.body
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss" scoped>
h1 {
  font-size: 2rem;
  margin: 0;
  text-align: center;
  color: #333;
  svg{
    color:#409EFF;
  }
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
}

#feed{
  border: 1px solid #e5e5ee;
}
</style>
