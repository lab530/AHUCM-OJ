<template>
    <div class="register">
      <b-row class="mt-5">
        <b-col md="6" offset-md="2" lg="6" offset-lg="3">
          <b-card title="登录">
            <b-form>
              <b-form-group label = "用户名">
                  <b-form-input
                  v-model="$v.user.UserName.$model"
                  type="text"
                  placeholder="输入您的用户名"
                  :state="validateState('UserName')"
                  ></b-form-input>
                  <b-form-invalid-feedback :state="validateState('UserName')">
                    手机号不符合要求
                  </b-form-invalid-feedback>
              </b-form-group>
              <b-form-group label = "密码">
                  <b-form-input
                  v-model="$v.user.UserPassword.$model"
                  type="password"
                  placeholder="输入您的密码"
                  :state="validateState('UserPassword')"
                  ></b-form-input>
                  <b-form-invalid-feedback :state="validateState('UserPassword')">
                    密码必须大于等于 6 位
                  </b-form-invalid-feedback>
              </b-form-group>
              <b-form-group>
                <b-button variant="outline-primary" @click="login()" block>登录</b-button>
              </b-form-group>
            </b-form>
          </b-card>
        </b-col>
      </b-row>
    </div>
</template>

<script>
import { required, minLength } from 'vuelidate/lib/validators';


import { mapActions } from 'vuex';

export default {
  data() {
    return {
      user: {
        UserName: '',
        UserPassword: '',
      },
      validation: null,
    };
  },
  validations: {
    user: {
      UserName: {
        required,
      },
      UserPassword: {
        required,
        minLength: minLength(6),
      },
    },
  },
  methods: {
    ...mapActions('userModule', { userLogin: 'login' }),
    validateState(name) {
      const { $dirty, $error } = this.$v.user[name];
      return $dirty ? !$error : null;
    },
    login() {
      // 验证数据
      this.$v.user.$touch();
      if (this.$v.user.$anyError) {
        return;
      }
      // 请求
      this.userLogin(this.user).then(() => {
        // 跳转页面
        this.$router.replace({ name: 'Home' });
      }).catch((err) => {
        this.$bvToast.toast(err.response.data.msg, {
          title: '数据验证错误',
          variant: 'danger',
          solid: true,
        });
      });
    },
  },
};
</script>

<style lang="scss" scoped>

</style>
