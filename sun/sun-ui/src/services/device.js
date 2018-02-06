var DeviceService = {
  methods: {
    deviceGetAll(hotspotId, success, error) {
      this.$http.get('https://' + this.$root.$options.api_host + '/api/devices' + (hotspotId && hotspotId != 0 ? '?hotspot=' + hotspotId : ''), {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    deviceGet(id, success, error) {
      this.$http.get('https://' + this.$root.$options.api_host + '/api/devices/' + id, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
  }
};
export default DeviceService;
