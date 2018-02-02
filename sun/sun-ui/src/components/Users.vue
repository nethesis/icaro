<template>
  <div>
    <h2>{{ msg }}</h2>
    <div v-if="isLoading" class="spinner spinner-lg"></div>
    <vue-good-table v-if="!isLoading" :perPage="25" :columns="columns" :rows="rows" :lineNumbers="false" :defaultSortBy="{field: 'username', type: 'asc'}"
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
            <strong>{{ $t('user.from') }}</strong>: {{ props.row.valid_from }}</div>
          <div>
            <strong>{{ $t('user.until') }}</strong>: {{ props.row.valid_until }}</div>
        </td>
        <td>
          <user-action details="false" :obj="props.row" :update="getAll"></user-action>
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
        isLoading: true,
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
            sortable: false
          },
          {
            label: this.$i18n.t('user.valid'),
            field: 'valid',
            filterable: true,
            sortable: false
          },
          {
            label: this.$i18n.t('action'),
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
        this.userGetAll(success => {
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
