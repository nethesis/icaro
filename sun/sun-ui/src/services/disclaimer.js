var DisclaimerService = {
  methods: {
    disclaimersByAccount(accountId, success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/disclaimers/accounts/" + accountId, {
            headers: {
              Token: (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    disclaimerCreate(body, accountId, success, error) {
      this.$http
        .post(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/disclaimers/accounts/" + accountId, body, {
            headers: {
              Token: (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    disclaimerDelete(disclaimerId, success, error) {
      this.$http
        .delete(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/disclaimers/" +
          disclaimerId, {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
  }
};
export default DisclaimerService;
