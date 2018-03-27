var HistoryService = {
  methods: {
    historiesGetAll(hotspotId, userId, unitId, dateFrom, dateTo, success, error) {
      this.$http.get(this.$root.$options.api_scheme + this.$root.$options.api_host + '/api/histories?' +
        (hotspotId && hotspotId != 0 ? '&hotspot=' + hotspotId : '') +
        (userId && userId != 0 ? '&user=' + userId : '') +
        (unitId && unitId != 0 ? '&unit=' + unitId : '') +
        (dateFrom && dateFrom != 0 ? '&from=' + dateFrom : '') +
        (dateTo && dateTo != 0 ? '&to=' + dateTo : ''), {
          headers: {
            'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
          }
        }).then(success, error);
    },
    historiesGet(id, success, error) {
      this.$http.get(this.$root.$options.api_scheme + this.$root.$options.api_host + '/api/histories/' + id, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
  }
};
export default HistoryService;
