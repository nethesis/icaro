import Vue from "vue"

var Filters = {
  byteFormat: function (size) {
    var result;

    switch (true) {
      case size === null || size === "" || isNaN(size):
        result = "-";
        break;

      case size >= 0 && size < 1024:
        result = size + " B";
        break;

      case size >= 1024 && size < Math.pow(1024, 2):
        result = Math.round(size / 1024 * 100) / 100 + " K";
        break;

      case size >= Math.pow(1024, 2) && size < Math.pow(1024, 3):
        result = Math.round(size / Math.pow(1024, 2) * 100) / 100 + " M";
        break;

      case size >= Math.pow(1024, 3) && size < Math.pow(1024, 4):
        result = Math.round(size / Math.pow(1024, 3) * 100) / 100 + " G";
        break;

      default:
        result = Math.round(size / Math.pow(1024, 4) * 100) / 100 + " T";
    }

    return result;
  },
  secondsInHour: function (value) {
    let hours = parseInt(Math.floor(value / 3600));
    let minutes = parseInt(Math.floor((value - hours * 3600) / 60));
    let seconds = parseInt((value - (hours * 3600 + minutes * 60)) % 60);

    let dHours = hours > 9 ? hours : "0" + hours;
    let dMins = minutes > 9 ? minutes : "0" + minutes;
    let dSecs = seconds > 9 ? seconds : "0" + seconds;

    return dHours + "h " + dMins + "m " + dSecs + "s";
  },
  capitalize: function (value) {
    return value && value.toString().charAt(0).toUpperCase() + value.toString().slice(1);
  },
  uppercase: function (value) {
    return value && value.toUpperCase()
  },
  isEmpty: function (value) {
    jQuery.isEmptyObject(value);
  },
  camelToSentence: function (value) {
    var result = value.replace(/([A-Z]+)/g, " $1").replace(/([A-Z][a-z])/g, " $1");
    return result.charAt(0).toUpperCase() + result.slice(1);
  },
  sanitize: function (value) {
    return value.replace(/^[^a-z]+|[^\w-]+/gi, "");
  }
};

for (var f in Filters) {
  Vue.filter(f, Filters[f])
}
