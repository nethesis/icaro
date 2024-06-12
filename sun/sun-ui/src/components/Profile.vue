<template>
  <div>
    <h2>{{ msg }}</h2>

    <div class="row row-cards-pf">
      <div class="col-xs-12 col-sm-12 col-md-6">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{ $t("profile.info") }}
              <div :class="[getLoginIcon(user.info.type), 'right']" data-toggle="tooltip" data-placement="left" :title="$t(user.info.type)"></div>
            </h2>
          </div>
          <div class="card-pf-body">
            <div class="list-details">
              <dt>{{ $t("profile.username") }}</dt>
              <dd>{{user.info.username}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("profile.name") }}</dt>
              <dd>{{user.info.name}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("profile.email") }}</dt>
              <dd>{{user.info.email}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("profile.type") }}</dt>
              <dd>{{$t(user.info.type)}}</dd>
            </div>
          </div>
          <div class="card-pf-footer">
            <div class="dropdown card-pf-time-frame-filter">
              <button class="btn btn-default" data-toggle="modal" data-target="#changePassModal">{{ $t("profile.change_password") }}</button>
            </div>
            <p>
              <a href="#" class="card-pf-link-with-icon">
              </a>
            </p>
          </div>
        </div>
      </div>

      <!-- privacy fields -->
      <div v-if="user.info.type == 'reseller'" class="col-xs-12 col-sm-12 col-md-6">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{ $t("profile.privacy_fields") }}
              <div v-if="!privacy.isLoading" class="fa fa-lock card-info-title right"></div>
              <div v-if="privacy.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
          <form v-if="!privacy.isLoading" class="form-horizontal" role="form" v-on:submit.prevent="updatePrivacyFields()">
            <div class="alert alert-info alert-dismissable">
              <span class="pficon pficon-info"></span>
              <strong>{{ $t('profile.info')}}.</strong> {{ $t('profile.privacy_info') }}.
            </div>
            <div class="card-pf-body">
                <div class="form-group">
                  <label class="col-sm-4 control-label">{{ $t("profile.privacy_name") }}</label>
                  <div class="col-sm-8">
                    <input v-model="user.info.privacy_name" required type="text" class="form-control">
                  </div>
                </div>
                <div class="form-group">
                  <label class="col-sm-4 control-label">{{ $t("profile.privacy_vat") }}</label>
                  <div class="col-sm-8">
                    <input v-model="user.info.privacy_vat" required type="number" class="form-control">
                  </div>
                </div>
                <div class="form-group">
                  <label class="col-sm-4 control-label">{{ $t("profile.privacy_address") }}</label>
                  <div class="col-sm-8">
                    <input v-model="user.info.privacy_address" required type="text" class="form-control">
                  </div>
                </div>
                <div class="form-group">
                  <label class="col-sm-4 control-label">{{ $t("profile.privacy_email") }}</label>
                  <div class="col-sm-8">
                    <input v-model="user.info.privacy_email" required type="email" class="form-control">
                  </div>
                </div>
                <div class="form-group">
                  <label class="col-sm-4 control-label">{{ $t("profile.privacy_dpo") }}</label>
                  <div class="col-sm-8">
                    <input v-model="user.info.privacy_dpo" type="text" class="form-control">
                  </div>
                </div>
                <div class="form-group">
                  <label class="col-sm-4 control-label">{{ $t("profile.privacy_dpo_mail") }}</label>
                  <div class="col-sm-8">
                    <input v-model="user.info.privacy_dpo_mail" type="email" class="form-control">
                  </div>
                </div>
            </div>
            <div class="card-pf-footer">
              <div class="dropdown card-pf-time-frame-filter">
                <button type="submit" class="btn btn-default">{{ $t("update") }}</button>
              </div>
              <p>
                <a href="#" class="card-pf-link-with-icon">
                </a>
              </p>
            </div>
          </form>
        </div>
      </div>

      <!-- SMS warning threshold -->
      <div class="col-xs-12 col-sm-12 col-md-6">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{ $t("profile.sms") }}
              <div v-if="!sms.isLoading" class="fa fa-commenting card-info-title right"></div>
              <div v-if="sms.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
          <div v-if="!sms.isLoading" class="card-pf-body">
            <form class="form-horizontal" role="form">
              <div class="form-group">
                <label class="col-sm-4 control-label">{{ $t("profile.sms_login_threshold") }}</label>
                <div class="col-sm-8">
                  <input v-model="sms.data.sms_threshold" required type="number" class="form-control">
                </div>
              </div>
            </form>
            <p>{{ $t("profile.sms_threshold_description") }}.</p>
          </div>
          <div class="card-pf-footer">
            <div class="dropdown card-pf-time-frame-filter">
              <button class="btn btn-default" v-on:click="updateSmsThreshold()">{{ $t("update") }}</button>
            </div>
            <p>
              <a href="#" class="card-pf-link-with-icon">
              </a>
            </p>
          </div>
        </div>
      </div>

      <!-- disclaimer -->
      <div
        v-if="disclaimers.data && disclaimers.data.length"
        class="col-xs-12 col-sm-12 col-md-6"
      >
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{ $t("profile.default_disclaimer") }}
              <div v-if="!disclaimers.isLoading" class="pficon pficon-catalog card-info-title right"></div>
              <div v-if="disclaimers.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
          <div v-if="!disclaimers.isLoading" class="card-pf-body">
            <form class="form-horizontal" role="form">
              <div class="form-group">
                <label class="col-sm-4 control-label">{{ $t("profile.privacy_disclaimer") }}</label>
                <div class="col-sm-8">
                  <select v-model="disclaimers.preferredPrivacyDisclaimerId" class="form-control">
                    <option value="">{{ $t("profile.default") }}</option>
                    <option v-for="privacyDisclaimer in disclaimers.privacyDisclaimers" :key="privacyDisclaimer.id" :value="privacyDisclaimer.id">
                      {{ privacyDisclaimer.title }}
                    </option>
                  </select>
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label">{{ $t("profile.tos_disclaimer") }}</label>
                <div class="col-sm-8">
                  <select v-model="disclaimers.preferredTosDisclaimerId" class="form-control">
                    <option value="">{{ $t("profile.default") }}</option>
                    <option v-for="tosDisclaimer in disclaimers.tosDisclaimers" :key="tosDisclaimer.id" :value="tosDisclaimer.id">
                      {{ tosDisclaimer.title }}
                    </option>
                  </select>
                </div>
              </div>
            </form>
          </div>
          <div class="card-pf-footer">
            <div class="dropdown card-pf-time-frame-filter">
              <span v-if="disclaimers.updateError" class="pficon pficon-error-circle-o bigger update-pref-feedback"></span>
              <span v-if="disclaimers.updateSuccess" class="pficon pficon-ok update-pref-feedback"></span>
              <button class="btn btn-default" v-on:click="updateDisclaimers()">{{ $t("update") }}</button>
            </div>
            <p>
              <a href="#" class="card-pf-link-with-icon">
              </a>
            </p>
          </div>
        </div>
      </div>
    </div>

    <div class="modal fade" id="changePassModal" tabindex="-1" role="dialog" aria-labelledby="changePassModalLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="changePassModalLabel">{{ $t("profile.change_password") }}</h4>
          </div>
          <form class="form-horizontal" role="form" v-on:submit.prevent="changePassword()">
            <div class="modal-body">
              <div class="form-group">
                <label class="col-sm-4 control-label" for="textInput-modal-markup">{{ $t("profile.new_password") }}</label>
                <div class="col-sm-8">
                  <input v-model="newPassword" required type="text" id="textInput-modal-markup" class="form-control" :placeholder="$t('profile.new_password')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="textInput2-modal-markup">{{ $t("profile.confirm_password") }}</label>
                <div class="col-sm-8">
                  <input v-model="confirmPassword" required type="text" id="textInput2-modal-markup" class="form-control" :placeholder="$t('profile.confirm_password')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="textInput2-modal-markup"></label>
                <div class="col-sm-8">
                  <span v-bind:class="[newPassword == confirmPassword && newPassword.length > 0 ? 'pass-confirm-ok' : 'pass-confirm-err', '']"></span>
                </div>
              </div>
              <div v-if="errors.password" class="alert alert-danger alert-dismissable">
                <span class="pficon pficon-error-circle-o"></span>
                <strong>{{ $t("profile.change_password_error_title") }}</strong>. {{ $t("profile.change_password_error_sub") }}.
              </div>
            </div>
            <div class="modal-footer">
              <span v-if="onAction" class="spinner spinner-sm spinner-inline modal-spinner"></span>
              <button type="button" class="btn btn-default" data-dismiss="modal">{{ $t("cancel") }}</button>
              <button :disabled="(newPassword != confirmPassword) || newPassword.length == 0" type="submit" class="btn btn-primary">{{ $t("update") }}</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import LoginService from "../services/login";
