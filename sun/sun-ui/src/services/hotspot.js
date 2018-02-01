var HotspotService = {
  methods: {
    hotspotGetAll(success, error) {
      this.$http.get('https://' + this.$root.$options.api_host + '/api/hotspots', {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    hotspotGet(id, success, error) {
      this.$http.get('https://' + this.$root.$options.api_host + '/api/hotspots/' + id, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    hotspotCreate(body, success, error) {
      this.$http.post('https://' + this.$root.$options.api_host + '/api/hotspots', body, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    hotspotModify(id, body, success, error) {
      this.$http.put('https://' + this.$root.$options.api_host + '/api/hotspots/' + id, body, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    hotspotDelete(id, success, error) {
      this.$http.delete('https://' + this.$root.$options.api_host + '/api/hotspots/' + id, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
  }
};
export default HotspotService;
