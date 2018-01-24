// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import VueResource from 'vue-resource'
import VueI18n from 'vue-i18n'
import App from './App'
import router from './router'

import 'semantic-css'
import 'semantic'

Vue.config.productionTip = true
Vue.use(VueResource)
Vue.use(VueI18n)

// get browser lang and init locales
var userLang = navigator.language || navigator.userLanguage;
userLang = userLang.replace('-', '_').split('_')[0];
var messages = require('./i18n/locale-' + userLang + '.json')
const i18n = new VueI18n({
  locale: userLang,
  messages,
})

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  i18n,
  components: {
    App
  },
  template: '<App/>',
  hotspot: {
    digest: '',
    uuid: '',
    sessionid: '',
    onError: false,
    preferences: {},
    socials: {}
  }
})