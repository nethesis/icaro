<template>
  <div id="app" class="ui container">
    <div v-show="loading" class="ui active dimmer">
      <div class="ui loader"></div>
    </div>
    <div v-show="!loading" class="ui center aligned">
      <h2>{{hotspot.preferences.captive_2_title}}</h2>
      <img :src="hotspot.preferences.captive_3_logo" class="ui centered image tiny">
    </div>
    <div v-show="!loading" class="ui segments route-container">
      <router-view class="ui segment"></router-view>
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

          setTimeout(function() {
            // container color and opacity
            $("div.ui.segment").css("background-color", success.body.preferences.captive_82_container_bg_color || '#fff9');

            // title color
            $("h2").css("color", success.body.preferences.captive_83_title_color || '#fff');

            // text color
            $("div.ui.segment h3").css("color", success.body.preferences.captive_84_text_color || '#f00');
            $("div.ui.segment p").css("color", success.body.preferences.captive_84_text_color || '#383838');

            // text syle
            $("h2").css("font-family", success.body.preferences.captive_85_text_style || 'Lato');
            $("div.ui.segment h3").css("font-family", success.body.preferences.captive_85_text_style || 'Lato');
            $("div.ui.segment p").css("font-family", success.body.preferences.captive_85_text_style || 'Lato');
            $("div.ui.segment a.green.button").css("font-family", success.body.preferences.captive_85_text_style || 'Lato');
          }, 300);
        }, function(error) {
          console.error(error)
          $("body").css("background-color", '#fff');
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
    height: 86px;
    resize: vertical;
  }
  
  .ui.segments {
    border: none !important;
    box-shadow: none !important;
    -webkit-box-shadow: none !important;
  }
</style>