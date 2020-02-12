<template>
  <div>
    <h2>{{msg}}</h2>

    <div v-if="(user.account_type == 'admin') || (user.account_type == 'reseller')" class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <label class="col-sm-2 control-label" for="textInput-markup">Hotspot</label>
      <div class="col-sm-4">
        <select v-on:change="getCurrentStats()" v-model="hotspotSearchId" class="form-control">
          <option v-for="hotspot in hotspots" v-bind:key="hotspot.id" v-bind:value="hotspot.id">
            {{ hotspot.name }} - {{ hotspot.description}}
          </option>
        </select>
      </div>
    </div>

    <div>
      <h2 :class="['graphs-container',(user.account_type == 'admin') || (user.account_type == 'reseller') ? 'title-graphs' : '']">{{
        $t('report.current_situation') }}</h2>
      <div class="row no-margin averages">
        <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6 border-right avg-cont">
          <h3>{{ $t('report.users_currently_connected') }}</h3>
          <h1>{{ averages.current[0] }}</h1>
        </div>
        <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6 avg-cont">
          <h3>{{ $t('report.total_daily_connections') }}</h3>
          <h1>{{ averages.current[1] }}</h1>
        </div>
      </div>
      <div class="graph-divider-blank"></div>
      <div class="row no-margin">
        <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12">
          <div v-if="loaders.current" class="spinner spinner-lg"></div>
          <vue-chart v-show="!loaders.current" :height="100" type="line" :options="charts.current.options" :data="charts.current"></vue-chart>
        </div>
      </div>
    </div>

    <div class="graph-divider"></div>

    <div>
      <h2 class="graphs-container title-graphs">{{ $t('report.history_situation') }}</h2>
      <div class="form-group select-search col-xs-12 col-sm-12 col-md-12 col-lg-12">
        <label class="col-sm-2 control-label" for="textInput-markup" id="lbl-date-range">
          <p>{{dataPoints.to.format('DD MMM YYYY')}}</p> -
          <p>{{dataPoints.from.format('DD MMM YYYY')}}</p>
        </label>
        <div class="col-sm-4">
          <select v-on:change="getHistoryStats()" class="form-control" v-model="dateRangeSearchId">
            <option v-for="dateRange in dateRanges" v-bind:key="dateRange.value" v-bind:value="dateRange.value">{{dateRange.display}}</option>
          </select>
        </div>
      </div>

      <h2 class="section-title col-xs-12 col-sm-12 col-md-12 col-lg-12">{{ $t('report.averages') }}</h2>

      <div class="row no-margin averages">
        <div class="col-xs-12 col-sm-12 col-md-4 col-lg-4 border-bottom avg-cont">
          <h3>{{ $t('report.average_duration_user') }}</h3>
          <div v-if="loaders.avg_user_duration" class="spinner spinner-lg"></div>
          <h1 v-if="!loaders.avg_user_duration">{{ averages.avg_user_duration[0] | secondsInHour }}</h1>
        </div>
        <div class="col-xs-12 col-sm-12 col-md-4 col-lg-4 border-center border-bottom avg-cont">
          <h3>{{ $t('report.medium_traffic_user_download') }}</h3>
          <div v-if="loaders.avg_user_traffic" class="spinner spinner-lg"></div>
          <h1 v-if="!loaders.avg_user_traffic">{{ averages.avg_user_traffic[1] | byteFormat }}</h1>
        </div>
        <div class="col-xs-12 col-sm-12 col-md-4 col-lg-4 border-bottom avg-cont">
          <h3>{{ $t('report.medium_traffic_user_upload') }}</h3>
          <div v-if="loaders.avg_user_traffic" class="spinner spinner-lg"></div>
          <h1 v-if="!loaders.avg_user_traffic">{{ averages.avg_user_traffic[0] | byteFormat }}</h1>
        </div>
      </div>

      <div class="row no-margin averages">
        <div class="col-xs-12 col-sm-12 col-md-4 col-lg-4 avg-cont">
          <h3>{{ $t('report.average_duration_connections') }}</h3>
          <div v-if="loaders.avg_conn_duration" class="spinner spinner-lg"></div>
          <h1 v-if="!loaders.avg_conn_duration">{{ averages.avg_conn_duration[0] | secondsInHour }}</h1>
        </div>
        <div class="col-xs-12 col-sm-12 col-md-4 col-lg-4 border-center avg-cont">
          <h3>{{ $t('report.medium_traffic_connections_download') }}</h3>
          <div v-if="loaders.avg_conn_traffic" class="spinner spinner-lg"></div>
          <h1 v-if="!loaders.avg_conn_traffic">{{ averages.avg_conn_traffic[1] | byteFormat }}</h1>
        </div>
        <div class="col-xs-12 col-sm-12 col-md-4 col-lg-4 avg-cont">
          <h3>{{ $t('report.medium_traffic_connections_upload') }}</h3>
          <div v-if="loaders.avg_conn_traffic" class="spinner spinner-lg"></div>
          <h1 v-if="!loaders.avg_conn_traffic">{{ averages.avg_conn_traffic[0] | byteFormat }}</h1>
        </div>
      </div>

      <div class="row no-margin">
        <h2 class="section-title col-xs-12 col-sm-12 col-md-12 col-lg-12">{{ $t('report.connection_traffic') }}</h2>
        <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6">
          <div v-if="loaders.sessions" class="spinner spinner-lg"></div>
          <vue-chart v-show="!loaders.sessions" type="line" :options="charts.sessions.options" :data="charts.sessions"></vue-chart>
        </div>
        <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6">
          <div v-if="loaders.traffic" class="spinner spinner-lg"></div>
          <vue-chart v-show="!loaders.traffic" type="line" :options="charts.traffic.options" :data="charts.traffic"></vue-chart>
        </div>
      </div>

      <div class="row no-margin">
        <h2 class="section-title col-xs-12 col-sm-12 col-md-12 col-lg-12">{{ $t('report.user_statistics') }}</h2>
        <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6">
          <div v-if="loaders.avg_user_traffic" class="spinner spinner-lg"></div>
          <vue-chart v-show="!loaders.avg_user_traffic" type="line" :options="charts.avg_user_traffic.options" :data="charts.avg_user_traffic"></vue-chart>
        </div>
        <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6">
          <div v-if="loaders.avg_user_duration" class="spinner spinner-lg"></div>
          <vue-chart v-show="!loaders.avg_user_duration" type="bar" :options="charts.avg_user_duration.options" :data="charts.avg_user_duration"></vue-chart>
        </div>
      </div>

      <div class="row no-margin">
        <h2 class="section-title col-xs-12 col-sm-12 col-md-12 col-lg-12">{{ $t('report.connections_statistics') }}</h2>
        <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6">
          <div v-if="loaders.avg_conn_traffic" class="spinner spinner-lg"></div>
          <vue-chart v-show="!loaders.avg_conn_traffic" type="line" :options="charts.avg_conn_traffic.options" :data="charts.avg_conn_traffic"></vue-chart>
        </div>
        <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6">
          <div v-if="loaders.avg_conn_duration" class="spinner spinner-lg"></div>
          <vue-chart v-show="!loaders.avg_conn_duration" type="bar" :options="charts.avg_conn_duration.options" :data="charts.avg_conn_duration"></vue-chart>
        </div>
      </div>
    </div>

    <div class="graph-divider"></div>

    <div>
      <h2 class="graphs-container title-graphs">{{ $t('report.sms_reports') }}</h2>
      <div class="row no-margin">
        <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6">
          <div v-if="loaders.sms_year" class="spinner spinner-lg"></div>
          <vue-chart v-show="!loaders.sms_year" type="line" :options="charts.sms_year.options" :data="charts.sms_year"></vue-chart>
        </div>
        <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6">
          <div v-if="loaders.sms_history" class="spinner spinner-lg"></div>
          <vue-chart v-show="!loaders.sms_history" type="line" :options="charts.sms_history.options" :data="charts.sms_history"></vue-chart>
        </div>
      </div>
    </div>

  </div>
