<template>
  <div>
    <h2>{{msg}}</h2>
    <div v-if="isChartLoading" class="spinner spinner-lg"></div>

    <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller') && !isChartLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <label class="col-sm-2 control-label" for="textInput-markup">Hotspot</label>
      <div class="col-sm-4">
        <select v-on:change="getSessionsByDate()" v-model="hotspotSearchId" class="form-control">
          <option v-for="hotspot in hotspots" v-bind:key="hotspot.id" v-bind:value="hotspot.id">
            {{ hotspot.name }} - {{ hotspot.description}}
          </option>
        </select>
      </div>
    </div>

    <div v-if="!isChartLoading">
      <h2 :class="['graphs-container',(user.account_type == 'admin') || (user.account_type == 'reseller') ? 'title-graphs' : '']">{{ $t('report.current_situation') }}</h2>
      <actual-report-statistics class="graphs-container adjust-top" :todayConnections="connections"></actual-report-statistics>
    </div>

    <div v-if="!isChartLoading">
      <h2 class="graphs-container title-graphs">{{ $t('report.history_situation') }}</h2>
      <div class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
        <label v-if="!isChartLoading" class="col-sm-2 control-label" for="textInput-markup" id="lbl-date-range">
          <p>{{dataPoints.from.format('DD MMM YYYY')}}</p> -
          <p>{{dataPoints.to.format('DD MMM YYYY')}}</p>
        </label>
        <div v-if="!isChartLoading" class="col-sm-4">
          <select v-on:change="getSessionsByDate()" class="form-control" v-model="dateRangeSearchId">
            <option v-for="dateRang in dateRangs" v-bind:key="dateRang.value" v-bind:value="dateRang.value">{{dateRang.display}}</option>
          </select>
        </div>
      </div>
      <report-statistics class="graphs-container" :chartLabels="labels" :chartDateRange="validDate" :newUsersReport="newUsers"
        :sessionsReport="sessions"></report-statistics>
    </div>
    <div v-if="!isChartLoading">
      <h2 class="graphs-container title-graphs">{{ $t('report.sms_reports') }}</h2>
      <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6">
      <sms-report-statistics class="graphs-container adjust-top" :isYear="true" :chartLabels="labelsSms" :chartDateRange="validDateSms" :smsData="sms"></sms-report-statistics>
      </div>
      <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6">
      <sms-report-statistics class="graphs-container adjust-top" :isYear="false" :chartLabels="labels" :chartDateRange="validDate" :smsData="sms"></sms-report-statistics>
      </div>
    </div>
  </div>
</template>
<script>
import UnitService from "../services/unit";
import StorageService from "../services/storage";
import HistoryService from "../services/history";
import SessionService from "../services/session";
import UserService from "../services/user";
import HotspotService from "../services/hotspot";
import StatsService from "../services/stats";

import ReportStatistics from "../components/details-view/ReportStatistics";
import ActualReportStatistics from "../components/details-view/ActualReportStatistics";
import SmsReportStatistics from "../components/details-view/SmsReportStatistics";

