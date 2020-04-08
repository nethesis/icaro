<template>
  <div>
    <h2>
      Hotspot
      <strong class="soft">{{ info.data.name }} - {{info.data.description}}</strong>
    </h2>
    <div class="row row-cards-pf">
      <div class="col-xs-12 col-sm-4 col-md-4 col-lg-3 max-height-80">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              <router-link class="card-link" :to="accountLink">
                <span class="pficon pficon-users card-info-title"></span>
                {{ $t("dashboard.accounts") }}
                <span
                  v-if="!totals.accounts.isLoading"
                  class="right"
                >
                  <strong class="soft">{{ totals.accounts.count}}</strong>
                </span>
              </router-link>
              <div v-if="totals.accounts.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
        </div>
      </div>
      <div class="col-xs-12 col-sm-4 col-md-4 col-lg-3 max-height-80">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              <router-link class="card-link" :to="unitLink">
                <span class="pficon pficon-connected card-info-title"></span>
                {{ $t("dashboard.units") }}
                <span
                  v-if="!totals.units.isLoading"
                  class="right"
                >
                  <strong class="soft">{{ totals.units.count}}</strong>
                </span>
              </router-link>
              <div v-if="totals.units.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
        </div>
      </div>
      <div class="col-xs-12 col-sm-4 col-md-4 col-lg-3 max-height-80">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              <router-link class="card-link" :to="userLink">
                <span class="fa fa-users card-info-title"></span>
                {{ $t("dashboard.users") }}
                <span
                  v-if="!totals.users.isLoading"
                  class="right"
                >
                  <strong class="soft">{{ totals.users.count}}</strong>
                </span>
              </router-link>
              <div v-if="totals.users.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
        </div>
      </div>
      <div class="col-xs-12 col-sm-4 col-md-4 col-lg-3 max-height-80">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              <router-link class="card-link" :to="deviceLink">
                <span class="fa fa-laptop card-info-title"></span>
                {{ $t("dashboard.devices") }}
                <span
                  v-if="!totals.devices.isLoading"
                  class="right"
                >
                  <strong class="soft">{{ totals.devices.count}}</strong>
                </span>
              </router-link>
              <div v-if="totals.devices.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
        </div>
      </div>
      <div class="col-xs-12 col-sm-4 col-md-4 col-lg-3 max-height-80">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              <router-link class="card-link" :to="sessionLink">
                <span class="fa fa-list card-info-title"></span>
                {{ $t("dashboard.sessions") }}
                <span
                  v-if="!totals.sessions.isLoading"
                  class="right"
                >
                  <strong class="soft">{{ totals.sessions.count}}</strong>
                </span>
              </router-link>
              <div v-if="totals.sessions.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
        </div>
      </div>
      <div class="col-xs-12 col-sm-6 col-md-6 col-lg-3 max-height-80">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              <span class="fa fa-commenting card-info-title"></span>
              {{ $t("dashboard.sms_sent") }}
              <div
                v-if="!totals.sms.isLoading"
                class="right"
              >
                <div
                  class="text-align-right"
                  :class="totals.sms.count < smsMaxCount || smsMaxCount == 0 ? 'soft' : 'red'"
                >
                  {{ totals.sms.count}} /
                  <b>{{smsMaxCount == 0 ? '-' : smsMaxCount}}</b>
                </div>
                <div v-if="smsThreshold > 0" class="small-text">
                  {{ $t("hotspot.warn_sms") }}: <b>{{ smsThreshold }}</b>
                </div>
              </div>
              <div v-if="totals.sms.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
        </div>
      </div>
    </div>
    <div class="row row-cards-pf">
      <div
        v-if="preferences.vouchersAvailable"
        :class="['col-xs-12 col-sm-12', user.account_type == 'admin' || user.account_type == 'reseller' ? 'col-md-12' : 'col-md-12']"
      >
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{ $t("hotspot.vouchers") }}
              <div v-if="vouchers.isLoading" class="spinner spinner-sm right"></div>
              <span @click="getVouchers()" class="fa fa-refresh right link"></span>
            </h2>
          </div>
          <div v-if="!vouchers.isLoading" class="card-pf-body">
            <div class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
              <label
                class="col-sm-3 control-label centered"
                for="textInput-markup"
              >{{$t('hotspot.code')}}</label>
              <div class="col-sm-9">
                <input v-model="vouchers.filters.code" type="text" class="form-control" />
              </div>
            </div>
            <div class="form-group col-xs-12 col-sm-12 col-md-12 col-lg-12 no-padding mg-left-20">
              <legend class="fields-section-header-pf" aria-expanded="true">
                <span
                  :class="['fa fa-angle-right field-section-toggle-pf', advancedFilters ? 'fa-angle-down' : '']"
                  @click="toggleAdvancedFilters()"
                ></span>
                <a
                  class="field-section-toggle-pf pointer"
                  @click="toggleAdvancedFilters()"
                >{{$t('hotspot.advanced_filters')}}</a>
              </legend>
            </div>
            <div v-show="advancedFilters">
              <div class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
                <label class="col-sm-3 control-label centered" for="textInput-markup">
                  {{$t('hotspot.duration')}}
                  ({{$t('hotspot.0_all')}})
                </label>
                <div class="col-sm-9">
                  <input v-model="vouchers.filters.duration" type="number" class="form-control" />
                </div>
              </div>
              <div class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
                <label
                  class="col-sm-3 control-label centered"
                  for="textInput-markup"
                >{{$t('hotspot.expires')}}</label>
                <div class="width-p50 flex">
                  <label class="label-dateFilter">{{$t('hotspot.from')}}</label>
                  <datepicker
                    :format="dateFormatter"
                    v-model="vouchers.filters.expiredStart"
                    :language="locale"
                    class="dateFilter"
                    input-class="datepicker-input"
                  ></datepicker>
                  <label class="label-dateFilter mg-left-20">{{$t('hotspot.to')}}</label>
                  <datepicker
                    :format="dateFormatter"
                    v-model="vouchers.filters.expiredEnd"
                    :language="locale"
                    class="dateFilter"
                    input-class="datepicker-input"
                  ></datepicker>
                  <button
                    @click="clearExpiredFilters()"
                    class="btn btn-default mg-left-20"
                  >{{$t('hotspot.clear')}}</button>
                </div>
              </div>
              <div class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
                <label
                  class="col-sm-3 control-label centered"
                  for="textInput-markup"
                >{{$t('hotspot.created')}}</label>
                <div class="width-p50 flex">
                  <label class="label-dateFilter">{{$t('hotspot.from')}}</label>
                  <datepicker
                    :format="dateFormatter"
                    v-model="vouchers.filters.createdStart"
                    :language="locale"
                    class="dateFilter"
                    input-class="datepicker-input"
                  ></datepicker>
                  <label class="label-dateFilter mg-left-20">{{$t('hotspot.to')}}</label>
                  <datepicker
                    :format="dateFormatter"
                    v-model="vouchers.filters.createdEnd"
                    :language="locale"
                    class="dateFilter"
                    input-class="datepicker-input"
                  ></datepicker>
                  <button
                    @click="clearCreatedFilters()"
                    class="btn btn-default mg-left-20"
                  >{{$t('hotspot.clear')}}</button>
                </div>
              </div>
              <div class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
                <label
                  class="col-sm-3 control-label centered"
                  for="textInput-markup"
                >{{$t('hotspot.auto_login')}}</label>
                <div class="col-sm-3 list-details">
                  <select v-model="vouchers.filters.auto_login" class="form-control">
                    <option value>-</option>
                    <option value="1">{{$t('hotspot.yes')}}</option>
                    <option value="0">{{$t('hotspot.no')}}</option>
                  </select>
                </div>
                <label
                  class="col-sm-3 control-label centered"
                  for="textInput-markup"
                >{{$t('hotspot.used')}}</label>
                <div class="col-sm-3 list-details">
                  <select v-model="vouchers.filters.used" class="form-control">
                    <option value>-</option>
                    <option value="1">{{$t('hotspot.yes')}}</option>
                    <option value="0">{{$t('hotspot.no')}}</option>
                  </select>
                </div>

                <label
                  class="col-sm-3 control-label centered"
                  for="textInput-markup"
                >{{$t('hotspot.reusable')}}</label>
                <div class="col-sm-3 list-details">
                  <select v-model="vouchers.filters.reusable" class="form-control">
                    <option value>-</option>
                    <option value="1">{{$t('hotspot.yes')}}</option>
                    <option value="0">{{$t('hotspot.no')}}</option>
                  </select>
                </div>

                <label
                  class="col-sm-3 control-label centered"
                  for="textInput-markup"
                >{{$t('hotspot.printed')}}</label>
                <div class="col-sm-3 list-details">
                  <select v-model="vouchers.filters.printed" class="form-control">
                    <option value>-</option>
                    <option value="1">{{$t('hotspot.yes')}}</option>
                    <option value="0">{{$t('hotspot.no')}}</option>
                  </select>
                </div>

                <label
                  class="col-sm-3 control-label centered"
                  for="textInput-markup"
                >{{$t('hotspot.type')}}</label>
                <div class="col-sm-3 list-details">
                  <select v-model="vouchers.filters.type" class="form-control">
                    <option value>-</option>
                    <option value="normal">{{$t('hotspot.normal')}}</option>
                    <option value="auth">{{$t('hotspot.auth')}}</option>
                  </select>
                </div>
              </div>
            </div>

            <div class="form-group col-xs-3 col-sm-3 col-md-3 col-lg-3">
              <button @click="getVouchers()" class="btn btn-primary">{{$t('hotspot.apply_filters')}}</button>
            </div>

            <vue-good-table
              :perPage="5"
              :paginate="true"
              :columns="columns"
              :rows="vouchers.data"
              :lineNumbers="false"
              :defaultSortBy="{field: 'created', type: 'desc'}"
              styleClass="table"
              :nextText="tableLangsTexts.nextText"
              :prevText="tableLangsTexts.prevText"
              :rowsPerPageText="tableLangsTexts.rowsPerPageText"
              :globalSearchPlaceholder="tableLangsTexts.globalSearchPlaceholder"
              :ofText="tableLangsTexts.ofText"
            >
              <template slot="table-row" slot-scope="props">
                <td :class="['fancy', 'td-voucher-'+props.row.type]">
                  <strong>{{ props.row.code | uppercase }}</strong>
                  <br />
                  {{ props.row.user_name}} ({{ props.row.user_mail }})
                </td>
                <td :class="['fancy', 'td-voucher-'+props.row.type]">
                  <strong>{{ $t('hotspot.'+props.row.type) }}</strong>
                </td>
                <td :class="['fancy', 'td-voucher-'+props.row.type]">
                  <span
                    :class="['pficon', props.row.auto_login ? 'pficon-ok' : 'pficon-error-circle-o']"
                  ></span>
                </td>
                <td :class="['fancy', 'td-voucher-'+props.row.type]">
                  <div>
                    <strong>{{ $t('user.kbps_down') }}</strong>
                    : {{ props.row.bandwidth_down || '-' }}
                  </div>
                  <div>
                    <strong>{{ $t('user.kbps_up') }}</strong>
                    : {{ props.row.bandwidth_up || '-' }}
                  </div>
                </td>
                <td :class="['fancy', 'td-voucher-'+props.row.type]">
                  <div>
                    <strong>{{ $t('user.traffic') }}</strong>
                    : {{ props.row.max_traffic || '-' | byteFormat }}
                  </div>
                  <div>
                    <strong>{{ $t('user.time') }}</strong>
                    : {{ props.row.max_time || '-' | secondsInHour }}
                  </div>
                </td>
                <td
                  :class="['fancy', 'td-voucher-'+props.row.type]"
                >{{ props.row.created | formatDateShort(true) }}</td>
                <td :class="['fancy', 'td-voucher-'+props.row.type]">
                  <span v-if="props.row.duration > 0">{{ props.row.duration }}</span>
                  <span v-if="props.row.duration == 0">-</span>
                </td>
                <td :class="['fancy', 'td-voucher-'+props.row.type]">
                  <span :class="['fa', checkVoucherUse(props.row.expires) ? 'fa-check green' : '']"></span>
                  {{ props.row.expires | formatDateShort(true) }}
                </td>
                <td :class="['fancy', 'td-voucher-'+props.row.type]">
                  {{props.row.remain_use == -1 ? $t('hotspot.limitless') : $t('hotspot.max_use') + ': ' }}
                  <strong
                    v-if="props.row.remain_use > 0"
                  >{{props.row.remain_use}}</strong>
                  <strong
                    class="red"
                    v-if="props.row.remain_use == 0"
                  >{{$t('hotspot.limit_reached')}}</strong>
                </td>
                <td :class="['fancy', 'td-voucher-'+props.row.type]">
                  <button
                    v-if="props.row.remain_use != 0"
                    v-on:click="printVoucher(props.row)"
                    :class="['btn', props.row.printed ? '' : 'btn-primary']"
                    type="button"
                  >
                    <span class="fa fa-print"></span>
                  </button>
                  <button
                    v-on:click="deleteVoucher(props.row.id)"
                    class="btn btn-danger"
                    type="button"
                  >
                    <span class="fa fa-remove"></span>
                  </button>
                </td>
              </template>
            </vue-good-table>
          </div>
          <div v-if="!vouchers.isLoading" class="card-pf-footer">
            <div class="dropdown card-pf-time-frame-filter">
              <button
                data-toggle="modal"
                data-target="#voucherModalCustom"
                class="btn btn-primary orange"
                type="button"
              >
                {{
                $t("hotspot.create_voucher_custom") }}
              </button>
              <button
                data-toggle="modal"
                data-target="#voucherModal"
                class="btn btn-primary"
                type="button"
              >
                {{
                $t("hotspot.create_voucher") }}
              </button>
              <button
                :disabled="vouchers.data.length == 0"
                v-on:click="printAllVoucher()"
                class="btn btn-default"
                type="button"
              >
                <span class="fa fa-print"></span>
                {{ $t("hotspot.print_all_voucher") }}
              </button>
              <button
                :disabled="vouchers.data.length == 0"
                v-on:click="exportCSVVoucher()"
                class="btn btn-default"
                type="button"
              >
                <span class="fa fa-list"></span>
                {{ $t("hotspot.export_csv") }}
              </button>
              <button
                :disabled="vouchers.data.length == 0"
                data-toggle="modal"
                data-target="#voucherDeleteAll"
                class="btn btn-danger"
                type="button"
              >{{ $t("hotspot.delete_all") }}</button>
            </div>
            <p>
              <a href="#" class="card-pf-link-with-icon"></a>
            </p>
          </div>
        </div>
      </div>
    </div>

    <div
      v-if="user.account_type == 'admin' || user.account_type == 'reseller'"
      class="row row-cards-pf"
    >
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{ $t("hotspot.preferences") }}
              <div v-if="preferences.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
          <form
            class="form-horizontal"
            role="form"
            v-on:submit.prevent="updatePreferences(preferences.global)"
          >
            <div v-if="!preferences.isLoading" class="card-pf-body">
              <div v-for="pref in preferences.global" :key="pref.key" class="form-group">
                <label class="col-sm-4 control-label">
                  {{$t('hotspot.'+pref.key)}}
                  <span :class="[getPrefIcon(pref.key)]"></span>
                </label>
                <div :class="pref.key == 'facebook_login_page' ? ['col-sm-4', 'col-lg-3'] : ['col-sm-2', 'col-lg-1']">
                  <input
                    v-model="pref.value"
                    :type="getInputType(pref.key, pref.value)"
                    id="textInput-markup"
                    class="form-control"
                  />
                </div>
                <div v-if="pref.key == 'sms_login_max'" class="col-sm-4">
                  <label class="control-label col-sm-3">{{$t('hotspot.add_sms_count')}}:</label>
                  <input
                    v-model="smsMaxCountAdd"
                    type="number"
                    id="textInput-markup"
                    class="form-control col-sm-1 special-input"
                  />
                  <button
                    type="button"
                    @click="addSMSCount(pref, smsMaxCountAdd)"
                    class="col-sm-2 btn btn-primary"
                  >{{$t('hotspot.add')}}</button>
                </div>
              </div>
            </div>
            <div v-if="!preferences.isLoading" class="card-pf-footer">
              <div class="dropdown card-pf-time-frame-filter">
                <button type="submit" class="btn btn-primary">{{ $t("update") }}</button>
              </div>
              <p>
                <a href="#" class="card-pf-link-with-icon"></a>
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
          <form
            class="form-horizontal"
            role="form"
            v-on:submit.prevent="updatePreferencesCaptive(preferences.captive)"
          >
            <div v-if="!preferences.isLoading" class="card-pf-body">
              <!-- redirect preference: captive_1_redir -->
              <div :class="[captivePrefsSorted[0].onError ? 'has-error' : '', 'form-group', 'mg-bottom-3r']">
                <label class="col-sm-4 control-label">
                  {{$t('hotspot.'+captivePrefsSorted[0].key)}}
                  <span :class="[getPrefIcon(captivePrefsSorted[0].key)]"></span>
                </label>
                <div class="col-sm-6">
                  <input
                    required
                    v-model="captivePrefsSorted[0].value"
                    :type="getInputType(captivePrefsSorted[0].key, captivePrefsSorted[0].value)"
                    class="form-control"
                  />
                </div>
              </div>
              <!-- text preferences: captive_2_title, captive_4_subtitle, captive_6_description -->
              <div class="row">
                <div
                  v-for="pref in [captivePrefsSorted[1], captivePrefsSorted[3], captivePrefsSorted[5]]"
                  :key="pref.key"
                  :class="[pref.onError ? 'has-error' : '', 'col-md-4', 'captive-portal-pref-inline', 'mg-bottom-3r']"
                >
                  <label class="block control-label">
                    {{$t('hotspot.'+pref.key)}}
                    <span :class="[getPrefIcon(pref.key)]"></span>
                  </label>
                  <div>
                    <input
                      v-model="pref.value"
                      :type="getInputType(pref.key, pref.value)"
                      class="form-control"
                    />
                  </div>
                </div>
              </div>
              <!-- image preferences: captive_3_logo, captive_5_banner, captive_81_bg_image -->
              <div class="row">
                <div
                  v-for="pref in [captivePrefsSorted[2], captivePrefsSorted[4], captivePrefsSorted[6]]"
                  :key="pref.key"
                  :class="[pref.onError ? 'has-error' : '', 'col-md-4', 'captive-portal-pref-inline', 'mg-bottom-3r']"
                >
                  <label class="block-centered control-label">
                    {{$t('hotspot.'+pref.key)}}
                    <span :class="[getPrefIcon(pref.key)]"></span>
                  </label>
                  <div :id="pref.key">
                    <picture-input
                      :ref="'prefInput-'+pref.key"
                      :prefill="urltoFile(pref.value, pref.key)"
                      :alertOnError="false"
                      @change="onPictureChanged(pref)"
                      @remove="onPictureRemoved(pref)"
                      @prefill="onPicturePrefill(pref)"
                      :width="100"
                      :height="100"
                      :crop="false"
                      :zIndex="1000"
                      :customStrings="uploadLangstexts"
                      :removable="pref.value ? true : false"
                      removeButtonClass="btn btn-danger"
                      buttonClass="btn btn-default"
                    ></picture-input>
                    <span v-if="pref.onError" class="help-block">{{$t('upload_file_exceed')}}</span>
                  </div>
                </div>
              </div>
              <!-- color preferences: captive_83_title_color, captive_84_text_color, captive_82_container_bg_color, captive_7_background -->
              <div class="row">
                <div
                  v-for="pref in [captivePrefsSorted[7], captivePrefsSorted[8], captivePrefsSorted[9], captivePrefsSorted[10]]"
                  :key="pref.key"
                  :class="[pref.onError ? 'has-error' : '', 'col-md-6', 'col-lg-3', 'captive-portal-pref-inline', 'mg-bottom-3r']"
                >
                  <label class="block-centered control-label">
                    {{$t('hotspot.'+pref.key)}}
                    <span :class="[getPrefIcon(pref.key)]"></span>
                  </label>
                  <div>
                    <sketch-picker
                      @input="onColorUpdated(pref)"
                      class="absolute-center"
                      v-model="pref.value"
                    />
                  </div>
                </div>
              </div>
              <!-- text style: captive_85_text_style -->
              <div :class="[captivePrefsSorted[11].onError ? 'has-error' : '', 'form-group', 'mg-bottom-3r']">
                <label class="col-sm-4 control-label">
                  {{$t('hotspot.'+captivePrefsSorted[11].key)}}
                  <span :class="[getPrefIcon(captivePrefsSorted[11].key)]"></span>
                </label>
                <div class="col-sm-6">
                  <select
                    v-on:change="textStyleChanged(captivePrefsSorted[11])"
                    v-model="captivePrefsSorted[11].value"
                    class="form-control"
                  >
                    <option v-for="textStyle in textStyles" v-bind:key="textStyle">
                      {{ textStyle }}
                    </option>
                  </select>
                </div>
              </div>
              <!-- captive portal preview -->
              <div class="form-group">
                <label class="block-centered control-label">{{$t('hotspot.captive_preview')}}</label>
                <div class="voucher-main">
                  <span class="label-dateFilter">
                    <input
                      v-model="preferences.currentDevice"
                      type="radio"
                      id="smartphone"
                      name="gender"
                      value="smartphone"
                    />
                    <label for="smartphone">
                      <span class="fa fa-mobile captive-mode"></span>
                    </label>
                  </span>
                  <span class="label-dateFilter">
                    <input
                      v-model="preferences.currentDevice"
                      type="radio"
                      id="tablet"
                      name="gender"
                      value="tablet"
                    />
                    <label for="tablet">
                      <span class="fa fa-tablet captive-mode"></span>
                    </label>
                  </span>
                  <span class="label-dateFilter">
                    <input
                      v-model="preferences.currentDevice"
                      type="radio"
                      id="desktop"
                      name="gender"
                      value="desktop"
                    />
                    <label for="desktop">
                      <span class="fa fa-desktop captive-mode"></span>
                    </label>
                  </span>
                </div>
                <div class>
                  <div v-if="preferences.isLoading" class="captive-preview">
                    <div class="spinner spinner-lg absolute-center"></div>
                  </div>
                  <captive-portal
                    v-if="!preferences.isLoading"
                    :obj="preferences.captive"
                    :class="['captive-preview', preferences.currentDevice]"
                  ></captive-portal>
                </div>
              </div>
            </div>
            <div v-if="!preferences.isLoading" class="card-pf-footer">
              <div class="dropdown card-pf-time-frame-filter">
                <button type="submit" class="btn btn-primary">{{ $t("update") }}</button>
              </div>
              <p>
                <a href="#" class="card-pf-link-with-icon"></a>
              </p>
            </div>
          </form>
        </div>
      </div>
    </div>

    <div v-if="totals.units.count > 0" class="row row-cards-pf">
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{ $t("hotspot.mac_auth") }}
              <div v-if="macAuth.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
          <div v-if="!macAuth.isLoading" class="card-pf-body">
            <vue-good-table
              :perPage="5"
              :paginate="true"
              :columns="columnsMAC"
              :rows="macAuth.data"
              :lineNumbers="false"
              :defaultSortBy="{field: 'name', type: 'asc'}"
              styleClass="table"
              :nextText="tableLangsTexts.nextText"
              :prevText="tableLangsTexts.prevText"
              :rowsPerPageText="tableLangsTexts.rowsPerPageText"
              :globalSearchPlaceholder="tableLangsTexts.globalSearchPlaceholder"
              :ofText="tableLangsTexts.ofText"
            >
              <template slot="table-row" slot-scope="props">
                <td class="fancy">
                  <strong>{{ props.row.name }}</strong>
                </td>
                <td class="fancy">{{ props.row.username }}</td>
                <td class="fancy">
                  <div>
                    <strong>{{ $t('user.kbps_down') }}</strong>
                    : {{ props.row.kbps_down || '-' }}
                  </div>
                  <div>
                    <strong>{{ $t('user.kbps_up') }}</strong>
                    : {{ props.row.kbps_up || '-' }}
                  </div>
                </td>
                <td class="fancy">{{ props.row.created | formatDate}}</td>
                <td>
                  <button
                    v-on:click="deleteMacAuth(props.row.id)"
                    class="btn btn-danger"
                    type="button"
                  >
                    <span class="fa fa-remove"></span>
                  </button>
                </td>
              </template>
            </vue-good-table>
          </div>
          <div v-if="!macAuth.isLoading" class="card-pf-footer">
            <div class="dropdown card-pf-time-frame-filter">
              <button
                data-toggle="modal"
                data-target="#macAuthModal"
                class="btn btn-primary"
                type="button"
              >
                {{
                $t("hotspot.new_mac_auth") }}
              </button>
            </div>
            <p>
              <a href="#" class="card-pf-link-with-icon"></a>
            </p>
          </div>
        </div>
      </div>
    </div>

    <div class="row row-cards-pf">
      <div
        v-if="user.account_type == 'admin' || user.account_type == 'reseller'"
        class="col-xs-12 col-sm-12 col-md-12 col-lg-12"
      >
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
              <dt>{{ $t("hotspot.business_name") }}</dt>
              <dd>{{info.data.business_name || '-'}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("hotspot.business_vat") }}</dt>
              <dd>{{info.data.business_vat || '-'}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("hotspot.business_address") }}</dt>
              <dd>{{info.data.business_address || '-'}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("hotspot.business_email") }}</dt>
              <dd>{{info.data.business_email || '-'}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("hotspot.business_dpo_recap") }}</dt>
              <dd>{{info.data.business_dpo || '-'}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("hotspot.business_dpo_mail_recap") }}</dt>
              <dd>{{info.data.business_dpo_mail || '-'}}</dd>
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
              <a href="#" class="card-pf-link-with-icon"></a>
            </p>
          </div>
        </div>
      </div>
    </div>

    <div class="row row-cards-pf">
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{ $t("hotspot.privacy") }}
              <div v-if="info.privacy.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
          <div v-if="!info.privacy.isLoading" class="card-pf-body">
            <div class="list-details">
              <dt>{{ $t("hotspot.marketing_info") }}</dt>
              <dd>
                <select
                  v-if="(user.account_type == 'admin' || user.account_type == 'reseller') && disclaimers.privacyDisclaimers.length > 1"
                  v-model="disclaimers.currentPrivacyDisclaimerId"
                  class="disclaimer-select"
                >
                  <option v-for="privacyDisclaimer in disclaimers.privacyDisclaimers" :key="privacyDisclaimer.id" :value="privacyDisclaimer.id">
                    {{ privacyDisclaimer.title }}
                  </option>
                </select>
                <textarea readonly id="marketings-info" :value="selectedPrivacyBody"></textarea>
              </dd>
            </div>
            <div class="list-details">
              <dt class="tos-title">{{ $t("hotspot.terms_info") }}</dt>
              <dd>
                <select
                  v-if="(user.account_type == 'admin' || user.account_type == 'reseller') && disclaimers.tosDisclaimers.length > 1"
                  v-model="disclaimers.currentTosDisclaimerId"
                  class="disclaimer-select"
                >
                  <option v-for="tosDisclaimer in disclaimers.tosDisclaimers" :key="tosDisclaimer.id" :value="tosDisclaimer.id">
                    {{ tosDisclaimer.title }}
                  </option>
                </select>
                <textarea readonly :value="selectedTosBody"></textarea>
              </dd>
            </div>
          </div>
          <div v-if="!info.privacy.isLoading" class="card-pf-footer">
            <div class="dropdown card-pf-time-frame-filter">
              <button
                @click="printPrivacy()"
                :class="['btn', user.account_type == 'admin' || user.account_type == 'reseller' ? 'btn-default' : 'btn-primary']"
              >{{$t('hotspot.print_privacy')}}</button>
              <span v-if="disclaimers.updateError" class="pficon pficon-error-circle-o bigger update-pref-feedback"></span>
              <span v-if="disclaimers.updateSuccess" class="pficon pficon-ok update-pref-feedback"></span>
              <button
                v-if="(user.account_type == 'admin' || user.account_type == 'reseller') && (disclaimers.privacyDisclaimers.length > 1 || disclaimers.tosDisclaimers.length > 1)"
                @click="updateDisclaimers()"
                :class="['btn', user.account_type == 'admin' || user.account_type == 'reseller' ? 'btn-primary' : 'btn-default']"
              >{{$t('update')}}</button>
            </div>
            <p>
              <a href="#" class="card-pf-link-with-icon"></a>
            </p>
          </div>
        </div>
      </div>
    </div>

    <div
      class="modal fade"
      id="voucherModal"
      tabindex="-1"
      role="dialog"
      aria-labelledby="myModalLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="myModalLabel">{{$t('hotspot.voucher_creation')}}</h4>
          </div>
          <form class="form-horizontal" v-on:submit.prevent="createVoucher(false)">
            <div class="modal-body">
              <div v-if="vouchers.typesAvailable == 'all'" class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.type')}}</label>
                <div class="col-sm-7">
                  <span class="span-radio">
                    <input
                      required
                      v-model="newVoucher.type"
                      class="form-check-input"
                      type="radio"
                      name="voucherType"
                      id="voucherType1"
                      value="normal"
                    />
                    <span class="voucher-type-span-normal"></span>
                    <label class="form-check-label" for="voucherType1">{{$t('hotspot.normal')}}</label>
                  </span>
                  <span class="span-radio">
                    <input
                      required
                      v-model="newVoucher.type"
                      class="form-check-input"
                      type="radio"
                      name="voucherType"
                      id="voucherType2"
                      value="auth"
                    />
                    <span class="voucher-type-span-auth"></span>
                    <label class="form-check-label" for="voucherType2">{{$t('hotspot.auth')}}</label>
                  </span>
                </div>
              </div>

              <div class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.voucher_creation_count')}}</label>
                <div class="col-sm-7">
                  <input
                    required
                    v-model="vouchersCount"
                    type="text"
                    id="textInput-modal-markup"
                    class="form-control"
                  />
                </div>
              </div>

              <div class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.auto_login')}}</label>
                <div class="col-sm-7">
                  <input
                    v-model="newVoucher.auto_login"
                    type="checkbox"
                    id="textInput-modal-markup"
                    class="form-control"
                  />
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.bandwidth_down')}}</label>
                <div class="col-sm-7">
                  <input
                    v-model="newVoucher.bandwidth_down"
                    type="number"
                    id="textInput-modal-markup"
                    class="form-control"
                  />
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.bandwidth_up')}}</label>
                <div class="col-sm-7">
                  <input
                    v-model="newVoucher.bandwidth_up"
                    type="number"
                    id="textInput-modal-markup"
                    class="form-control"
                  />
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.daily_traffic_limit')}}</label>
                <div class="col-sm-7">
                  <input
                    v-model="newVoucher.max_traffic"
                    type="number"
                    id="textInput-modal-markup"
                    class="form-control"
                  />
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.daily_time_limit')}}</label>
                <div class="col-sm-7">
                  <input
                    v-model="newVoucher.max_time"
                    type="number"
                    id="textInput-modal-markup"
                    class="form-control"
                  />
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.validity')}}</label>
                <div class="col-sm-7">
                  <span class="span-radio">
                    <input
                      required
                      v-model="newVoucher.time"
                      class="form-check-input"
                      type="radio"
                      name="timeExpiration"
                      id="timeExpiration1"
                      value="duration"
                    />
                    <label class="form-check-label" for="timeExpiration1">{{$t('hotspot.duration')}}</label>
                  </span>
                  <span class="span-radio">
                    <input
                      required
                      v-model="newVoucher.time"
                      class="form-check-input"
                      type="radio"
                      name="timeExpiration"
                      id="timeExpiration2"
                      value="expiration"
                    />
                    <label
                      class="form-check-label"
                      for="timeExpiration2"
                    >{{$t('hotspot.expiration')}}</label>
                  </span>
                </div>
              </div>
              <div v-if="newVoucher.time == 'duration'" class="form-group">
                <label class="col-sm-5 control-label" for="textInput-modal-markup"></label>
                <div class="col-sm-7">
                  <input
                    v-model="newVoucher.duration"
                    type="number"
                    id="textInput-modal-markup"
                    class="form-control"
                  />
                </div>
              </div>
              <div v-if="newVoucher.time == 'expiration'" class="form-group">
                <label class="col-sm-5 control-label" for="textInput-modal-markup"></label>
                <div class="col-sm-7">
                  <datepicker
                    :format="dateFormatter"
                    v-model="newVoucher.expiration"
                    :language="locale"
                  ></datepicker>
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.mode')}}</label>
                <div class="col-sm-7">
                  <span class="span-radio">
                    <input
                      required
                      v-model="newVoucher.limitless"
                      class="form-check-input"
                      type="radio"
                      name="exampleRadios"
                      id="exampleRadios1"
                      value="true"
                    />
                    <label class="form-check-label" for="exampleRadios1">{{$t('hotspot.limitless')}}</label>
                  </span>
                  <span class="span-radio">
                    <input
                      required
                      v-model="newVoucher.limitless"
                      class="form-check-input"
                      type="radio"
                      name="exampleRadios"
                      id="exampleRadios2"
                      value="false"
                    />
                    <label
                      class="form-check-label"
                      for="exampleRadios2"
                    >{{$t('hotspot.with_limit')}}</label>
                  </span>
                </div>
              </div>
              <div v-if="newVoucher.limitless == 'false'" class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.max_use')}}</label>
                <div class="col-sm-7">
                  <input
                    v-model="newVoucher.remain_use"
                    type="number"
                    id="textInput-modal-markup"
                    class="form-control"
                  />
                </div>
              </div>
              <div v-if="newVoucher.type == 'auth'" class="divider"></div>
              <div v-if="newVoucher.type == 'auth'" class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.user_name')}}</label>
                <div class="col-sm-7">
                  <input
                    :required="newVoucher.type == 'auth' ? true : false"
                    v-model="newVoucher.user_name"
                    type="text"
                    id="textInput-modal-markup"
                    class="form-control"
                  />
                </div>
              </div>
              <div v-if="newVoucher.type == 'auth'" class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.user_mail')}}</label>
                <div class="col-sm-7">
                  <input
                    v-model="newVoucher.user_mail"
                    type="email"
                    id="textInput-modal-markup"
                    class="form-control"
                  />
                </div>
              </div>
            </div>
            <div class="modal-footer">
              <span
                v-if="vouchers.isCreating"
                class="spinner spinner-sm spinner-inline modal-spinner"
              ></span>
              <button type="button" class="btn btn-default" data-dismiss="modal">{{$t('cancel')}}</button>
              <button type="submit" class="btn btn-primary">{{$t('save')}}</button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <div
      class="modal fade"
      id="voucherModalCustom"
      tabindex="-1"
      role="dialog"
      aria-labelledby="myModalLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="myModalLabel">{{$t('hotspot.voucher_creation')}}</h4>
          </div>
          <form class="form-horizontal" v-on:submit.prevent="createVoucher(true)">
            <div class="modal-body">
              <div v-if="vouchers.typesAvailable == 'all'" class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.type')}}</label>
                <div class="col-sm-7">
                  <span class="span-radio">
                    <input
                      required
                      v-model="newVoucher.type"
                      class="form-check-input"
                      type="radio"
                      name="voucherType"
                      id="voucherType1"
                      value="normal"
                    />
                    <span class="voucher-type-span-normal"></span>
                    <label class="form-check-label" for="voucherType1">{{$t('hotspot.normal')}}</label>
                  </span>
                  <span class="span-radio">
                    <input
                      required
                      v-model="newVoucher.type"
                      class="form-check-input"
                      type="radio"
                      name="voucherType"
                      id="voucherType2"
                      value="auth"
                    />
                    <span class="voucher-type-span-auth"></span>
                    <label class="form-check-label" for="voucherType2">{{$t('hotspot.auth')}}</label>
                  </span>
                </div>
              </div>

              <div class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.voucher_custom_code')}}</label>
                <div class="col-sm-7">
                  <input
                    required
                    v-model="newVoucher.code"
                    type="text"
                    id="textInput-modal-markup"
                    class="form-control"
                    value="1"
                  />
                </div>
              </div>

              <div class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.auto_login')}}</label>
                <div class="col-sm-7">
                  <input
                    v-model="newVoucher.auto_login"
                    type="checkbox"
                    id="textInput-modal-markup"
                    class="form-control"
                  />
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.bandwidth_down')}}</label>
                <div class="col-sm-7">
                  <input
                    v-model="newVoucher.bandwidth_down"
                    type="number"
                    id="textInput-modal-markup"
                    class="form-control"
                  />
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.bandwidth_up')}}</label>
                <div class="col-sm-7">
                  <input
                    v-model="newVoucher.bandwidth_up"
                    type="number"
                    id="textInput-modal-markup"
                    class="form-control"
                  />
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.daily_traffic_limit')}}</label>
                <div class="col-sm-7">
                  <input
                    v-model="newVoucher.max_traffic"
                    type="number"
                    id="textInput-modal-markup"
                    class="form-control"
                  />
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.daily_time_limit')}}</label>
                <div class="col-sm-7">
                  <input
                    v-model="newVoucher.max_time"
                    type="number"
                    id="textInput-modal-markup"
                    class="form-control"
                  />
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.validity')}}</label>
                <div class="col-sm-7">
                  <span class="span-radio">
                    <input
                      required
                      v-model="newVoucher.time"
                      class="form-check-input"
                      type="radio"
                      name="timeExpiration"
                      id="timeExpiration1"
                      value="duration"
                    />
                    <label class="form-check-label" for="timeExpiration1">{{$t('hotspot.duration')}}</label>
                  </span>
                  <span class="span-radio">
                    <input
                      required
                      v-model="newVoucher.time"
                      class="form-check-input"
                      type="radio"
                      name="timeExpiration"
                      id="timeExpiration2"
                      value="expiration"
                    />
                    <label
                      class="form-check-label"
                      for="timeExpiration2"
                    >{{$t('hotspot.expiration')}}</label>
                  </span>
                </div>
              </div>
              <div v-if="newVoucher.time == 'duration'" class="form-group">
                <label class="col-sm-5 control-label" for="textInput-modal-markup"></label>
                <div class="col-sm-7">
                  <input
                    v-model="newVoucher.duration"
                    type="number"
                    id="textInput-modal-markup"
                    class="form-control"
                  />
                </div>
              </div>
              <div v-if="newVoucher.time == 'expiration'" class="form-group">
                <label class="col-sm-5 control-label" for="textInput-modal-markup"></label>
                <div class="col-sm-7">
                  <datepicker
                    :format="dateFormatter"
                    v-model="newVoucher.expiration"
                    :language="locale"
                  ></datepicker>
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.mode')}}</label>
                <div class="col-sm-7">
                  <span class="span-radio">
                    <input
                      required
                      v-model="newVoucher.limitless"
                      class="form-check-input"
                      type="radio"
                      name="exampleRadios"
                      id="exampleRadios1"
                      value="true"
                    />
                    <label class="form-check-label" for="exampleRadios1">{{$t('hotspot.limitless')}}</label>
                  </span>
                  <span class="span-radio">
                    <input
                      required
                      v-model="newVoucher.limitless"
                      class="form-check-input"
                      type="radio"
                      name="exampleRadios"
                      id="exampleRadios2"
                      value="false"
                    />
                    <label
                      class="form-check-label"
                      for="exampleRadios2"
                    >{{$t('hotspot.with_limit')}}</label>
                  </span>
                </div>
              </div>
              <div v-if="newVoucher.limitless == 'false'" class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.max_use')}}</label>
                <div class="col-sm-7">
                  <input
                    v-model="newVoucher.remain_use"
                    type="number"
                    id="textInput-modal-markup"
                    class="form-control"
                  />
                </div>
              </div>
              <div v-if="newVoucher.type == 'auth'" class="divider"></div>
              <div v-if="newVoucher.type == 'auth'" class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.user_name')}}</label>
                <div class="col-sm-7">
                  <input
                    :required="newVoucher.type == 'auth' ? true : false"
                    v-model="newVoucher.user_name"
                    type="text"
                    id="textInput-modal-markup"
                    class="form-control"
                  />
                </div>
              </div>
              <div v-if="newVoucher.type == 'auth'" class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.user_mail')}}</label>
                <div class="col-sm-7">
                  <input
                    v-model="newVoucher.user_mail"
                    type="email"
                    id="textInput-modal-markup"
                    class="form-control"
                  />
                </div>
              </div>
            </div>
            <div class="modal-footer">
              <span
                v-if="vouchers.isCreating"
                class="spinner spinner-sm spinner-inline modal-spinner"
              ></span>
              <button type="button" class="btn btn-default" data-dismiss="modal">{{$t('cancel')}}</button>
              <button type="submit" class="btn btn-primary">{{$t('save')}}</button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <div
      class="modal fade"
      id="voucherDeleteAll"
      tabindex="-1"
      role="dialog"
      aria-labelledby="myModalLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4
              class="modal-title"
              id="HSdeleteVoucherModalLabel"
            >{{ $t("hotspot.delete_all_vouchers") }}</h4>
          </div>
          <form class="form-horizontal" role="form" v-on:submit.prevent="deleteAllVouchers()">
            <div class="modal-body">
              <div class="alert alert-warning alert-dismissable">
                <span class="pficon pficon-warning-triangle-o"></span>
                <strong>{{ $t("hotspot.warning_delete_title") }}</strong>
                . {{ $t("hotspot.warning_delete_vouchers") }}.
              </div>
            </div>
            <div class="modal-footer">
              <span
                v-if="vouchers.isDeleting"
                class="spinner spinner-sm spinner-inline modal-spinner"
              ></span>
              <button type="button" class="btn btn-default" data-dismiss="modal">{{ $t("cancel") }}</button>
              <button type="submit" class="btn btn-danger">{{ $t("delete") }}</button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <div
      class="modal fade"
      id="macAuthModal"
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
            <h4 class="modal-title">{{ $t("hotspot.new_mac_auth") }}</h4>
          </div>
          <form class="form-horizontal" role="form" v-on:submit.prevent="createMacAuth()">
            <div class="modal-body">
              <div class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="ACtextInput2-modal-markup"
                >{{ $t("unit.unit") }}</label>
                <div class="col-sm-7">
                  <select v-model="newMACAuth.unit" class="form-control">
                    <option v-for="u in totals.units.data" v-bind:key="u.id" v-bind:value="u">
                      {{ u.name }} -
                      {{ u.description }}
                    </option>
                  </select>
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{ $t("unit.description") }}</label>
                <div class="col-sm-7">
                  <input
                    required
                    v-model="newMACAuth.name"
                    type="text"
                    id="textInput-modal-markup"
                    class="form-control"
                    :placeholder="$t('unit.description')"
                  />
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput2-modal-markup"
                >{{ $t("unit.mac_address") }}</label>
                <div class="col-sm-7">
                  <input
                    required
                    v-model="newMACAuth.username"
                    pattern="^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})"
                    type="text"
                    id="textInput2-modal-markup"
                    class="form-control"
                    placeholder="00:11:22:AA:BB:CC"
                  />
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.bandwidth_down')}}</label>
                <div class="col-sm-7">
                  <input v-model="newMACAuth.kbps_down" type="number" class="form-control" />
                </div>
              </div>
              <div class="form-group">
                <label
                  class="col-sm-5 control-label"
                  for="textInput-modal-markup"
                >{{$t('hotspot.bandwidth_up')}}</label>
                <div class="col-sm-7">
                  <input v-model="newMACAuth.kbps_up" type="number" class="form-control" />
                </div>
              </div>
            </div>
            <div class="modal-footer">
              <span
                v-show="newMACAuth.onAction"
                class="spinner spinner-sm spinner-inline modal-spinner"
              ></span>
              <button type="button" class="btn btn-default" data-dismiss="modal">{{ $t("cancel") }}</button>
              <button type="submit" class="btn btn-primary">{{ $t("create") }}</button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <div
      v-for="voucher in vouchers.usable"
      v-bind:key="voucher.id"
      class="card-pf"
      id="voucher-coupon"
    >
      <img class="voucher-logo" src="/static/logo.png" />
      <h2 class="card-pf-title voucher-name" id="test">{{info.data.name}}</h2>
      <h3 class="card-pf-title voucher-desc">{{info.data.description}}</h3>
      <div class="card-pf-body voucher-main">
        <p>
          <strong class="voucher-code">{{voucher.code | uppercase}}</strong>
        </p>
      </div>
      <div class="card-pf-footer voucher-details">
        <div class="card-pf-time-frame-filter voucher-valid">
          {{ $t("hotspot.valid") }}:
          <strong v-if="voucher.duration > 0">{{voucher.duration}}</strong>
          <span v-if="voucher.duration > 0">{{ $t("hotspot.days") }}</span>
          <strong v-if="voucher.duration == 0">{{voucher.expires | formatDate(true)}}</strong>
        </div>
        <div class="card-pf-time-frame-filter voucher-max-use">
          <b v-if="voucher.type == 'auth'">{{$t('hotspot.'+voucher.type+'_print')}}</b>
          <span v-if="voucher.type == 'auth'">-</span>
          {{voucher.remain_use == -1 ? $t('hotspot.limitless') : $t('hotspot.max_use') + ': ' }}
          <strong
            v-if="voucher.remain_use != -1"
          >{{voucher.remain_use}}</strong>
        </div>
        <div class="card-pf-time-frame-filter voucher-traffic">
          {{ $t("hotspot.traffic") }}:
          <strong>{{voucher.max_traffic | byteFormat}}</strong>
        </div>
        <div class="card-pf-time-frame-filter voucher-time">
          {{ $t("hotspot.time") }}:
          <strong>{{voucher.max_time | secondsInHour}}</strong>
        </div>

        <div class="voucher-download">
          <span class="fa fa-arrow-down"></span> Download:
          <span v-if="voucher.bandwidth_down === 0">{{ $t("hotspot.unlimited") }}</span>
          <span v-else>
            <strong>{{voucher.bandwidth_down}}</strong> Kb/s
          </span>
        </div>
        <div class="voucher-upload">
          <span class="fa fa-arrow-up"></span> Upload:
          <span v-if="voucher.bandwidth_up === 0">{{ $t("hotspot.unlimited") }}</span>
          <span v-else>
            <strong>{{voucher.bandwidth_up}}</strong> Kb/s
          </span>
        </div>
      </div>
    </div>
    <div id="voucherPrint" v-show="showVoucherPrint" class="voucherPrintContainer">
      <div
        v-for="index in Math.ceil(vouchersToPrint.length / 2)"
        :key="index"
        :class="[ index % 5 == 0 && index != 0 ? 'pagebreak' : '']"
      >
        <!-- two vouchers columns -->
        <div v-for="i in [ (index*2-2), (index*2-1) ]" :key="i" class="voucher-print">
          <div v-if="i < vouchersToPrint.length">
            <table class="width-100">
              <tr>
                <td class="voucher-print-q1"></td>
                <td class="voucher-print-q2">
                  <img :src="preferences.captive[2].value" class="voucher-print-logo" />
                  <div class="voucher-print-hs-title" :style="{ 'font-size': voucherTitleSize + 'px' }">{{ info.data.description }}</div>
                  <table class="width-75">
                    <tr>
                      <td class="padding-left-2 border-left">
                        <div class="voucher-print-code">{{ vouchersToPrint[i].code | uppercase }}</div>
                      </td>
                    </tr>
                    <tr>
                      <td class="padding-left-2 border-left line-height-5">
                        <div class="voucher-print-data">
                          <span class="fa fa-clock-o"></span>
                          <span v-if="vouchersToPrint[i].duration > 0">
                            {{ vouchersToPrint[i].duration }}
                            {{ $t("hotspot.days") }}
                          </span>
                          <span v-else>{{ vouchersToPrint[i].expires | formatDate(true) }}</span>
                        </div>
                        <div class="voucher-print-data">
                          <span v-show="vouchersToPrint[i].user_name" class="pficon pficon-user"></span>
                          {{ vouchersToPrint[i].user_name | truncate(12) }}
                        </div>
                        <div class="voucher-print-data">
                          <span class="pficon pficon-history"></span>
                          <span v-if="vouchersToPrint[i].remain_use != -1">
                            {{ vouchersToPrint[i].remain_use }}
                            {{ $t("hotspot.max_use") }}
                          </span>
                          <span v-else>{{ $t("hotspot.reusable") }}</span>
                        </div>
                        <div class="voucher-print-data"></div>
                      </td>
                    </tr>
                  </table>
                </td>
              </tr>
              <tr>
                <td></td>
                <td class="voucher-print-q4">
                  <table class="width-75">
                    <tr>
                      <td class="padding-left-2 border-left">
                        <div class="voucher-print-data">
                          <span class="fa fa-arrow-circle-o-down"></span>
                          <span v-if="vouchersToPrint[i].bandwidth_down != 0">
                            {{ vouchersToPrint[i].bandwidth_down }}
                            Kbps
                          </span>
                          <span v-else>{{ $t("hotspot.unlimited") }}</span>
                        </div>
                        <div class="voucher-print-data">
                          <span class="fa fa-tachometer"></span>
                          <span
                            v-if="vouchersToPrint[i].max_traffic != 0"
                          >{{ vouchersToPrint[i].max_traffic | byteFormat}}</span>
                          <span v-else>{{ $t("hotspot.unlimited") }}</span>
                        </div>
                      </td>
                    </tr>
                    <tr>
                      <td class="padding-left-2-2">
                        <div class="voucher-print-data align-top">
                          <span class="fa fa-arrow-circle-o-up"></span>
                          <span v-if="vouchersToPrint[i].bandwidth_up != 0">
                            {{ vouchersToPrint[i].bandwidth_up }}
                            Kbps
                          </span>
                          <span v-else>{{ $t("hotspot.unlimited") }}</span>
                        </div>
                        <div class="voucher-print-data align-top">
                          <span class="pficon pficon-pending"></span>
                          <span
                            v-if="vouchersToPrint[i].max_time != 0"
                          >{{ vouchersToPrint[i].max_time | secondsInHour }}</span>
                          <span v-else>{{ $t("hotspot.unlimited") }}</span>
                        </div>
                      </td>
                    </tr>
                  </table>
                </td>
              </tr>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import HotspotService from "../../services/hotspot";
