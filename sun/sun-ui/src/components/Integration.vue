<template>
  <div>
    <h2>{{ msg }}</h2>
    <div v-if="isLoading" class="spinner spinner-lg"></div>
    <div
      v-if="(user.account_type == 'admin') || (user.account_type == 'reseller') && !isLoading"
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
    <div v-if="!isLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <div class="col-sm-3">
        <button class="btn btn-primary" @click="getMaps()">{{$t('integrations.refresh')}}</button>
      </div>
    </div>

    <div v-if="!isLoading" class="form-group col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <div class="alert alert-info alert-dismissable">
        <span class="pficon pficon-info"></span>
        <strong>{{$t('integrations.info')}}:</strong>
        {{$t('integrations.info_desc')}}.
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
                  <span class="fa fa-refresh"></span>
                  {{maps[i.id] && maps[i.id].last_sync | formatDate(true,true)}}
                  <button
                    :disabled="!maps[i.id]"
                    @click="createMap(i.id)"
                    class="btn btn-sm btn-default"
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
      accountType: this.get("loggedUser").account_type,
      hotspots: [],
      integrations: [],
      maps: {},
      isAdmin: this.get("loggedUser").account_type == "admin",
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
              last_sync: result[i].last_sync
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
          this.getMaps();
        },
        error => {
          console.error(error);
          this.isLoading = false;
        }
      );
    },
    deleteMap(integrationId) {
      this.isLoading = true;
      this.mapDelete(
        this.hotspotSearchId,
        integrationId,
        success => {
          this.getMaps();
        },
        error => {
          console.error(error);
          this.isLoading = false;
        }
      );
    }
  }
};
</script>

<style>
.logo-int {
  width: 50px;
}

.adjust-list-view {
  padding-right: 0px !important;
}

.adjust-main-info {
  padding-bottom: 0px !important;
  padding-top: 10px !important;
}

.adjust-margin-top {
  margin-top: 0px;
}
</style>