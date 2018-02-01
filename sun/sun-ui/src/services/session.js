var SessionService = {
  methods: {
    sessionGetAll(success, error) {
      this.$http.get('https://' + this.$root.$options.api_host + '/api/sessions', {
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
