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

      <!-- custom disclaimers -->
      <div v-if="isAdmin && info.data.type == 'reseller'" class="col-xs-12 col-sm-12 col-md-6">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{$t("account.add_disclaimer")}}
              <div
                v-if="!disclaimers.isLoading"
                class="pficon pficon-catalog right"
                data-toggle="tooltip"
                data-placement="left"
              ></div>
              <div v-if="disclaimers.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
          <form role="form" v-on:submit.prevent="addDisclaimer()">
            <div v-if="!disclaimers.isLoading" class="card-pf-body">
              <div class="list-details">
                <dt>{{ $t("account.disclaimers") }}</dt>
              </div>
              <ul v-if="disclaimers.data.length" class="disclaimer-list">
                <li v-for="(disclaimer, dk) in disclaimers.data" :key="dk">
                  <span>{{ disclaimer.title }} ({{ $t("account." + disclaimer.type) }})</span>
                  <button
                    type="button"
                    class="btn btn-danger disclaimer-delete-btn"
                    @click="showDeleteDisclaimerModal(disclaimer)"
                  >{{$t('delete')}}</button>
                </li>
              </ul>
              <div v-else class="no-disclaimer">
                {{$t('account.no_disclaimer')}}
              </div>
              <!-- add new disclaimer -->
              <div class="list-details">
                <dt>{{ $t("account.disclaimer_type") }}</dt>
                <span class="disclaimer-type">
                  <input
                    v-model="disclaimers.newDisclaimer.type"
                    type="radio"
                    id="privacy"
                    name="disclaimer-type"
                    value="privacy"
                  />
                  <label for="privacy">
                    <span>{{ $t("account.privacy") }}</span>
                  </label>
                </span>
                <span class="disclaimer-type">
                  <input
                    v-model="disclaimers.newDisclaimer.type"
                    type="radio"
                    id="tos"
                    name="disclaimer-type"
                    value="tos"
                  />
                  <label for="tos">
                    <span>{{ $t("account.tos") }}</span>
                  </label>
                </span>
              </div>
              <div class="list-details">
                <dt>{{ $t("account.disclaimer_title") }}</dt>
                <input v-model="disclaimers.newDisclaimer.title" required class="form-control">
              </div>
              <div class="list-details">
                <dt>{{ $t("account.disclaimer_body") }}</dt>
                <textarea v-model="disclaimers.newDisclaimer.body" required class="form-control disclaimer-body"></textarea>
              </div>
            </div>
            <div v-if="!disclaimers.isLoading" class="card-pf-footer">
              <div class="dropdown card-pf-time-frame-filter">
                <button class="btn btn-primary">{{$t('account.add')}}</button>
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
          <form role="form" v-on:submit.prevent="updateSMSCount()">
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

    <!-- delete disclaimer modal -->
    <div class="modal fade" id="deleteDisclaimerModal" tabindex="-1" role="dialog" data-backdrop="static">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h4 class="modal-title">{{$t('account.delete_disclaimer')}}</h4>
          </div>
          <form class="form-horizontal" v-on:submit.prevent="deleteDisclaimer()">
            <div class="modal-body">
              <div class="alert alert-warning mg-bottom-20">
                <span class="pficon pficon-warning-triangle-o"></span>
                <span>
                  <strong>{{$t('warning')}}:</strong>
                  {{$t('account.delete_disclaimer_confirm')}}
                  <strong>{{ disclaimerToDelete.title}} ({{ $t("account." + disclaimerToDelete.type) }})</strong>.
                </span>
              </div>
              <label>{{$t('are_you_sure')}}?</label>
            </div>
            <div class="modal-footer">
              <button class="btn btn-default" type="button" @click="hideDeleteDisclaimerModal()">{{$t('cancel')}}</button>
              <button
                class="btn btn-danger"
                type="submit"
              >{{$t('delete')}}</button>
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
import DisclaimerService from "../../services/disclaimer";

export default {
  name: "AccountsDetails",
  mixins: [
    AccountService,
    StatsService,
    StorageService,
    UtilService,
    IntegrationService,
    DisclaimerService
  ],
  components: {
    accountAction: AccountAction
  },
  data() {
    // get account info
    this.getInfo();

    // get sms info
    this.getActualSMSCount();

    // get privacy and tos custom disclaimers
    this.getDisclaimers();

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
      disclaimers: {
        data: {},
        isLoading: true,
        newDisclaimer: {
          title: "",
          type: "privacy",
          body: ""
        },
      },
      integrations: [],
      maps: {},
      isAdmin: this.get("loggedUser").account_type == "admin",
      disclaimerToDelete: {},
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
    getDisclaimers() {
      this.disclaimersByAccount(
        this.$route.params.id,
        success => {
          this.disclaimers.data = success.body;
          this.disclaimers.isLoading = false;
        },
        error => {
          console.error(error.body);
          this.disclaimers.isLoading = false;
          this.disclaimers.data = {};
        }
      );
    },
    addDisclaimer() {
      this.disclaimers.isLoading = true;

      this.disclaimerCreate(
        {
          title: this.disclaimers.newDisclaimer.title,
          type: this.disclaimers.newDisclaimer.type,
          body: this.disclaimers.newDisclaimer.body,
        },
        this.$route.params.id,
        success => {
          this.disclaimers.newDisclaimer = {
            title: "",
            type: "privacy",
            body: ""
          };
          this.disclaimers.isLoading = false;
          this.getDisclaimers();
        },
        error => {
          this.disclaimers.isLoading = false;
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
    },
    showDeleteDisclaimerModal(disclaimer) {
      this.disclaimerToDelete = disclaimer;
      $('#deleteDisclaimerModal').modal("show");
    },
    hideDeleteDisclaimerModal(disclaimer) {
      this.disclaimerToDelete = {};
      $('#deleteDisclaimerModal').modal("hide");
    },
    deleteDisclaimer() {
      this.disclaimerDelete(
        this.disclaimerToDelete.id,
        success => {
          this.getDisclaimers();
        },
        error => {
          console.error(error.body);
        }
      );
      this.disclaimerToDelete = {};
      $('#deleteDisclaimerModal').modal("hide");
    },
  }
};
</script>

<style>
.img-int {
  height: 20px !important;
}

.disclaimer-type {
  margin-right: 2rem;
}

.disclaimer-body {
  min-height: 10rem;
}

.no-disclaimer {
  margin-bottom: 5px;
}

.disclaimer-list {
  list-style-type: none;
  padding: 0;
}

.disclaimer-list li {
  display: inline-block;
  padding: 0.5rem;
  width: 100%;
}

.disclaimer-list li:hover {
  background-color: #C1E6E4;
}

.disclaimer-delete-btn {
  float: right;
}

.mg-bottom-20 {
  margin-bottom: 20px !important;
}
</style>