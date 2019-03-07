<template>
  <div>
    <h2>{{ msg }}</h2>
    <div v-if="isLoading" class="spinner spinner-lg"></div>
    <form class="form-horizontal">
      <div
        v-if="(user.account_type == 'admin') || (user.account_type == 'reseller') && !isLoading"
        class="form-group"
      >
        <label v-if="!isLoading" class="col-sm-3 control-label" for="textInput-markup">Hotspot</label>
        <div v-if="!isLoading" class="col-sm-4">
          <select v-on:change="getMarketingInfo()" v-model="hotspotSearchId" class="form-control">
            <option
              v-for="hotspot in hotspots"
              v-bind:key="hotspot.id"
              v-bind:value="hotspot.id"
            >{{ hotspot.name }} - {{ hotspot.description}}</option>
          </select>
        </div>
      </div>
    </form>

    <form class="form-horizontal" v-on:submit.prevent="saveMarketingInfo()">
      <h3 class="subtitle">{{ $t('marketing.preferences') }}</h3>

      <div v-if="!isLoading && isLoadingMarketing" class="form-group">
        <label class="col-sm-2 control-label" for="textInput-markup">
          <div class="spinner spinner-sm"></div>
        </label>
      </div>
      <div v-if="!isLoadingMarketing" class="form-group">
        <label
          v-if="!isLoadingMarketing"
          class="col-sm-3 control-label"
          for="textInput-markup"
        >{{$t('marketing.enable_service')}}</label>
        <div v-if="!isLoadingMarketing" class="col-sm-4">
          <input v-model="marketingPrefs.marketing_1_enabled" class="form-control" type="checkbox">
        </div>
      </div>

      <div v-if="!isLoadingMarketing && marketingPrefs.marketing_1_enabled" class="form-group">
        <label v-if="!isLoadingMarketing" class="col-sm-3 control-label" for="textInput-markup">
          {{$t('marketing.email_feedback')}}
          <span class="fa fa-facebook-square login-pref-option"></span>
        </label>
        <div v-if="!isLoadingMarketing" class="col-sm-4">
          <input
            :required="marketingPrefs.marketing_1_enabled"
            v-model="marketingPrefs.marketing_2_feedback_email"
            class="form-control"
            type="text"
          >
        </div>
      </div>

      <div v-if="!isLoadingMarketing && marketingPrefs.marketing_1_enabled" class="form-group">
        <label
          v-if="!isLoadingMarketing"
          class="col-sm-3 control-label"
          for="textInput-markup"
        >{{$t('marketing.email_first')}}</label>
        <div v-if="!isLoadingMarketing" class="col-sm-4">
          <input
            v-model="marketingPrefs.marketing_3_first_email_enabled"
            class="form-control"
            type="checkbox"
          >
        </div>
      </div>

      <div
        v-if="!isLoadingMarketing && marketingPrefs.marketing_3_first_email_enabled && marketingPrefs.marketing_1_enabled"
        class="form-group"
      >
        <label
          v-if="!isLoadingMarketing"
          class="col-sm-3 control-label"
          for="textInput-markup"
        >{{$t('marketing.email_first_after')}}</label>
        <div v-if="!isLoadingMarketing" class="col-sm-4">
          <input
            v-model="marketingPrefs.marketing_6_first_after"
            class="form-control"
            type="number"
          >
        </div>
      </div>

      <div v-if="!isLoadingMarketing && marketingPrefs.marketing_1_enabled" class="form-group">
        <label
          v-if="!isLoadingMarketing"
          class="col-sm-3 control-label"
          for="textInput-markup"
        >{{$t('marketing.email_second')}}</label>
        <div v-if="!isLoadingMarketing" class="col-sm-4">
          <input
            v-model="marketingPrefs.marketing_4_second_email_enabled"
            class="form-control"
            type="checkbox"
          >
        </div>
      </div>

      <div
        v-if="!isLoadingMarketing && marketingPrefs.marketing_4_second_email_enabled && marketingPrefs.marketing_1_enabled"
        class="form-group"
      >
        <label
          class="col-sm-3 control-label"
          for="textInput-modal-markup"
        >{{$t('marketing.email_second_after')}}</label>
        <div class="col-sm-7">
          <span class="span-radio">
            <input
              v-model="marketingPrefs.marketing_7_second_after"
              required
              class="form-check-input"
              type="radio"
              name="timeExpiration"
              id="timeExpiration1"
              value="days"
            >
            <label class="form-check-label" for="timeExpiration1">{{$t('marketing.after_days')}}</label>
          </span>
          <span class="span-radio">
            <input
              v-model="marketingPrefs.marketing_7_second_after"
              required
              class="form-check-input"
              type="radio"
              name="timeExpiration"
              id="timeExpiration2"
              value="expiration"
            >
            <label
              class="form-check-label"
              for="timeExpiration2"
            >{{$t('marketing.after_expiration')}}</label>
          </span>
        </div>
      </div>

      <div
        v-if="!isLoadingMarketing && marketingPrefs.marketing_4_second_email_enabled && marketingPrefs.marketing_7_second_after == 'days' && marketingPrefs.marketing_1_enabled"
        class="form-group"
      >
        <label class="col-sm-3 control-label" for="textInput-modal-markup"></label>
        <div class="col-sm-4">
          <input
            v-model="marketingPrefs.marketing_8_second_after_days"
            type="number"
            id="textInput-modal-markup"
            class="form-control"
          >
        </div>
      </div>

      <div v-if="!isLoadingMarketing && marketingPrefs.marketing_1_enabled" class="form-group">
        <label
          v-if="!isLoadingMarketing"
          class="col-sm-3 control-label"
          for="textInput-markup"
        >{{$t('marketing.sms')}}</label>
        <div v-if="!isLoadingMarketing" class="col-sm-4">
          <input
            v-model="marketingPrefs.marketing_5_sms_enabled"
            class="form-control"
            type="checkbox"
          >
        </div>
      </div>

      <div v-if="!isLoadingMarketing && marketingPrefs.marketing_1_enabled" class="form-group">
        <label
          v-if="!isLoadingMarketing"
          class="col-sm-3 control-label"
          for="textInput-markup"
        >{{$t('marketing.threshold')}}</label>
        <div v-if="!isLoadingMarketing" class="col-sm-4">
          <input v-model="marketingPrefs.marketing_9_threshold" class="form-control" type="number">
        </div>
      </div>

      <div v-if="!isLoadingMarketing && marketingPrefs.marketing_1_enabled" class="form-group">
        <label
          v-if="!isLoadingMarketing"
          class="col-sm-3 control-label"
          for="textInput-markup"
        >{{$t('marketing.first_external_url')}}</label>
        <div v-if="!isLoadingMarketing" class="col-sm-4">
          <input
            placeholder="https://www.tripadvisor.it/MYRESTAURANT"
            v-model="marketingPrefs.marketing_10_first_url"
            class="form-control"
            type="url"
          >
        </div>
      </div>
      <div v-if="!isLoadingMarketing && marketingPrefs.marketing_1_enabled" class="form-group">
        <label
          v-if="!isLoadingMarketing"
          class="col-sm-3 control-label"
          for="textInput-markup"
        >{{$t('marketing.second_external_url')}}</label>
        <div v-if="!isLoadingMarketing" class="col-sm-4">
          <input
            placeholder="https://www.booking.com/hotel/MYHOTEL"
            v-model="marketingPrefs.marketing_11_second_url"
            class="form-control"
            type="url"
          >
        </div>
      </div>
      <div v-if="!isLoadingMarketing && marketingPrefs.marketing_1_enabled" class="form-group">
        <label
          v-if="!isLoadingMarketing"
          class="col-sm-3 control-label"
          for="textInput-markup"
        >{{$t('marketing.third_external_url')}}</label>
        <div v-if="!isLoadingMarketing" class="col-sm-4">
          <input v-model="marketingPrefs.marketing_12_third_url" class="form-control" type="url">
        </div>
      </div>

      <div v-if="!isLoadingMarketing" class="form-group">
        <label v-if="!isLoadingMarketing" class="col-sm-3 control-label" for="textInput-markup"></label>
        <div class="col-sm-3">
          <span
            v-show="marketingPrefs.isLoading"
            class="spinner spinner-sm spinner-inline modal-spinner"
          ></span>
          <button
            :disabled="marketingPrefs.isLoading"
            class="btn btn-primary"
            type="submit"
          >{{$t('marketing.save')}}</button>
        </div>
      </div>
    </form>
  </div>
