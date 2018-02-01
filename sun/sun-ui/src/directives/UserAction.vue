<template>
  <div>
    <div class="dropdown">
      <button class="btn btn-default dropdown-toggle" type="button" id="dropdownMenu1" data-toggle="dropdown">
        {{ $t('action') }}
        <span class="caret"></span>
      </button>
      <ul class="dropdown-menu dropdown-menu-right" role="menu" aria-labelledby="dropdownMenu1">
        <li v-if="details == 'true'" role="presentation">
          <a role="menuitem" tabindex="-1" :href="'#/users/'+ obj.id">
            <span class="fa-info action-icon-menu"></span>
            {{ $t('details') }}
          </a>
        </li>
        <li role="presentation">
          <a v-on:click="setCurrentObj(obj)" role="menuitem" tabindex="-1" href="" data-toggle="modal" :data-target="'#UsmodifyModal'+obj.id">
            <span class="fa fa-edit action-icon-menu"></span>
            {{ $t('modify') }}
          </a>
        </li>
        <li role="presentation" class="divider"></li>
        <li role="presentation">
          <a v-on:click="setCurrentObj(obj)" role="menuitem" tabindex="-1" href="" data-toggle="modal" :data-target="'#UsdeleteModal'+obj.id">
            <span class="fa fa-remove action-icon-menu"></span>
            {{ $t('delete') }}
          </a>
        </li>
      </ul>
    </div>
    <div class="modal fade" :id="'UsmodifyModal'+obj.id" tabindex="-1" role="dialog" aria-labelledby="UsmodifyModalLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="UsmodifyModalLabel">{{ $t("modify") }} {{currentObj.name}}</h4>
          </div>
          <form class="form-horizontal" role="form" v-on:submit.prevent="modifyUser(currentObj)">
            <div class="modal-body">
              <div class="form-group">
                <label class="col-sm-4 control-label" for="textInput-modal-markup">{{ $t("user.name") }}</label>
                <div class="col-sm-8">
                  <input required v-model="currentObj.name" type="text" id="textInput-modal-markup" class="form-control" :placeholder="$t('user.name')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="textInput2-modal-markup">{{ $t("user.email") }}</label>
                <div class="col-sm-8">
                  <input required v-model="currentObj.email" type="text" id="textInput2-modal-markup" class="form-control" :placeholder="$t('user.email')">
                </div>
              </div>

              <div class="form-group">
                <label class="col-sm-4 control-label" for="textInput2-modal-markup">{{ $t("user.kbps_down") }}</label>
                <div class="col-sm-8">
                  <input v-model="currentObj.kbps_down" type="number" id="textInput2-modal-markup" class="form-control" :placeholder="$t('user.kbps_down')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="textInput2-modal-markup">{{ $t("user.kbps_up") }}</label>
                <div class="col-sm-8">
                  <input v-model="currentObj.kbps_up" type="number" id="textInput2-modal-markup" class="form-control" :placeholder="$t('user.kbps_up')">
                </div>
              </div>

              <div class="form-group">
                <label class="col-sm-4 control-label" for="textInput2-modal-markup">{{ $t("user.valid_from") }}</label>
                <div class="col-sm-8">
                  <input v-model="currentObj.valid_from" type="date" id="textInput2-modal-markup" class="form-control" :placeholder="$t('user.valid_from')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="textInput2-modal-markup">{{ $t("user.valid_until") }}</label>
                <div class="col-sm-8">
                  <input type="text" class="form-control bootstrap-datepicker">
                  <!-- <input v-model="currentObj.valid_until" type="date" id="textInput2-modal-markup" class="form-control" :placeholder="$t('user.valid_until')"> -->
                </div>
              </div>

              <div v-if="errors.update" class="alert alert-danger alert-dismissable">
                <span class="pficon pficon-error-circle-o"></span>
                <strong>{{ $t("user.update_error_title") }}</strong>. {{ $t("user.update_error_sub") }}.
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
    <div class="modal fade" :id="'UsdeleteModal'+obj.id" tabindex="-1" role="dialog" aria-labelledby="UsdeleteModalLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="UsdeleteModalLabel">{{ $t("delete") }} {{currentObj.username}}</h4>
          </div>
          <form class="form-horizontal" role="form" v-on:submit.prevent="deleteUser(currentObj)">
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
  import UserService from '../services/user';
  import StorageService from '../services/storage';
  export default {
    name: 'UserAction',
    props: ['details', 'obj', 'update'],
    mixins: [UserService, StorageService],
    data() {
      var currentObj = {}
      var errors = {
        update: false,
        delete: false
      }

      this.initGraphics()

      return {
        errors: errors,
        currentObj: currentObj
      }
    },
    methods: {
      initGraphics() {
        $(document).ready(function () {
          console.log("ready!");
          $('.bootstrap-datepicker').datepicker({
            autoclose: true,
            todayBtn: "linked",
            todayHighlight: true
          });
        });


      },
      setCurrentObj(obj) {
        this.currentObj = Object.assign({}, obj);
      },
      modifyUser(obj) {
        this.execModify(obj.id, {
          name: obj.name,
          email: obj.email,
          kbps_down: obj.kbps_down,
          kbps_up: obj.kbps_up,
          valid_from: obj.valid_from,
          valid_until: obj.valid_until
        }, success => {
          $('#UsmodifyModal' + obj.id).modal('toggle');
          this.update()
        }, error => {
          this.errors.update = true
          console.log(error.body.message);
        })
      },
      deleteUser(obj) {
        this.execDelete(obj.id, success => {
          $('#UsdeleteModal' + obj.id).modal('toggle');
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
