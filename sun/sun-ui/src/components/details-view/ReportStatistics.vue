<template>
    <div class="row">
        <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12 page-header">
            <h1>{{ $t('report.history_situation') }}</h1>
        </div>
        <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12">
            <div class="row">
                <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12 section-title">
                    <h2>{{ $t('report.total_traffic') }}</h2>
                </div>
                <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12">
                    <vue-chart type="bar" :height="100" :options="trafficChartData.options" :data="trafficChartData"></vue-chart>
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
                <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12 section-title">
                    <h2>{{ $t('report.user_statistics') }}</h2>
                </div>
                <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6">
                    <vue-chart type="bar" :options="avgTrafficUserChart.options" :data="avgTrafficUserChart"></vue-chart>
                </div>
                <div class="col-xs-12 col-sm-12 col-md-6 col-lg-6">
                    <vue-chart type="bar" :options="avgDurationUserChart.options" :data="avgDurationUserChart"></vue-chart>
                </div>
            </div>
        </div>
        <div class="col-xs-12 col-sm-12 col-md-12 col-lg-12 section-title">
            <h2>{{ $t('report.connections_statistics') }}</h2>
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
    import VueChart from 'vue-chart-js'
    import moment from 'moment'
    import {
        extendMoment
    } from 'moment-range'
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
                    datasets: [{
                            label: this.$i18n.t('report.upload'),
                            data: [],
                            backgroundColor: '#008ae6'
                        },
                        {
                            label: this.$i18n.t('report.download'),
                            data: [],
                            backgroundColor: '#66d9ff'
                        }
                    ],
                    options: {
                        tooltips: {
                            callbacks: {
                                label: function(item) {
                                    const prettyBytes = require('pretty-bytes');
                                    return prettyBytes(item.yLabel);
                                }
                            }
                        },
                        title: {
                            display: true
                        },
                        scales: {
                            yAxes: [{
                                scaleLabel: {
                                    display: true,
                                    labelString: this.$i18n.t('report.size')
                                },
                                ticks: {
                                    maxTicksLimit: 5,
                                    callback: function(item) {
                                        const prettyBytes = require('pretty-bytes');
                                        return prettyBytes(item);
                                    },
                                    beginAtZero: true
                                }
                            }]
                        }
                    }
                },
                avgDurationUserChart: {
                    labels: this.chartLabels,
                    datasets: [{
                        label: this.$i18n.t('report.duration'),
                        data: [],
                        backgroundColor: '#444'
                    }],
                    options: {
                        legend: {
                            display: false
                        },
                        tooltips: {
                            callbacks: {
                                label: function(item) {
                                    let hours = parseInt(Math.floor(item.yLabel / 3600))
                                    let minutes = parseInt(Math.floor((item.yLabel - hours * 3600) / 60))
                                    let seconds = parseInt((item.yLabel - (hours * 3600 + minutes * 60)) % 60)
                                    
                                    let dHours = hours > 9 ? hours : '0' + hours
                                    let dMins = minutes > 9 ? minutes : '0' + minutes
                                    let dSecs = seconds > 9 ? seconds : '0' + seconds

                                    if (dHours === "00" && dMins !== "00") {
                                        return dMins + 'm ' + dSecs + 's'
                                    } else if (dHours === "00" && dMins === "00") {
                                        return dSecs + 's'
                                    } else {
                                        return dHours + 'h ' + dMins + 'm ' + dSecs + 's'
                                    }
                                }
                            }
                        },
                        title: {
                            display: true,
                            text: this.$i18n.t('report.average_duration_user')
                        },
                        scales: {
                            yAxes: [{
                                scaleLabel: {
                                    display: true,
                                    labelString: this.$i18n.t('report.time')
                                },
                                ticks: {
                                    maxTicksLimit: 5,
                                    callback: function(item) {
                                        let hours = parseInt(Math.floor(item / 3600))
                                        let minutes = parseInt(Math.floor((item - hours * 3600) / 60))
                                        let seconds = parseInt((item - (hours * 3600 + minutes * 60)) % 60)
                                        let dHours = hours > 9 ? hours : '0' + hours
                                        let dMins = minutes > 9 ? minutes : '0' + minutes
                                        let dSecs = seconds > 9 ? seconds : '0' + seconds
                                        if (dHours === "00") {
                                            return dMins + 'm ';
                                        } else {
                                            return dHours + 'h ' + dMins + 'm ';
                                        }
                                    },
                                    beginAtZero: true
                                }
                            }]
                        }
                    }
                },
                avgTrafficUserChart: {
                    labels: this.chartLabels,
                    datasets: [{
                            label: this.$i18n.t('report.upload'),
                            data: [],
                            backgroundColor: '#35a3d7'
                        },
                        {
                            label: this.$i18n.t('report.download'),
                            data: [],
                            backgroundColor: '#66b83c'
                        }
                    ],
                    options: {
                        tooltips: {
                            callbacks: {
                                label: function(item) {
                                    const prettyBytes = require('pretty-bytes');
                                    return prettyBytes(item.yLabel);
                                }
                            }
                        },
                        title: {
                            display: true,
                            text: this.$i18n.t('report.medium_traffic_user')
                        },
                        scales: {
                            yAxes: [{
                                scaleLabel: {
                                    display: true,
                                    labelString: this.$i18n.t('report.size')
                                },
                                ticks: {
                                    maxTicksLimit: 5,
                                    callback: function(item) {
                                        const prettyBytes = require('pretty-bytes');
                                        return prettyBytes(item);
                                    },
                                    beginAtZero: true
                                }
                            }]
                        }
                    }
                },
                avgDurationSessionChart: {
                    labels: this.chartLabels,
                    datasets: [{
                        label: this.$i18n.t('report.duration'),
                        data: [],
                        backgroundColor: '#444'
                    }],
                    options: {
                        legend: {
                            display: false
                        },
                        tooltips: {
                            callbacks: {
                                label: function(item) {
                                    let hours = parseInt(Math.floor(item.yLabel / 3600))
                                    let minutes = parseInt(Math.floor((item.yLabel - hours * 3600) / 60))
                                    let seconds = parseInt((item.yLabel - (hours * 3600 + minutes * 60)) % 60)

                                    let dHours = hours > 9 ? hours : '0' + hours
                                    let dMins = minutes > 9 ? minutes : '0' + minutes
                                    let dSecs = seconds > 9 ? seconds : '0' + seconds

                                    if (dHours === "00" && dMins !== "00") {
                                        return dMins + 'm ' + dSecs + 's'
                                    } else if (dHours === "00" && dMins === "00") {
                                        return dSecs + 's'
                                    } else {
                                        return dHours + 'h ' + dMins + 'm ' + dSecs + 's'
                                    }
                                }
                            }
                        },
                        scales: {
                            yAxes: [{
                                scaleLabel: {
                                    display: true,
                                    labelString: this.$i18n.t('report.time')
                                },
                                ticks: {
                                    maxTicksLimit: 5,
                                    callback: function(item) {
                                        let hours = parseInt(Math.floor(item / 3600))
                                        let minutes = parseInt(Math.floor((item - hours * 3600) / 60))
                                        let seconds = parseInt((item - (hours * 3600 + minutes * 60)) % 60)

                                        let dHours = hours > 9 ? hours : '0' + hours
                                        let dMins = minutes > 9 ? minutes : '0' + minutes
                                        let dSecs = seconds > 9 ? seconds : '0' + seconds

                                        if (dHours === "00") {
                                            return dMins + 'm ';
                                        } else {
                                            return dHours + 'h ' + dMins + 'm ';
                                        }
                                    },
                                    beginAtZero: true
                                }
                            }]
                        },
                        title: {
                            display: true,
                            text: this.$i18n.t('report.average_duration_connections')
                        }
                    }
                },
                avgTrafficSessionChart: {
                    labels: this.chartLabels,
                    datasets: [{
                            label: this.$i18n.t('report.upload'),
                            data: [],
                            backgroundColor: '#35a3d7'
                        },
                        {
                            label: this.$i18n.t('report.download'),
                            data: [],
                            backgroundColor: '#66b83c'
                        }
                    ],
                    options: {
                        tooltips: {
                            callbacks: {
                                label: function(item) {
                                    const prettyBytes = require('pretty-bytes');
                                    return prettyBytes(item.yLabel);
                                }
                            }
                        },
                        title: {
                            display: true,
                            text: this.$i18n.t('report.medium_traffic_connections')
                        },
                        scales: {
                            yAxes: [{
                                scaleLabel: {
                                    display: true,
                                    labelString: this.$i18n.t('report.size')
                                },
                                ticks: {
                                    maxTicksLimit: 5,
                                    callback: function(item) {
                                        const prettyBytes = require('pretty-bytes');
                                        return prettyBytes(item);
                                    },
                                    beginAtZero: true
                                }
                            }]
                        }
                    }
                },
                avgDurationConnection: 0,
                avgDownloadTrafficConnection: 0,
                avgUploadTrafficConnection: 0,
                avgDurationUser: 0,
                avgDownloadTrafficUser: 0,
                avgUploadTrafficUser: 0,
            }
        },
        created() {
            this.fillTotalTrafficChart()
            this.implementUserChart()
            this.implementSessionChart()
            this.implementAvgTable();
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
                let hours = parseInt(Math.floor(item / 3600))
                let minutes = parseInt(Math.floor((item - hours * 3600) / 60))
                let seconds = parseInt((item - (hours * 3600 + minutes * 60)) % 60)
                let dHours = hours > 9 ? hours : '0' + hours
                let dMins = minutes > 9 ? minutes : '0' + minutes
                let dSecs = seconds > 9 ? seconds : '0' + seconds
                if (dHours === "00" && dMins !== "00") {
                    return dMins + 'm ' + dSecs + 's'
                } else if (dHours === "00" && dMins === "00") {
                    return dSecs + 's'
                } else {
                    return dHours + 'h ' + dMins + 'm ' + dSecs + 's'
                }
            },
            fillTotalTrafficChart() {
                this.chartDateRange.forEach(date => {
                    let kbps_up = 0
                    let kbps_down = 0
                    this.newUsersReport.map(function(user) {
                        if (user.valid_from.substring(0, 10) === date) {
                            kbps_up += user.kbps_up
                            kbps_down += user.kbps_down
                        } else {
                            kbps_up += 0
                            kbps_down += 0
                        }
                    })
                    this.trafficChartData.datasets[0].data.push(kbps_up)
                    this.trafficChartData.datasets[1].data.push(kbps_down)
                })
            },
            implementUserChart() {
                this.chartDateRange.forEach(date => {
                    let moment_range = extendMoment(moment)
                    let index = 0
                    let kbps_up = [0]
                    let kbps_down = [0]
                    let duration_array = [0]
                    let all_sessions = this.sessionsReport
                    let session_counter = 0
                    let duration = 0
                    this.newUsersReport.map(function(user) {
                        // Check if user is active between dates
                        if (user.valid_from.substring(0, 10) <= date) {
                            // iterate through session
                            all_sessions.map(function(session) {
                                // Check if user has any session
                                if (session.user_id === user.id) {
                                    // check if user session start is on chart date
                                    if (session.start_time.substring(0, 10) === date) {
                                        // if user session is on date, check if session has been stopped on that day
                                        if (session.stop_time.substring(0, 10) === date) {
                                            duration += moment.range(moment(session.start_time), moment(session.stop_time)).diff('seconds')

                                            // or session has been stopped on next day
                                        } else if (session.stop_time.substring(0, 10) > date) {
                                            duration += moment.range(moment(session.start_time), moment(date).endOf('day')).diff('seconds')
                                        }
                                        duration_array[session_counter] += duration
                                        duration_array.push(0)
                                        session_counter++

                                        // check if session has started one day before, but has stopped on actual day
                                    } else if (session.stop_time.substring(0, 10) === date && session.start_time.substring(0, 10) < date) {
                                        duration += moment.range(moment(date).startOf('day'), moment(session.stop_time)).diff('seconds')
                                        duration_array[session_counter] += duration
                                        duration_array.push(0)
                                        session_counter++
                                    } else {
                                        duration_array[session_counter] += 0
                                    }
                                }
                            })
                        }

                        if (user.valid_from.substring(0, 10) === date) {
                            kbps_up[index] += user.kbps_up
                            kbps_down[index] += user.kbps_down
                            kbps_up.push(0)
                            kbps_down.push(0)
                            index++
                        } else {
                            kbps_up[index] += 0
                            kbps_down[index] += 0
                        }
                    })

                    if (duration_array.length > 1) {
                        duration_array.pop()
                    }
                    if (kbps_up.length > 1) {
                        kbps_up.pop()
                    }
                    if (kbps_down.length > 1) {
                        kbps_down.pop()
                    }

                    this.avgTrafficUserChart.datasets[0].data.push(this.calculateAVG(kbps_up))
                    this.avgTrafficUserChart.datasets[1].data.push(this.calculateAVG(kbps_down))
                    this.avgDurationUserChart.datasets[0].data.push(this.calculateAVG(duration_array))
                })
            },
            implementSessionChart() {
                this.chartDateRange.forEach(date => {
                    let index = 0
                    let duration = [0]
                    let kbps_up = [0]
                    let kbps_down = [0]
                    this.sessionsReport.map(function(session) {
                        // Check if session has been started between chart date
                        if (session.start_time.substring(0, 10) === date) {
                            kbps_up[index] += session.bytes_up
                            kbps_down[index] += session.bytes_down
                            duration[index] += session.duration
                            kbps_up.push(0)
                            kbps_down.push(0)
                            duration.push(0)
                            index++
                        } else {
                            kbps_up[index] += 0
                            kbps_down[index] += 0
                            duration[index] += 0
                        }
                    })
                    if (duration.length > 1) {
                        duration.pop()
                    }
                    if (kbps_up.length > 1) {
                        kbps_up.pop()
                    }
                    if (kbps_down.length > 1) {
                        kbps_down.pop()
                    }

                    this.avgTrafficSessionChart.datasets[0].data.push(this.calculateAVG(kbps_up))
                    this.avgTrafficSessionChart.datasets[1].data.push(this.calculateAVG(kbps_down))
                    this.avgDurationSessionChart.datasets[0].data.push(this.calculateAVG(duration))
                })
            },
            implementAvgTable() {
                const prettyBytes = require('pretty-bytes');
                this.avgDurationConnection = this.secondsToHR(this.calculateAVG(this.avgDurationSessionChart.datasets[0].data));
                this.avgUploadTrafficConnection = prettyBytes(this.calculateAVG(this.avgTrafficSessionChart.datasets[0].data));
                this.avgDownloadTrafficConnection = prettyBytes(this.calculateAVG(this.avgTrafficSessionChart.datasets[1].data));
                this.avgDurationUser = this.secondsToHR(this.calculateAVG(this.avgDurationUserChart.datasets[0].data));
                this.avgUploadTrafficUser = prettyBytes(this.calculateAVG(this.avgTrafficUserChart.datasets[0].data));
                this.avgDownloadTrafficUser = prettyBytes(this.calculateAVG(this.avgTrafficUserChart.datasets[1].data));
            }
        }
    }
