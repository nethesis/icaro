var Filters = {
  byteFormat: require('vue-filters-kit/filters/byteFormatter'),
  secondsInHour: function (value) {
    let hours = parseInt(Math.floor(value / 3600));
    let minutes = parseInt(Math.floor((value - (hours * 3600)) / 60));
    let seconds = parseInt((value - ((hours * 3600) + (minutes * 60))) % 60);

    let dHours = (hours > 9 ? hours : '0' + hours);
    let dMins = (minutes > 9 ? minutes : '0' + minutes);
    let dSecs = (seconds > 9 ? seconds : '0' + seconds);

    return dHours + "h " + dMins + "m " + dSecs + "s";
  },
  formatDate: function (value) {
    var moment = require("patternfly/node_modules/moment/moment.js")
    if (+new Date(value) > 0)
      return moment(String(value)).format('DD MMMM YYYY, HH:mm')
    else
      return '-'
  },
}

export default Filters
