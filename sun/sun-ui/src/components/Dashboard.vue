<template>
  <div>
    <h2>{{ msg }}</h2>

    <div class="row row-cards-pf">

      <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller')" class="col-xs-12 col-sm-6 col-md-3 adjust-height">
        <div class="card-pf card-pf-accented card-pf-aggregate-status">
          <h2 class="card-pf-title">
            <a href="#/accounts" class="pficon pficon-users card-link"></a>
            <a href="#/accounts" class="card-link">{{ $t("dashboard.accounts") }}</a>
          </h2>
          <div class="card-pf-body">
            <p class="card-pf-aggregate-status-notifications">
              <span class="card-pf-aggregate-status-notification">
                <div v-if="totals.accounts.isLoading" class="spinner spinner-sm"></div>
                <div v-if="!totals.accounts.isLoading">{{ totals.accounts.count }}</div>
              </span>
            </p>
          </div>
        </div>
      </div>

      <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller')" class="col-xs-12 col-sm-6 col-md-3 adjust-height">
        <div class="card-pf card-pf-accented card-pf-aggregate-status">
          <h2 class="card-pf-title">
            <a href="#/hotspots" class="fa fa-wifi card-link"></a>
            <a href="#/hotspots" class="card-link">{{ $t("dashboard.hotspots") }}</a>
          </h2>
          <div class="card-pf-body">
            <p class="card-pf-aggregate-status-notifications">
              <span class="card-pf-aggregate-status-notification">
                <div v-if="totals.hotspots.isLoading" class="spinner spinner-sm"></div>
                <div v-if="!totals.hotspots.isLoading">{{ totals.hotspots.count }}</div>
              </span>
            </p>
          </div>
        </div>
      </div>

      <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller') || (user.account_type == 'customer')" class="col-xs-12 col-sm-6 col-md-3 adjust-height">
        <div class="card-pf card-pf-accented card-pf-aggregate-status">
          <h2 class="card-pf-title">
            <a href="#/units" class="pficon pficon-connected card-link"></a>
            <a href="#/units" class="card-link">{{ $t("dashboard.units") }}</a>
          </h2>
          <div class="card-pf-body">
            <p class="card-pf-aggregate-status-notifications">
              <span class="card-pf-aggregate-status-notification">
                <div v-if="totals.units.isLoading" class="spinner spinner-sm"></div>
                <div v-if="!totals.units.isLoading">{{ totals.units.count }} <span v-if="isAuth0()">/ <strong class="soft">{{ user.subscription.subscription_plan.max_units }}</strong></span></div>
              </span>
            </p>
          </div>
        </div>
      </div>

      <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller') || (user.account_type == 'customer') || (user.account_type == 'desk')"
        class="col-xs-12 col-sm-6 col-md-3 adjust-height">
        <div class="card-pf card-pf-accented card-pf-aggregate-status">
          <h2 class="card-pf-title">
            <a href="#/users" class="fa fa-users card-link"></a>
            <a href="#/users" class="card-link">{{ $t("dashboard.users") }}</a>
          </h2>
          <div class="card-pf-body">
            <p class="card-pf-aggregate-status-notifications">
              <span class="card-pf-aggregate-status-notification">
                <div v-if="totals.users.isLoading" class="spinner spinner-sm"></div>
                <div v-if="!totals.users.isLoading">{{ totals.users.count }}</div>
              </span>
            </p>
          </div>
        </div>
      </div>

      <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller') || (user.account_type == 'customer') || (user.account_type == 'desk')"
        class="col-xs-12 col-sm-6 col-md-3 adjust-height">
        <div class="card-pf card-pf-accented card-pf-aggregate-status">
          <h2 class="card-pf-title">
            <a href="#/devices" class="fa fa-laptop card-link"></a>
            <a href="#/devices" class="card-link">{{ $t("dashboard.devices") }}</a>
          </h2>
          <div class="card-pf-body">
            <p class="card-pf-aggregate-status-notifications">
              <span class="card-pf-aggregate-status-notification">
                <div v-if="totals.devices.isLoading" class="spinner spinner-sm"></div>
                <div v-if="!totals.devices.isLoading">{{ totals.devices.count }}</div>
              </span>
            </p>
          </div>
        </div>
      </div>

      <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller') || (user.account_type == 'customer')" class="col-xs-12 col-sm-6 col-md-3 adjust-height">
        <div class="card-pf card-pf-accented card-pf-aggregate-status">
          <h2 class="card-pf-title">
            <a href="#/sessions" class="fa fa-list card-link"></a>
            <a href="#/sessions" class="card-link">{{ $t("dashboard.sessions") }}</a>
          </h2>
          <div class="card-pf-body">
            <p class="card-pf-aggregate-status-notifications">
              <span class="card-pf-aggregate-status-notification">
                <div v-if="totals.sessions.isLoading" class="spinner spinner-sm"></div>
                <div v-if="!totals.sessions.isLoading">{{ totals.sessions.count }}</div>
              </span>
            </p>
          </div>
        </div>
      </div>

      <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller')" class="col-xs-12 col-sm-6 col-md-3 adjust-height">
        <div class="card-pf card-pf-accented card-pf-aggregate-status">
          <h2 class="card-pf-title">
            <span class="fa fa-commenting"></span>
            <span class="">{{ $t("dashboard.sms_sent") }}</span>
          </h2>
          <div class="card-pf-body">
            <p class="card-pf-aggregate-status-notifications">
              <span class="card-pf-aggregate-status-notification">
                <div v-if="totals.sms.isLoading" class="spinner spinner-sm"></div>
                <div v-if="!totals.sms.isLoading"><span :class="[totals.sms.count >= totals.sms.max_count ? 'red' : '']">{{ totals.sms.count }}</span> / <strong class="soft"><span :class="[totals.sms.count >= totals.sms.max_count ? 'red' : '']">{{ totals.sms.max_count }}</span></strong></div>
              </span>
            </p>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script>
  import AccountService from "../services/account";
  import StatsService from "../services/stats";
  import StorageService from "../services/storage";

  export default {
    name: "Dashboard",
    mixins: [AccountService, StatsService, StorageService],
    data() {
      // get totals
      this.getTotals();

      // get subscription info
      this.getSubscriptionInfo()

      return {
        msg: "Dashboard",
        totals: {
          accounts: {
            isLoading: true,
            count: 0
          },
          units: {
            isLoading: true,
            count: 0
          },
          hotspots: {
            isLoading: true,
            count: 0
          },
          users: {
            isLoading: true,
            count: 0
          },
          devices: {
            isLoading: true,
            count: 0
          },
          sessions: {
            isLoading: true,
            count: 0
          },
          sms: {
            isLoading: true,
            count: 0,
            max_count: 0
          }
        },
        user: this.get("loggedUser") || null
      };
    },
    methods: {
      isAuth0() {
        return this.get("auth0User");
      },
      getSubscriptionInfo() {
        this.accountGet(this.get("loggedUser").id, success => {
          this.user.subscription = success.body.subscription;
        }, error => {

        });
      },
      getTotals() {
        this.statsHotspotsTotal(
          success => {
            this.totals.hotspots.count = success.body.total;
            this.totals.hotspots.isLoading = false;
          },
          error => {
            console.error(error.body);
            this.totals.hotspots.isLoading = false;
          }
        );
        this.statsUnitsTotal(
          success => {
            this.totals.units.count = success.body.total;
            this.totals.units.isLoading = false;
          },
          error => {
            console.error(error.body);
            this.totals.units.isLoading = false;
          }
        );
        this.statsAccountsTotal(
          success => {
            this.totals.accounts.count = success.body.total;
            this.totals.accounts.isLoading = false;
          },
          error => {
            console.error(error.body);
            this.totals.accounts.isLoading = false;
          }
        );
        this.statsDevicesTotal(
          success => {
            this.totals.devices.count = success.body.total;
            this.totals.devices.isLoading = false;
          },
          error => {
            console.error(error.body);
            this.totals.devices.isLoading = false;
          }
        );
        this.statsUsersTotal(
          success => {
            this.totals.users.count = success.body.total;
            this.totals.users.isLoading = false;
          },
          error => {
            console.error(error.body);
            this.totals.users.isLoading = false;
          }
        );
        this.statsSessionsTotal(
          success => {
            this.totals.sessions.count = success.body.total;
            this.totals.sessions.isLoading = false;
          },
          error => {
            console.error(error.body);
            this.totals.sessions.isLoading = false;
          }
        );
        this.statsSMSTotalForAccount(
        success => {
          this.totals.sms.count = success.body.sms_count;
          this.totals.sms.max_count = success.body.sms_max_count;
          this.totals.sms.isLoading = false;
        },
        error => {
          console.error(error.body);
          this.totals.sms.isLoading = false;
        }
      );
      }
    }
  };

</script>

<style scoped>
  .adjust-height {
    max-height: 109px;
    min-height: 109px;
  }

</style>