</script>

<style scoped>
    .section-title {
        margin: 0px 0px 10px -15px;
    }
    .page-header h1 {
        margin-top: 50px;
        margin-bottom: -10px;
        font-size: 40px;
    }
    .average-table {
        border-top: 1px solid gray;
        margin: 20px;
    }
    .average-table h1 {
        font-size: 40px;
        font-weight: bold;
    }
    .average-table h3 {
        font-size: 18px;
    }
    .first-row>div:first-child,
    .first-row>div:nth-child(2),
    .second-row>div:first-child,
    .second-row>div:nth-child(2) {
        border-right: 1px dashed gray;
    }
    .second-row {
        border-top: 1px dashed gray;
    }
    .first-row>div,
    .second-row>div {
        min-height: 135px;
        text-align: center;
    }
    .first-row>div h3,
    .second-row>div h3 {
        max-height: 40px;
        height: 40px;
        font-size: 18px;
    }
    @media only screen and (max-width: 1200px) {
        .first-row>div h3,
        .second-row>div h3 {
            font-size: 17px;
        }
        .first-row>div h1,
        .second-row>div h1 {
            font-size: 32px;
        }
    }
    @media only screen and (max-width: 1030px) {
        .first-row>div h3,
        .second-row>div h3 {
            font-size: 17px;
        }
        .first-row>div h1,
        .second-row>div h1 {
            font-size: 26px;
        }
    }
    @media only screen and (max-width: 925px) {
        .first-row>div h3,
        .second-row>div h3 {
            font-size: 14px;
        }
        .first-row>div h1,
        .second-row>div h1 {
            font-size: 20px;
        }
        .first-row>div h3,
        .second-row>div h3 {
            max-height: 28px;
            height: 28px;
        }
    }
    @media only screen and (max-width: 768px) {
        .first-row>div h3,
        .second-row>div h3 {
            font-size: 18px;
        }
        .first-row>div h1,
        .second-row>div h1 {
            font-size: 35px;
        }
        .first-row>div:first-child,
        .first-row>div:nth-child(2),
        .second-row>div:first- .first-row>div h3,
        .second-row>div h3 {
            max-height: 25px;
            height: 25px;
        }
        .first-row>div:first-child,
        .first-row>div:nth-child(2),
        .second-row>div:first-child,
        .second-row>div:nth-child(2) {
            border-right: 0;
        } .page-header h1 {
            font-size: 35px;
        }
        .second-row {
            border-top: 0;
        }
        .page-header h1 {
            font-size: 35px;
        }
    }

</style>
