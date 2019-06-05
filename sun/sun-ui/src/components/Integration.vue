<template>
  <div>
    <h2>{{ msg }}</h2>
    <div v-if="isLoading" class="spinner spinner-lg"></div>
    <div v-if="!isLoading && integrations.length == 0" class="blank-slate-pf" id>
      <div class="blank-slate-pf-icon">
        <span class="fa fa-cubes"></span>
      </div>
      <h1>{{$t('integrations.not_found')}}</h1>
      <p>{{$t('integrations.not_found_desc')}}.</p>
    </div>
    <div
      v-if="(user.account_type == 'admin' || user.account_type == 'reseller') && !isLoading && integrations.length > 0"
      class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12"
    >
      <label v-if="!isLoading" class="col-sm-2 control-label" for="textInput-markup">Hotspot</label>
      <div v-if="!isLoading" class="col-sm-4">
        <select v-on:change="getMaps()" v-model="hotspotSearchId" class="form-control">
          <option v-bind:value="0">-</option>
          <option
            v-for="hotspot in hotspots"
            v-bind:key="hotspot.id"
            v-bind:value="hotspot.id"
          >{{ hotspot.name }}</option>
        </select>
      </div>
    </div>
    <div
      v-if="!isLoading && integrations.length > 0"
      class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12"
    >
      <div class="col-sm-3">
        <button class="btn btn-primary" @click="getMaps()">{{$t('integrations.refresh')}}</button>
      </div>
    </div>

    <div
      v-if="!isLoading && hotspotSearchId > 0 && integrations.length > 0"
      class="form-group col-xs-12 col-sm-12 col-md-12 col-lg-12"
    >
      <div class="alert alert-info">
        <span class="pficon pficon-info"></span>
        <strong>{{$t('integrations.info')}}:</strong>
        {{$t('integrations.info_desc')}}.
      </div>
    </div>

    <div v-if="!isLoading && isError" class="form-group col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <div class="alert alert-danger alert-dismissable">
        <button type="button" class="close" data-dismiss="alert" aria-label="Close">
          <span class="pficon pficon-close"></span>
        </button>
        <span class="pficon pficon-error-circle-o"></span>
        <strong>{{$t('integrations.error')}}.</strong>
        {{$t('integrations.error_desc')}}.
      </div>
    </div>

    <div
      v-if="!isLoading && hotspotSearchId > 0"
      class="list-group list-view-pf list-view-pf-view col-xs-12 col-sm-12 col-md-12 col-lg-12 adjust-list-view adjust-margin-top"
    >
      <div
        v-for="(i,ik) in integrations"
        :key="ik"
        :class="['list-group-item', maps[i.id] ? '' : 'disabled']"
      >
        <div class="list-group-item-header">
          <div class="list-view-pf-actions">
            <button
              @click="maps[i.id] ? deleteMap(i.id) : createMap(i.id)"
              :class="['btn', maps[i.id] ? 'btn-default' : 'btn-primary']"
            >{{ maps[i.id] ? $t('integrations.disable') : $t('integrations.enable')}}</button>
          </div>
          <div class="list-view-pf-main-info adjust-main-info">
            <div class="list-view-pf-left">
              <img class="logo-int" :src="i.logo">
            </div>
            <div class="list-view-pf-body">
              <div class="list-view-pf-description">
                <div class="list-group-item-heading">
                  <b>{{i.name}}</b>
                </div>
                <div class="list-group-item-text">{{i.description}}</div>
              </div>
              <div class="list-view-pf-additional-info">
                <div class="list-view-pf-additional-info-item">
                  <span class="fa fa-globe"></span>
                  <a target="_blank" :href="i.site">{{i.site}}</a>
                </div>
                <div class="list-view-pf-additional-info-item">
                  <span
                    :class="[maps[i.id] && maps[i.id].last_sync ? 'fa fa-refresh' : 'pficon pficon-warning-triangle-o']"
                  ></span>
                  <span
                    v-if="maps[i.id] && maps[i.id].last_sync"
                  >{{maps[i.id] && maps[i.id].last_sync | formatDate(false,true)}}</span>
                  <span v-else>{{$t('integrations.failed')}}</span>
                  <button
                    :disabled="!maps[i.id]"
                    @click="createMap(i.id)"
                    class="btn btn-sm btn-primary adjust-space"
                  >{{$t('integrations.force_sync')}}</button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>


<script>
import IntegrationService from "../services/integration";
import StorageService from "../services/storage";
import HotspotService from "../services/hotspot";
import UtilService from "../services/util";

export default {
  name: "Accounts",
  mixins: [IntegrationService, StorageService, UtilService, HotspotService],
  data() {
    return {
      msg: this.$i18n.t("menu.integrations"),
      isLoading: true,
      isError: false,
      hotspots: [],
      integrations: [],
      maps: {},
      hotspotSearchId: 0,
      user: this.get("loggedUser") || null
    };
  },
  mounted() {
    if (this.$route.params.hotspotId !== undefined) {
      this.hotspotSearchId = this.$route.params.hotspotId;
    }

    // get account list
    var context = this;
    this.getAllHotspots(function() {
      context.getIntegrations();
      context.getMaps();
    });
  },
  methods: {
    getAllHotspots(callback) {
      this.hotspotGetAll(
        null,
        null,
        null,
        success => {
          this.hotspots = success.body.data;
          var hsId = this.get("integrations_hotspot_id") || this.hotspots[0].id;
          this.hotspotSearchId = hsId;
          if (
            this.$parent.user.info.type == "customer" ||
            this.$parent.user.info.type == "desk"
          ) {
            hsId = this.$parent.user.info.hotspot_id;
          }

          this.isLoading = false;

          callback();
        },
        error => {
          console.error(error);

          callback();
        }
      );
    },
    getIntegrations() {
      this.isLoading = true;
      this.integrationGetAll(
        success => {
          this.integrations = success.body;
          this.isLoading = false;
        },
        error => {
          console.error(error);
          this.isLoading = false;
          this.integrations = [];
        }
      );
    },
    getMaps() {
      this.isLoading = true;
      this.set("integrations_hotspot_id", this.hotspotSearchId);

      this.mapGetAll(
        this.hotspotSearchId,
        success => {
          var result = success.body;
          var maps = {};

          for (var i in result) {
            maps[result[i].integration_id] = {
              last_sync:
                +new Date(result[i].last_sync) > 0 ? result[i].last_sync : false
            };
          }

          this.maps = maps;
          this.isLoading = false;
        },
        error => {
          console.error(error);
          this.isLoading = false;
          this.maps = [];
        }
      );
    },
    createMap(integrationId) {
      this.isLoading = true;
      this.mapCreate(
        this.hotspotSearchId,
        integrationId,
        success => {
          this.isError = false;
          this.getMaps();
        },
        error => {
          console.error(error);
          this.isLoading = false;
          this.isError = true;
        }
      );
    },
    deleteMap(integrationId) {
      this.isLoading = true;
      this.mapDelete(
        this.hotspotSearchId,
        integrationId,
        success => {
          this.isError = false;
          this.getMaps();
        },
        error => {
          console.error(error);
          this.isLoading = false;
          this.isError = true;
        }
      );
    }
  }
};
</script>

<style>
.logo-int {
  height: 50px !important;
}

.adjust-list-view {
  padding-right: 0px !important;
}

.adjust-main-info {
  padding-bottom: 8px !important;
  padding-top: 8px !important;
}

.adjust-margin-top {
  margin-top: 0px !important;
}

.adjust-space {
  margin-left: 8px !important;
}
</style>