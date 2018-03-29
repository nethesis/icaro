<template>
  <div>
    <h2>{{msg}}</h2>
    <div v-if="isChartLoading" class="spinner spinner-lg"></div>
    <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller') && !isChartLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <label v-if="!isChartLoading" class="col-sm-2 control-label" for="textInput-markup">Hotspot</label>
      <div v-if="!isChartLoading" class="col-sm-4">
        <select v-on:change="getSessionsByDate()" v-model="hotspotSearchId" class="form-control">
          <option value="0">-</option>
          <option v-for="hotspot in hotspots" v-bind:key="hotspot.id" v-bind:value="hotspot.id">
            {{ hotspot.name }}
          </option>
        </select>
      </div>
    </div>
    <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller') && !isChartLoading" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
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
    <div v-if="!isChartLoading">
      <div class="panel-heading">
        <h1 class="panel-title">{{ $t('report.statistic') }}</h1>
      </div>
      <div class="panel-body">
        <vue-chart type="line" :width="150" :heigth="150" :options="options" :data="chartData"></vue-chart>
        <br>
        <report-statistics 
        :chartLabels="chartData.labels"
        :chartDateRange="validDate"
        :newUsersReport="newUsers"
        :sessionsReport="sessions"
        ></report-statistics>
      </div>
    </div>
  </div>
</template>
<script>
  import UnitService from "../services/unit";
  import StorageService from "../services/storage";
  import HistoryService from '../services/history'
  import UserService from '../services/user'
  import HotspotService from "../services/hotspot";
  import ReportStatistics from '../components/details-view/ReportStatistics'

  import VueChart from "vue-chart-js";
  import moment from "moment";
  import {
    extendMoment
  } from 'moment-range';
  export default {
    name: "Reports",
    components: {
      VueChart,
      ReportStatistics
    },
    mixins: [HistoryService, StorageService, UserService, HotspotService],
    data() {
      this.getAllHotspots();
      return {
        range: null,
        sessionToShow: [],
        newUsersToShow: [],
        isChartLoading: true,
        msg: this.$i18n.t('report.reports'),
        dataPoints: {
          to: moment().subtract(7, 'day').startOf('day'),
          from: moment().startOf("day"),
        },
        dateRangs: [{
            display: this.$i18n.t('report.last_7_days'),
            value: 1
          },
          {
            display: this.$i18n.t('report.last_15_days'),
            value: 2
          },
          {
            display: this.$i18n.t('report.last_month'),
            value: 3
          }
        ],
        validDate: [],
        dateRangeSearchId: 1,
        chartData: {
          labels: [],
          datasets: [{
              label: this.$i18n.t('report.session'),
              data: [],
              backgroundColor: '#e6e6ff'
            },
            {
              label: this.$i18n.t('report.new_user'),
              data: [],
              backgroundColor: '#cceeff'
            }
          ],
        },
        options: {
          elements: {
            line: {
              tension: 0
            }
          },
          scales: {
            xAxes: [{
              gridLines: {
                zeroLineColor: 'transparent'
              }
            }],
            yAxes: [{
              ticks: {
                beginAtZero: true,
                stepSize: 1,
              },
              gridLines: {
                display: false,
                drawBorder: false
              }
            }]
          }
        },
        newUsers: [],
        sessions: [],
        hotspots: [],
        hotspotSearchId: 0,
        user: this.get("loggedUser") || null
      };
    },
    mounted() {
      this.getSessionsByDate();
    },
    methods: {
      getSessionsByDate() {
        const moment1 = extendMoment(moment);
        this.isChartLoading = true;
        this.chartData.labels = [];
        this.chartData.datasets[0].data = [];
        this.chartData.datasets[1].data = [];


        switch (this.dateRangeSearchId) {
          case 1:
            this.range = moment1.range(moment().subtract(8, 'day').startOf('day').toDate(), moment().subtract(1, 'day')
              .endOf("day").toDate());
            this.dataPoints.from = moment().subtract(8, 'day').startOf('day');
            this.dataPoints.to = moment().subtract(1, 'day').endOf("day");
            break;
          case 2:
            this.range = moment1.range(moment().subtract(15, 'day').startOf('day').toDate(), moment().subtract(1, 'day')
              .endOf("day"));
            this.dataPoints.from = moment().subtract(15, 'day').startOf('day');
            this.dataPoints.to = moment().subtract(1, 'day').endOf("day");
            break;
          case 3:
            this.range = moment1.range(moment().subtract(4, 'week').startOf('week').toDate(), moment().subtract(1,
              'day').endOf("day"));
            this.dataPoints.from = moment().subtract(4, 'week').startOf('week');
            this.dataPoints.to = moment().subtract(1, 'day').endOf("day");
          default:
            break;
        }
        let selectedRange = Array.from(this.range.by('day'));

        this.validDate = selectedRange.map(function (item) {
          return item.format('YYYY-MM-DD');
        })

        let valueToDisplay = selectedRange.map(function (item) {
          return item.format('DD MMM YYYY');
        })
        this.chartData.labels = valueToDisplay;
        this.getSession();
      },
      // Get All Session between date range
      getSession() {
        this.historiesGetAll(
          this.hotspotSearchId,
          "",
          "",
          this.dataPoints.from.toISOString(),
          this.dataPoints.to.toISOString(),
          success => {
            this.sessions = success.body;
            this.getNewUsers()
          },
          error => {
            this.sessions = [];
            this.getNewUsers()
            console.log(error);
          }
        );
      },
      // Get Users
      getNewUsers() {
        this.userGetAll(this.hotspotSearchId, success => {
          this.newUsers = success.body;
          this.implementDataInChart();
          this.isChartLoading = false;
        }, error => {
          this.isChartLoading = false
          this.newUsers = []
          this.implementDataInChart();
          console.log(error)
        })
      },
      getAllHotspots() {
        this.hotspotGetAll(
          success => {
            this.hotspots = success.body;
            $('[data-toggle="tooltip"]').tooltip();
          },
          error => {
            console.log(error);
          }
        );
      },
      // Implement Session and newUser between date range in Chart
      implementDataInChart() {
        this.validDate.forEach(element => {
          this.sessionToShow = this.sessions.map(function (item) {
            return item.start_time.substring(0, 10) === element;
          }).filter(function (item) {
            return item == true;
          }).length;
          this.chartData.datasets[0].data.push(this.sessionToShow);
        });

        this.validDate.forEach(element => {
          this.newUsersToShow = this.newUsers.map(function (item) {
            return item.valid_from.substring(0, 10) === element;
          }).filter(function (item) {
            return item == true;
          }).length;
          this.chartData.datasets[1].data.push(this.newUsersToShow);
        });
      },
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

</style>
