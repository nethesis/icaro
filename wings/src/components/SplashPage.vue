<template>
  <div class="ui">
    <div v-if="!hotspot.onError">
      <div v-if="!hotspot.loaded" class="ui active centered inline text loader">{{ $t("splash.retrieve_info") }}...
      </div>
      <div v-if="hotspot.loaded">
        <h3 :style="textStyle">{{hotspot.preferences.captive_4_subtitle}}</h3>
        <img class="ui centered small image" :src="hotspot.preferences.captive_5_banner" />
        <p :style="textStyle" v-html="hotspot.preferences.captive_6_description"></p>
        <div class="ui divider"></div>
        <div class="ui checkbox">
          <input v-model="hotspot.agree" type="checkbox" id="agree-terms" />
          <label for="agree-terms" class="inline">
            <span class="terms-space" :style="textStyle">{{$t('splash.i_agree')}}</span>
          </label><span class="terms-space" :style="textStyle">(<a class="terms" @click="showModal()" :style="textStyle">{{$t('splash.tos_and_privacy')}}</a>)</span>
        </div>
        <div class="ui divider"></div>
        <router-link :to="hotspot.agree ? '/login' : '/'"
          :class="['big ui green button', !hotspot.agree ? 'disabled' : '']" :style="buttonStyle">{{ $t("splash.start_navigate") }}
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

    <v-dialog class="terms-dialog"/>
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
          this.textColor = success.body.preferences.captive_84_text_color || '#4A4A4A';
          this.textFont = success.body.preferences.captive_85_text_style || 'Roboto';
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
          agree: false,
          textColor: '#4A4A4A',
          textFont: 'Roboto',
        }
      }
    },
    computed: {
      textStyle: function () {
        return {
          color: this.textColor,
          'font-family': this.textFont
        }
      },
      buttonStyle: function () {
        return {
          'font-family': this.textFont
        }
      }
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

  a.terms {
    text-decoration: underline !important;
    cursor: pointer;
  }

  .terms-space {
    margin-right: 3px;
    cursor: default;
  }

  .terms-dialog {
    color: #4A4A4A;
  }

  label.inline {
    display: inline;
  }
</style>