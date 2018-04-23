// The following line loads the standalone build of Vue instead of the runtime-only build,
// so you don't have to do: import Vue from 'vue/dist/vue'
// This is done with the browser options. For the config, see package.json
import Vue from "vue";
import VueResource from "vue-resource";
import VueI18n from "vue-i18n";
import VueGoodTable from "vue-good-table";

import App from "./App.vue";
import router from "./routes/router";
import languages from "./i18n/lang";
import filters from "./filters/filters";

window.$ = window.jQuery = require("patternfly/node_modules/jquery/dist/jquery.min.js");
require("patternfly/node_modules/bootstrap/dist/js/bootstrap.min.js");
require("patternfly/dist/js/patternfly.min");
require("patternfly/dist/css/patternfly.min.css");
require("patternfly/dist/css/patternfly-additions.min.css");

Vue.config.productionTip = true;
Vue.use(VueResource);
Vue.use(VueI18n);
Vue.use(VueGoodTable);

// configure i18n
var langConf = languages.initLang();
const i18n = new VueI18n({
  locale: langConf.locale,
  messages: langConf.messages
});
var moment = require("patternfly/node_modules/moment/moment.js");
moment.locale(langConf.locale);

// import filters
Vue.filter("byteFormat", filters.byteFormat);
Vue.filter("secondsInHour", filters.secondsInHour);
Vue.filter("formatDate", filters.formatDate);

// init Vue app
new Vue({
  el: "#app",
  router,
  i18n,
  render: (h) => h(App),
  api_host: window.location.host,
  api_scheme: window.location.protocol + '//',
  currentLocale: langConf.locale,
  moment: moment
});
