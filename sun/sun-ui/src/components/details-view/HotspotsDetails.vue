<template>
  <div>
    <h2>Hotspot
      <strong class="soft">{{ info.data.name }} - {{info.data.description}}</strong>
    </h2>

    <div class="row row-cards-pf">
      <div class="col-xs-12 col-sm-6 col-md-6 col-lg-3">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              <router-link class="card-link" :to="accountLink">
                <span class="pficon pficon-users card-info-title"></span>
                {{ $t("dashboard.accounts") }}
                <span v-if="!totals.accounts.isLoading" class="right">
                  <strong class="soft">{{ totals.accounts.count}}</strong>
                </span>
              </router-link>
              <div v-if="totals.accounts.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
        </div>
      </div>
      <div class="col-xs-12 col-sm-6 col-md-6 col-lg-3">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              <router-link class="card-link" :to="unitLink">
                <span class="pficon pficon-connected card-info-title"></span>
                {{ $t("dashboard.units") }}
                <span v-if="!totals.units.isLoading" class="right">
                  <strong class="soft">{{ totals.units.count}}</strong>
                </span>
              </router-link>
              <div v-if="totals.units.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
        </div>
      </div>
      <div class="col-xs-12 col-sm-4 col-md-4 col-lg-2">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              <router-link class="card-link" :to="userLink">
                <span class="fa fa-users card-info-title"></span>
                {{ $t("dashboard.users") }}
                <span v-if="!totals.users.isLoading" class="right">
                  <strong class="soft">{{ totals.users.count}}</strong>
                </span>
              </router-link>
              <div v-if="totals.users.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
        </div>
      </div>
      <div class="col-xs-12 col-sm-4 col-md-4 col-lg-2">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              <router-link class="card-link" :to="deviceLink">
                <span class="fa fa-laptop card-info-title"></span>
                {{ $t("dashboard.devices") }}
                <span v-if="!totals.devices.isLoading" class="right">
                  <strong class="soft">{{ totals.devices.count}}</strong>
                </span>
              </router-link>
              <div v-if="totals.devices.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
        </div>
      </div>
      <div class="col-xs-12 col-sm-4 col-md-4 col-lg-2">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              <router-link class="card-link" :to="sessionLink">
                <span class="fa fa-list card-info-title"></span>
                {{ $t("dashboard.sessions") }}
                <span v-if="!totals.sessions.isLoading" class="right">
                  <strong class="soft">{{ totals.sessions.count}}</strong>
                </span>
              </router-link>
              <div v-if="totals.sessions.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
        </div>
      </div>
      <div class="col-xs-12 col-sm-12 col-md-6">
      </div>
    </div>
    <div class="row row-cards-pf">
      <div v-if="preferences.vouchersAvailable" :class="['col-xs-12 col-sm-12', user.account_type == 'admin' || user.account_type == 'reseller' ? 'col-md-12' : 'col-md-12']">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{ $t("hotspot.vouchers") }}
              <div v-if="vouchers.isLoading" class="spinner spinner-sm right"></div>
              <span @click="getVouchers()" class="fa fa-refresh right link"></span>
            </h2>
          </div>
          <div v-if="!vouchers.isLoading" class="card-pf-body">
            <vue-good-table :perPage="5" :paginate="true" :columns="columns" :rows="vouchers.data" :lineNumbers="false" :defaultSortBy="{field: 'expires', type: 'asc'}"
              styleClass="table" :nextText="tableLangsTexts.nextText" :prevText="tableLangsTexts.prevText" :rowsPerPageText="tableLangsTexts.rowsPerPageText"
              :globalSearchPlaceholder="tableLangsTexts.globalSearchPlaceholder" :ofText="tableLangsTexts.ofText">
              <template slot="table-row" slot-scope="props">
                <td class="fancy">
                  <strong>{{ props.row.code }}</strong>
                </td>
                <td class="fancy">
                  <span :class="['pficon', props.row.auto_login ? 'pficon-ok' : 'pficon-error-circle-o']"></span>
                </td>
                <td class="fancy">
                  <div>
                    <strong>{{ $t('user.kbps_down') }}</strong>: {{ props.row.bandwidth_down || '-' }}</div>
                  <div>
                    <strong>{{ $t('user.kbps_up') }}</strong>: {{ props.row.bandwidth_up || '-' }}</div>
                </td>
                <td class="fancy">{{ props.row.duration }} ({{$t('hotspot.days')}})</td>
                <td class="fancy">
                  <span :class="['fa', checkVoucherUse(props.row.expires) ? 'fa-check green' : 'fa-minus']"></span>
                </td>
                <td>
                  <button v-on:click="printVoucher(props.row)" class="btn btn-primary" type="button">
                    <span class="fa fa-print"></span>
                  </button>
                  <button v-on:click="deleteVoucher(props.row.id)" class="btn btn-danger" type="button">
                    <span class="fa fa-remove"></span>
                  </button>
                </td>
              </template>
            </vue-good-table>
          </div>
          <div v-if="!vouchers.isLoading" class="card-pf-footer">
            <div class="dropdown card-pf-time-frame-filter">
              <button data-toggle="modal" data-target="#voucherModal" class="btn btn-primary" type="button">{{ $t("hotspot.create_voucher") }}</button>
              <button v-on:click="printAllVoucher()" class="btn btn-default" type="button">
                <span class="fa fa-print"></span>
                {{ $t("hotspot.print_all_voucher") }}
              </button>
            </div>
            <p>
              <a href="#" class="card-pf-link-with-icon">
              </a>
            </p>
          </div>
        </div>
      </div>
    </div>

    <div v-if="user.account_type == 'admin' || user.account_type == 'reseller'" class="row row-cards-pf">
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{ $t("hotspot.preferences") }}
              <div v-if="preferences.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
          <form class="form-horizontal" role="form" v-on:submit.prevent="updatePreferences(preferences.global)">
            <div v-if="!preferences.isLoading" class="card-pf-body">
              <div v-for="pref in preferences.global" :key="pref.key" class="form-group">
                <label class="col-sm-4 control-label" for="textInput-markup">{{$t('hotspot.'+pref.key)}}
                  <span :class="[getPrefIcon(pref.key)]"></span>
                </label>
                <div class="col-sm-6">
                  <input v-model="pref.value" :type="getInputType(pref.key, pref.value)" id="textInput-markup" class="form-control">
                </div>
              </div>
            </div>
            <div v-if="!preferences.isLoading" class="card-pf-footer">
              <div class="dropdown card-pf-time-frame-filter">
                <button type="submit" class="btn btn-primary">{{ $t("update") }}</button>
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

    <div v-if="displayCaptivePortalOptions" class="row row-cards-pf">
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{ $t("hotspot.preferences_captive") }}
              <div v-if="preferences.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
          <form class="form-horizontal" role="form" v-on:submit.prevent="updatePreferencesCaptive(preferences.captive)">
            <div class="card-pf-body">
              <div v-for="pref in preferences.captive" :key="pref.key" :class="[pref.onError ? 'has-error' : '', 'form-group']">
                <label class="col-sm-4 control-label" for="textInput-markup">{{$t('hotspot.'+pref.key)}}
                  <span :class="[getPrefIcon(pref.key)]"></span>
                </label>
                <div class="col-sm-6">
                  <picture-input v-if="pref.key == 'captive_3_logo' || pref.key == 'captive_5_banner'" :ref="'prefInput-'+pref.key" :prefill="urltoFile(pref.value, pref.key)"
                    :alertOnError="false" @change="onChanged(pref)" :width="100" :height="100"
                    :crop="false" :zIndex="1000" :customStrings="uploadLangstexts" removeButtonClass="btn btn-danger" buttonClass="btn btn-default">

                  </picture-input>
                  <span v-if="pref.onError" class="help-block">{{$t('upload_file_exceed')}}</span>

                  <vue-editor :editorToolbar="customToolbar" v-if="pref.key == 'captive_6_description'" v-model="pref.value"></vue-editor>

                  <sketch-picker @input="onUpdate(pref.value)" class="absolute-center" v-if="pref.key == 'captive_7_background'" v-model="pref.value"
                  />

                  <input required v-if="pref.key != 'captive_6_description' && pref.key != 'captive_3_logo' && pref.key != 'captive_5_banner' && pref.key != 'captive_7_background'"
                    v-model="pref.value" :type="getInputType(pref.key, pref.value)" class="form-control">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-4 control-label" for="textInput-markup">
                  {{$t('hotspot.captive_preview')}}
                </label>
                <div class="col-sm-6">
                  <div v-if="preferences.isLoading" class="captive-preview">
                    <div class="spinner spinner-lg absolute-center"></div>
                  </div>
                  <captive-portal v-if="!preferences.isLoading" :obj="preferences.captive"></captive-portal>
                </div>
              </div>
            </div>
            <div v-if="!preferences.isLoading" class="card-pf-footer">
              <div class="dropdown card-pf-time-frame-filter">
                <button type="submit" class="btn btn-primary">{{ $t("update") }}</button>
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

    <div v-if="(user.account_type == 'admin' || user.account_type == 'reseller') && totals.units.count > 0" class="row row-cards-pf">
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{ $t("hotspot.mac_auth") }}
              <div v-if="macAuth.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
          <div v-if="!macAuth.isLoading" class="card-pf-body">
            <vue-good-table :perPage="5" :paginate="true" :columns="columnsMAC" :rows="macAuth.data" :lineNumbers="false" :defaultSortBy="{field: 'name', type: 'asc'}"
              styleClass="table" :nextText="tableLangsTexts.nextText" :prevText="tableLangsTexts.prevText" :rowsPerPageText="tableLangsTexts.rowsPerPageText"
              :globalSearchPlaceholder="tableLangsTexts.globalSearchPlaceholder" :ofText="tableLangsTexts.ofText">
              <template slot="table-row" slot-scope="props">
                <td class="fancy">
                  <strong>{{ props.row.name }}</strong>
                </td>
                <td class="fancy">{{ props.row.username }}</td>
                <td class="fancy">
                  <div>
                    <strong>{{ $t('user.kbps_down') }}</strong>: {{ props.row.kbps_down || '-' }}</div>
                  <div>
                    <strong>{{ $t('user.kbps_up') }}</strong>: {{ props.row.kbps_up || '-' }}</div>
                </td>
                <td class="fancy">{{ props.row.created | formatDate}}</td>
                <td>
                  <button v-on:click="deleteMacAuth(props.row.id)" class="btn btn-danger" type="button">
                    <span class="fa fa-remove"></span>
                  </button>
                </td>
              </template>
            </vue-good-table>
          </div>
          <div v-if="!macAuth.isLoading" class="card-pf-footer">
            <div class="dropdown card-pf-time-frame-filter">
              <button data-toggle="modal" data-target="#macAuthModal" class="btn btn-primary" type="button">{{ $t("hotspot.new_mac_auth") }}</button>
            </div>
            <p>
              <a href="#" class="card-pf-link-with-icon">
              </a>
            </p>
          </div>
        </div>
      </div>
    </div>

    <div class="row row-cards-pf">
      <div v-if="user.account_type == 'admin' || user.account_type == 'reseller'" class="col-xs-12 col-sm-12 col-md-12 col-lg-12">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{ $t("hotspot.info") }}
              <div v-if="info.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
          <div v-if="!info.isLoading" class="card-pf-body">
            <div class="list-details">
              <dt>{{ $t("hotspot.owner") }}</dt>
              <dd>
                <a :href="'#/accounts/' + info.data.account_id">{{info.data.account_name}}</a>
              </dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("hotspot.name") }}</dt>
              <dd>{{info.data.name}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("hotspot.description") }}</dt>
              <dd>{{info.data.description}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("hotspot.created") }}</dt>
              <dd>{{info.data.created | formatDate}}</dd>
            </div>
          </div>
          <div v-if="!info.isLoading" class="card-pf-footer">
            <div class="dropdown card-pf-time-frame-filter">
              <hotspot-action details="false" :obj="info.data" :update="getInfo"></hotspot-action>
            </div>
            <p>
              <a href="#" class="card-pf-link-with-icon">
              </a>
            </p>
          </div>
        </div>
      </div>
    </div>

    <div class="modal fade" id="voucherModal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="myModalLabel">{{$t('hotspot.voucher_creation')}}</h4>
          </div>
          <div class="modal-body">
            <form class="form-horizontal">
              <div class="form-group">
                <label class="col-sm-5 control-label" for="textInput-modal-markup">{{$t('hotspot.voucher_creation_count')}}</label>
                <div class="col-sm-7">
                  <input v-model="vouchersCount" type="text" id="textInput-modal-markup" class="form-control">
                </div>
              </div>

              <div class="form-group">
                <label class="col-sm-5 control-label" for="textInput-modal-markup">{{$t('hotspot.auto_login')}}</label>
                <div class="col-sm-7">
                  <input v-model="newVoucher.auto_login" type="checkbox" id="textInput-modal-markup" class="form-control">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-5 control-label" for="textInput-modal-markup">{{$t('hotspot.bandwidth_down')}}</label>
                <div class="col-sm-7">
                  <input v-model="newVoucher.bandwidth_down" type="number" id="textInput-modal-markup" class="form-control">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-5 control-label" for="textInput-modal-markup">{{$t('hotspot.bandwidth_up')}}</label>
                <div class="col-sm-7">
                  <input v-model="newVoucher.bandwidth_up" type="number" id="textInput-modal-markup" class="form-control">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-5 control-label" for="textInput-modal-markup">{{$t('hotspot.duration')}} ({{$t('hotspot.days')}})</label>
                <div class="col-sm-7">
                  <input v-model="newVoucher.duration" type="number" id="textInput-modal-markup" class="form-control">
                </div>
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-default" data-dismiss="modal">{{$t('cancel')}}</button>
            <button v-on:click="createVoucher()" type="button" class="btn btn-primary">{{$t('save')}}</button>
          </div>
        </div>
      </div>
    </div>

    <div class="modal fade" id="macAuthModal" tabindex="-1" role="dialog" aria-labelledby="HScreateModalLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title">{{ $t("hotspot.new_mac_auth") }}</h4>
          </div>
          <form class="form-horizontal" role="form" v-on:submit.prevent="createMacAuth()">
            <div class="modal-body">
              <div class="form-group">
                <label class="col-sm-5 control-label" for="ACtextInput2-modal-markup">{{ $t("unit.unit") }}</label>
                <div class="col-sm-7">
                  <select v-model="newMACAuth.unit" class="form-control">
                    <option v-for="u in totals.units.data" v-bind:key="u.id" v-bind:value="u">
                      {{ u.name }}
                    </option>
                  </select>
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-5 control-label" for="textInput-modal-markup">{{ $t("unit.description") }}</label>
                <div class="col-sm-7">
                  <input required v-model="newMACAuth.name" type="text" id="textInput-modal-markup" class="form-control" :placeholder="$t('unit.description')">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-5 control-label" for="textInput2-modal-markup">{{ $t("unit.mac_address") }}</label>
                <div class="col-sm-7">
                  <input required v-model="newMACAuth.username" pattern="^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})" type="text" id="textInput2-modal-markup"
                    class="form-control" placeholder="00:11:22:AA:BB:CC">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-5 control-label" for="textInput-modal-markup">{{$t('hotspot.bandwidth_down')}}</label>
                <div class="col-sm-7">
                  <input v-model="newMACAuth.kbps_down" type="number" class="form-control">
                </div>
              </div>
              <div class="form-group">
                <label class="col-sm-5 control-label" for="textInput-modal-markup">{{$t('hotspot.bandwidth_up')}}</label>
                <div class="col-sm-7">
                  <input v-model="newMACAuth.kbps_up" type="number" class="form-control">
                </div>
              </div>
            </div>
            <div class="modal-footer">
              <span v-show="newMACAuth.onAction" class="spinner spinner-sm spinner-inline modal-spinner"></span>
              <button type="button" class="btn btn-default" data-dismiss="modal">{{ $t("cancel") }}</button>
              <button type="submit" class="btn btn-primary">{{ $t("create") }}</button>
            </div>
          </form>
        </div>
      </div>
    </div>

     <div v-if="!info.isLoading" v-for="voucher in vouchers.data" v-bind:key="voucher.id" class="card-pf" id="voucher-coupon">
            <img class="voucher-logo" src="/static/logo.png" />
            <h2 class="card-pf-title voucher-name" id="test">
                {{info.data.name}}
            </h2>
            <h3 class="card-pf-title voucher-desc">
                {{info.data.description}}
            </h3>
            <div class="card-pf-body voucher-main">
                <p>
                    <strong class="voucher-code">{{voucher.code}}</strong>
                </p>
            </div>
            <div class="card-pf-footer voucher-details">
                <div class="card-pf-time-frame-filter voucher-valid">
                    {{ $t("hotspot.valid") }}:
                    <strong>{{voucher.duration}}</strong> <span> {{ $t("hotspot.days") }} </span>
                </div>
                <div class="voucher-download">
                    <span class="fa fa-arrow-down"></span> Download: <strong>{{voucher.bandwidth_down}} </strong> <span> Kb/s</span>
                </div>
                <div class="voucher-upload">
                    <span class="fa fa-arrow-up"></span> Upload: <strong>{{voucher.bandwidth_up}} </strong> <span>Kb/s</span>
                </div>
            </div>
        </div>
  </div>
