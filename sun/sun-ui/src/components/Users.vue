<template>
  <div>
    <h2>{{ msg }}</h2>
    <div v-if="isLoading" class="spinner spinner-lg"></div>
    <div
      v-if="(user.account_type == 'admin') || (user.account_type == 'reseller') && !isLoading"
      class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12"
    >
      <label v-if="!isLoading" class="col-sm-2 control-label" for="textInput-markup">Hotspot</label>
      <div v-if="!isLoading" class="col-sm-4">
        <select v-on:change="getAll(true)" v-model="hotspotSearchId" class="form-control">
          <option
            v-for="hotspot in hotspots"
            v-bind:key="hotspot.id"
            v-bind:value="hotspot.id"
          >{{ hotspot.name }} - {{ hotspot.description}}</option>
        </select>
      </div>
    </div>
    <div v-if="!isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <label class="col-sm-2 control-label" for="textInput-markup">{{$t('user.show_expired')}}</label>
      <div class="col-sm-4">
        <input
          @click="toggleExpire()"
          v-model="hotspotShowExpired"
          type="checkbox"
          id="textInput2-modal-markup"
          class="form-control"
        >
      </div>
    </div>
    <div v-if="!isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <label
        class="col-sm-2 control-label"
        for="textInput-markup"
      >{{$t('user.show_marketing_auth')}}</label>
      <div class="col-sm-4">
        <input
          @click="toggleMarketing()"
          v-model="hotspotShowMarketing"
          type="checkbox"
          id="textInput2-modal-markup"
          class="form-control"
        >
      </div>
    </div>
    <div v-if="!isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <div class="col-sm-3">
        <button class="btn btn-primary" @click="getAll()">{{$t('session.refresh')}}</button>
      </div>
    </div>
    <div v-if="!isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <div class="col-sm-12">
        <button
          @click="exportCSVUsers()"
          class="btn btn-primary export-btn"
        >{{$t('session.export_csv')}}</button>
      </div>
      <div
        class="result-list adjust-results"
      >{{total}} {{total == 1 ? $t('result') : $t('results')}}</div>
    </div>
    <div v-if="!isLoading">
      <form v-on:submit.prevent="searchFn($event)">
        <input
          class="form-control input-lg search-table-input"
          type="text"
          :placeholder="tableLangsTexts.globalSearchPlaceholder"
        >
      </form>
    </div>
    <div v-if="!isLoading && isLoadingTable" class="spinner spinner-lg spinner-adjust"></div>
    <div
      v-if="!isLoadingTable && !isLoading && exportError"
      class="alert alert-danger alert-dismissable alert-export"
    >
      <span class="pficon pficon-error-circle-o"></span>
      <strong>{{$t('session.export_error')}}</strong>
      . {{$t('session.export_error_details')}}.
    </div>
    <vue-good-table
      v-if="!isLoadingTable && !isLoading"
      @perPageChanged="handlePerPage"
      :customRowsPerPageDropdown="[25,50,100]"
      :perPage="hotspotPerPage"
      :columns="columns"
      :rows="rows"
      :lineNumbers="false"
      :defaultSortBy="{field: 'username', type: 'asc'}"
      :globalSearch="true"
      :globalSearchFn="searchFn"
      :paginate="false"
      styleClass="table"
      :nextText="tableLangsTexts.nextText"
      :prevText="tableLangsTexts.prevText"
      :rowsPerPageText="tableLangsTexts.rowsPerPageText"
      :globalSearchPlaceholder="tableLangsTexts.globalSearchPlaceholder"
      :ofText="tableLangsTexts.ofText"
    >
      <template slot="table-row" slot-scope="props">
        <td
          :class="[isExpired(props.row.valid_until) ? 'disabled' : '', 'fancy']"
        >{{ props.row.name }}</td>
        <td
          :class="[isExpired(props.row.valid_until) ? 'disabled' : '', 'fancy']"
        >{{ props.row.email || '-' }}</td>
        <td :class="[isExpired(props.row.valid_until) ? 'disabled' : '', 'fancy']">
          <strong>{{ props.row.account_type }}</strong>
        </td>
        <td
          :title="props.row.reason.length > 0 ? $t('user.'+props.row.reason) : '-'"
          :class="[isExpired(props.row.valid_until) ? 'disabled' : '', 'fancy']"
          v-if="reasonView"
        >
          <span :class="getReasonIcon(props.row.reason)">{{props.row.reason ? '' : '-'}}</span>
        </td>
        <td
          :class="[isExpired(props.row.valid_until) ? 'disabled' : '', 'fancy']"
          v-if="reasonView"
        >
          <flag class="adjust-flag adjust-size-icon" :iso="props.row.country"/>
          <span>{{props.row.country ? '' : '-'}}</span>
        </td>
        <td class="fancy">
          <span
            :class="['pficon', props.row.marketing_auth ? 'pficon-ok' : 'pficon-error-circle-o']"
          ></span>
        </td>
        <td v-if="surveyView" class="fancy">
          <span :class="['pficon', props.row.survey_auth ? 'pficon-ok' : 'pficon-error-circle-o']"></span>
        </td>
        <td class="fancy">
          <span :class="['pficon', props.row.auto_login ? 'pficon-ok' : 'pficon-error-circle-o']"></span>
        </td>
        <td :class="[isExpired(props.row.valid_until) ? 'disabled' : '', 'fancy']">
          <div>
            <strong>{{ $t('user.kbps_down') }}</strong>
            : {{ props.row.kbps_down || '-' }}
          </div>
          <div>
            <strong>{{ $t('user.kbps_up') }}</strong>
            : {{ props.row.kbps_up || '-' }}
          </div>
        </td>
        <td :class="[isExpired(props.row.valid_until) ? 'disabled' : '', 'fancy']">
          <div>
            <strong>{{ $t('user.traffic') }}</strong>
            : {{ props.row.max_navigation_traffic | byteFormat }}
          </div>
          <div>
            <strong>{{ $t('user.time') }}</strong>
            : {{ props.row.max_navigation_time | secondsInHour }}
          </div>
        </td>
        <td :class="[isExpired(props.row.valid_until) ? 'disabled' : '', 'fancy']">
          <div>
            <strong>{{ $t('user.from') }}</strong>
            : {{ props.row.valid_from | formatDate }}
          </div>
          <div>
            <strong>{{ $t('user.until') }}</strong>
            : {{ props.row.valid_until | formatDate }}
            <span
              :class="['pficon', isExpired(props.row.valid_until) ? 'pficon-error-circle-o': 'pficon-ok']"
            ></span>
          </div>
        </td>
        <td>
          <user-action
            :expired="hotspotShowExpired"
            details="false"
            :obj="props.row"
            :update="getAll"
          ></user-action>
        </td>
      </template>
    </vue-good-table>
    <div v-if="!isLoadingTable && !isLoading" class="right paginator">
      <span class="page-count">
        <b>{{hotspotPage}}</b>
        {{tableLangsTexts.ofText}} {{total / hotspotPerPage | adjustPage}} (
        <b>{{hotspotPerPage}}</b>
        {{tableLangsTexts.rowsPerPageText}})
      </span>
      <button
        :disabled="availablePrevPage()"
        @click="prevPage()"
        class="btn btn-default"
      >{{tableLangsTexts.prevText}}</button>
      <button
        :disabled="availableNextPage()"
        @click="nextPage()"
        class="btn btn-default"
      >{{tableLangsTexts.nextText}}</button>
    </div>
  </div>
