<template>
    <div class="problemlist">
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
                  <div class="col-num">{{ 1000 + data.item.ID }}</div>
                </template>  
                
                <template #cell(用户)="data">
                  <div class="col-title" @click="navigateToProblem(data.item.ID)">{{ data.item.user.user_name }}</div>
                </template>

                <template #cell(题目编号)="data">
                  <div class="col-num" @click="navigateToProblem(data.item.ID)">{{ 1000 + data.item.problem_id }}</div>
                </template>


                <template #cell(结果)="data">
                  <div class="col-status" >
                    <b-button v-if="data.item.status == 0" variant="warning" size="sm">
                        <b-icon icon="arrow-counterclockwise" animation="spin-reverse"></b-icon> Pending
                    </b-button>
                    <b-button v-else-if="data.item.status == 1" variant="warning" size="sm">
                        <b-icon icon="arrow-counterclockwise" animation="spin-reverse"></b-icon> PR
                    </b-button>
                    <b-button v-else-if="data.item.status == 2" variant="warning" size="sm">
                        <b-icon icon="arrow-counterclockwise" animation="spin-reverse"></b-icon> Compiling
                    </b-button>
                    <b-button v-else-if="data.item.status == 3" variant="warning" size="sm">
                        <b-icon icon="arrow-counterclockwise" animation="spin-reverse"></b-icon> Running
                    </b-button>
                    <b-button v-else-if="data.item.status == 4" variant="success" size="sm">
                        <b-icon icon="check2"></b-icon> AC
                    </b-button> 
                    <b-button v-else-if="data.item.status == 5" variant="info" size="sm">
                        <b-icon icon="calendar3-week"></b-icon> PE
                    </b-button> 
                    <b-button v-else-if="data.item.status == 6" variant="danger" size="sm">
                        <b-icon icon="x-circle"></b-icon> WA
                    </b-button> 
                    <b-button v-else-if="data.item.status == 7" variant="dark" size="sm">
                        <b-icon icon="hourglass-split"></b-icon> TLE
                    </b-button> 
                    <b-button v-else-if="data.item.status == 8" variant="dark" size="sm">
                        <b-icon icon="file-earmark"></b-icon> MLE
                    </b-button> 
                    <b-button v-else-if="data.item.status == 9" variant="info" size="sm">
                        <b-icon icon="collection"></b-icon> OLE
                    </b-button> 
                    <b-button v-else-if="data.item.status == 10" variant="dark" size="sm">
                        <b-icon icon="exclamation-triangle"></b-icon> RE
                    </b-button> 
                    <b-button v-else-if="data.item.status == 11" variant="info" size="sm">
                        <b-icon icon="bug"></b-icon> CE
                    </b-button> 
                    <b-button v-else variant="outline-primary" size="sm">
                        <b-icon icon="question"></b-icon> UE
                    </b-button> 
                  </div>
                </template>

                <template #cell(内存)="data">
                  <div class="col-num" >{{ data.item.memo_limit }}</div>
                </template>

                <template #cell(耗时)="data">
                  <div class="col-num" >{{ data.item.time_limit }} ms</div>
                </template>

                <template #cell(语言)="data">
                  <div class="col-title" @click="navigateToProblem(data.item.ID)">{{ data.item.lang }}</div>
                </template>

                <template #cell(提交时间)="data">
                  <div class="col-time">    
                    {{
                        formatDate(new Date(data.item.submit_time)) +
                        ' ' +
                        formatTime(new Date(data.item.submit_time))
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
</template>
<script>
import { mapState, mapActions } from 'vuex';
export default {
  data() {
    return {
      // data: [],
      perPage: 50,
      currentPage: 1,
      fields:['提交编号', '用户', '题目编号','结果', '内存', '耗时', '语言', '提交时间' ],
      items: [],
    }
  },
  computed: {
    ...mapState('submitModule', ['SubmitList']),
    rows() {
      return this.items.length
    }
  },
  created() {
    this.fetchSubmitList()
        .then(() => {
            this.items = this.SubmitList.data.data;
        })
        .catch((error) => {
        this.$bvToast.toast(error.response.data.msg, {
            title: '数据验证错误',
            variant: 'danger',
            solid: true,
        });
        });
    },

  methods: {
    ...mapActions('submitModule', ['fetchSubmitList']),
    navigateToProblem(problemId) {
      this.$router.push('/problem?pid=' + problemId);
    },
    navigateToProblemList(id) {
      const path = '/problemset?category=' + id;
      if (this.$route.path !== path) {
        this.$router.push(path).catch(error => {
          if (error.name !== 'NavigationDuplicated') {
            console.warn('导航出错:', error);
          }
        });
      }
    },
    formatDate(date) {
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');
        return `${year}/${month}/${day}`;
    },
    formatTime(date) {
        const hours = String(date.getHours()).padStart(2, '0');
        const minutes = String(date.getMinutes()).padStart(2, '0');
        const seconds = String(date.getSeconds()).padStart(2, '0');
        return `${hours}:${minutes}:${seconds}`;
    }
  }
}
</script>
<style scoped>
.problemlist{
    width: 70%;
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
.col-title {
  color: #006eff;
  cursor: pointer;
}
.col-time{
    width: 150px;
}
.col-status{
    width: 100px;
    overflow: hidden;
    font-size: 2px;
}
</style>