var StatsService = {
  methods: {
    statsHotspotsTotal(success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
            this.$root.$options.api_host +
            "/api/stats/hotspots/total",
          {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
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
            "/api/stats/units/total",
          {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
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
            "/api/stats/accounts/total",
          {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
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
            "/api/stats/devices/total",
          {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
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
            "/api/stats/users/total",
          {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
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
            "/api/stats/sessions/total",
          {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    statsSMSTotal(success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
            this.$root.$options.api_host +
            "/api/stats/sms",
          {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    statsSMSSent(hotspotId, success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
            this.$root.$options.api_host +
            "/api/stats/sms/" + hotspotId,
          {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    }
  }
};
export default StatsService;
