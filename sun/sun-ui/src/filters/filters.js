var Filters = {
  byteFormat: function(size) {
    var result;

    switch (true) {
      case size === null || size === "" || isNaN(size) || size === 0:
        result = "-";
        break;

      case size >= 0 && size < 1024:
        result = size + " B";
        break;

      case size >= 1024 && size < Math.pow(1024, 2):
        result = Math.round((size / 1024) * 100) / 100 + " KB";
        break;

      case size >= Math.pow(1024, 2) && size < Math.pow(1024, 3):
        result = Math.round((size / Math.pow(1024, 2)) * 100) / 100 + " MB";
        break;

      case size >= Math.pow(1024, 3) && size < Math.pow(1024, 4):
        result = Math.round((size / Math.pow(1024, 3)) * 100) / 100 + " GB";
        break;

      default:
        result = Math.round((size / Math.pow(1024, 4)) * 100) / 100 + " TB";
    }

    return result;
  },
  secondsInHour: function(value) {
    if (value <= 0 || value == "-" || value == null) {
      return "-";
    }
    let hours = parseInt(Math.floor(value / 3600));
    let minutes = parseInt(Math.floor((value - hours * 3600) / 60));
    let seconds = parseInt((value - (hours * 3600 + minutes * 60)) % 60);

    let dHours = hours > 9 ? hours : "0" + hours;
    let dMins = minutes > 9 ? minutes : "0" + minutes;
    let dSecs = seconds > 9 ? seconds : "0" + seconds;

    let returnString = "";
    if (dHours != "00") returnString += dHours + "h ";
    if (dMins != "00") returnString += dMins + "m ";
    if (dSecs != "00") returnString += dSecs + "s ";

    return returnString;
  },
  formatDate: function(value, withHour, empty) {
    var moment = require("patternfly/node_modules/moment/moment.js");
    if (+new Date(value) > 0)
      return moment(String(value)).format(
        withHour ? "DD MMMM YYYY" : "DD MMMM YYYY, HH:mm"
      );
    else return empty ? "" : "-";
  },
  adjustPage: function(value) {
    return Math.ceil(value);
  }
};

export default Filters;
