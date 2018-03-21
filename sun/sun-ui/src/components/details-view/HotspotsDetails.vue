<template>
  <div>
    <h2>Hotspot
      <strong class="soft">{{ info.data.name }}</strong>
    </h2>

    <div class="row row-cards-pf">
      <div class="col-xs-12 col-sm-6 col-md-6 col-lg-3">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              <span class="pficon pficon-users card-info-title"></span>
              {{ $t("dashboard.accounts") }}
              <span v-if="!totals.accounts.isLoading" class="right">
                <strong class="soft">{{ totals.accounts.count}}</strong>
              </span>
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
              <span class="pficon fa-mobile card-info-title"></span>
              {{ $t("dashboard.devices") }}
              <span v-if="!totals.devices.isLoading" class="right">
                <strong class="soft">{{ totals.devices.count}}</strong>
              </span>
              <div v-if="totals.devices.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
        </div>
      </div>
      <div class="col-xs-12 col-sm-4 col-md-4 col-lg-2">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              <span class="fa fa-list card-info-title"></span>
              {{ $t("dashboard.sessions") }}
              <span v-if="!totals.sessions.isLoading" class="right">
                <strong class="soft">{{ totals.sessions.count}}</strong>
              </span>
              <div v-if="totals.sessions.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
        </div>
      </div>
      <div class="col-xs-12 col-sm-12 col-md-6">
      </div>
    </div>
    <div class="row row-cards-pf">
      <div v-if="user.account_type == 'admin' || user.account_type == 'reseller'" class="col-xs-12 col-sm-12 col-md-6">
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
      <div v-if="preferences.vouchersAvailable" :class="['col-xs-12 col-sm-12', user.account_type == 'admin' || user.account_type == 'reseller' ? 'col-md-6' : 'col-md-12']">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{ $t("hotspot.vouchers") }}
              <div v-if="vouchers.isLoading" class="spinner spinner-sm right"></div>
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
                <td class="fancy">{{ props.row.expires | formatDate }}</td>
                <td>
                  <button v-on:click="printVoucher(props.row.code)" class="btn btn-primary" type="button">
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

    <div class="row row-cards-pf">
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
                    :alertOnError="false" @change="onChanged(pref)" :width="100" :height="100" accept="image/jpeg, image/png"
                    :crop="false" :zIndex="1000" :customStrings="uploadLangstexts" removeButtonClass="btn btn-danger" buttonClass="btn btn-default">

                  </picture-input>
                  <span v-if="pref.onError" class="help-block">{{$t('upload_file_exceed')}}</span>

                  <vue-editor :editorToolbar="customToolbar" v-if="pref.key == 'captive_6_description'" v-model="pref.value"></vue-editor>

                  <sketch-picker @input="onUpdate(pref.value)" class="absolute-center" v-if="pref.key == 'captive_7_background'" v-model="pref.value"
                  />

                  <input v-if="pref.key != 'captive_6_description' && pref.key != 'captive_3_logo' && pref.key != 'captive_5_banner' && pref.key != 'captive_7_background'"
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
            </form>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-default" data-dismiss="modal">Cancel</button>
            <button v-on:click="createVoucher()" type="button" class="btn btn-primary">Save</button>
          </div>
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

      return {
        info: {
          isLoading: true,
          data: {}
        },
        vouchers: {
          isLoading: true,
          data: []
        },
        vouchersCount: 1,
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
          }
        },
        columns: [{
            label: this.$i18n.t('hotspot.code'),
            field: 'code',
            filterable: false,
            sortable: false,
          }, {
            label: this.$i18n.t('hotspot.expires'),
            field: 'expires',
            filterable: false,
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
        customToolbar: [
          ['bold', 'italic', 'underline'],
          ['image', 'code-block']
        ],      
        userLink:{
          name:'Users',
          params:{
            hotspotId: this.$route.params.id
          }
        },
        unitLink:{
          name: 'Units',
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
          this.totals.units.count = success.body.length
          this.totals.units.isLoading = false
        }, error => {
          console.log(error.body)
          this.totals.units.isLoading = false
        })
        this.userGetAll(this.$route.params.id, success => {
          this.totals.users.count = success.body.length
          this.totals.users.isLoading = false
        }, error => {
          console.log(error.body)
          this.totals.users.isLoading = false
        })
        this.deviceGetAll(this.$route.params.id, success => {
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
            $('#captive-preview').css('background-color', backgroundColor);
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
        var v = 0
        var doc = new jsPDF("portrait", "mm", "a4");
        doc.setFontSize(22);
        doc.text(20, 18, '-'.repeat(20))
        doc.setFontSize(15);
        doc.text(20, (v + 1) * (2.5) + 20, this.$i18n.t('hotspot.voucher_code'));
        doc.setFontSize(22);
        doc.text(20, (v + 1) * (2.5) + 30, voucher);
        doc.text(20, (v + 1) * (2.5) + 38, '-'.repeat(20))
        doc.autoPrint();
        window.open(doc.output('bloburl'), '_blank');
      },
      printAllVoucher() {
        var doc = new jsPDF("portrait", "mm", "a4");
        var pageHeight = doc.internal.pageSize.height;

        for (var v in this.vouchers.data) {
          var voucher = this.vouchers.data[v]
          if (parseInt(v) != 0 && parseInt(v) % 11 == 0) {
            doc.addPage();
          }

          doc.setFontSize(22);
          doc.text(20, 18, '-'.repeat(20))
          doc.setFontSize(15);
          doc.text(20, (((v % 11) + 1) * (2.5) + 20) + (22.5 * (v % 11)), this.$i18n.t('hotspot.voucher_code'));
          doc.setFontSize(22);
          doc.text(20, (((v % 11) + 1) * (2.5) + 30) + (22.5 * (v % 11)), voucher.code);
          doc.text(20, (((v % 11) + 1) * (2.5) + 38) + (22.5 * (v % 11)), '-'.repeat(20))
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
    }
  }

</script>

<style scoped>
textarea {
  width: 100%;
  min-height: 180px;
  resize: vertical;
}

.card-link{
  color: black;
}
.card-link:hover{
  color: #26a9a3;
}

</style>
