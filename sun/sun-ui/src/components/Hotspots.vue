<template>
  <div>
    <h2>{{ msg }}</h2>
    <div v-if="isLoading" class="spinner spinner-lg"></div>
    <div v-if="rows.length > 0 && !privacyFilled()" class="alert alert-warning alert-dismissable create-hotspot-warning">
      <span class="pficon pficon-warning-triangle-o"></span>
      <strong>{{ $t('hotspot.warning_privacy') }}!</strong> {{ $t('hotspot.warning_privacy_text') }} <a href="#/profile" class="alert-link">{{ $t('hotspot.profile_link') }}</a>
    </div>
    <button
      v-if="rows.length > 0 && !isLoading && !isAdmin"
      data-toggle="modal"
      data-target="#HScreateModal"
      class="btn btn-primary btn-lg create-hotspot"
      :disabled="!privacyFilled()"
    >{{ $t('hotspot.create_new') }}</button>
    <div v-if="rows.length == 0 && !isLoading && !isAdmin" class="blank-slate-pf" id>
      <div class="blank-slate-pf-icon">
        <span class="fa fa-wifi"></span>
      </div>
      <h1>{{ $t('hotspot.no_hotspot_found') }}</h1>
      <p>{{ $t('hotspot.no_hotspot_found_sub') }}.</p>
      <div class="blank-slate-pf-main-action">
        <button
          data-toggle="modal"
          data-target="#HScreateModal"
          class="btn btn-primary btn-lg"
          :disabled="!privacyFilled()"
        >{{ $t('hotspot.create_new') }}</button>
      </div>
      <br/>
      <div v-if="!privacyFilled()" class="alert alert-warning alert-dismissable">
        <span class="pficon pficon-warning-triangle-o"></span>
        <strong>{{ $t('hotspot.warning_privacy') }}!</strong> {{ $t('hotspot.warning_privacy_text') }} <a href="#/profile" class="alert-link">{{ $t('hotspot.profile_link') }}</a>
      </div>
    </div>
    <div
      v-if="rows.length > 0 && !isLoading"
      class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12"
    >
      <div class="col-sm-3">
        <button class="btn btn-primary" @click="getAll()">{{$t('session.refresh')}}</button>
      </div>
    </div>
    <div
      v-if="!isLoading && rows.length > 0"
      class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12"
    >
      <div class="result-list">{{total}} {{total == 1 ? $t('result') : $t('results')}}</div>
    </div>
    <div v-if="!isLoading && rows.length > 0">
      <form v-on:submit.prevent="searchFn($event)">
        <input
          class="form-control input-lg search-table-input"
          type="text"
          :placeholder="tableLangsTexts.globalSearchPlaceholder"
        >
      </form>
    </div>
    <div v-if="!isLoading && isLoadingTable" class="spinner spinner-lg spinner-adjust"></div>
    <vue-good-table
      v-if="!isLoading && !isLoadingTable && rows.length > 0"
      @perPageChanged="handlePerPage"
      :customRowsPerPageDropdown="[25,50,100]"
      :perPage="hotspotPerPage"
      :columns="columns"
      :rows="rows"
      :lineNumbers="false"
      :defaultSortBy="{field: 'name', type: 'asc'}"
      :globalSearch="true"
      :paginate="false"
      styleClass="table"
      :nextText="tableLangsTexts.nextText"
      :prevText="tableLangsTexts.prevText"
      :rowsPerPageText="tableLangsTexts.rowsPerPageText"
      :globalSearchPlaceholder="tableLangsTexts.globalSearchPlaceholder"
      :ofText="tableLangsTexts.ofText"
    >
      <template slot="table-row" slot-scope="props">
        <td>
          <a :href="'#/hotspots/'+ props.row.id">{{ props.row.name }}</a>
        </td>
        <td class="fancy">{{ props.row.description }}</td>
        <td v-if="props.row.business_name.length > 0" class="fancy">{{ props.row.business_name }}</td>
        <td v-if="props.row.business_name.length == 0" class="fancy">
          <span class="pficon pficon-error-circle-o"></span>
          <span class="red">{{$t('hotspot.missing_business_name')}}</span>
        </td>
        <td v-if="props.row.business_vat.length > 0" class="fancy">{{ props.row.business_vat }}</td>
        <td v-if="props.row.business_vat.length == 0" class="fancy">
          <span class="pficon pficon-error-circle-o"></span>
          <span class="red">{{$t('hotspot.missing_business_vat')}}</span>
        </td>
        <td
          v-if="props.row.business_address.length > 0"
          class="fancy"
        >{{ props.row.business_address }}</td>
        <td v-if="props.row.business_address.length == 0" class="fancy">
          <span class="pficon pficon-error-circle-o"></span>
          <span class="red">{{$t('hotspot.missing_business_address')}}</span>
        </td>
        <td v-if="props.row.business_email.length > 0" class="fancy">{{ props.row.business_email }}</td>
        <td v-if="props.row.business_email.length == 0" class="fancy">
          <span class="pficon pficon-error-circle-o"></span>
          <span class="red">{{$t('hotspot.missing_business_email')}}</span>
        </td>
        <td class="fancy">{{ props.row.business_dpo || '-' }}</td>
        <td class="fancy">{{ props.row.business_dpo_mail || '-' }}</td>
        <td class="fancy">{{ props.row.created | formatDate }}</td>
        <td>
          <hotspot-action details="true" :obj="props.row" :update="getAll"></hotspot-action>
        </td>
      </template>
    </vue-good-table>
    <div v-if="!isLoadingTable && !isLoading && rows.length > 0" class="right paginator">
      <span class="page-count">
        <b>{{hotspotPage}}</b>
        {{tableLangsTexts.ofText}} {{total / hotspotPerPage | adjustPage}} (
        <b>{{hotspotPerPage}}</b>
        {{tableLangsTexts.rowsPerPageText}})
      </span>
      <button
        :disabled="availablePrevPage()"
        @click="prevPage()"
        class="btn btn-default"
      >{{tableLangsTexts.prevText}}</button>
      <button
        :disabled="availableNextPage()"
        @click="nextPage()"
        class="btn btn-default"
      >{{tableLangsTexts.nextText}}</button>
    </div>

    <div
      class="modal fade"
      id="HScreateModal"
      tabindex="-1"
      role="dialog"
      aria-labelledby="HScreateModalLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="HScreateModalLabel">{{ $t("create") }}</h4>
          </div>
          <form class="form-horizontal" role="form" v-on:submit.prevent="createHotspot()">
            <div class="modal-body">
              <div class="form-group">
                <label
                  class="col-sm-4 control-label"
                  for="textInput-modal-markup"
                >{{ $t("hotspot.name") }}</label>
                <div class="col-sm-8">
                  <input
                    required
                    v-model="newObj.name"
                    pattern="^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$"
                    type="text"
                    id="textInput-modal-markup"
                    class="form-control"
                    :placeholder="$t('hotspot.hostname_hotspot')"
                  >
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-4 control-label"
                  for="textInput2-modal-markup"
                >{{ $t("hotspot.description") }}</label>
                <div class="col-sm-8">
                  <input
                    required
                    v-model="newObj.description"
                    type="text"
                    id="textInput2-modal-markup"
                    class="form-control"
                    :placeholder="$t('hotspot.description')"
                  >
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-4 control-label"
                  for="textInput2-modal-markup"
                >{{ $t("hotspot.business_name") }}</label>
                <div class="col-sm-8">
                  <input
                    required
                    v-model="newObj.business_name"
                    type="text"
                    id="textInput2-modal-markup"
                    class="form-control"
                    :placeholder="$t('hotspot.business_name')"
                  >
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-4 control-label"
                  for="textInput2-modal-markup"
                >{{ $t("hotspot.business_vat") }}</label>
                <div class="col-sm-8">
                  <input
                    required
                    v-model="newObj.business_vat"
                    type="text"
                    id="textInput2-modal-markup"
                    class="form-control"
                    :placeholder="$t('hotspot.business_vat')"
                  >
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-4 control-label"
                  for="textInput2-modal-markup"
                >{{ $t("hotspot.business_address") }}</label>
                <div class="col-sm-8">
                  <input
                    required
                    v-model="newObj.business_address"
                    type="text"
                    id="textInput2-modal-markup"
                    class="form-control"
                    :placeholder="$t('hotspot.business_address')"
                  >
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-4 control-label"
                  for="textInput2-modal-markup"
                >{{ $t("hotspot.business_email") }}</label>
                <div class="col-sm-8">
                  <input
                    required
                    v-model="newObj.business_email"
                    type="email"
                    id="textInput2-modal-markup"
                    class="form-control"
                    :placeholder="$t('hotspot.business_email')"
                  >
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-4 control-label"
                  for="textInput2-modal-markup"
                >{{ $t("hotspot.business_dpo") }}</label>
                <div class="col-sm-8">
                  <input
                    v-model="newObj.business_dpo"
                    type="text"
                    id="textInput2-modal-markup"
                    class="form-control"
                    :placeholder="$t('hotspot.business_dpo')"
                  >
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-4 control-label"
                  for="textInput2-modal-markup"
                >{{ $t("hotspot.business_dpo_mail") }}</label>
                <div class="col-sm-8">
                  <input
                    v-model="newObj.business_dpo_mail"
                    type="text"
                    id="textInput2-modal-markup"
                    class="form-control"
                    :placeholder="$t('hotspot.business_dpo_mail')"
                  >
                </div>
              </div>
              <!-- privacy dislaimer -->
              <div v-if="disclaimers.privacyDisclaimers && disclaimers.privacyDisclaimers.length" class="form-group">
                <label
                  class="col-sm-4 control-label"
                  for="textInput2-modal-markup"
                >{{ $t("hotspot.privacy_disclaimer") }}</label>
                <div class="col-sm-8">
                  <select v-model="disclaimers.selectedPrivacyDisclaimerId" class="form-control">
                    <option value="">{{ $t("hotspot.disclaimer_default") }}</option>
                    <option v-for="privacyDisclaimer in disclaimers.privacyDisclaimers" :key="privacyDisclaimer.id" :value="privacyDisclaimer.id">
                      {{ privacyDisclaimer.title }}
                    </option>
                  </select>
                </div>
              </div>
              <!-- tos dislaimer -->
              <div v-if="disclaimers.tosDisclaimers && disclaimers.tosDisclaimers.length" class="form-group">
                <label
                  class="col-sm-4 control-label"
                  for="textInput2-modal-markup"
                >{{ $t("hotspot.tos_disclaimer") }}</label>
                <div class="col-sm-8">
                  <select v-model="disclaimers.selectedTosDisclaimerId" class="form-control">
                    <option value="">{{ $t("hotspot.disclaimer_default") }}</option>
                    <option v-for="tosDisclaimer in disclaimers.tosDisclaimers" :key="tosDisclaimer.id" :value="tosDisclaimer.id">
                      {{ tosDisclaimer.title }}
                    </option>
                  </select>
                </div>
              </div>
              <div v-if="errors.create" class="alert alert-danger alert-dismissable">
                <span class="pficon pficon-error-circle-o"></span>
                <strong>{{ $t("hotspot.create_error_title") }}</strong>
                . {{ $t("hotspot.create_error_sub") }}.
              </div>
            </div>
            <div class="modal-footer">
              <span v-if="newObj.onAction" class="spinner spinner-sm spinner-inline modal-spinner"></span>
              <button type="button" class="btn btn-default" data-dismiss="modal">{{ $t("cancel") }}</button>
              <button type="submit" class="btn btn-primary">{{ $t("create") }}</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import HotspotService from "../services/hotspot";
