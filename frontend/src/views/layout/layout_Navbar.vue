<template>
    <div>
    <b-navbar toggleable="lg" type="dark" variant="info">
        <b-container>
            <b-navbar-brand @click="$router.push({name:'Home'})" style="cursor: pointer;">AHUCM Online Judge</b-navbar-brand>

            <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

            <b-navbar-nav>
              <b-nav-item href="#">新闻/公告</b-nav-item>
              <b-nav-item @click="navigateTo('/problemset')" :disabled="isChildRoute('/problem')">问题</b-nav-item>
              <b-nav-item href="#">竞赛</b-nav-item>
              <b-nav-item @click="navigateTo('/blog')" :disabled="isChildRoute('/blog')">博客</b-nav-item>
              <b-nav-item href="#">常见问题</b-nav-item>
            </b-navbar-nav>
            <!-- Right aligned nav items -->
            <b-navbar-nav class="ml-auto">

            <b-nav-item-dropdown right v-if="userInfo">
            <template #button-content>
                <b-avatar :src="require('../../../../backend/static/icon/' + computedUserIcon)"> </b-avatar>
                &nbsp;
                <em>{{ userInfo.user_name }}</em>
            </template>
            <b-dropdown-item @click="navigateTo('/profile')">个人中心</b-dropdown-item>
            <b-dropdown-item @click="navigateTo('/admin')">后台管理</b-dropdown-item>
            <b-dropdown-item @click="logout">注销</b-dropdown-item>
            
            </b-nav-item-dropdown>
            <div v-if="!userInfo">
              <b-nav-item
                class = "inline"
                v-if="$route.name != 'login'"
                @click="$router.replace({name:'login'})">登录</b-nav-item>
              <b-nav-item
                class = "inline"
                v-if="$route.name != 'register'"
                @click="$router.replace({name:'register'})">注册</b-nav-item>
            </div>
            </b-navbar-nav>
        </b-container>
    </b-navbar>
    </div>
</template>
<script>
import { mapState, mapActions } from 'vuex';
import storageService from '../../service/storageService';
export default {
  data() {
    return {
      userIcon:"init.jpg",
      validation: null,
      token:'',
    };
  },
  computed: {
    ...mapState('userModule', {
      userInfo: (state) => state.userInfo,
    }),
    computedUserIcon() {
      if (this.userInfo.user_icon != "") {
        return this.userInfo.user_icon;
      }
      return 'init.jpg'; // 默认值
    },
  },
  methods: {
    ...mapActions('userModule', ['logout']),
    ...mapActions('userModule', ['getInfo']),
    navigateTo(path) {
      if (!this.isChildRoute(path)) {
        if(this.$route.path != path) this.$router.push(path);
      }
    },
    isChildRoute(path) {
      // 判断当前路由是否是子页面的路径
      return this.$route.path.startsWith(path + '/');
    },
    manage(){
      console.log("yes")
    }
  },
  mounted(){
    console.log(1)
    if(this.token && this.token.length > 0){
      this.getInfo("hello").catch((error) => {
        console.log(3)
        this.$bvToast.toast(error.response?.data.msg, {
            title: '信息错误',
            variant: 'danger',
            solid: true,
          });
          console.log(4)
      });
    }
  },
  created(){
    this.token = storageService.get(storageService.USER_TOKEN);
  }
};
</script>
<style lang="scss" scoped>
.inline{
  display: inline-flex;
  padding-right: 10px;
}
</style>
