var DeviceService = {
  methods: {
    deviceGetAll(success, error) {
      this.$http.get('https://' + this.$root.$options.api_host + '/api/devices', {
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
