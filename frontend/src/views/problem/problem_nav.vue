<template>
<div class = "problemCard">
    <div class = "problemImg">
        <div id="app">
            <ul class="nav">
                <li :class="{'active' : idx === 0}">
                    <a @click="ShowContent(0)">问题详情</a>
                </li>
                <li :class="{'active' : idx === 1}">
                    <a @click="ShowContent(1)">提交代码</a>
                </li>
                <div class="line"></div>
            </ul>
        </div>
        <router-view />
    </div>
</div>
</template>
<script>
import { getCurrentURL, getCurrentQueryParams, getPathName } from '../../helper/getUrl';
export default {
    props: {  
        pid: {  
        type: String, // 或者 Number，取决于你的pid格式  
        required: true,  
        },
    },  
    data() {
    return {
        idx: 0,
        labels: ['问题详情', '提交代码'],
        currentURL: getCurrentURL(),
        queryParams: getCurrentQueryParams(),
        pathName: getPathName(),
    }
  },
  created() {
    if (Object.keys(this.queryParams).length === 0) {
      setTimeout(() => {
        this.goBack();
      }, 3000);
    }
  },
  methods: {
    goBack() {
      this.$router.go(-1);
    },
    ShowContent(index){
        this.idx = index;
        let path = '/problem';
        if(index === 1){
            path = '/submit';
        }
        if (this.$route.path !== path) {  
            this.$router.push(path+'?pid='+this.pid);  
        }  
    },
  }
}
</script>
<style scoped>
*{
    margin: 0;
    padding: 0;
    list-style: none;
    box-sizing: border-box;
    text-decoration: none;
}
.problemCard{
    width: 60%;
    margin: 0 auto;
    display: flex;
    /* border: 1px red solid; */
}
.problemImg{
    width: 100%;
}
.toolCard{
    width: 30%;
    display: block;
}
.problemCount{
    width: 100%;
    border: 1px green solid;
    justify-content: center;
}

.nav{
    width: 200px;
    height: 40px;
    /* border: 1px solid grey; */
    display: flex;
    flex-wrap: wrap;
    align-content: center;
    font-size: 14px;
    position: relative;
}

.nav li{
    width: 100%;
    flex: 1;
    text-align: center;
}

.nav a{
    width: 100%;
    color: #979797;
    transition: all .5s;
    text-decoration: none;
}

.active a{
    color: #000;
}

.line {
    position: absolute;
    bottom: 0;
    left: 20px;
    height: 3px;
    width: 60px;
    background-color: grey;
    transition: all .5s;
}

.nav li:nth-child(1).active ~ .line{
    transform: translateX(0);
    background-color: grey;
}

.nav li:nth-child(2).active ~ .line{
    transform: translateX(100px);
    background-color: grey;
}
</style>