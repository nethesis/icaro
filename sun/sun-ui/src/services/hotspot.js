var HotspotService = {
  methods: {
    hotspotGetAll(page, limit, q, success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/hotspots?" +
          (q ? "&q=" + q : "") +
          (page ? "&page=" + page : "") +
          (limit ? "&limit=" + limit : ""), {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    hotspotGet(id, success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/hotspots/" +
          id, {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    hotspotCreate(body, success, error) {
      this.$http
        .post(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/hotspots",
          body, {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    hotspotModify(id, body, success, error) {
      this.$http
        .put(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/hotspots/" +
          id,
          body, {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    hotspotDelete(id, success, error) {
      this.$http
        .delete(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/hotspots/" +
          id, {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    hotspotGetVouchers(filters, hotspotId, success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/vouchers/" +
          hotspotId + "?" +
          (filters && filters.code != 0 ? "&code=" + filters.code : "") +
          (filters && filters.duration != 0 ? "&duration=" + filters.duration : "") +
          (filters && filters.auto_login ? "&auto_login=" + filters.auto_login : "") +
          (filters && filters.used ? "&used=" + filters.used : "") +
          (filters && filters.reusable ? "&reusable=" + filters.reusable : "") +
          (filters && filters.printed ? "&printed=" + filters.printed : ""), {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    hotspotCreateVoucher(body, success, error) {
      this.$http
        .post(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/vouchers",
          body, {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    hotspotUpdateVoucher(id, body, success, error) {
      this.$http
        .put(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/vouchers/" + id,
          body, {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    hotspotVoucherDelete(id, success, error) {
      this.$http
        .delete(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/vouchers/" +
          id, {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    hotspotPrivacy(uuid, success, error) {
      this.$http
        .get(
          this.$root.$options.wax_url +
          "/wax/privacy/" +
          uuid
        )
        .then(success, error);
    }
  }
};
export default HotspotService;
