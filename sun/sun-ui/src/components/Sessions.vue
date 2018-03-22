<template>
  <div>
    <h2>{{ $t("sessions") }}</h2>
    <div v-if="isLoading" class="spinner spinner-lg"></div>
    <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller') && !isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <label class="col-sm-2 control-label" for="textInput-markup">Hotspot</label>
      <div class="col-sm-4">
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
      <label class="col-sm-2 control-label" for="textInput-markup">{{$t('session.unit')}}</label>
      <div class="col-sm-4">
        <select v-on:change="getAll()" v-model="hotspotUnitId" class="form-control">
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
      <div class="col-sm-12">
        <button @click="exportCSV()" class="btn btn-primary export-btn">{{$t('session.export_csv')}}</button>
      </div>
    </div>
    <vue-good-table v-if="!isLoading" :perPage="25" :columns="columns" :rows="rows" :lineNumbers="false" :defaultSortBy="{field: 'duration', type: 'asc'}"
      :globalSearch="true" :globalSearchFn="searchFn" :paginate="true" styleClass="table" :nextText="tableLangsTexts.nextText"
      :prevText="tableLangsTexts.prevText" :rowsPerPageText="tableLangsTexts.rowsPerPageText" :globalSearchPlaceholder="tableLangsTexts.globalSearchPlaceholder"
      :ofText="tableLangsTexts.ofText">
      <template slot="table-row" slot-scope="props">
        <td class="fancy">
          <a :href="'#/sessions/'+ props.row.id">{{ extractUnit(props.row.unit_id) && extractUnit(props.row.unit_id).description || '-'}}</a>
        </td>
        <td class="fancy">
          <a :href="'#/sessions/'+ props.row.id">{{ extractUser(props.row.user_id) && extractUser(props.row.user_id).name || '-'}}</a>
        </td>
        <td class="fancy">{{ props.row.bytes_up | byteFormat}}</td>
        <td class="fancy">{{ props.row.bytes_down | byteFormat }}</td>
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
  </div>
</template>

