<template>
  <div class="ui form">
    <div v-if="voucherAvailable && !voucherValidated && !authorized && !dedaloError">
      <h3 :style="textStyle">{{ $t("login.voucher_title") }}</h3>
      <div class="inline field" v-bind:class="{ error: badInput }">
        <label :style="textStyle">Voucher</label>
        <div class="ui big left icon input">
          <input
            v-model="authCode"
            class="pw-input"
            :type="voucherVisible ? 'email' : 'password'"
            :placeholder="$t('login.insert_voucher')"
          >
          <i class="braille icon"></i>
          <!-- toggle voucher visibility -->
          <button
            tabindex="-1"
            type="button"
            :class="['ui', 'button', voucherVisible ? '' : 'button-dark' ]"
            @click="toggleVoucherVisibility()"
            :style="buttonStyle"
          >
            <i class="eye icon no-margin"></i>
          </button>
        </div>
      </div>
      <button :style="buttonStyle" v-on:click="validateCode()" class="ui big button">{{ $t("login.validate_code") }}</button>
      <div v-if="badCode" class="ui tiny icon negative message">
        <i class="remove icon"></i>
        <div class="content">
          <div class="header">{{ $t("login.error_voucher_code") }}</div>
          <p>{{ $t("login.error_voucher_code_sub") }}</p>
        </div>
      </div>
    </div>
    <div
      v-if="!voucherAvailable || (voucherAvailable && voucherValidated) && !authorized && !dedaloError"
    >
      <h3 :style="textStyle">{{ $t("login.choose_login") }}</h3>
      <div class="ui relaxed list">
        <div v-if="hotspot.preferences.facebook_login == 'true'" class="item">
          <div @click="changeRoute('/login/facebook', false)" class="ui facebook button big fluid" :style="buttonStyle">
            <i class="facebook icon"></i>
            Facebook
          </div>
        </div>
        <div v-if="hotspot.preferences.instagram_login == 'true'" class="item">
          <div
            @click="changeRoute('/login/instagram', false)"
            class="ui instagram button big fluid"
            :style="buttonStyle"
          >
            <i class="instagram icon"></i>
            Instagram
          </div>
        </div>
        <div v-if="hotspot.preferences.linkedin_login == 'true'" class="item">
          <div @click="changeRoute('/login/linkedin', false)" class="ui linkedin button big fluid" :style="buttonStyle">
            <i class="linkedin icon"></i>
            LinkedIn
          </div>
        </div>
      </div>
      <div class="ui divider"></div>
      <div class="ui relaxed list">
        <div v-if="hotspot.preferences.sms_login == 'true'" class="item">
          <div @click="changeRoute('/login/sms', false)" class="ui button yellow big fluid" :style="buttonStyle">
            <i class="talk icon"></i>
            SMS
          </div>
        </div>
        <div v-if="hotspot.preferences.email_login == 'true'" class="item">
          <div @click="changeRoute('/login/email', false)" class="ui button red big fluid" :style="buttonStyle">
            <i class="mail icon"></i>
            Email
          </div>
        </div>
      </div>
      <div v-if="hotspot.preferences.temp_code_login == 'true'" class="ui divider"></div>
      <div class="ui relaxed list">
        <div v-if="hotspot.preferences.temp_code_login == 'true'" class="item">
          <div @click="changeRoute('/login', true)" class="ui button teal big fluid" :style="buttonStyle">
            <i class="barcode icon"></i>
            {{ $t("login.code") }}
          </div>
        </div>
      </div>
    </div>
    <div v-if="dedaloError" class="ui icon negative message">
      <div class="content">
        <div class="header">{{ $t("social.auth_error") }}</div>
        <p>{{ $t("social.auth_error_sub") }}</p>
      </div>
    </div>
    <button
      v-if="dedaloError"
      v-on:click="back()"
      class="ui big button green"
      :style="buttonStyle"
    >{{ $t("login.back") }}</button>

    <div v-if="authorized" class="ui icon positive message">
      <i class="check icon"></i>
      <div class="content">
        <div class="header">{{ $t("social.auth_success") }}</div>
      </div>
    </div>

    <div v-if="authorized && hotspot.preferences.marketing_0_reason_country == 'true' && hotspot.preferences.check_marketing == 'false'">
      <h3 :style="textStyle">{{ $t("login.additional_info") }}</h3>
      <div class="inline field">
        <label :style="textStyle">{{ $t("login.country") }}</label>
        <div class="ui big left icon input">
          <select v-model="additionalCountry">
            <option value="-">--</option>
            <option :value="c.code" v-for="c in countries" v-bind:key="c.code">{{c.name}}</option>
          </select>
        </div>
      </div>
      <div class="inline field">
        <label :style="textStyle">{{ $t("login.reason") }}</label>
        <div class="ui big left icon input">
          <select v-model="additionalReason">
            <option value="-">--</option>
            <option value="business">{{$t("login.business")}}</option>
            <option value="family">{{$t("login.family")}}</option>
            <option value="other">{{$t("login.other")}}</option>
          </select>
        </div>
      </div>
    </div>

    <div
      :class="hotspot.preferences.marketing_0_reason_country == 'true' ? 'adjust-top-big' : ''"
      v-if="authorized"
    >
      <div class="conditions-surveys">
        <div class="ui inline" v-if="hotspot.preferences.check_marketing == 'false'">
          <input id="conditions" v-model="conditions" type="checkbox" class="ui checkbox field">
          <label :style="textStyle" for="conditions">{{ $t("login.disclaimer_privacy_accept") }}</label>
        </div>
        <div v-if="hotspot.preferences.marketing_1_enabled == 'true' && hotspot.preferences.check_marketing == 'false'" class="ui inline">
          <input id="surveys" v-model="surveys" type="checkbox" class="ui checkbox field">
          <label :style="textStyle" for="surveys">{{ $t("login.disclaimer_survey_accept") }}</label>
        </div>
      </div>
      <div>
        <button :style="buttonStyle" v-on:click="navigate()" class="ui big button green">{{ $t("login.navigate") }}</button>
      </div>
    </div>
  </div>
