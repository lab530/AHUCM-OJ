<template>
    <div>
        <div id="app">
        <div class="nav">
            <ul>
            <li :class="{'active' : idx === 0}">
                <a @click="navigateTo('/problemset', 0)">问题列表</a>
            </li>
            <li :class="{'active' : idx === 1}">
                <a @click="navigateTo('/problem/status', 1)">提交状态</a>
            </li>
            <li :class="{'active' : idx === 2}">
                <a @click="navigateTo('/problem/rank', 2)">当前排名</a>
            </li>
            <div class="nav-box"></div>
            </ul>
        </div>
        </div>
        <router-view/>
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
        } else {
            const currentQuery = this.$route.query;
            const hasCategoryParam = 'category' in currentQuery;
            if(hasCategoryParam) this.$router.push(path);
        }
    },
},
};
</script>

<style scoped>
* {
    margin: 0;
    padding: 0;
    list-style: none;
    box-sizing: border-box;
    text-decoration: none;
}

.body {
    display: flex;
    justify-content: center;
    background-color: #8da1f8;
}

.nav {
    margin: 10px auto;
    width: 70%;
    position: relative;
    top: 0;
    box-shadow: 0 5px 20px rgba(0, 0, 0, .2);
    border-radius: 10px;
    background-color: #fff;
}

.nav ul {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    right: 5%;
}

.nav ul .nav-box {
    position: absolute;
    bottom: 0;
    left: 0;
    width: calc((100%/3)*1);
    height: 5px;
    border-radius: 2px;
    background-color: rgb(112, 184, 184);
    transition: transform.5s;
}

.nav ul li {
    width: 100%;
    text-align: center;
    cursor: pointer;
}

.nav ul li a {
    color: rgb(70, 100, 180);
    font: 100 15px '';
    /* display: block; */
    /* width: 100%; */
    height: 100%;
    line-height: 50px;
    text-decoration: none;
    color: gray;
}

.nav li:nth-child(1).active~.nav-box {
    transform: translateX(calc((100%/3)*0));
    background-color: rgb(112, 184, 184);
}

.nav li:nth-child(2).active~.nav-box {
    transform: translateX(calc((100%/3)*3));
    background-color: rgb(112, 184, 184);
}

.nav li:nth-child(3).active~.nav-box {
    transform: translateX(calc((100%/3)*6));
    background-color: rgb(112, 184, 184);
}
</style>