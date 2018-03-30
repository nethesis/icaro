<template>
  <div>
    <h2>{{ msg }}</h2>
    <div v-if="isLoading" class="spinner spinner-lg"></div>
    <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller') && !isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <label v-if="!isLoading" class="col-sm-2 control-label" for="textInput-markup">Hotspot</label>
      <div v-if="!isLoading" class="col-sm-4">
        <select v-on:change="getAll()" v-model="hotspotSearchId" class="form-control">
          <option value="0">-</option>
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
    <vue-good-table v-if="!isLoading" @perPageChanged="handlePerPage" :customRowsPerPageDropdown="[25,50,100]" :perPage="hotspotPerPage" :columns="columns" :rows="rows" :lineNumbers="false" :defaultSortBy="{field: 'name', type: 'asc'}"
      :globalSearch="true" :paginate="true" styleClass="table" :nextText="tableLangsTexts.nextText" :prevText="tableLangsTexts.prevText"
      :rowsPerPageText="tableLangsTexts.rowsPerPageText" :globalSearchPlaceholder="tableLangsTexts.globalSearchPlaceholder"
      :ofText="tableLangsTexts.ofText">
      <template slot="table-row" slot-scope="props">
        <td class="fancy">
          <a :href="'#/units/'+ props.row.id">{{ props.row.name || '-' }}</a>
        </td>
        <td class="fancy">{{ props.row.description }}</td>
        <td>
          <strong>{{ props.row.mac_address }}</strong>
        </td>
        <td class="fancy">{{ props.row.uuid || '-' }}</td>
        <td>
          <unit-action details="true" :obj="props.row" :update="getAll"></unit-action>
        </td>
      </template>
    </vue-good-table>
  </div>
</template>
<script>
  import UnitService from "../services/unit";
  import StorageService from "../services/storage";
  import HotspotService from "../services/hotspot";
  import UtilService from "../services/util";

  import UnitAtion from "../directives/UnitAction";

  export default {
    name: "Unit",
    mixins: [UnitService, StorageService, UtilService, HotspotService],
    components: {
      unitAction: UnitAtion
    },
    data() {
      return {
        msg: this.$i18n.t("menu.units"),
        isLoading: true,
        columns: [{
            label: this.$i18n.t("unit.name"),
            field: "name",
            filterable: true
          },
          {
            label: this.$i18n.t("unit.description"),
            field: "description",
            filterable: true
          },
          {
            label: this.$i18n.t("unit.mac_address"),
            field: "mac_address",
            filterable: true
          },
          {
            label: this.$i18n.t("unit.uuid"),
            field: "uuid",
            filterable: true
          },
          {
            label: this.$i18n.t('action'),
            field: '',
            sortable: false
          },
        ],
        rows: [],
        tableLangsTexts: this.tableLangs(),
        hotspots: [],
        hotspotSearchId: this.get('units_hotspot_id') || 0,
        hotspotPerPage: this.get('units_per_page') || 25,
        user: this.get("loggedUser") || null
      };
    },
    mounted() {
      if (this.$route.params.hotspotId !== undefined) {
        this.hotspotSearchId = this.$route.params.hotspotId;
      }
      // get unit list
      this.getAll();
      this.getAllHotspots();
    },
    methods: {
      handlePerPage(evt) {
        this.set('units_per_page', evt.currentPerPage)
      },
      getAllHotspots() {
        this.hotspotGetAll(
          success => {
            this.hotspots = success.body;
            $('[data-toggle="tooltip"]').tooltip();
            this.isLoading = false;
          },
          error => {
            console.log(error);
          }
        );
      },
      getAll() {
        this.set('units_hotspot_id', this.hotspotSearchId || this.get('units_hotspot_id') || 0)

        this.unitGetAll(
          this.hotspotSearchId,
          success => {
            this.rows = success.body;
            this.isLoading = false;
          },
          error => {
            this.isLoading = false;
            this.rows = [];
            console.log(error);
          }
        );
      }
    }
  };

</script>
<style scoped>


</style>
