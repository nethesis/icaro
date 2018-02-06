<template>
  <div>
    <h2>{{ msg }}</h2>

    <div class="row row-cards-pf">

      <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller')" class="col-xs-12 col-sm-6 col-md-3">
        <div class="card-pf card-pf-accented card-pf-aggregate-status">
          <h2 class="card-pf-title">
            <span class="fa pficon-users"></span>
            {{ $t("dashboard.accounts") }}
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

      <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller')" class="col-xs-12 col-sm-6 col-md-3">
        <div class="card-pf card-pf-accented card-pf-aggregate-status">
          <h2 class="card-pf-title">
            <span class="fa fa-wifi"></span>
            {{ $t("dashboard.hotspots") }}
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

      <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller') || (user.account_type == 'customer')" class="col-xs-12 col-sm-6 col-md-3">
        <div class="card-pf card-pf-accented card-pf-aggregate-status">
          <h2 class="card-pf-title">
            <span class="fa pficon-connected"></span>
            {{ $t("dashboard.units") }}
          </h2>
          <div class="card-pf-body">
            <p class="card-pf-aggregate-status-notifications">
              <span class="card-pf-aggregate-status-notification">
                <div v-if="totals.units.isLoading" class="spinner spinner-sm"></div>
                <div v-if="!totals.units.isLoading">{{ totals.units.count }}</div>
              </span>
            </p>
          </div>
        </div>
      </div>

      <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller') || (user.account_type == 'customer') || (user.account_type == 'desk')"
        class="col-xs-12 col-sm-6 col-md-3">
        <div class="card-pf card-pf-accented card-pf-aggregate-status">
          <h2 class="card-pf-title">
            <span class="fa fa-users"></span>
            {{ $t("dashboard.users") }}
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
        class="col-xs-12 col-sm-6 col-md-3">
        <div class="card-pf card-pf-accented card-pf-aggregate-status">
          <h2 class="card-pf-title">
            <span class="fa fa-mobile"></span>
            {{ $t("dashboard.devices") }}
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

      <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller') || (user.account_type == 'customer')" class="col-xs-12 col-sm-6 col-md-3">
        <div class="card-pf card-pf-accented card-pf-aggregate-status">
          <h2 class="card-pf-title">
            <span class="fa fa-list"></span>
            {{ $t("dashboard.sessions") }}
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

    </div>
  </div>
</template>

<script>
  import StatsService from '../services/stats';
  import StorageService from '../services/storage';

  export default {
    name: 'Dashboard',
    mixins: [StatsService, StorageService],
    data() {

      // get totals
      this.getTotals();

      return {
        msg: 'Dashboard',
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
        },
        user: this.get('loggedUser') || null,
      }
    },
    methods: {
      getTotals() {
        this.statsHotspotsTotal(success => {
          this.totals.hotspots.count = success.body.total
          this.totals.hotspots.isLoading = false
        }, error => {
          console.log(error.body)
          this.totals.hotspots.isLoading = false
        })
        this.statsUnitsTotal(success => {
          this.totals.units.count = success.body.total
          this.totals.units.isLoading = false
        }, error => {
          console.log(error.body)
          this.totals.units.isLoading = false
        })
        this.statsAccountsTotal(success => {
          this.totals.accounts.count = success.body.total
          this.totals.accounts.isLoading = false
        }, error => {
          console.log(error.body)
          this.totals.accounts.isLoading = false
        })
        this.statsDevicesTotal(success => {
          this.totals.devices.count = success.body.total
          this.totals.devices.isLoading = false
        }, error => {
          console.log(error.body)
          this.totals.devices.isLoading = false
        })
        this.statsUsersTotal(success => {
          this.totals.users.count = success.body.total
          this.totals.users.isLoading = false
        }, error => {
          console.log(error.body)
          this.totals.users.isLoading = false
        })
        this.statsSessionsTotal(success => {
          this.totals.sessions.count = success.body.total
          this.totals.sessions.isLoading = false
        }, error => {
          console.log(error.body)
          this.totals.sessions.isLoading = false
        })
      },

    }
  }

</script>

<style scoped>


</style>