import StorageService from "../services/storage";
import UtilService from "../services/util";
import DisclaimerService from "../services/disclaimer";
import PreferenceService from "../services/preference";

import HotspotAction from "../directives/HotspotAction.vue";

export default {
  name: "Hotspots",
  mixins: [HotspotService, StorageService, UtilService, DisclaimerService, PreferenceService],
  components: {
    hotspotAction: HotspotAction
  },
  data() {
    return {
      msg: this.$i18n.t("menu.hotspots"),
      isLoading: true,
      isLoadingTable: true,
      columns: [
        {
          label: this.$i18n.t("hotspot.name"),
          field: "name",
          filterable: true
        },
        {
          label: this.$i18n.t("hotspot.description"),
          field: "description",
          filterable: true
        },
        {
          label: this.$i18n.t("hotspot.business_name"),
          field: "business_name",
          filterable: true
        },
        {
          label: this.$i18n.t("hotspot.business_vat"),
          field: "business_vat",
          filterable: true
        },
        {
          label: this.$i18n.t("hotspot.business_address"),
          field: "business_address",
          filterable: true
        },
        {
          label: this.$i18n.t("hotspot.business_email"),
          field: "business_email",
          filterable: true
        },
        {
          label: this.$i18n.t("hotspot.business_dpo_recap"),
          field: "business_dpo",
          filterable: true
        },
        {
          label: this.$i18n.t("hotspot.business_dpo_mail_recap"),
          field: "business_dpo_mail",
          filterable: true
        },
        {
          label: this.$i18n.t("hotspot.created"),
          field: "created",
          sortable: false
        },
        {
          label: this.$i18n.t("action"),
          field: "",
          sortable: false
        }
      ],
      rows: [],
      tableLangsTexts: this.tableLangs(),
      newObj: {
        name: "",
        description: ""
      },
      errors: {
        create: false
      },
      hotspotPerPage: 25,
      hotspotPage: 1,
      total: 0,
      isAdmin: this.get("loggedUser").account_type == "admin",
      searchString: "",
      user: {
        login: this.get("loggedUser") || null,
        info: this.$parent.user.info
      },
      disclaimers: {
        data: {},
        preferredPrivacyDisclaimerId: "",
        preferredTosDisclaimerId: "",
        selectedPrivacyDisclaimerId: "",
        selectedTosDisclaimerId: "",
        privacyDisclaimers: [],
        tosDisclaimers: [],
      }
    };
  },
  mounted() {
    // get hotspot list
    this.getAll();
    this.getDisclaimers();
  },
  methods: {
    privacyFilled() {
      return (
        this.user.info.privacy_name.length > 0 &&
        this.user.info.privacy_vat.length > 0 &&
        this.user.info.privacy_address.length > 0 &&
        this.user.info.privacy_email.length > 0
      )
    },
    handlePerPage(evt) {
      this.set("hotspots_per_page", evt.currentPerPage);
    },
    searchFn(evt) {
      var elem = evt.srcElement || evt.target;
      this.searchString = elem[0].value;
      this.getAll(true);
    },
    getAll(reset) {
      this.isLoadingTable = true;

      if (reset) {
        this.hotspotPage = 1;
        this.total = 0;
      }

      this.hotspotGetAll(
        this.hotspotPage,
        this.hotspotPerPage,
        encodeURIComponent(this.searchString),
        success => {
          this.rows = success.body.data;
          this.total = success.body.total;
          this.isLoading = false;
          this.isLoadingTable = false;
        },
        error => {
          this.isLoading = false;
          this.isLoadingTable = false;
          this.rows = [];
          console.error(error);
        }
      );
    },
    createHotspot() {
      this.newObj.onAction = true;
      this.newObj.privacy_disclaimer_id = this.disclaimers.selectedPrivacyDisclaimerId ? parseInt(this.disclaimers.selectedPrivacyDisclaimerId) : 0;
      this.newObj.tos_disclaimer_id = this.disclaimers.selectedTosDisclaimerId ? parseInt(this.disclaimers.selectedTosDisclaimerId) : 0;

      this.hotspotCreate(
        this.newObj,
        success => {
          this.newObj.onAction = false;
          this.newObj.name = "";
          this.newObj.description = "";
          $("#HScreateModal").modal("toggle");
          this.getAll();
        },
        error => {
          this.newObj.onAction = false;
          this.errors.create = true;
          console.error(error.body.message);
        }
      );
    },
    prevPage() {
      this.isLoading = true;
      if (this.hotspotPage != 1) {
        this.hotspotPage--;
      }
      this.getAll();
    },
    nextPage() {
      this.isLoading = true;
      if (this.hotspotPage != Math.ceil(this.total / this.hotspotPerPage)) {
        this.hotspotPage++;
      }
      this.getAll();
    },
    availablePrevPage() {
      return this.hotspotPage == 1;
    },
    availableNextPage() {
      return this.hotspotPage == Math.ceil(this.total / this.hotspotPerPage);
    },
    getDisclaimers() {
      this.disclaimersByAccount(
        this.user.login.id,
        success => {
          this.disclaimers.data = success.body;
          this.privacy_disclaimers = [];
          this.tos_disclaimers = [];

          this.disclaimers.data.forEach((disclaimer) => {
            if (disclaimer.type === "privacy") {
              this.disclaimers.privacyDisclaimers.push(disclaimer);
            } else {
              this.disclaimers.tosDisclaimers.push(disclaimer);
            }
          });
          this.getPreferredDisclaimers();
        },
        error => {
          console.error(error.body);
        }
      );
    },
    getPreferredDisclaimers() {
      this.accountPrefGet(
        this.user.login.id,
        success => {
          this.disclaimers.preferredPrivacyDisclaimerId = "";
          this.disclaimers.preferredTosDisclaimerId = "";
          this.disclaimers.selectedPrivacyDisclaimerId = "";
          this.disclaimers.selectedTosDisclaimerId = "";

          success.body.forEach((pref) => {
            if (pref.key === "custom_disclaimers_privacy") {
              this.disclaimers.preferredPrivacyDisclaimerId = pref.value;
              this.disclaimers.selectedPrivacyDisclaimerId = pref.value;
            } else if (pref.key === "custom_disclaimers_terms") {
              this.disclaimers.preferredTosDisclaimerId = pref.value;
              this.disclaimers.selectedTosDisclaimerId = pref.value;
            }
          });
        },
        error => {
          console.error(error.body);
        }
      );
    },
  }
};
</script>

<style>
.create-hotspot {
  float: right;
  margin-top: -52px;
  margin-right: 35px;
}
</style>