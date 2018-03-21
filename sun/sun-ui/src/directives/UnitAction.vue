<template>
  <div>
    <div class="dropdown">
      <button class="btn btn-default dropdown-toggle" type="button" id="dropdownMenu1" data-toggle="dropdown">
        {{ $t('action') }}
        <span class="caret"></span>
      </button>
      <ul class="dropdown-menu dropdown-menu-right" role="menu" aria-labelledby="dropdownMenu1">
        <li v-if="details == 'true'" role="presentation">
          <a role="menuitem" tabindex="-1" :href="'#/units/'+ obj.id">
            <span class="fa fa-info action-icon-menu"></span>
            {{ $t('details') }}
          </a>
        </li>
        <li role="presentation" class="divider"></li>
        <li role="presentation">
          <a v-on:click="setCurrentObj(obj)" role="menuitem" tabindex="-1" href="" data-toggle="modal" :data-target="'#UndeleteModal'+obj.id">
            <span class="fa fa-remove action-icon-menu"></span>
            {{ $t('delete') }}
          </a>
        </li>
      </ul>
    </div>
    <div class="modal fade" :id="'UndeleteModal'+obj.id" tabindex="-1" role="dialog" aria-labelledby="UsdeleteModalLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="UsdeleteModalLabel">{{ $t("delete") }} {{currentObj.username}}</h4>
          </div>
          <form class="form-horizontal" role="form" v-on:submit.prevent="deleteUnit(currentObj)">
            <div class="modal-body">
              <div class="alert alert-warning alert-dismissable">
                <span class="pficon pficon-warning-triangle-o"></span>
                <strong>{{ $t("user.warning_delete_title") }}</strong>. {{ $t("user.warning_delete_sub") }}
              </div>
              <div v-if="errors.delete" class="alert alert-danger alert-dismissable">
                <span class="pficon pficon-error-circle-o"></span>
                <strong>{{ $t("user.delete_error_title") }}</strong>. {{ $t("user.delete_error_sub") }}.
              </div>
            </div>
            <div class="modal-footer">
              <span v-if="currentObj.onAction" class="spinner spinner-sm spinner-inline modal-spinner"></span>
              <button type="button" class="btn btn-default" data-dismiss="modal">{{ $t("cancel") }}</button>
              <button type="submit" class="btn btn-danger">{{ $t("delete") }}</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
  import UnitService from '../services/unit';
  import StorageService from '../services/storage';

  import datePicker from 'vue-bootstrap-datetimepicker';

  export default {
    name: 'UnitAction',
    props: ['details', 'obj', 'update'],
    mixins: [UnitService, StorageService],
    components: {
      datePicker
    },
    data() {
      var currentObj = {}
      var errors = {
        delete: false
      }

      return {
        errors: errors,
        currentObj: currentObj,
      }
    },
    methods: {
      setCurrentObj(obj) {
        this.currentObj = Object.assign({}, obj);
      },
      deleteUnit(obj) {
        this.currentObj.onAction = true
        this.unitDelete(obj.id, success => {
          this.currentObj.onAction = false
          $('#UndeleteModal' + obj.id).modal('toggle');
          this.update()
        }, error => {
          this.currentObj.onAction = false
          this.errors.delete = true
          console.log(error.body.message);
        })
      }
    }
  }

</script>
<style>


</style>
