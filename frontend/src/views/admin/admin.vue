<template>

<div>
<div class = "header">AHUCM Online Judge Backstage management </div>
  <div class = "manage">
      <div class="nav">
        <div class = "list">
          <div @click="navigateTo('/', 0)"> OJ 首页 </div>
          <div :class="{'active' : idx === 1}" @click="navigateTo('/admin', 1)">管理中心</div>
          <div :class="{'active' : idx === 2}" @click="navigateTo('/admin/new', 2)"><p>新闻管理</p></div>
          <div :class="{'active' : idx === 3}" @click="navigateTo('/admin/user', 3)"><p>用户管理</p></div>
          <div :class="{'active' : idx === 4}" @click="navigateTo('/admin/problem', 4)"><p>问题管理</p></div>
          <div :class="{'active' : idx === 5}" @click="navigateTo('/admin/contest', 5)"><p>竞赛管理</p></div>
          <div :class="{'active' : idx === 6}" @click="navigateTo('/admin/system', 6)"><p>系统管理</p></div>
        </div>
      </div>
      <router-view class="router"/>
  </div>
</div>


</template>
<script>

export default {
data() {
  return {
    idx: 0
  }
},

methods: {
    navigateTo(path, index) {
        if (this.$route.path !== path) {
            this.$router.push(path);
            this.idx = index;
        }
    },
  },
};
</script>
<style scoped>
.header {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  width: 100%;
  height: 50px;
  background-color: #2f4796;
  color: white;
  font-family: "Arial", "Microsoft YaHei", "黑体", "宋体", sans-serif;
  display: flex;
  justify-content: center;
  font-size: 30px;
  line-height: 50px;
  z-index: 9999;
}
.manage {
  margin-top: 50px;
  height: 100%;
  display: flex;
}
.list {
  margin-top: 50px;
  cursor: pointer;
}
.nav {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  display: flex;
  flex-direction: column; /* 将主轴方向设置为竖直方向 */
  width: 200px;
  height: 100vh;
  background-color: #2f4796;
}

.nav .list div {
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 20px;
  width: 200px;
  height: 50px;
  color: white;
  font-family: "Arial", "Microsoft YaHei", "黑体", "宋体", sans-serif;
  transition: 0.3s;
}

.nav .list p {
  margin: 0; /* 去除 <p> 元素的默认外边距 */
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%; /* 填充父元素的高度 */
}

.nav .list div:hover {
  background-color: #708de9;
}
.active {
    background-color: #708de9;
}
.router {
  margin-left: 200px; /* 导航栏的宽度 */
  width: calc(100% - 200px); /* 子组件宽度减去导航栏宽度 */
}
</style>