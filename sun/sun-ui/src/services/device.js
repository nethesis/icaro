var DeviceService = {
  methods: {
    deviceGetAll(hotspotId, userId, success, error) {
      this.$http.get(this.$root.$options.api_scheme + this.$root.$options.api_host + '/api/devices?' +
        (hotspotId && hotspotId != 0 ? '&hotspot=' + hotspotId : '') +
        (userId && userId != 0 ? '&user=' + userId : ''), {
          headers: {
            'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
          }
        }).then(success, error);
    },
    deviceGet(id, success, error) {
      this.$http.get(this.$root.$options.api_scheme + this.$root.$options.api_host + '/api/devices/' + id, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
  }
};
export default DeviceService;
