<template>
  <div class="register">
    <b-row class="mt-5">
      <b-col
        md="6"
        offset-md="2"
        lg="6"
        offset-lg="3"
      >
        <b-card title="设置注册信息">
          <b-form>
            <!-- 用户名 -->
            <b-form-group label="用户名">
              <b-form-input
                class="inline"
                v-model="$v.user.UserName.$model"
                type="text"
                placeholder="用户名（学号）*"
                :state="validateState('UserName')"
              >
                <b-form-invalid-feedback :state="validateState('UserName')">
                  用户名长度必须大于 4 位
                </b-form-invalid-feedback>
              </b-form-input>
            </b-form-group>

            <!-- 昵称 -->
            <b-form-group label="昵称">
              <b-form-input
                v-model="$v.user.UserNickname.$model"
                type="text"
                placeholder="输入您的昵称"
                :state="validateState('UserNickname')"
              >
              </b-form-input>
            </b-form-group>

            <b-form-group label="邮箱">
              <b-form-input
                v-model="$v.user.UserEmail.$model"
                type="email"
                placeholder="输入您的电子邮箱"
                :state="validateState('UserEmail')"
              ></b-form-input>
              <b-form-invalid-feedback :state="validateState('UserEmail')">
                邮箱不符合要求
              </b-form-invalid-feedback>
            </b-form-group>

            <!-- <b-form-group label = "验证码">
                <div>
                  <b-form inline>
                    <b-form-input

                      class="mb-4 mr-sm-4 mb-sm-0"
                      placeholder="请输入验证码"
                    ></b-form-input>
                    <b-button variant="primary">发送</b-button>
                  </b-form>
                </div>
              </b-form-group> -->

            <!-- 密码 -->
            <b-form-group label="密码">
              <b-form-input
                v-model="$v.user.UserPassword.$model"
                type="password"
                placeholder="输入您的密码 *"
                :state="validateState('UserPassword')"
              ></b-form-input>
              <b-form-invalid-feedback :state="validateState('UserPassword')">
                密码必须大于等于 6 位
              </b-form-invalid-feedback>
            </b-form-group>
            <!-- 确认密码 -->
            <b-form-group label="确认密码">
              <b-form-input
                v-model="$v.user.confirmPassword.$model"
                type="password"
                placeholder="请再次您的密码 *"
                :state="validateState('confirmPassword')"
              ></b-form-input>
              <b-form-invalid-feedback :state="validateState('confirmPassword')">
                两次密码输入不一致
              </b-form-invalid-feedback>
            </b-form-group>

            <b-form-group>
              <b-button
                variant="outline-primary"
                @click="register()"
                block
              >register</b-button>
            </b-form-group>
          </b-form>
        </b-card>
      </b-col>
    </b-row>
  </div>
</template>

<script>
import { required, minLength } from 'vuelidate/lib/validators';

import customValidator from '@/helper/validator';

import { mapActions } from 'vuex';


export default {
  data() {
    return {
      user: {
        UserName: '',
        UserNickname:'',
        UserEmail: '',
        UserPassword: '',
        UserIcon:'',
        confirmPassword:'',
      },
      validation: null,
    };
  },
  validations: {
    user: {
      UserName: {
        required,
        minLength: minLength(4),
      },
      UserNickname: {

      },
      UserEmail: {
        required,
        email: customValidator.emailValidator,
      },
      UserPassword: {
        required,
        minLength: minLength(6),
      },
      confirmPassword: {
        required,
      }
    },
  },
  methods: {
    ...mapActions('userModule', { userRegister: 'register'}),
    validateState(name) {
      const { $dirty, $error } = this.$v.user[name];
      if( name == 'confirmPassword' ) {
        if(!$dirty) return null;
        else {
          if(this.user.UserPassword == this.user.confirmPassword) return true;
          else return false;
        }
      }
      return $dirty ? !$error : null;
    },
    async register() {
      // 验证数据
      this.$v.user.$touch();
      if (this.$v.user.$anyError) {
          return;
      }
      try {
          // 请求
          await this.userRegister(this.user);
          // 跳转请求成功
          this.$router.replace({ name: 'Home' });
      } catch (err) {
          console.log(err);
          this.$bvToast.toast(err.response ? err.response.data.msg : '未知错误', {
              title: '数据验证错误',
              variant: 'danger',
              solid: true,
          });
      }
      console.log('register');
    }
  },
};
</script>

<style lang="scss" scoped>
.register {
  max-width: 1300px;
  margin: 100px auto;
}
</style>