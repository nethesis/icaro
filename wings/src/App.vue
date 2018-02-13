<template>
  <div id="app" class="ui container">
    <div class="ui center aligned">
      <h2>{{hotspot.preferences.captive_2_title}}</h2>
      <img :src="hotspot.preferences.captive_3_logo" class="ui centered image tiny">
    </div>
    <div class="ui segments route-container">
      <router-view class="ui segment"></router-view>
    </div>
  </div>
</template>

<script>
  import AuthMixin from './mixins/auth';
  export default {
    name: 'App',
    mixins: [AuthMixin],
    data() {
      this.validateParams()
      this.getPreferences({
        digest: this.$root.$options.hotspot.digest,
        uuid: this.$root.$options.hotspot.uuid,
        sessionid: this.$root.$options.hotspot.sessionid,
      }, success => {
        this.$root.$options.hotspot.preferences = success.body.preferences
        this.$root.$options.hotspot.socials = success.body.socials
        this.hotspot.name = success.body.hotspot_name
        this.hotspot.preferences = success.body.preferences
        $("body").css("background-color", success.body.preferences.captive_7_background || '#2a87be');
      }, error => {
        console.error(error)
        $("body").css("background-color", '#2a87be');
      })
      return {
        hotspot: {
          name: '',
          preferences: {}
        }
      }
    },
    methods: {
      validateParams() {
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
          }
        } else {
          this.$root.$options.hotspot.onError = true
        }
      }
    }
  }
</script>

<style>
  html {
    height: initial !important;
  }

  body {
    background: #2a87be;
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

  .ui.segments {
    color: #5a5a5a !important;
  }
</style>