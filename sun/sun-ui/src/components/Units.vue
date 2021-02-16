<template>
  <div>
    <h2>{{ msg }}</h2>
    <div v-if="isLoading" class="spinner spinner-lg"></div>
    <div
      v-if="
        user.account_type == 'admin' ||
        (user.account_type == 'reseller' && !isLoading)
      "
      class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12"
    >
      <label
        v-if="!isLoading"
        class="col-sm-2 control-label"
        for="textInput-markup"
        >Hotspot</label
      >
      <div v-if="!isLoading" class="col-sm-4">
        <select
          v-on:change="getAll(true)"
          v-model="hotspotSearchId"
          class="form-control"
        >
          <option
            v-for="hotspot in hotspots"
            v-bind:key="hotspot.id"
            v-bind:value="hotspot.id"
          >
            {{ hotspot.name }} - {{ hotspot.description }}
          </option>
        </select>
      </div>
    </div>
    <div
      v-if="!isLoading"
      class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12"
    >
      <div class="col-sm-3">
        <button class="btn btn-primary" @click="getAll()">
          {{ $t("session.refresh") }}
        </button>
      </div>
    </div>
    <div
      v-if="!isLoading"
      class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12"
    >
      <div class="result-list">
        {{ total }} {{ total == 1 ? $t("result") : $t("results") }}
      </div>
    </div>
    <div v-if="!isLoading">
      <form v-on:submit.prevent="searchFn($event)">
        <input
          class="form-control input-lg search-table-input"
          type="text"
          :placeholder="tableLangsTexts.globalSearchPlaceholder"
        />
      </form>
    </div>
    <div
      v-if="!isLoading && isLoadingTable"
      class="spinner spinner-lg spinner-adjust"
    ></div>
    <vue-good-table
      v-if="!isLoadingTable && !isLoading"
      @perPageChanged="handlePerPage"
      :customRowsPerPageDropdown="[25, 50, 100]"
      :perPage="hotspotPerPage"
      :columns="columns"
      :rows="rows"
      :lineNumbers="false"
      :defaultSortBy="{ field: 'name', type: 'asc' }"
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
        <td class="fancy">
          <a :href="'#/units/' + props.row.id">{{ props.row.name || "-" }}</a>
        </td>
        <td class="fancy">{{ props.row.description }}</td>
        <td>
          <strong>{{ props.row.mac_address }}</strong>
        </td>
        <td class="fancy">{{ props.row.uuid || "-" }}</td>
        <td class="fancy">
          <div
            v-if="unitStates[props.row.id].isLoading"
            class="spinner spinner-sm"
          ></div>
          <span
            class="font"
            v-if="
              !unitStates[props.row.id].isLoading &&
              unitStates[props.row.id].length > 0
            "
            ><span class="fa fa-list"></span>
            {{ unitStates[props.row.id].length }} {{ $t("unit.active") }}</span
          >
          <span
            v-if="
              !unitStates[props.row.id].isLoading &&
              unitStates[props.row.id].length == 0
            "
            class="font"
            ><span class="fa fa-list"></span>{{ $t("unit.inactive") }}</span
          >
        </td>
        <td>
          <unit-action
            details="true"
            :obj="props.row"
            :update="getAll"
          ></unit-action>
        </td>
      </template>
    </vue-good-table>
    <div v-if="!isLoadingTable && !isLoading" class="right paginator">
      <span class="page-count">
        <b>{{ hotspotPage }}</b> {{ tableLangsTexts.ofText }}
        {{ (total / hotspotPerPage) | adjustPage }} (
        <b>{{ hotspotPerPage }}</b> {{ tableLangsTexts.rowsPerPageText }})</span
      >
      <button
        :disabled="availablePrevPage()"
        @click="prevPage()"
        class="btn btn-default"
      >
        {{ tableLangsTexts.prevText }}
      </button>
      <button
        :disabled="availableNextPage()"
        @click="nextPage()"
        class="btn btn-default"
      >
        {{ tableLangsTexts.nextText }}
      </button>
    </div>
  </div>
