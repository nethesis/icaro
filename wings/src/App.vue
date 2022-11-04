<template>
  <div id="app" class="ui segment container" :style="containerStyle">
    <div v-show="loading" class="ui active dimmer inverted">
      <div class="ui loader"></div>
    </div>
    <div v-show="!loading" class="ui center aligned">
      <h2 :style="titleStyle">{{hotspot.preferences.captive_2_title}}</h2>
      <img :src="hotspot.preferences.captive_3_logo" class="ui centered image tiny">
    </div>
    <div v-show="!loading" class="ui route-container">
      <router-view class="ui"></router-view>
    </div>
  </div>
</template>

<script>
  import AuthMixin from './mixins/auth';
  export default {
    name: 'App',
    mixins: [AuthMixin],
    data: function () {
      var loading = true
      if (this.validateParams()) {
        this.getPreferences({
          digest: this.$root.$options.hotspot.digest,
          uuid: this.$root.$options.hotspot.uuid,
          sessionid: this.$root.$options.hotspot.sessionid,
        }, function (success) {
          this.$root.$options.hotspot.disclaimers = success.body.disclaimers
          this.$root.$options.hotspot.preferences = success.body.preferences
          this.$root.$options.hotspot.socials = success.body.socials
          this.$root.$options.hotspot.integrations = success.body.integrations
          this.$root.$options.session = {}
          this.hotspot.name = success.body.hotspot_name
          this.hotspot.preferences = success.body.preferences
          this.loading = false

          // background color
          $("body").css("background-color", success.body.preferences.captive_7_background || '#2a87be');

          // background image
          if (success.body.preferences.captive_81_bg_image) {
            $("body").css("height", "100vh");
            $("body").css("background-size", "cover");
            $("body").css("background-position", "center");
            $("body").css("background-image", 'url("' + success.body.preferences.captive_81_bg_image + '")');
          }

          this.titleColor = success.body.preferences.captive_83_title_color || '#4A4A4A'
          this.textFont = success.body.preferences.captive_85_text_style || 'Roboto';
          this.containerBgColor = success.body.preferences.captive_82_container_bg_color || '#ffffffdb';

          // wifi4eu head tags if enabled
          if (this.hotspot.preferences.wifi4eu_enabled === "true") {
            // define script with vars
            var script1 = document.createElement('script');
            script1.type = 'text/javascript';
            script1.textContent = 'var wifi4euTimerStart = Date.now();\nvar wifi4euNetworkIdentifier = \''
              + this.hotspot.preferences.wifi4eu_id +'\';\nvar wifi4euLanguage = \''
              + this.hotspot.preferences.wifi4eu_lang +'\';'
              + (this.hotspot.preferences.wifi4eu_zdebug === "true" ? '\nvar selftestModus = true;' : '');

            // define script of wifi4eu
            var script2 = document.createElement('script');
            script2.type = 'text/javascript';
            script2.src = 'https://collection.wifi4eu.ec.europa.eu/wifi4eu.min.js';

            // create img for banner
            var img = document.createElement('img')
            img.id = 'wifi4eubanner';
            img.width = window.innerWidth;
            img.height = img.width / 3.50;

            // append script to head
            document.head.appendChild(script1);
            document.head.appendChild(script2);

            // fire onload event to start wifi4eu script
            setTimeout(function() {
              document.body.prepend(img);
              document.body.style.height = "100%";
              dispatchEvent(new Event('load'));
              dispatchEvent(new Event('DOMContentLoaded'));
            }, 1000)
          }
        }, function(error) {
          console.error(error)
          this.loading = false
        })
      } else {
        if (!this.$route.query.state) {
          window.location.replace("http://1.0.0.1")
        }
        loading = false
      }

      return {
        hotspot: {
          name: '',
          preferences: {}
        },
        loading: loading,
        session: {
          loginDest: '',
          voucherCode: ''
        },
        titleColor: '#4A4A4A',
        textFont: 'Roboto',
        containerBgColor: '#ffffffdb'
      }
    },
    computed: {
      titleStyle: function () {
        return {
          color: this.titleColor,
          'font-family': this.textFont
        }
      },
      containerStyle: function () {
        return {
          'background-color': this.containerBgColor
        }
      }
    },
    methods: {
      validateParams: function () {
        var state = false
        if (Object.keys(this.$route.query).length > 0) {
          if (!this.$route.query.digest) {
            this.$root.$options.hotspot.onError = true
          } else if (!this.$route.query.uuid) {
            this.$root.$options.hotspot.onError = true
          } else if (!this.$route.query.sessionid) {
            this.$root.$options.hotspot.onError = true
          } else {
            this.$root.$options.hotspot.digest = this.$route.query.digest
            this.$root.$options.hotspot.uuid = this.$route.query.uuid
            this.$root.$options.hotspot.sessionid = this.$route.query.sessionid
            this.$root.$options.hotspot.uamip = this.$route.query.uamip
            this.$root.$options.hotspot.uamport = this.$route.query.uamport
            this.$root.$options.hotspot.nasid = this.$route.query.nasid
            state = true
          }
        } else {
          this.$root.$options.hotspot.onError = true
        }
        return state
      }
    }
  }
</script>

<style>
  html {
    height: initial !important;
  }

  body {
    background: #444;
  }

  img {
    margin-top: 20px !important;
    margin-bottom: 20px !important;
  }

  #app {
    -webkit-font-smoothing: antialiased !important;
    -moz-osx-font-smoothing: grayscale !important;
    text-align: center !important;
    color: #eff7fc !important;
    margin-top: 30px !important;
    margin-bottom: 30px !important;
  }

  .route-container {
    margin-top: 30px !important;
  }

  .ui.small.image {
    width: 450px !important;
  }

  .no-margin {
    margin: 0 !important;
  }

  .button-dark {
    background-color: #999999 !important;
  }

  .terms-text {
    margin: 0px;
    width: 100%;
    height: 200px;
    resize: vertical;
  }

  .ui.segment {
    border: none !important;
    box-shadow: none !important;
    -webkit-box-shadow: none !important;
  }

  @media (min-width: 723px) {
    .ui.container {
      width: 60% !important;
    }
  }

  @media (min-width: 933px) {
    .ui.container {
      width: 50% !important;
    }
  }

  @media (min-width: 1127px) {
    .ui.container {
      width: 40% !important;
    }
  }

  #wifi4eubanner {
    padding: 20px;
  }
</style>