</template>

<script>
import UserService from "../services/user";
import StorageService from "../services/storage";
import HotspotService from "../services/hotspot";
import UtilService from "../services/util";
import MarketingService from "../services/marketing";

import UserAction from "../directives/UserAction.vue";

export default {
  name: "Users",
  mixins: [
    UserService,
    StorageService,
    UtilService,
    HotspotService,
    MarketingService
  ],
  components: {
    userAction: UserAction
  },
  data() {
    return {
      msg: this.$i18n.t("menu.users"),
      isLoading: true,
      isLoadingTable: true,
      columns: [
        {
          label: this.$i18n.t("user.name"),
          field: "name",
          filterable: true
        },
        {
          label: this.$i18n.t("user.email"),
          field: "email",
          filterable: true
        },
        {
          label: this.$i18n.t("user.type"),
          field: "account_type",
          filterable: true
        },
        {
          label: this.$i18n.t("user.reason"),
          field: "account_reason",
          filterable: true
        },
        {
          label: this.$i18n.t("user.country"),
          field: "account_country",
          filterable: true
        },
        {
          label: this.$i18n.t("user.marketing_auth"),
          field: "marketing_auth",
          filterable: true
        },
        {
          label: this.$i18n.t("user.survey_auth"),
          field: "survey_auth",
          filterable: true
        },
        {
          label: this.$i18n.t("user.auto_login"),
          field: "auto_login",
          filterable: true
        },
        {
          label: this.$i18n.t("user.bandwidth_limit"),
          field: "bandwidth",
          filterable: true,
          sortable: false
        },
        {
          label: this.$i18n.t("user.daily_limit"),
          field: "limit",
          filterable: true,
          sortable: false
        },
        {
          label: this.$i18n.t("user.valid"),
          field: "valid",
          filterable: true,
          sortable: false
        },
        {
          label: this.$i18n.t("action"),
          field: "",
          sortable: false
        }
      ],
      rows: [],
      tableLangsTexts: this.tableLangs(),
      hotspots: [],
      hotspotSearchId: 0,
      hotspotShowExpired: this.get("users_show_expired") || false,
      hotspotShowMarketing: this.get("users_show_marketing") || false,
      hotspotPerPage: 25,
      hotspotPage: 1,
      total: 0,
      user: this.get("loggedUser") || null,
      searchString: "",
      exportError: false,
      reasonView: true,
      surveyView: true
    };
  },
  mounted() {
    if (this.$route.params.hotspotId !== undefined) {
      this.hotspotSearchId = this.$route.params.hotspotId;
    }
    // get user list
    var context = this;
    this.getAllHotspots(function() {
      context.getAll();
      context.getMarketingInfo();
    });
  },
  methods: {
    getReasonIcon(reason) {
      switch (reason) {
        case "business":
          return "fa fa-briefcase adjust-size-icon";
          break;
        case "family":
          return "fa fa-child adjust-size-icon";
          break;
        case "other":
          return "fa fa-question adjust-size-icon";
          break;
      }
    },
    handlePerPage(evt) {
      this.set("users_per_page", evt.currentPerPage);
    },
    searchFn(evt) {
      var elem = evt.srcElement || evt.target;
      this.searchString = elem[0].value;
      this.getAll(true);
    },
    isExpired(date) {
      return new Date().toISOString() > date;
    },
    getAllHotspots(callback) {
      this.hotspotGetAll(
        null,
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
          if (this.$route.params.hotspotId === undefined) {
            this.hotspotSearchId = hsId;
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
    getMarketingInfo() {
      if (
        this.$parent.user.info.type == "customer" ||
        this.$parent.user.info.type == "desk"
      ) {
        this.marketingPrefGet(
          this.$parent.user.info.hotspot_id,
          success => {
            var reasonView = success.body.filter(function(p) {
              return p.key == "marketing_0_reason_country";
            })[0];
            var surveyView = success.body.filter(function(p) {
              return p.key == "marketing_1_enabled";
            })[0];
            this.columns[3].hidden = reasonView.value == "false";
            this.columns[4].hidden = reasonView.value == "false";
            this.columns[6].hidden = surveyView.value == "false";
            this.reasonView = reasonView.value == "true";
            this.surveyView = surveyView.value == "true";
            this.$forceUpdate();
          },
          error => {
            console.error(error.body);
          }
        );
      }
    },
    toggleExpire() {
      this.isLoading = true;
      this.hotspotShowExpired = !this.hotspotShowExpired;
      this.set("users_show_expired", this.hotspotShowExpired);
      this.getAll(true);
    },
    toggleMarketing() {
      this.isLoading = true;
      this.hotspotShowMarketing = !this.hotspotShowMarketing;
      this.set("users_show_marketing", this.hotspotShowMarketing);
      this.getAll(true);
    },
    getAll(reset) {
      this.isLoadingTable = true;

      if (reset) {
        this.hotspotPage = 1;
        this.total = 0;
      }

      this.set(
        "users_hotspot_id",
        this.hotspotSearchId || this.get("users_hotspot_id") || 0
      );

      this.userGetAll(
        this.hotspotSearchId,
        null,
        this.hotspotShowExpired,
        this.hotspotPage,
        this.hotspotPerPage,
        encodeURIComponent(this.searchString),
        success => {
          this.rows = [];
          if (this.hotspotShowExpired) {
            for (var s in success.body.data_users) {
              var res = success.body.data_users[s];
              this.rows.push(res);
            }
            for (var s in success.body.data_user_histories) {
              var res = success.body.data_user_histories[s];
              this.rows.push(res);
            }
          } else {
            for (var s in success.body.data) {
              var res = success.body.data[s];
              this.rows.push(res);
            }
          }
          this.total = success.body.total;
          this.isLoading = false;
          this.isLoadingTable = false;
        },
        error => {
          this.isLoading = false;
          this.isLoadingTable = false;
          this.rows = [];
          console.error(error);
        },
        this.hotspotShowMarketing
      );
    },
    exportCSVUsers() {
      this.isLoadingTable = true;
      // get users
      this.userGetAll(
        this.hotspotSearchId,
        null,
        this.hotspotShowExpired,
        null,
        null,
        encodeURIComponent(this.searchString),
        success => {
          var data_export = [];
          if (this.hotspotShowExpired) {
            for (var s in success.body.data_users) {
              var res = success.body.data_users[s];
              data_export.push(res);
            }
            for (var s in success.body.data_user_histories) {
              var res = success.body.data_user_histories[s];
              data_export.push(res);
            }
          } else {
            for (var s in success.body.data) {
              var res = success.body.data[s];
              data_export.push(res);
            }
          }

          if (data_export.length < 5000) {
            var usersRows = JSON.parse(JSON.stringify(data_export));
            for (var r in usersRows) {
              // get only email verified users
              if (
                usersRows[r].account_type == "email" &&
                !usersRows[r].email_verified
              ) {
                delete usersRows[r];
              }
            }

            var columns = this.columns.slice();

            delete columns[3];
            delete columns[4];
            delete columns[5];
            delete columns[6];
            delete columns[7];

            var csv = this.createCSV(columns, usersRows);
            this.isLoadingTable = false;
            this.downloadCSV(csv.cols, csv.rows, "users");
          } else {
            this.isLoadingTable = false;
            this.exportError = true;
          }
        },
        error => {
          this.isLoading = false;
          this.isLoadingTable = false;
          console.error(error);
        },
        this.hotspotShowMarketing
      );
    },
    prevPage() {
      this.isLoading = true;
      if (this.hotspotPage != 1) {
        this.hotspotPage--;
      }
      this.getAll();
    },
    nextPage() {
      this.isLoading = true;
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
.adjust-results {
  right: 45px !important;
  top: 0px;
}
.adjust-flag {
  margin-right: 5px !important;
}
.adjust-size-icon {
  font-size: 18px !important;
  margin-top: 8px;
}
</style>