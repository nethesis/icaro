// The following line loads the standalone build of Vue instead of the runtime-only build,
// so you don't have to do: import Vue from 'vue/dist/vue'
// This is done with the browser options. For the config, see package.json
import Vue from 'vue'
import VueResource from 'vue-resource'
import VueI18n from 'vue-i18n'

import App from './App.vue'
import router from './routes/router'
import languages from './i18n/lang'

Vue.config.productionTip = true
Vue.use(VueResource)
Vue.use(VueI18n)

// configure i18n
var langConf = languages.initLang()
const i18n = new VueI18n({
  locale: langConf.locale,
  messages: langConf.messages,
})

// init Vue app
new Vue({
  el: '#app',
  router,
  i18n,
  render: (h) => h(App)
})