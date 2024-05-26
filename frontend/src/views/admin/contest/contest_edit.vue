<template>
    <div style="he">
      <b-alert show>           
        <b-button variant="outline-primary" size="sm" class = "back" @click="back()">
          <b-icon icon="arrow-left-short"></b-icon> 返回
        </b-button> 编辑竞赛
      </b-alert>
      <div class="box">
        <b-card>
            <div>
                <b-form  @submit.stop.prevent>
                    <label for="feedback-user">标题：</label>
                    <b-form-input 
                        v-model="contest.title" 
                        :state="titleState" 
                        id="feedback-user"
                        ref="input"
                        placeholder="请输入竞赛标题">
                    </b-form-input>
                </b-form>
            </div>
            <br>
            <div class = "time">
                <div class="start">
                    <b-form-group label="选择比赛开始时间：">
                    <input
                        type="datetime-local"
                        v-model="contest.start_at"
                    >
                    </b-form-group>
                </div>
                <div class="end">
                    <b-form-group label="选择比赛结束时间：">
                    <input
                        type="datetime-local"
                        v-model="contest.end_at"
                    >
                    </b-form-group>
                </div>
            </div>
              <div class="addProblem">
                <label for="feedback-user">题目编号（每个编号用 , 分隔开）：</label>
                <b-form-input v-model="problemIDs" placeholder="输入题号用英文逗号分割开。没有显示的问题将不会被添加到比赛当中"></b-form-input>
                <ul class="problemList">
                  <li v-for="problem in filteredProblems" :key="problem.ID" class="listItem">
                    <a v-if="problem.title" :href="problem.url" target="_blank">{{ problem.ID }}: {{ problem.title }}</a>
                  </li>
                </ul>
              </div>
            <div>
                <label for="feedback-user"
                class="mt-2">竞赛描述:
                </label>
                <mavon-editor v-model="contest.description"/>
            </div>
            <div class="Edit">
                <div class="isPublic">
                    比赛是否公开：
                    <b-form-select v-model="contest.Public" :options="options" class="mb-3">
                    <!-- This slot appears above the options from 'options' prop -->
                    <template #first>
                        <b-form-select-option :value="null" disabled>-- 选择比赛是否公开 --</b-form-select-option>
                    </template>
                    
                    </b-form-select>
                </div>
                <div class="password">
                    竞赛密码：
                    <b-form-input v-model="contest.password" type="text"></b-form-input>
                </div>
            </div>
            <div v-show="!contest.Public">
              私人竞赛参赛名单:
              <b-form-textarea v-model="contest.Participants" debounce="500" rows="4" max-rows="5"
              placeholder="注：私人竞赛只能通过后台添加名单的方式参加竞赛，请输入参加竞赛的用户名用英文,分隔开。例如: PlumYu,Test "
              ></b-form-textarea>
            </div>
            <div>
            <b-button block variant="primary" class="mt-3" @click="update()">更新竞赛</b-button>
            </div>
        </b-card>
      </div>
    </div>
