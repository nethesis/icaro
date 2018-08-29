<template>
  <div>
    <h2>{{ msg }}</h2>
    <div v-if="isLoading" class="spinner spinner-lg"></div>
    <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller') && !isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <label v-if="!isLoading" class="col-sm-2 control-label" for="textInput-markup">Hotspot</label>
      <div v-if="!isLoading" class="col-sm-4">
        <select v-on:change="getAll(false, true)" v-model="hotspotSearchId" class="form-control">
          <option v-for="hotspot in hotspots" v-bind:key="hotspot.id" v-bind:value="hotspot.id">
            {{ hotspot.name }} - {{ hotspot.description}}
          </option>
        </select>
      </div>
    </div>
    <div v-if="!isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <label class="col-sm-2 control-label" for="textInput-markup">{{$t('session.user')}}</label>
      <div class="col-sm-4">
        <select v-on:change="getAll(false, true)" v-model="hotspotUserId" class="form-control">
          <option value="0">-</option>
          <option v-for="user in users" v-bind:key="user.id" v-bind:value="user.id">
            {{ user.name }}
          </option>
        </select>
      </div>
    </div>
    <div v-if="!isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <label class="col-sm-2 control-label" for="textInput-markup">{{$t('session.unit')}}</label>
      <div class="col-sm-4">
        <select v-on:change="getAll(false, true)" v-model="hotspotUnitId" class="form-control">
          <option value="0">-</option>
          <option v-for="unit in units" v-bind:key="unit.id" v-bind:value="unit.id">
            {{ unit.name }} - {{ unit.description }}
          </option>
        </select>
      </div>
    </div>
    <div v-if="!isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <label class="col-sm-2 control-label" for="textInput-markup">{{$t('session.from')}}</label>
      <div class="col-sm-4">
        <datepicker :format="dateFormatter" @input="getAll()" v-model="hotspotDateFrom" :language="locale"></datepicker>
      </div>
    </div>
    <div v-if="!isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <label class="col-sm-2 control-label" for="textInput-markup">{{$t('session.to')}}</label>
      <div class="col-sm-4">
        <datepicker :format="dateFormatter" @input="getAll()" v-model="hotspotDateTo" :language="locale"></datepicker>
      </div>
    </div>
    <div v-if="!isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <div class="col-sm-3">
        <a class="link" @click="clearFilters()">{{$t('session.clear_filters')}}</a>
      </div>
    </div>
    <div v-if="!isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <div class="col-xs-3 col-sm-3 col-md-3 col-lg-3">
        <button class="btn btn-primary" @click="getAll(true)">{{$t('session.refresh')}}</button>
      </div>
    </div>

    <ul v-if="!isLoading" class="nav nav-tabs nav-tabs-pf adjust-tabs" id="myTab" role="tablist">
      <li class="nav-item">
        <a @click="handleTab('active')" class="nav-link active" id="active-tab-parent" data-toggle="tab" href="#active-tab" role="tab"
          aria-controls="active" aria-selected="true">{{$t('session.active')}}</a>
      </li>
      <li class="nav-item">
        <a @click="handleTab('history')" class="nav-link" id="history-tab-parent" data-toggle="tab" href="#history-tab" role="tab"
          aria-controls="history" aria-selected="false">{{$t('session.history')}}</a>
      </li>
    </ul>
    <div class="tab-pane fade" id="active-tab" role="tabpanel" aria-labelledby="active-tab">
      <div v-if="!isLoading && isLoadingTable && activeTab == 'active'" class="spinner spinner-lg spinner-adjust"></div>
      <div v-if="!isLoadingTable && !isLoading && activeTab == 'active'" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
        <div class="col-sm-12">
          <button @click="exportCSVActive()" class="btn btn-primary export-btn">{{$t('session.export_csv')}}</button>
        </div>
        <div class="result-list adjust-results">{{totalActive}} {{totalActive == 1 ? $t('result') : $t('results')}}</div>
      </div>
      <div v-if="!isLoadingTable && !isLoading && activeTab == 'active'">
        <form v-on:submit.prevent="searchFn($event, 'active')">
          <input class="form-control input-lg search-table-input" type="text" :placeholder="tableLangsTexts.globalSearchPlaceholder">
        </form>
      </div>
      <div v-if="!isLoadingTable && !isLoading && activeTab == 'active' && exportError" class="alert alert-danger alert-dismissable alert-export">
        <span class="pficon pficon-error-circle-o"></span>
        <strong>{{$t('session.export_error')}}</strong>. {{$t('session.export_error_details')}}.
      </div>
      <vue-good-table v-if="!isLoadingTable && !isLoading && activeTab == 'active'" @perPageChanged="handlePerPage" :customRowsPerPageDropdown="[25,50,100]"
        :perPage="hotspotPerPage" :columns="columns_active" :rows="rows_active" :lineNumbers="false" :defaultSortBy="{field: 'duration', type: 'asc'}"
        :globalSearch="true" :paginate="false" styleClass="table" :nextText="tableLangsTexts.nextText" :prevText="tableLangsTexts.prevText"
        :rowsPerPageText="tableLangsTexts.rowsPerPageText" :globalSearchPlaceholder="tableLangsTexts.globalSearchPlaceholder"
        :ofText="tableLangsTexts.ofText">
        <template slot="table-row" slot-scope="props">
          <td class="fancy">
            <a :href="'#/units/'+ props.row.unit_id">{{ props.row.unit.name || '-'}}</a>
          </td>
          <td class="fancy">
            <a :href="'#/sessions/'+ props.row.id">{{ props.row.username || '-'}}</a>
          </td>
          <td class="fancy">{{ props.row.bytes_down | byteFormat}}</td>
          <td class="fancy">{{ props.row.bytes_up | byteFormat }}</td>
          <td class="fancy">{{ props.row.duration | secondsInHour }}</td>
          <td class="fancy">{{ props.row.start_time | formatDate }}</td>
          <td class="fancy">{{ props.row.update_time | formatDate }}</td>
          <td>
            <a :href="'#/sessions/'+ props.row.id">
              <span class="fa fa-angle-right details-arrow"></span>
            </a>
          </td>
        </template>
      </vue-good-table>
      <div v-if="!isLoadingTable && !isLoading&& activeTab == 'active'" class="right paginator">
        <span class="page-count">
          <b>{{hotspotPageActive}}</b> {{tableLangsTexts.ofText}} {{totalActive / hotspotPerPage | adjustPage}} (
          <b>{{hotspotPerPage}}</b> {{tableLangsTexts.rowsPerPageText}})</span>
        <button :disabled="availablePrevPageActive()" @click="prevPageActive()" class="btn btn-default">{{tableLangsTexts.prevText}}</button>
        <button :disabled="availableNextPageActive()" @click="nextPageActive()" class="btn btn-default">{{tableLangsTexts.nextText}}</button>
      </div>
    </div>
    <div class="tab-pane fade" id="history-tab" role="tabpanel" aria-labelledby="history-tab">
      <div v-if="!isLoading && isLoadingTable && activeTab == 'history'" class="spinner spinner-lg spinner-adjust"></div>
      <div v-if="!isLoadingTable && !isLoading && activeTab == 'history'" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
        <div class="col-sm-12">
          <button @click="exportCSVHistory()" class="btn btn-primary export-btn">{{$t('session.export_csv')}}</button>
        </div>
        <div class="result-list adjust-results">{{totalHistory}} {{totalHistory == 1 ? $t('result') : $t('results')}}</div>
      </div>
      <div v-if="!isLoadingTable && !isLoading && activeTab == 'history'">
        <form v-on:submit.prevent="searchFn($event, 'history')">
          <input class="form-control input-lg search-table-input" type="text" :placeholder="tableLangsTexts.globalSearchPlaceholder">
        </form>
      </div>
      <div v-if="!isLoadingTable && !isLoading && activeTab == 'history' && exportError" class="alert alert-danger alert-dismissable alert-export">
        <span class="pficon pficon-error-circle-o"></span>
        <strong>{{$t('session.export_error')}}</strong>. {{$t('session.export_error_details')}}.
      </div>
      <vue-good-table v-if="!isLoadingTable && !isLoading && activeTab == 'history'" @perPageChanged="handlePerPage" :customRowsPerPageDropdown="[25,50,100]"
        :perPage="hotspotPerPage" :columns="columns_history" :rows="rows_history" :lineNumbers="false" :defaultSortBy="{field: 'duration', type: 'asc'}"
        :globalSearch="true" :paginate="false" styleClass="table" :nextText="tableLangsTexts.nextText" :prevText="tableLangsTexts.prevText"
        :rowsPerPageText="tableLangsTexts.rowsPerPageText" :globalSearchPlaceholder="tableLangsTexts.globalSearchPlaceholder"
        :ofText="tableLangsTexts.ofText">
        <template slot="table-row" slot-scope="props">
          <td class="fancy">
            <a :href="'#/units/'+ props.row.unit_id">{{ props.row.unit.name || '-'}}</a>
          </td>
          <td class="fancy">
            <a :href="'#/sessions/'+ props.row.id">{{ props.row.username || '-'}}</a>
          </td>
          <td class="fancy">{{ props.row.bytes_down | byteFormat}}</td>
          <td class="fancy">{{ props.row.bytes_up | byteFormat }}</td>
          <td class="fancy">{{ props.row.duration | secondsInHour }}</td>
          <td class="fancy">{{ props.row.start_time | formatDate }}</td>
          <td class="fancy">{{ props.row.stop_time | formatDate }}</td>
          <td>
            <a :href="'#/sessions/'+ props.row.id">
              <span class="fa fa-angle-right details-arrow"></span>
            </a>
          </td>
        </template>
      </vue-good-table>
      <div v-if="!isLoadingTable && !isLoading&& activeTab == 'history'" class="right paginator">
        <span class="page-count">
          <b>{{hotspotPageHistory}}</b> {{tableLangsTexts.ofText}} {{totalHistory / hotspotPerPage | adjustPage}} (
          <b>{{hotspotPerPage}}</b> {{tableLangsTexts.rowsPerPageText}})</span>
        <button :disabled="availablePrevPageHistory()" @click="prevPageHistory()" class="btn btn-default">{{tableLangsTexts.prevText}}</button>
        <button :disabled="availableNextPageHistory()" @click="nextPageHistory()" class="btn btn-default">{{tableLangsTexts.nextText}}</button>
      </div>
    </div>
  </div>
