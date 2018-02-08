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
              <span class="pficon pficon-connected card-info-title"></span>
              {{ $t("dashboard.units") }}
              <span v-if="!totals.units.isLoading" class="right">
                <strong class="soft">{{ totals.units.count}}</strong>
              </span>
              <div v-if="totals.units.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
        </div>
      </div>
      <div class="col-xs-12 col-sm-4 col-md-4 col-lg-2">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              <span class="fa fa-users card-info-title"></span>
              {{ $t("dashboard.users") }}
              <span v-if="!totals.users.isLoading" class="right">
                <strong class="soft">{{ totals.users.count}}</strong>
              </span>
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
      <div class="col-xs-12 col-sm-12 col-md-6">
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
              <dd>{{info.data.account_id}}</dd>
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
      <div class="col-xs-12 col-sm-12 col-md-6">
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
                  <button v-on:click="deleteVoucher(props.row.id)" class="btn btn-danger" type="button">{{ $t("hotspot.delete_voucher") }}</button>
                </td>
              </template>
            </vue-good-table>
          </div>
          <div v-if="!vouchers.isLoading" class="card-pf-footer">
            <div class="dropdown card-pf-time-frame-filter">
              <button v-on:click="createVoucher()" class="btn btn-primary" type="button">{{ $t("hotspot.create_voucher") }}</button>
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
      <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{ $t("hotspot.preferences") }}
              <div v-if="preferences.isLoading" class="spinner spinner-sm right"></div>
            </h2>
          </div>
          <form class="form-horizontal" role="form" v-on:submit.prevent="updatePreferences(preferences.data)">
            <div v-if="!preferences.isLoading" class="card-pf-body">
              <div v-for="pref in preferences.data" :key="pref.key" class="form-group">
                <label class="col-sm-4 control-label" for="textInput-markup">{{pref.key}}</label>
                <div class="col-sm-6">
                  <input v-model="pref.value" :type="getInputType(pref.value)" id="textInput-markup" class="form-control">
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

  import HotspotAction from '../../directives/HotspotAction.vue';

  export default {
    name: 'HotspotDetails',
    mixins: [HotspotService, PreferenceService, AccountService, UnitService, UserService, DeviceService, SessionService,
      StorageService, UtilService
    ],
    components: {
      hotspotAction: HotspotAction
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
        preferences: {
          isLoading: true,
          data: {}
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
      }
    },
    methods: {
      createVoucher() {
        this.vouchers.isLoading = true
        this.hotspotCreateVoucher({
          hotspot_id: parseInt(this.$route.params.id),
          code: this.generateVoucher(),
        }, success => {
          this.vouchers.isLoading = false
          this.getVouchers()
        }, error => {
          console.log(error.body)
          this.vouchers.isLoading = false
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
        this.hotspotGetVouchers(success => {
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
        this.sessionGetAll(this.$route.params.id, success => {
          this.totals.sessions.count = success.body.length
          this.totals.sessions.isLoading = false
        }, error => {
          console.log(error.body)
          this.totals.sessions.isLoading = false
        })
      },
      getPreferences() {
        this.hsPrefGet(this.$route.params.id, success => {
          for (var p in success.body) {
            var pref = success.body[p]
            if (pref.value === "true") {
              pref.value = true
            }
            if (pref.value === "false") {
              pref.value = false
            }
          }
          this.preferences.data = success.body
          this.preferences.isLoading = false
        }, error => {
          console.log(error.body)
        })
      },
      updatePreferences() {
        this.preferences.isLoading = true
        // create promises array
        var promises = []
        for (var i in this.preferences.data) {
          var pref = this.preferences.data[i]
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
      }
    }
  }

</script>

<style scoped>


</style>
