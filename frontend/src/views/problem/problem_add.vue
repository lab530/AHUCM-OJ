<template>
   <div class="box">
    <b-alert show>           
      <b-button variant="outline-primary" size="sm" class = "back" @click="back()">
        <b-icon icon="arrow-left-short"></b-icon> 返回
      </b-button> 问题添加
    </b-alert>
      <div class = "problem-add">
          <b-card>
              <div>
                  <b-form  @submit.stop.prevent>
                      <label for="feedback-user">标题：</label>
                      <b-form-input 
                          v-model="problem.Title" 
                          :state="titleState" 
                          id="feedback-user"
                          ref="input"
                          placeholder="请输入题目标题">
                      </b-form-input>
                      <b-form-invalid-feedback :state="titleState">
                          题目标题不少于 1 个字符
                      </b-form-invalid-feedback>
                      <b-form-valid-feedback :state="titleState">
                          Looks Good.
                      </b-form-valid-feedback>
                  </b-form>
              </div>
              <div>
                  <label class="mt-2">题目时间限制和内存限制：</label>
                  <b-form inline>
                      <b-form-input 
                      class="mb-2 mr-sm-2 mb-sm-0"
                      v-model="problem.TimeLimit"
                      placeholder="单位 （ms）">
                      </b-form-input>
                      <div class="mb-2 mr-sm-2 mb-sm-0"> ms</div>
                      <div>&nbsp;&nbsp;&nbsp;</div>
                      <b-form-input 
                      class="mb-2 mr-sm-2 mb-sm-0"
                      v-model="problem.MemoLimit"
                      placeholder="单位 （MB）">
                      </b-form-input>
                      <div class="mb-2 mr-sm-2 mb-sm-0"> mb</div>
                  </b-form>
              </div>

              <div>
                  <label for="feedback-user"
                  class="mt-2">题目描述:
                  </label>
                  <mavon-editor v-model="problem.Description"/>
              </div>
              <div>
                  <label for="feedback-user"
                  class="mt-2">输入描述:
                  </label>
                  <mavon-editor v-model="problem.Input"/>
              </div>
              <div>
                  <label for="feedback-user"
                  class="mt-2">输出描述:
                  </label>
                  <mavon-editor v-model="problem.Output"/>
              </div>
              <div>
                  <label for="feedback-user"
                  class="mt-2">样例输入:
                  </label>
                  <b-form-textarea
                      id="textarea"
                      v-model="problem.SimpleInput"
                      placeholder="没有可不填。"
                      rows="3"
                      max-rows="6"
                  ></b-form-textarea>
              </div>
              <div>
                  <label for="feedback-user"
                  class="mt-2">样例输出:
                  </label>
                  <b-form-textarea
                      id="textarea"
                      v-model="problem.SimpleOutput"
                      placeholder="没有可不填。"
                      rows="3"
                      max-rows="6"
                  ></b-form-textarea>
              </div>
              <div>
                  <label for="feedback-user"
                  class="mt-2">数据范围与提示:
                  </label>
                  <mavon-editor v-model="problem.Illustrate"/>
              </div>
              <div>
              <b-button block variant="primary" class="mt-3" @click="add()">添加题目</b-button>
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
      problem:{
        Title: '',
        UserId: '',
        Description:'',
        Input:'',
        Output:'',
        SimpleInput:'',
        SimpleOutput:'',
        Illustrate:'',
        TimeLimit:'',
        MemoLimit:'',
      }
    }
  },
  computed: {
    ...mapState('userModule', {
      userInfo: (state) => state.userInfo,
    }),
    titleState() {
      if(this.problem.Title.length === 0) return null;
      return this.problem.Title.length > 0;
    },
  },
  methods:{
    ...mapActions('problemModule', { addProblem:'add' }),
    add(){
      this.problem.UserId = this.userInfo.id;
      console.log(this.problem)
      this.addProblem(this.problem).then(() => {
        // 请求成功
        console.log("problem add");
        this.$router.push('/admin/problem');
      }).catch((error) => {
        this.$bvToast.toast(error.response.data.msg, {
          title: '数据验证错误',
          variant: 'danger',
          solid: true,
        });
      });      
    },
    back(){
      this.$router.go(-1)
    },
  },
  mounted(){
    this.$nextTick(() => {
      this.$refs.input.focus();
    })
  }
}
</script>

<style scoped>
.box{
  margin-left: 200px;
}
.title{
  display: flex;
  justify-content: center;
}
.problem-add{
  width: 1000px;
  margin: 0 auto;
  box-shadow: 1px;
}
.blog-li{
  padding: 10px;
  display: flex;
}
.cover-img{
  width: 200px;
  height: 200px;
}
.blog-content{
  width: 65;
  padding: 20px;
  position: relative;
}
.writer-info{
position: absolute;
right: 0;
bottom: 0;

}
.img-box{
width: 35%;
}

</style>