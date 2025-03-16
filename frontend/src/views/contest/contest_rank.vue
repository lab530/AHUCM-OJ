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
        <div class="Title">Contest - {{1000 + contest.ID}} {{ contest.title }} - 竞赛排名</div>
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
      </div>
      <div class="bgc">
        <b-table
          id="ranking-table"
          :items="rankingData"
          :fields="fields"
          striped
          responsive
          hover
          class="text-center"
          head-variant="dark"
          body-class="cell-content"
          :style="{
                '--bs-table-height': '60px'
            }"
        >
          <template #cell(rank)="data">
            {{ data.index + 1 }}
          </template>
          <template #cell(item)="data">
            <div class="cell-content">
              <div class="main-text">
                <span>{{ data.field.key }}</span>
              </div>
              <div class="sub-text">
                <span v-if="data.value == '+'">+1</span>
                <span v-else>-</span>
              </div>
            </div>
          </template>
        </b-table>
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
            letterProblemsMap: {},
            problemList:[],
            fields: [
            { key: 'rank', label: 'Rank' },
            { key: 'name', label: 'Name' },
            { key:'solve', label:'Solve'},
            { key: 'penalty', label: 'Penalty' },          
            ],
            rankingData: [],
            data:[],
        };
    },
    methods:{
    ...mapActions('contestModule',['GetContestProblem']),
    ...mapActions('contestModule',['getContestInfo']),
    ...mapMutations('contestModule', ['setContestInfo']),
    ...mapActions('contestModule',['GetRankInfo']),
    formatDate,
    getTotalDuration,
    isBeforeStart,
    isDuring,
    processedData(data){
        let UserMap = new Map();
        let index = 0;
        data.forEach(item => {
            let user_name = item.name;
            if(!UserMap.has(user_name)){
                let userData = {  
                    name: user_name,  
                    solve: 0,  
                    penalty: 0,  
                };
                this.rankingData.push(userData);
                UserMap.set(user_name, index);
                index ++;
            };
            console.log(UserMap)
            let  ProblemID = this.letterProblemsMap[item.problem_id];
            let  pass = item.accepted;
            let idx = UserMap.get(user_name);
            console.log(idx)
            if(pass == true){
                console.log(item)
                this.rankingData[idx].solve++;
                console.log(1)
                this.rankingData[idx].penalty += item.penalty;
                console.log(2)
                let s = this.formatTime(item.penalty);
                if(item.penalty_count != 0) s += `(+${item.penalty_count})`;
                this.rankingData[idx][ProblemID] = s;
                console.log(3)
            } else {
                console.log(item)
                this.rankingData[idx][ProblemID] = `(-${item.penalty_count})`;
            }
        })
        this.rankingData = this.rankingDataSort(this.rankingData);
    },
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
    async getContestRank() {
      try {
          const response = await this.GetRankInfo();
          console.log(response.data.data.data);
          this.data = response.data.data.data;
          this.processedData(this.data);
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
    formatTime(seconds) {  
        // 将秒转换为分钟和秒  
        const minutes = Math.floor(seconds / 60);  
        const remainingSeconds = seconds % 60;  
        const formattedMinutes = minutes.toString().padStart(2, '0'); // 两位数格式  
        const formattedSeconds = remainingSeconds.toString().padStart(2, '0'); // 两位数格式  
        return `${formattedMinutes}:${formattedSeconds}`;  
    },  
    generateProblemLettersMap() {
        this.letterProblemsMap = {};
        this.problemList.forEach((item, index) => {
        this.letterProblemsMap[item.ID] = String.fromCharCode(65 + index);
        this.fields.push({key:String.fromCharCode(65 + index), label:String.fromCharCode(65 + index)})
        });
    },
    rankingDataSort(data) {
      return data.sort((a, b) => {
        // 先根据题解数进行比较
        if (b.solve - a.solve !== 0) {
          return b.solve - a.solve;
        }
        // 如果题解数相同,则根据总罚时进行比较
        if (a.penalty - b.penalty !== 0) {
          return a.penalty - b.penalty;
        }
        // 如果题解数和总罚时都相同,则保持原顺序
        return 0;
      }).map(item => {  
        // 假设 penalty 是秒数  
        item.penalty = this.formatTime(item.penalty);  
        // 如果还有其他时间字段需要格式化，可以在这里继续添加  
        // ...  
        return item;  
    });
    },
  },
  async created() {
    try {
        const response = await this.GetContestProblem();
        console.log(response);
        this.problemList = response.data.data.data;
        console.log(this.problemList);
        this.generateProblemLettersMap();

        this.contest = this.contestInfo;

        if (Object.keys(this.contest).length === 0) {
            // contest 为空对象
            await this.getInfo();
        }
        await this.getContestRank();
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
.cell-content {
  display: flex;
  flex-direction: column;
  justify-content: center;
  height: 100%;
}

.main-text {
  font-size: 1rem;
}

.sub-text {
  font-size: 0.8rem;
  line-height: 1;
}

.bgc {
  background-image: url("../../../public/bg.png");
  background-repeat: repeat;
  background-size: 150px 150px;
}
</style>
