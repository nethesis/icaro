var StatsService = {
  methods: {
    statsHotspotsTotal(success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/stats/hotspots/total", {
            headers: {
              Token: (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    statsUnitsTotal(success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/stats/units/total", {
            headers: {
              Token: (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    statsAccountsTotal(success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/stats/accounts/total", {
            headers: {
              Token: (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    statsDevicesTotal(success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/stats/devices/total", {
            headers: {
              Token: (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    statsUsersTotal(success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/stats/users/total", {
            headers: {
              Token: (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    statsSessionsTotal(success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/stats/sessions/total", {
            headers: {
              Token: (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    statsSMSTotalForAccount(success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/stats/sms/accounts", {
            headers: {
              Token: (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    statsSMSTotalForAccountByAccount(accountId, success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/stats/sms/accounts/" + accountId, {
            headers: {
              Token: (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    updateSMSTotalForAccountByAccount(body, accountId, success, error) {
      this.$http
        .post(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/stats/sms/accounts/" + accountId, body, {
            headers: {
              Token: (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    statsSMSSentByHotspot(hotspotId, success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          (hotspotId > 0 ? "/api/stats/sms/hotspots/" + hotspotId : "/api/stats/sms/hotspots"), {
            headers: {
              Token: (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    reportsHistory(graph, hotspotId, range, success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/stats/reports/" + graph + "/graph?" +
          (hotspotId != 0 ? "&hotspot=" + hotspotId : "") +
          (range != 0 ? "&range=" + range : ""), {
            headers: {
              Token: (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    }
  }
};
export default StatsService;
