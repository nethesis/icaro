var AccountService = {
    methods: {
      accountGetAll(success, error) {
        this.$http.get('https://' + this.$root.$options.api_host + '/api/accounts', {
          headers: {
            'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
          }
        }).then(success, error);
      },
      accountGet(id) {
        this.$http.get('https://' + this.$root.$options.api_host + '/api/accounts/' + id, {
          headers: {
            'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
          }
        }).then(success, error);
      },
      accountModify(id, body, success, error) {
        this.$http.put('https://' + this.$root.$options.api_host + '/api/accounts/' + id, body, {
          headers: {
            'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
          }
        }).then(success, error);
      },
      accountDelete(id, success, error) {
        this.$http.delete('https://' + this.$root.$options.api_host + '/api/accounts/' + id, {
          headers: {
            'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
          }
        }).then(success, error);
      },
      accountChangePassword(id, password, success, error) {
        this.$http.put('https://' + this.$root.$options.api_host + '/api/accounts/' + id, { password: password }, {
          headers: {
            'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
          }
        }).then(success, error);
      },
      accountCreate(body, success, error) {
        this.$http.post('https://' + this.$root.$options.api_host + '/api/accounts', body, {
          headers: {
            'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
          }
        }).then(success, error);
      },
    }
  };
  export default AccountService;