</template>

<script>
import AuthMixin from "./../mixins/auth";
export default {
  name: "LoginPage",
  mixins: [AuthMixin],
  data: function() {
    var voucherAvailable = false;
    var badCode = false;
    var badInput = false;
    var authorized = false;
    var dedaloError = false;

    var authCode = "";
    if (this.$root.$options.hotspot.preferences.voucher_login == "true") {
      voucherAvailable = true;
    }

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

        if (this.$route.query.integration_done && this.$route.query.code) {
          this.voucherAvailable = true;
          this.voucherValidated = false;
          this.authorized = false;
          this.dedaloError = false;
          this.authCode = this.$route.query.code;
          this.validateCode();
        }
      },
      function(error) {
        console.error(error);
      }
    );

    return {
      hotspot: {
        preferences: this.$root.$options.hotspot.preferences,
        disclaimers: this.$root.$options.hotspot.disclaimers
      },
      voucherAvailable: voucherAvailable,
      voucherValidated: this.$root.$options.hotspot.voucherValidated,
      authCode: authCode,
      badCode: badCode,
      badInput: badInput,
      authorized: authorized,
      dedaloError: dedaloError,
      userId: 0,
      conditions: false,
      surveys: false,
      countries: require("./../i18n/countries.json"),
      additionalCountry: "-",
      additionalReason: "-",
      voucherVisible: true,
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
      this.voucherAvailable = this.$root.$options.hotspot.preferences.voucher_login;
      this.voucherValidated = false;
      this.$root.$options.hotspot.voucherValidated = false;
      this.authCode = "";
      this.badCode = false;
      this.badInput = false;
      this.authorized = false;
      this.dedaloError = false;
    },
    changeRoute: function(path, withCode) {
      if (withCode) {
        this.voucherAvailable = true;
        this.voucherValidated = false;
        this.authCode = "";
        this.$root.$options.hotspot.voucherValidated = false;
        this.badCode = false;
      } else {
        this.$router.push({
          path: path
        });
      }
    },
    validateCode: function() {
      this.badCode = false;
      if (this.authCode.length == 0) {
        this.badInput = true;
        return;
      }
      var params = this.extractParams();

      // make request to wax
      var url = this.createWaxURL(this.authCode, params, "voucher");

      // get user id
      this.$http.get(url).then(
        function(responseAuth) {
          this.$root.$options.session["voucherCode"] = responseAuth.body.code;
          this.userId = responseAuth.body.user_db_id;

          if (
            responseAuth.body.type == "auth" &&
            this.hotspot.preferences.temp_code_login == "true"
          ) {
            // do login immediately
            var context = this;

            // check integrations
            if (
              this.$root.$options.hotspot.integrations &&
              this.$root.$options.hotspot.integrations[0] &&
              this.$root.$options.hotspot.integrations[0]
                .post_auth_redirect_url &&
              this.$root.$options.hotspot.integrations[0].post_auth_redirect_url
                .length > 0 &&
              !this.$route.query.integration_done
            ) {
              // go to post_auth_redirect_url
              var redirectUrl = this.$root.$options.hotspot.integrations[0]
                .post_auth_redirect_url;

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
                "&validated=true";

              var pathname = window.location.pathname;

              window.location.replace(redirectUrl + pathname + query);
            } else {
              this.doDedaloLogin(
                {
                  id: responseAuth.body.code,
                  password: responseAuth.body.code || ""
                },
                function(responseDedalo) {
                  if (
                    responseDedalo &&
                    responseDedalo.body &&
                    responseDedalo.body.clientState == 1
                  ) {
                    context.authorized = true;
                    context.dedaloError = false;
                    context.$forceUpdate();
                  } else {
                    context.authorized = false;
                    context.dedaloError = true;
                  }
                },
                function(error) {
                  this.authorized = false;
                  this.dedaloError = true;
                  console.error(error);
                }
              );
            }
          } else if (
            responseAuth.body.type == "auth" &&
            this.hotspot.preferences.temp_code_login == "false"
          ) {
            this.voucherValidated = false;
            this.$root.$options.hotspot.voucherValidated = false;
            this.badCode = true;
          } else {
            this.voucherValidated = true;
            this.$root.$options.hotspot.voucherValidated = true;
            this.badCode = false;
          }
        },
        function(error) {
          this.voucherValidated = false;
          this.$root.$options.hotspot.voucherValidated = false;
          this.badCode = true;
          console.error(error);
        }
      );
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
    toggleVoucherVisibility: function() {
      this.voucherVisible = !this.voucherVisible;
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

.divider {
  margin: 20px;
}

.item {
  margin: 10px;
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
