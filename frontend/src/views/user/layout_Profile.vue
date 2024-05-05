<template>
    <div class = "container">
        <b-card class = "first-card">
            <div class = "second-container">
                <div class="userInfo">
                    <div>
                        <img :src="require('../../../../backend/static/icon/' + userIcon)"  class="icon" alt="用户信息获取失败">
                    </div>
                    <div>
                        <h3> Hello! <span v-if="userInfo">{{ userInfo.user_name }}</span></h3>
                        <h5>{{ userInfo.user_nickname }}</h5>
                        <h5>{{ userInfo.user_email }}</h5>
                    </div>
                    <b-button @click="$router.push({name:'editProfile'})" variant="outline-primary">修改个人信息</b-button>
                </div>
                <div class="second-card">Self info</div>
            </div>
        </b-card>
    </div>
</template>
<script>
import { mapState } from 'vuex';
export default {
    data () {
        return {
            userIcon:"init.jpg" 
        }
    },
    computed: {
        ...mapState({
            userInfo: (state) => state.userModule.userInfo,
        }),
        getUserIcon() {
            // this.userIcon = userInfo.user_icon;
            return require(`./${this.userIcon}`);
        }
    },
    created() {
    if (this.userInfo.user_icon) {
      this.userIcon = this.userInfo.user_icon;
    }
  },
};
</script>
<style scoped>
 .icon{
    width: 200px;
    height: 200px;
    border: 1px red solid;
}

.container
{
    width: 50%;
    margin: 10px auto;
    
}

.userInfo {
    width: 40%;
}

.second-card{
    width: 60%;
    border: 1px red solid;
}
.second-container{
    display: flex;
}

</style>
