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
      <label class="col-sm-2 control-label" for="textInput-markup">{{$t('session.user')}}</label>
      <div class="col-sm-4">
        <select v-on:change="getAll()" v-model="hotspotUserId" class="form-control">
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
    <vue-good-table v-if="!isLoading" @perPageChanged="handlePerPage" :customRowsPerPageDropdown="[25,50,100]" :perPage="hotspotPerPage"
      :columns="columns" :rows="rows" :lineNumbers="false" :defaultSortBy="{field: 'name', type: 'asc'}" :globalSearch="true"
      :paginate="true" styleClass="table" :nextText="tableLangsTexts.nextText" :prevText="tableLangsTexts.prevText" :rowsPerPageText="tableLangsTexts.rowsPerPageText"
      :globalSearchFn="searchFn" :globalSearchPlaceholder="tableLangsTexts.globalSearchPlaceholder" :ofText="tableLangsTexts.ofText">
      <template slot="table-row" slot-scope="props">
        <td class="fancy">
          <a :href="'#/hotspots/'+ props.row.hotspot_id">{{ extractHotspot(props.row.hotspot_id) && extractHotspot(props.row.hotspot_id).name || '-' }}</a>
        </td>
        <td class="fancy">
          {{ extractUser(props.row.user_id) && extractUser(props.row.user_id).name || '-' }}
        </td>
        <td class="fancy">{{ props.row.mac_address }}</td>
        <td>
          <strong>{{ props.row.ip_address }}</strong>
        </td>
        <td class="fancy">{{ props.row.created || '-' | formatDate}}</td>
      </template>
    </vue-good-table>
  </div>
</template>
<script>
  import DeviceService from "../services/device";
  import StorageService from "../services/storage";
  import HotspotService from "../services/hotspot";
  import UtilService from "../services/util";
  import UserService from '../services/user';

  import UnitAtion from "../directives/UnitAction";

  export default {
    name: "Unit",
    mixins: [DeviceService, StorageService, UtilService, HotspotService, UserService],
    components: {
      unitAction: UnitAtion
    },
    data() {
      return {
        msg: this.$i18n.t("menu.devices"),
        isLoading: true,
        columns: [{
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
          },
        ],
        rows: [],
        tableLangsTexts: this.tableLangs(),
        hotspots: [],
        hotspotSearchId: this.get('devices_hotspot_id') || 0,
        hotspotUserId: this.get('devices_user_id') || 0,
        hotspotPerPage: this.get('devices_per_page') || 25,
        user: this.get("loggedUser") || null,
        users: [],
      };
    },
    mounted() {
      if (this.$route.params.hotspotId !== undefined) {
        this.hotspotSearchId = this.$route.params.hotspotId;
      }
      // get unit list
      this.getAll();
      this.getAllHotspots();
      this.getAllUsers()
    },
    methods: {
      handlePerPage(evt) {
        this.set('devices_per_page', evt.currentPerPage)
      },
      searchFn(row, col, cellValue, searchTerm) {
        var value = cellValue.toString().toLowerCase()
        if (col.field == 'hotspot_id') {
          value = this.extractHotspot(cellValue).name.toLowerCase()
        }

        if (col.field == 'user_id') {
          value = this.extractUser(cellValue).name.toLowerCase()
        }
        return value.includes(searchTerm.toLowerCase())
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
        this.set('devices_hotspot_id', this.hotspotSearchId || this.get('devices_hotspot_id') || 0)
        this.set('devices_user_id', this.hotspotUserId || this.get('devices_user_id') || 0)

        // preload users and units
        this.getAllUsers()

        this.rows = []
        this.deviceGetAll(
          this.hotspotSearchId,
          this.hotspotUserId,
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
      },
      getAllUsers() {
        this.userGetAll(this.hotspotSearchId, success => {
          this.users = success.body
        }, error => {
          console.log(error)
          this.users = {}
        })
      },
      extractHotspot(id) {
        for (var h in this.hotspots) {
          if (this.hotspots[h].id == id) {
            return this.hotspots[h]
            break;
          }
        }
      },
      extractUser(id) {
        for (var u in this.users) {
          if (this.users[u].id == id) {
            return this.users[u]
            break;
          }
        }
      },
      clearFilters() {
        this.hotspotSearchId = 0
        this.hotspotUserId = 0
        this.set('sessions_hotspot_id', 0)
        this.set('sessions_user_id', 0)
        this.getAll()
      },
    }
  };

</script>
<style scoped>


</style>
