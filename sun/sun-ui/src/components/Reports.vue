<template>
  <div>
    <h2>{{ $t("sessions") }}</h2>
    <div v-if="isLoading" class="spinner spinner-lg"></div>
    <vue-good-table v-if="!isLoading" :perPage="25" :columns="columns" :rows="rows" :lineNumbers="false" :defaultSortBy="{field: 'sessionname', type: 'asc'}"
      :globalSearch="true" :paginate="true" styleClass="table" :nextText="tableLangsTexts.nextText" :prevText="tableLangsTexts.prevText"
      :rowsPerPageText="tableLangsTexts.rowsPerPageText" :globalSearchPlaceholder="tableLangsTexts.globalSearchPlaceholder"
      :ofText="tableLangsTexts.ofText">
      <template slot="table-row" slot-scope="props">
        <td class="fancy">{{ props.row.bytes_up }}</td>
        <td class="fancy">{{ props.row.bytes_down }}</td>
        <td class="fancy">{{ props.row.duration }}</td>
        <td class="fancy">{{ props.row.auth_time }}</td>
        <td class="fancy">{{ props.row.start_time }}</td>
        <td class="fancy">{{ props.row.update_time }}</td>
        <td class="fancy">{{ props.row.stop_time }}</td>
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
  import SessionService from '../services/session';

  export default {
    name: 'Reports',
    mixins: [SessionService, StorageService, UtilService],
    data() {
      // get session list
      this.getAll()

      return {
        msg: 'Reports',
        isLoading: true,
        columns: [{
            label: this.$i18n.t('session.bytes_up'),
            field: 'bytes_up',
            filterable: true,
          }, {
            label: this.$i18n.t('session.bytes_down'),
            field: 'bytes_down',
            filterable: true,
          }, {
            label: this.$i18n.t('session.duration'),
            field: 'duration',
            filterable: true,
          },
          {
            label: this.$i18n.t('session.auth_time'),
            field: 'auth_time',
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
            label: this.$i18n.t('session.stop_time'),
            field: 'stop_time',
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
        tableLangsTexts: this.tableLangs()
      }
    },
    methods: {
      getAll() {
        this.sessionGetAll(success => {
          this.rows = success.body
          this.isLoading = false
        }, error => {
          this.isLoading = false
          console.log(error)
        })
      }
    }
  }

</script>
