var MarketingService = {
  methods: {
    marketingPrefGet(id, success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/marketing/" +
          id, {
            headers: {
              Token: (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    marketingPrefModify(id, body, success, error) {
      this.$http
        .put(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/marketing/" +
          id,
          body, {
            headers: {
              Token: (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    testMail(type, body, success, error) {
      this.$http
        .post(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/marketing/testmail/" +
          type,
          body, {
            headers: {
              Token: (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    }
  }
};
export default MarketingService;