</template>
<script>
import { mapState, mapActions } from 'vuex';
import {convertToBackendTime} from '../../../helper/getTime'
export default {
    data() {
    return {
      problemIDs: '',
      problems: [],
      contest:{
        ID:'',
        description:'',
        end_at: '',
        password: '',
        Public: '',
        start_at:'', 
        title:'',
        UserId:'',
        Participants:'',
        ProblemIDs:'',
      },
      options: [
          { value: 1, text: '公开' },
          { value: 0, text: '私有' }
      ],
      idOffset: 1000, // 前端输入框 ID 与后端的差异
    }
  },
  created(){
    this.GetContestDetail().then((response) => {
        this.contest = response.data.data.data.data
        this.contest.Participants = response.data.data.data.users
        this.problemIDs = response.data.data.data.problems
        this.contest.start_at = this.removeTimeZone(this.contest.start_at)
        this.contest.end_at = this.removeTimeZone(this.contest.end_at)
        
      }).catch((error) => {
        this.$bvToast.toast(error.response.data.msg, {
            title: '消息通知',
            variant: 'danger',
            toaster:'b-toaster-bottom-right',
            solid: true,
            appendToast:true,
        });
    });
  },
  mounted(){
    this.$nextTick(() => {
    this.$refs.input.focus();
    });
    if (!this.problemList || this.problemList.length === 0) {
      this.fetchProblemList().then(() => {
        this.processProblemList();
      });
    } else {
      this.processProblemList();
    }
  },
  computed: {
    ...mapState('problemModule', {
      problemList:(state) => state.problemList
    }),
    ...mapState('userModule', {
      userInfo: (state) => state.userInfo,
    }),
    titleState() {
      if(this.contest.title.length === 0) return null;
      return this.contest.title.length > 0;
    },
    parsedProblemIDs() {
      return this.problemIDs.split(',').map(id => id.trim()); 
    },
    filteredProblems() {
      const problemIDs = this.problemIDs.split(',').map(id => id.trim()); // 解析题号字符串为题号数组
      return this.problems.filter(problem => problemIDs.includes(problem.ID.toString()));
    },
  },
  methods: {
    ...mapActions('adminModule', ['GetContestDetail']),
    ...mapActions('contestModule',['editContest']),
    ...mapActions('problemModule', ['fetchProblemList']),
    back(){
      this.$router.go(-1)
    },
    removeTimeZone(dateTimeString) {  
        // 假设 dateTimeString 是 "2024-05-24T21:45:00+08:00" 这样的字符串  
        return dateTimeString.split('+')[0]; // 使用 '+' 作为分隔符，并取第一个部分  
    },  
    update(){
      if(!this.contest.start_at || !this.contest.end_at ||  this.contest.start_at > this.contest.end_at){
        this.$bvToast.toast("请输入正确的比赛结束时间", {
              title: '消息通知',
              variant: 'danger',
              toaster:'b-toaster-bottom-right',
              solid: true,
              appendToast:true,
          });
      } else {
        this.contest.start_at = convertToBackendTime(this.contest.start_at);
        this.contest.end_at = convertToBackendTime(this.contest.end_at);
        this.contest.ProblemIDs = this.problemIDs
        this.editContest(this.contest).then((res) => {
          // 请求成功
          console.log("Success")
          this.$bvToast.toast(res.data.msg, {
              title: '消息通知',
              variant: 'success',
              toaster:'b-toaster-bottom-right',
              solid: true,
              appendToast:true,
            });

          setTimeout(() => {  
              this.$router.push('/admin/contest');  
          }, 500);
        }).catch((error) => {
          this.$bvToast.toast(error.response.data.msg, {
              title: '消息通知',
              variant: 'danger',
              toaster:'b-toaster-bottom-right',
              solid: true,
              appendToast:true,
          });
        });
      }
    },
    processProblemList() {
      // 获取当前主机名
      const hostname = window.location.hostname;
      // 获取当前端口号
      const port = window.location.port;
      // 构建基础 URL
      const baseURL = `http://${hostname}:${port}`;

      // 定义函数来生成问题的 URL
      const buildProblemURL = (problemID) => `${baseURL}/problem?pid=${problemID}`;
      // 将数据转换为适合前端使用的格式
      this.problems = this.problemList.data.data.data.map((problem) => {
        return {
          ID: problem.ID + this.idOffset,
          title: problem.title,
          url: buildProblemURL(problem.ID),
        };
      });
    },
  },  
}
</script>
<style scoped>
.back{
  margin-right: 10px;
}

.box{
  width: 1200px;
  height: 100%;
  margin: 50px auto;
}

.save-container {
  display: flex;
  justify-content: flex-end;
  align-items: flex-end;
  margin-top: 10px;
}
.time{
    display: flex;
}
.start{

}
.end{
    margin-left: 100px;
}
.isPublic{
    width:200px;
}
.Edit{
    margin-top: 25px;
    display: flex;
}
.problemList {
  margin-top: 10px;
  list-style-type: none;
  padding-left: 0;
}

.listItem {
  text-align: left;
}
.password{
  margin-left: 200px;
}
</style>