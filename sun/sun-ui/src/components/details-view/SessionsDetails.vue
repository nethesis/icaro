<template>
  <div>
    <h2>Session
      <strong class="soft">{{ info.data.session_key }}</strong>
    </h2>

    <div class="row row-cards-pf">
      <div class="col-xs-12 col-sm-12 col-md-6">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">
              {{info.data.session_key}}
              <div v-if="info.isLoading" class="spinner spinner-sm right"></div>
              <div v-if="!info.isLoading" class="fa fa-list right"></div>
            </h2>
          </div>
          <div v-if="!info.isLoading" class="card-pf-body">
            <div class="list-details">
              <dt>{{ $t("session.bytes_up") }}</dt>
              <dd>{{info.data.bytes_up | byteFormat}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("session.bytes_down") }}</dt>
              <dd>{{info.data.bytes_down | byteFormat}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("session.duration") }}</dt>
              <dd>{{info.data.duration | secondsInHour}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("session.auth_time") }}</dt>
              <dd>{{info.data.auth_time}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("session.start_time") }}</dt>
              <dd>{{info.data.start_time}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("session.stop_time") }}</dt>
              <dd>{{info.data.stop_time}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("session.update_time") }}</dt>
              <dd>{{info.data.update_time}}</dd>
            </div>
          </div>
        </div>
      </div>

      <div class="col-xs-12 col-sm-12 col-md-6">
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">{{ $t('session.hotspot') }}
              <div v-if="info.hotspot.isLoading" class="spinner spinner-sm right"></div>
              <div v-if="!info.hotspot.isLoading" class="fa fa-wifi right"></div>
            </h2>
          </div>
          <div class="card-pf-body">
            <div class="list-details">
              <dt>{{ $t("session.hotspot_name") }}</dt>
              <dd>{{info.hotspot.name || ""}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("session.hotspot_description") }}</dt>
              <dd>{{info.hotspot.description || ""}}</dd>
            </div>
          </div>
        </div>
        <div class="card-pf card-pf-accented">
          <div class="card-pf-heading">
            <h2 class="card-pf-title">{{ $t('session.unit') }}
              <div v-if="info.unit.isLoading" class="spinner spinner-sm right"></div>
              <div v-if="!info.unit.isLoading" class="pficon pficon-connected right"></div>
            </h2>
          </div>
          <div class="card-pf-body">
            <div class="list-details">
              <dt>{{ $t("session.unit_uuid") }}</dt>
              <dd>{{info.unit.uuid || ""}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("session.unit_name") }}</dt>
              <dd>{{info.unit.description || ""}}</dd>
            </div>
            <div class="list-details">
              <dt>{{ $t("session.mac_address") }}</dt>
              <dd>{{info.unit.mac_address || ""}}</dd>
            </div>
          </div>
        </div>
        <div class="col-xs-12 col-sm-12 col-md-6">
          <div class="card-pf card-pf-accented">
            <div class="card-pf-heading">
              <h2 class="card-pf-title">{{ $t('session.user') }}
                <div v-if="info.user.isLoading" class="spinner spinner-sm right"></div>
                <div v-if="!info.user.isLoading" :class="[getUserTypeIcon(info.user.account_type), 'right']"></div>
              </h2>
            </div>
            <div class="card-pf-body">
              <div class="list-details">
                <dt>{{ $t("session.user_name") }}</dt>
                <dd>{{ info.user.name || "" }} {{ "(" + info.user.username + ")" || ""}}</dd>
              </div>
              <div class="list-details">
                <dt>{{ $t("session.created") }}</dt>
                <dd>{{info.user.created || ""}}</dd>
              </div>
            </div>
          </div>
        </div>
        <div class="col-xs-12 col-sm-12 col-md-6">
          <div class="card-pf card-pf-accented">
            <div class="card-pf-heading">
              <h2 class="card-pf-title">{{ $t('session.device') }}
                <div v-if="info.device.isLoading" class="spinner spinner-sm right"></div>
                <div v-if="!info.device.isLoading" class="fa fa-mobile right"></div>
              </h2>
            </div>
            <div class="card-pf-body">
              <div class="list-details">
                <dt>{{ $t("session.mac_address") }}</dt>
                <dd>{{info.device.mac_address || ""}}</dd>
              </div>
              <div class="list-details">
                <dt>{{ $t("session.created") }}</dt>
                <dd>{{info.device.created || ""}}</dd>
              </div>
            </div>
          </div>
        </div>
      </div>

    </div>

  </div>
</template>

<script>
  import SessionService from '../../services/session';
  import HotspotService from '../../services/hotspot';
  import UserService from '../../services/user';
  import DeviceService from '../../services/device';
  import UnitService from '../../services/unit';
  import StorageService from '../../services/storage';
  import UtilService from '../../services/util';

  export default {
    name: 'SessionsDetails',
    mixins: [SessionService, StorageService, UnitService, UtilService, HotspotService, DeviceService, UserService],
    filters: {
      byteFormat: require('vue-filters-kit/filters/byteFormatter'),
      secondsInHour: function (value) {
        let hours = parseInt(Math.floor(value / 3600));
        let minutes = parseInt(Math.floor((value - (hours * 3600)) / 60));
        let seconds = parseInt((value - ((hours * 3600) + (minutes * 60))) % 60);

        let dHours = (hours > 9 ? hours : '0' + hours);
        let dMins = (minutes > 9 ? minutes : '0' + minutes);
        let dSecs = (seconds > 9 ? seconds : '0' + seconds);

        return dHours + "h " + dMins + "m " + dSecs + "s";
      }
    },
    data() {
      // get account info
      this.getInfo()

      return {
        info: {
          isLoading: true,
          unit: {},
          hotspot: {},
          device: {},
          user: {},
          data: {}
        },
      }
    },
    // enable tooltips after rendering
    updated: function () {
      $('[data-toggle="tooltip"]').tooltip()
    },
    methods: {
      getInfo() {
        this.sessionGet(this.$route.params.id, success => {
          this.info.data = success.body
          this.info.isLoading = false
          this.getUnitInfo(success.body.unit_id)
          this.getUserInfo(success.body.user_id)
          this.getHotspotInfo(success.body.hotspot_id)
          this.getDeviceInfo(success.body.device_id)
        }, error => {
          this.info.isLoading = false
          console.log(error.body)
        })
      },
      getUnitInfo(unitId) {
        this.unitGet(unitId, success => {
          this.info.unit = success.body
          this.info.unit.isLoading = false
        }, error => {
          this.info.unit.isLoading = false
          console.log(error.body)
        })
      },
      getUserInfo(userId) {
        this.userGet(userId, success => {
          this.info.user = success.body
          this.info.user.isLoading = false
        }, error => {
          this.info.user.isLoading = false
          console.log(error.body)
        })
      },
      getHotspotInfo(hotspotId) {
        this.hotspotGet(hotspotId, success => {
          this.info.hotspot = success.body
          this.info.hotspot.isLoading = false
        }, error => {
          this.info.hotspot.isLoading = false
          console.log(error.body)
        })
      },
      getDeviceInfo(deviceId) {
        this.deviceGet(deviceId, success => {
          this.info.device = success.body
          this.info.device.isLoading = false
        }, error => {
          this.info.device.isLoading = false
          console.log(error.body)
        })
      },

    }
  }

</script>