</template>
<script>
import MarketingService from "../services/marketing";
import UtilService from "../services/util";
import StorageService from "../services/storage";
import HotspotService from "../services/hotspot";

export default {
  name: "Unit",
  mixins: [HotspotService, UtilService, MarketingService, StorageService],
  data() {
    return {
      msg: this.$i18n.t("menu.marketing"),
      isLoading: true,
      isLoadingMarketing: true,
      marketingPrefs: {
        marketing_1_enabled: false,
        marketing_2_feedback_email: "",
        marketing_3_first_email_enabled: false,
        marketing_4_second_email_enabled: false,
        marketing_6_first_after: false,
        marketing_7_second_after: "expiration",
        marketing_8_second_after_days: 4,
        marketing_5_sms_enabled: false,
        marketing_9_threshold: 3,
        marketing_10_first_url: "",
        marketing_11_second_url: "",
        marketing_12_third_url: "",
        isLoading: false
      },
      hotspotSearchId: 0,
      hotspots: [],
      user: this.get("loggedUser") || null
    };
  },
  mounted() {
    // get hotspots list
    var context = this;
    this.getAllHotspots(function() {
      context.getMarketingInfo();
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
          var hsId = this.get("marketing_hotspot_id") || this.hotspots[0].id;
          this.hotspotSearchId = hsId;
          this.isLoading = false;

          callback();
        },
        error => {
          console.error(error);
          callback();
        }
      );
    },
    getMarketingInfo() {
      this.isLoadingMarketing = true;
      this.marketingPrefGet(
        this.hotspotSearchId,
        success => {
          for (var i in success.body) {
            var pref = success.body[i];
            if (pref.value == "true") pref.value = true;

            if (pref.value == "false") pref.value = false;
            this.marketingPrefs[pref.key] = pref.value;
          }

          this.marketingPrefs.isLoading = false;
          this.isLoadingMarketing = false;
          this.set("marketing_hotspot_id", this.hotspotSearchId);
        },
        error => {
          console.error(error.body);
          this.isLoadingMarketing = false;
        }
      );
    },
    saveMarketingInfo() {
      this.marketingPrefs.isLoading = true;

      var promises = [];
      for (var i in this.marketingPrefs) {
        var pref = {
          key: i,
          value: this.marketingPrefs[i]
        };

        if (pref.key != "isLoading") {
          promises.push(
            new Promise((resolve, reject) => {
              if (typeof pref.value == "boolean") {
                pref.value = pref.value.toString();
              }

              this.marketingPrefModify(
                this.hotspotSearchId,
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
      }

      // exec promises
      var context = this;
      Promise.all(promises).then(function(response) {
        context.marketingPrefs.isLoading = true;
        context.getMarketingInfo();
      });
    }
  }
};
</script>
<style scoped>
.subtitle {
  padding-left: 15px !important;
}
</style>