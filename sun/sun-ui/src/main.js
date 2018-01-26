// The following line loads the standalone build of Vue instead of the runtime-only build,
// so you don't have to do: import Vue from 'vue/dist/vue'
// This is done with the browser options. For the config, see package.json
import Vue from 'vue'
import VueResource from 'vue-resource'
import VueI18n from 'vue-i18n'
import App from './App.vue'
import router from './routes/router'

Vue.config.productionTip = true
Vue.use(VueResource)
Vue.use(VueI18n)

// get browser lang and init locales
var messages = {}
var userLang = navigator.language || navigator.userLanguage;
userLang = userLang.replace('-', '_').split('_')[0];

// try loading browser locale
try {
  messages = require('./i18n/locale-' + userLang + '.json')
} catch (e) {
  console.log('locale: '+ userLang +' not found. fallback to en')
  messages = require('./i18n/locale-en.json')
  userLang = 'en'
}

// configure i18n
const i18n = new VueI18n({
  locale: userLang,
  messages,
})

// init Vue app
new Vue({
  el: '#app',
  router,
  i18n,
  render: (h) => h(App)
})