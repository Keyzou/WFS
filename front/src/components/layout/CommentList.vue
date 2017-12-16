<template>
  <div id="comments">
    <comment v-bind:key="index" v-bind:comment="comment" v-for="(comment, index) in comments"></comment>
  </div>
</template>
<script>
import Comment from '@/components/layout/Comment'
export default {
  props: ['comments'],
  name: 'CommentList',
  components: {
    Comment
  }
}
</script>
<style lang="scss" scoped>
</style>

