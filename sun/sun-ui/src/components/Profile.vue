<template>
  <div>
    <h2>{{ msg }}</h2>

    <div class="row row-cards-pf">
      <div class="col-xs-12 col-sm-12 col-md-6">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{ $t("profile.info") }}
            </h2>
          </div>
          <div class="card-pf-body">
            <p>
              <dt>{{ $t("profile.username") }}</dt>
              <dd>{{user.info.username}}</dd>
            </p>
            <p>
              <dt>{{ $t("profile.name") }}</dt>
              <dd>{{user.info.name}}</dd>
            </p>
            <p>
              <dt>{{ $t("profile.email") }}</dt>
              <dd>{{user.info.email}}</dd>
            </p>
            <p>
              <dt>{{ $t("profile.type") }}</dt>
              <dd>{{user.info.type}}</dd>
            </p>
          </div>
          <div class="card-pf-footer">
            <div class="dropdown card-pf-time-frame-filter">
              <button class="btn btn-default" data-toggle="modal" data-target="#changePassModal">{{ $t("profile.change_password") }}</button>
            </div>
            <p>
              <a href="#" class="card-pf-link-with-icon">

              </a>
            </p>
          </div>
        </div>

      </div>
    </div>

    <div class="modal fade" id="changePassModal" tabindex="-1" role="dialog" aria-labelledby="changePassModalLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="changePassModalLabel">{{ $t("profile.change_password") }}</h4>
          </div>
          <form class="form-horizontal" role="form" v-on:submit.prevent="changePassword()">
            <div class="modal-body">
              <div class="form-group">
                <label class="col-sm-4 control-label" for="textInput-modal-markup">{{ $t("profile.new_password") }}</label>
                <div class="col-sm-8">
                  <input v-model="newPassword" required type="text" id="textInput-modal-markup" class="form-control" :placeholder="$t('profile.new_password')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="textInput2-modal-markup">{{ $t("profile.confirm_password") }}</label>
                <div class="col-sm-8">
                  <input v-model="confirmPassword" required type="text" id="textInput2-modal-markup" class="form-control" :placeholder="$t('profile.confirm_password')">
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
                <strong>{{ $t("profile.change_password_error_title") }}</strong>. {{ $t("profile.change_password_error_sub") }}.
              </div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-default" data-dismiss="modal">{{ $t("cancel") }}</button>
              <button :disabled="newPassword != confirmPassword" type="submit" class="btn btn-primary">{{ $t("update") }}</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import LoginService from '../services/login';
  import StorageService from '../services/storage';

  export default {
    name: 'Profile',
    mixins: [LoginService, StorageService],
    data() {
      var user = {
        login: this.get('loggedUser') || null,
        info: this.$parent.user.info
      }
      var newPassword, confirmPassword = ''

      var errors = {
        password: false
      }

      return {
        msg: 'Profile',
        user: user,
        newPassword: newPassword,
        confirmPassword: confirmPassword,
        errors: errors
      }
    },
    methods: {
      changePassword() {
        this.execChangePassword(this.newPassword, this.user.login.id, success => {
          $('#changePassModal').modal('toggle');
        }, error => {
          this.errors.password = true
          console.log(error.body.message);
        })
      }
    }
  }
</script>

<style scoped>
  .pass-confirm-err {
    height: 5px;
    background: #cc0000;
    width: 100px;
    display: block;
  }

  .pass-confirm-ok {
    height: 5px;
    background: #3f9c35;
    width: 100px;
    display: block;
  }
</style>