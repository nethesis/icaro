<template>
  <div>
    <h2>{{ msg }}</h2>
    <vue-good-table :perPage="25" :columns="columns" :rows="rows" :lineNumbers="false" :defaultSortBy="{field: 'username', type: 'asc'}"
      :globalSearch="true" :paginate="true" styleClass="table" :nextText="tableLangsTexts.nextText" :prevText="tableLangsTexts.prevText"
      :rowsPerPageText="tableLangsTexts.rowsPerPageText" :globalSearchPlaceholder="tableLangsTexts.globalSearchPlaceholder"
      :ofText="tableLangsTexts.ofText">
      <template slot="table-row" slot-scope="props">
        <td>
          <strong>{{ props.row.username }}</strong>
        </td>
        <td class="fancy">{{ props.row.name }}</td>
        <td class="fancy">{{ props.row.email || '-' }}</td>
        <td>
          <strong>{{ props.row.account_type }}</strong>
        </td>
        <td class="fancy">
          <div>
            <strong>{{ $t('user.kbps_down') }}</strong>: {{ props.row.kbps_down }}</div>
          <div>
            <strong>{{ $t('user.kbps_up') }}</strong>: {{ props.row.kbps_up }}</div>
        </td>
        <td class="fancy">
          <div>
            <strong>{{ $t('user.valid_from') }}</strong>: {{ props.row.valid_from }}</div>
          <div>
            <strong>{{ $t('user.valid_until') }}</strong>: {{ props.row.valid_until }}</div>
        </td>
        <td>
          <user-action details="false" :obj="props.row" :update="getAll"></user-action>
        </td>
        <td>
          <a :href="'#/users/'+ props.row.id">
            <span class="fa fa-angle-right details-arrow"></span>
          </a>
        </td>
      </template>
    </vue-good-table>
  </div>
</template>

<script>
  import UserService from '../services/user';
  import StorageService from '../services/storage';
  import UtilService from '../services/util';

  import UserAction from '../directives/UserAction.vue';

  export default {
    name: 'Users',
    mixins: [UserService, StorageService, UtilService],
    components: {
      userAction: UserAction
    },
    data() {
      // get user list
      this.getAll()

      return {
        msg: 'Users',
        columns: [{
            label: this.$i18n.t('user.username'),
            field: 'username',
            filterable: true,
          }, {
            label: this.$i18n.t('user.name'),
            field: 'name',
            filterable: true,
          },
          {
            label: this.$i18n.t('user.email'),
            field: 'email',
            filterable: true,
          },
          {
            label: this.$i18n.t('user.type'),
            field: 'type',
            filterable: true,
          },
          {
            label: this.$i18n.t('user.bandwidth'),
            field: 'bandwidth',
            filterable: true,
          },
          {
            label: this.$i18n.t('user.valid'),
            field: 'valid',
            filterable: true,
          },
          {
            label: this.$i18n.t('action'),
            field: '',
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
        this.execGetAll(success => {
          this.rows = success.body
        }, error => {
          console.log(error)
        })
      }
    }
  }

</script>
