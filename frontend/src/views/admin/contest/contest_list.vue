<template>
  <div class = "box">
      <div class="contestlist">
          <b-card>
              <b-button @click="$router.push('/admin/contest/add')" variant="primary">添加竞赛</b-button>
          </b-card>
          <b-card class="contest_card">
                  <b-table
                  id="my-table"
                  striped hover 
                  :items="items"
                  :fields="fields"
                  :per-page="perPage"
                  :current-page="currentPage"
                  >
                  <template #cell(比赛编号)="data">
                    <div class="col-id">{{ 1000 + data.item.ID }}</div>
                  </template>  
                  
                  <template #cell(比赛名称)="data">
                    <div class="col-name"><a href="#">{{ data.item.title }}</a></div>
                  </template>

                  <template #cell(开始时间)="data">
                    <div class = "status">
                        {{ formatDate(data.item.start_at) }} 
                    </div>
                  </template>
                
                  <template #cell(结束时间)="data">
                    <div class = "status">
                        {{ formatDate(data.item.end_at) }} 
                    </div>
                  </template>

                  <template #cell(开放)="data">
                    <div class="col-num" v-if=data.item.public>公开</div>
                    <div class="col-num" v-else>私有</div>
                  </template>

                  <template #cell(创建人)="data">
                    <div class = "col-num" v-if="data.item.user.user_name">{{ data.item.user.user_name }}</div>
                    <div v-else>佚名</div>
                  </template>
                  
                  <template #cell(操作)="data">
                    <div class = "col-op" @click="NavigateTo(data.item.ID)">编辑</div>
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
          </b-card>
      </div>
    </div>
</template>
<script>

import { mapState, mapActions } from 'vuex';
import { isDuring, isBeforeStart, formatDate, getTotalDuration } from '../../../helper/getTime';
export default {
  data() {
    return {
      // data: [],
      perPage: 10,
      currentPage: 1,
      fields:['比赛编号', '比赛名称', '开始时间','结束时间','开放','操作'],
      items: [],
    }
  },
  computed: {
    ...mapState('contestModule', ['contestList']),
    rows() {
      return this.items.length
    },
  },
  methods: {
    ...mapActions('contestModule', ['fetchContestList']),
    navigateToProblem(problemId) {
      this.$router.push('/admin/problemdetail?pid=' + problemId);
    },
    navigateToCase(problemId) {
      this.$router.push('/admin/testcase?pid=' + problemId);
    },
    navigateToEdit(problemId){
      this.$router.push('/admin/problem/edit?pid=' + problemId);
    },
    formatDate,
    getTotalDuration,
    isBeforeStart,
    isDuring,
    NavigateTo(id){
      this.$router.push('/admin/contest/edit?cid='+id);
    }
  },
  created() {
    this.fetchContestList().then((response) => {
      // 在获取问题列表数据后，更新 data 数组
      // this.data = this.problemList.data.data.data;
      this.items = response.data;
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
};
</script>

<style scoped>
.box {
    margin-left: 200px;
}
.contestlist{
  margin: 40px auto;
  max-width: 1400px;
}
.contest_card{
    margin-top: 40px;
    width: 100%;
}
.status {
  width: 150px;
}
.col-id{
  width: 60px;
}
.col-name{
  width: 200px;
  overflow: hidden;
  white-space: nowrap; /* 防止文本换行 */
}
.clo-num{
  width: 50px;
}
.col-create{
  width: 100px;
}
.col-op{
  width:50px;
  color:  #006eff;
  cursor: pointer;
}
</style>