<template>
    <div class="box">
        <b-card class="problem">
            <div v-if="isBeforeStart(contest)" class="beforestart">
                比赛开始前
            </div>
            <div v-else>
                <div class="Title">Contest - {{1000 + contest.ID}} {{ contest.title }} - 提交记录</div>
                <hr>
                <div class="time">
                    <div class="contesttime">
                        <div style="display: flex;">Start Time:&nbsp;<p style="color:firebrick">{{ formatDate(contest.start_at) }}&nbsp;&nbsp;</p></div>
                        <div style="display: flex;">End Time:&nbsp;<p style="color:firebrick"> {{ formatDate(contest.end_at) }}</p></div>
                    </div>
                    <div class="curstatus">
                        <div class="curtime">Current Time:&nbsp;<p style="color:firebrick"> <CurrentTime/></p>&nbsp;&nbsp;</div> 
                        <div class="status">Status:&nbsp;
                            <p style="color: firebrick;" v-if="isDuring(contest)">Running</p>
                            <p style="color: firebrick;" v-else>Ended</p>
                            &nbsp;
                            <p style="color: green;" v-if="contest.public==true">Public</p>
                            <p style="color: firebrick;" v-else>Private</p>
                            </div>
                    </div> 
                </div>
                <div class="Button">
                    <b-button variant="success" @click="Problem()">题目列表</b-button>
                    <b-button variant="primary" @click="Statistics()">提交统计</b-button>
                    <b-button variant="info" @click="Rank()">竞赛排名</b-button>
                    <b-button variant="warning" @click="Submit()">提交记录</b-button>
                </div>
                <div class="problemList">
                    <b-card>
                        <div class="overflow-auto">
                            <b-table
                            id="my-table"
                            striped hover 
                            :items="items"
                            :fields="fields"
                            :per-page="perPage"
                            :current-page="currentPage"
                            >
                            <template #cell(提交编号)="data">
                            <div class="col-snum">{{ 1000 + data.item.ID }}</div>
                            </template>  
                            
                            <template #cell(用户)="data">
                            <div class="col-user" @click="navigateToUser(data.item.user.ID)">{{ data.item.user.user_name }}</div>
                            </template>

                            <template #cell(题目编号)="data">
                            <div class="col-pnum" @click="navigateToProblem(data.item.problem_id)">{{ letterProblemsMap[data.item.problem_id] }}</div>
                            </template>


                            <template #cell(结果)="data">
                            <div class="col-status" >
                                <b-button v-if="data.item.status == 0" variant="warning" size="sm">
                                    <b-icon icon="arrow-counterclockwise" animation="spin-reverse"></b-icon> Pending
                                </b-button>
                                <b-button v-else-if="data.item.status == 1" variant="warning" size="sm">
                                    <b-icon icon="arrow-counterclockwise" animation="spin-reverse"></b-icon> PendingRejudge
                                </b-button>
                                <b-button v-else-if="data.item.status == 2" variant="warning" size="sm">
                                    <b-icon icon="arrow-counterclockwise" animation="spin-reverse"></b-icon> Compiling
                                </b-button>
                                <b-button v-else-if="data.item.status == 3" variant="warning" size="sm">
                                    <b-icon icon="arrow-counterclockwise" animation="spin-reverse"></b-icon> Running
                                </b-button>
                                <b-button v-else-if="data.item.status == 4" variant="success" size="sm">
                                    <b-icon icon="check2"></b-icon> Accepted
                                </b-button> 
                                <b-button v-else-if="data.item.status == 5" variant="info" size="sm">
                                    <b-icon icon="calendar3-week"></b-icon> PresentationError
                                </b-button> 
                                <b-button v-else-if="data.item.status == 6" variant="danger" size="sm">
                                    <b-icon icon="x-circle"></b-icon> WrongAnswer
                                </b-button> 
                                <b-button v-else-if="data.item.status == 7" variant="dark" size="sm">
                                    <b-icon icon="hourglass-split"></b-icon> TimeLimitExceeded
                                </b-button> 
                                <b-button v-else-if="data.item.status == 8" variant="dark" size="sm">
                                    <b-icon icon="file-earmark"></b-icon> MemoryLimitExceeded
                                </b-button> 
                                <b-button v-else-if="data.item.status == 9" variant="info" size="sm">
                                    <b-icon icon="collection"></b-icon> OutputLimitExceeded
                                </b-button> 
                                <b-button v-else-if="data.item.status == 10" variant="dark" size="sm">
                                    <b-icon icon="exclamation-triangle"></b-icon> RuntimeError
                                </b-button> 
                                <b-button v-else-if="data.item.status == 11" variant="info" size="sm">
                                    <b-icon icon="bug"></b-icon> CompileError
                                </b-button> 
                                <b-button v-else variant="outline-primary" size="sm">
                                    <b-icon icon="question"></b-icon> UnknownError
                                </b-button> 
                            </div>
                            </template>

                            <template #cell(内存)="data">
                            <div class="col-m" v-if="data.item.status > 3" >
                                <div v-if="(!isUserInfoEmpty && userInfo.id == data.item.user.ID) || userInfo.permission_id > 0" >{{ data.item.memo_limit }}</div>
                                <div v-else>---</div>
                            </div>
                            </template>

                            <template #cell(耗时)="data">
                            <div class="col-t" v-if="data.item.status > 3" >
                            <div v-if="(!isUserInfoEmpty && userInfo.id == data.item.user.ID) || userInfo.permission_id > 0" >{{ data.item.time_limit }}ms</div>
                            <div v-else>---</div>
                            </div>
                            </template>

                            <template #cell(语言)="data">
                            <div class="col-lang" @click="navigateToProblem(data.item.ID)">{{ data.item.lang }}</div>
                            </template>

                            <template #cell(提交时间)="data">
                            <div class="col-time">    
                                {{
                                    formatDate(new Date(data.item.submit_time))
                                }}
                            </div>
                            </template>

                            </b-table>
                            <b-pagination
                            align="center"
                            v-model="currentPage"
                            :total-rows="rows"
                            :per-page="perPage"
                            first-text="First"
                            prev-text="Prev"
                            next-text="Next"
                            last-text="Last"
                            ></b-pagination>
                        </div>
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
            fields:['提交编号', '用户', '题目编号','结果', '内存', '耗时', '语言', '提交时间' ],
            problemList:[],
            items:[],
            problemLettersMap:{},
            letterProblemsMap:{},
            perPage: 50,
            currentPage: 1,
            submit:{},
        }
    },
    methods:{
    ...mapActions('contestModule',['GetContestProblem']),
    ...mapActions('contestModule',['getContestInfo']),
    ...mapMutations('contestModule', ['setContestInfo']),
    ...mapActions('contestModule',['GetContestSubmit']),
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
    getInfo() {
        this.getContestInfo().then((response) => {
            this.contest = response.data.data.data;
            // this.contest.start_at = formatDate(this.contest.start_at);
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
    generateProblemLettersMap() {
        this.problemLettersMap = {};
        this.problemList.forEach((item, index) => {
        this.problemLettersMap[String.fromCharCode(65 + index)] = item.ID;
        });
        this.letterProblemsMap = {};
        this.problemList.forEach((item, index) => {
        this.letterProblemsMap[item.ID] = String.fromCharCode(65 + index);
        });
    },
  },
  created(){
    this.GetContestProblem().then((response) => {
        this.problemList = response.data.data.data;
        // this.contest.start_at = formatDate(this.contest.start_at);
        this.generateProblemLettersMap();
        
        this.contest = this.contestInfo;
        console.log(this.problemLettersMap['A']); // 输出 ID
        console.log(this.problemLettersMap['B']); // 输出 ID
        if (Object.keys(this.contest).length === 0) {
            // contest 为空对象
            this.getInfo();
        }
    }).catch((error) => {
    this.$bvToast.toast(error.response.data.msg, {
        title: '数据验证错误',
        variant: 'danger',
        toaster:'b-toaster-bottom-right',
        solid: true,
        appendToast:true,
        });
    });
    this.GetContestSubmit().then((response) => {
        this.items = response.data.data.data;
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
  computed:{
    ...mapState('userModule', {
      userInfo: (state) => state.userInfo,
    }),
    ...mapState('contestModule', {
      contestInfo: (state) => state.contestInfo,
    }),
    rows() {
      return this.items.length
    },
    isUserInfoEmpty() {
        return !this.userInfo || Object.keys(this.userInfo).length === 0;
    }
  },
}    
</script>
<style scoped>
.box {
    width: 1200px;
    margin: 20px auto;
}
.problem{
    width: 1200px;
}
.Title{
    display: flex;
    justify-content: center;
    font-size: 22px;
}
.downtime{
    widows: 10px;
}
.curtime{
    display: flex;
}
.contesttime{
    display: flex;
    height: 20px;
    justify-content: center;
}
.status{
    display: flex;
}
.curstatus{
    display: flex;
    height: 20px;
    justify-content: center;
}
.beforestart{
    color:green; 
    font-size:30px; 
    display:flex;
    justify-content: center;
}
.problemList{
    margin-top: 20px;
}
.Button{
    margin: 10px;
    display: flex;
    justify-content: center;
}
.col-snum{
    width: 60px;
    display: flex;
    justify-content: center;
    flex-wrap: nowrap;
    overflow: hidden;
}
.col-user{
    width: 100px;
    display: flex;
    flex-wrap: nowrap;
    overflow: hidden;
    color: #006eff;
    cursor: pointer;
}
.col-pnum{
    width: 60px;
    display: flex;
    justify-content: center;
    flex-wrap: nowrap;
    overflow: hidden;
    color: #006eff;
    cursor: pointer;
}
.col-title{
    width: 200px;
    text-wrap: nowrap;
    overflow: hidden;
    color: #006eff;
    cursor: pointer;
}
.col-solve .col-submit{
    width: 50px;
    text-wrap: nowrap;
    display: flex;
    justify-content: center;
}
.col-title {
  color: #006eff;
  cursor: pointer;
}
.col-time{
    width: 150px;
}
.col-status{
    width: 180px;
    font-size: 2px;
}
.col-m .col-t .col-lang{
    width: 60px;
}
</style>
