<template>
  <div id="app">
    <!-- login view -->
    <div v-if="!isLogged" class="login-pf">
      <div>
        <span id="badge">
          <img src="/static/logo-light.png" alt=" logo" />
        </span>
        <div class="container">
          <div class="row">
            <div class="col-sm-12 col-brand">
              <div id="brand">
                <h1>{{appName}}</h1>
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
                    <input autocomplete="current-password" required v-model="password" type="password" class="form-control" id="inputPassword"
                      :placeholder="$t('login.insert_password')" tabindex="2">
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
                <strong>{{ $t("login.welcome") }} {{appName}}</strong>
              </p>
            </div>
            <!--/.col-*-->
          </div>
          <!--/.row-->
          <div class="right">
            <a class="unsplash" href="https://unsplash.com/@mojoblogs?utm_medium=referral&amp;utm_campaign=photographer-credit&amp;utm_content=creditBadge"
              target="_blank" rel="noopener noreferrer" title="Download free do whatever you want high-resolution photos from Zara Walker">
              <span style="display:inline-block;padding:2px 3px;">
                <svg xmlns="http://www.w3.org/2000/svg" style="height:12px;width:auto;position:relative;vertical-align:middle;top:-1px;fill:white;"
                  viewBox="0 0 32 32">
                  <title>unsplash-logo</title>
                  <path d="M20.8 18.1c0 2.7-2.2 4.8-4.8 4.8s-4.8-2.1-4.8-4.8c0-2.7 2.2-4.8 4.8-4.8 2.7.1 4.8 2.2 4.8 4.8zm11.2-7.4v14.9c0 2.3-1.9 4.3-4.3 4.3h-23.4c-2.4 0-4.3-1.9-4.3-4.3v-15c0-2.3 1.9-4.3 4.3-4.3h3.7l.8-2.3c.4-1.1 1.7-2 2.9-2h8.6c1.2 0 2.5.9 2.9 2l.8 2.4h3.7c2.4 0 4.3 1.9 4.3 4.3zm-8.6 7.5c0-4.1-3.3-7.5-7.5-7.5-4.1 0-7.5 3.4-7.5 7.5s3.3 7.5 7.5 7.5c4.2-.1 7.5-3.4 7.5-7.5z"></path>
                </svg>
              </span>
              <span style="display:inline-block;padding:2px 3px;">Zara Walker</span>
            </a>
          </div>
          <div>Copyright &copy; {{currentYear()}} {{companyName}}</div>
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
            <img class="navbar-brand-icon" src="/static/logo-light.png" alt="" />
            <p class="navbar-brand-name">{{appName}}</p>
          </a>
        </div>
        <nav class="collapse navbar-collapse">
          <ul class="nav navbar-nav navbar-right navbar-iconic navbar-utility">

            <li class="dropdown">
              <a href="#" class="dropdown-toggle nav-item-iconic" id="dropdownMenu2" data-toggle="dropdown" aria-haspopup="true" aria-expanded="true">
                <p class="login-main-name">{{ user.info.name }}</p>
                <p class="login-main-type">
                  <span v-bind:class="[getLoginIcon(user.info.type), 'login-main-icon']"></span>
                  {{ user.info.type }}
                  <span class="caret"></span>
                </p>
              </a>
              <ul class="dropdown-menu" aria-labelledby="dropdownMenu2">
                <li>
                  <a href="#/profile">{{ $t("menu.profile") }}</a>
                </li>
                <li>
                  <a v-on:click="doLogout()" href="#">Logout</a>
                </li>
              </ul>
            </li>
            <li class="dropdown">
              <a href="#" class="dropdown-toggle nav-item-iconic" id="dropdownMenu1" data-toggle="dropdown" aria-haspopup="true" aria-expanded="true">
                <span :title="$t('help')" class="fa pficon-help"></span>
                <span class="caret"></span>
              </a>
              <ul class="dropdown-menu" aria-labelledby="dropdownMenu1">
                <li>
                  <a target="blank" :href="helpUrl">{{ $t("help") }}</a>
                </li>
                <li>
                  <a href="#" data-toggle="modal" data-target="#about-modal">{{ $t("menu.about") }}</a>
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

          <li v-bind:class="[getCurrentPath('') ? 'active' : '', 'list-group-item']">
            <a href="#/">
              <span class="fa fa-dashboard"></span>
              <span class="list-group-item-value">{{ $t('menu.dashboard') }}</span>
            </a>
          </li>

          <li></li>

          <li v-bind:class="[getCurrentPath('hotspots') ? 'active' : '', 'list-group-item']">
            <a :href="user.info.type == 'admin' || user.info.type == 'reseller' ? '#/hotspots' : '#/hotspots/'+user.info.hotspot_id">
              <span class="fa fa-wifi"></span>
              <span class="list-group-item-value">{{ user.info.type == 'admin' || user.info.type == 'reseller' ? $t('menu.hotspots') : $t('menu.hotspot') }}</span>
            </a>
          </li>
          <li v-bind:class="[getCurrentPath('units') ? 'active' : '', 'list-group-item']" v-if="(user.info.type == 'admin') || (user.info.type == 'reseller') ||  (user.info.type == 'customer')">
            <a href="#/units">
              <span class="fa pficon-connected" data-toggle="tooltip" title="Adipscing"></span>
              <span class="list-group-item-value">{{ $t("menu.units") }}</span>
            </a>
          </li>
          <li v-bind:class="[getCurrentPath('users') ? 'active' : '', 'list-group-item']">
            <a href="#/users">
              <span class="fa fa-users"></span>
              <span class="list-group-item-value">{{ $t("menu.users") }}</span>

            </a>
          </li>
          <li v-bind:class="[getCurrentPath('devices') ? 'active' : '', 'list-group-item']" v-if="(user.info.type == 'admin') || (user.info.type == 'reseller') ||  (user.info.type == 'customer')">
            <a href="#/devices">
              <span class="fa fa-laptop"></span>
              <span class="list-group-item-value">{{ $t("menu.devices") }}</span>

            </a>
          </li>
          <li v-bind:class="[getCurrentPath('sessions') ? 'active' : '', 'list-group-item']" v-if="(user.info.type == 'admin') || (user.info.type == 'reseller') ||  (user.info.type == 'customer')">
            <a href="#/sessions">
              <span class="fa fa-list" data-toggle="tooltip" title="Adipscing"></span>
              <span class="list-group-item-value">{{ $t("menu.sessions") }}</span>
            </a>
          </li>

          <li v-bind:class="[getCurrentPath('reports') ? 'active' : '', 'list-group-item']" v-if="(user.info.type == 'admin')">
            <a href="#/reports">
              <span class="fa fa-area-chart" data-toggle="tooltip" title="Adipscing"></span>
              <span class="list-group-item-value">{{ $t("menu.reports") }}</span>
            </a>
          </li>
          <li></li>

          <li v-bind:class="[getCurrentPath('accounts') ? 'active' : '', 'list-group-item']" v-if="(user.info.type == 'admin') || (user.info.type == 'reseller')">
            <a href="#/accounts">
              <span class="fa pficon-users"></span>
              <span class="list-group-item-value">{{ $t("menu.accounts") }}</span>

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
                    <span class="list-group-item-value">{{ $t("menu.profile") }}</span>
                  </a>
                </li>

                <li class="list-group-item">
                  <a v-on:click="doLogout()" href="#">
                    <span class="list-group-item-value">{{ $t("menu.logout") }}</span>
                  </a>
                </li>
              </ul>
            </div>
          </li>
          <li class="list-group-item secondary-nav-item-pf mobile-nav-item-pf visible-xs-block" data-target="#help-secondary">
            <a href="#">
              <span class="pficon pficon-help" data-toggle="tooltip" title="" data-original-title="Help"></span>
              <span class="list-group-item-value">{{ $t("menu.help") }}</span>
            </a>
            <div id="help-secondary" class="nav-pf-secondary-nav">
              <div class="nav-item-pf-header">
                <a href="#" class="secondary-collapse-toggle-pf" data-toggle="collapse-secondary-nav"></a>
                <span>{{ $t("help") }}</span>
              </div>
              <ul class="list-group">
                <li class="list-group-item">
                  <a target="blank" href="https://github.com/nethesis/icaro">
                    <span class="list-group-item-value">{{ $t("menu.help") }}</span>
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
            <h1>{{appName}}</h1>
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
                <li class="space-line"></li>
                <li>
                  <strong>Privacy</strong>
                  <a target="blank" href="/privacy">Link</a>
                </li>
              </ul>
            </div>
          </div>
          <div class="modal-footer">
            <img class="about-logo" src="/static/logo-light.png" alt="Patternfly Symbol">
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
  import UtilService from './services/util';
  import {
    setTimeout
  } from 'timers';

  export default {
    name: 'app',
    mixins: [LoginService, StorageService, UtilService],
    created() {
      document.title = CONFIG.APP_NAME
    },
    data() {
      // is logged
      var isLogged = false
      var user = {
        login: this.get('loggedUser') || null,
        info: {}
      }

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
        isLogged: isLogged,
        errors: errors,
        appName: CONFIG.APP_NAME,
        helpUrl: CONFIG.HELP_URL,
        companyName: CONFIG.COMPANY_NAME
      }
    },
    methods: {
      currentYear() {
        return new Date().getFullYear()
      },
      getCurrentPath(route) {
        return this.$route.path.split('/')[1] === route
      },
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

<style src="./styles/main.css">

</style>
