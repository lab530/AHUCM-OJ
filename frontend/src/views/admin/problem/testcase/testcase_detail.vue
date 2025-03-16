<template>
  <div>
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
    <div class="box">
      <b-card>
        <div v-if="!content">
          <b-button variant="outline-primary">
            <b-icon
              icon="arrow-counterclockwise"
              animation="spin-reverse"
            ></b-icon>&nbsp;数据加载中
          </b-button>
        </div>
        <b-form-textarea
          v-else
          id="textarea-no-resize"
          rows="20"
          no-resize
          v-model="content"
        ></b-form-textarea>
        <div class="save-container">
          <b-button
            variant="outline-primary"
            size="sm"
            class="save"
            @click="save()"
          >
            <b-icon icon="save"></b-icon> 保存文件
          </b-button>
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
      content: '',
    }
  },
  async created() {
    try {
        const response = await this.GetCaseInfo();
        this.content = response.data.data.data;
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
    ...mapActions('adminModule', ['GetCaseInfo']),
    ...mapActions('adminModule', ['UpdateCase']),
    GoEdit(name){
      this.$router.push('/admin/testdetail?pid='+ this.pid + '&fname=' +  name);
    },
    back(){
      this.$router.go(-1)
    },
    async save() {
      try {
          await this.UpdateCase(this.content);
          this.$bvToast.toast("文件上传成功", {
              title: '上传成功',
              variant: 'success',
              toaster: 'b-toaster-bottom-right',
              solid: true,
              appendToast: true,
          });
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
    }
  },
}
</script>
<style scoped>
.back {
  margin-right: 10px;
}

.box {
  width: 1200px;
  height: 600px;
  margin: 20px auto;
}

.save-container {
  display: flex;
  justify-content: flex-end;
  align-items: flex-end;
  margin-top: 10px;
}
</style>