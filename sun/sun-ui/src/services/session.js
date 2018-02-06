var SessionService = {
  methods: {
    sessionGetAll(hotspotId, success, error) {
      this.$http.get('https://' + this.$root.$options.api_host + '/api/sessions' + (hotspotId && hotspotId != 0 ? '?hotspot=' + hotspotId : ''), {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    sessionGet(id, success, error) {
      this.$http.get('https://' + this.$root.$options.api_host + '/api/sessions/' + id, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
  }
};
export default SessionService;
