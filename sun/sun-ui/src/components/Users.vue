<template>
  <div>
    <h2>{{ msg }}</h2>
    <div v-if="isLoading" class="spinner spinner-lg"></div>
    <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller') && !isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <label v-if="!isLoading" class="col-sm-2 control-label" for="textInput-markup">Hotspot</label>
      <div v-if="!isLoading" class="col-sm-4">
        <select v-on:change="getAll()" v-model="hotspotSearchId" class="form-control">
          <option v-for="hotspot in hotspots" v-bind:key="hotspot.id" v-bind:value="hotspot.id">
            {{ hotspot.name }} - {{ hotspot.description}}
          </option>
        </select>
      </div>
    </div>
    <div v-if="!isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <label class="col-sm-2 control-label" for="textInput-markup">{{$t('user.show_expired')}}</label>
      <div class="col-sm-4">
        <input @click="toggleExpire()" v-model="hotspotShowExpired" type="checkbox" id="textInput2-modal-markup" class="form-control">
      </div>
    </div>
    <div v-if="!isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <div class="col-sm-3">
        <button class="btn btn-primary" @click="getAll()">{{$t('session.refresh')}}</button>
      </div>
    </div>
    <div v-if="!isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <div class="col-sm-12">
          <button @click="exportCSVUsers()" class="btn btn-primary export-btn">{{$t('session.export_csv')}}</button>
        </div>
      <div class="result-list adjust-results">{{rows.length}} {{rows.length == 1 ? $t('result') : $t('results')}}</div>
    </div>
    <vue-good-table v-if="!isLoading" @perPageChanged="handlePerPage" :customRowsPerPageDropdown="[25,50,100]" :perPage="hotspotPerPage"
      :columns="columns" :rows="rows" :lineNumbers="false" :defaultSortBy="{field: 'username', type: 'asc'}" :globalSearch="true"
      :paginate="true" styleClass="table" :nextText="tableLangsTexts.nextText" :prevText="tableLangsTexts.prevText" :rowsPerPageText="tableLangsTexts.rowsPerPageText"
      :globalSearchPlaceholder="tableLangsTexts.globalSearchPlaceholder" :ofText="tableLangsTexts.ofText">
      <template slot="table-row" slot-scope="props">
        <td :class="[isExpired(props.row.valid_until) ? 'disabled' : '', 'fancy']">{{ props.row.name }}</td>
        <td :class="[isExpired(props.row.valid_until) ? 'disabled' : '', 'fancy']">{{ props.row.email || '-' }}</td>
        <td :class="[isExpired(props.row.valid_until) ? 'disabled' : '', 'fancy']">
          <strong>{{ props.row.account_type }}</strong>
        </td>
        <td class="fancy">
          <span :class="['pficon', props.row.auto_login ? 'pficon-ok' : 'pficon-error-circle-o']"></span>
        </td>
        <td :class="[isExpired(props.row.valid_until) ? 'disabled' : '', 'fancy']">
          <div>
            <strong>{{ $t('user.kbps_down') }}</strong>: {{ props.row.kbps_down || '-' }}</div>
          <div>
            <strong>{{ $t('user.kbps_up') }}</strong>: {{ props.row.kbps_up || '-' }}</div>
        </td>
        <td :class="[isExpired(props.row.valid_until) ? 'disabled' : '', 'fancy']">
          <div>
            <strong>{{ $t('user.traffic') }}</strong>: {{ props.row.max_navigation_traffic | byteFormat }}</div>
          <div>
            <strong>{{ $t('user.time') }}</strong>: {{ props.row.max_navigation_time | secondsInHour }}
          </div>
        </td>
        <td :class="[isExpired(props.row.valid_until) ? 'disabled' : '', 'fancy']">
          <div>
            <strong>{{ $t('user.from') }}</strong>: {{ props.row.valid_from | formatDate }}</div>
          <div>
            <strong>{{ $t('user.until') }}</strong>: {{ props.row.valid_until | formatDate }}
            <span :class="['pficon', isExpired(props.row.valid_until) ? 'pficon-error-circle-o': 'pficon-ok']"></span>
          </div>
        </td>
        <td>
          <user-action details="false" :obj="props.row" :update="getAll"></user-action>
        </td>
      </template>
    </vue-good-table>
  </div>
</template>

<script>
import UserService from "../services/user";
import StorageService from "../services/storage";
import HotspotService from "../services/hotspot";
import UtilService from "../services/util";

import UserAction from "../directives/UserAction.vue";

export default {
  name: "Users",
  mixins: [UserService, StorageService, UtilService, HotspotService],
  components: {
    userAction: UserAction
  },
  data() {
    return {
      msg: this.$i18n.t("menu.users"),
      isLoading: true,
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
      hotspotPerPage: this.get("users_per_page") || 25,
      user: this.get("loggedUser") || null
    };
  },
  mounted() {
    if (this.$route.params.hotspotId !== undefined) {
      this.hotspotSearchId = this.$route.params.hotspotId;
    }
    // get user list
    var context = this
    this.getAllHotspots(function() {
      context.getAll();
    });
  },
  methods: {
    handlePerPage(evt) {
      this.set("users_per_page", evt.currentPerPage);
    },
    isExpired(date) {
      return new Date().toISOString() > date;
    },
    getAllHotspots(callback) {
      this.hotspotGetAll(
        success => {
          this.hotspots = success.body;
          var hsId = this.get("users_hotspot_id") || this.hotspots[0].id;
          if (
            this.$parent.user.info.type == "customer" ||
            this.$parent.user.info.type == "desk"
          ) {
            hsId = this.$parent.user.info.hotspot_id;
          }
          this.hotspotSearchId = hsId;

          $('[data-toggle="tooltip"]').tooltip();
          this.isLoading = false;

          callback()
        },
        error => {
          console.log(error);
          callback()
        }
      );
    },
    toggleExpire() {
      this.isLoading = true;
      this.set("users_show_expired", !this.hotspotShowExpired);
      this.getAll();
    },
    getAll() {
      this.set(
        "users_hotspot_id",
        this.hotspotSearchId || this.get("users_hotspot_id") || 0
      );

      this.userGetAll(
        this.hotspotSearchId,
        null,
        success => {
          this.rows = [];
          for (var s in success.body) {
            var res = success.body[s];

            if (!this.isExpired(res.valid_until)) {
              this.rows.push(res);
            } else if (
              this.isExpired(res.valid_until) &&
              this.hotspotShowExpired
            ) {
              this.rows.push(res);
            }
          }
          this.isLoading = false;
        },
        error => {
          this.isLoading = false;
          this.rows = [];
          console.log(error);
        }
      );
    },
    exportCSVUsers() {
      var usersRows = JSON.parse(JSON.stringify(this.rows));
      for (var r in usersRows) {
        // get only email verified users
        if (
          usersRows[r].account_type == "email" &&
          !usersRows[r].email_verified
        ) {
          delete usersRows[r];
        }

        // check marketing authorization
        if (usersRows[r] && !usersRows[r].marketing_auth) {
          delete usersRows[r];
        }
      }

      delete this.columns[3];
      delete this.columns[4];
      delete this.columns[5];
      delete this.columns[6];
      delete this.columns[7];

      var csv = this.createCSV(this.columns, usersRows);
      this.downloadCSV(csv.cols, csv.rows, "users");
    }
  }
};
</script>

<style>
.adjust-results {
  right: 45px !important;
  top: 0px;
}
</style>