import VueChart from "vue-chart-js";
import moment from "moment";
import { extendMoment } from "moment-range";
export default {
  name: "Reports",
  components: {
    VueChart,
    ReportStatistics,
    ActualReportStatistics,
    SmsReportStatistics
  },
  mixins: [
    HistoryService,
    StorageService,
    UserService,
    HotspotService,
    SessionService,
    StatsService
  ],
  data() {
    return {
      range: null,
      isChartLoading: true,
      msg: this.$i18n.t("report.reports"),
      dataPoints: {
        to: moment()
          .subtract(7, "day")
          .startOf("day"),
        from: moment().startOf("day")
      },
      dateRangs: [
        {
          display: this.$i18n.t("report.last_7_days"),
          value: 1
        },
        {
          display: this.$i18n.t("report.last_15_days"),
          value: 2
        },
        {
          display: this.$i18n.t("report.last_month"),
          value: 3
        }
      ],
      validDate: [],
      validDateSms: [],
      dateRangeSearchId: this.get("reports_date_range_id") || 1,
      newUsers: [],
      sessions: [],
      hotspots: [],
      connections: [],
      sms: [],
      labels: [],
      labelsSms: [],
      hotspotSearchId: 0,
      user: this.get("loggedUser") || null
    };
  },
  mounted() {
    if (this.$route.params.hotspotId !== undefined) {
      this.hotspotSearchId = this.$route.params.hotspotId;
    }

    var context = this;
    this.getAllHotspots(function() {
      context.getSessionsByDate();
      context.getSmsData();
    });
  },
  methods: {
    getSessionsByDate() {
      this.set(
        "reports_hotspot_id",
        this.hotspotSearchId || this.get("reports_hotspot_id") || 0
      );
      this.set(
        "reports_date_range_id",
        this.dateRangeSearchId || this.get("reports_date_range_id") || 1
      );

      const moment1 = extendMoment(moment);
      this.isChartLoading = true;
      this.labels = [];
      this.labelsSms = [];
      switch (this.dateRangeSearchId) {
        case 1:
          this.range = moment1.range(
            moment()
              .utc()
              .subtract(8, "day")
              .startOf("day")
              .toDate(),
            moment()
              .utc()
              .subtract(1, "day")
              .endOf("day")
              .toDate()
          );
          this.dataPoints.from = moment()
            .utc()
            .subtract(8, "day")
            .startOf("day");
          this.dataPoints.to = moment()
            .utc()
            .subtract(1, "day")
            .endOf("day");
          break;
        case 2:
          this.range = moment1.range(
            moment()
              .utc()
              .subtract(15, "day")
              .startOf("day")
              .toDate(),
            moment()
              .utc()
              .subtract(1, "day")
              .endOf("day")
          );
          this.dataPoints.from = moment()
            .utc()
            .subtract(15, "day")
            .startOf("day");
          this.dataPoints.to = moment()
            .utc()
            .subtract(1, "day")
            .endOf("day");
          break;
        case 3:
          this.range = moment1.range(
            moment()
              .utc()
              .subtract(4, "week")
              .startOf("week")
              .toDate(),
            moment()
              .utc()
              .subtract(1, "day")
              .endOf("day")
          );
          this.dataPoints.from = moment()
            .utc()
            .subtract(4, "week")
            .startOf("week");
          this.dataPoints.to = moment()
            .utc()
            .subtract(1, "day")
            .endOf("day");
        default:
          break;
      }
      let rangeYear = moment1.range(
        moment()
          .utc()
          .startOf("year"),
        moment().utc()
      );
      let selectedRange = Array.from(this.range.by("day"));
      let selectedRangeSms = Array.from(rangeYear.by("month"));

      this.validDate = selectedRange.map(function(item) {
        return item.format("YYYY-MM-DD");
      });
      this.validDateSms = selectedRangeSms.map(function(item) {
        return item.format("YYYY-MM-DD");
      });

      let valueToDisplay = selectedRange.map(function(item) {
        return item.format("DD MMM");
      });
      let valueToDisplaySms = selectedRangeSms.map(function(item) {
        return item.format("MMM");
      });

      this.labels = valueToDisplay;
      this.labelsSms = valueToDisplaySms;

      this.getSession();
      this.getSmsData();
    },
    // Get All Session between date range
    getSession() {
      this.historiesGetAll(
        this.hotspotSearchId,
        "",
        "",
        this.dataPoints.from.toISOString(),
        this.dataPoints.to.toISOString(),
        null,
        null,
        success => {
          this.sessions = success.body.data;
          this.getNewUsers();
        },
        error => {
          this.sessions = [];
          this.getNewUsers();
          console.error(error);
        }
      );
    },
    // Get Users
    getNewUsers() {
      this.userGetAll(
        this.hotspotSearchId,
        null,
        false,
        null,
        null,
        success => {
          this.newUsers = success.body.data;
          this.getTodayUsersLogin();
        },
        error => {
          this.newUsers = [];
          this.getTodayUsersLogin();
          console.error(error);
        }
      );
    },
    getAllHotspots(callback) {
      this.hotspotGetAll(
        null,
        null,
        success => {
          this.hotspots = success.body.data;

          var hsId = this.get("reports_hotspot_id") || this.hotspots[0].id;
          if (
            this.$parent.user.info.type == "customer" ||
            this.$parent.user.info.type == "desk"
          ) {
            hsId = this.$parent.user.info.hotspot_id;
          }
          this.hotspotSearchId = hsId;

          $('[data-toggle="tooltip"]').tooltip();

          callback();
        },
        error => {
          console.error(error);
          callback();
        }
      );
    },
    getTodayUsersLogin() {
      this.sessionGetAll(
        this.hotspotSearchId,
        "",
        "",
        moment()
          .utc()
          .startOf("day")
          .toISOString(),
        "",
        null,
        null,
        success => {
          this.connections = success.body.data;
          this.isChartLoading = false;
        },
        error => {
          this.connections = [];
          this.isChartLoading = false;
          console.error(error);
        }
      );
    },
    getSmsData() {
      this.statsSMSSentByHotspot(
        this.hotspotSearchId,
        success => {
          this.sms = success.body;
          this.isChartLoading = false;
        },
        error => {
          console.error(error.body);
          this.sms = [];
          this.isChartLoading = false;
        }
      );
    }
  }
};
</script>
<style scoped>
.control-label p {
  display: inline;
}

#lbl-date-range {
  padding-right: 8px;
}

.graphs-container {
  margin: 10px;
}

.title-graphs {
  margin-top: 60px;
}

.adjust-top {
  margin-top: -30px;
}
</style>
