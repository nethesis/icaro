<template>
  <div>
    <h2>{{ msg }}</h2>
    <vue-good-table :columns="columns" :rows="rows" :lineNumbers="false" :defaultSortBy="{field: 'name', type: 'asc'}" :globalSearch="true"
      :paginate="true" styleClass="table" :nextText="tableLangsTexts.nextText" :prevText="tableLangsTexts.prevText" :rowsPerPageText="tableLangsTexts.rowsPerPageText"
      :globalSearchPlaceholder="tableLangsTexts.globalSearchPlaceholder" :ofText="tableLangsTexts.ofText">
      <template slot="table-row" scope="props">
        <td>
          <strong>{{ props.row.name }}</strong>
        </td>
        <td class="fancy">{{ props.row.description }}</td>
        <td class="fancy">{{ props.row.created }}</td>
        <td>
          <a :href="'#/hotspots/'+ props.row.id" class="btn btn-primary" type="button">{{ $t('hotspot.details') }}</a>
        </td>
      </template>
    </vue-good-table>
  </div>
</template>

<script>
  import HotspotService from '../services/hotspot';
  import StorageService from '../services/storage';
  import UtilService from '../services/util';

  export default {
    name: 'Hotspots',
    mixins: [HotspotService, StorageService, UtilService],
    data() {
      // get hotspot list
      this.getAll()

      return {
        msg: 'Hotspots',
        columns: [{
            label: this.$i18n.t('hotspot.name'),
            field: 'name',
            filterable: true,
          },
          {
            label: this.$i18n.t('hotspot.description'),
            field: 'description',
            filterable: true,
          },
          {
            label: this.$i18n.t('hotspot.created'),
            field: 'created',
            sortable: false
          },
          {
            label: this.$i18n.t('hotspot.action'),
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

<style>
  .global-search {
    padding-left: 0px !important;
  }

  .global-search-icon {
    display: none;
  }

  .good-table {
    padding-left: 20px;
    padding-right: 20px;
  }

  .table thead input {
    background-image: none;
    border-radius: 1px !important;
  }

  .table thead {
    background-image: none;
  }

  .table thead tr th:first-child {
    padding-left: 25px !important;
  }

  .table tbody tr td:first-child {
    padding-left: 25px !important;
  }

  .table thead th {
    color: #464646 !important;
    border: none !important;
    background: #f5f5f5 !important;
    border-bottom: 2px solid #39a5dd !important;
    font-size: 14px !important;
  }

  .table tbody tr {
    background-color: white;
    height: 47px;
  }

  .table tbody tr:hover {
    background-color: #def3ff;
  }

  .table tbody tr td {
    line-height: 31px;
  }

</style>
