<template>
  <div class="row">
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12">

      <div class="row">
        <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12">
          <h2 class="section-title">{{ $t('report.connection_traffic') }}</h2>
        </div>
        <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6">
          <vue-chart type="line" :options="totalTrafficChart.options" :data="totalTrafficChart"></vue-chart>
        </div>
        <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6">
          <vue-chart type="bar" :options="trafficChartData.options" :data="trafficChartData"></vue-chart>
        </div>
      </div>

      <div class="row average-table">
        <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12 first-row">
          <div class="col-xs-12 col-sm-4 col-md-4 col-lg-4">
            <h3>{{ $t('report.average_duration_connections') }}</h3>
            <h1>{{ avgDurationConnection }}</h1>
          </div>
          <div class="col-xs-12 col-sm-4 col-md-4 col-lg-4">
            <h3>{{ $t('report.medium_traffic_connections_download') }}</h3>
            <h1>{{ avgDownloadTrafficConnection }}</h1>
          </div>
          <div class="col-xs-12 col-sm-4 col-md-4 col-lg-4">
            <h3>{{ $t('report.medium_traffic_connections_upload') }}</h3>
            <h1>{{ avgUploadTrafficConnection }}</h1>
          </div>
        </div>
        <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12 second-row">
          <div class="col-xs-12 col-sm-4 col-md-4 col-lg-4">
            <h3>{{ $t('report.average_duration_user') }}</h3>
            <h1>{{ avgDurationUser }}</h1>
          </div>
          <div class="col-xs-12 col-sm-4 col-md-4 col-lg-4">
            <h3>{{ $t('report.medium_traffic_user_download') }}</h3>
            <h1>{{ avgDownloadTrafficUser }}</h1>
          </div>
          <div class="col-xs-12 col-sm-4 col-md-4 col-lg-4">
            <h3>{{ $t('report.medium_traffic_user_upload') }}</h3>
            <h1>{{ avgUploadTrafficUser }}</h1>
          </div>
        </div>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <div class="row">
        <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12">
          <h2 class="section-title">{{ $t('report.user_statistics') }}</h2>
        </div>
        <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6">
          <vue-chart type="bar" :options="avgTrafficUserChart.options" :data="avgTrafficUserChart"></vue-chart>
        </div>
        <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6">
          <vue-chart type="bar" :options="avgDurationUserChart.options" :data="avgDurationUserChart"></vue-chart>
        </div>
      </div>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <h2 class="section-title last-title">{{ $t('report.connections_statistics') }}</h2>
    </div>
    <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12">
      <div class="row">
        <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6">
          <vue-chart type="bar" :options="avgTrafficSessionChart.options" :data="avgTrafficSessionChart"></vue-chart>
        </div>
        <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6">
          <vue-chart type="bar" :options="avgDurationSessionChart.options" :data="avgDurationSessionChart"></vue-chart>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import VueChart from "vue-chart-js";