</template>

<script>
import StorageService from "../services/storage";
import HotspotService from "../services/hotspot";
import StatsService from "../services/stats";

import VueChart from "vue-chart-js";
import moment from "moment";

export default {
  name: "Reports",
  components: {
    VueChart
  },
  mixins: [StorageService, HotspotService, StatsService],
  data() {
    return {
      loaders: {
        current: true,
        sessions: true,
        traffic: true,
        avg_user_traffic: true,
        avg_user_duration: true,
        avg_conn_traffic: true,
        avg_conn_duration: true,
        sms_year: true,
        sms_history: true
      },
      msg: this.$i18n.t("report.reports"),
      user: this.get("loggedUser") || null,
      dataPoints: {
        to: moment()
          .subtract(7, "day")
          .startOf("day"),
        from: moment().startOf("day")
      },
      dateRanges: [
        {
          display: this.$i18n.t("report.last_7_days"),
          value: 7
        },
        {
          display: this.$i18n.t("report.last_15_days"),
          value: 15
        },
        {
          display: this.$i18n.t("report.last_month"),
          value: 30
        }
      ],
      averages: {
        current: [0, 0],
        avg_user_traffic: [0, 0],
        avg_user_duration: [0],
        avg_conn_traffic: [0, 0],
        avg_conn_duration: [0]
      },
      dateRangeSearchId: this.get("reports_date_range_id") || 7,
      hotspots: [],
      hotspotSearchId: 0,
      charts: {
        current: this.chartConfig("current"),
        sessions: this.chartConfig("sessions"),
        traffic: this.chartConfig("traffic"),
        avg_user_traffic: this.chartConfig("avg_user_traffic"),
        avg_user_duration: this.chartConfig("avg_user_duration"),
        avg_conn_traffic: this.chartConfig("avg_conn_traffic"),
        avg_conn_duration: this.chartConfig("avg_conn_duration"),
        sms_year: this.chartConfig("sms_year"),
        sms_history: this.chartConfig("sms_history")
      }
    };
  },
  mounted() {
    if (this.$route.params.hotspotId !== undefined) {
      this.hotspotSearchId = this.$route.params.hotspotId;
    } else {
      var hsId = this.get("selected_hotspot_id");
      this.hotspotSearchId = hsId;
    }

    var context = this;
    this.getAllHotspots(function() {
      context.getCurrentStats();
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

          var hsId = this.get("selected_hotspot_id") || this.hotspots[0].id;
          if (
            this.$parent.user.info.type == "customer" ||
            this.$parent.user.info.type == "desk"
          ) {
            hsId = this.$parent.user.info.hotspot_id;
          }
          this.hotspotSearchId = hsId;

          callback();
        },
        error => {
          console.error(error);
          callback();
        }
      );
    },
    getCurrentStats() {
      this.isChartLoading = true;

      // save filters
      this.set(
        "selected_hotspot_id",
        this.hotspotSearchId || this.get("selected_hotspot_id") || 0
      );

      this.getHistoryStats();
    },
    getHistoryStats() {
      this.isChartLoading = true;

      // save filters
      this.set(
        "reports_date_range_id",
        this.dateRangeSearchId || this.get("reports_date_range_id") || 1
      );

      switch (this.dateRangeSearchId) {
        case 7:
          this.dataPoints.to = moment().subtract(7, "days");
          break;

        case 15:
          this.dataPoints.to = moment().subtract(15, "days");
          break;

        case 30:
          this.dataPoints.to = moment().subtract(1, "month");
          break;
      }

      this.reportGraph("current");
      this.reportGraph("sessions");
      this.reportGraph("traffic");
      this.reportGraph("avg_user_traffic");
      this.reportGraph("avg_user_duration");
      this.reportGraph("avg_conn_traffic");
      this.reportGraph("avg_conn_duration");
      this.reportGraph("sms_year");
      this.reportGraph("sms_history");
    },
    reportGraph(graph) {
      this.loaders[graph] = true;
      this.reportsHistory(
        graph,
        this.hotspotSearchId,
        this.dateRangeSearchId,
        success => {
          this.charts[graph].labels = success.body.labels || [];

          for (var l in success.body.labels) {
            var label = success.body.labels[l];

            var indexSet0 =
              success.body.set0 &&
              success.body.set0.labels &&
              success.body.set0.labels.indexOf(label);

            if (indexSet0 == -1) {
              this.charts[graph].datasets[0].data[l] = 0;
            } else {
              if (indexSet0 != null)
                this.charts[graph].datasets[0].data[l] =
                  success.body.set0.data[indexSet0];
            }

            if (this.charts[graph].datasets[1]) {
              var indexSet1 =
                success.body.set1 &&
                success.body.set1.labels &&
                success.body.set1.labels.indexOf(label);

              if (indexSet1 == -1) {
                this.charts[graph].datasets[1].data[l] = 0;
              } else {
                if (indexSet1 != null)
                  this.charts[graph].datasets[1].data[l] =
                    success.body.set1.data[indexSet1];
              }
            }

            // averages
            this.averages[graph] = success.body.avg || [];
          }
          this.loaders[graph] = false;
        },
        error => {
          console.error(error);
          this.loaders[graph] = false;
        }
      );
    },
    chartConfig(chart) {
      var context = this;
      var config = {};
      switch (chart) {
        case "current":
          config = {
            labels: [],
            datasets: [
              {
                label: this.$i18n.t("report.today_logins"),
                data: [],
                fill: false,
                backgroundColor: "#00b9e4",
                borderColor: "#00b9e4"
              },
              {
                label: this.$i18n.t("report.today_connections"),
                data: [],
                steppedLine: false,
                fill: false,
                showLines: true,
                backgroundColor: "#703fec",
                borderColor: "#703fec"
              }
            ],
            options: {
              scales: {
                yAxes: [
                  {
                    ticks: {
                      maxTicksLimit: 5,
                      callback: function(item) {
                        if (item % 1 === 0) {
                          return item;
                        }
                      },
                      beginAtZero: true
                    },
                    gridLines: {
                      display: false,
                      drawBorder: false
                    }
                  }
                ]
              },
              title: {
                display: true,
                text: this.$i18n.t("report.connection_logins")
              },
              elements: {
                line: {
                  tension: 0
                }
              }
            }
          };
          break;

        case "sessions":
          config = {
            labels: [],
            datasets: [
              {
                label: this.$i18n.t("report.session"),
                data: [],
                backgroundColor: "#a30000",
                fill: false,
                borderColor: "#a30000"
              },
              {
                label: this.$i18n.t("report.user"),
                data: [],
                backgroundColor: "#ec7a08",
                fill: false,
                borderColor: "#ec7a08"
              }
            ],
            options: {
              elements: {
                line: {
                  tension: 0
                }
              },
              title: {
                display: true,
                text: this.$i18n.t("report.connection_newlogins")
              },
              scales: {
                xAxes: [
                  {
                    gridLines: {
                      zeroLineColor: "transparent"
                    }
                  }
                ],
                yAxes: [
                  {
                    ticks: {
                      maxTicksLimit: 5,
                      callback: function(item) {
                        if (item % 1 === 0) {
                          return item;
                        }
                      },
                      beginAtZero: true
                    },
                    gridLines: {
                      display: false,
                      drawBorder: false
                    }
                  }
                ]
              }
            }
          };
          break;

        case "traffic":
          config = {
            labels: [],
            datasets: [
              {
                label: this.$i18n.t("report.upload"),
                data: [],
                backgroundColor: "#f0ab00"
              },
              {
                label: this.$i18n.t("report.download"),
                data: [],
                backgroundColor: "#92d400"
              }
            ],
            options: {
              tooltips: {
                callbacks: {
                  label: function(item) {
                    const prettyBytes = require("pretty-bytes");
                    return prettyBytes(item.yLabel);
                  }
                }
              },
              title: {
                display: true,
                text: this.$i18n.t("report.total_traffic")
              },
              scales: {
                yAxes: [
                  {
                    scaleLabel: {
                      display: true,
                      labelString: this.$i18n.t("report.size")
                    },
                    ticks: {
                      maxTicksLimit: 5,
                      callback: function(item) {
                        const prettyBytes = require("pretty-bytes");
                        return prettyBytes(item);
                      },
                      beginAtZero: true
                    }
                  }
                ]
              }
            }
          };
          break;

        case "avg_user_traffic":
          config = {
            labels: [],
            datasets: [
              {
                label: this.$i18n.t("report.upload"),
                data: [],
                backgroundColor: "#3f9c35"
              },
              {
                label: this.$i18n.t("report.download"),
                data: [],
                backgroundColor: "#007a87"
              }
            ],
            options: {
              tooltips: {
                callbacks: {
                  label: function(item) {
                    const prettyBytes = require("pretty-bytes");
                    return prettyBytes(item.yLabel);
                  }
                }
              },
              title: {
                display: true,
                text: this.$i18n.t("report.medium_traffic_user")
              },
              scales: {
                yAxes: [
                  {
                    scaleLabel: {
                      display: true,
                      labelString: this.$i18n.t("report.size")
                    },
                    ticks: {
                      maxTicksLimit: 5,
                      callback: function(item) {
                        const prettyBytes = require("pretty-bytes");
                        return prettyBytes(item);
                      },
                      beginAtZero: true
                    }
                  }
                ]
              }
            }
          };
          break;
        case "avg_user_duration":
          config = {
            labels: [],
            datasets: [
              {
                label: this.$i18n.t("report.duration"),
                data: [],
                backgroundColor: "#4d5258"
              }
            ],
            options: {
              legend: {
                display: false
              },
              tooltips: {
                callbacks: {
                  label: function(item) {
                    return context.$options.filters.secondsInHour(item.yLabel);
                  }
                }
              },
              title: {
                display: true,
                text: this.$i18n.t("report.average_duration_user")
              },
              scales: {
                yAxes: [
                  {
                    scaleLabel: {
                      display: true,
                      labelString: this.$i18n.t("report.time")
                    },
                    ticks: {
                      maxTicksLimit: 5,
                      callback: function(item) {
                        return context.$options.filters.secondsInHour(item);
                      },
                      beginAtZero: true
                    }
                  }
                ]
              }
            }
          };
          break;

        case "avg_conn_traffic":
          config = {
            labels: [],
            datasets: [
              {
                label: this.$i18n.t("report.upload"),
                data: [],
                backgroundColor: "#3f9c35"
              },
              {
                label: this.$i18n.t("report.download"),
                data: [],
                backgroundColor: "#007a87"
              }
            ],
            options: {
              tooltips: {
                callbacks: {
                  label: function(item) {
                    const prettyBytes = require("pretty-bytes");
                    return prettyBytes(item.yLabel);
                  }
                }
              },
              title: {
                display: true,
                text: this.$i18n.t("report.medium_traffic_connections")
              },
              scales: {
                yAxes: [
                  {
                    scaleLabel: {
                      display: true,
                      labelString: this.$i18n.t("report.size")
                    },
                    ticks: {
                      maxTicksLimit: 5,
                      callback: function(item) {
                        const prettyBytes = require("pretty-bytes");
                        return prettyBytes(item);
                      },
                      beginAtZero: true
                    }
                  }
                ]
              }
            }
          };
          break;
        case "avg_conn_duration":
          config = {
            labels: [],
            datasets: [
              {
                label: this.$i18n.t("report.duration"),
                data: [],
                backgroundColor: "#4d5258"
              }
            ],
            options: {
              legend: {
                display: false
              },
              tooltips: {
                callbacks: {
                  label: function(item) {
                    return context.$options.filters.secondsInHour(item.yLabel);
                  }
                }
              },
              scales: {
                yAxes: [
                  {
                    scaleLabel: {
                      display: true,
                      labelString: this.$i18n.t("report.time")
                    },
                    ticks: {
                      maxTicksLimit: 5,
                      callback: function(item) {
                        return context.$options.filters.secondsInHour(item);
                      },
                      beginAtZero: true
                    }
                  }
                ]
              },
              title: {
                display: true,
                text: this.$i18n.t("report.average_duration_connections")
              }
            }
          };
          break;

        case "sms_year":
          config = {
            labels: [],
            datasets: [
              {
                label: this.$i18n.t("report.sms"),
                data: [],
                backgroundColor: "#8b0000",
                fill: false,
                borderColor: "#8b0000"
              }
            ],
            options: {
              elements: {
                line: {
                  tension: 0
                }
              },
              title: {
                display: true,
                text: this.$i18n.t("report.sent_this_year")
              },
              scales: {
                xAxes: [
                  {
                    gridLines: {
                      zeroLineColor: "transparent"
                    }
                  }
                ],
                yAxes: [
                  {
                    ticks: {
                      maxTicksLimit: 5,
                      callback: function(item) {
                        if (item % 1 === 0) {
                          return item;
                        }
                      },
                      beginAtZero: true
                    },
                    gridLines: {
                      display: false,
                      drawBorder: false
                    }
                  }
                ]
              }
            }
          };
          break;

        case "sms_history":
          config = {
            labels: [],
            datasets: [
              {
                label: this.$i18n.t("report.sms"),
                data: [],
                backgroundColor: "#ec7a08",
                fill: false,
                borderColor: "#ec7a08"
              }
            ],
            options: {
              elements: {
                line: {
                  tension: 0
                }
              },
              title: {
                display: true,
                text: this.$i18n.t("report.sent_this_period")
              },
              scales: {
                xAxes: [
                  {
                    gridLines: {
                      zeroLineColor: "transparent"
                    }
                  }
                ],
                yAxes: [
                  {
                    ticks: {
                      maxTicksLimit: 5,
                      callback: function(item) {
                        if (item % 1 === 0) {
                          return item;
                        }
                      },
                      beginAtZero: true
                    },
                    gridLines: {
                      display: false,
                      drawBorder: false
                    }
                  }
                ]
              }
            }
          };
          break;
      }
      return config;
    }
  }
};
</script>

<style scoped>
.control-label p {
  display: inline;
}

.no-margin {
  margin: 0px !important;
}

.graph-divider {
  border-top: 1px solid #d1d1d1;
}

.graph-divider-blank {
  margin-bottom: 20px;
}

.averages {
  text-align: center;
}

.border-center {
  border-left: 1px dashed;
  border-right: 1px dashed;
}

.border-bottom {
  border-bottom: 1px dashed;
}

.border-right {
  border-right: 1px dashed;
}

@media only screen and (max-width: 991px) {
  .border-bottom {
    border-bottom: 1px;
  }

  .border-center {
    border-left: 0px;
    border-right: 0px;
  }

  .border-right {
    border-right: 0px;
  }
}

.avg-cont {
  min-height: 94px;
}
</style>