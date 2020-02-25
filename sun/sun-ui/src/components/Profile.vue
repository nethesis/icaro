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
      <!-- Whatsapp warning threshold -->
      <div class="col-xs-12 col-sm-12 col-md-6">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{ $t("profile.whatsapp") }}
              <div v-if="!whatsapp.isLoading" class="fa fa-whatsapp card-info-title right"></div>
              <div v-if="whatsapp.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
          <div v-if="!whatsapp.isLoading" class="card-pf-body">
            <form class="form-horizontal" role="form">
              <div class="form-group">
                <label class="col-sm-4 control-label">{{ $t("profile.whatsapp_login_threshold") }}</label>
                <div class="col-sm-8">
                  <input v-model="whatsapp.data.whatsapp_threshold" required type="number" class="form-control">
                </div>
              </div>
            </form>
            <p>{{ $t("profile.whatsapp_threshold_description") }}.</p>
          </div>
          <div class="card-pf-footer">
            <div class="dropdown card-pf-time-frame-filter">
              <button class="btn btn-default" v-on:click="updateWhatsappThreshold()">{{ $t("update") }}</button>
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

export default {
  name: "Profile",
  mixins: [LoginService, StorageService, UtilService, StatsService],
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
      whatsapp: {
        isLoading: true,
        data: {}
      }
    };
  },
  mounted() {
    this.getSmsThreshold();
    this.getWhatsappThreshold();
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
    updateWhatsappThreshold() {
      this.whatsapp.isLoading = true;

      this.updateWhatsappThresholdForAccountByAccount(
        {
          whatsapp_threshold: parseInt(this.whatsapp.data.whatsapp_threshold) || 0
        },
        this.user.login.id,
        success => {
          this.whatsapp.isLoading = false;
          this.whatsapp.data = success.body;
          this.getWhatsappThreshold();
        },
        error => {
          this.whatsapp.isLoading = false;
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
    getWhatsappThreshold() {
      this.statsWhatsappTotalForAccountByAccount(
        this.user.login.id,
        success => {
          this.whatsapp.isLoading = false;
          this.whatsapp.data = success.body;
        },
        error => {
          this.whatsapp.isLoading = false;
          console.error(error.body);
        }
      );
    }
  }
};
</script>

<style>

</style>
