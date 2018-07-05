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
          <div v-if="!isAuth0()" class="card-pf-footer">
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

      <div v-if="isAuth0()" class="col-xs-12 col-sm-12 col-md-6">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{ $t("profile.subscription") }}
              <div :class="[isExpired(user.info.subscription.valid_until) ? 'pficon pficon-error-circle-o' : 'pficon pficon-ok', 'right']"
                data-toggle="tooltip" data-placement="left" :title="isExpired(user.info.subscription.valid_until) ? $t('profile.invalid') : $t('profile.valid')"></div>
            </h2>
          </div>
          <div class="card-pf-body">
            <div class="list-details">
              <dt>{{ $t("profile.name") }}</dt>
              <dd>{{user.info.subscription.subscription_plan.name}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("profile.description") }}</dt>
              <dd>{{user.info.subscription.subscription_plan.description}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("profile.valid_from") }}</dt>
              <dd>{{user.info.subscription.valid_from | formatDate}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("profile.valid_until") }}</dt>
              <dd>{{user.info.subscription.valid_until | formatDate}}</dd>
            </div>
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
    mounted: function () {
      $('[data-toggle="tooltip"]').tooltip();
    },
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
        onAction: false
      };
    },
    methods: {
      isAuth0() {
        return this.get("auth0User");
      },
      isExpired(date) {
        return new Date().toISOString() > date;
      },
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
      }
    }
  };
</script>

<style>
</style>