</template>
<script>
import UnitService from "../services/unit";
import StorageService from "../services/storage";
import HotspotService from "../services/hotspot";
import UtilService from "../services/util";
import moment from "moment";
import UnitAtion from "../directives/UnitAction";

export default {
  name: "Unit",
  mixins: [UnitService, StorageService, UtilService, HotspotService],
  components: {
    unitAction: UnitAtion,
  },
  data() {
    return {
      msg: this.$i18n.t("menu.units"),
      isLoading: true,
      isLoadingTable: true,
      columns: [
        {
          label: this.$i18n.t("unit.name"),
          field: "name",
          filterable: true,
        },
        {
          label: this.$i18n.t("unit.description"),
          field: "description",
          filterable: true,
        },
        {
          label: this.$i18n.t("unit.mac_address"),
          field: "mac_address",
          filterable: true,
        },
        {
          label: this.$i18n.t("unit.uuid"),
          field: "uuid",
          filterable: true,
        },
        {
          label: this.$i18n.t("unit.status"),
          field: "status",
          filterable: true,
        },
        {
          label: this.$i18n.t("action"),
          field: "",
          sortable: false,
        },
      ],
      rows: [],
      tableLangsTexts: this.tableLangs(),
      hotspots: [],
      hotspotSearchId: 0,
      hotspotPerPage: 25,
      hotspotPage: 1,
      total: 0,
      user: this.get("loggedUser") || null,
      searchString: "",
      unitStates: {},
    };
  },
  mounted() {
    if (this.$route.params.hotspotId !== undefined) {
      this.hotspotSearchId = this.$route.params.hotspotId;
    } else {
      var hsId = this.get("selected_hotspot_id");
      this.hotspotSearchId = hsId;
    }
    // get unit list
    var context = this;
    this.getAllHotspots(function () {
      context.getAll();
    });
  },
  methods: {
    handlePerPage(evt) {
      this.set("units_per_page", evt.currentPerPage);
    },
    searchFn(evt) {
      var elem = evt.srcElement || evt.target;
      this.searchString = elem[0].value;
      this.getAll(true);
    },
    getAllHotspots(callback) {
      this.hotspotGetAll(
        null,
        null,
        null,
        (success) => {
          this.hotspots = success.body.data;
          var hsId = this.get("selected_hotspot_id") || this.hotspots[0].id;
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
        (error) => {
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
        "selected_hotspot_id",
        this.hotspotSearchId || this.get("selected_hotspot_id") || 0
      );

      this.unitGetAll(
        this.hotspotSearchId,
        this.hotspotPage,
        this.hotspotPerPage,
        encodeURIComponent(this.searchString),
        (success) => {
          this.rows = success.body.data;
          this.total = success.body.total;
          this.isLoading = false;
          this.isLoadingTable = false;
          var m = moment;
          var context = this;
          for (var i = 0; i < success.body.data.length; i++) {
            var unit = success.body.data[i];
            context.unitStates[unit.id] = {
              isLoading: true,
              length: 0,
            };
            this.unitGetStatus(
              1,
              25,
              unit.id,
              m().subtract(30, "minutes").toISOString(),
              m().toISOString(),
              (success) => {
                const url = new URL(success.url);
                var parsed = url.pathname;
                var urlParams = new URLSearchParams(parsed);
                var unitID = urlParams.get("unit");
                context.unitStates[unitID].isLoading = false;
                context.unitStates[unitID].length = success.body.data.length;
                context.$forceUpdate();
              },
              (error) => {
                var url = new URL(error.url);
                var parsed = url.pathname;
                var urlParams = new URLSearchParams(parsed);
                var unitID = urlParams.get("unit");
                context.unitStates[unitID].isLoading = false;
                context.unitStates[unitID].length = 0;
                context.$forceUpdate();
              }
            );
          }
        },
        (error) => {
          this.isLoading = false;
          this.isLoadingTable = false;
          this.rows = [];
          console.error(error);
        }
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
    },
  },
};
</script>
<style scoped>
</style>