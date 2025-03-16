<template>
  <div class="box">
    <b-card class="problem">
      <div
        v-if="isBeforeStart(contest)"
        class="beforestart"
      >
        比赛开始前
      </div>
      <div v-else>
        <div class="Title">Contest - {{1000 + contest.ID}} {{ contest.title }} - 提交统计</div>
        <hr>
        <div class="time">
          <div class="contesttime">
            <div style="display: flex;">Start Time:&nbsp;<p style="color:firebrick">{{ formatDate(contest.start_at) }}&nbsp;&nbsp;</p>
            </div>
            <div style="display: flex;">End Time:&nbsp;<p style="color:firebrick"> {{ formatDate(contest.end_at) }}</p>
            </div>
          </div>
          <div class="curstatus">
            <div class="curtime">Current Time:&nbsp;<p style="color:firebrick">
                <CurrentTime />
              </p>&nbsp;&nbsp;</div>
            <div class="status">Status:&nbsp;
              <p
                style="color: firebrick;"
                v-if="isDuring(contest)"
              >Running</p>
              <p
                style="color: firebrick;"
                v-else
              >Ended</p>
              &nbsp;
              <p
                style="color: green;"
                v-if="contest.public==true"
              >Public</p>
              <p
                style="color: firebrick;"
                v-else
              >Private</p>
            </div>
          </div>
        </div>
        <div class="Button">
          <b-button
            variant="success"
            @click="Problem()"
          >题目列表</b-button>
          <b-button
            variant="primary"
            @click="Statistics()"
          >提交统计</b-button>
          <b-button
            variant="info"
            @click="Rank()"
          >竞赛排名</b-button>
          <b-button
            variant="warning"
            @click="Submit()"
          >提交记录</b-button>

        </div>
        <div class="problemList">
          <b-card>
            <div
              class="overflow-auto"
              v-if="this.items"
            >
              <b-table
                id="my-table"
                striped
                hover
                :items="items"
                :fields="fields"
              >
                <template #cell(题目编号)="data">
                  <div
                    class="col-num"
                    @click="navigateToProblem(data.item.ID)"
                  >Problem {{ String.fromCharCode(65 + data.index) }}</div>
                </template>

                <template #cell(标题)="data">
                  <div
                    class="col-title"
                    @click="navigateToProblem(data.item.ID)"
                  >{{ data.item.title }}</div>
                </template>

                <template #cell(解决)>
                  <div class="col-solve">0</div>
                </template>

                <template #cell(提交)>
                  <div class="col-submit">0</div>
                </template>

              </b-table>
            </div>
            <div v-else> 暂无题目数据 </div>
          </b-card>
        </div>
      </div>
    </b-card>
  </div>
</template>

<script>
import { mapActions, mapState, mapMutations } from 'vuex'
import CurrentTime from '@/components/CurrentTime.vue';
import { isDuring, isBeforeStart, formatDate, getTotalDuration } from '../../helper/getTime';
export default {
    components:{
        CurrentTime,
    },
    data() {
        return {
            contest:{},
            password:'',
            status:0,
            fields:['题目编号', '标题', '解决','提交'],
            problem:'',
            items:[],
            problemLettersMap: {},
        }
    },
    methods:{
    ...mapActions('contestModule',['GetContestProblem']),
    ...mapActions('contestModule',['getContestInfo']),
    ...mapMutations('contestModule', ['setContestInfo']),
    formatDate,
    getTotalDuration,
    isBeforeStart,
    isDuring,
    Problem(){
        const queryParams = window.location.search;
        if (this.$route.path !== "/contestdetail") {
            this.$router.push("/contestdetail" + queryParams);
        }
    },
    Submit(){
        const queryParams = window.location.search;
        if (this.$route.path !== "/contestsubmit") {
            this.$router.push("/contestsubmit" + queryParams);
        }
    },
    Statistics(){
        const queryParams = window.location.search;
        if (this.$route.path !== "/conteststatics") {
            this.$router.push("/conteststatics" + queryParams);
        }
    },
    Rank(){
        const queryParams = window.location.search;
        if (this.$route.path !== "/contestrank") {
            this.$router.push("/contestrank" + queryParams);
        }
    },
    navigateToProblem(pid){
        const queryParams =  window.location.search;
        if (isDuring(this.contestInfo)){
            window.open("/problem" + queryParams + '&pid=' + pid, '_blank');
        } else {
            window.open("/problem?pid=" + pid);
        }
    },
    async getInfo() {
        try {
            const response = await this.getContestInfo();
            this.contest = response.data.data.data;
            // this.contest.start_at = formatDate(this.contest.start_at);
            this.setContestInfo(response.data.data.data);
        } catch (error) {
            console.log(error);
            this.$bvToast.toast(error.response ? error.response.data.msg : '未知错误', {
                title: '数据验证错误',
                variant: 'danger',
                toaster: 'b-toaster-bottom-right',
                solid: true,
                appendToast: true,
            });
        }
    },
    generateProblemLettersMap() {
        this.problemLettersMap = {};
        this.items.forEach((item, index) => {
        this.problemLettersMap[String.fromCharCode(65 + index)] = item.ID;
        });
    }
  },
  async created() {
    try {
        const response = await this.GetContestProblem();
        this.items = response.data.data.data;
        // this.contest.start_at = formatDate(this.contest.start_at);
        this.generateProblemLettersMap();
        this.contest = this.contestInfo;

        if (Object.keys(this.contest).length === 0) {
            // contest 为空对象
            await this.getInfo();
        }
    } catch (error) {
        console.log(error);
        this.$bvToast.toast(error.response ? error.response.data.msg : '未知错误', {
            title: '数据验证错误',
            variant: 'danger',
            toaster: 'b-toaster-bottom-right',
            solid: true,
            appendToast: true,
        });
    }
},
  computed:{
    ...mapState('userModule', {
      userInfo: (state) => state.userInfo,
    }),
    ...mapState('contestModule', {
      contestInfo: (state) => state.contestInfo,
    }),
  },
}    
</script>
<style scoped>
.box {
  width: 1200px;
  margin: 20px auto;
}
.problem {
  width: 1200px;
}
.Title {
  display: flex;
  justify-content: center;
  font-size: 22px;
}
.downtime {
  widows: 10px;
}
.curtime {
  display: flex;
}
.contesttime {
  display: flex;
  height: 20px;
  justify-content: center;
}
.status {
  display: flex;
}
.curstatus {
  display: flex;
  height: 20px;
  justify-content: center;
}
.beforestart {
  color: green;
  font-size: 30px;
  display: flex;
  justify-content: center;
}
.problemList {
  margin-top: 20px;
}
.Button {
  margin: 10px;
  display: flex;
  justify-content: center;
}
.col-num {
  width: 80px;
  display: flex;
  color: #006eff;
  cursor: pointer;
}
.col-title {
  width: 200px;
  text-wrap: nowrap;
  overflow: hidden;
  color: #006eff;
  cursor: pointer;
}
.col-solve .col-submit {
  width: 50px;
  text-wrap: nowrap;
  display: flex;
  justify-content: center;
}
</style>
