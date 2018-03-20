<template>
  <div>
    <h2>{{ $t('unit.unit') }}
      <strong class="soft">{{ info.unit.uuid }}</strong>
    </h2>

      <div class="row row-cards-pf">
          <div class="col-xs-12 col-sm-12 col-md-6">
            <div class="card-pf card-pf-accented">
                <div class="card-pf-heading">
                <h2 class="card-pf-title">{{ $t('unit.unit') }}
                    <div v-if="info.isLoading" class="spinner spinner-sm right"></div>
                    <div v-if="!info.unit.isLoading" class="pficon pficon-connected right"></div>
                </h2>
                </div>
                <div class="card-pf-body">
                <div class="list-details">
                    <dt>{{ $t("unit.uuid") }}</dt>
                    <dd>{{info.unit.uuid || ""}}</dd>
                </div>
                <div class="list-details">
                    <dt>{{ $t("unit.description") }}</dt>
                    <dd>{{info.unit.description || ""}}</dd>
                </div>
                <div class="list-details">
                    <dt>{{ $t("unit.mac_address") }}</dt>
                    <dd>{{info.unit.mac_address || ""}}</dd>
                </div>
                </div>
            </div>
          </div>
          <div class="col-xs-12 col-sm-12 col-md-6">
            <div class="card-pf card-pf-accented">
                <div class="card-pf-heading">
                    <h2 class="card-pf-title">{{ $t('unit.hotspot_name') }}
                    <div v-if="info.isLoading" class="spinner spinner-sm right"></div>
                    <div v-if="!info.hotspot.isLoading" class="fa fa-wifi right"></div>
                    </h2>
                </div>
                <div class="card-pf-body">
                    <div class="list-details">
                    <dt>{{ $t("unit.hotspot_name") }}</dt>
                    <dd>{{info.hotspot.name || ""}}</dd>
                    </div>
                    <div class="list-details">
                    <dt>{{ $t("unit.hotspot_description") }}</dt>
                    <dd>{{info.hotspot.description || ""}}</dd>
                    </div>
                </div>
            </div>
          </div>
      </div>
  </div>
</template>
<script>
  import UnitService from '../../services/unit';
  import HotspotService from '../../services/hotspot';
  import UtilService from '../../services/util';
  import StorageService from '../../services/storage';
    

export default {
  name: 'UnitsDetails',
  mixins: [UnitService, HotspotService, UtilService, StorageService],
  data() {
      //get unit info
      this.getInfo()

      return {
          info: {
            isLoading: true,
            unit: {},
            hotspot: {},
          }
      }
  },
  methods: {
      getInfo() {
          this.unitGet(this.$route.params.id, success => {
              this.info.unit = success.body
              this.info.isLoading = false
              this.getHotspotInfo(success.body.hotspot_id)
          }, error => {
              this.info.unit.isLoading = false
              console.log(error.body);
          })
      },
      getHotspotInfo(hotspotId) {
          this.hotspotGet(hotspotId, success => {
              this.info.hotspot = success.body
              this.info.isLoading = false
          }, error => {
              this.info.isLoading = false
              console.log(error.body);
          })
      }
  }
}
</script>