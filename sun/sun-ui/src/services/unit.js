var UnitService = {
  methods: {
    unitGetAll(hotspotId, success, error) {
      this.$http.get('https://' + this.$root.$options.api_host + '/api/units' + (hotspotId && hotspotId != 0 ? '?hotspot=' + hotspotId : ''), {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    unitGet(id, success, error) {
      this.$http.get('https://' + this.$root.$options.api_host + '/api/units/' + id, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    unitModify(id, body, success, error) {
      this.$http.put('https://' + this.$root.$options.api_host + '/api/units/' + id, body, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    unitDelete(id, success, error) {
      this.$http.delete('https://' + this.$root.$options.api_host + '/api/units/' + id, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
  }
};
export default UnitService;
