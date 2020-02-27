<template>
  <div class="ui form">
    <div v-if="!dedaloRequested">
      <div
        v-if="choosedMode && !this.$route.query.num"
        class="field"
        v-bind:class="{ error: errors.badInput }"
      >
        <label :style="textStyle">{{ $t("sms.prefix") }}</label>
        <div class="ui big left icon input">
          <select v-model="authPrefix">
            <option :value="c.dial_code" v-for="c in countries" v-bind:key="c.code">
              {{c.name}}
              ({{c.dial_code}})
            </option>
          </select>
        </div>
      </div>
      <div v-if="choosedMode" class="field">
        <label :style="textStyle">{{ $t("sms.number") }}</label>
        <div class="ui big left icon input">
          <input v-model="authSMS" type="tel" :placeholder="$t('sms.insert_number')" />
          <i class="talk icon"></i>
        </div>
      </div>
      <p v-if="!codeRequested && !choosedMode" class="not-code-exp">
        <label :style="textStyle">{{$t("sms.not_code_explain") }}</label>
        <br />
        <br />
        <button v-on:click="chooseMode()" class="ui button blue request-code" :style="buttonStyle">
          {{
          $t("sms.not_have_code") }}
        </button>
      </p>

      <p v-if="!codeRequested && !choosedMode" class="not-code-exp">
        <label :style="textStyle">{{$t("sms.not_code_explain_else") }}</label>
        <br />
        <br />
        <button
          v-if="!codeRequested && !choosedMode"
          v-on:click="chooseMode(true)"
          class="ui button blue request-code"
          :style="buttonStyle"
        >
          {{
          $t("sms.have_code") }}
        </button>
      </p>

      <button
        v-if="!codeRequested && choosedMode"
        v-on:click="getCode(true)"
        class="ui button blue request-code"
        :style="buttonStyle"
      >
        {{
        $t("sms.get_code") }}
      </button>
      <div v-if="errors.badNumber" class="ui tiny icon negative message">
        <i class="remove icon"></i>
        <div class="content">
          <div class="header">{{ $t("sms.error_code") }}</div>
          <p>{{ $t("sms.error_code_sub") }}</p>
        </div>
      </div>
      <div v-if="codeRequested">
        <div class="field">
          <label :style="textStyle">{{ $t("sms.code") }}</label>
          <div class="ui big left icon input">
            <input
              v-model="authCode"
              class="pw-input"
              :type="passwordVisible ? 'number' : 'password'"
              :placeholder="$t('sms.insert_your_code')"
            />
            <i class="braille icon"></i>
            <!-- toggle password visibility -->
            <button
              tabindex="-1"
              type="button"
              :class="['ui', 'button', passwordVisible ? '' : 'button-dark' ]"
              @click="togglePasswordVisibility()"
            >
              <i class="eye icon no-margin"></i>
            </button>
          </div>
        </div>
        <div v-if="bannerShow" class="ui compact message info no-margin-top">
          <div class="content">
            <div class="header">{{$t('sms.wait')}}</div>
            <p>{{$t('sms.we_are_sending_sms_code')}}</p>
          </div>
        </div>
      </div>
      <div class="ui divider"></div>
      <button
        :style="buttonStyle"
        v-on:click="back()"
        class="big ui red button"
      >{{ $t("login.back") }}</button>
      <button
        v-on:click="execLogin()"
        :disabled="isDisabled()"
        class="big ui green button"
        :style="buttonStyle"
      >{{ $t("sms.start_navigate") }}</button>
    </div>
    <div v-if="dedaloRequested">
      <div v-if="!authorized && !errors.dedaloError" class="ui active centered inline text loader">
        {{
        $t("sms.auth_progress") }}...
      </div>
      <div v-if="authorized" class="ui icon positive message">
        <i class="check icon"></i>
        <div class="content">
          <div class="header">{{ $t("sms.auth_success") }}</div>
        </div>
      </div>
      <div
        v-if="authorized && hotspot.preferences.marketing_0_reason_country == 'true' && userId != 0"
      >
        <h3>{{ $t("login.additional_info") }}</h3>
        <div class="field">
          <label :style="textStyle">{{ $t("login.country") }}</label>
          <div class="ui big left icon input">
            <select v-model="additionalCountry">
              <option :value="c.code" v-for="c in countries" v-bind:key="c.code">{{c.name}}</option>
            </select>
          </div>
        </div>
        <div class="field">
          <label :style="textStyle">{{ $t("login.reason") }}</label>
          <div class="ui big left icon input">
            <select v-model="additionalReason">
              <option value="business">{{$t("login.business")}}</option>
              <option value="family">{{$t("login.family")}}</option>
              <option value="other">{{$t("login.other")}}</option>
            </select>
          </div>
        </div>
      </div>
      <div v-if="errors.dedaloError" class="ui icon negative message">
        <i class="remove icon"></i>
        <div class="content">
          <div class="header">{{ $t("sms.auth_error") }}</div>
          <p>{{ $t("sms.auth_error_sub") }}</p>
        </div>
      </div>
      <div
        :class="hotspot.preferences.marketing_0_reason_country == 'true' ? 'adjust-top-big' : ''"
        v-if="authorized"
      >
        <div class="conditions-surveys">
          <div class="ui inline">
            <input id="conditions" v-model="conditions" type="checkbox" class="ui checkbox field" />
            <label :style="textStyle" for="conditions">{{ $t("login.disclaimer_privacy_accept") }}</label>
          </div>
          <div v-if="hotspot.preferences.marketing_1_enabled == 'true'" class="ui inline">
            <input id="surveys" v-model="surveys" type="checkbox" class="ui checkbox field" />
            <label :style="textStyle" for="surveys">{{ $t("login.disclaimer_survey_accept") }}</label>
          </div>
        </div>
        <div>
          <button
            :style="buttonStyle"
            v-on:click="navigate()"
            class="ui big button green"
          >{{ $t("login.navigate") }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import AuthMixin from "./../../mixins/auth";
import { setTimeout } from "timers";
export default {
  name: "SMSPage",
  mixins: [AuthMixin],
  data: function() {
    var params = this.extractParams();

    this.getPreferences(
      params,
      function(success) {
        this.$parent.hotspot.disclaimers = success.body.disclaimers;
        this.$root.$options.hotspot.disclaimers = success.body.disclaimers;
        this.$root.$options.hotspot.integrations = success.body.integrations;
        this.hotspot.disclaimers = success.body.disclaimers;
        this.hotspot.preferences = success.body.preferences;
        this.textColor =
          success.body.preferences.captive_84_text_color || "#4A4A4A";
        this.textFont =
          success.body.preferences.captive_85_text_style || "Roboto";

        if (this.$route.query.integration_done) {
          var context = this;
          context.dedaloRequested = true;
          context.authorized = false;
          context.errors.dedaloError = false;
          setTimeout(function() {
            context.execLogin();
          }, 1000);
        }
      },
      function(error) {
        console.error(error);
      }
    );

    var countries = require("./../../i18n/countries.json");
    var authPrefix = !this.$route.query.num ? "+39" : "";

    return {
      authorized: false,
      choosedMode:
        this.$route.query.num && this.$route.query.code ? true : false,
      codeRequested: this.$route.query.code || false,
      bannerShow: false,
      dedaloRequested: false,
      authPrefix: authPrefix,
      authSMS: this.$route.query.num || "",
      authCode: this.$route.query.code || "",
      authReset: this.$route.query.code || false,
      userId: this.$route.query.user || 0,
      resetDone: false,
      errors: {
        badNumber: false,
        badCode: false,
        dedaloError: false,
        badInput: false
      },
      countries: countries,
      hotspot: {
        disclaimers: this.$root.$options.hotspot.disclaimers,
        preferences: this.$root.$options.hotspot.preferences
      },
      conditions: false,
      surveys: false,
      countries: require("./../../i18n/countries.json"),
      additionalCountry: "-",
      additionalReason: "-",
      passwordVisible: true,
      textColor: "#4A4A4A",
      textFont: "Roboto"
    };
  },
  computed: {
    textStyle: function() {
      return {
        color: this.textColor,
        "font-family": this.textFont
      };
    },
    buttonStyle: function() {
      return {
        "font-family": this.textFont
      };
    }
  },
  methods: {
    back() {
      this.$router.push({
        path: "/login"
      });
    },
    isDisabled: function() {
      return this.authSMS.length == 0 || this.authCode.length == 0;
    },
    setPrefix: function(prefix) {
      this.authPrefix = prefix;
    },
    chooseMode: function(haveCode) {
      this.choosedMode = true;
      if (haveCode) {
        this.codeRequested = true;
      }
      setTimeout(function() {
        $(".ui.dropdown").dropdown();
      }, 100);
    },
    getCode: function(reset) {
      this.errors.badNumber = false;
      this.bannerShow = true;
      if (!(this.authPrefix + this.authSMS).startsWith("+")) {
        this.errors.badInput = true;
        return;
      }
      var params = this.extractParams();

      // make request to wax
      var url = this.createWaxURL(
        this.authPrefix + this.authSMS,
        params,
        "sms",
        reset
      );

      // get user id
      this.$http.get(url).then(
        function(responseAuth) {
          this.codeRequested = true;
          this.authReset = responseAuth.body.exists;
          this.resetDone = responseAuth.body.reset;
          this.userId = responseAuth.body.user_db_id;
        },
        function(error) {
          this.codeRequested = false;
          this.errors.badNumber = true;
          console.error(error);
        }
      );
    },
    execLogin: function() {
      this.dedaloRequested = true;
      this.authorized = false;
      this.errors.dedaloError = false;
      this.errors.badCode = false;

      if (
        this.$root.$options.hotspot.integrations &&
        this.$root.$options.hotspot.integrations[0] &&
        this.$root.$options.hotspot.integrations[0].post_auth_redirect_url &&
        this.$root.$options.hotspot.integrations[0].post_auth_redirect_url
          .length > 0 &&
        !this.$route.query.integration_done
      ) {
        // go to post_auth_redirect_url
        var redirectUrl = this.$root.$options.hotspot.integrations[0]
          .post_auth_redirect_url;

        var params = this.extractParams();

        var query =
          "?digest=" +
          params.digest +
          "&uuid=" +
          params.uuid +
          "&sessionid=" +
          params.sessionid +
          "&uamip=" +
          params.uamip +
          "&uamport=" +
          params.uamport +
          "&user=" +
          this.userId +
          "&nasid=" +
          params.nasid +
          "&code=" +
          this.authCode +
          "&num=" +
          encodeURIComponent(this.authPrefix + this.authSMS);

        var pathname = window.location.pathname;

        window.location.replace(redirectUrl + pathname + query);
      } else {
        // exec dedalo login
        this.doDedaloLogin(
          {
            id: this.authPrefix + this.authSMS,
            password: this.authCode || ""
          },
          function(responseDedalo) {
            if (responseDedalo.body.clientState == 1) {
              this.authorized = true;
              this.errors.dedaloError = false;
            } else {
              this.authorized = false;
              this.errors.dedaloError = true;
              this.errors.badCode = true;
            }
          },
          function(error) {
            this.authorized = false;
            this.errors.dedaloError = true;
            console.error(error);
          }
        );
      }
    },
    navigate() {
      if (this.conditions && this.surveys) {
        this.accept();
      } else {
        this.deleteInfo();
      }
    },
    deleteInfo: function() {
      var context = this;
      // extract code and state
      var params = this.extractParams();

      // delete marketing
      if (!this.conditions) {
        this.deleteMarketingInfo(
          this.userId,
          params,
          function(success) {
            context.accept();
          },
          function(error) {
            console.error(error);
            if (error.status == 404) {
              context.accept();
            }
          }
        );
      }

      if (!this.surveys) {
        this.deleteSurveyInfo(
          this.userId,
          params,
          function(success) {
            context.accept();
          },
          function(error) {
            console.error(error);
            if (error.status == 404) {
              context.accept();
            }
          }
        );
      }
    },
    accept: function() {
      var context = this;

      // extract code and state
      var params = this.extractParams();

      if (this.hotspot.preferences.marketing_0_reason_country == "true") {
        this.addAdditionalInfo(
          this.userId,
          params,
          {
            reason: this.additionalReason,
            country: this.additionalCountry
          },
          function(success) {
            // open redir url
            window.location.replace(
              context.$root.$options.hotspot.preferences.captive_1_redir
            );
          },
          function(error) {
            console.error(error);
            if (error.status == 404) {
              // open redir url
              window.location.replace(
                context.$root.$options.hotspot.preferences.captive_1_redir
              );
            }
          }
        );
      } else {
        // open redir url
        window.location.replace(
          context.$root.$options.hotspot.preferences.captive_1_redir
        );
      }
    },
    togglePasswordVisibility: function() {
      this.passwordVisible = !this.passwordVisible;
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

.auth-code-cont {
  margin-top: 15px !important;
  margin: 0;
}

.select-state {
  max-width: 270px;
  margin: 0 auto !important;
  margin-bottom: 1em !important;
  margin-top: 0.25em !important;
}

.request-code {
  margin-bottom: 10px !important;
  margin-top: 5px !important;
}

.no-margin-top {
  margin-top: 0px;
}

.not-code-exp {
  margin-top: 15px !important;
}

.adjust-top {
  margin-top: 10px !important;
}

.adjust-top-big {
  margin-top: 20px !important;
}

.adjust-checkbox {
  display: block !important;
}

.adjust-button {
  margin-top: 0px !important;
}

.conditions-surveys {
  display: inline-block !important;
  text-align: left !important;
}

.pw-input {
  width: 75% !important;
}
</style>