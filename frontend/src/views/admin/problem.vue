<template>
    <div class = "box">
        <div class="problemlist">
            <b-card>
                <b-button @click="$router.push('/admin/problem/add')" variant="primary">添加问题</b-button>
            </b-card>
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
                    <template #cell(题目编号)="data">
                    <div class="col-num">{{ 1000 + data.item.ID }}</div>
                    </template>  
                    
                    <template #cell(题目标题)="data">
                      <div class="col-title" @click="navigateToProblem(data.item.ID)">{{ data.item.title }}</div>
                    </template>
                    <template #cell(编辑题目)="data">
                      <div class="col-title" @click="navigateToEdit(data.item.ID)">编辑</div>
                    </template>
                    <template #cell(编辑测试点)="data">
                        <div class="col-title" @click="navigateToCase(data.item.ID)">TestCase</div>
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
</template>
<script>
import { mapState, mapActions } from 'vuex';
export default {
  data() {
    return {
      // data: [],
      perPage: 50,
      currentPage: 1,
      fields:['题目编号', '题目标题','编辑题目', '编辑测试点'],
      items: [],
    }
  },
  computed: {
    ...mapState('problemModule', ['problemList']),
    rows() {
      return this.items.length
    }
  },
  created() {
    this.fetchProblemList().then(() => {
      // 在获取问题列表数据后，更新 data 数组
      // this.data = this.problemList.data.data.data;
      this.items = this.problemList.data.data.data;
    }).catch((error) => {
      this.$bvToast.toast(error.response.data.msg, {
          title: '数据验证错误',
          variant: 'danger',
          solid: true,
        });
    });
  },
  methods: {
    ...mapActions('problemModule', ['fetchProblemList']),
    navigateToProblem(problemId) {
      this.$router.push('/admin/problemdetail?pid=' + problemId);
    },
    navigateToCase(problemId) {
      this.$router.push('/admin/testcase?pid=' + problemId);
    },
    navigateToEdit(problemId){
      this.$router.push('/admin/problem/edit?pid=' + problemId);
    },
  },
}
</script>
<style scoped>
.box {
    margin-left: 200px;
}
.problemlist{
    width: 100%; /* 子组件宽度减去导航栏宽度 */
    margin: auto;
}

.col-num {
  width: 50px; /* 设置题目编号列的宽度 */
  padding-left: 10px;
}

.tag-a {
  text-decoration: none;
  color: white;
}
.col-title{
  color: #006eff;
  cursor: pointer;
}
</style>