<template>
  <div>
    <h2>{{ msg }}</h2>
    <div v-if="isLoading" class="spinner spinner-lg"></div>
    <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller') && !isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <label v-if="!isLoading" class="col-sm-2 control-label" for="textInput-markup">Hotspot</label>
      <div v-if="!isLoading" class="col-sm-4">
        <select v-on:change="getAll(true)" v-model="hotspotSearchId" class="form-control">
          <option v-for="hotspot in hotspots" v-bind:key="hotspot.id" v-bind:value="hotspot.id">
            {{ hotspot.name }} - {{ hotspot.description}}
          </option>
        </select>
      </div>
    </div>
    <div v-if="!isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <label class="col-sm-2 control-label" for="textInput-markup">{{$t('session.user')}}</label>
      <div class="col-sm-4">
        <select v-on:change="getAll(true)" v-model="hotspotUserId" class="form-control">
          <option value="0">-</option>
          <option v-for="user in users" v-bind:key="user.id" v-bind:value="user.id">
            {{ user.name }}
          </option>
        </select>
      </div>
    </div>
    <div v-if="!isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <div class="col-sm-3">
        <a class="link" @click="clearFilters()">{{$t('session.clear_filters')}}</a>
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
    <div v-if="!isLoading">
      <form v-on:submit.prevent="searchFn($event)">
        <input class="form-control input-lg search-table-input" type="text" :placeholder="tableLangsTexts.globalSearchPlaceholder">
      </form>
    </div>
    <div v-if="!isLoading && isLoadingTable" class="spinner spinner-lg spinner-adjust"></div>
    <vue-good-table v-if="!isLoadingTable && !isLoading" @perPageChanged="handlePerPage" :customRowsPerPageDropdown="[25,50,100]" :perPage="hotspotPerPage"
      :columns="columns" :rows="rows" :lineNumbers="false" :defaultSortBy="{field: 'name', type: 'asc'}" :globalSearch="true"
      :paginate="false" styleClass="table" :nextText="tableLangsTexts.nextText" :prevText="tableLangsTexts.prevText" :rowsPerPageText="tableLangsTexts.rowsPerPageText"
      :ofText="tableLangsTexts.ofText">
      <template slot="table-row" slot-scope="props">
        <td class="fancy">
          <a :href="'#/hotspots/'+ props.row.hotspot_id">{{ props.row.hotspot.name +' - '+ props.row.hotspot.description || '-' }}</a>
        </td>
        <td class="fancy">
          {{ props.row.user.name || '-' }}
        </td>
        <td class="fancy">{{ props.row.mac_address }}</td>
        <td>
          <strong>{{ props.row.ip_address }}</strong>
        </td>
        <td class="fancy">{{ props.row.created || '-' | formatDate}}</td>
      </template>
    </vue-good-table>
    <div v-if="!isLoadingTable && !isLoading" class="right paginator">
      <span class="page-count">
        <b>{{hotspotPage}}</b> {{tableLangsTexts.ofText}} {{total / hotspotPerPage | adjustPage}} (
        <b>{{hotspotPerPage}}</b> {{tableLangsTexts.rowsPerPageText}})</span>
      <button :disabled="availablePrevPage()" @click="prevPage()" class="btn btn-default">{{tableLangsTexts.prevText}}</button>
      <button :disabled="availableNextPage()" @click="nextPage()" class="btn btn-default">{{tableLangsTexts.nextText}}</button>
    </div>
  </div>
</template>
<script>
import DeviceService from "../services/device";
import StorageService from "../services/storage";
import HotspotService from "../services/hotspot";
import UtilService from "../services/util";
import UserService from "../services/user";

import UnitAtion from "../directives/UnitAction";

export default {
  name: "Unit",
  mixins: [
    DeviceService,
    StorageService,
    UtilService,
    HotspotService,
    UserService
  ],
  components: {
    unitAction: UnitAtion
  },
  data() {
    return {
      msg: this.$i18n.t("menu.devices"),
      isLoading: true,
      isLoadingTable: true,
      columns: [
        {
          label: this.$i18n.t("device.hotspot"),
          field: "hotspot_id",
          filterable: true,
          sortable: false
        },
        {
          label: this.$i18n.t("device.user"),
          field: "user_id",
          filterable: true,
          sortable: false
        },
        {
          label: this.$i18n.t("device.mac_address"),
          field: "mac_address",
          filterable: true
        },
        {
          label: this.$i18n.t("device.ip_address"),
          field: "ip_address",
          filterable: true
        },
        {
          label: this.$i18n.t("device.created"),
          field: "created",
          filterable: true
        }
      ],
      rows: [],
      tableLangsTexts: this.tableLangs(),
      hotspots: [],
      hotspotSearchId: 0,
      hotspotUserId: this.get("devices_user_id") || 0,
      hotspotPerPage: 25,
      hotspotPage: 1,
      total: 0,
      user: this.get("loggedUser") || null,
      users: [],
      searchString: ""
    };
  },
  mounted() {
    if (this.$route.params.hotspotId !== undefined) {
      this.hotspotSearchId = this.$route.params.hotspotId;
    }
    // get unit list
    var context = this;
    this.getAllHotspots(function() {
      context.getAllUsers(function() {
        context.getAll();
      });
    });
  },
  methods: {
    handlePerPage(evt) {
      this.set("devices_per_page", evt.currentPerPage);
    },
    searchFn(evt) {
      this.searchString = evt.srcElement[0].value;
      this.getAll(true);
    },
    getAllHotspots(callback) {
      this.hotspotGetAll(
        null,
        null,
        null,
        success => {
          this.hotspots = success.body.data;

          var hsId = this.get("devices_hotspot_id") || this.hotspots[0].id;
          if (
            this.$parent.user.info.type == "customer" ||
            this.$parent.user.info.type == "desk"
          ) {
            hsId = this.$parent.user.info.hotspot_id;
          }
          this.hotspotSearchId = hsId;

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
    getAll(reset) {
      this.isLoadingTable = true;
      if (reset) {
        this.hotspotPage = 1;
        this.total = 0;
      }

      this.set(
        "devices_hotspot_id",
        this.hotspotSearchId || this.get("devices_hotspot_id") || 0
      );
      this.set(
        "devices_user_id",
        this.hotspotUserId || this.get("devices_user_id") || 0
      );

      this.rows = [];
      this.deviceGetAll(
        this.hotspotSearchId,
        this.hotspotUserId,
        this.hotspotPage,
        this.hotspotPerPage,
        encodeURIComponent(this.searchString),
        success => {
          this.rows = success.body.data;
          this.total = success.body.total;
          this.isLoading = false;
          this.isLoadingTable = false;
        },
        error => {
          this.isLoading = false;
          this.isLoadingTable = false;
          this.rows = [];
          console.error(error);
        }
      );
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
    clearFilters() {
      this.hotspotUserId = 0;
      this.set("devices_user_id", 0);
      this.getAll();
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
<style scoped>
</style>