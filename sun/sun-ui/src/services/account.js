var AccountService = {
  methods: {
    accountGetAll(hotspotId, success, error) {
      this.$http.get(this.$root.$options.api_scheme + this.$root.$options.api_host + '/api/accounts?' +
      (hotspotId && hotspotId != 0 ? '&hotspot=' + hotspotId : ''), {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    accountGet(id, success, error) {
      this.$http.get(this.$root.$options.api_scheme + this.$root.$options.api_host + '/api/accounts/' + id, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    accountModify(id, body, success, error) {
      this.$http.put(this.$root.$options.api_scheme + this.$root.$options.api_host + '/api/accounts/' + id, body, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    accountDelete(id, success, error) {
      this.$http.delete(this.$root.$options.api_scheme + this.$root.$options.api_host + '/api/accounts/' + id, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    accountChangePassword(id, password, success, error) {
      this.$http.put(this.$root.$options.api_scheme + this.$root.$options.api_host + '/api/accounts/' + id, {
        password: password
      }, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    accountCreate(body, success, error) {
      this.$http.post(this.$root.$options.api_scheme + this.$root.$options.api_host + '/api/accounts', body, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
  }
};
export default AccountService;
