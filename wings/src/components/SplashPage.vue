<template>
  <div class="ui segment">
    <div v-if="!hotspot.onError">
      <div v-if="!hotspot.loaded" class="ui active centered inline text loader">{{ $t("splash.retrieve_info") }}...</div>
      <div v-if="hotspot.loaded">
        <h3>{{hotspot.preferences.captive_subtitle}}</h3>
        <img class="ui centered small image" :src="hotspot.preferences.captive_banner">
        <p>
          {{hotspot.preferences.captive_description}}
        </p>
        <div class="ui divider"></div>
        <router-link to='/login' class="big ui green button">
          {{ $t("splash.start_navigate") }}
        </router-link>
      </div>
    </div>
    <div v-if="hotspot.onError" class="ui icon negative message">
      <i class="warning icon"></i>
      <div class="content">
        <div class="header">
          {{ $t("splash.main_error") }}
        </div>
        <p>{{ $t("splash.main_error_sub") }}...</p>
      </div>
    </div>
  </div>

</template>

<script>
  import AuthMixin from './../mixins/auth';
  export default {
    name: 'SplashPage',
    mixins: [AuthMixin],
    data() {
      var loaded = false
      this.getPreferences({
        digest: this.$root.$options.hotspot.digest,
        uuid: this.$root.$options.hotspot.uuid,
        sessionid: this.$root.$options.hotspot.sessionid,
      }, success => {
        this.hotspot.preferences = success.body.preferences
        this.hotspot.loaded = true
      }, error => {
        this.hotspot.loaded = false
        console.error(error)
      })
      return {
        hotspot: {
          onError: this.$root.$options.hotspot.onError,
          preferences: {},
          loaded: loaded
        }
      }
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  h1,
  h2 {
    font-weight: normal;
  }

  ul {
    list-style-type: none;
    padding: 0;
  }

  li {
    display: inline-block;
    margin: 0 10px;
  }

  a {
    color: #42b983;
  }
</style>