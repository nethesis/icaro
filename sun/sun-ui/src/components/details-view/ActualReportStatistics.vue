<template>
    <div class="row">
        <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12 page-header">
            <h1>{{ $t('report.current_situation') }}</h1>
        </div>
        <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12 actual-report-header">
            <div class="col-xs-4 col-sm-4 col-md-4 col-lg-4">
                <h3>{{ $t('report.users_currently_connected') }}</h3>
                <h1>{{ userCurrentlyLogin }}</h1>
            </div>
            <div class="col-xs-4 col-sm-4 col-md-4 col-lg-4">
                <h3>{{ $t('report.total_daily_connections') }}</h3>
                <h1>{{ totalDailyConnection }}</h1>
            </div>
            <div class="col-xs-4 col-sm-4 col-md-4 col-lg-4">
                <h3>{{ $t('report.new_daily_logins') }}</h3>
                <h1>{{ newDailyLogins }}</h1>
            </div>
        </div>
        <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12s">
            <vue-chart :height="100" type="line" :options="actualStatisticChart.options" :data="actualStatisticChart"></vue-chart>
        </div>
    </div>
</template>

<script>
    import VueChart from "vue-chart-js";
    import moment from "moment";
    import {
        extendMoment
    } from "moment-range";
    import UtilService from "../../services/util";
    export default {
        name: "ActualReport",
        components: {
            VueChart
        },
        props: {
            todayConnections: {
                type: Array
            }
        },
        mixins: [UtilService],
        data() {
            return {
                idsArray: [],
                newLogins: [],
                userCurrentlyLogin: 0,
                totalDailyConnection: 0,
                newDailyLogins: 0,
                dateRange: {
                    from: null,
                    to: null
                },
                actualStatisticChart: {
                    labels: [],
                    datasets: [{
                            label: this.$i18n.t('report.today_logins'),
                            data: [],
                            fill: false,
                            backgroundColor: "#6fbc45",
                            borderColor: "#6fbc45"
                        },
                        {
                            label: this.$i18n.t('report.today_connections'),
                            data: [],
                            steppedLine: false,
                            fill: false,
                            showLines: true,
                            backgroundColor: "#40b1dd",
                            borderColor: "#40b1dd"
                        }
                    ],
                    options: {
                        scales: {
                            yAxes: [{
                                ticks: {
                                    maxTicksLimit: 5,
                                    callback: function(item) {
                                        if (item % 1 === 0) {
                                            return item;
                                        }
                                    },
                                    beginAtZero: true,
                                },
                                gridLines: {
                                    display: false,
                                    drawBorder: false
                                }
                            }]
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
                }
            };
        },
        created() {
            this.configureChart();
            this.implementToChart();
        },
        methods: {
            configureChart() {
                const momentRange = extendMoment(moment);
                this.dateRange.from = moment().startOf("day");
                this.dateRange.to = moment();
                
                this.actualStatisticChart.labels = Array.from(
                    momentRange.range(this.dateRange.from, this.dateRange.to).by("hours")
                ).map(function(date) {
                    if (date._d.getHours().toString().length === 1) {
                        return "0" + date._d.getHours().toString();
                    } else {
                        return date._d.getHours().toString();
                    }
                });
            },
            implementToChart() {
                this.todayConnections.forEach(connetion => {
                    if (this.idsArray.indexOf(connetion.user_id) === -1)
                        this.idsArray.push(connetion.user_id);
                })
                this.idsArray.forEach(id => {
                    let earlierLogin = "24";
                    let newLogins = [0];
                    this.todayConnections.map(function(connection) {
                        if (id === connection.user_id) {
                            if (connection.start_time.substring(11, 13) < earlierLogin) {
                                earlierLogin = connection.start_time.substring(11, 13);
                                newLogins[0] = connection;
                            }
                        }
                    })
                    this.newLogins.push(newLogins[0]);
                })
                this.actualStatisticChart.labels.forEach(date => {
                    let newLogins = 0;
                    let newLoginCounter = 0;
                    let connections = 0;
                    let userConnected = 0;
                    // Insert newLogins on chart
                    this.todayConnections.map(function(connection) {
                        if (connection.start_time.substring(11, 13) === date) {
                            connections++;
                            if (connection.stop_time === null) {
                                userConnected++;
                            }
                        }
                    });
                    this.newLogins.map(function(connection) {
                        if (connection.start_time.substring(11, 13) === date) {
                            newLogins++;
                        }
                    })
                    this.totalDailyConnection += connections;
                    this.userCurrentlyLogin += userConnected;
                    this.newDailyLogins += newLogins;
                    this.actualStatisticChart.datasets[0].data.push(newLogins);
                    this.actualStatisticChart.datasets[1].data.push(connections);
                });
            }
        }
    };
</script>

<style scoped>
    .page-header h1 {
        margin-top: 50px;
        font-size: 30px;
    }
    .actual-report-header {
        margin: 20px;
        text-align: center;
    }
    .actual-report-header h1 {
        font-size: 40px;
        font-weight: bold;
    }
    .actual-report-header>div:nth-child(2) {
        border-left: 1px dashed gray;
        border-right: 1px dashed gray;
    }
    .actual-report-header h3{
        max-height: 35px;
        height: 35px;
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
