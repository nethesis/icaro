<template>
  <div>
    <h2>Manager
      <strong class="soft">{{ info.data.username }}</strong>
    </h2>

    <div class="row row-cards-pf">
      <div class="col-xs-12 col-sm-12 col-md-6">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{info.data.name}}
              <div v-if="!info.isLoading" :class="[getLoginIcon(info.data.type), 'right']" data-toggle="tooltip" data-placement="left"
                :title="$t(info.data.type)"></div>
              <div v-if="info.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
          <div v-if="!info.isLoading" class="card-pf-body">
            <div v-if="info.data.hotspot_id" class="list-details">
              <dt>{{ $t("account.hotspot") }}</dt>
              <dd>
                <a :href="'#/hotspots/' + info.data.hotspot_id">{{info.data.hotspot_name}}</a>
              </dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("account.email") }}</dt>
              <dd>{{info.data.email}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("account.uuid") }}</dt>
              <dd>{{info.data.uuid}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("account.created") }}</dt>
              <dd>{{info.data.created | formatDate}}</dd>
            </div>
          </div>
          <div v-if="!info.isLoading" class="card-pf-footer">
            <div class="dropdown card-pf-time-frame-filter">
              <account-action details="false" :obj="info.data" :update="getInfo"></account-action>
            </div>
            <p>
              <a href="#" class="card-pf-link-with-icon">
              </a>
            </p>
          </div>
        </div>
      </div>

      <div v-if="isAdmin" class="col-xs-12 col-sm-12 col-md-6">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{$t("account.add_sms")}}
              <div v-if="!sms.isLoading" class="fa fa-commenting right" data-toggle="tooltip" data-placement="left"></div>
              <div v-if="sms.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
          <form role="form" v-on:submit.prevent="updateSMSCount()">
            <div v-if="!sms.isLoading" class="card-pf-body">
              <div class="list-details">
                <dt>{{ $t("account.sms_sent") }}</dt>
                <dd>{{sms.data.sms_count}}</dd>
              </div>
              <div class="list-details">
                <dt>{{ $t("account.sms_total") }}</dt>
                <dd>{{sms.data.sms_max_count}}</dd>
              </div>
              <div class="list-details">
                <dt>{{ $t("account.sms_to_add") }}</dt>
                <input v-model="sms.data.sms_to_add" required class="form-control">
              </div>
            </div>
            <div v-if="!sms.isLoading" class="card-pf-footer">
              <div class="dropdown card-pf-time-frame-filter">
                <button class="btn btn-primary">{{$t('update')}}</button>
              </div>
              <p>
                <a href="#" class="card-pf-link-with-icon">
                </a>
              </p>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import AccountService from "../../services/account";
import StorageService from "../../services/storage";
import StatsService from "../../services/stats";
import UtilService from "../../services/util";
import AccountAction from "../../directives/AccountAction.vue";

export default {
  name: "AccountsDetails",
  mixins: [AccountService, StatsService, StorageService, UtilService],
  components: {
    accountAction: AccountAction
  },
  data() {
    // get account info
    this.getInfo();

    // get sms info
    this.getActualSMSCount();

    return {
      info: {
        isLoading: true,
        data: {}
      },
      sms: {
        isLoading: true,
        data: {}
      },
      isAdmin: this.get("loggedUser").account_type == "admin"
    };
  },
  // enable tooltips after rendering
  updated: function() {
    $('[data-toggle="tooltip"]').tooltip();
  },
  methods: {
    getInfo() {
      this.accountGet(
        this.$route.params.id,
        success => {
          this.info.data = success.body;
          this.info.isLoading = false;
        },
        error => {
          this.info.isLoading = false;
          console.error(error.body);
        }
      );
    },
    getActualSMSCount() {
      this.statsSMSTotalForAccountByAccount(
        this.$route.params.id,
        success => {
          this.sms.isLoading = false;
          this.sms.data = success.body;
        },
        error => {
          this.sms.isLoading = false;
          console.error(error.body);
        }
      );
    },
    updateSMSCount() {
      this.sms.isLoading = true;

      this.updateSMSTotalForAccountByAccount(
        {
          sms_to_add: parseInt(this.sms.data.sms_to_add) || 0
        },
        this.$route.params.id,
        success => {
          this.sms.isLoading = false;
          this.sms.data = success.body;
          this.getActualSMSCount()
        },
        error => {
          this.sms.isLoading = false;
          console.error(error.body);
        }
      );
    }
  }
};
</script>