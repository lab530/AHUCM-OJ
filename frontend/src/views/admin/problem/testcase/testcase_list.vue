<template>
  <div class="box">
    <b-alert show>
      <b-button
        variant="outline-primary"
        size="sm"
        class="back"
        @click="back()"
      >
        <b-icon icon="arrow-left-short"></b-icon> 返回
      </b-button> 编辑测试点
    </b-alert>
    <div class="upload">
      <b-form-file
        v-model="selectedFiles"
        multiple
        class="chooseBox"
      >
        <template
          slot="file-name"
          slot-scope="{ names }"
        >
          <b-badge variant="dark">{{ names[0] }}</b-badge>
          <b-badge
            v-if="names.length > 1"
            variant="dark"
            class="ml-1"
          >
            + {{ names.length - 1 }} More files
          </b-badge>
        </template>
      </b-form-file>
      <b-button
        variant="outline-primary"
        @click="uploadFiles"
      >
        <b-icon icon="cloud-upload"></b-icon> 上传
      </b-button>
    </div>
    <div class="problemlist">
      <b-card>
        <div
          v-if="items"
          class="overflow-auto"
        >
          <b-table
            id="my-table"
            striped
            hover
            :items="items"
            :fields="fields"
            :per-page="perPage"
            :current-page="currentPage"
          >
            <template #cell(名字)="data">
              <div class="col-num">{{ data.item.name }}</div>
            </template>

            <template #cell(修改时间)="data">
              <div class="col-num">{{ data.item.mod_time }}</div>
            </template>

            <template #cell(大小)="data">
              <div class="col-num">{{ data.item.bytes }} bytes</div>
            </template>

            <template #cell(文件类型)="data">
              <div class="col-num">{{data.item.file_type}}</div>
            </template>

            <template #cell(操作)="data">
              <div class="operate">
                <div class="col-title">
                  <b-button
                    size="sm"
                    variant="outline-primary"
                    @click="GoEdit(data.item.name)"
                  >
                    <b-icon icon="code-square"></b-icon> 编辑
                  </b-button>
                </div>

                <div class="col-title">
                  <b-button
                    size="sm"
                    variant="outline-primary"
                    @click="DeleteCase( data.item.name )"
                  >
                    <b-icon icon="trash"></b-icon> 删除
                  </b-button>
                </div>
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
        <div v-else>暂无测试点。</div>
      </b-card>
    </div>
  </div>
</template>
<script>
import { mapState, mapActions } from 'vuex';
export default {
  props: ['pid'],
  data() {
    return {
      perPage: 50,
      currentPage: 1,
      fields:['名字', '大小','修改时间', '文件类型', '操作'],
      items: [],
      selectedFiles: [], // 存储用户选择的文件
    }
  },
  computed: {
    ...mapState('adminModule', ['TestCaseList']),
    rows() {
      return this.items.length
    },
  },
  async created() {
    try {
        await this.fetchTestCaseList();
        // 在获取问题列表数据后，更新 data 数组
        this.items = this.TestCaseList.data;
        // console.log(this.items)
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
  methods: {
    ...mapActions('adminModule', ['fetchTestCaseList']),
    ...mapActions('adminModule', ['UploadFile']),
    ...mapActions('adminModule', ['DeleteTestCase']),
    async uploadFiles() {
      try {
          await this.UploadFile(this.selectedFiles);
          this.$bvToast.toast("文件上传成功", {
              title: '上传成功',
              variant: 'success',
              toaster: 'b-toaster-bottom-right',
              solid: true,
              appendToast: true,
          });
          this.selectedFiles = []; // 清空用户选择的文件
          await this.UpdateList();
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

    async UpdateList() {
        try {
            await this.fetchTestCaseList();
            // 在获取问题列表数据后，更新 data 数组
            this.items = this.TestCaseList.data;
        } catch (error) {
            console.log(error);
            this.$bvToast.toast(error.response ? error.response.data.msg : '未知错误', {
                title: '更新列表失败',
                variant: 'danger',
                toaster: 'b-toaster-bottom-right',
                solid: true,
                appendToast: true,
            });
        }
    },

    async DeleteCase(name) {
        try {
            await this.DeleteTestCase(name);
            this.$bvToast.toast("文件删除成功", {
                title: '删除成功',
                variant: 'success',
                toaster: 'b-toaster-bottom-right',
                solid: true,
                appendToast: true,
            });
            await this.UpdateList();
        } catch (error) {
            console.log(error);
            this.$bvToast.toast(error.response ? error.response.data.msg : '未知错误', {
                title: '删除失败',
                variant: 'danger',
                toaster: 'b-toaster-bottom-right',
                solid: true,
                appendToast: true,
            });
        }
    },
    GoEdit(name){
      this.$router.push('/admin/testdetail?pid='+ this.pid + '&fname=' +  name);
    },
    back(){
      this.$router.go(-1)
    }
  },
}
</script>
<style scoped>
.back {
  margin-right: 10px;
}

.box {
  margin-left: 200px;
}
.problemlist {
  width: 100%; /* 子组件宽度减去导航栏宽度 */
  margin: auto;
}

.col-num {
  width: 150px; /* 设置题目编号列的宽度 */
  overflow: hidden;
}

.tag-a {
  text-decoration: none;
  color: white;
}
.col-title {
  color: #006eff;
  cursor: pointer;
}
.operate {
  width: 200px;
  display: flex;
  justify-content: space-between;
}
.upload {
  margin-left: 2px;
  display: flex;
  margin-bottom: 20px;
}
.chooseBox {
  width: 600px;
}
</style>