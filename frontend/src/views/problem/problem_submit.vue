<template>
    <div class="submit">
        <h3>Submit:</h3>
        <b-card>
            <div v-if="!submission.File">
                <p>请选择编程语言:</p>
                <b-form-select v-model="submission.Lang" class = "choseLang">
                <option v-for="option in lang" :value="option" :key="option">{{ option }}</option>
                </b-form-select>
                <div>
                    <p style="margin-top:10px">源码:</p>
                    <b-form-textarea
                        id="textarea-no-resize"
                        v-model="submission.Code"
                        placeholder="请您将代码复制在此框内。如果选择上传文件提交代码，此框不需要填写，编程语言也不需要选择。仅支持下拉框下面的编程语言。"
                        rows="12"
                        no-resize
                    ></b-form-textarea>
                </div>
            </div>
            <div v-if="!(submission.Lang || submission.Code)">
                <p style="margin-top:10px">或 上传文件:</p>
                <b-form-file v-model="submission.File" class="mt-3" plain></b-form-file>
            </div>
                <div class="sBut">
                    <b-button @click="submit" variant="outline-secondary">提交代码</b-button>
                </div>
        </b-card>
    </div>
</template>
<script>
import { mapState, mapActions } from 'vuex';
import { formatCurrentDateTime } from '../../helper/getTime'
import { getProblemIdFromURL } from '../../helper/getUrl'
export default {
    data() {
        return {
        lang:[],
        submission:{
            UserId:'',
            Lang:'',
            ProblemId:'',
            Time:'',
            Code:'',
            File: null,
        }
        }
    },
  computed: {
    ...mapState('submitModule', ['languages']),
    ...mapState('userModule', {
      userInfo: (state) => state.userInfo,
    }),
  },
  created() {
    this.GetLang().then(() => {
        this.lang = this.languages
    }).catch((error) => {
      this.$bvToast.toast(error.response.data.msg, {
            title: '请求错误',
            variant: 'danger',
            toaster:'b-toaster-bottom-right',
            solid: true,
            appendToast:true,
        });
    });
  },
  methods: {
    ...mapActions('submitModule', ['GetLang']),
    ...mapActions('submitModule', ['fetchSubmit']),
    submit(){
        this.submission.UserId = this.userInfo.id;
        this.submission.Time = formatCurrentDateTime();
        this.submission.ProblemId = getProblemIdFromURL('pid');
        if (this.submission.File || this.submission.Code.length * this.submission.Lang.length){
            this.fetchSubmit(this.submission).then((response) => {
                this.$bvToast.toast(response.data.msg, {
                title: '请求成功',
                variant: 'success',
                toaster:'b-toaster-bottom-right',
                solid: true,
                appendToast:true,
                });
                setTimeout(() => {
                // 一秒后进行页面跳转
                this.$router.push('/problem/status');
                }, 1000);
            }).catch((error) => {
            this.$bvToast.toast(error.response.data.msg, {
                title: '请求错误',
                variant: 'danger',
                toaster:'b-toaster-bottom-right',
                solid: true,
                appendToast:true,
                });
            });
        } else {
            this.$bvToast.toast("请输入正确格式", {
            title: '数据验证错误',
            variant: 'danger',
            toaster:'b-toaster-bottom-right',
            solid: true,
            appendToast:true,
            });
        }
    }
  },
}
</script>
<style scoped>
.submit{
    margin: 30px auto;
    width: 70%;
}
.choseLang{
    width: 40%;
}
p {
    margin-top: 10px;
}
.sBut{
    margin-top: 10px;
    display: flex;
    justify-content: center;
}
</style>