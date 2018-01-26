<template>
  <div id="app">
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
          <img class="navbar-brand-icon" src="src/assets/logo.png" alt="" />
          <p class="navbar-brand-name">Icaro Hotspot Manager</p>
        </a>
      </div>
      <nav class="collapse navbar-collapse">
        <ul class="nav navbar-nav navbar-right navbar-iconic navbar-utility">

          <li class="dropdown">
            <a href="#" class="dropdown-toggle nav-item-iconic" id="dropdownMenu2" data-toggle="dropdown" aria-haspopup="true" aria-expanded="true">
              <span title="Username" class="fa pficon-user"></span>
              Brian Johnson
              <span class="caret"></span>
            </a>
            <ul class="dropdown-menu" aria-labelledby="dropdownMenu2">
              <li>
                <a href="#/profile">{{ $t("dashboard.profile") }}</a>
              </li>
              <li>
                <a href="#">Logout</a>
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
                <a href="#">{{ $t("dashboard.help") }}</a>
              </li>
              <li>
                <a href="#">{{ $t("dashboard.about") }}</a>
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
            <span class="list-group-item-value">Brian Johnson</span>
          </a>
          <div id="user-secondary" class="nav-pf-secondary-nav">
            <div class="nav-item-pf-header">
              <a href="#" class="secondary-collapse-toggle-pf" data-toggle="collapse-secondary-nav"></a>
              <span>Brian Johnson</span>
            </div>

            <ul class="list-group">
              <li class="list-group-item">
                <a href="#/profile">
                  <span class="list-group-item-value">{{ $t("dashboard.profile") }}</span>
                </a>
              </li>

              <li class="list-group-item">
                <a href="#">
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
                <a href="#">
                  <span class="list-group-item-value">{{ $t("dashboard.help") }}</span>
                </a>
              </li>
              <li class="list-group-item">
                <a href="#">
                  <span class="list-group-item-value">{{ $t("dashboard.about") }}</span>
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
</template>

<script>
  window.$ = window.jQuery = require('jquery')
  require('bootstrap')
  require('../node_modules/patternfly/dist/js/patternfly')

  export default {
    name: 'app',
    mounted: function () {
      // init vertical menu setup
      $().setupVerticalNavigation(true);
    },
    data() {
      // get current path to highlight menu item
      var currentPath = this.$route.path

      return {
        currentPath: currentPath
      }
    }
  }
</script>

<style>
  @import url('../node_modules/patternfly/dist/css/patternfly.min.css');
  @import url('../node_modules/patternfly/dist/css/patternfly-additions.min.css');

  h2 {
    padding-left: 15px;
    padding-right: 15px;
    padding-bottom: 15px;
  }

  a {
    text-decoration: none !important;
  }

  .container-fluid {
    padding: 0px;
  }

  .container-cards-pf {
    margin-top: 0px;
  }

  .list-group-item {
    text-decoration: none !important;
  }
</style>