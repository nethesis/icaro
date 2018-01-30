<template>
  <div id="app">
    <!-- login view -->
    <div v-if="!isLogged" class="login-pf">
      <div>
        <span id="badge">
          <img src="./assets/logo.png" alt=" logo" />
        </span>
        <div class="container">
          <div class="row">
            <div class="col-sm-12 col-brand">
              <div id="brand">
                <h1>Icaro Hotspot Manager</h1>
              </div>
            </div>
            <!--/.col-*-->
            <div class="col-sm-7 col-md-6 col-lg-5 login">
              <form class="form-horizontal" role="form" v-on:submit.prevent="doLogin()">
                <div v-bind:class="[errors.username ? 'has-error' : '', 'form-group']">
                  <label for="inputUsername" class="col-sm-2 col-md-2 control-label">Username</label>
                  <div class="col-sm-10 col-md-10">
                    <input autocomplete="username" required v-model="username" type="text" class="form-control" id="inputUsername" :placeholder="$t('login.insert_username')"
                      tabindex="1">
                    <span v-if="errors.username" class="help-block">{{ $t("login.user_error") }}</span>
                  </div>
                </div>
                <div v-bind:class="[errors.password ? 'has-error' : '', 'form-group']">
                  <label for="inputPassword" class="col-sm-2 col-md-2 control-label">Password</label>
                  <div class="col-sm-10 col-md-10">
                    <input autocomplete="current-password" required v-model="password" type="password" class="form-control" id="inputPassword" :placeholder="$t('login.insert_password')"
                      tabindex="2">
                    <span v-if="errors.password" class="help-block">{{ $t("login.password_error") }}</span>
                  </div>
                </div>
                <div class="form-group">
                  <div class="col-xs-8 col-sm-offset-2 col-sm-6 col-md-offset-2 col-md-6">
                    <span class="help-block"></span>
                  </div>
                  <div class="col-xs-4 col-sm-4 col-md-4 submit">
                    <button type="submit" class="btn btn-primary btn-lg" tabindex="4">Log In</button>
                  </div>
                </div>
              </form>
            </div>
            <!--/.col-*-->
            <div class="col-sm-5 col-md-6 col-lg-7 details">
              <p>
                <strong>{{ $t("login.welcome") }} Icaro Hotspot Manager</strong>
              </p>
            </div>
            <!--/.col-*-->
          </div>
          <!--/.row-->
        </div>
        <!--/.container-->
      </div>
    </div>
    <!-- end login view -->
    <!-- logged view -->
    <div v-if="isLogged" class="cards-pf">
      <!-- top bar -->
      <nav class="navbar navbar-pf-vertical">
        <div class="navbar-header">
          <button type="button" class="navbar-toggle">
            <span class="sr-only"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
          </button>
          <a href="/" class="navbar-brand">
            <img class="navbar-brand-icon" src="./assets/logo.png" alt="" />
            <p class="navbar-brand-name">Icaro Hotspot Manager</p>
          </a>
        </div>
        <nav class="collapse navbar-collapse">
          <ul class="nav navbar-nav navbar-right navbar-iconic navbar-utility">

            <li class="dropdown">
              <a href="#" class="dropdown-toggle nav-item-iconic" id="dropdownMenu2" data-toggle="dropdown" aria-haspopup="true" aria-expanded="true">
                <p class="login-main-name">{{ user.info.name }}</p>
                <p class="login-main-type">
                  <span v-bind:class="[getLoginIcon(), 'login-main-icon']"></span>
                  {{ user.info.type }}
                  <span class="caret"></span>
                </p>
              </a>
              <ul class="dropdown-menu" aria-labelledby="dropdownMenu2">
                <li>
                  <a href="#/profile">{{ $t("dashboard.profile") }}</a>
                </li>
                <li>
                  <a v-on:click="doLogout()" href="#">Logout</a>
                </li>
              </ul>
            </li>
            <li class="dropdown">
              <a href="#" class="dropdown-toggle nav-item-iconic" id="dropdownMenu1" data-toggle="dropdown" aria-haspopup="true" aria-expanded="true">
                <span :title="$t('dashboard.help')" class="fa pficon-help"></span>
                <span class="caret"></span>
              </a>
              <ul class="dropdown-menu" aria-labelledby="dropdownMenu1">
                <li>
                  <a target="blank" href="https://github.com/nethesis/icaro">{{ $t("dashboard.help") }}</a>
                </li>
                <li>
                  <a href="#" data-toggle="modal" data-target="#about-modal">{{ $t("dashboard.about") }}</a>
                </li>
              </ul>
            </li>
          </ul>
        </nav>

      </nav>
      <!-- end top bar -->

      <!-- left menu -->
      <div class="nav-pf-vertical nav-pf-vertical-with-sub-menus nav-pf-persistent-secondary">
        <ul class="list-group">

          <li v-bind:class="[currentPath == '/' ? 'active' : '', 'list-group-item']">
            <a href="#/">
              <span class="fa fa-dashboard" data-toggle="tooltip" title="Dashboard"></span>
              <span class="list-group-item-value">Dashboard</span>
            </a>
          </li>
          <li v-bind:class="[currentPath == '/hotspots' ? 'active' : '', 'list-group-item']">
            <a href="#/hotspots">
              <span class="fa fa-wifi" data-toggle="tooltip" title="Dolor"></span>
              <span class="list-group-item-value">Hotspots</span>

            </a>
          </li>
          <li v-bind:class="[currentPath == '/users' ? 'active' : '', 'list-group-item']">
            <a href="#/users">
              <span class="fa fa-users" data-toggle="tooltip" title="Dolor"></span>
              <span class="list-group-item-value">{{ $t("dashboard.users") }}</span>

            </a>
          </li>
          <li v-bind:class="[currentPath == '/reports' ? 'active' : '', 'list-group-item']">
            <a href="#/reports">
              <span class="fa fa-list" data-toggle="tooltip" title="Adipscing"></span>
              <span class="list-group-item-value">Report</span>
            </a>
          </li>

          <li></li>

          <li v-bind:class="[currentPath == '/accounts' ? 'active' : '', 'list-group-item']">
            <a href="#/accounts">
              <span class="fa pficon-users" data-toggle="tooltip" title="Dolor"></span>
              <span class="list-group-item-value">Accounts</span>

            </a>
          </li>
          <li v-bind:class="[currentPath == '/preferences' ? 'active' : '', 'list-group-item']">
            <a href="#/preferences">
              <span class="fa fa-gear" data-toggle="tooltip" title="Lorem"></span>
              <span class="list-group-item-value">{{ $t("dashboard.preferences") }}</span>
            </a>
          </li>

          <li class="list-group-item secondary-nav-item-pf mobile-nav-item-pf visible-xs-block" data-target="#user-secondary">
            <a href="#">
              <span class="pficon pficon-user" data-toggle="tooltip" title="" data-original-title="User"></span>
              <span class="list-group-item-value">{{ user.info.name }}</span>
            </a>
            <div id="user-secondary" class="nav-pf-secondary-nav">
              <div class="nav-item-pf-header">
                <a href="#" class="secondary-collapse-toggle-pf" data-toggle="collapse-secondary-nav"></a>
                <span>{{ user.info.name }}</span>
              </div>

              <ul class="list-group">
                <li class="list-group-item">
                  <a href="#/profile">
                    <span class="list-group-item-value">{{ $t("dashboard.profile") }}</span>
                  </a>
                </li>

                <li class="list-group-item">
                  <a v-on:click="doLogout()" href="#">
                    <span class="list-group-item-value">Logout</span>
                  </a>
                </li>
              </ul>
            </div>
          </li>
          <li class="list-group-item secondary-nav-item-pf mobile-nav-item-pf visible-xs-block" data-target="#help-secondary">
            <a href="#">
              <span class="pficon pficon-help" data-toggle="tooltip" title="" data-original-title="Help"></span>
              <span class="list-group-item-value">{{ $t("dashboard.help") }}</span>
            </a>
            <div id="help-secondary" class="nav-pf-secondary-nav">
              <div class="nav-item-pf-header">
                <a href="#" class="secondary-collapse-toggle-pf" data-toggle="collapse-secondary-nav"></a>
                <span>{{ $t("dashboard.help") }}</span>
              </div>
              <ul class="list-group">
                <li class="list-group-item">
                  <a target="blank" href="https://github.com/nethesis/icaro">
                    <span class="list-group-item-value">{{ $t("dashboard.help") }}</span>
                  </a>
                </li>
              </ul>
            </div>
          </li>

        </ul>

      </div>
      <!-- end left menu -->

      <!-- main view -->
      <div class="container-fluid container-cards-pf container-pf-nav-pf-vertical nav-pf-persistent-secondary">
        <router-view></router-view>
      </div>
      <!-- end main view -->
    </div>
    <!-- end logged view -->

    <!-- modals -->
    <div class="modal fade" id="about-modal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content about-modal-pf">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
          </div>
          <div class="modal-body">
            <h1>Icaro Hotspot Manager</h1>
            <div class="product-versions-pf">
              <ul class="list-unstyled">
                <li>
                  <strong>Sun UI</strong>
                  <a target="blank" href="https://github.com/nethesis/icaro/tree/master/sun/sun-ui">GitHub</a>
                </li>
                <li>
                  <strong>Sun API</strong>
                  <a target="blank" href="https://github.com/nethesis/icaro/tree/master/sun/sun-api">GitHub</a>
                </li>
                <li>
                  <strong>Wax</strong>
                  <a target="blank" href="https://github.com/nethesis/icaro/tree/master/wax">GitHub</a>
                </li>
                <li>
                  <strong>Wings</strong>
                  <a target="blank" href="https://github.com/nethesis/icaro/tree/master/wings">GitHub</a>
                </li>
                <li>
                  <strong>Dedalo</strong>
                  <a target="blank" href="https://github.com/nethesis/icaro/tree/master/dedalo">GitHub</a>
                </li>
              </ul>
            </div>
          </div>
          <div class="modal-footer">
            <img class="about-logo" src="./assets/logo.png" alt="Patternfly Symbol">
          </div>
        </div>
      </div>
    </div>
    <!-- end modals -->
  </div>
