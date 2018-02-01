<template>
  <div>
    <h2>{{ msg }}</h2>
    <button v-if="rows.length > 0"  v-on:click="initNewAccount()" data-toggle="modal" data-target="#ACcreateModal" class="btn btn-primary btn-lg create-account"> {{ $t('account.create_new') }} </button>
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
          <strong>{{ props.row.type }}</strong>
        </td>
	<td class="fancy">{{ props.row.created }}</td>
        <td>
          <account-action details="false" :obj="props.row" :update="getAll"></account-action>
        </td>
        <td>
          <a :href="'#/accounts/'+ props.row.id">
            <span class="fa fa-angle-right details-arrow"></span>
          </a>
        </td>
      </template>
    </vue-good-table>

    <div class="modal fade" id="ACcreateModal" tabindex="-1" role="dialog" aria-labelledby="ACModifyModalLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="ACModifyModalLabel">{{ $t("account.create_title") }}</h4>
          </div>
          <form class="form-horizontal" role="form" v-on:submit.prevent="execCreate(newObj)">
            <div class="modal-body">
              <div class="form-group">
                <label class="col-sm-4 control-label" for="accuuid">{{ $t("account.uuid") }}</label>
                <div class="col-sm-8">
                  <input required v-model="newObj.uuid" readonly type="text" id="accuuid" class="form-control" :placeholder="$t('account.uuid')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="accusername">{{ $t("account.username") }}</label>
                <div class="col-sm-8">
                  <input required v-model="newObj.username" type="text" id="accusername" class="form-control" :placeholder="$t('account.username')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="accountname">{{ $t("account.name") }}</label>
                <div class="col-sm-8">
                  <input required v-model="newObj.name" type="text" id="accountname" class="form-control" :placeholder="$t('account.name')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="ACtextInput2-modal-markup">{{ $t("account.email") }}</label>
                <div class="col-sm-8">
                  <input required v-model="newObj.email" type="text" id="ACtextInput2-modal-markup" class="form-control" :placeholder="$t('account.email')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="ACtextInput2-modal-markup">{{ $t("account.type") }}</label>
                <div class="col-sm-8">
                  <select v-model="newObj.type" class="bootstrap-select">
                    <option disabled value="">{{ $t("account.select_one") }}</option>
                    <option value="customer">{{ $t("account.type_customer") }}</option>
                    <option value="desk">{{ $t("account.type_desk") }}</option>
                  </select>
                </div>
              </div>
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
              <div v-if="errors.update" class="alert alert-danger alert-dismissable">
                <span class="pficon pficon-error-circle-o"></span>
                <strong>{{ $t("account.update_error_title") }}</strong>. {{ $t("account.update_error_sub") }}.
              </div>
              <div v-if="errors.create" class="alert alert-danger alert-dismissable">
                <span class="pficon pficon-error-circle-o"></span>
                <strong>{{ $t("account.create_error_title") }}</strong>. {{ $t("account.create_error_sub") }}.
              </div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-default" data-dismiss="modal">{{ $t("cancel") }}</button>
              <button :disabled="(newPassword != confirmPassword) || (confirmPassword.length == 0) || (newPassword.length == 0)" type="submit" class="btn btn-primary">{{ $t("create") }}</button>
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
  import UtilService from '../services/util';

  import AccountAction from '../directives/AccountAction.vue';

  export default {
    name: 'Accounts',
    mixins: [AccountService, StorageService, UtilService],
    components: {
      accountAction: AccountAction
    },
    data() {
      // get account list
      this.getAll()

      var newObj = {
        uuid: '',
        username: '',
        name: '',
        email: '',
        type: '',
        password: '' 
      }

      var errors = {
        create: false,
      }

      var newPassword = '';
      var confirmPassword = '';


      return {
        msg: 'Accounts',
	isLoading: true,
        columns: [{
            label: this.$i18n.t('account.username'),
            field: 'username',
            filterable: true,
          }, {
            label: this.$i18n.t('account.name'),
            field: 'name',
            filterable: true,
          },
          {
            label: this.$i18n.t('account.email'),
            field: 'email',
            filterable: true,
          },
          {
            label: this.$i18n.t('account.type'),
            field: 'type',
            filterable: true,
          },
          {
            label: this.$i18n.t('account.created'),
            field: 'created',
            filterable: false,
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
        newPassword: newPassword,
        confirmPassword: confirmPassword,
        errors: errors
      }
    },
    methods: {
      initNewAccount() {
          this.newObj.uuid = this.generateUUID();
      },
      generateUUID() {
        var d = new Date().getTime();
        var uuid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
          var r = (d + Math.random()*16)%16 | 0;
          d = Math.floor(d/16);
          return (c=='x' ? r : (r&0x3|0x8)).toString(16);
        });
        return uuid;
      },
      execCreate(obj) {
        this.accountCreate(obj, success => {
          $('#ACcreateModall').modal('toggle');
          this.getAll()
        }, error => {
          console.log(error.body.message);
        })
      },
      getAll() {
        this.accountGetAll(success => {
          this.rows = success.body
          this.isLoading = false;
        }, error => {
          console.log(error)
        })
      },

    }
  }

</script>

<style>
  .create-account {
    float: right;
    margin-top: -52px;
    margin-right: 35px;
  }

</style>