import StorageService from "../services/storage";
import UtilService from "../services/util";
import StatsService from "../services/stats";
import DisclaimerService from "../services/disclaimer";
import PreferenceService from "../services/preference";
import AccountService from '../services/account';

export default {
  name: "Profile",
  mixins: [LoginService, StorageService, UtilService, StatsService, DisclaimerService, PreferenceService, AccountService],
  data() {
    return {
      msg: this.$i18n.t("menu.profile"),
      user: {
        login: this.get("loggedUser") || null,
        info: this.$parent.user.info
      },
      newPassword: "",
      confirmPassword: "",
      errors: {
        password: false
      },
      onAction: false,
      sms: {
        isLoading: true,
        data: {}
      },
      privacy: {
        isLoading: false,
      },
      disclaimers: {
        data: {},
        isLoading: true,
        preferredPrivacyDisclaimerId: "",
        preferredTosDisclaimerId: "",
        privacyDisclaimers: [],
        tosDisclaimers: [],
        updateSuccess: false,
        updateError: false,
      },
    };
  },
  mounted() {
    this.getSmsThreshold();
    this.getDisclaimers();
  },
  methods: {
    changePassword() {
      this.onAction = true;
      this.execChangePassword(
        this.newPassword,
        this.user.login.id,
        success => {
          this.onAction = false;
          $("#changePassModal").modal("toggle");
        },
        error => {
          this.onAction = false;
          this.errors.password = true;
          console.error(error.body.message);
        }
      );
    },
    updateSmsThreshold() {
      this.sms.isLoading = true;

      this.updateSMSThresholdForAccountByAccount(
        {
          sms_threshold: parseInt(this.sms.data.sms_threshold) || 0
        },
        this.user.login.id,
        success => {
          this.sms.isLoading = false;
          this.sms.data = success.body;
          this.getSmsThreshold();
        },
        error => {
          this.sms.isLoading = false;
          console.error(error.body);
        }
      );
    },
    getSmsThreshold() {
      this.statsSMSTotalForAccountByAccount(
        this.user.login.id,
        success => {
          this.sms.isLoading = false;
          this.sms.data = success.body;
        },
        error => {
          this.sms.isLoading = false;
          console.error(error.body);
        }
      );
    },
    getDisclaimers() {
      this.disclaimersByAccount(
        this.user.login.id,
        success => {
          this.disclaimers.data = success.body;
          this.privacy_disclaimers = [];
          this.tos_disclaimers = [];

          this.disclaimers.data.forEach((disclaimer) => {
            if (disclaimer.type === "privacy") {
              this.disclaimers.privacyDisclaimers.push(disclaimer);
            } else {
              this.disclaimers.tosDisclaimers.push(disclaimer);
            }
          });
          this.getPreferredDisclaimers();
        },
        error => {
          this.disclaimers.isLoading = false;
          console.error(error.body);
        }
      );
    },
    getPreferredDisclaimers() {
      this.accountPrefGet(
        this.user.login.id,
        success => {
          this.disclaimers.preferredPrivacyDisclaimerId = "";
          this.disclaimers.preferredTosDisclaimerId = "";

          success.body.forEach((pref) => {
            if (pref.key === "custom_disclaimers_privacy") {
              this.disclaimers.preferredPrivacyDisclaimerId = pref.value;
            } else if (pref.key === "custom_disclaimers_terms") {
              this.disclaimers.preferredTosDisclaimerId = pref.value;
            }
          });
          this.disclaimers.isLoading = false;
        },
        error => {
          this.disclaimers.isLoading = false;
          console.error(error.body);
        }
      );
    },
    updateDisclaimers() {
      this.updateDisclaimer("custom_disclaimers_privacy", this.disclaimers.preferredPrivacyDisclaimerId);
      this.updateDisclaimer("custom_disclaimers_terms", this.disclaimers.preferredTosDisclaimerId);
    },
    updateDisclaimer(prefKey, prefValue) {
      this.accountPrefModify(
        this.user.login.id,
        {
          key: prefKey,
          value: prefValue.toString(),
          account_id: this.user.login.id,
        },
        success => {
          this.disclaimers.updateSuccess = true;
          setTimeout(() => {
            this.disclaimers.updateSuccess = false;
          }, 2000);
        },
        error => {
          console.error(error.body);

          this.disclaimers.updateError = true;
          setTimeout(() => {
            this.disclaimers.updateError = false;
          }, 2000);
        }
      );
    },
    updatePrivacyFields() {
      this.privacy.isLoading = true;
      this.accountModify(this.user.login.id, {
        privacy_name: this.user.info.privacy_name,
        privacy_vat:  this.user.info.privacy_vat,
        privacy_address:  this.user.info.privacy_address,
        privacy_email:  this.user.info.privacy_email,
        privacy_dpo:  this.user.info.privacy_dpo,
        privacy_dpo_mail:  this.user.info.privacy_dpo_mail,
      },
      success => {
        this.privacy.isLoading = false;
      },
      error => {
        this.privacy.isLoading = false;
        console.error(error.body.message);
      });
    }
  }
};
</script>

<style>
.update-pref-feedback {
  vertical-align: middle;
  margin-right: 0.4em;
  font-size: 140%;
}
</style>
