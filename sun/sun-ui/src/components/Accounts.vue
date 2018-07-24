<template>
  <div>
    <h2>{{ msg }}</h2>
    <button v-on:click="initNewAccount()" data-toggle="modal" data-target="#ACcreateModal" class="btn btn-primary btn-lg create-account">
      {{ $t('account.create_new') }} </button>
    <div v-if="isLoading" class="spinner spinner-lg"></div>
    <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller') && !isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <label v-if="!isLoading" class="col-sm-2 control-label" for="textInput-markup">Hotspot</label>
      <div v-if="!isLoading" class="col-sm-4">
        <select v-on:change="getAll()" v-model="hotspotSearchId" class="form-control">
          <option v-bind:value="0">-</option>
          <option v-for="hotspot in hotspots" v-bind:key="hotspot.id" v-bind:value="hotspot.id">
            {{ hotspot.name }}
          </option>
        </select>
      </div>
    </div>
    <div v-if="!isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <div class="col-sm-3">
        <button class="btn btn-primary" @click="getAll()">{{$t('session.refresh')}}</button>
      </div>
    </div>
    <div v-if="!isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <div class="result-list">{{total}} {{total == 1 ? $t('result') : $t('results')}}</div>
    </div>
    <vue-good-table v-if="!isLoading" @perPageChanged="handlePerPage" :customRowsPerPageDropdown="[25,50,100]" :perPage="hotspotPerPage"
      :columns="columns" :rows="rows" :lineNumbers="false" :defaultSortBy="{field: 'username', type: 'asc'}" :globalSearch="true"
      :paginate="false" styleClass="table" :nextText="tableLangsTexts.nextText" :prevText="tableLangsTexts.prevText" :rowsPerPageText="tableLangsTexts.rowsPerPageText"
      :globalSearchPlaceholder="tableLangsTexts.globalSearchPlaceholder" :ofText="tableLangsTexts.ofText">
      <template slot="table-row" slot-scope="props">
        <td>
          <a :href="'#/accounts/'+ props.row.id">
            {{ props.row.username }}
          </a>
        </td>
        <td class="fancy">{{ props.row.name }}</td>
        <td class="fancy">{{ props.row.email || '-' }}</td>
        <td>
          <span :class="getLoginIcon(props.row.type)" data-toggle="tooltip" data-placement="left" :title="$t(props.row.type)"></span>
          {{$t(props.row.type)}}
        </td>
        <td v-if="(user.account_type == 'admin')" class="fancy">{{ props.row.subscription.subscription_plan.name || '-' }}</td>
        <td class="fancy">{{ props.row.created | formatDate }}</td>
        <td>
          <account-action details="false" :obj="props.row" :update="getAll"></account-action>
        </td>
      </template>
    </vue-good-table>
    <div v-if="!isLoading" class="right paginator">
      <span class="page-count">
        <b>{{hotspotPage}}</b> {{tableLangsTexts.ofText}} {{total / hotspotPerPage | adjustPage}} (
        <b>{{hotspotPerPage}}</b> {{tableLangsTexts.rowsPerPageText}})</span>
      <button :disabled="availablePrevPage()" @click="prevPage()" class="btn btn-default">{{tableLangsTexts.prevText}}</button>
      <button :disabled="availableNextPage()" @click="nextPage()" class="btn btn-default">{{tableLangsTexts.nextText}}</button>
    </div>

    <div class="modal fade" id="ACcreateModal" tabindex="-1" role="dialog" aria-labelledby="ACModifyModalLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="ACModifyModalLabel">{{ $t("account.create_title") }}</h4>
          </div>
          <form class="form-horizontal" role="form" v-on:submit.prevent="createAccount(newObj)">
            <div class="modal-body">
              <div class="form-group" v-if="isAdmin">
                <label class="col-sm-4 control-label" for="accuuid">{{ $t("account.uuid") }}</label>
                <div class="col-sm-8">
                  <input required v-model="newObj.uuid" type="text" id="accuuid" class="form-control" :placeholder="$t('account.uuid')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="accusername">{{ $t("account.username") }}</label>
                <div class="col-sm-8">
                  <input required v-model="newObj.username" type="text" id="accusername" class="form-control" :placeholder="$t('account.username')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="accountname">{{ $t("account.name") }}</label>
                <div class="col-sm-8">
                  <input required v-model="newObj.name" type="text" id="accountname" class="form-control" :placeholder="$t('account.name')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="ACtextInput2-modal-markup">{{ $t("account.email") }}</label>
                <div class="col-sm-8">
                  <input required v-model="newObj.email" type="email" id="ACtextInput2-modal-markup" class="form-control" :placeholder="$t('account.email')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="ACtextInput2-modal-markup">{{ $t("account.type") }}</label>
                <div class="col-sm-8">
                  <select v-model="newObj.type" class="form-control">
                    <option value="customer" selected>{{ $t("account.type_customer") }}</option>
                    <option value="desk">{{ $t("account.type_desk") }}</option>
                    <option v-if="accountType == 'admin'" value="reseller">{{ $t("account.type_reseller") }}</option>
                  </select>
                </div>
              </div>
              <div class="form-group" v-if="(isAdmin && newObj.type == 'reseller')">
                <label class="col-sm-4 control-label" for="accuuid">{{ $t("account.subscription_plan_name") }}</label>
                <div class="col-sm-8">
                  <select v-model="newObj.subscription_plan_id" class="form-control">
                    <option v-for="plan in plans" v-bind:key="plan.id" v-bind:value="plan.id">
                      {{ plan.name }}
                    </option>
                  </select>
                </div>
              </div>
              <div class="form-group" v-if="(newObj.type == 'customer' || newObj.type == 'desk')">
                <label class="col-sm-4 control-label" for="ACtextInput2-modal-markup">{{ $t("account.hotspot") }}</label>
                <div class="col-sm-8">
                  <select v-model="newObj.hotspot_id" class="form-control">
                    <option v-for="hotspot in hotspots" v-bind:key="hotspot.id" v-bind:value="hotspot.id">
                      {{ hotspot.name }}
                    </option>
                  </select>
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="textInput-modal-markup">{{ $t("account.new_password") }}</label>
                <div class="col-sm-8">
                  <input v-model="newPassword" required type="text" id="textInput-modal-markup" class="form-control" :placeholder="$t('account.new_password')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="textInput2-modal-markup">{{ $t("account.confirm_password") }}</label>
                <div class="col-sm-8">
                  <input v-model="confirmPassword" required type="text" id="textInput2-modal-markup" class="form-control" :placeholder="$t('account.confirm_password')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="textInput2-modal-markup"></label>
                <div class="col-sm-8">
                  <span v-bind:class="[newPassword == confirmPassword && newPassword.length > 0 ? 'pass-confirm-ok' : 'pass-confirm-err', '']"></span>
                </div>
              </div>
              <div v-if="errors.create" class="alert alert-danger alert-dismissable">
                <span class="pficon pficon-error-circle-o"></span>
                <strong>{{ $t("account.create_error_title") }}</strong>.
                <span>{{ errors.status != 409 ? $t("account.create_error_sub") : $t("account.duplicate_error") }}.</span>
              </div>
            </div>
            <div class="modal-footer">
              <span v-if="newObj.onAction" class="spinner spinner-sm spinner-inline modal-spinner"></span>
              <button type="button" class="btn btn-default" data-dismiss="modal">{{ $t("cancel") }}</button>
              <button :disabled="(newPassword != confirmPassword) || (confirmPassword.length == 0) || (newPassword.length == 0) || (!newObj.hotspot_id && newObj.type != 'reseller')"
                type="submit" class="btn btn-primary">{{ $t("create") }}</button>
            </div>
          </form>
        </div>
      </div>
    </div>

  </div>