<script>
  import StorageService from '../services/storage';
  import UtilService from '../services/util';
  import HotspotService from '../services/hotspot';
  import SessionService from '../services/session';
  import UserService from '../services/user';
  import UnitService from '../services/unit';

  import Datepicker from 'vuejs-datepicker';

  export default {
    name: 'Reports',
    mixins: [SessionService, StorageService, UtilService, HotspotService, UserService, UnitService],
    components: {
      Datepicker
    },
    mounted() {
      // get session list
      this.getAll()
      this.getAllHotspots()
      this.getAllUsers()
      this.getAllUnits()
    },
    data() {
      return {
        msg: 'Reports',
        isLoading: true,
        locale: this.$root.$options.currentLocale,
        columns: [{
            label: this.$i18n.t('session.unit'),
            field: 'unit_id',
            filterable: true,
          }, {
            label: this.$i18n.t('session.user'),
            field: 'user_id',
            filterable: true,
          }, {
            label: this.$i18n.t('session.bytes_up'),
            field: 'bytes_up',
            type: 'number',
            filterable: true,
          }, {
            label: this.$i18n.t('session.bytes_down'),
            field: 'bytes_down',
            type: 'number',
            filterable: true,
          }, {
            label: this.$i18n.t('session.duration'),
            field: 'duration',
            type: 'number',
            filterable: true,
          },
          {
            label: this.$i18n.t('session.start_time'),
            field: 'start_time',
            filterable: true,
          },
          {
            label: this.$i18n.t('session.update_time'),
            field: 'update_time',
            filterable: true,
            sortable: false
          },
          {
            label: '',
            field: '',
            sortable: false
          },
        ],
        rows: [],
        tableLangsTexts: this.tableLangs(),
        hotspots: [],
        hotspotSearchId: this.get('sessions_hotspot_id') || 0,
        hotspotUserId: this.get('sessions_user_id') || 0,
        hotspotUnitId: this.get('sessions_unit_id') || 0,
        hotspotDateFrom: this.get('sessions_date_from') || new Date(Date.now() - 12096e5).toISOString(),
        hotspotDateTo: this.get('sessions_date_to') || new Date().toISOString(),
        user: this.get('loggedUser') || null,
        users: [],
        units: []
      }
    },
    methods: {
      searchFn(row, col, cellValue, searchTerm) {
        var value = cellValue.toString().toLowerCase()
        if (col.field == 'unit_id') {
          value = this.extractUnit(cellValue).description.toLowerCase()
        }

        if (col.field == 'user_id') {
          value = this.extractUser(cellValue).name.toLowerCase()
        }
        return value.startsWith(searchTerm.toLowerCase())
      },
      dateFormatter(date) {
        return this.$root.$options.moment(date).format('DD MMMM YYYY');
      },
      getAllHotspots() {
        this.hotspotGetAll(success => {
          this.hotspots = success.body
          $('[data-toggle="tooltip"]').tooltip()
          this.isLoading = false;
        }, error => {
          console.log(error)
        })
      },
      getAll() {
        // save to storage
        this.set('sessions_hotspot_id', this.hotspotSearchId || this.get('sessions_hotspot_id') || 0)
        this.set('sessions_user_id', this.hotspotUserId || this.get('sessions_user_id') || 0)
        this.set('sessions_unit_id', this.hotspotUnitId || this.get('sessions_unit_id') || 0)
        this.set('sessions_date_from', this.hotspotDateFrom || this.get('sessions_date_from') || new Date(Date.now() -
          12096e5).toISOString())
        this.set('sessions_date_to', this.hotspotDateTo || this.get('sessions_date_to') || new Date().toISOString())

        // preload users and units
        this.getAllUsers()
        this.getAllUnits()

        // get all sessions
        this.sessionGetAll(this.hotspotSearchId, this.hotspotUserId, this.hotspotUnitId, new Date(this.hotspotDateFrom)
          .toISOString(), new Date(this.hotspotDateTo).toISOString(), success => {
            this.rows = success.body
            this.isLoading = false
          }, error => {
            this.isLoading = false
            this.rows = []

            console.log(error)
          })
      },
      getAllUsers() {
        this.userGetAll(this.hotspotSearchId, success => {
          this.users = success.body
        }, error => {
          console.log(error)
          this.users = {}
        })
      },
      getAllUnits() {
        this.unitGetAll(this.hotspotSearchId, success => {
          this.units = success.body
        }, error => {
          console.log(error)
          this.units = {}
        })
      },
      extractUnit(id) {
        for (var u in this.units) {
          if (this.units[u].id == id) {
            return this.units[u]
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
        this.hotspotUnitId = 0
        this.hotspotDateFrom = new Date(Date.now() - 12096e5).toISOString()
        this.hotspotDateTo = new Date().toISOString()
        this.set('sessions_hotspot_id', 0)
        this.set('sessions_user_id', 0)
        this.set('sessions_unit_id', 0)
        this.set('sessions_date_from', new Date(Date.now() - 12096e5).toISOString())
        this.set('sessions_date_to', new Date().toISOString())
        this.getAll()
      },
      exportCSV() {
        var newRows = JSON.parse(JSON.stringify(this.rows))
        for (var r in newRows) {
          newRows[r].unit_id = this.extractUnit(newRows[r].unit_id).description
          newRows[r].user_id = this.extractUser(newRows[r].user_id).name
          newRows[r].bytes_up = this.$options.filters['byteFormat'](newRows[r].bytes_up)
          newRows[r].bytes_down = this.$options.filters['byteFormat'](newRows[r].bytes_down)
          newRows[r].duration = this.$options.filters['secondsInHour'](newRows[r].duration)
          newRows[r].start_time = this.$options.filters['formatDate'](newRows[r].start_time)
          newRows[r].update_time = this.$options.filters['formatDate'](newRows[r].update_time)
        }
        var csv = this.createCSV(this.columns, newRows)
        this.downloadCSV(csv.cols, csv.rows, 'sessions')
      }
    }
  }

</script>
