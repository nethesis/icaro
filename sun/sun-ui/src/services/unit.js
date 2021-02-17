var UnitService = {
  methods: {
    unitGetAll(hotspotId, page, limit, q, success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/units?" +
          (q ? "&q=" + q : "") +
          (page ? "&page=" + page : "") +
          (limit ? "&limit=" + limit : "") +
          (hotspotId && hotspotId != 0 ? "&hotspot=" + hotspotId : ""), {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    unitGet(id, success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/units/" +
          id, {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    unitModify(id, body, success, error) {
      this.$http
        .put(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/units/" +
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
    unitDelete(id, success, error) {
      this.$http
        .delete(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/units/" +
          id, {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    unitGetStatus(page, limit, unitId, dateFrom, dateTo, success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
            this.$root.$options.api_host +
            "/api/sessions?" +
            (page ? "&page=" + page : "") +
            (limit ? "&limit=" + limit : "") +
            (unitId ? "&unit=" + unitId : "") +
            (dateFrom ? "&from=" + dateFrom : "") +
            (dateTo ? "&to=" + dateTo : ""), {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || "",
            }
          }
        )
        .then(success, error);
  
      }
  }
};
export default UnitService;
