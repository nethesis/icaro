// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import VueResource from 'vue-resource'
import VueI18n from 'vue-i18n'
import VModal from 'vue-js-modal'
import VueAwesomeCountdown from 'vue-awesome-countdown'

import App from './App'
import router from './router'

import 'semantic-css'
import 'semantic'

Vue.config.productionTip = true
Vue.use(VueResource)
Vue.use(VueI18n)
Vue.use(VModal, {
  dialog: true
})

Vue.use(VueAwesomeCountdown, 'vac')

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

// init Vue app
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
    uamip: '',
    uamport: '',
    onError: false,
    preferences: {},
    socials: {},
    voucherValidated: false
  }
})
