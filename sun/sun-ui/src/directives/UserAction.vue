<template>
  <div>
    <div class="dropup">
      <button
        class="btn btn-default dropdown-toggle"
        type="button"
        id="dropdownMenu1"
        data-toggle="dropdown"
      >
        {{ $t('action') }}
        <span class="caret"></span>
      </button>
      <ul class="dropdown-menu dropdown-menu-right" role="menu" aria-labelledby="dropdownMenu1">
        <li v-if="details == 'true'" role="presentation">
          <a role="menuitem" tabindex="-1" :href="'#/users/'+ obj.id">
            <span class="fa fa-info action-icon-menu"></span>
            {{ $t('details') }}
          </a>
        </li>
        <li role="presentation">
          <a
            v-on:click="setCurrentObj(obj)"
            role="menuitem"
            tabindex="-1"
            href
            data-toggle="modal"
            :data-target="'#UsmodifyModal'+obj.id"
          >
            <span class="fa fa-edit action-icon-menu"></span>
            {{ $t('modify') }}
          </a>
        </li>
        <li role="presentation" class="divider"></li>
        <li role="presentation">
          <a
            v-on:click="setCurrentObj(obj)"
            role="menuitem"
            tabindex="-1"
            href
            data-toggle="modal"
            :data-target="'#UsdeleteModal'+obj.id"
          >
            <span class="fa fa-remove action-icon-menu"></span>
            {{ $t('delete') }}
          </a>
        </li>
      </ul>
    </div>
    <div
      class="modal fade"
      :id="'UsmodifyModal'+obj.id"
      tabindex="-1"
      role="dialog"
      aria-labelledby="UsmodifyModalLabel"
      aria-hidden="true"
    >
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
                <label
                  class="col-sm-4 control-label"
                  for="textInput-modal-markup"
                >{{ $t("user.name") }}</label>
                <div class="col-sm-8">
                  <input
                    required
                    v-model="currentObj.name"
                    type="text"
                    id="textInput-modal-markup"
                    class="form-control"
                    :placeholder="$t('user.name')"
                  >
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-4 control-label"
                  for="textInput2-modal-markup"
                >{{ $t("user.email") }}</label>
                <div class="col-sm-8">
                  <input
                    v-model="currentObj.email"
                    type="text"
                    id="textInput2-modal-markup"
                    class="form-control"
                    :placeholder="$t('user.email')"
                  >
                </div>
              </div>

              <div class="form-group">
                <label
                  class="col-sm-4 control-label"
                  for="textInput2-modal-markup"
                >{{ $t("user.kbps_down") }}</label>
                <div class="col-sm-8">
                  <input
                    v-model="currentObj.kbps_down"
                    type="number"
                    id="textInput2-modal-markup"
                    class="form-control"
                    :placeholder="$t('user.kbps_down')"
                  >
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-4 control-label"
                  for="textInput2-modal-markup"
                >{{ $t("user.kbps_up") }}</label>
                <div class="col-sm-8">
                  <input
                    v-model="currentObj.kbps_up"
                    type="number"
                    id="textInput2-modal-markup"
                    class="form-control"
                    :placeholder="$t('user.kbps_up')"
                  >
                </div>
              </div>

              <div class="form-group">
                <label
                  class="col-sm-4 control-label"
                  for="textInput2-modal-markup"
                >{{ $t("user.traffic_limit") }}</label>
                <div class="col-sm-8">
                  <input
                    v-model="currentObj.max_navigation_traffic"
                    type="number"
                    id="textInput2-modal-markup"
                    class="form-control"
                    :placeholder="$t('user.traffic_limit')"
                  >
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-4 control-label"
                  for="textInput2-modal-markup"
                >{{ $t("user.time_limit") }}</label>
                <div class="col-sm-8">
                  <input
                    v-model="currentObj.max_navigation_time"
                    type="number"
                    id="textInput2-modal-markup"
                    class="form-control"
                    :placeholder="$t('user.time_limit')"
                  >
                </div>
              </div>

              <div class="form-group">
                <label
                  class="col-sm-4 control-label"
                  for="textInput2-modal-markup"
                >{{ $t("user.auto_login") }}</label>
                <div class="col-sm-8">
                  <input
                    v-model="currentObj.auto_login"
                    type="checkbox"
                    id="textInput2-modal-markup"
                    class="form-control"
                    :placeholder="$t('user.auto_login')"
                  >
                </div>
              </div>

              <div class="form-group">
                <label
                  class="col-sm-4 control-label"
                  for="textInput2-modal-markup"
                >{{ $t("user.valid_from") }}</label>
                <div class="col-sm-8">
                  <date-picker v-model="currentObj.valid_from" :config="dateConfig"></date-picker>
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-4 control-label"
                  for="textInput2-modal-markup"
                >{{ $t("user.valid_until") }}</label>
                <div class="col-sm-8">
                  <date-picker v-model="currentObj.valid_until" :config="dateConfig"></date-picker>
                </div>
              </div>

              <div v-if="errors.update" class="alert alert-danger alert-dismissable">
                <span class="pficon pficon-error-circle-o"></span>
                <strong>{{ $t("user.update_error_title") }}</strong>
                . {{ $t("user.update_error_sub") }}.
              </div>
            </div>
            <div class="modal-footer">
              <span
                v-if="currentObj.onAction"
                class="spinner spinner-sm spinner-inline modal-spinner"
              ></span>
              <button type="button" class="btn btn-default" data-dismiss="modal">{{ $t("cancel") }}</button>
              <button type="submit" class="btn btn-primary">{{ $t("update") }}</button>
            </div>
          </form>
        </div>
      </div>
    </div>
    <div
      class="modal fade"
      :id="'UsdeleteModal'+obj.id"
      tabindex="-1"
      role="dialog"
      aria-labelledby="UsdeleteModalLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="UsdeleteModalLabel">{{ $t("delete") }} {{currentObj.name}}</h4>
          </div>
          <form class="form-horizontal" role="form" v-on:submit.prevent="deleteUser(currentObj)">
            <div class="modal-body">
              <div class="alert alert-warning alert-dismissable">
                <span class="pficon pficon-warning-triangle-o"></span>
                <strong>{{ $t("user.warning_delete_title") }}</strong>
                . {{ $t("user.warning_delete_sub") }}
              </div>
              <div v-if="errors.delete" class="alert alert-danger alert-dismissable">
                <span class="pficon pficon-error-circle-o"></span>
                <strong>{{ $t("user.delete_error_title") }}</strong>
                . {{ $t("user.delete_error_sub") }}.
              </div>
            </div>
            <div class="modal-footer">
              <span
                v-if="currentObj.onAction"
                class="spinner spinner-sm spinner-inline modal-spinner"
              ></span>
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
import UserService from "../services/user";
import StorageService from "../services/storage";

