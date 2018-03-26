<template>
  <div>
    <h2>{{msg}}</h2>
    <div v-if="isLoading" class="spinner spinner-lg"></div>
    <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller')" class="row select-search">
      <label class="col-sm-2 control-label" for="textInput-markup">Hotspot</label>
      <div class="col-sm-4">
        <select v-on:change="getSessionsByDate()" v-model="hotspotSearchId" class="form-control">
          <option value="0">-</option>
          <option v-for="hotspot in hotspots" v-bind:key="hotspot.id" v-bind:value="hotspot.id">
            {{ hotspot.name }}
          </option>
        </select>
      </div>
    </div>
    <br>
    <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller')" class="row select-search">
       <label class="col-sm-2 control-label" for="textInput-markup">
        <p>{{dataPoints.from.format('DD MMM')}}</p> - <p>{{dataPoints.to.format('DD MMM')}}</p>
      </label>
      <div class="col-sm-4">
       <select v-on:change="getSessionsByDate()" class="form-control" v-model="dateRangeSearchId">
          <option 
                v-for="dateRang in dateRangs" 
                v-bind:key="dateRang.value" 
                v-bind:value="dateRang.value">{{dateRang.display}}</option>
        </select>
      </div>
    </div>
    <br>
    <div v-if="isChartLoading" class="spinner spinner-lg"></div>
    <div v-if="!isChartLoading">
      <div class="panel-heading">
          <h1 class="panel-title">Statistics</h1>
      </div>
      <div class="panel-body">
        <vue-chart type="line" :options="options" :data="chartData"></vue-chart>
      </div>
    </div>
  </div>
</template>
<script>
import UnitService from "../services/unit";
import StorageService from "../services/storage";
import SessionService from '../services/session'
import UserService from '../services/user'
import HotspotService from "../services/hotspot";

import VueChart from "vue-chart-js";
import moment from "moment";
import { extendMoment }  from 'moment-range';
export default {
  name: "Reports",
  components: {
    VueChart,
  },
    mixins: [SessionService, StorageService, UserService, HotspotService],
  data() {
    this.getAllHotspots();    
    return {
      range: null,
      sessionToShow: [],
      newUsersToShow: [],
      isLoading: true,
      isChartLoading: true,
      msg: "Reports",
      dataPoints: {
        to: moment().subtract(7, 'day').startOf('day'),
        from: moment().startOf("day"),
      },
      dateRangs: [
        {
          display: "Last 7 days",
          value: 1
        },
        {
          display: "Last 15 days",
          value: 2
        },
        {
          display: "Last Month",
          value: 3
        }
      ],
      validDate: [],
      dateRangeSearchId: 1,
      dateRangeValue: "Yesterday",
      chartData: {
        labels: [],
        datasets: [
          {
            label: "Sessions",
            data: [],
            backgroundColor:'#e6e6ff'
          },
           {
            label: "New Users",
            data: [],
            backgroundColor:'#cceeff'
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
                  beginAtZero:true,
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
        this.range = moment1.range(moment().subtract(7, 'day').startOf('day').toDate(), moment().startOf("day").toDate());
        this.dataPoints.from = moment().subtract(7, 'day').startOf('day');
        this.dataPoints.to = moment().startOf("day");
          break;
        case 2:
        this.range = moment1.range(moment().subtract(15, 'day').startOf('day').toDate(), moment().endOf("day").toDate());
         this.dataPoints.from = moment().subtract(15, 'day').startOf('day');
        this.dataPoints.to =  moment().endOf("day");
          break;
        case 3:
        this.range = moment1.range(moment().subtract(4, 'week').startOf('week').toDate(), moment().endOf('day').toDate());
        this.dataPoints.from = moment().subtract(4, 'week').startOf('week');
        this.dataPoints.to =  moment().endOf('day');
        default:
          break;
      }
      let selectedRange = Array.from(this.range.by('day'));

      this.validDate = selectedRange.map(function(item){
        return item.toISOString();
      })

      let valueToDisplay = selectedRange.map(function(item){
        return item.format('dd, MMM DD');
      })
      this.chartData.labels = valueToDisplay;
      this.getSession();
      this.getNewUsers();
    },
    // Get All Session between date range
    getSession() {
      this.sessionGetAll(
          this.hotspotSearchId,
          "",
          "",
          this.dataPoints.to.toISOString(),
          this.dataPoints.from.toISOString(),
          success => {
            this.sessions = success.body;
            this.isLoading = false;
            this.implementSessionInChart(this.sessions);
          },
          error => {
            this.isLoading = false;
            this.isChartLoading = false;
            this.sessions = this.rows;
            console.log(error);
          }
        );
    },
    // Get Users
    getNewUsers() {
      this.userGetAll(this.hotspotSearchId, success => {
          this.newUsers = success.body;
          this.implementUsersInChart(this.newUsers);
        }, error => {
          this.isLoading = false
          this.isChartLoading = false;
          this.newUsers = []
          console.log(error)
        })
    },
    getAllHotspots() {
      this.hotspotGetAll(
        success => {
          this.hotspots = success.body;
          $('[data-toggle="tooltip"]').tooltip();
          this.isLoading = false;
        },
        error => {
          console.log(error);
        }
      );
    },
    // Implement Session between date range in Chart
    implementSessionInChart(data) {
      this.validDate.forEach(element => {
        this.sessionToShow = data.map(function(item){
          return item.start_time.includes(element.substring(0, element.indexOf('T')));
        }).filter(function(item) {
              return item == true;
        }).length;
      this.chartData.datasets[0].data.push(this.sessionToShow);
      });
    },
    // Implement Users between date range in Chart
    implementUsersInChart(data) {
      this.validDate.forEach(element => {
        this.newUsersToShow = data.map(function(item){
          return item.valid_from.includes(element.substring(0, element.indexOf('T')));
        }).filter(function(item) {
              return item == true;
        }).length;
      this.chartData.datasets[1].data.push(this.newUsersToShow);
      });
      this.isChartLoading = false;
    }
  }
};
</script>
  <style scoped>
  .control-label p {
    display: inline;
  }
</style>