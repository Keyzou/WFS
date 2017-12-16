<template>
  <div id="app">
    <el-container>
      <el-header style="height: 76px;">
        <header-component></header-component>
      </el-header>
      <el-container>
        <el-aside style="width: auto;">
          <aside-component></aside-component>
        </el-aside>
        <el-container id="content">
          <el-main>
            <router-view></router-view>
          </el-main>
          <el-footer style="text-align: center; color: #bbb">Copyright &copy; Dean CH&Eacute;RIF - 2017. Projet r&eacute;alis&eacute; pour le module SPT 1</el-footer>
        </el-container>
      </el-container>
    </el-container>
  </div>
</template>

<script>

import HeaderComponent from '@/components/layout/HeaderComponent'
import AsideComponent from '@/components/layout/AsideComponent'

export default {
  name: 'app',
  components: {
    HeaderComponent,
    AsideComponent
  }
}
</script>

<style lang="scss">
@import url('https://fonts.googleapis.com/css?family=Nunito:400,700');

.fa-icon {
  vertical-align: middle;
  height: 24px;
}

.el-card__header{
    background: #409EFF;
    color: mix(#f5faff, #409EFF, 90%);
    font-weight: bold;
  }

#app{
  max-width: 95rem;
  box-shadow: 0 0 10px rgba(33,37,41,.05);
  margin: auto;
}

body{
  background: #f1f3f5;
  padding: 0;
  margin: 0;
}

#content{
  padding: 0;
  margin: 0;
  background: #f8f9fa;
}

.el-header{
  padding: 0;
}

.el-menu-vertical{
  overflow: hidden;
  min-height: 970px;
}

.el-alert{
  padding: 1rem;
}

.container{
  width: 1024px;
  margin: auto;
}

</style>