import PreferenceService from "../../services/preference";
import AccountService from "../../services/account";
import UnitService from "../../services/unit";
import UserService from "../../services/user";
import DeviceService from "../../services/device";
import SessionService from "../../services/session";
import StatsService from "../../services/stats";
import StorageService from "../../services/storage";
import UtilService from "../../services/util";
import DisclaimerService from "../../services/disclaimer";

import jsPDF from "jspdf";
import CaptivePortal from "../../directives/CaptivePortal.vue";
import HotspotAction from "../../directives/HotspotAction.vue";
import PictureInput from "vue-picture-input";
import arrow from "../../../static/arrows.js";
import { Sketch } from "vue-color";
import { VueEditor } from "vue2-editor";
import { setTimeout } from "timers";
import Datepicker from "vuejs-datepicker";
import moment from "moment";

export default {
  name: "HotspotDetails",
  mixins: [
    HotspotService,
    PreferenceService,
    AccountService,
    UnitService,
    UserService,
    DeviceService,
    SessionService,
    StatsService,
    StorageService,
    UtilService,
    DisclaimerService
  ],
  components: {
    hotspotAction: HotspotAction,
    PictureInput,
    "sketch-picker": Sketch,
    captivePortal: CaptivePortal,
    VueEditor,
    Datepicker
  },
  mounted() {
    // get hotspot info
    this.getInfo();

    // get totals
    this.getTotals();

    // get preferences
    this.getPreferences();

    // get vouchers
    this.getVouchers();

    // get mac auths
    this.getAllMACAuth();

    // get sms count
    this.getSmsCount();

    // set selected hotspot in local storage
    this.set("selected_hotspot_id", parseInt(this.$route.params.id) || this.get("selected_hotspot_id") || 0);
  },
  data() {
    return {
      info: {
        isLoading: true,
        data: {},
        privacy: {
          isLoading: true,
          data: {}
        }
      },
      vouchers: {
        isLoading: true,
        isDeleting: false,
        isCreating: false,
        data: [],
        usable: [],
        filters: {
          code: "",
          duration: 0,
          expiredStart: "",
          expiredEnd: "",
          createdStart: "",
          createdEnd: "",
          auto_login: "",
          used: "",
          reusable: "",
          printed: "",
          type: ""
        },
        typesAvailable: "normal"
      },
      macAuth: {
        isLoading: true,
        data: []
      },
      vouchersCount: 1,
      newVoucher: {
        code: "",
        auto_login: true,
        bandwidth_down: 0,
        bandwidth_up: 0,
        max_traffic: 0,
        max_time: 0,
        duration: 7,
        remain_use: 3,
        limitless: "true",
        type: "normal",
        user_name: "",
        user_mail: "",
        time: "duration",
        expiration: moment(Date.now() + 12096e5)
          .utc()
          .startOf("day")
          .toISOString()
      },
      newMACAuth: {
        name: "",
        username: "",
        kbps_down: 0,
        kbps_up: 0,
        onAction: false
      },
      preferences: {
        isLoading: true,
        currentDevice: "smartphone",
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
        },
        sms: {
          isLoading: true,
          count: 0
        }
      },
      columns: [
        {
          label: this.$i18n.t("hotspot.code"),
          field: "code",
          filterable: false,
          sortable: false
        },
        {
          label: this.$i18n.t("hotspot.type"),
          field: "type",
          filterable: false,
          sortable: false
        },
        {
          label: this.$i18n.t("hotspot.auto_login"),
          field: "auto_login",
          filterable: false,
          sortable: false
        },
        {
          label: this.$i18n.t("hotspot.bandwidth_limit"),
          field: "bandwidth",
          filterable: false,
          sortable: false
        },
        {
          label: this.$i18n.t("hotspot.daily_limit"),
          field: "limit",
          filterable: false,
          sortable: false
        },
        {
          label: this.$i18n.t("hotspot.created"),
          field: "created",
          filterable: false,
          sortable: true
        },
        {
          label: this.$i18n.t("hotspot.duration"),
          field: "duration",
          filterable: false,
          sortable: false
        },
        {
          label: this.$i18n.t("hotspot.expiration"),
          field: "expires",
          filterable: false,
          sortable: true
        },
        {
          label: this.$i18n.t("hotspot.mode"),
          field: "remain_use",
          filterable: false,
          sortable: false
        },
        {
          label: "",
          field: "",
          sortable: false
        }
      ],
      columnsMAC: [
        {
          label: this.$i18n.t("unit.description"),
          field: "name",
          filterable: false,
          sortable: false
        },
        {
          label: this.$i18n.t("unit.mac_address"),
          field: "username",
          filterable: false
        },
        {
          label: this.$i18n.t("hotspot.bandwidth_limit"),
          field: "bandwidth",
          filterable: false,
          sortable: false
        },
        {
          label: this.$i18n.t("hotspot.created"),
          field: "created",
          filterable: false,
          sortable: false
        },
        {
          label: "",
          field: "",
          sortable: false
        }
      ],
      tableLangsTexts: this.tableLangs(),
      uploadLangstexts: this.uploadImageLangs(),
      user: this.get("loggedUser"),
      displayCaptivePortalOptions:
        (this.get("loggedUser") &&
          this.get("loggedUser").subscription &&
          this.get("loggedUser").subscription.subscription_plan &&
          this.get("loggedUser").subscription.subscription_plan
            .wings_customization) ||
        this.get("loggedUser").account_type == "admin",
      customToolbar: [
        ["bold", "italic", "underline"],
        ["image", "code-block"]
      ],
      userLink: {
        name: "Users",
        params: {
          hotspotId: this.$route.params.id
        }
      },
      unitLink: {
        name: "Units",
        params: {
          hotspotId: this.$route.params.id
        }
      },
      sessionLink: {
        name: "Sessions",
        params: {
          hotspotId: this.$route.params.id
        }
      },
      accountLink: {
        name: "Accounts",
        params: {
          hotspotId: this.$route.params.id
        }
      },
      deviceLink: {
        name: "Devices",
        params: {
          hotspotId: this.$route.params.id
        }
      },
      locale: this.$root.$options.currentLocale,
      smsMaxCount: 0,
      smsMaxCountAdd: 0,
      smsThreshold: 0,
      showVoucherPrint: false,
      vouchersToPrint: [],
      advancedFilters: false,
      textStyles: [
        "Lato",
        "Roboto",
        "Hind",
        "Fira Sans Extra Condensed",
        "Dosis",
        "EB Garamond",
        "Playfair Display",
        "Playfair Display SC",
      ],
      captivePrefsSortMap: {
        "captive_1_redir": 0,
        "captive_2_title": 1,
        "captive_3_logo": 2,
        "captive_4_subtitle": 3,
        "captive_5_banner": 4,
        "captive_6_description": 5,
        "captive_7_background": 10,
        "captive_81_bg_image": 6,
        "captive_82_container_bg_color": 9,
        "captive_83_title_color": 7,
        "captive_84_text_color": 8,
        "captive_85_text_style": 11
      },
      captivePrefsSorted: [],
      disclaimers: {
        currentPrivacyDisclaimerId: 0,
        currentTosDisclaimerId: 0,
        privacyDisclaimers: [],
        tosDisclaimers: [],
        updateSuccess: false,
        updateError: false,
      },
    };
  },
  computed: {
    voucherTitleSize: function () {
      var size = 450 / this.info.data.description.length;

      if (size > 20) {
        size = 20;
      }
      return size;
    },
    selectedPrivacyBody: function() {
      let selectedPrivacyDisclaimer = this.disclaimers.privacyDisclaimers.find((element) => {
        return element.id == this.disclaimers.currentPrivacyDisclaimerId;
      });
      return selectedPrivacyDisclaimer.body;
    },
    selectedTosBody: function() {
      let selectedTosDisclaimer = this.disclaimers.tosDisclaimers.find((element) => {
        return element.id == this.disclaimers.currentTosDisclaimerId;
      });
      return selectedTosDisclaimer.body;
    },
  },
  methods: {
    dateFormatter(date) {
      return moment(date).format("DD MMMM YYYY");
    },
    getPrevieHTML(value) {
      return '<img src="' + value + '"></img>';
    },
    getPrefIcon(pref) {
      return this.getPrefTypeIcon(pref);
    },
    addSMSCount(pref, amount) {
      pref.value = (parseInt(pref.value) + parseInt(amount)).toString();
      this.updatePreferences();
    },
    getSmsCount() {
      this.statsSMSSentByHotspot(
        this.$route.params.id,
        success => {
          this.totals.sms.count = success.body.length;
          this.totals.sms.isLoading = false;
        },
        error => {
          console.error(error.body);
          this.totals.sms.data = 0;
          this.totals.sms.isLoading = false;
        }
      );
    },
    createVoucher(custom) {
      this.vouchers.isCreating = true;

      var promises = [];
      var context = this;
      for (var i = 0; i < this.vouchersCount; i++) {
        promises.push(
          new Promise(function(resolve, reject) {
            context.hotspotCreateVoucher(
              {
                hotspot_id: parseInt(context.$route.params.id),
                code: custom
                  ? context.newVoucher.code
                  : context.generateVoucher(),
                auto_login: context.newVoucher.auto_login,
                bandwidth_down: parseInt(context.newVoucher.bandwidth_down),
                bandwidth_up: parseInt(context.newVoucher.bandwidth_up),
                max_traffic:
                  parseInt(context.newVoucher.max_traffic) * 1024 * 1024,
                max_time: parseInt(context.newVoucher.max_time) * 60,
                time: context.newVoucher.time,
                duration: parseInt(context.newVoucher.duration),
                expiration: moment(moment(context.newVoucher.expiration)).diff(
                  moment(new Date()),
                  "days"
                ),
                type: context.newVoucher.type,
                user_name:
                  context.newVoucher.type == "auth"
                    ? context.newVoucher.user_name + (i > 0 ? "-" + i : "")
                    : null,
                user_mail:
                  context.newVoucher.type == "auth"
                    ? context.newVoucher.user_mail
                    : null,
                remain_use:
                  context.newVoucher.limitless == "true"
                    ? -1
                    : parseInt(context.newVoucher.remain_use)
              },
              success => {
                resolve();
              },
              error => {
                reject();
              }
            );
          })
        );
      }
      Promise.all(promises)
        .then(function() {
          context.vouchers.isCreating = false;
          context.newVoucher.user_name = "";
          context.newVoucher.user_mail = "";
          context.getVouchers();
          $("#voucherModal").modal("hide");
          $("#voucherModalCustom").modal("hide");
        })
        .catch(function(err) {
          console.error(err);
          context.vouchers.isCreating = false;
        });
    },
    deleteVoucher(id) {
      this.vouchers.isLoading = true;
      this.hotspotVoucherDelete(
        id,
        success => {
          this.vouchers.isLoading = false;
          this.getVouchers();
        },
        error => {
          console.error(error.body);
          this.vouchers.isLoading = false;
          this.getVouchers();
        }
      );
    },
    deleteAllVouchers() {
      var context = this;
      var promises = [];
      this.vouchers.isDeleting = true;

      for (var v in this.vouchers.data) {
        var voucher = this.vouchers.data[v];
        promises.push(
          new Promise((resolve, reject) => {
            context.hotspotVoucherDelete(
              voucher.id,
              success => {
                resolve(success);
              },
              error => {
                reject(error);
              }
            );
          })
        );
      }

      Promise.all(promises)
        .then(function() {
          context.vouchers.isDeleting = false;
          context.getVouchers();
          $("#voucherDeleteAll").modal("toggle");
        })
        .catch(function(err) {
          console.error(err);
          context.vouchers.isDeleting = false;
        });
    },
    getVouchers() {
      this.vouchers.data = [];
      this.vouchers.usable = [];
      this.hotspotGetVouchers(
        {
          code: this.vouchers.filters.code,
          duration: this.vouchers.filters.duration,
          auto_login: this.vouchers.filters.auto_login,
          used: this.vouchers.filters.used,
          reusable: this.vouchers.filters.reusable,
          printed: this.vouchers.filters.printed,
          type: this.vouchers.filters.type,
          expiredStart: this.vouchers.filters.expiredStart,
          expiredEnd: this.vouchers.filters.expiredEnd,
          createdStart: this.vouchers.filters.createdStart,
          createdEnd: this.vouchers.filters.createdEnd
        },
        this.$route.params.id,
        success => {
          this.vouchers.data = success.body;
          for (var r in this.vouchers.data) {
            if (this.vouchers.data[r].remain_use != 0) {
              this.vouchers.usable.push(this.vouchers.data[r]);
            }
          }
          this.vouchers.isLoading = false;
        },
        error => {
          console.error(error.body);
          this.vouchers.data = [];
          this.vouchers.isLoading = false;
        }
      );
    },
    checkVoucherUse(value) {
      return +new Date(value) > 0;
    },
    getInfo() {
      this.hotspotGet(
        this.$route.params.id,
        success => {
          this.info.data = success.body;
          this.info.isLoading = false;

          // get privacy info
          this.getPrivacyInfo(this.info.data.uuid);
        },
        error => {
          console.error(error.body);
        }
      );
    },
    getTotals() {
      this.accountGetAll(
        this.$route.params.id,
        null,
        null,
        null,
        success => {
          this.totals.accounts.count = success.body.total;
          this.totals.accounts.isLoading = false;
        },
        error => {
          console.error(error.body);
          this.totals.accounts.isLoading = false;
        }
      );
      this.unitGetAll(
        this.$route.params.id,
        null,
        null,
        null,
        success => {
          this.totals.units.data = success.body.data;
          this.totals.units.count = success.body.total;
          this.totals.units.isLoading = false;
        },
        error => {
          console.error(error.body);
          this.totals.units.isLoading = false;
        }
      );
      this.userGetAll(
        this.$route.params.id,
        null,
        false,
        null,
        null,
        null,
        success => {
          this.totals.users.count = success.body.total;
          this.totals.users.isLoading = false;
        },
        error => {
          console.error(error.body);
          this.totals.users.isLoading = false;
        }
      );
      this.deviceGetAll(
        this.$route.params.id,
        "",
        null,
        null,
        null,
        success => {
          this.totals.devices.count = success.body.total;
          this.totals.devices.isLoading = false;
        },
        error => {
          console.error(error.body);
          this.totals.devices.isLoading = false;
        }
      );
      this.sessionGetAll(
        this.$route.params.id,
        null,
        null,
        null,
        null,
        null,
        null,
        null,
        success => {
          this.totals.sessions.count = success.body.total;
          this.totals.sessions.isLoading = false;
        },
        error => {
          console.error(error.body);
          this.totals.sessions.isLoading = false;
        }
      );
    },
    getPreferences() {
      this.hsPrefGet(
        this.$route.params.id,
        success => {
          var globalPref = [];
          var captivePref = [];
          var vouchersAvailable = false;
          var voucherFlag = 0;

          for (var p in success.body) {
            var pref = success.body[p];
            if (pref.value === "true") {
              pref.value = true;
            }
            if (pref.value === "false") {
              pref.value = false;
            }

            if (
              (pref.key == "voucher_login" && pref.value) ||
              (pref.key == "temp_code_login" && pref.value)
            ) {
              vouchersAvailable = pref.value || vouchersAvailable;

              if (pref.key == "voucher_login" && pref.value) {
                this.vouchers.typesAvailable = "normal";
                this.newVoucher.type = "normal";
                voucherFlag++;
              }

              if (pref.key == "temp_code_login" && pref.value) {
                this.vouchers.typesAvailable = "auth";
                this.newVoucher.type = "auth";
                voucherFlag++;
              }

              if (voucherFlag == 2) {
                this.vouchers.typesAvailable = "all";
                this.newVoucher.type = "normal";
              }
            }

            if (pref.key == "captive_7_background") {
              this.preferences.backgroundColor = pref.value;
            }

            if (pref.key == "captive_81_bg_image") {
              this.preferences.backgroundImage = pref.value;
            }

            if (pref.key == "captive_82_container_bg_color") {
              this.preferences.containerBgColor = pref.value;
            }

            if (pref.key == "captive_83_title_color") {
              this.preferences.titleColor = pref.value;
            }

            if (pref.key == "captive_84_text_color") {
              this.preferences.textColor = pref.value;
            }

            if (pref.key == "captive_85_text_style") {
              this.preferences.textStyle = pref.value;
            }

            if (pref.key == "sms_login_max") {
              this.smsMaxCount = pref.value;
            }

            if (pref.key == "sms_login_threshold") {
              this.smsThreshold = pref.value;
            }

            if (pref.key.startsWith("captive")) {
              captivePref.push(pref);
            } else {
              globalPref.push(pref);
            }

            if (pref.key == "CoovaChilli-Max-Total-Octets") {
              pref.value = parseInt(pref.value) / 1024 / 1024;
            }

            if (pref.key == "CoovaChilli-Max-Navigation-Time") {
              pref.value = parseInt(pref.value) / 60;
            }
          }

          this.preferences.vouchersAvailable = vouchersAvailable;
          this.preferences.global = globalPref;
          this.preferences.captive = captivePref;
          this.captivePrefsSorted = this.sortCaptivePreferences(captivePref),
          this.preferences.isLoading = false;
          var context = this;
          setTimeout(function() {
            context.updateCaptivePreview();
          }, 0);
        },
        error => {
          console.error(error.body);
        }
      );
    },
    sortCaptivePreferences(captivePrefs) {
      var sortedPrefs = [];

      for (var pref of captivePrefs) {
        sortedPrefs[this.captivePrefsSortMap[pref.key]] = pref;
      }
      return sortedPrefs;
    },
    updateCaptivePreview() {
      window.$("#captive-preview").css("background-color", this.preferences.backgroundColor);
      window.$("#captive-preview div.ui.segment").css("background-color", this.preferences.containerBgColor);
      window.$("#captive-preview h2").css("color", this.preferences.titleColor);
      window.$("#captive-preview h3").css("color", this.preferences.textColor);
      window.$("#captive-preview p").css("color", this.preferences.textColor);
      window.$("#captive-preview .terms-space").css("color", this.preferences.textColor);
      window.$("#captive-preview a.terms").css("color", this.preferences.textColor);

      if (this.preferences.backgroundImage) {
        window.$("#captive-preview").css("background-image", 'url("' + this.preferences.backgroundImage + '")');
        window.$("#captive-preview").css("background-size", "cover");
        window.$("#captive-preview").css("background-position", "center");
      } else {
        window.$("#captive-preview").css("background-image", "none");
      }

      // font syle
      window.$("#captive-preview h2").css("font-family", this.preferences.textStyle);
      window.$("#captive-preview h3").css("font-family", this.preferences.textStyle);
      window.$("#captive-preview p").css("font-family", this.preferences.textStyle);
      window.$("#captive-preview a.green.button").css("font-family", this.preferences.textStyle);
      window.$("#captive-preview .terms-space").css("font-family", this.preferences.textStyle);

      // set patterned background to captive images
      setTimeout(function() {
        $('div.preview-container canvas').css('background', 'repeating-linear-gradient(-45deg, LightGray, LightGray 5px, Gray 5px, Gray 10px)');
      }, 500);
    },
    updatePreferences() {
      this.preferences.isLoading = true;
      // create promises array
      var promises = [];
      for (var i in this.preferences.global) {
        var pref = this.preferences.global[i];
        promises.push(
          new Promise((resolve, reject) => {
            if (typeof pref.value == "boolean") {
              pref.value = pref.value.toString();
            }

            if (pref.key == "CoovaChilli-Max-Total-Octets") {
              pref.value = (pref.value * 1024 * 1024).toString();
            }

            if (pref.key == "CoovaChilli-Max-Navigation-Time") {
              pref.value = (pref.value * 60).toString();
            }

            this.hsPrefModify(
              this.$route.params.id,
              pref,
              success => {
                resolve(success);
              },
              error => {
                reject(error);
              }
            );
          })
        );
      }

      // exec promises
      var context = this;
      Promise.all(promises).then(function(response) {
        context.preferences.isLoading = false;
        context.getPreferences();
      });
    },
    getHexValue(value) {
      if (typeof value === 'string' || value instanceof String) {
        return value;
      } else {
        var rgb_hex = value.hex;
        var a_hex = Math.floor(value.a * 255).toString(16).padStart(2, '0');
        return rgb_hex + a_hex;
      }
    },
    updatePreferencesCaptive() {
      this.preferences.isLoading = true;
      // create promises array
      var promises = [];
      for (var i in this.preferences.captive) {
        var pref = this.preferences.captive[i];
        promises.push(
          new Promise((resolve, reject) => {
            if (typeof pref.value == "boolean") {
              pref.value = pref.value.toString();
            }
            if (pref.key == "captive_7_background" || pref.key == "captive_82_container_bg_color" || pref.key == "captive_83_title_color" || pref.key == "captive_84_text_color") {
              pref.value = this.getHexValue(pref.value);
            }
            this.hsPrefModify(
              this.$route.params.id,
              pref,
              success => {
                resolve(success);
              },
              error => {
                reject(error);
              }
            );
          })
        );
      }

      // exec promises
      var context = this;
      Promise.all(promises).then(function(response) {
        context.preferences.isLoading = false;
        context.getPreferences();
      });
    },
    printVouchers(voucherList) {
      var context = this;
      var promises = [];

      for (var index = 0; index < voucherList.length; index++) {
        promises.push(
          new Promise(function(resolve, reject) {
            context.hotspotUpdateVoucher(
              voucherList[index].id,
              {
                printed: true
              },
              success => {
                resolve(success);
              },
              error => {
                reject(error);
              }
            );
          })
        );

        Promise.all(promises)
          .then(function() {
          })
          .catch(function(err) {
            console.error(err);
          });
      }
      this.vouchersToPrint = voucherList;
      this.showVoucherPrint = true;

      setTimeout(async function() {
        context.printDiv("voucherPrint");
      }, 500);
    },
    printDiv(div_id) {
      $("body").html($("#" + div_id).html());
      scroll(0, 0);
      var context = this;

      window.onafterprint = function(event) {
        // show "Go back to Hotspot Manager" button in Firefox
        var browser = navigator.userAgent.toLowerCase();
        if (browser.indexOf("firefox") > -1) {
          $("body").prepend(
            '<div class="go-back-to-hotspot-manager"><button id="goBack" class="btn btn-primary btn-lg" type="button"><span class="fa fa-arrow-left"></span></button></div>'
          );

          $("#goBack").click(function() {
            location.reload();
          });

          setTimeout(function() {
            location.reload();
          }, 2000);
        } else {
          location.reload();
        }
      };

      window.print();
    },
    printVoucher(voucher) {
      for (let i = 0; i < this.vouchers.usable.length; i++) {
        // print vocuher only if it's usable
        if (this.vouchers.usable[i].id === voucher.id) {
          this.printVouchers([voucher]);
          break;
        }
      }
    },
    printAllVoucher() {
      this.printVouchers(this.vouchers.usable);
    },
    exportCSVVoucher() {
      var voucherRows = JSON.parse(JSON.stringify(this.vouchers.usable));
      for (var r in voucherRows) {
        voucherRows[r].type = this.$i18n.t("hotspot." + voucherRows[r].type);

        var banDown =
          voucherRows[r].bandwidth_down > 0
            ? this.$options.filters["byteFormat"](
                voucherRows[r].bandwidth_down * 1024
              )
            : "-";

        var banUp =
          voucherRows[r].bandwidth_up > 0
            ? this.$options.filters["byteFormat"](
                voucherRows[r].bandwidth_up * 1024
              )
            : "-";

        voucherRows[r].bandwidth = "Down: " + banDown + " | Up: " + banUp;

        var traffic =
          voucherRows[r].max_traffic > 0
            ? this.$options.filters["byteFormat"](voucherRows[r].max_traffic)
            : "-";

        var time =
          voucherRows[r].max_time > 0
            ? this.$options.filters["secondsInHour"](voucherRows[r].max_time)
            : "-";

        voucherRows[r].limit = "Traffic: " + traffic + " | Time: " + time;

        voucherRows[r].auto_login = voucherRows[r].auto_login
          ? this.$i18n.t("hotspot.enabled")
          : this.$i18n.t("hotspot.disabled");

        voucherRows[r].duration =
          voucherRows[r].duration + " (" + this.$i18n.t("hotspot.days") + ")";

        voucherRows[r].expires = this.checkVoucherUse(voucherRows[r].expires)
          ? this.$i18n.t("hotspot.expires") +
            ": " +
            this.$options.filters["formatDate"](voucherRows[r].expires)
          : "-";

        voucherRows[r].remain_use =
          voucherRows[r].remain_use == -1
            ? this.$i18n.t("hotspot.limitless")
            : this.$i18n.t("hotspot.max_use") +
              ": " +
              voucherRows[r].remain_use;
      }
      var csv = this.createCSV(this.columns, voucherRows);
      this.downloadCSV(csv.cols, csv.rows, "vouchers");
    },
    onPicturePrefill(pref) {
      // set patterned background to captive images
      setTimeout(function() {
        $('div.preview-container canvas').css('background', 'repeating-linear-gradient(-45deg, LightGray, LightGray 5px, Gray 5px, Gray 10px)');
      }, 20);
    },
    onPictureRemoved(pref) {
      if (pref.key == 'captive_81_bg_image') {
        pref.value = "";
        this.preferences.backgroundImage = pref.value;
      } else if (pref.key == 'captive_3_logo') {
        pref.value = "";
      } else if (pref.key == 'captive_5_banner') {
        pref.value = "";
      }
      this.updateCaptivePreview();
      this.$forceUpdate();
    },
    textStyleChanged(pref) {
      this.preferences.textStyle = pref.value;
      this.updateCaptivePreview();
      this.$forceUpdate();
    },
    onPictureChanged(pref) {
      if (this.$refs["prefInput-" + pref.key][0].image.length > 921083) {
        pref.onError = true;
        this.$refs["prefInput-" + pref.key][0].image = pref.value;
      } else {
        pref.onError = false;
        pref.value = this.$refs["prefInput-" + pref.key][0].image;

        if (pref.key == 'captive_81_bg_image') {
          this.preferences.backgroundImage = pref.value;
        }

        setTimeout(function() {
          $('div#' + pref.key + ' canvas').css('background', 'repeating-linear-gradient(-45deg, LightGray, LightGray 5px, Gray 5px, Gray 10px)');
        }, 200);
      }
      this.updateCaptivePreview();
      this.$forceUpdate();
    },
    onColorUpdated(pref) {
      if (pref.key == 'captive_7_background') {
        this.preferences.backgroundColor = this.getHexValue(pref.value);
      } else if (pref.key == 'captive_82_container_bg_color') {
        this.preferences.containerBgColor = this.getHexValue(pref.value);
      } else if (pref.key == 'captive_83_title_color') {
        this.preferences.titleColor = this.getHexValue(pref.value);
      } else if (pref.key == 'captive_84_text_color') {
        this.preferences.textColor = this.getHexValue(pref.value);
      }
      this.updateCaptivePreview();
      this.$forceUpdate();
    },
    getAllMACAuth() {
      this.userGetAll(
        this.$route.params.id,
        "mac",
        false,
        null,
        null,
        null,
        success => {
          this.macAuth.data = success.body.data;
          this.macAuth.isLoading = false;
        },
        error => {
          this.macAuth.isLoading = false;
          this.macAuth.data = [];
          console.error(error);
        }
      );
    },
    createMacAuth() {
      this.newMACAuth.onAction = true;
      var md5 = require("md5");
      var digest = md5(this.newMACAuth.unit.secret + this.newMACAuth.unit.uuid);
      this.userMACCreate(
        this.newMACAuth,
        digest,
        this.newMACAuth.unit.uuid,
        success => {
          $("#macAuthModal").modal("toggle");
          this.newMACAuth = {
            name: "",
            username: "",
            kbps_down: 0,
            kbps_up: 0,
            onAction: false
          };
          this.getAllMACAuth();
          this.newMACAuth.onAction = false;
        },
        error => {
          console.error(error.body);
          this.newMACAuth.onAction = false;
        }
      );
    },
    deleteMacAuth(id) {
      this.userDelete(
        id,
        success => {
          this.getAllMACAuth();
        },
        error => {
          console.error(error.body);
        }
      );
    },
    getPrivacyInfo(uuid) {
      this.hotspotPrivacy(
        uuid,
        success => {
          this.info.privacy.isLoading = false;
          this.disclaimers.privacyDisclaimers = success.body.privacy_disclaimers;
          this.disclaimers.tosDisclaimers = success.body.tos_disclaimers;
          this.disclaimers.currentPrivacyDisclaimerId = success.body.privacy_selected_id;
          this.disclaimers.currentTosDisclaimerId = success.body.tos_selected_id;
        },
        error => {
          console.error(error.body);
          this.info.privacy.isLoading = false;
        }
      );
    },
    printPrivacy() {
      var finalPrivacy = "data:text/plain;charset=utf-8," + this.selectedPrivacyBody;
      var encodedUri = encodeURI(finalPrivacy);
      var link = document.createElement("a");
      link.setAttribute("href", encodedUri);
      link.setAttribute("download", "privacy" + ".txt");
      link.click();

      var finalTos = "data:text/plain;charset=utf-8," + this.selectedTosBody;
      var encodedUri = encodeURI(finalTos);
      var link = document.createElement("a");
      link.setAttribute("href", encodedUri);
      link.setAttribute("download", "tos" + ".txt");
      link.click();
    },
    clearExpiredFilters() {
      this.vouchers.filters.expiredStart = "";
      this.vouchers.filters.expiredEnd = "";
    },
    clearCreatedFilters() {
      this.vouchers.filters.createdStart = "";
      this.vouchers.filters.createdEnd = "";
    },
    toggleAdvancedFilters() {
      this.advancedFilters = !this.advancedFilters;
      this.$forceUpdate();
    },
    updateDisclaimers() {
      this.hotspotModify(
        this.$route.params.id,
        {
          privacy_disclaimer_id: this.disclaimers.currentPrivacyDisclaimerId,
          tos_disclaimer_id: this.disclaimers.currentTosDisclaimerId
        },
        success => {
          this.disclaimers.updateSuccess = true;

          setTimeout(() => {
            this.disclaimers.updateSuccess = false;
          }, 2000);
          this.getPrivacyInfo(this.info.data.uuid);
        },
        error => {
          console.error(error.body.message);

          this.disclaimers.updateError = true;
          setTimeout(() => {
            this.disclaimers.updateError = false;
          }, 2000);
        }
      );
    }
  }
};
</script>

