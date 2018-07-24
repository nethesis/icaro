var SessionService = {
  methods: {
    sessionGetAll(hotspotId, userId, unitId, dateFrom, dateTo, page, limit, success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/sessions?" +
          (page ? "&page=" + page : "") +
          (limit ? "&limit=" + limit : "") +
          (hotspotId && hotspotId != 0 ? "&hotspot=" + hotspotId : "") +
          (userId && userId != 0 ? "&user=" + userId : "") +
          (unitId && unitId != 0 ? "&unit=" + unitId : "") +
          (dateFrom && dateFrom != 0 ? "&from=" + dateFrom : "") +
          (dateTo && dateTo != 0 ? "&to=" + dateTo : ""), {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    sessionGet(id, success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/sessions/" +
          id, {
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
export default SessionService;
