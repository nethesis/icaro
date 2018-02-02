<template>
  <div>
    <h2>{{ msg }}</h2>

    <div class="row row-cards-pf">
      <div class="col-xs-12 col-sm-6 col-md-3">
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
      <div class="col-xs-12 col-sm-6 col-md-3">
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

    </div>
  </div>
</template>

<script>
  import HotspotService from '../services/hotspot';
  import AccountService from '../services/account';
  import StorageService from '../services/storage';

  export default {
    name: 'Dashboard',
    mixins: [HotspotService, AccountService, StorageService],
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
          hotspots: {
            isLoading: true,
            count: 0
          }
        }
      }
    },
    methods: {
      getTotals() {
        this.hotspotGetAll(success => {
          this.totals.hotspots.count = success.body.length
          this.totals.hotspots.isLoading = false
        }, error => {
          console.log(error.body)
          this.totals.hotspots.isLoading = false
        })
        this.accountGetAll(success => {
          this.totals.accounts.count = success.body.length
          this.totals.accounts.isLoading = false
        }, error => {
          console.log(error.body)
          this.totals.accounts.isLoading = false
        })
      },

    }
  }
</script>

<style scoped>

</style>