</template>

<script>
  import LoginService from './services/login';
  import StorageService from './services/storage';
  import {
    setTimeout
  } from 'timers';

  export default {
    name: 'app',
    mixins: [LoginService, StorageService],
    data() {
      // is logged
      var isLogged = false
      var user = {
        login: this.get('loggedUser') || null,
        info: {}
      }

      // get current path to highlight menu item
      var currentPath = this.$route.path

      var errors = {
        username: false,
        password: false,
      }

      if (user.login) {
        this.getInfo(user.login.id, response => {
          if (response) {
            this.user.info = response
            this.isLogged = true
            this.initGraphics()
          } else {
            this.isLogged = false
            this.showBody()
          }
        })
      } else {
        this.showBody()
      }

      return {
        username: '',
        password: '',
        user: user,
        currentPath: currentPath,
        isLogged: isLogged,
        errors: errors
      }
    },
    methods: {
      getLoginIcon() {
        var icon = 'fa fa-user'
        switch (this.user.info.type) {
          case 'admin':
            icon = 'fa fa-graduation-cap'
            break;
          case 'reseller':
            icon = 'fa fa-user'
            break;
          case 'customer':
            icon = 'fa fa-briefcase'
            break;
          case 'desk':
            icon = 'fa fa-coffee'
            break;
        }
        return icon
      },
      getInfo(id, callback) {
        this.execGetInfo(id, success => {
          callback(success.body)
        }, error => {
          callback(null)
        })
      },
      doLogin() {
        this.errors.username = false
        this.errors.password = false
        this.execLogin({
          username: this.username,
          password: this.password
        }, success => {
          // extract loggedUser info
          var loggedUser = success.body
          // save to localstorage
          this.set('loggedUser', loggedUser)

          // get user info
          this.getInfo(loggedUser.id, response => {
            if (response) {
              this.user.info = response
              this.isLogged = true
              this.initGraphics()
            } else {
              this.isLogged = false
            }
          })

          // change route
          this.isLogged = true
          this.$router.push('/')
        }, error => {
          if (error.body.message == 'No username found!') {
            this.errors.username = true
          }
          if (error.body.message == 'Password is invalid') {
            this.errors.password = true
          }
          console.log(error.body.message);
        })
      },
      doLogout() {
        this.execLogout(success => {
          // save to localstorage
          this.delete('loggedUser')

          // change route
          this.isLogged = false
          this.$router.push('/')
          this.resetGraphics()
        }, error => {
          console.log(error.body.message);
        })
      },
      initGraphics() {
        $('body').addClass('logged')
        $('body').removeClass('not-logged')
        setTimeout(function () {
          $().setupVerticalNavigation(true);
        }, 1000);
        this.showBody()
      },
      resetGraphics() {
        $('body').addClass('not-logged')
        $('body').removeClass('logged')
        window.location.reload()
      },
      showBody() {
        $('body').show()
        $('body').addClass('not-logged')
      }
    }
  }

