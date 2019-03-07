import Vue from 'vue'
import VueResource from 'vue-resource'
import VueI18n from 'vue-i18n'

//import SuiVue from 'semantic-ui-vue'
import 'semantic-ui-css/semantic.min.css'

window.$ = require('jquery')

import App from './App.vue'
import router from './router'

import "./filters/filters";

Vue.config.productionTip = true
Vue.use(VueResource)
Vue.use(VueI18n)
//Vue.use(SuiVue);

// get browser lang and init locales
var messages = {}
var userLang = navigator.language || navigator.userLanguage;
userLang = userLang.replace('-', '_').split('_')[0];

// try loading browser locale
try {
  messages = require('./i18n/locale-' + userLang + '.json')
} catch (e) {
  console.info('locale: ' + userLang + ' not found. fallback to en')
  messages = require('./i18n/locale-en.json')
  userLang = 'en'
}

// configure i18n
var i18n = new VueI18n({
  locale: userLang,
  messages,
})

import UtilService from "./services/util"
Vue.mixin(UtilService)

new Vue({
  router,
  i18n,
  render: function (h) {
    return h(App)
  },
  api_host: window.location.host,
  api_scheme: window.location.protocol + '//',
}).$mount('#app')