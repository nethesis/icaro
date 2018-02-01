<template>
  <div>
    <div class="dropdown">
      <button class="btn btn-default dropdown-toggle" type="button" id="dropdownMenu1" data-toggle="dropdown">
        {{ $t('action') }}
        <span class="caret"></span>
      </button>
      <ul class="dropdown-menu dropdown-menu-right" role="menu" aria-labelledby="dropdownMenu1">
        <li v-if="details == 'true'" role="presentation">
          <a role="menuitem" tabindex="-1" :href="'#/hotspots/'+ obj.id">
            <span class="fa-info action-icon-menu"></span>
            {{ $t('details') }}
          </a>
        </li>
        <li role="presentation">
          <a v-on:click="setCurrentObj(obj)" role="menuitem" tabindex="-1" href="" data-toggle="modal" :data-target="'#HSmodifyModal'+obj.id">
            <span class="fa fa-edit action-icon-menu"></span>
            {{ $t('modify') }}
          </a>
        </li>
        <li role="presentation" class="divider"></li>
        <li role="presentation">
          <a v-on:click="setCurrentObj(obj)" role="menuitem" tabindex="-1" href="" data-toggle="modal" :data-target="'#HSdeleteModal'+obj.id">
            <span class="fa fa-remove action-icon-menu"></span>
            {{ $t('delete') }}
          </a>
        </li>
      </ul>
    </div>
    <div class="modal fade" :id="'HSmodifyModal'+obj.id" tabindex="-1" role="dialog" aria-labelledby="HSmodifyModalLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="HSmodifyModalLabel">{{ $t("modify") }} {{currentObj.name}}</h4>
          </div>
          <form class="form-horizontal" role="form" v-on:submit.prevent="modifyHotspot(currentObj)">
            <div class="modal-body">
              <div class="form-group">
                <label class="col-sm-4 control-label" for="textInput-modal-markup">{{ $t("hotspot.name") }}</label>
                <div class="col-sm-8">
                  <input disabled required v-model="currentObj.name" type="text" id="textInput-modal-markup" class="form-control" :placeholder="$t('hotspot.name')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="textInput2-modal-markup">{{ $t("hotspot.description") }}</label>
                <div class="col-sm-8">
                  <input required v-model="currentObj.description" type="text" id="textInput2-modal-markup" class="form-control" :placeholder="$t('hotspot.description')">
                </div>
              </div>
              <div v-if="errors.update" class="alert alert-danger alert-dismissable">
                <span class="pficon pficon-error-circle-o"></span>
                <strong>{{ $t("hotspot.update_error_title") }}</strong>. {{ $t("hotspot.update_error_sub") }}.
              </div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-default" data-dismiss="modal">{{ $t("cancel") }}</button>
              <button type="submit" class="btn btn-primary">{{ $t("update") }}</button>
            </div>
          </form>
        </div>
      </div>
    </div>
    <div class="modal fade" :id="'HSdeleteModal'+obj.id" tabindex="-1" role="dialog" aria-labelledby="HSdeleteModalLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="HSdeleteModalLabel">{{ $t("delete") }} {{currentObj.name}}</h4>
          </div>
          <form class="form-horizontal" role="form" v-on:submit.prevent="deleteHotspot(currentObj)">
            <div class="modal-body">
              <div class="alert alert-warning alert-dismissable">
                <span class="pficon pficon-warning-triangle-o"></span>
                <strong>{{ $t("hotspot.warning_delete_title") }}</strong>. {{ $t("hotspot.warning_delete_sub") }}.
              </div>
              <div v-if="errors.delete" class="alert alert-danger alert-dismissable">
                <span class="pficon pficon-error-circle-o"></span>
                <strong>{{ $t("hotspot.delete_error_title") }}</strong>. {{ $t("hotspot.delete_error_sub") }}.
              </div>
            </div>
            <div class="modal-footer">
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
  import HotspotService from '../services/hotspot';
  import StorageService from '../services/storage';

  export default {
    name: 'HotspotAction',
    props: ['details', 'obj', 'update'],
    mixins: [HotspotService, StorageService],
    data() {
      var currentObj = {}
      var errors = {
        update: false,
        delete: false
      }
      return {
        errors: errors,
        currentObj: currentObj
      }
    },
    methods: {
      setCurrentObj(obj) {
        this.currentObj = Object.assign({}, obj);
      },
      modifyHotspot(obj) {
        this.execModify(obj.id, {
          description: obj.description
        }, success => {
          $('#HSmodifyModal' + obj.id).modal('toggle');
          this.update()
        }, error => {
          this.errors.update = true
          console.log(error.body.message);
        })
      },
      deleteHotspot(obj) {
        this.execDelete(obj.id, success => {
          $('#HSdeleteModal' + obj.id).modal('toggle');
          this.update()
        }, error => {
          this.errors.delete = true
          console.log(error.body.message);
        })
      }
    }
  }

</script>
<style>


</style>
