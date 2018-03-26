var PreferenceService = {
  methods: {
    hsPrefGet(id, success, error) {
      this.$http.get(this.$root.$options.api_scheme + this.$root.$options.api_host + '/api/preferences/hotspots/' + id, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    hsPrefModify(id, body, success, error) {
      this.$http.put(this.$root.$options.api_scheme + this.$root.$options.api_host + '/api/preferences/hotspots/' + id, body, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
  }
};
export default PreferenceService;
