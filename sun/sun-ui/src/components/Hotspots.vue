<template>
  <div>
    <h2>{{ msg }}</h2>
    <div v-if="isLoading" class="spinner spinner-lg"></div>
    <button v-if="rows.length > 0 && !isLoading" data-toggle="modal" data-target="#HScreateModal" class="btn btn-primary btn-lg create-hotspot">
      {{ $t('hotspot.create_new') }} </button>
    <div v-if="rows.length == 0 && !isLoading" class="blank-slate-pf " id="">
      <div class="blank-slate-pf-icon">
        <span class="fa fa-wifi"></span>
      </div>
      <h1>
        {{ $t('hotspot.no_hotspot_found') }}
      </h1>
      <p>
        {{ $t('hotspot.no_hotspot_found_sub') }}.
      </p>
      <div class="blank-slate-pf-main-action">
        <button data-toggle="modal" data-target="#HScreateModal" class="btn btn-primary btn-lg"> {{ $t('hotspot.create_new') }} </button>
      </div>
    </div>
    <vue-good-table v-if="rows.length > 0 && !isLoading" :perPage="25" :columns="columns" :rows="rows" :lineNumbers="false" :defaultSortBy="{field: 'name', type: 'asc'}"
      :globalSearch="true" :paginate="true" styleClass="table" :nextText="tableLangsTexts.nextText" :prevText="tableLangsTexts.prevText"
      :rowsPerPageText="tableLangsTexts.rowsPerPageText" :globalSearchPlaceholder="tableLangsTexts.globalSearchPlaceholder"
      :ofText="tableLangsTexts.ofText">
      <template slot="table-row" slot-scope="props">
        <td>
          <a :href="'#/hotspots/'+ props.row.id">
            <strong>{{ props.row.name }}</strong>
          </a>
        </td>
        <td class="fancy">{{ props.row.description }}</td>
        <td class="fancy">{{ props.row.created }}</td>
        <td>
          <hotspot-action details="true" :obj="props.row" :update="getAll"></hotspot-action>
        </td>
        <td>
          <a :href="'#/hotspots/'+ props.row.id">
            <span class="fa fa-angle-right details-arrow"></span>
          </a>
        </td>
      </template>
    </vue-good-table>
    <div class="modal fade" id="HScreateModal" tabindex="-1" role="dialog" aria-labelledby="HScreateModalLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="HScreateModalLabel">{{ $t("create") }}</h4>
          </div>
          <form class="form-horizontal" role="form" v-on:submit.prevent="createHotspot()">
            <div class="modal-body">
              <div class="form-group">
                <label class="col-sm-4 control-label" for="textInput-modal-markup">{{ $t("hotspot.name") }}</label>
                <div class="col-sm-8">
                  <input required v-model="newObj.name" pattern="^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$" type="text" id="textInput-modal-markup" class="form-control" :placeholder="$t('hotspot.name')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="textInput2-modal-markup">{{ $t("hotspot.description") }}</label>
                <div class="col-sm-8">
                  <input required v-model="newObj.description" type="text" id="textInput2-modal-markup" class="form-control" :placeholder="$t('hotspot.description')">
                </div>
              </div>
              <div v-if="errors.create" class="alert alert-danger alert-dismissable">
                <span class="pficon pficon-error-circle-o"></span>
                <strong>{{ $t("hotspot.create_error_title") }}</strong>. {{ $t("hotspot.create_error_sub") }}.
              </div>
            </div>
            <div class="modal-footer">
              <span v-if="newObj.onAction" class="spinner spinner-sm spinner-inline modal-spinner"></span>
              <button type="button" class="btn btn-default" data-dismiss="modal">{{ $t("cancel") }}</button>
              <button type="submit" class="btn btn-primary">{{ $t("create") }}</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import HotspotService from '../services/hotspot';
  import StorageService from '../services/storage';
  import UtilService from '../services/util';

  import HotspotAction from '../directives/HotspotAction.vue';

  export default {
    name: 'Hotspots',
    mixins: [HotspotService, StorageService, UtilService],
    components: {
      hotspotAction: HotspotAction
    },
    data() {
      // get hotspot list
      this.getAll()

      var newObj = {
        name: '',
        description: ''
      }
      var errors = {
        create: false,
      }

      return {
        msg: 'Hotspots',
        isLoading: true,
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
        tableLangsTexts: this.tableLangs(),
        newObj: newObj,
        errors: errors
      }
    },
    methods: {
      getAll() {
        this.hotspotGetAll(success => {
          this.rows = success.body
          this.isLoading = false
        }, error => {
          this.isLoading = false
          console.log(error)
        })
      },
      createHotspot() {
        this.newObj.onAction = true
        this.hotspotCreate(this.newObj, success => {
          this.newObj.onAction = false
          this.newObj.name = ""
          this.newObj.description = ""
          $('#HScreateModal').modal('toggle');
          this.getAll()
        }, error => {
          this.newObj.onAction = false
          this.errors.create = true
          console.log(error.body.message);
        })
      }
    }
  }

</script>

<style>
  .create-hotspot {
    float: right;
    margin-top: -52px;
    margin-right: 35px;
  }

</style>
