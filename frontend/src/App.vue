<template>
  <div id="app">
    <navbar v-show="!isAdminRoute"/>
    <router-view/>
    <!-- <bottom-time class="bottom-component"/> -->
  </div>
</template>

<script>
import Navbar from './views/layout/layout_Navbar.vue';
import BottomTime from './views/bottom/bottom_time.vue';

export default {
  components: {
    Navbar,
    BottomTime
  },
  data() {
    return {};
  },
  computed: {
    isAdminRoute() {
      return this.$route.path.startsWith('/admin'); // 检查当前路由路径是否以 "admin/" 开头
    }
  }
};
</script>

<style scoped>

.app {
  position: relative;
  display: flex;  
  justify-content: center; /* 水平居中 */  
}
.bottom-component {
  /* position: fixed; 移除固定定位 */  
width: 100%;
  /* 可以添加 margin 或 padding 来调整它与页面内容的间距 */  
  margin-top: 20px; /* 例如，添加上边距 */  
}
</style>