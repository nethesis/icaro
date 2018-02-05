<template>
  <div>
    <div class="dropdown">
      <button class="btn btn-default dropdown-toggle" type="button" id="ACdropdownMenu1" data-toggle="dropdown">
        {{ $t('action') }}
        <span class="caret"></span>
      </button>
      <ul class="dropdown-menu dropdown-menu-right" role="menu" aria-labelledby="dropdownMenu1">
        <li v-if="details == 'true'" role="presentation">
          <a role="menuitem" tabindex="-1" :href="'#/accounts/'+ obj.id">
            <span class="fa-info action-icon-menu"></span>
            {{ $t('details') }}
          </a>
        </li>
        <li role="presentation">
          <a v-on:click="setCurrentObj(obj)" role="menuitem" tabindex="-1" href="" data-toggle="modal" :data-target="'#ACModifyModal'+obj.id">
            <span class="fa fa-edit action-icon-menu"></span>
            {{ $t('modify') }}
          </a>
        </li>
        <li role="presentation">
          <a v-on:click="setCurrentObj(obj)" role="menuitem" tabindex="-1" href="" data-toggle="modal" :data-target="'#ACChangePasswordModal'+obj.id">
            <span class="fa fa-key action-icon-menu"></span>
            {{ $t('change_password') }}
          </a>
        </li>
        <li role="presentation" class="divider"></li>
        <li role="presentation">
          <a v-on:click="setCurrentObj(obj)" role="menuitem" tabindex="-1" href="" data-toggle="modal" :data-target="'#ACDeleteModal'+obj.id">
            <span class="fa fa-remove action-icon-menu"></span>
            {{ $t('delete') }}
          </a>
        </li>
      </ul>
    </div>
    <div class="modal fade" :id="'ACModifyModal'+obj.id" tabindex="-1" role="dialog" aria-labelledby="ACModifyModalLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="ACModifyModalLabel">{{ $t("modify") }} {{currentObj.accountname}}</h4>
          </div>
          <form class="form-horizontal" role="form" v-on:submit.prevent="modifyAccount(currentObj)">
            <div class="modal-body">
              <div class="form-group">
                <label class="col-sm-4 control-label" for="accountname">{{ $t("account.name") }}</label>
                <div class="col-sm-8">
                  <input required v-model="currentObj.name" type="text" id="accountname" class="form-control" :placeholder="$t('account.name')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="ACtextInput2-modal-markup">{{ $t("account.email") }}</label>
                <div class="col-sm-8">
                  <input required v-model="currentObj.email" type="text" id="ACtextInput2-modal-markup" class="form-control" :placeholder="$t('account.email')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="ACtextInput2-modal-markup">{{ $t("account.type") }}</label>
                <div class="col-sm-8">
                  <select v-model="currentObj.type" class="form-control">
                    <option value="customer">{{ $t("account.type_customer") }}</option>
                    <option value="desk">{{ $t("account.type_desk") }}</option>
                  </select>
                </div>
              </div>
              <div v-if="errors.update" class="alert alert-danger alert-dismissable">
                <span class="pficon pficon-error-circle-o"></span>
                <strong>{{ $t("account.update_error_title") }}</strong>. {{ $t("account.update_error_sub") }}.
              </div>
            </div>
            <div class="modal-footer">
              <span v-if="currentObj.onAction" class="spinner spinner-sm spinner-inline modal-spinner"></span>
              <button type="button" class="btn btn-default" data-dismiss="modal">{{ $t("cancel") }}</button>
              <button type="submit" class="btn btn-primary">{{ $t("update") }}</button>
            </div>
          </form>
        </div>
      </div>
    </div>
    <div class="modal fade" :id="'ACDeleteModal'+obj.id" tabindex="-1" role="dialog" aria-labelledby="ACDeleteModalLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="ACDeleteModalLabel">{{ $t("delete") }} {{currentObj.accountname}}</h4>
          </div>
          <form class="form-horizontal" role="form" v-on:submit.prevent="deleteAccount(currentObj)">
            <div class="modal-body">
              <div class="alert alert-warning alert-dismissable">
                <span class="pficon pficon-warning-triangle-o"></span>
                <strong>{{ $t("account.warning_delete_title") }}</strong>. {{ $t("account.warning_delete_sub") }}
              </div>
              <div v-if="errors.delete" class="alert alert-danger alert-dismissable">
                <span class="pficon pficon-error-circle-o"></span>
                <strong>{{ $t("account.delete_error_title") }}</strong>. {{ $t("account.delete_error_sub") }}.
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
    <div class="modal fade" :id="'ACChangePasswordModal'+obj.id" tabindex="-1" role="dialog" aria-labelledby="ACChangePassModalLabel"
      aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="changePassModalLabel">{{ $t("account.change_password") }}</h4>
          </div>
          <form class="form-horizontal" role="form" v-on:submit.prevent="changePassword(currentObj)">
            <div class="modal-body">
              <div class="form-group">
                <label class="col-sm-4 control-label" for="textInput-modal-markup">{{ $t("account.new_password") }}</label>
                <div class="col-sm-8">
                  <input v-model="newPassword" required type="text" id="textInput-modal-markup" class="form-control" :placeholder="$t('account.new_password')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="textInput2-modal-markup">{{ $t("account.confirm_password") }}</label>
                <div class="col-sm-8">
                  <input v-model="confirmPassword" required type="text" id="textInput2-modal-markup" class="form-control" :placeholder="$t('account.confirm_password')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="textInput2-modal-markup"></label>
                <div class="col-sm-8">
                  <span v-bind:class="[newPassword == confirmPassword ? 'pass-confirm-ok' : 'pass-confirm-err', '']"></span>
                </div>
              </div>
              <div v-if="errors.password" class="alert alert-danger alert-dismissable">
                <span class="pficon pficon-error-circle-o"></span>
                <strong>{{ $t("profile.change_password_error_title") }}</strong>. {{ $t("account.change_password_error_sub") }}.
              </div>
            </div>
            <div class="modal-footer">
              <span v-if="currentObj.onAction" class="spinner spinner-sm spinner-inline modal-spinner"></span>
              <button type="button" class="btn btn-default" data-dismiss="modal">{{ $t("cancel") }}</button>
              <button :disabled="(newPassword != confirmPassword) || (confirmPassword.length == 0) || (newPassword.length == 0)" type="submit"
                class="btn btn-primary">{{ $t("update") }}</button>

            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
  import AccountService from '../services/account';
  import StorageService from '../services/storage';
  export default {
    name: 'AccountAction',
    props: ['details', 'obj', 'update'],
    mixins: [AccountService, StorageService],
    data() {
      var currentObj = {}
      var newPassword, confirmPassword = ''

      var errors = {
        update: false,
        delete: false,
        password: false
      }

      return {
        errors: errors,
        currentObj: currentObj,
        newPassword: newPassword,
        confirmPassword: confirmPassword,
        errors: errors
      }
    },
    methods: {
      setCurrentObj(obj) {
        this.currentObj = Object.assign({}, obj);
      },
      modifyAccount(obj) {
        this.currentObj.onAction = true
        this.accountModify(obj.id, {
          name: obj.name,
          email: obj.email,
          type: obj.type
        }, success => {
          this.currentObj.onAction = false
          $('#ACModifyModal' + obj.id).modal('toggle');
          this.update()
        }, error => {
          this.currentObj.onAction = false
          this.errors.update = true
          console.log(error.body.message);
        })
      },
      deleteAccount(obj) {
        this.currentObj.onAction = true
        this.accountDelete(obj.id, success => {
          this.currentObj.onAction = false
          $('#ACDeleteModal' + obj.id).modal('toggle');
          this.update()
        }, error => {
          this.currentObj.onAction = false
          this.errors.delete = true
          console.log(error.body.message);
        })
      },
      changePassword(obj) {
        this.currentObj.onAction = true
        this.accountChangePassword(obj.id, this.newPassword, success => {
          this.currentObj.onAction = false
          $('#ACChangePasswordModal' + obj.id).modal('toggle');
        }, error => {
          this.currentObj.onAction = false
          this.errors.password = true
          console.log(error.body.message);
        })
      }

    }
  }

</script>
