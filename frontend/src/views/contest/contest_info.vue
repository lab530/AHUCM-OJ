<template>
  <div>
    <b-card
      class="box"
      v-if="contest.title"
    >
      <div class="title">{{contest.title}}</div>
      <b-alert
        show
        class="time"
      >
        <CountdownTimer
          title="距离比赛开始"
          :targetDate="contest.start_at"
        />
      </b-alert>
      <div class="desc">比赛描述：</div>
      <mavon-editor
        v-model="contest.description"
        class="content-show"
        :subfield="false"
        defaultOpen="preview"
        :editable="false"
        :toolbarsFlag="false"
        :boxShadow="false"
      />
      <div class="button">
        <div>
          <b-button
            class="mt-3"
            variant="success"
            @click="Detail()"
          >竞赛详情</b-button>
          <b-button
            class="mt-3"
            variant="info"
            @click="RankList()"
          >竞赛排名</b-button>
        </div>
        <div
          v-if='contest.password == "True"'
          class="password"
        >
          <b-input-group class="mt-3">
            <b-form-input
              v-model="password"
              placeholder="请输入竞赛密码"
            ></b-form-input>
            <b-input-group-append>
              <b-button
                variant="outline-success"
                @click="SubmitPassword()"
              >提交</b-button>
            </b-input-group-append>
          </b-input-group>
        </div>
      </div>
    </b-card>
    <b-card v-else>
      Loading
    </b-card>
  </div>
</template>
<script>
import { mapActions, mapState, mapMutations } from 'vuex'
import CountdownTimer from '@/components/CountdownTimer.vue';
import { formatDate} from '../../helper/getTime';
export default {
    components:{
        CountdownTimer,
    },
    data() {
        return {
            contest:{},
            password:'',
        }
  },
  methods:{
    ...mapActions('contestModule',['getContestInfo']),
    ...mapActions('contestModule',['SubmitContestPassword']),
    ...mapActions('contestModule',['VerityContest']),
    ...mapMutations('contestModule', ['setContestInfo']),
    formatDate,
    GetContestInfo(){
        this.getContestInfo().then((response) => {
            this.contest = response.data.data.data;
            // this.contest.start_at = formatDate(this.contest.start_at);
            console.log(this.contest)
            this.changeDateForm();
            this.setContestInfo(response.data.data.data)
        }).catch((error) => {
        this.$bvToast.toast(error.response.data.msg, {
            title: '数据验证错误',
            variant: 'danger',
            toaster:'b-toaster-bottom-right',
            solid: true,
            appendToast:true,
            });
        });
    },
    Detail(){
        if(this.contest.public && this.contest.password == "False"){
            const queryParams = window.location.search;
            this.$router.push("/contestdetail" + queryParams);
        } else {
            if(this.userInfo){
                let UserId = this.userInfo.id
                this.VerityContest( UserId ).then((response) => {
                    this.$bvToast.toast(response.data.msg, {
                    title: '数据验证成功',
                    variant: 'success',
                    toaster:'b-toaster-bottom-right',
                    solid: true,
                    appendToast:true,
                    });
                    const queryParams = window.location.search;
                    setTimeout(() => {
                        this.$router.push("/contestdetail" + queryParams);
                    }, 500)
                }).catch((error) => {
                this.$bvToast.toast(error.response.data.msg, {
                    title: '数据验证错误',
                    variant: 'danger',
                    toaster:'b-toaster-bottom-right',
                    solid: true,
                    appendToast:true,
                    });
                });
            } else {
                this.$bvToast.toast("游客仅提供输入密码的方式查看", {
                title: '错误',
                variant: 'danger',
                toaster:'b-toaster-bottom-right',
                solid: true,
                appendToast:true,
                });
            }
        }
    },
    RankList(){
        const queryParams = window.location.search;
        this.$router.push("/contestrank" + queryParams);
    },
    SubmitPassword(){
        let UserId = -1, ContestPassword = this.password;
        if(this.userInfo){
            UserId = this.userInfo.id
        }
        this.SubmitContestPassword( { UserId, ContestPassword }).then((response) => {
            this.$bvToast.toast(response.data.msg, {
            title: '数据验证成功',
            variant: 'success',
            toaster:'b-toaster-bottom-right',
            solid: true,
            appendToast:true,
            });
            const queryParams = window.location.search;
            setTimeout(() => {
                    this.$router.push("/contestdetail" + queryParams);
                }, 500)
        }).catch((error) => {
        this.$bvToast.toast(error.response.data.msg, {
            title: '数据验证错误',
            variant: 'danger',
            toaster:'b-toaster-bottom-right',
            solid: true,
            appendToast:true,
            });
        });
    },
    changeDateForm(){
        const dateString = this.contest.start_at;
        // 进行后续的日期处理和分割操作
        const parts = dateString.split('T');
        const datePart = parts[0];
        const timePart = parts[1].split('+')[0];
        const [year, month, day] = datePart.split('-');
        const [hours, minutes, seconds] = timePart.split(':');
        const targetDate = new Date(year, month - 1, day, hours, minutes, seconds);
        this.contest.start_at = targetDate.toString();
    }
  },
  created(){
    this.GetContestInfo();
  },
  computed:{
    ...mapState('userModule', {
      userInfo: (state) => state.userInfo,
    }),
  },
}
</script>
<style scoped>
.box {
  margin: 50px auto;
  max-width: 1100px;
}
.title {
  font-size: 30px;
  display: flex;
  justify-content: center;
}
.desc {
  margin-bottom: 10px;
  font-size: 20px;
}
/deep/ .v-note-wrapper .v-note-panel .v-note-show .v-show-content,
/deep/ .v-note-wrapper .v-note-panel .v-note-show .v-show-content-html {
  height: 100%;
  padding: 0;
  overflow-y: auto;
  box-sizing: border-box;
  overflow-x: hidden;
}
.password {
  width: 250px;
}
.button {
  display: flex;
  justify-content: space-between;
}
.time {
  display: flex;
  justify-content: center;
}
</style>
