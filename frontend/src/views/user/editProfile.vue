<template>
  <div class = "container">
      <b-card class = "first-card" title="修改个人信息">
        <hr>
          <div class = "second-container">
              <div class="userInfo">
                  <div>
                      <img :src="require('../../../../backend/static/icon/' + user.icon)"   class="icon" alt="用户信息获取失败">
                  </div>
                  <div>
                      <h3> Hello!</h3>
                      <h3><span v-if="userInfo && userInfo.user_name">{{ userInfo.user_name }}</span></h3>
                  </div>
              </div>
              <div class="second-card">上传新的头像:<br>
                <div>
                  <b-form-file
                    v-model="user.UserIcon"
                    :state="Boolean(user.UserIcon)"
                    placeholder="上传您的新头像，不更改可以不上传, 大小不超过 200 KB"
                  ></b-form-file>
                </div>
                <div role="group">
                  <label for="input-live1">昵称:</label>
                  <b-form-input
                    id="input-live1"
                    v-model="user.UserNickname"
                    :state="nameState"
                    placeholder="输入您的新昵称，不更改可以不填写"
                    trim
                  ></b-form-input>

                  <!-- This will only be shown if the preceding input has an invalid state -->
                  <b-form-invalid-feedback id="input-live-feedback">
                   昵称长度为 6 - 12 个字符
                  </b-form-invalid-feedback>
                </div>
                <div role="group">
                  <label for="input-live2">邮箱:</label>
                  <b-form-input
                    id="input-live2"
                    v-model="user.UserEmail"
                    :state="emailState"
                    placeholder="输入您的新邮箱，不更改可以不填写"
                    trim
                  ></b-form-input>

                  <!-- This will only be shown if the preceding input has an invalid state -->
                  <b-form-invalid-feedback id="input-live-feedback">
                    输入正确的邮箱格式
                  </b-form-invalid-feedback>
                </div>
                <div role="group">
                  <label for="input-live3">原始密码:</label>
                  <b-form-input
                    type="password"
                    v-model="user.UserPassword"
                    placeholder="输入您的新密码，不更改以下可以都不填写"
                    trim
                  ></b-form-input>

                </div>
                <div role="group">
                  <label for="input-live4">新的密码:</label>
                  <b-form-input
                    id="input-live4"
                    v-model="user.NewPassword"
                    type="password"
                    :state="passwordState"
                    placeholder="不更改密码可以不填写"
                    trim
                  ></b-form-input>

                  <!-- This will only be shown if the preceding input has an invalid state -->
                  <b-form-invalid-feedback id="input-live-feedback">
                    密码长度为 6 - 12 个字符
                  </b-form-invalid-feedback>
                </div>
                <div role="group">
                  <label for="input-live5">确认密码:</label>
                  <b-form-input
                    v-model="user.confirmPassword"
                    :state="confirmState"
                    placeholder="不更改密码可以不填写"
                    trim
                  ></b-form-input>

                  <!-- This will only be shown if the preceding input has an invalid state -->
                  <b-form-invalid-feedback id="input-live-feedback">
                    两次密码输入不一致
                  </b-form-invalid-feedback>
                </div>
                <br>
                <div
                >
                <b-button block variant="primary" @click="editProfile">更新个人信息</b-button>
                </div>
              </div>
          </div>
      </b-card>
  </div>
</template>
<script>
import { mapState, mapActions } from 'vuex';
import customValidator from '@/helper/validator';
export default {
  data() {
      return {
          user: {
          UserName: '',
          UserNickname:'',
          UserEmail: '',
          UserPassword:'',
          NewPassword: '',
          UserIcon:[],
          icon:'init.jpg',
          confirmPassword:'',
          IconUpload:'',
        },
      };
    },
  created() {
    if (this.userInfo) {
      this.user.icon = this.userInfo.user_icon ? this.userInfo.user_icon : "init.jpg";
      this.user.UserName = this.userInfo.user_name;
    }
  },
  computed: {
    ...mapState({
      userInfo: (state) => state.userModule.userInfo,
      validation() {
        return this.user.UserNickname.length >= 4 && this.user.UserNickname.length < 13
      },
      nameState() {
        if(this.user.UserNickname.length == 0) return null;
        return this.user.UserNickname.length > 3 ? true : false
      },
      emailState() {
        if(this.user.UserEmail.length == 0) return null;
        const invalid = customValidator.emailValidator(this.user.UserEmail);
        return invalid
      },
      passwordState() {
        if(this.user.UserPassword.length == 0) return null;
        return this.user.UserPassword.length >= 6 && this.user.UserPassword.length < 13
      },
      confirmState(){
        if(this.user.confirmPassword.length == 0) return null;
        return this.user.UserPassword == this.user.confirmPassword;
      }
  }),
  },
  methods: {
    ...mapActions('userModule', { userEdit: 'edit'}),
    editProfile(){
    this.userEdit(this.user).then(() => {
        // 跳转请求成功
          this.$router.replace({ name: 'profile' });
      }).catch((err) => {
        console.log(err)
        this.$bvToast.toast(err.response.data.msg, {
          title: '数据验证错误',
          variant: 'danger',
          solid: true,
        });
      });
      console.log('edit');
    },
  }
};
</script>
<style scoped>
.icon{
  width: 200px;
  height: 200px;
}

.container
{
  width: 50%;
  margin: 10px auto;
  
}

.userInfo {
  width: 40%;
}


.first-card{
  
}
.second-card{
  width: 60%;
}
.second-container{
  display: flex;
}

</style>