</script>

<style>
  body {
    overflow: hidden !important;
    display: none;
  }

  h2 {
    padding-left: 15px;
    padding-right: 15px;
    padding-bottom: 15px;
  }

  a {
    text-decoration: none !important;
  }

  .not-logged {
    background-color: #010101;
  }

  .logged {
    background-color: #f5f5f5
  }

  #brand {
    top: -100px !important;
  }

  .has-error .checkbox,
  .has-error .checkbox-inline,
  .has-error .control-label,
  .has-error .help-block,
  .has-error .radio,
  .has-error .radio-inline,
  .has-error.checkbox label,
  .has-error.checkbox-inline label,
  .has-error.radio label,
  .has-error.radio-inline label {
    color: #c00 !important;
  }

  .col-brand {
    height: 15px !important;
  }

  .container-fluid {
    padding: 0px !important;
  }

  .container-cards-pf {
    margin-top: 0px !important;
  }

  .list-group-item {
    text-decoration: none !important;
  }

  .login-main-name {
    margin-top: -10px !important;
    font-size: 16px !important;
  }

  .login-main-icon {
    margin-right: 6px !important;
  }

  .login-main-type {
    margin-top: -4px !important;
    text-align: right !important;
  }

  .about-logo {
    height: 50px !important;
  }

  .modal-footer {
    margin-top: 0px !important;
  }

  .alert {
    margin-bottom: 0px !important;
  }

</style>