</template>

<script>
import StorageService from "../services/storage";
import UtilService from "../services/util";
import HotspotService from "../services/hotspot";
import SessionService from "../services/session";
import HistoryService from "../services/history";
import UserService from "../services/user";
import UnitService from "../services/unit";

import Datepicker from "vuejs-datepicker";
import moment from "moment";

export default {
  name: "Reports",
  mixins: [
    SessionService,
    StorageService,
    UtilService,
    HotspotService,
    UserService,
    UnitService,
    HistoryService
  ],
  components: {
    Datepicker
  },
  data() {
    var activeTab = this.get("sessions_active_tab") || "active";
    setTimeout(function() {
      window.$("#" + activeTab + "-tab-parent").click();
    }, 500);

    return {
      msg: this.$i18n.t("menu.sessions"),
      isLoading: true,
      isLoadingTable: true,
      locale: this.$root.$options.currentLocale,
      columns_active: [
        {
          label: this.$i18n.t("session.unit"),
          field: "unit_id",
          filterable: true,
          sortable: false
        },
        {
          label: this.$i18n.t("session.user"),
          field: "username",
          filterable: true,
          sortable: false
        },
        {
          label: this.$i18n.t("session.bytes_down"),
          field: "bytes_down",
          type: "number",
          filterable: true
        },
        {
          label: this.$i18n.t("session.bytes_up"),
          field: "bytes_up",
          type: "number",
          filterable: true
        },
        {
          label: this.$i18n.t("session.duration"),
          field: "duration",
          type: "number",
          filterable: true
        },
        {
          label: this.$i18n.t("session.start_time"),
          field: "start_time",
          filterable: true
        },
        {
          label: this.$i18n.t("session.update_time"),
          field: "update_time",
          filterable: true,
          sortable: false
        },
        {
          label: "",
          field: "",
          sortable: false
        }
      ],
      columns_history: [
        {
          label: this.$i18n.t("session.unit"),
          field: "unit_id",
          filterable: true,
          sortable: false
        },
        {
          label: this.$i18n.t("session.user"),
          field: "username",
          filterable: true,
          sortable: false
        },
        {
          label: this.$i18n.t("session.bytes_down"),
          field: "bytes_down",
          type: "number",
          filterable: true
        },
        {
          label: this.$i18n.t("session.bytes_up"),
          field: "bytes_up",
          type: "number",
          filterable: true
        },
        {
          label: this.$i18n.t("session.duration"),
          field: "duration",
          type: "number",
          filterable: true
        },
        {
          label: this.$i18n.t("session.start_time"),
          field: "start_time",
          filterable: true
        },
        {
          label: this.$i18n.t("session.stop_time"),
          field: "stop_time",
          filterable: true,
          sortable: false
        },
        {
          label: "",
          field: "",
          sortable: false
        }
      ],
      rows_active: [],
      rows_history: [],
      activeTab: activeTab,
      tableLangsTexts: this.tableLangs(),
      hotspots: [],
      hotspotSearchId: 0,
      hotspotUserId: this.get("sessions_user_id") || 0,
      hotspotUnitId: this.get("sessions_unit_id") || 0,
      hotspotDateFrom: moment(Date.now() - 12096e5)
        .utc()
        .startOf("day")
        .toISOString(),
      hotspotDateTo: moment(Date.now())
        .endOf("day")
        .utc()
        .toISOString(),
      hotspotPerPage: 25,
      hotspotPageActive: 1,
      totalActive: 0,
      hotspotPageHistory: 1,
      totalHistory: 0,
      user: this.get("loggedUser") || null,
      users: [],
      units: [],
      searchStringActive: "",
      searchStringHistory: "",
      exportError: false
    };
  },
  mounted() {
    if (this.$route.params.hotspotId !== undefined) {
      this.hotspotSearchId = this.$route.params.hotspotId;
    }
    // get session list
    var context = this;
    this.getAllHotspots(function() {
      context.getAllUsers(function() {
        context.getAllUnits(function() {
          context.getAll(true);
        });
      });
    });
  },
  methods: {
    handleTab(tab) {
      this.exportError = false;
      this.activeTab = tab;
      this.getAll();
    },
    handlePerPage(evt) {
      this.set("sessions_per_page", evt.currentPerPage);
    },
    searchFn(evt, type) {
      var elem = evt.srcElement || evt.target;
      if (type == "active") {
        this.searchStringActive = elem[0].value;
      }

      if (type == "history") {
        this.searchStringHistory = elem[0].value;
      }
      this.getAll(true);
    },
    dateFormatter(date) {
      return moment(date).format("DD MMMM YYYY");
    },
    getAllHotspots(callback) {
      this.hotspotGetAll(
        null,
        null,
        null,
        success => {
          this.hotspots = success.body.data;
          var hsId = this.get("sessions_hotspot_id") || this.hotspots[0].id;
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
          this.isLoading = false;
        }
      );
    },
    getAll(refresh, reset) {
      this.isLoadingTable = true;
      // save to storage
      this.set(
        "sessions_active_tab",
        this.activeTab || this.get("sessions_active_tab") || 0
      );
      this.set(
        "sessions_hotspot_id",
        this.hotspotSearchId || this.get("sessions_hotspot_id") || 0
      );
      this.set(
        "sessions_user_id",
        this.hotspotUserId || this.get("sessions_user_id") || 0
      );
      this.set(
        "sessions_unit_id",
        this.hotspotUnitId || this.get("sessions_unit_id") || 0
      );

      if (refresh) {
        if (
          this.hotspotDateTo ===
          moment()
            .utc()
            .endOf("day")
            .toISOString()
        ) {
          this.hotspotDateTo = moment(Date.now())
            .endOf("day")
            .utc()
            .toISOString();
        }
      }

      if (reset) {
        this.hotspotPageActive = 1;
        this.hotspotPageHistory = 1;
        this.totalActive = 0;
        this.totalHistory = 0;
      }

      // get all sessions
      if (this.activeTab == "active") {
        this.sessionGetAll(
          this.hotspotSearchId,
          this.hotspotUserId,
          this.hotspotUnitId,
          new Date(this.hotspotDateFrom).toISOString(),
          new Date(this.hotspotDateTo).toISOString(),
          this.hotspotPageActive,
          this.hotspotPerPage,
          encodeURIComponent(this.searchStringActive),
          success => {
            this.rows_active = success.body.data;
            this.totalActive = success.body.total;
            this.isLoadingTable = false;
            this.isLoading = false;
            var activeTab = this.get("sessions_active_tab") || "active";
            setTimeout(function() {
              window.$("#" + activeTab + "-tab-parent").click();
            }, 250);
          },
          error => {
            this.isLoadingTable = false;
            this.isLoading = false;
            this.rows_active = [];
            console.error(error);
          }
        );
      } else {
        this.historiesGetAll(
          this.hotspotSearchId,
          this.hotspotUserId,
          this.hotspotUnitId,
          new Date(this.hotspotDateFrom).toISOString(),
          new Date(this.hotspotDateTo).toISOString(),
          this.hotspotPageHistory,
          this.hotspotPerPage,
          encodeURIComponent(this.searchStringHistory),
          success => {
            this.rows_history = success.body.data;
            this.totalHistory = success.body.total;
            this.isLoadingTable = false;
            this.isLoading = false;
            var activeTab = this.get("sessions_active_tab") || "active";
            setTimeout(function() {
              window.$("#" + activeTab + "-tab-parent").click();
            }, 250);
          },
          error => {
            this.isLoadingTable = false;
            this.isLoading = false;
            this.rows_history = [];
            console.error(error);
          }
        );
      }
    },
    getAllUsers(callback) {
      this.userGetAll(
        this.hotspotSearchId,
        null,
        false,
        null,
        null,
        null,
        success => {
          this.users = success.body.data;
          callback();
        },
        error => {
          console.error(error);
          this.users = {};
          callback();
        }
      );
    },
    getAllUnits(callback) {
      this.unitGetAll(
        this.hotspotSearchId,
        null,
        null,
        null,
        success => {
          this.units = success.body.data;
          callback();
        },
        error => {
          console.error(error);
          this.units = {};
          callback();
        }
      );
    },
    clearFilters() {
      this.hotspotUserId = 0;
      this.hotspotUnitId = 0;
      this.hotspotDateFrom = moment(Date.now() - 12096e5)
        .utc()
        .startOf("day")
        .toISOString();
      this.hotspotDateTo = moment(Date.now())
        .endOf("day")
        .utc()
        .toISOString();
      this.set("sessions_user_id", 0);
      this.set("sessions_unit_id", 0);
      this.getAll();
    },
    exportCSVActive() {
      this.isLoadingTable = true;
      this.sessionGetAll(
        this.hotspotSearchId,
        this.hotspotUserId,
        this.hotspotUnitId,
        new Date(this.hotspotDateFrom).toISOString(),
        new Date(this.hotspotDateTo).toISOString(),
        null,
        null,
        encodeURIComponent(this.searchStringActive),
        success => {
          var data_export = success.body.data;
          if (data_export.length < 5000) {
            var newRows = JSON.parse(JSON.stringify(data_export));
            for (var r in newRows) {
              newRows[r].unit_id = newRows[r].unit.description || "-";
              newRows[r].user_id = newRows[r].username || "-";
              newRows[r].bytes_up = this.$options.filters["byteFormat"](
                newRows[r].bytes_up
              );
              newRows[r].bytes_down = this.$options.filters["byteFormat"](
                newRows[r].bytes_down
              );
              newRows[r].duration = this.$options.filters["secondsInHour"](
                newRows[r].duration
              );
              newRows[r].start_time = this.$options.filters["formatDate"](
                newRows[r].start_time
              );
              newRows[r].update_time = this.$options.filters["formatDate"](
                newRows[r].update_time
              );
            }
            var csv = this.createCSV(this.columns_active, newRows);
            this.isLoadingTable = false;
            this.downloadCSV(csv.cols, csv.rows, "sessions_active");
          } else {
            this.isLoadingTable = false;
            this.exportError = true;
          }
        },
        error => {
          this.isLoadingTable = false;
          this.isLoading = false;
          this.rows_active = [];
          console.error(error);
        }
      );
    },
    exportCSVHistory() {
      this.isLoadingTable = true;
      this.historiesGetAll(
        this.hotspotSearchId,
        this.hotspotUserId,
        this.hotspotUnitId,
        new Date(this.hotspotDateFrom).toISOString(),
        new Date(this.hotspotDateTo).toISOString(),
        null,
        null,
        encodeURIComponent(this.searchStringHistory),
        success => {
          var data_export = success.body.data;
          if (data_export.length < 5000) {
            var newRows = JSON.parse(JSON.stringify(data_export));
            for (var r in newRows) {
              newRows[r].unit_id = newRows[r].unit.description || "-";
              newRows[r].user_id = newRows[r].username || "-";
              newRows[r].bytes_up = this.$options.filters["byteFormat"](
                newRows[r].bytes_up
              );
              newRows[r].bytes_down = this.$options.filters["byteFormat"](
                newRows[r].bytes_down
              );
              newRows[r].duration = this.$options.filters["secondsInHour"](
                newRows[r].duration
              );
              newRows[r].start_time = this.$options.filters["formatDate"](
                newRows[r].start_time
              );
              newRows[r].stop_time = this.$options.filters["formatDate"](
                newRows[r].stop_time
              );
            }
            var csv = this.createCSV(this.columns_history, newRows);
            this.isLoadingTable = false;
            this.downloadCSV(csv.cols, csv.rows, "sessions_history");
          } else {
            this.isLoadingTable = false;
            this.exportError = true;
          }
        },
        error => {
          this.isLoadingTable = false;
          this.isLoading = false;
          this.rows_history = [];
          console.error(error);
        }
      );
    },
    prevPageActive() {
      this.isLoading = true;
      if (this.hotspotPageActive != 1) {
        this.hotspotPageActive--;
      }
      this.getAll();
    },
    nextPageActive() {
      this.isLoading = true;
      if (
        this.hotspotPageActive !=
        Math.ceil(this.totalActive / this.hotspotPerPage)
      ) {
        this.hotspotPageActive++;
      }
      this.getAll();
    },
    availablePrevPageActive() {
      return this.hotspotPageActive == 1;
    },
    availableNextPageActive() {
      return (
        this.hotspotPageActive ==
        Math.ceil(this.totalActive / this.hotspotPerPage)
      );
    },
    prevPageHistory() {
      this.isLoading = true;
      if (this.hotspotPageHistory != 1) {
        this.hotspotPageHistory--;
      }
      this.getAll();
    },
    nextPageHistory() {
      this.isLoading = true;
      if (
        this.hotspotPageHistory !=
        Math.ceil(this.totalHistory / this.hotspotPerPage)
      ) {
        this.hotspotPageHistory++;
      }
      this.getAll();
    },
    availablePrevPageHistory() {
      return this.hotspotPageHistory == 1;
    },
    availableNextPageHistory() {
      return (
        this.hotspotPageHistory ==
        Math.ceil(this.totalHistory / this.hotspotPerPage)
      );
    }
  }
};
</script>
<style scoped>
.adjust-tabs {
  margin-left: 15px !important;
  margin-right: 15px !important;
  border-bottom: 1px solid #cccccc !important;
}

.adjust-results {
  right: 45px !important;
  top: 0px;
}
</style>