<style scoped>
textarea {
  width: 100%;
  min-height: 180px;
  resize: vertical;
}

#voucher-coupon {
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

.voucher-traffic {
  font-size: 18px;
}

.voucher-time {
  font-size: 18px;
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

.voucher-max-use {
  font-size: 18px;
}

.voucher-type-span-normal {
  width: 15px;
  height: 15px;
  display: inline-block;
  background: #f9d67a;
  margin-right: 5px;
  margin-left: 3px;
}

.voucher-type-span-auth {
  width: 15px;
  height: 15px;
  display: inline-block;
  background: #c8eb79;
  margin-right: 5px;
  margin-left: 3px;
}

.td-voucher-normal {
  background: #f9d67a;
}

.td-voucher-auth {
  background: #c8eb79;
}

@media print {
  @page {
    margin: 15mm 10mm;
  }

  .pagebreak {
    break-after: always !important;
    page-break-after: always !important;
  }
}

.voucherPrintContainer {
  width: 210mm;
  padding: 5mm;
}

.voucher-print {
  display: inline-block;
  width: 49%;
  border: 0.2mm dashed gray;
  padding: 2mm;
  vertical-align: top;
}

.voucher-print-logo {
  width: 22%;
  float: right;
}

.voucher-print-hs-title {
  font-weight: 400;
  font-size: 5mm;
  margin-bottom: 3mm;
}

.voucher-print-code {
  font-weight: 800;
  font-size: 6mm;
}

.voucher-print-data {
  display: inline-block;
  width: 30mm;
  font-size: 3mm;
}

.width-100 {
  width: 100%;
}

.width-75 {
  width: 75%;
}

.voucher-print-q1 {
  width: 3mm;
  border-bottom: 0.3mm solid black;
}

.voucher-print-q2 {
  padding: 0;
}

.voucher-print-q4 {
  border-top: 0.3mm solid black;
}

.padding-left-2 {
  padding-left: 2mm;
}

.padding-left-2-2 {
  padding-left: 2.2mm;
}

.border-left {
  border-left: 0.3mm solid black;
}

.line-height-5 {
  line-height: 5.5mm;
}

.align-top {
  vertical-align: top;
}

.label-dateFilter {
  margin-right: 10px;
}

.dateFilter {
  display: inline-block;
  vertical-align: middle;
  flex-grow: 1;
}

.flex {
  display: flex;
}

.mg-left-20 {
  margin-left: 20px;
}

.width-p50 {
  width: 50%;
  padding-left: 20px;
}

.no-padding {
  padding: 0;
}

.pointer {
  cursor: pointer;
}

.orange {
  background-color: #ec7a08 !important;
  background-image: linear-gradient(
    to bottom,
    #f39d3c 0,
    #ec7a08 100%
  ) !important;
  border-color: #b35c00 !important;
}

label.block {
  display: block;
  text-align: left;
}

label.block-centered {
  display: block;
  text-align: center;
  margin-bottom: 0.5rem;
}

.captive-portal-pref-inline {
  padding-left: 4rem;
  padding-right: 4rem;
}

.mg-bottom-3r {
  margin-bottom: 3rem;
}

.captive-preview {
  margin: auto;
  border: none;
}

.smartphone {
  width: 350px;
  min-height: 450px;
  margin: auto;
  border: none;
}

.tablet {
  width: 450px;
  min-height: 500px;
  margin: auto;
  border: none;
}

.desktop {
  width: 750px;
  min-height: 400px;
  margin: auto;
  border: none;
}

.captive-mode {
  font-size: 30px;
}

.max-height-80 {
  max-height: 80px;
}

.small-text {
  font-size: 12px;
}

.disclaimer-select {
  margin-top: 0.5rem;
  margin-bottom: 1rem;
}

.tos-title {
  margin-top: 1.5rem;
}

.update-pref-feedback {
  vertical-align: middle;
  margin-left: 0.4em;
  margin-right: 0.4em;
  font-size: 140%;
}
</style>