</template>


<script>
import AccountService from "../services/account";
import StorageService from "../services/storage";
import HotspotService from "../services/hotspot";
import SubscriptionService from "../services/subscription";
import UtilService from "../services/util";

import AccountAction from "../directives/AccountAction.vue";

export default {
  name: "Accounts",
  mixins: [
    AccountService,
    StorageService,
    UtilService,
    HotspotService,
    SubscriptionService
  ],
  components: {
    accountAction: AccountAction
  },
  data() {
    var newObj = {
      uuid: "",
      username: "",
      name: "",
      email: "",
      type: "",
      password: "",
      hotspot_id: 0,
      subscription_plan_id: 0
    };

    var errors = {
      create: false
    };

    var newPassword = "";
    var confirmPassword = "";

    return {
      msg: this.$i18n.t("menu.accounts"),
      isLoading: true,
      accountType: this.get("loggedUser").account_type,
      columns: [
        {
          label: this.$i18n.t("account.username"),
          field: "username",
          filterable: true
        },
        {
          label: this.$i18n.t("account.name"),
          field: "name",
          filterable: true
        },
        {
          label: this.$i18n.t("account.email"),
          field: "email",
          filterable: true
        },
        {
          label: this.$i18n.t("account.type"),
          field: "type",
          filterable: true
        },
        {
          label: this.$i18n.t("account.subscription_plan_name"),
          field: "subscription.subscription_plan.name",
          filterable: true,
          hidden: this.get("loggedUser").account_type != "admin" ? true : false
        },
        {
          label: this.$i18n.t("account.created"),
          field: "created",
          filterable: false
        },
        {
          label: this.$i18n.t("action"),
          field: "",
          sortable: false
        }
      ],
      rows: [],
      hotspots: [],
      plans: [],
      tableLangsTexts: this.tableLangs(),
      newObj: newObj,
      newPassword: newPassword,
      confirmPassword: confirmPassword,
      errors: errors,
      isAdmin: this.get("loggedUser").account_type == "admin",
      hotspotSearchId: 0,
      hotspotPerPage: 25,
      hotspotPage: 1,
      total: 0,
      user: this.get("loggedUser") || null
    };
  },
  mounted() {
    if (this.$route.params.hotspotId !== undefined) {
      this.hotspotSearchId = this.$route.params.hotspotId;
    }

    // get account list
    var context = this;
    this.getAllHotspots(function() {
      context.getAll();
      context.getAllSubscriptionPlans();
    });
  },
  methods: {
    handlePerPage(evt) {
      this.set("accounts_per_page", evt.currentPerPage);
    },
    initNewAccount() {
      this.newObj.uuid = this.generateUUID();
      this.newObj.password = this.generatePassword();
      this.newPassword = this.newObj.password;
      this.confirmPassword = this.newObj.password;
      this.newObj.type = this.accountType == "admin" ? "reseller" : "customer";
      this.newObj.username = "";
      this.newObj.name = "";
      this.newObj.email = "";
      this.newObj.hotspot_id = 0;
      this.newObj.subscription_plan_id = 0;
    },
    createAccount(obj) {
      this.newObj.onAction = true;
      obj.password = this.newPassword;
      this.accountCreate(
        obj,
        success => {
          this.newObj.onAction = false;
          $("#ACcreateModal").modal("toggle");
          this.getAll();
        },
        error => {
          this.newObj.onAction = false;
          this.errors.create = true;
          this.errors.status = error.status;
          console.error(error);
        }
      );
    },
    getAllHotspots(callback) {
      this.hotspotGetAll(
        null,
        null,
        success => {
          this.hotspots = success.body.data;
          var hsId = this.get("users_hotspot_id") || this.hotspots[0].id;
          if (
            this.$parent.user.info.type == "customer" ||
            this.$parent.user.info.type == "desk"
          ) {
            hsId = this.$parent.user.info.hotspot_id;
          }

          $('[data-toggle="tooltip"]').tooltip();
          this.isLoading = false;

          callback();
        },
        error => {
          console.error(error);

          callback();
        }
      );
    },
    getAllSubscriptionPlans() {
      this.subscriptionPlansGetAll(
        success => {
          this.plans = success.body;
          $('[data-toggle="tooltip"]').tooltip();
          this.isLoading = false;
        },
        error => {
          console.error(error);
        }
      );
    },
    getAll() {
      this.accountGetAll(
        this.hotspotSearchId,
        this.hotspotPage,
        this.hotspotPerPage,
        success => {
          this.rows = success.body.data;
          this.total = success.body.total;
          this.isLoading = false;
        },
        error => {
          console.error(error);
          this.rows = [];
          this.isLoading = false;
        }
      );
    },
    prevPage() {
      if (this.hotspotPage != 1) {
        this.hotspotPage--;
      }
      this.getAll();
    },
    nextPage() {
      if (this.hotspotPage != Math.ceil(this.total / this.hotspotPerPage)) {
        this.hotspotPage++;
      }
      this.getAll();
    },
    availablePrevPage() {
      return this.hotspotPage == 1;
    },
    availableNextPage() {
      return this.hotspotPage == Math.ceil(this.total / this.hotspotPerPage);
    }
  }
};
</script>

<style>
.create-account {
  float: right;
  margin-top: -52px;
  margin-right: 35px;
}
</style>