<template>
  <div class="ui segment">
    <div v-if="!hotspot.onError">
      <div v-if="!hotspot.loaded" class="ui active centered inline text loader">{{ $t("splash.retrieve_info") }}...
      </div>
      <div v-if="hotspot.loaded">
        <h3>{{hotspot.preferences.captive_4_subtitle}}</h3>
        <img class="ui centered small image" :src="hotspot.preferences.captive_5_banner" />
        <p v-html="hotspot.preferences.captive_6_description"></p>
        <div class="ui divider"></div>
        <div class="ui checkbox">
          <input v-model="hotspot.agree" type="checkbox" name="example" />
          <label>
            <span class="terms-space">{{$t('splash.i_agree')}}</span><a class="terms"
              @click="showModal()">{{$t('splash.agree_terms')}}</a>
          </label>
        </div>
        <div class="ui divider"></div>
        <router-link :to="hotspot.agree ? '/login' : '/'"
          :class="['big ui green button', !hotspot.agree ? 'disabled' : '']">{{ $t("splash.start_navigate") }}
        </router-link>
      </div>
    </div>
    <div v-if="hotspot.onError" class="ui icon negative message">
      <i class="warning icon"></i>
      <div class="content">
        <div class="header">{{ $t("splash.main_error") }}</div>
        <p>{{ $t("splash.main_error_sub") }}...</p>
      </div>
    </div>

    <v-dialog />
  </div>
</template>

<script>
  import AuthMixin from "./../mixins/auth";
  export default {
    name: "SplashPage",
    mixins: [AuthMixin],
    data: function () {
      var loaded = false;
      this.getPreferences({
          digest: this.$root.$options.hotspot.digest,
          uuid: this.$root.$options.hotspot.uuid,
          sessionid: this.$root.$options.hotspot.sessionid
        },
        function (success) {
          this.hotspot.preferences = success.body.preferences;
          this.hotspot.loaded = true;
        },
        function (error) {
          this.hotspot.loaded = false;
          console.error(error);
        }
      );
      return {
        hotspot: {
          onError: this.$root.$options.hotspot.onError,
          preferences: {},
          loaded: loaded,
          agree: false
        }
      };
    },
    methods: {
      showModal() {
        this.$modal.show("dialog", {
          title: this.$i18n.t("splash.modal_terms"),
          text: [
            "<div>" + this.$i18n.t("splash.tos") + "</div>",
            '<textarea readonly class="terms-text">' +
            this.$root.$options.hotspot.disclaimers.terms_of_use +
            "</textarea>",
            "<div></div>",
            "<div>" + this.$i18n.t("splash.privacy") + "</div>",
            '<textarea readonly class="terms-text">' +
            this.$root.$options.hotspot.disclaimers.marketing_use +
            "</textarea>"
          ].join(""),
          buttons: [{
            title: "Close"
          }]
        });
      }
    }
  };
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

  .terms {
    cursor: pointer;
  }

  .terms-space {
    margin-right: 3px;
  }
</style>