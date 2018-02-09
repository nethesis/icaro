<template>
  <div>
    <h2>{{ $t("sessions") }}</h2>
    <div v-if="isLoading" class="spinner spinner-lg"></div>
    <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller')" class="form-group select-search">
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
    <vue-good-table v-if="!isLoading" :perPage="25" :columns="columns" :rows="rows" :lineNumbers="false" :defaultSortBy="{field: 'duration', type: 'asc'}"
      :globalSearch="true" :paginate="true" styleClass="table" :nextText="tableLangsTexts.nextText" :prevText="tableLangsTexts.prevText"
      :rowsPerPageText="tableLangsTexts.rowsPerPageText" :globalSearchPlaceholder="tableLangsTexts.globalSearchPlaceholder"
      :ofText="tableLangsTexts.ofText">
      <template slot="table-row" slot-scope="props">
        <td class="fancy">{{ props.row.bytes_up | byteFormat}}</td>
        <td class="fancy">{{ props.row.bytes_down | byteFormat }}</td>
        <td class="fancy">{{ props.row.duration | secondsInHour }}</td>
        <td class="fancy">{{ props.row.start_time | formatDate }}</td>
        <td class="fancy">{{ props.row.stop_time | formatDate }}</td>
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

  export default {
    name: 'Reports',
    mixins: [SessionService, StorageService, UtilService, HotspotService],
    data() {
      // get session list
      this.getAll()
      this.getAllHotspots();

      return {
        msg: 'Reports',
        isLoading: true,
        columns: [{
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
            label: this.$i18n.t('session.stop_time'),
            field: 'stop_time',
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
        hotspotSearchId: 0,
        user: this.get('loggedUser') || null,
      }
    },
    methods: {
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
        this.sessionGetAll(this.hotspotSearchId, success => {
          this.rows = success.body
          this.isLoading = false
        }, error => {
          this.isLoading = false
          this.rows = []
          console.log(error)
        })
      }
    }
  }

</script>