</template>

<script>
  import HotspotService from '../../services/hotspot';
  import PreferenceService from '../../services/preference';
  import AccountService from '../../services/account';
  import UnitService from '../../services/unit';
  import UserService from '../../services/user';
  import DeviceService from '../../services/device';
  import SessionService from '../../services/session';
  import StorageService from '../../services/storage';
  import UtilService from '../../services/util';

  import jsPDF from 'jspdf'
  import CaptivePortal from '../../directives/CaptivePortal.vue'
  import HotspotAction from '../../directives/HotspotAction.vue';
  import PictureInput from 'vue-picture-input'
  import arrow from '../../../static/arrows.js';
  import {
    Sketch
  } from 'vue-color'
  import {
    VueEditor
  } from 'vue2-editor'
  import {
    setTimeout
  } from 'timers';

  export default {
    name: 'HotspotDetails',
    mixins: [HotspotService, PreferenceService, AccountService, UnitService, UserService, DeviceService, SessionService,
      StorageService, UtilService
    ],
    components: {
      hotspotAction: HotspotAction,
      PictureInput,
      'sketch-picker': Sketch,
      captivePortal: CaptivePortal,
      VueEditor
    },
    data() {
      // get hotspot info
      this.getInfo()

      // get totals
      this.getTotals()

      // get preferences
      this.getPreferences()

      // get vouchers
      this.getVouchers()

      // get mac auths
      this.getAllMACAuth()

      return {
        info: {
          isLoading: true,
          data: {}
        },
        vouchers: {
          isLoading: true,
          data: []
        },
        macAuth: {
          isLoading: true,
          data: []
        },
        vouchersCount: 1,
        newVoucher: {
          auto_login: true,
          bandwidth_down: 0,
          bandwidth_up: 0,
          duration: 7
        },
        newMACAuth: {
          name: '',
          username: '',
          kbps_down: 0,
          kbps_up: 0,
          onAction: false
        },
        preferences: {
          isLoading: true,
          global: {},
          captive: {}
        },
        totals: {
          accounts: {
            isLoading: true,
            count: 0
          },
          units: {
            isLoading: true,
            count: 0,
            data: []
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
          }
        },
        columns: [{
            label: this.$i18n.t('hotspot.code'),
            field: 'code',
            filterable: false,
            sortable: false,
          }, {
            label: this.$i18n.t('hotspot.auto_login'),
            field: 'auto_login',
            filterable: false,
          }, {
            label: this.$i18n.t('hotspot.bandwidth_limit'),
            field: 'bandwidth',
            filterable: false,
            sortable: false,
          }, {
            label: this.$i18n.t('hotspot.duration'),
            field: 'duration',
            filterable: false,
            sortable: false,
          },
          {
            label: this.$i18n.t('hotspot.used'),
            field: 'expires',
            filterable: false,
            sortable: false,
          },
          {
            label: '',
            field: '',
            sortable: false
          },
        ],
        columnsMAC: [{
            label: this.$i18n.t('unit.description'),
            field: 'name',
            filterable: false,
            sortable: false,
          }, {
            label: this.$i18n.t('unit.mac_address'),
            field: 'username',
            filterable: false,
          }, {
            label: this.$i18n.t('hotspot.bandwidth_limit'),
            field: 'bandwidth',
            filterable: false,
            sortable: false,
          }, {
            label: this.$i18n.t('hotspot.created'),
            field: 'created',
            filterable: false,
            sortable: false,
          },
          {
            label: '',
            field: '',
            sortable: false
          },
        ],
        tableLangsTexts: this.tableLangs(),
        uploadLangstexts: this.uploadImageLangs(),
        user: this.get("loggedUser"),
        displayCaptivePortalOptions: this.get("loggedUser") && this.get("loggedUser").subscription && this.get(
            "loggedUser").subscription.subscription_plan && this.get("loggedUser").subscription.subscription_plan.wings_customization ||
          this.get("loggedUser").account_type == "admin",
        customToolbar: [
          ['bold', 'italic', 'underline'],
          ['image', 'code-block']
        ],
        userLink: {
          name: 'Users',
          params: {
            hotspotId: this.$route.params.id
          }
        },
        unitLink: {
          name: 'Units',
          params: {
            hotspotId: this.$route.params.id
          }
        },
        sessionLink: {
          name: 'Sessions',
          params: {
            hotspotId: this.$route.params.id
          }
        },
        accountLink: {
          name: 'Accounts',
          params: {
            hotspotId: this.$route.params.id
          }
        },
        deviceLink: {
          name: 'Devices',
          params: {
            hotspotId: this.$route.params.id
          }
        }
      }
    },
    methods: {
      getPrevieHTML(value) {
        return '<img src="' + value + '"></img>'
      },
      getPrefIcon(pref) {
        return this.getPrefTypeIcon(pref)
      },
      createVoucher() {
        this.vouchers.isLoading = true

        var promises = []
        var context = this
        for (var i = 0; i < this.vouchersCount; i++) {
          promises.push(new Promise(function (resolve, reject) {
            context.hotspotCreateVoucher({
              hotspot_id: parseInt(context.$route.params.id),
              code: context.generateVoucher(),
              auto_login: context.newVoucher.auto_login,
              bandwidth_down: parseInt(context.newVoucher.bandwidth_down),
              bandwidth_up: parseInt(context.newVoucher.bandwidth_up),
              duration: parseInt(context.newVoucher.duration),
            }, success => {
              resolve()
            }, error => {
              console.log(error.body)
              reject()
            })
          }))
        }
        Promise.all(promises).then(function () {
          context.vouchers.isLoading = false
          context.getVouchers()
          $('#voucherModal').modal('hide')
        }).catch(function (err) {
          console.error(err)
          context.vouchers.isLoading = false
        })
      },
      deleteVoucher(id) {
        this.vouchers.isLoading = true
        this.hotspotVoucherDelete(id, success => {
          this.vouchers.isLoading = false
          this.getVouchers()
        }, error => {
          console.log(error.body)
          this.vouchers.isLoading = false
          this.getVouchers()
        })
      },
      getVouchers() {
        this.hotspotGetVouchers(this.$route.params.id, success => {
          this.vouchers.data = success.body
          this.vouchers.isLoading = false
        }, error => {
          console.log(error.body)
          this.vouchers.data = []
          this.vouchers.isLoading = false
        })
      },
      checkVoucherUse(value) {
        return +new Date(value) > 0
      },
      getInfo() {
        this.hotspotGet(this.$route.params.id, success => {
          this.info.data = success.body
          this.info.isLoading = false
        }, error => {
          console.log(error.body)
        })
      },
      getTotals() {
        this.accountGetAll(this.$route.params.id, success => {
          this.totals.accounts.count = success.body.length
          this.totals.accounts.isLoading = false
        }, error => {
          console.log(error.body)
          this.totals.accounts.isLoading = false
        })
        this.unitGetAll(this.$route.params.id, success => {
          this.totals.units.data = success.body
          this.totals.units.count = success.body.length
          this.totals.units.isLoading = false
        }, error => {
          console.log(error.body)
          this.totals.units.isLoading = false
        })
        this.userGetAll(this.$route.params.id, null, success => {
          this.totals.users.count = success.body.length
          this.totals.users.isLoading = false
        }, error => {
          console.log(error.body)
          this.totals.users.isLoading = false
        })
        this.deviceGetAll(this.$route.params.id, "", success => {
          this.totals.devices.count = success.body.length
          this.totals.devices.isLoading = false
        }, error => {
          console.log(error.body)
          this.totals.devices.isLoading = false
        })
        this.sessionGetAll(this.$route.params.id, null, null, null, null, success => {
          this.totals.sessions.count = success.body.length
          this.totals.sessions.isLoading = false
        }, error => {
          console.log(error.body)
          this.totals.sessions.isLoading = false
        })
      },
      getPreferences() {
        this.hsPrefGet(this.$route.params.id, success => {
          var globalPref = []
          var captivePref = []
          var backgroundColor = ''

          for (var p in success.body) {
            var pref = success.body[p]
            if (pref.value === "true") {
              pref.value = true
            }
            if (pref.value === "false") {
              pref.value = false
            }

            if (pref.key == 'voucher_login' && pref.value) {
              this.preferences.vouchersAvailable = true
            } else {
              this.preferences.vouchersAvailable = false
            }

            if (pref.key == 'captive_7_background') {
              backgroundColor = pref.value
            }

            if (pref.key.startsWith('captive')) {
              captivePref.push(pref)
            } else {
              globalPref.push(pref)
            }
          }

          this.preferences.global = globalPref
          this.preferences.captive = captivePref
          this.preferences.isLoading = false
          setTimeout(function () {
            window.$('#captive-preview').css('background-color', backgroundColor);
          }, 0)
        }, error => {
          console.log(error.body)
        })
      },
      updatePreferences() {
        this.preferences.isLoading = true
        // create promises array
        var promises = []
        for (var i in this.preferences.global) {
          var pref = this.preferences.global[i]
          promises.push(new Promise((resolve, reject) => {
            if (typeof pref.value == "boolean") {
              pref.value = pref.value.toString()
            }
            this.hsPrefModify(this.$route.params.id, pref, success => {
              resolve(success)
            }, error => {
              reject(error)
            })
          }))
        }

        // exec promises
        var context = this;
        Promise.all(promises).then(function (response) {
          context.preferences.isLoading = false
          context.getPreferences()
        })
      },
      updatePreferencesCaptive() {
        this.preferences.isLoading = true
        // create promises array
        var promises = []
        for (var i in this.preferences.captive) {
          var pref = this.preferences.captive[i]
          promises.push(new Promise((resolve, reject) => {
            if (typeof pref.value == "boolean") {
              pref.value = pref.value.toString()
            }
            if (pref.key == 'captive_7_background') {
              pref.value = pref.value.hex || pref.value
            }
            this.hsPrefModify(this.$route.params.id, pref, success => {
              resolve(success)
            }, error => {
              reject(error)
            })
          }))
        }

        // exec promises
        var context = this;
        Promise.all(promises).then(function (response) {
          context.preferences.isLoading = false
          context.getPreferences()
        })
      },
      printVoucher(voucher) {
        let index
        for (let i = 0; i < this.vouchers.data.length; i++) {
          if(this.vouchers.data[i].id === voucher.id){
            index = i;
          }
        }

        var doc = new jsPDF("portrait", "mm", "a4");
        var halfWidth = Math.round(doc.internal.pageSize.width / 2);
        var fifthHeight = Math.round(doc.internal.pageSize.height / 5);
        
        doc.setDrawColor(17, 17, 17)
        doc.line(0, fifthHeight, halfWidth, fifthHeight);
        doc.line(halfWidth, 0, halfWidth, fifthHeight);
        doc.addImage(this.preferences.captive[2].value, 2, 2, 20, 20);
        doc.fromHTML(document.getElementsByClassName('voucher-desc')[index], 50, -3, {
            width: 55
        });
        doc.fromHTML(document.getElementsByClassName('voucher-main')[index], 35, 18);
        doc.setLineWidth(0.3)
        doc.setDrawColor(158, 160, 163);
        doc.line(5,  35, halfWidth - 5, + 35)
        doc.addImage(arrow.down, 4, 41, 3, 3);
        doc.fromHTML(document.getElementsByClassName('voucher-download')[index], 8, 36);
        doc.addImage(arrow.up, 4, 49, 3, 3);
        doc.fromHTML(document.getElementsByClassName('voucher-upload')[index], 8, 44);
        doc.fromHTML(document.getElementsByClassName('voucher-valid')[index], 65, 36);

        doc.autoPrint();
        window.open(doc.output('bloburl'), '_blank');
          
      },
      printAllVoucher() {
        var doc = new jsPDF("portrait", "mm", "a4");
        var width = doc.internal.pageSize.width;
        var height = doc.internal.pageSize.height;
        var halfWidth = Math.round(width / 2);
        var fifthHeight = Math.round(height / 5);
        var pageNumber = (15 % 8 !== 0) ? (Math.floor(15 / 8) + 1) : (15 / 8);
        var row = 0;
        var cordinates = {
            y: 0,
            x: 0,
            width: 0,
            height: 0,
        }

        for (var index = 0; index < this.vouchers.data.length; index++) {
          if (index % 10 === 0 && index !== 0) {
              console.log('index: ', index);
              doc.addPage();
              row = 0;
              cordinates = {
                  y: 0,
                  x: 0,
                  width: 0,
                  height: 0,
              }
          }
            doc.setLineWidth(0.3)
            doc.setDrawColor(17, 17, 17)
            // Left column
            if (index % 2 === 0) {
              doc.addImage(this.preferences.captive[2].value, 2, cordinates.y + 2, 20, 20);
              doc.fromHTML(document.getElementsByClassName('voucher-desc')[index], 50, cordinates.y + (-3), {
                  width: 55
                  
              });
              doc.fromHTML(document.getElementsByClassName('voucher-main')[index], 35, cordinates.y + 18);
              doc.setLineWidth(0.3)
              doc.setDrawColor(158, 160, 163);
              doc.line(5, cordinates.y + 35, halfWidth - 5, cordinates.y + 35)
              doc.addImage(arrow.down, 4, cordinates.y + 41, 3, 3);
              doc.fromHTML(document.getElementsByClassName('voucher-download')[index], 8, cordinates.y + 36);
              doc.addImage(arrow.up, 4, cordinates.y + 49, 3, 3);
              doc.fromHTML(document.getElementsByClassName('voucher-upload')[index], 8, cordinates.y + 44);
              doc.fromHTML(document.getElementsByClassName('voucher-valid')[index], 65, cordinates.y + 36);

              doc.setDrawColor(17, 17, 17)

              if(index % 10 === 8){
                doc.line(halfWidth, cordinates.y, halfWidth, cordinates.y + fifthHeight);
              }else{
                doc.line(0, cordinates.y + fifthHeight, halfWidth, cordinates.y + fifthHeight);
                doc.line(halfWidth, cordinates.y, halfWidth, cordinates.y + fifthHeight);
              }

              
              // Right column
            } else {
              cordinates.x = halfWidth;
              doc.addImage(this.preferences.captive[2].value, cordinates.x + 2, cordinates.y + 2, 20, 20);
              doc.fromHTML(document.getElementsByClassName('voucher-desc')[index], cordinates.x + 50, cordinates.y + (-3), {
                  width: 55
              });
              doc.fromHTML(document.getElementsByClassName('voucher-main')[index], cordinates.x + 35, cordinates.y + 18);
              doc.setLineWidth(0.3)
              doc.setDrawColor(158, 160, 163);
              doc.line(cordinates.x + 5, cordinates.y + 35, (cordinates.x * 2) - 5, cordinates.y + 35)
              doc.addImage(arrow.down, cordinates.x + 4, cordinates.y + 41, 3, 3);
              doc.fromHTML(document.getElementsByClassName('voucher-download')[index], cordinates.x + 8, cordinates.y + 36);
              doc.addImage(arrow.up, cordinates.x + 4, cordinates.y + 49, 3, 3);
              doc.fromHTML(document.getElementsByClassName('voucher-upload')[index], cordinates.x + 8, cordinates.y + 44);
              doc.fromHTML(document.getElementsByClassName('voucher-valid')[index], cordinates.x + 65, cordinates.y + 36);
              
              row++;
              cordinates.width = cordinates.y;
              cordinates.y = fifthHeight * row;
              doc.setDrawColor(17, 17, 17)

              if(!(index % 10 === 9)){
                doc.line(halfWidth, cordinates.y, width, cordinates.y);
              }
            }
        }
        doc.autoPrint();
        window.open(doc.output('bloburl'), '_blank');
      },
      onChanged(pref) {
        if (this.$refs['prefInput-' + pref.key][0].image.length > 655360) {
          pref.onError = true
          this.$refs['prefInput-' + pref.key][0].image = pref.value
        } else {
          pref.onError = false
          pref.value = this.$refs['prefInput-' + pref.key][0].image
        }
        this.$forceUpdate()
      },
      onUpdate(value) {
        $('#captive-preview').css('background-color', value.hex);
      },
      getAllMACAuth() {
        this.userGetAll(this.$route.params.id, 'mac', success => {
          this.macAuth.data = success.body
          this.macAuth.isLoading = false
        }, error => {
          this.macAuth.isLoading = false
          this.macAuth.data = []
          console.log(error)
        })
      },
      createMacAuth() {
        this.newMACAuth.onAction = true
        var md5 = require('md5');
        var digest = md5(this.newMACAuth.unit.secret + this.newMACAuth.unit.uuid)
        console.log(this.newMACAuth)
        this.userMACCreate(this.newMACAuth, digest, this.newMACAuth.unit.uuid, success => {
          $('#macAuthModal').modal('toggle');
          this.newMACAuth = {
            name: '',
            username: '',
            kbps_down: 0,
            kbps_up: 0,
            onAction: false
          }
          this.getAllMACAuth()
          this.newMACAuth.onAction = false
        }, error => {
          console.log(error.body);
          this.newMACAuth.onAction = false
        })
      },
      deleteMacAuth(id) {
        this.userDelete(id, success => {
          this.getAllMACAuth()
        }, error => {
          console.log(error.body);
        })
      }
    }
  }

</script>

<style scoped>


  textarea {
    width: 100%;
    min-height: 180px;
    resize: vertical;
  }

  #voucher-coupon{
    display: none;
  }
  .voucher-logo {
      height: 25px;
      margin: 15px;
      margin-left: -5px;
  }
  .voucher-name {
      position: absolute;
      top: -8px;
      left: 90px;
      max-width: 200px;
      color: goldenrod;
      
  }
  .voucher-desc {
      position: absolute;
      top: 15px;
      left: 90px;
      font-size: 21px;
      color: #757575;
      max-width: 200px;
  }
  .voucher-code {
      font-size: 30px;
      margin-top: 5px;
  }
  .voucher-valid {
      font-size: 18px;
      margin-top: 6px;
  }
  .voucher-main {
      text-align: center;
      padding-bottom: 10px;
      margin-top: 15px;
  }
  .voucher-details {
      padding: 10px;
      background: #fff;
  }
  .voucher-upload {
      font-size: 18px;
  }
  .voucher-download {
      font-size: 18px;
  }


</style>
