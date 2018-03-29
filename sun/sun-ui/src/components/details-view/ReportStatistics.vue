<template>
    <div class="row">
        <div class="col-sm-10 section-title">
           <h1>{{ $t('report.total_traffic') }}</h1>
        </div>
        <div class="col-sm-11">
            <vue-chart type="bar" :options="trafficChartData.options" :data="trafficChartData"></vue-chart>
        </div>
        <div class="col-sm-10 section-title">
          <h1>{{ $t('report.user_statistics') }}</h1>
        </div>
        <div class="col-sm-12">
          <div class="row">
              <div class="col-sm-6">
                <vue-chart type="bar"  :options="medianUserChart.options" :data="medianUserChart"></vue-chart>
              </div>
              <div class="col-sm-6">
                <vue-chart type="bar"  :options="avgUserChart.options" :data="avgUserChart"></vue-chart>
              </div>
          </div>
        </div>
        <div class="col-sm-10 section-title">
           <h1>{{ $t('report.connections_statistics') }}</h1>
        </div>
        <div class="col-sm-12">
            <div class="row">
              <div class="col-sm-6">
                <vue-chart type="bar" class="test_chart" :options="medianSessionChart.options" :data="medianSessionChart"></vue-chart>
              </div>
              <div class="col-sm-6">
                <vue-chart type="bar"  :options="avgSessionChart.options" :data="avgSessionChart"></vue-chart>
              </div>
            </div>
        </div>
    </div>

</template>
<script>
import VueChart from "vue-chart-js";

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
            backgroundColor: "#008ae6"
          },
          {
            label: this.$i18n.t("report.download"),
            data: [],
            backgroundColor: "#66d9ff"
          }
        ],
        options: {
          title: {
            display: true,
            text: 'Mbyte'
          },
          scales: {
            yAxes: {
              ticks: {
                beginAtZero: true
              }
            }
          }
        }
      },
      avgUserChart: {
        labels: this.chartLabels,
        datasets: [
          {
            label: this.$i18n.t("report.upload"),
            data: [],
            backgroundColor: "#579f14"
          },
          {
            label: this.$i18n.t("report.download"),
            data: [],
            backgroundColor: "#444"
          }
        ],
        options: {
          title: {
            display: true,
            text: this.$i18n.t("report.average_duration_user")
          },
          scales: {
            yAxes: {
              ticks: {
                beginAtZero: true
              }
            }
          }
        }
      },
      medianUserChart: {
        labels: this.chartLabels,
        datasets: [
          {
            label: this.$i18n.t("report.upload"),
            data: [],
            backgroundColor: "#35a3d7"
          },
          {
            label: this.$i18n.t("report.download"),
            data: [],
            backgroundColor: "#66b83c"
          }
        ],
        options: {
          title: {
            display: true,
            text: this.$i18n.t("report.medium_traffic_user")
          },
          scales: {
            yAxes: {
              ticks: {
                beginAtZero: true
              }
            }
          }
        }
      },
      medianSessionChart: {
        labels: this.chartLabels,
        datasets: [
          {
            label: this.$i18n.t("report.upload"),
            data: [],
            backgroundColor: "#35a3d7"
          },
          {
            label: this.$i18n.t("report.download"),
            data: [],
            backgroundColor: "#66b83c"
          }
        ],
        options: {
          title: {
            display: true,
            text: this.$i18n.t("report.medium_traffic_connections")
          },
          scales: {
            yAxes: {
              ticks: {
                beginAtZero: true
              }
            }
          }
        }
      },
      avgSessionChart: {
        labels: this.chartLabels,
       datasets: [
          {
            label: this.$i18n.t("report.upload"),
            data: [],
            backgroundColor: "#579f14"
          },
          {
            label: this.$i18n.t("report.download"),
            data: [],
            backgroundColor: "#444"
          }
        ],
        options: {
          title: {
            display: true,
            text: this.$i18n.t("report.average_duration_connections")
          },
          scales: {
            yAxes: {
              ticks: {
                beginAtZero: true
              }
            }
          }
        }
      }
    };
  },
  created() {
    this.fillTotalTrafficChart();
    this.implementUserChart();
    this.implementSessionChart();
  },
  methods: {
    calculateAVG(input) {
      let sum = 0;
      for (let i = 0; i < input.length; i++) {
        sum += input[i];
      }
      return sum / input.length;
    },
    calculateMEDIUM(input) {
      let sortedArray = input.sort();

      let half = Math.floor(sortedArray.length / 2);

      if (sortedArray.length % 2) {
        return sortedArray[half];
      } else {
        return (sortedArray[half - 1] + sortedArray[half]) / 2.0;
      }
    },
    fillTotalTrafficChart() {
      this.chartDateRange.forEach(date => {
        let kbps_up = 0;
        let kbps_down = 0;
        this.newUsersReport.map(function(user) {
          if (user.valid_from.substring(0, 10) === date) {
            kbps_up += user.kbps_up/1000;
            kbps_down += user.kbps_down/1000;
          } else {
            kbps_up += 0;
            kbps_down += 0;
          }
        });
        this.trafficChartData.datasets[0].data.push(kbps_up);
        this.trafficChartData.datasets[1].data.push(kbps_down);
      });
    },
    implementUserChart() {
      this.chartDateRange.forEach(date => {
        let index = 0;
        let kbps_up = [0];
        let kbps_down = [0];
        this.newUsersReport.map(function(user) {
          if (user.valid_from.substring(0, 10) === date) {
            kbps_up[index] += user.kbps_up/1000;
            kbps_down[index] += user.kbps_down/1000;
            kbps_up.push(0);
            kbps_down.push(0);
            index++;
          } else {
            kbps_up[index] += 0;
            kbps_down[index] += 0;
          }
        });
        if (kbps_up.length > 1) {
          kbps_up.pop();
        }
        if (kbps_down.length > 1) {
          kbps_down.pop();
        }

        this.medianUserChart.datasets[0].data.push(
          this.calculateMEDIUM(kbps_up)
        );
        this.medianUserChart.datasets[1].data.push(
          this.calculateMEDIUM(kbps_down)
        );

        this.avgUserChart.datasets[0].data.push(this.calculateAVG(kbps_up));
        this.avgUserChart.datasets[1].data.push(this.calculateAVG(kbps_down));
      });
    },
    implementSessionChart() {
      this.chartDateRange.forEach(date => {
        let index = 0;
        let kbps_up = [0];
        let kbps_down = [0];
        this.sessionsReport.map(function(session) {
          if (session.start_time.substring(0, 10) === date) {
            kbps_up[index] += session.bytes_up;
            kbps_down[index] += session.bytes_down;
            kbps_up.push(0);
            kbps_down.push(0);
            index++;
          } else {
            kbps_up[index] += 0;
            kbps_down[index] += 0;
          }
        });
        if (kbps_up.length > 1) {
          kbps_up.pop();
        }
        if (kbps_down.length > 1) {
          kbps_down.pop();
        }

        this.medianSessionChart.datasets[0].data.push(
          this.calculateMEDIUM(kbps_up)
        );
        this.medianSessionChart.datasets[1].data.push(
          this.calculateMEDIUM(kbps_down)
        );

        this.avgSessionChart.datasets[0].data.push(this.calculateAVG(kbps_up));
        this.avgSessionChart.datasets[1].data.push(
          this.calculateAVG(kbps_down)
        );
      });
    }
  }
};
</script>
<style scoped>
.section-title {
  margin: 30px 0px 10px 20px;
} 
</style>
