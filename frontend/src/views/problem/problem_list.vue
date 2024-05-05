<template>
    <div class="problemlist">
        <b-card>
            <div class="overflow-auto" v-if="this.items">
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
                
                <template #cell(标题)="data">
                  <div class="col-title" @click="navigateToProblem(data.item.ID)">{{ data.item.title }}</div>
                </template>
                <template #cell(分类)="data">
                    <b-tag variant="primary" v-for="category in data.item.ProblemCategories"  :key="category.ID" no-remove>
                      <div  @click="navigateToProblemList(category.Category.ID)"> {{ category.Category['category-name'] }}</div>
                      </b-tag>
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
            <div v-else> 暂无数据 </div>
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
      fields:['题目编号', '标题', '分类'],
      items: [],
    }
  },
  computed: {
    ...mapState('problemModule', ['problemList']),
    rows() {
      return this.items.length
    }
  },
  watch: {
    '$route': {
      immediate: true,
      handler() {
        this.fetchProblemList().then(() => {
          this.items = this.problemList.data.data.data;
        }).catch((error) => {
          this.$bvToast.toast(error.response.data.msg, {
            title: '数据验证错误',
            variant: 'danger',
            solid: true,
          });
        });
      }
    }
  },
  methods: {
    ...mapActions('problemModule', ['fetchProblemList']),
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
</style>