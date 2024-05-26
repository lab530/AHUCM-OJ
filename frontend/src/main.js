import Vue from 'vue';
import { BootstrapVue, IconsPlugin } from 'bootstrap-vue';
import Vuelidate from 'vuelidate';
import axios from 'axios';
import VueAxios from 'vue-axios';
import App from './App.vue';
import router from './router';
import store from './store';
import mavonEditor from 'mavon-editor'
import 'mavon-editor/dist/css/index.css'
import './assets/scss/index.scss';


Vue.config.productionTip = false;
// Vue.use(hljs)
Vue.use(mavonEditor)
Vue.use(BootstrapVue);
Vue.use(IconsPlugin);
Vue.use(Vuelidate);
Vue.use(VueAxios, axios);


new Vue({
  router,
  store,
  render: (h) => h(App),
}).$mount('#app');
