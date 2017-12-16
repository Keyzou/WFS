<template>
  <div id="feed">
    <post v-bind:key="index" v-bind:post="post" v-for="(post, index) in posts"></post>
  </div>
</template>
<script>
import Post from '@/components/layout/Post'
export default {
  props: ['posts'],
  name: 'PostList',
  components: {
    Post
  }
}
</script>
