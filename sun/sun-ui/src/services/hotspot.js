var HotspotService = {
  methods: {
    execGetAll(success, error) {
      this.$http.get('https://' + this.$root.$options.api_host + '/api/hotspots', {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    execGet(id) {
      this.$http.get('https://' + this.$root.$options.api_host + '/api/hotspots/' + id, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    execCreate() {

    },
    execModify(id, body, success, error) {
      this.$http.put('https://' + this.$root.$options.api_host + '/api/hotspots/' + id, body, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    execDelete(id, success, error) {
      this.$http.delete('https://' + this.$root.$options.api_host + '/api/hotspots/' + id, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
  }
};
export default HotspotService;
