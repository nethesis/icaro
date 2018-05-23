<template>
  <div class="row">
      <div class="row">
        <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12">
          <vue-chart type="line" :options="totalTrafficChart.options" :data="totalTrafficChart"></vue-chart>
        </div>
      </div>

  </div>
</template>

<script>
import VueChart from "vue-chart-js";
import moment from "moment";
import { extendMoment } from "moment-range";

import UtilService from "../../services/util";
import StorageService from "../../services/storage";
import StatsService from "../../services/stats";
import filters from "../../filters/filters";

export default {
  name: "ActualReport",
  components: {
    VueChart
  },
  props: {
    smsData: {
      type: Array
    },
    chartDateRange: {
      type: Array
    },
    chartLabels: {
      type: Array
    },
    isYear: {
      type: Boolean
    }
  },
  mixins: [UtilService, StatsService, StorageService],
  data() {
    return {
      totalTrafficChart: {
        labels: this.chartLabels,
        datasets: [
          {
            label: this.$i18n.t("report.sms"),
            data: [],
            backgroundColor: this.isYear ? "#8b0000" : "#ec7a08",
            fill: false,
            borderColor: this.isYear ? "#8b0000" : "#ec7a08"
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
            text: this.isYear
              ? this.$i18n.t("report.sent_this_year")
              : this.getFilterText()
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
      smsToShow: []
    };
  },
  created() {
    this.implementDataInChart(this.isYear);
  },
  methods: {
    getFilterText() {
      var range = this.get("reports_date_range_id");
      switch (range) {
        case 1:
          return this.$i18n.t("report.sent_last_7_days");
          break;
        case 2:
          return this.$i18n.t("report.sent_last_15_days");
          break;
        case 3:
          return this.$i18n.t("report.sent_last_month");
          break;
      }
    },
    implementDataInChart(isYear) {
      let newUsers = 0;
      this.chartDateRange.forEach(element => {
        this.smsToShow = this.smsData
          .map(function(item) {
            if (isYear) {
              return moment(element).month() == moment(item.sent).month();
            } else {
              return item.sent.substring(0, 10) === element;
            }
          })
          .filter(function(item) {
            return item == true;
          }).length;
        this.totalTrafficChart.datasets[0].data.push(this.smsToShow);
      });
    }
  }
};
</script>

<style scoped>
.page-header h1 {
  margin-top: 0px;
  font-size: 30px;
}

.actual-report-header {
  margin: 20px;
  text-align: center;
}

.actual-report-header h1 {
  font-size: 25px;
  font-weight: bold;
  margin-top: 5px;
}

.actual-report-header > div:nth-child(2) {
  border-left: 1px dashed gray;
  border-right: 1px dashed gray;
}

.actual-report-header h3 {
  max-height: 35px;
  height: 35px;
  margin-bottom: -5px;
}

@media only screen and (min-width: 1025px) {
  .actual-report-header h3 {
    font-size: 18px;
  }
}

@media only screen and (max-width: 925px) {
  .actual-report-header h3 {
    font-size: 15px;
  }
}

@media only screen and (max-width: 768px) {
  .actual-report-header h3 {
    font-size: 13px;
  }
}
</style>
