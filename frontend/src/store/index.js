import Vue from 'vue';
import Vuex from 'vuex';
import userModule from './module/user'
import problemModule from './module/problem'
import submitModule from './module/submit';
import adminModule from './module/admin';

Vue.use(Vuex);

export default new Vuex.Store({
  strict: process.env.NODE_ENV !== 'production',
  state: {
  },
  getters: {
  },
  mutations: {
  },
  actions: {
  },
  modules: {
    userModule,
    problemModule,
    submitModule,
    adminModule,
  },
});