import moment from "moment";
import { extendMoment } from "moment-range";
export default {
  props: {
    chartLabels: {
      type: Array
    },
    chartDateRange: {
      type: Array
    },
    newUsersReport: {
      type: Array
    },
    sessionsReport: {
      type: Array
    }
  },
  watch: {
    newUsersReport: function(val) {
      this.newUsersReport = val;
      this.getNewUserSessions();
      this.implementDataInChart();
      this.fillTotalTrafficChart();
      this.implementUserChart();
      this.implementSessionChart();
      this.implementAvgTable();
    },
    sessionsReport: function(val) {
      this.sessionsReport = val;
      this.getNewUserSessions();
      this.implementDataInChart();
      this.fillTotalTrafficChart();
      this.implementUserChart();
      this.implementSessionChart();
      this.implementAvgTable();
    }
  },
  components: {
    VueChart
  },
  data() {
    return {
      trafficChartData: {
        labels: this.chartLabels,
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
      },
      avgDurationUserChart: {
        labels: this.chartLabels,
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
                let hours = parseInt(Math.floor(item.yLabel / 3600));
                let minutes = parseInt(
                  Math.floor((item.yLabel - hours * 3600) / 60)
                );
                let seconds = parseInt(
                  (item.yLabel - (hours * 3600 + minutes * 60)) % 60
                );

                let dHours = hours > 9 ? hours : "0" + hours;
                let dMins = minutes > 9 ? minutes : "0" + minutes;
                let dSecs = seconds > 9 ? seconds : "0" + seconds;

                if (dHours === "00" && dMins !== "00") {
                  return dMins + "m " + dSecs + "s";
                } else if (dHours === "00" && dMins === "00") {
                  return dSecs + "s";
                } else {
                  return dHours + "h " + dMins + "m " + dSecs + "s";
                }
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
                    let hours = parseInt(Math.floor(item / 3600));
                    let minutes = parseInt(
                      Math.floor((item - hours * 3600) / 60)
                    );
                    let seconds = parseInt(
                      (item - (hours * 3600 + minutes * 60)) % 60
                    );
                    let dHours = hours > 9 ? hours : "0" + hours;
                    let dMins = minutes > 9 ? minutes : "0" + minutes;
                    let dSecs = seconds > 9 ? seconds : "0" + seconds;
                    if (dHours === "00") {
                      return dMins + "m ";
                    } else {
                      return dHours + "h " + dMins + "m ";
                    }
                  },
                  beginAtZero: true
                }
              }
            ]
          }
        }
      },
      avgTrafficUserChart: {
        labels: this.chartLabels,
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
      },
      avgDurationSessionChart: {
        labels: this.chartLabels,
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
                let hours = parseInt(Math.floor(item.yLabel / 3600));
                let minutes = parseInt(
                  Math.floor((item.yLabel - hours * 3600) / 60)
                );
                let seconds = parseInt(
                  (item.yLabel - (hours * 3600 + minutes * 60)) % 60
                );

                let dHours = hours > 9 ? hours : "0" + hours;
                let dMins = minutes > 9 ? minutes : "0" + minutes;
                let dSecs = seconds > 9 ? seconds : "0" + seconds;

                if (dHours === "00" && dMins !== "00") {
                  return dMins + "m " + dSecs + "s";
                } else if (dHours === "00" && dMins === "00") {
                  return dSecs + "s";
                } else {
                  return dHours + "h " + dMins + "m " + dSecs + "s";
                }
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
                    let hours = parseInt(Math.floor(item / 3600));
                    let minutes = parseInt(
                      Math.floor((item - hours * 3600) / 60)
                    );
                    let seconds = parseInt(
                      (item - (hours * 3600 + minutes * 60)) % 60
                    );

                    let dHours = hours > 9 ? hours : "0" + hours;
                    let dMins = minutes > 9 ? minutes : "0" + minutes;
                    let dSecs = seconds > 9 ? seconds : "0" + seconds;

                    if (dHours === "00") {
                      return dMins + "m ";
                    } else {
                      return dHours + "h " + dMins + "m ";
                    }
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
      },
      avgTrafficSessionChart: {
        labels: this.chartLabels,
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
      },
      totalTrafficChart: {
        labels: this.chartLabels,
        datasets: [
          {
            label: this.$i18n.t("report.session"),
            data: [],
            backgroundColor: "#a30000",
            fill: false,
            borderColor: "#a30000"
          },
          {
            label: this.$i18n.t("report.new_user"),
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
      },
      sessionToCalculate: [],
      sessionToShow: [],
      newUsersSession: [],
      avgDurationConnection: 0,
      avgDownloadTrafficConnection: 0,
      avgUploadTrafficConnection: 0,
      avgDurationUser: 0,
      avgDownloadTrafficUser: 0,
      avgUploadTrafficUser: 0
    };
  },
  methods: {
    calculateAVG(input) {
      let sum = 0;
      for (let i = 0; i < input.length; i++) {
        sum += input[i];
      }
      return sum / input.length;
    },
    secondsToHR(item) {
      let hours = parseInt(Math.floor(item / 3600));
      let minutes = parseInt(Math.floor((item - hours * 3600) / 60));
      let seconds = parseInt((item - (hours * 3600 + minutes * 60)) % 60);
      let dHours = hours > 9 ? hours : "0" + hours;
      let dMins = minutes > 9 ? minutes : "0" + minutes;
      let dSecs = seconds > 9 ? seconds : "0" + seconds;
      if (dHours === "00" && dMins !== "00") {
        return dMins + "m " + dSecs + "s";
      } else if (dHours === "00" && dMins === "00") {
        return dSecs + "s";
      } else {
        return dHours + "h " + dMins + "m " + dSecs + "s";
      }
    },
    removeDuplicates(myArr, prop) {
      return myArr.filter((obj, pos, arr) => {
        return arr.map(mapObj => mapObj[prop]).indexOf(obj[prop]) === pos;
      });
    },
    getNewUserSessions() {
      this.sessionsReport.forEach(element => {
        this.sessionToCalculate.push(Object.assign({}, element));
      });

      let newUsers = [];
      let id = [];
      let date = [];
      let isInserted = false;
      this.sessionToCalculate.map(function(session) {
        if (newUsers.length === 0) {
          newUsers.push(session);
          id.push(session.user_id);
          date.push(session.start_time.substring(0, 10) + session.user_id);
        } else {
          isInserted = false;
          newUsers.forEach(user => {
            if (!isInserted) {
              if (
                session.user_id === user.user_id &&
                user.start_time.substring(0, 10) ===
                  session.start_time.substring(0, 10)
              ) {
                user.bytes_up += session.bytes_up;
                user.bytes_down += session.bytes_down;
                user.duration += session.duration;
                isInserted = true;
              } else {
                if (id.includes(session.user_id) === false) {
                  newUsers.push(session);
                  id.push(session.user_id);
                  date.push(
                    session.start_time.substring(0, 10) + session.user_id
                  );
                  isInserted = true;
                } else if (
                  date.includes(
                    session.start_time.substring(0, 10) + session.user_id
                  ) === false
                ) {
                  newUsers.push(session);
                  date.push(
                    session.start_time.substring(0, 10) + session.user_id
                  );
                  isInserted = true;
                }
              }
            }
          });
        }
      });

      this.newUsersSession = newUsers;
    },
    // Implement Session and newUser between date range in Chart
    implementDataInChart() {
      let newUsers = 0;
      this.chartDateRange.forEach(element => {
        this.sessionToShow = this.sessionsReport
          .map(function(item) {
            return item.start_time.substring(0, 10) === element;
          })
          .filter(function(item) {
            return item == true;
          }).length;

        this.totalTrafficChart.datasets[0].data.push(this.sessionToShow);
        this.$set(this.totalTrafficChart.datasets[0], 0, {
          labels: this.chartLabels
        });
      });
      this.chartDateRange.forEach(element => {
        newUsers = 0;
        this.newUsersSession.map(function(session) {
          if (session.start_time.substring(0, 10) === element) newUsers++;
        });
        this.totalTrafficChart.datasets[1].data.push(newUsers);
        this.$set(this.totalTrafficChart.datasets[1], 0, {
          labels: this.chartLabels
        });
      });
    },
    fillTotalTrafficChart() {
      this.chartDateRange.forEach(date => {
        let bpsUp = 0;
        let bpsDown = 0;
        this.sessionsReport.map(function(session) {
          // check if user session start is on chart date
          if (session.start_time.substring(0, 10) === date) {
            // if user session is on date, check if session has been stopped on that day
            if (session.stop_time.substring(0, 10) === date) {
              bpsUp += session.bytes_up;
              bpsDown += session.bytes_down;
              // or session has been stopped on next day
            } else if (session.stop_time.substring(0, 10) > date) {
              bpsUp += session.bytes_up;
              bpsDown += session.bytes_down;
            }
            // check if session has started one day before, but has stopped on actual day
          } else if (
            session.stop_time.substring(0, 10) === date &&
            session.start_time.substring(0, 10) < date
          ) {
            bpsUp += session.bytes_up;
            bpsDown += session.bytes_down;
          } else {
            bpsUp += 0;
            bpsDown += 0;
          }
        });
        this.trafficChartData.datasets[0].data.push(bpsUp);
        this.trafficChartData.datasets[1].data.push(bpsDown);
        this.$set(this.trafficChartData.datasets[0], 0, {
          labels: this.chartLabels
        });
        this.$set(this.trafficChartData.datasets[1], 0, {
          labels: this.chartLabels
        });
      });
    },

    implementUserChart() {
      this.chartDateRange.forEach(date => {
        let bpsUpArray = [0];
        let bpsDownArray = [0];
        let durationArray = [0];
        let sessionCounter = 0;
        this.newUsersSession.map(function(session) {
          if (session.start_time.substring(0, 10) === date) {
            bpsUpArray[sessionCounter] += session.bytes_up;
            bpsDownArray[sessionCounter] += session.bytes_down;
            durationArray[sessionCounter] += session.duration;
            bpsUpArray.push(0);
            bpsDownArray.push(0);
            durationArray.push(0);
            sessionCounter++;
          } else {
            bpsUpArray[sessionCounter] += 0;
            bpsDownArray[sessionCounter] += 0;
            durationArray[sessionCounter] += 0;
          }
        });

        if (durationArray.length > 1) {
          durationArray.pop();
        }
        if (bpsUpArray.length > 1) {
          bpsUpArray.pop();
        }
        if (bpsDownArray.length > 1) {
          bpsDownArray.pop();
        }
        this.avgTrafficUserChart.datasets[0].data.push(
          this.calculateAVG(bpsUpArray)
        );
        this.avgTrafficUserChart.datasets[1].data.push(
          this.calculateAVG(bpsDownArray)
        );
        this.avgDurationUserChart.datasets[0].data.push(
          this.calculateAVG(durationArray)
        );
        this.$set(this.avgTrafficUserChart.datasets[0], 0, {
          labels: this.chartLabels
        });
        this.$set(this.avgTrafficUserChart.datasets[1], 0, {
          labels: this.chartLabels
        });
        this.$set(this.avgDurationUserChart.datasets[0], 0, {
          labels: this.chartLabels
        });
      });
    },
    implementSessionChart() {
      this.chartDateRange.forEach(date => {
        let index = 0;
        let duration = [0];
        let bpsUp = [0];
        let bpsDown = [0];
        this.sessionsReport.map(function(session) {
          // Check if session has been started between chart date
          if (session.start_time.substring(0, 10) === date) {
            bpsUp[index] += session.bytes_up;
            bpsDown[index] += session.bytes_down;
            duration[index] += session.duration;
            bpsUp.push(0);
            bpsDown.push(0);
            duration.push(0);
            index++;
          } else {
            bpsUp[index] += 0;
            bpsDown[index] += 0;
            duration[index] += 0;
          }
        });
        if (duration.length > 1) {
          duration.pop();
        }
        if (bpsUp.length > 1) {
          bpsUp.pop();
        }
        if (bpsDown.length > 1) {
          bpsDown.pop();
        }
        this.avgTrafficSessionChart.datasets[0].data.push(
          this.calculateAVG(bpsUp)
        );
        this.avgTrafficSessionChart.datasets[1].data.push(
          this.calculateAVG(bpsDown)
        );
        this.avgDurationSessionChart.datasets[0].data.push(
          this.calculateAVG(duration)
        );
        this.$set(this.avgTrafficSessionChart.datasets[0], 0, {
          labels: this.chartLabels
        });
        this.$set(this.avgTrafficSessionChart.datasets[1], 0, {
          labels: this.chartLabels
        });
        this.$set(this.avgDurationSessionChart.datasets[0], 0, {
          labels: this.chartLabels
        });
      });
    },
    implementAvgTable() {
      const prettyBytes = require("pretty-bytes");
      this.avgDurationConnection = this.secondsToHR(
        this.calculateAVG(this.avgDurationSessionChart.datasets[0].data)
      );
      this.avgUploadTrafficConnection = prettyBytes(
        this.calculateAVG(this.avgTrafficSessionChart.datasets[0].data)
      );
      this.avgDownloadTrafficConnection = prettyBytes(
        this.calculateAVG(this.avgTrafficSessionChart.datasets[1].data)
      );
      this.avgDurationUser = this.secondsToHR(
        this.calculateAVG(this.avgDurationUserChart.datasets[0].data)
      );
      this.avgUploadTrafficUser = prettyBytes(
        this.calculateAVG(this.avgTrafficUserChart.datasets[0].data)
      );
      this.avgDownloadTrafficUser = prettyBytes(
        this.calculateAVG(this.avgTrafficUserChart.datasets[1].data)
      );

      this.$set(this.avgDurationSessionChart.datasets[0], 0, {
        labels: this.chartLabels
      });
      this.$set(this.avgTrafficSessionChart.datasets[0], 0, {
        labels: this.chartLabels
      });
      this.$set(this.avgTrafficSessionChart.datasets[1], 0, {
        labels: this.chartLabels
      });
      this.$set(this.avgDurationUserChart.datasets[0], 0, {
        labels: this.chartLabels
      });
      this.$set(this.avgTrafficUserChart.datasets[0], 0, {
        labels: this.chartLabels
      });
      this.$set(this.avgTrafficUserChart.datasets[1], 0, {
        labels: this.chartLabels
      });
    }
  }
};
</script>

