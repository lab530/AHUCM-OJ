<template>
    <b-card class = "box">
        <div>
            <BottomTime />
        </div>
        <b-card class="contest_card">
            <div>
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
                  <div class="col-name" @click="navigateToContest(data.item.ID)">{{ data.item.title }}</div>
                </template>

                <template #cell(比赛状态)="data">
                  <div class="col-status">
                    <div v-if="isBeforeStart(data.item)">  
                      开始于 {{ formatDate(data.item.start_at) }}，总赛时 {{ getTotalDuration(data.item.start_at, data.item.end_at) }}  
                    </div>  
                    <div v-else-if="isDuring(data.item)" class="status-ongoing">  
                      <p style="color:brown; display:inline;">进行中</p>，结束于 {{ formatDate(data.item.end_at) }}  
                    </div>  
                    <div v-else>  
                      已结束 {{ formatDate(data.item.end_at) }}  
                    </div>
                    </div>  
                </template>
               
                <template #cell(开放)="data">
                  <div class="col-open" v-if=data.item.public>公开</div>
                  <div class="col-open" v-else>私有</div>
                </template>

                <template #cell(创建人)="data">
                  <div class = "col-create" v-if="data.item.user.user_name">{{ data.item.user.user_name }}</div>
                  <div v-else class = "col-create">YesterCafe</div>
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

    </b-card>
</template>

<script>
import BottomTime from '../bottom/bottom_time.vue';
import { mapState, mapActions } from 'vuex';
import { isDuring, isBeforeStart, formatDate, getTotalDuration } from '../../helper/getTime';
export default {
  components: {
    BottomTime,
  },
  data() {
    return {
      // data: [],
      perPage: 10,
      currentPage: 1,
      fields:['比赛编号', '比赛名称', '比赛状态','开放','创建人'],
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
    navigateToContest(contestId) {
      this.$router.push('/contest?cid=' + contestId);
    },
    formatDate,
    getTotalDuration,
    isBeforeStart,
    isDuring,
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
.box{
    width: 1300px;
    margin: 10px auto;
    background-image: url("../../../public/bg.png");
    background-repeat: repeat;
    background-size: 150px 150px; /* 自定义宽度为800像素，高度为400像素 */
}

.contest_card{
    margin-top: 15px;
}

.col-id{
  width: 85px;
}

.col-name{
  width: 200px;
  color: #006eff;
  cursor: pointer;
  overflow: hidden;
  white-space: nowrap; /* 防止文本换行 */
}

.col-status{
  width: 400px;
  overflow: hidden;
}

.col-open{
  width: 50px;
}

.col-create{
  width: 100px;
  white-space: nowrap; /* 防止文本换行 */
  overflow: hidden;
}
</style>
