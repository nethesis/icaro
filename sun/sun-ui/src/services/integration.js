var IntegrationService = {
  methods: {
    integrationGetAll(success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
            this.$root.$options.api_host +
            "/api/integrations",
          {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    mapGetAll(hotspotId, success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
            this.$root.$options.api_host +
            "/api/hotspot_integrations/" +
            hotspotId,
          {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    mapAccountsGetAll(accountId, success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
            this.$root.$options.api_host +
            "/api/account_integrations/" +
            accountId,
          {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    mapCreate(hotspotId, integrationId, success, error) {
      this.$http
        .put(
          this.$root.$options.api_scheme +
            this.$root.$options.api_host +
            "/api/hotspot_integrations/" +
            hotspotId +
            "/" +
            integrationId,
          null,
          {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    mapAccountsCreate(accountId, integrationId, success, error) {
      this.$http
        .post(
          this.$root.$options.api_scheme +
            this.$root.$options.api_host +
            "/api/account_integrations/" +
            accountId +
            "/" +
            integrationId,
          null,
          {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    mapDelete(hotspotId, integrationId, success, error) {
      this.$http
        .delete(
          this.$root.$options.api_scheme +
            this.$root.$options.api_host +
            "/api/hotspot_integrations/" +
            hotspotId +
            "/" +
            integrationId,
          {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    mapAccountsDelete(accountId, integrationId, success, error) {
      this.$http
        .delete(
          this.$root.$options.api_scheme +
            this.$root.$options.api_host +
            "/api/account_integrations/" +
            accountId +
            "/" +
            integrationId,
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
export default IntegrationService;
