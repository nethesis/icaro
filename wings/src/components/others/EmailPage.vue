<template>
  <div class="ui form">
    <div v-if="!dedaloRequested">
      <div v-if="choosedMode" class="field" v-bind:class="{ error: errors.badInput }">
        <label :style="textStyle">Email</label>
        <div class="ui big left icon input">
          <input v-model="authEmail" type="email" :placeholder="$t('email.insert_email')" />
          <i class="mail icon"></i>
        </div>
      </div>
      <p v-if="!codeRequested && !choosedMode" class="not-code-exp">
        <label :style="textStyle">{{$t("email.not_code_explain") }}</label>
        <br />
        <br />
        <button v-on:click="chooseMode()" class="ui blue button request-code" :style="buttonStyle">
          {{
          $t("email.not_have_code") }}
        </button>
      </p>

      <p v-if="!codeRequested && !choosedMode" class="not-code-exp">
        <label :style="textStyle">{{$t("email.not_code_explain_else") }}</label>
        <br />
        <br />
        <button
          v-if="!codeRequested && !choosedMode"
          v-on:click="chooseMode(true)"
          class="ui blue button request-code"
          :style="buttonStyle"
        >
          {{
          $t("email.have_code") }}
        </button>
      </p>

      <button
        v-if="!codeRequested && choosedMode"
        v-on:click="getCode(true)"
        class="ui blue button request-code"
        :style="buttonStyle"
      >
        {{
        $t("email.get_code") }}
      </button>
      <div v-if="errors.badMail" class="ui tiny icon negative message">
        <i class="remove icon"></i>
        <div class="content">
          <div class="header">{{ $t("email.error_code") }}</div>
          <p>{{ $t("email.error_code_sub") }}</p>
        </div>
      </div>
      <div v-if="codeRequested">
        <div class="field">
          <label :style="textStyle">{{ $t("email.code") }}</label>
          <div class="ui big left icon input">
            <input
              v-model="authCode"
              class="pw-input"
              :type="passwordVisible ? 'number' : 'password'"
              :placeholder="$t('email.insert_your_code')"
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
            <div class="header">{{$t('email.wait')}}</div>
            <p>{{$t('email.we_are_sending_email_code')}}</p>
          </div>
        </div>
      </div>
      <div class="ui divider"></div>
      <button v-on:click="back()" class="big ui red button" :style="buttonStyle">{{ $t("login.back") }}</button>
      <button
        v-on:click="execLogin()"
        :disabled="isDisabled()"
        class="big ui green button"
        :style="buttonStyle"
      >{{ $t("email.start_navigate") }}</button>
    </div>
    <div v-if="dedaloRequested">
      <div v-if="!authorized && !errors.dedaloError" class="ui active centered inline text loader" :style="textStyle">
        {{
        $t("email.auth_progress") }}...
      </div>
      <div v-if="authorized" class="ui icon positive message">
        <i class="check icon"></i>
        <div class="content">
          <div class="header">{{ $t("email.auth_success") }}</div>
        </div>
      </div>
      <div
        v-if="authorized && hotspot.preferences.marketing_0_reason_country == 'true' && userId != 0"
      >
        <h3 :style="textStyle">{{ $t("login.additional_info") }}</h3>
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
          <div class="header">{{ $t("email.auth_error") }}</div>
          <p>{{ $t("email.auth_error_sub") }}</p>
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
          <button :style="buttonStyle" v-on:click="navigate()" class="ui big button green">{{ $t("login.navigate") }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import AuthMixin from "./../../mixins/auth";
import { setTimeout } from "timers";
export default {
  name: "EmailPage",
  mixins: [AuthMixin],
  mounted() {
    var params = this.extractParams();

    this.getPreferences(
      params,
      function(success) {
        this.$parent.hotspot.disclaimers = success.body.disclaimers;
        this.$root.$options.hotspot.disclaimers = success.body.disclaimers;
        this.$root.$options.hotspot.integrations = success.body.integrations;
        this.hotspot.disclaimers = success.body.disclaimers;
        this.hotspot.preferences = success.body.preferences;
        this.textColor = success.body.preferences.captive_84_text_color || '#4A4A4A';
        this.textFont = success.body.preferences.captive_85_text_style || 'Roboto';

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
  },
  data: function() {
    return {
      authorized: false,
      choosedMode:
        this.$route.query.email && this.$route.query.code ? true : false,
      codeRequested: this.$route.query.code || false,
      bannerShow: this.$route.query.code || false,
      dedaloRequested: false,
      authEmail: this.$route.query.email || "",
      authCode: this.$route.query.code || "",
      authReset: this.$route.query.code || false,
      userId: this.$route.query.user || 0,
      resetDone: false,
      errors: {
        badMail: false,
        badCode: false,
        dedaloError: false,
        badInput: false
      },
      hotspot: {
        disclaimers: this.$root.$options.hotspot.disclaimers,
        preferences: this.$root.$options.hotspot.preferences
      },
      conditions: false,
      surveys: false,
      countries: require("./../../i18n/countries.json"),
      additionalCountry: "-",
      additionalReason: "-",
      iOS: /iPad|iPhone|iPod/.test(navigator.userAgent) && !window.MSStream,
      passwordVisible: true,
      textColor: '#4A4A4A',
      textFont: 'Roboto',
    };
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
    back() {
      this.$router.push({
        path: "/login"
      });
    },
    isDisabled: function() {
      return this.authEmail.length == 0 || this.authCode.length == 0;
    },
    chooseMode: function(haveCode) {
      this.choosedMode = true;
      if (haveCode) {
        this.codeRequested = true;
      }
    },
    getCode: function(reset) {
      this.errors.badMail = false;
      this.bannerShow = true;
      if (this.authEmail.indexOf("@") == -1) {
        this.errors.badInput = true;
        return;
      }
      var params = this.extractParams();

      // make request to wax
      var url = this.createWaxURL(this.authEmail, params, "email", reset);

      // get user id
      this.$http.get(url).then(
        function(responseAuth) {
          this.authReset = responseAuth.body.exists;
          this.resetDone = responseAuth.body.reset;
          this.userId = responseAuth.body.user_db_id;

          // check if user already exists
          if (this.authReset && !(this.resetDone && this.resetDone == "true")) {
            this.codeRequested = true;
          } else {
            // open temp session for the user
            this.doTempSession(
              this.authEmail,
              null,
              function(responseTmp) {
                // if apple
                if (this.iOS) {
                  var origin = "http://conncheck." + window.location.host;
                  var pathname = window.location.pathname;
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
                    "&code=.&email=" +
                    this.authEmail;
                  window.location.replace(origin + pathname + query);
                } else {
                  this.codeRequested = true;
                }
              },
              function(error) {
                this.codeRequested = false;
                this.errors.badMail = true;
                console.error(error);
              }
            );
          }
        },
        function(error) {
          this.codeRequested = false;
          this.errors.badMail = true;
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
          "&email=" +
          this.authEmail;

        var pathname = window.location.pathname;

        this.doDedaloLogout(
          function(responseDedaloLogout) {
            window.location.replace(redirectUrl + pathname + query);
          },
          function(error) {
            this.authorized = false;
            this.errors.dedaloError = true;
            console.error(error);
          }
        );
      } else {
        // exec logout
        this.doDedaloLogout(
          function(responseDedaloLogout) {
            // exec dedalo login
            this.doDedaloLogin(
              {
                id: this.authEmail,
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

.request-code {
  margin-bottom: 10px !important;
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