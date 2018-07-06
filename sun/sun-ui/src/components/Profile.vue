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
          <div class="card-pf-footer">
            <div class="dropdown card-pf-time-frame-filter">
              <renew-button :obj="user.info" :status="billingEmpty" :update="getSubscriptionInfo"></renew-button>
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
              {{$t('profile.billing')}}
              <div v-if="isBillingLoading" class="spinner spinner-sm loader-spinner right"></div>
              <div v-if="!isBillingLoading" class="fa fa-list right"></div>
            </h2>
          </div>
          <form class="form-horizontal" v-on:submit.prevent="updateBillingInfo()">
            <div class="card-pf-body">
              <div class="form-group">
                <label class="col-sm-3 control-label" for="textInput-markup">{{$t('profile.name')}}</label>
                <div class="col-sm-9">
                  <input v-model="billingInfo.name" required type="text" id="textInput-markup" class="form-control">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-3 control-label" for="textInput-markup">{{$t('profile.address')}}</label>
                <div class="col-sm-9">
                  <input v-model="billingInfo.address" required type="text" id="textInput-markup" class="form-control">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-3 control-label" for="textInput-markup">{{$t('profile.city')}}</label>
                <div class="col-sm-9">
                  <input v-model="billingInfo.city" required type="text" id="textInput-markup" class="form-control">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-3 control-label" for="textInput-markup">{{$t('profile.postal_code')}}</label>
                <div class="col-sm-9">
                  <input v-model="billingInfo.postal_code" required type="text" id="textInput-markup" class="form-control">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-3 control-label" for="textInput-markup">{{$t('profile.country')}}</label>
                <div class="col-sm-9">
                  <select required v-model="billingInfo.country" class="form-control">
                    <option v-for="c in countries" v-bind:key="c.id" v-bind:value="c.country">
                      {{ c.country }}
                    </option>
                  </select>
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-3 control-label" for="textInput-markup">{{$t('profile.type')}}</label>
                <div class="col-sm-9">
                  <span class="span-radio">
                    <input required v-model="billingType" class="form-check-input" type="radio" name="exampleRadios" id="exampleRadios1" value="business">
                    <label class="form-check-label" for="exampleRadios1">
                      {{$t('profile.type_business')}}
                    </label>
                  </span>
                  <span class="span-radio">
                    <input required v-model="billingType" class="form-check-input" type="radio" name="exampleRadios" id="exampleRadios2" value="person">
                    <label class="form-check-label" for="exampleRadios2">
                      {{$t('profile.type_person')}}
                    </label>
                  </span>
                </div>
              </div>
              <div v-if="billingType == 'business'" class="form-group">
                <label class="col-sm-3 control-label" for="textInput-markup">{{$t('profile.vat')}}</label>
                <div class="col-sm-9">
                  <input v-model="billingInfo.vat" required type="text" id="textInput-markup" class="form-control">
                </div>
              </div>
            </div>
            <div class="card-pf-footer">
              <div class="dropdown card-pf-time-frame-filter">
                <button :disabled="isSaving" type="submit" class="btn btn-primary">Save</button>
              </div>
              <p>
                <a href="#" class="card-pf-link-with-icon">
                </a>
              </p>
            </div>
          </form>

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
import AccountService from "../services/account";
import LoginService from "../services/login";
import StorageService from "../services/storage";
import UtilService from "../services/util";
import StatsService from "../services/stats";
import SubscriptionService from "../services/subscription";

import RenewButton from "./details-view/RenewButton.vue";

export default {
  name: "Profile",
  mixins: [
    AccountService,
    LoginService,
    StorageService,
    UtilService,
    StatsService,
    SubscriptionService
  ],
  components: {
    renewButton: RenewButton
  },
  mounted: function() {
    $('[data-toggle="tooltip"]').tooltip();
  },
  data() {
    var updateBilling = false;
    if (this.$parent.action == "updateBilling") {
      updateBilling = true;
      this.delete("query_params");
    }

    // get subscription info
    this.getSubscriptionInfo()

    // get country list
    this.getCountryList();

    // get billing info
    this.getBillingInfo();

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
      updateBilling: updateBilling,
      billingInfo: {},
      billingType: "person",
      billingEmpty: true,
      countries: [],
      isSaving: false,
      isBillingLoading: true
    };
  },
  methods: {
    isAuth0() {
      return this.get("auth0User");
    },
    isExpired(date) {
      return new Date().toISOString() > date;
    },
    getSubscriptionInfo() {
      this.accountGet(this.get("loggedUser").id, success => {
        this.user.info.subscription = success.body.subscription;
      }, error => {

      });
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
          console.log(error.body.message);
        }
      );
    },
    getBillingInfo() {
      this.isBillingLoading = true;
      this.subscriptionBillingsGetAll(
        success => {
          this.isBillingLoading = false;
          this.billingEmpty = false;
          this.billingInfo = success.body;
          if (success.body.vat && success.body.vat.length > 0) {
            this.billingType = "business";
          } else {
            this.billingType = "person";
          }
        },
        error => {
          this.billingEmpty = true;
          this.billingType = "person";
          this.isBillingLoading = false;
          console.error(error);
        }
      );
    },
    getCountryList() {
      this.subscriptionTaxesGetAll(
        success => {
          this.countries = success.body;
        },
        error => {
          console.error(error);
        }
      );
    },
    updateBillingInfo() {
      this.isSaving = true;
      if (this.billingType == "person") {
        this.billingInfo.vat = "";
      }

      this.subscriptionBillingUpdate(
        this.billingEmpty,
        success => {
          this.isSaving = false;
          this.updateBilling = false;
          this.billingEmpty = false;
        },
        error => {
          this.isSaving = false;
          console.error(error);
        }
      );
    }
  }
};
</script>

<style>
</style>