<style scoped>
.section-title {
  margin-bottom: -10px;
  margin-top: 30px;
}

.last-title {
  margin-top: 60px;
}

.total-traffic-chart {
  margin-top: -15px;
}

.page-header h1 {
  margin-top: 50px;
  margin-bottom: -10px;
  font-size: 30px;
}

.average-table {
  margin: 20px;
  margin-top: 60px;
}

.average-table h1 {
  font-size: 25px;
  font-weight: bold;
  margin-top: 5px;
}

.average-table h3 {
  font-size: 18px;
}

.first-row > div:first-child,
.first-row > div:nth-child(2),
.second-row > div:first-child,
.second-row > div:nth-child(2) {
  border-right: 1px dashed gray;
}

.second-row {
  border-top: 1px dashed gray;
}

.first-row > div,
.second-row > div {
  min-height: 110px;
  text-align: center;
}

.first-row > div h3,
.second-row > div h3 {
  max-height: 35px;
  height: 35px;
  font-size: 18px;
}

@media only screen and (max-width: 1200px) {
  .first-row > div h3,
  .second-row > div h3 {
    font-size: 17px;
  }
  .first-row > div h1,
  .second-row > div h1 {
    font-size: 32px;
  }
}

@media only screen and (max-width: 1030px) {
  .first-row > div h3,
  .second-row > div h3 {
    font-size: 17px;
  }
  .first-row > div h1,
  .second-row > div h1 {
    font-size: 26px;
  }
}

@media only screen and (max-width: 925px) {
  .first-row > div h3,
  .second-row > div h3 {
    font-size: 14px;
  }
  .first-row > div h1,
  .second-row > div h1 {
    font-size: 20px;
  }
  .first-row > div h3,
  .second-row > div h3 {
    max-height: 28px;
    height: 28px;
  }
}

@media only screen and (max-width: 768px) {
  .first-row > div h3,
  .second-row > div h3 {
    font-size: 18px;
  }
  .first-row > div h1,
  .second-row > div h1 {
    font-size: 35px;
  }
  .first-row > div:first-child,
  .first-row > div:nth-child(2),
  .second-row > div:first- .first-row > div h3,
  .second-row > div h3 {
    max-height: 25px;
    height: 25px;
  }
  .first-row > div:first-child,
  .first-row > div:nth-child(2),
  .second-row > div:first-child,
  .second-row > div:nth-child(2) {
    border-right: 0;
  }
  .second-row {
    border-top: 0;
  }
}
</style>
