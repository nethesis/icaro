<template>
  <div>
    <h2>
      Manager
      <strong class="soft">{{ info.data.username }}</strong>
    </h2>

    <div class="row row-cards-pf">
      <div class="col-xs-12 col-sm-12 col-md-6">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{info.data.name}}
              <div
                v-if="!info.isLoading"
                :class="[getLoginIcon(info.data.type), 'right']"
                data-toggle="tooltip"
                data-placement="left"
                :title="$t(info.data.type)"
              ></div>
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
              <a href="#" class="card-pf-link-with-icon"></a>
            </p>
          </div>
        </div>
      </div>

      <div v-if="isAdmin" class="col-xs-12 col-sm-12 col-md-6">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{$t("account.add_sms")}}
              <div
                v-if="!sms.isLoading"
                class="fa fa-commenting right"
                data-toggle="tooltip"
                data-placement="left"
              ></div>
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
                <a href="#" class="card-pf-link-with-icon"></a>
              </p>
            </div>
          </form>
        </div>
      </div>

      <div v-if="isAdmin" class="col-xs-12 col-sm-12 col-md-6">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{$t("account.add_whatsapp")}}
              <div
                v-if="!whatsapp.isLoading"
                class="fa fa-whatsapp right"
                data-toggle="tooltip"
                data-placement="left"
              ></div>
              <div v-if="whatsapp.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
          <form role="form" v-on:submit.prevent="updateWhatsappCount()">
            <div v-if="!whatsapp.isLoading" class="card-pf-body">
              <div class="list-details">
                <dt>{{ $t("account.whatsapp_sent") }}</dt>
                <dd>{{whatsapp.data.whatsapp_count}}</dd>
              </div>
              <div class="list-details">
                <dt>{{ $t("account.whatsapp_total") }}</dt>
                <dd>{{whatsapp.data.whatsapp_max_count}}</dd>
              </div>
              <div class="list-details">
                <dt>{{ $t("account.whatsapp_to_add") }}</dt>
                <input v-model="whatsapp.data.whatsapp_to_add" required class="form-control">
              </div>
            </div>
            <div v-if="!whatsapp.isLoading" class="card-pf-footer">
              <div class="dropdown card-pf-time-frame-filter">
                <button class="btn btn-primary">{{$t('update')}}</button>
              </div>
              <p>
                <a href="#" class="card-pf-link-with-icon"></a>
              </p>
            </div>
          </form>
        </div>
      </div>

      <div
        v-if="isAdmin && info.data.type == 'reseller'"
        v-for="(i,ik) in integrations"
        :key="ik"
        class="col-xs-12 col-sm-12 col-md-6"
      >
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{$t('account.integration')}}:
              <b>{{i.name}}</b>
              <div class="right">
                <img class="img-int" :src="i.logo">
              </div>
            </h2>
          </div>
          <form role="form">
            <div class="card-pf-body">
              <div class="list-details">
                <dt>{{ $t("account.description") }}</dt>
                <dd>{{i.description}}</dd>
              </div>
              <div class="list-details">
                <dt>{{ $t("account.site") }}</dt>
                <dd>
                  <a target="_blank" :href="i.site">{{i.site}}</a>
                </dd>
              </div>
            </div>
            <div class="card-pf-footer">
              <div class="dropdown card-pf-time-frame-filter">
                <button
                  :class="['btn', maps[i.id] ? 'btn-default' : 'btn-primary']"
                  @click="maps[i.id] ? deleteIntegration(i.id) : createIntegration(i.id)"
                >{{maps[i.id] ? $t('account.disable') : $t('account.enable')}}</button>
              </div>
              <p>
                <a href="#" class="card-pf-link-with-icon"></a>
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
import IntegrationService from "../../services/integration";
import AccountAction from "../../directives/AccountAction.vue";

export default {
  name: "AccountsDetails",
  mixins: [
    AccountService,
    StatsService,
    StorageService,
    UtilService,
    IntegrationService
  ],
  components: {
    accountAction: AccountAction
  },
  data() {
    // get account info
    this.getInfo();

    // get sms info
    this.getActualSMSCount();

    // get whatsapp info
    this.getActualWhatsappCount();

    // integrations
    this.getIntegrations();
    this.getMaps();

    return {
      info: {
        isLoading: true,
        data: {}
      },
      sms: {
        isLoading: true,
        data: {}
      },
      whatsapp: {
        isLoading: true,
        data: {}
      },
      integrations: [],
      maps: {},
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
    getIntegrations() {
      this.integrationGetAll(
        success => {
          this.integrations = success.body;
        },
        error => {
          console.error(error);
          this.integrations = [];
        }
      );
    },
    getMaps() {
      this.mapAccountsGetAll(
        this.$route.params.id,
        success => {
          var result = success.body;
          var maps = {};

          for (var i in result) {
            maps[result[i].integration_id] = true;
          }

          this.maps = maps;
        },
        error => {
          console.error(error);
          this.maps = [];
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
    getActualWhatsappCount() {
      this.statsWhatsappTotalForAccountByAccount(
        this.$route.params.id,
        success => {
          this.whatsapp.isLoading = false;
          this.whatsapp.data = success.body;
        },
        error => {
          this.whatsapp.isLoading = false;
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
          this.getActualSMSCount();
        },
        error => {
          this.sms.isLoading = false;
          console.error(error.body);
        }
      );
    },
    updateWHatsappCount() {
      this.whatsapp.isLoading = true;

      this.updateWhatsappTotalForAccountByAccount(
        {
          whatsapp_to_add: parseInt(this.whatsapp.data.whatsapp_to_add) || 0
        },
        this.$route.params.id,
        success => {
          this.whatsapp.isLoading = false;
          this.whatsapp.data = success.body;
          this.getActualWhatsappCount();
        },
        error => {
          this.whatsapp.isLoading = false;
          console.error(error.body);
        }
      );
    },
    createIntegration(integrationId) {
      this.mapAccountsCreate(
        this.$route.params.id,
        integrationId,
        success => {
          this.getMaps();
        },
        error => {
          console.error(error.body);
        }
      );
    },
    deleteIntegration(integrationId) {
      this.mapAccountsDelete(
        this.$route.params.id,
        integrationId,
        success => {
          this.getMaps();
        },
        error => {
          console.error(error.body);
        }
      );
    }
  }
};
</script>

<style>
.img-int {
  height: 20px !important;
}
</style>