import datePicker from "vue-bootstrap-datetimepicker";

var moment = require("moment");

export default {
  name: "UserAction",
  props: ["details", "obj", "update", "expired"],
  mixins: [UserService, StorageService],
  components: {
    datePicker
  },
  data() {
    var currentObj = {};
    var errors = {
      update: false,
      delete: false
    };

    return {
      errors: errors,
      currentObj: currentObj,
      dateConfig: {
        format: "YYYY-MM-DD HH:mm:ss",
        useCurrent: false,
        sideBySide: true,
        locale: this.$root.$options.currentLocale
      }
    };
  },
  methods: {
    setCurrentObj(obj) {
      this.currentObj = Object.assign({}, obj);
      this.currentObj.max_navigation_time =
        this.currentObj.max_navigation_time / 60;
      this.currentObj.max_navigation_traffic =
        this.currentObj.max_navigation_traffic / 1024 / 1024;
      this.currentObj.valid_from = moment(String(this.currentObj.valid_from));
      this.currentObj.valid_until = moment(String(this.currentObj.valid_until));
    },
    modifyUser(obj) {
      this.currentObj.onAction = true;
      this.userModify(
        obj.id,
        {
          name: obj.name,
          email: obj.email,
          kbps_down: parseInt(obj.kbps_down),
          kbps_up: parseInt(obj.kbps_up),
          max_navigation_time: parseInt(obj.max_navigation_time) * 60,
          max_navigation_traffic:
            parseInt(obj.max_navigation_traffic) * 1024 * 1024,
          auto_login: obj.auto_login || false,
          valid_from: new Date(obj.valid_from).toISOString(),
          valid_until: new Date(obj.valid_until).toISOString()
        },
        success => {
          this.currentObj.onAction = false;
          $("#UsmodifyModal" + obj.id).modal("toggle");
          this.update();
        },
        error => {
          this.currentObj.onAction = false;
          this.errors.update = true;
          console.error(error.body.message);
        },
        this.expired
      );
    },
    deleteUser(obj) {
      this.currentObj.onAction = true;
      this.userDelete(
        obj.id,
        success => {
          this.currentObj.onAction = false;
          $("#UsdeleteModal" + obj.id).modal("toggle");
          this.update();
        },
        error => {
          this.currentObj.onAction = false;
          this.errors.delete = true;
          console.error(error.body.message);
        },
        this.expired
      );
    }
  }
};
</script>
<style>
</style>
