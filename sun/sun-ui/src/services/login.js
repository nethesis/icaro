var LoginService = {
  methods: {
    execLogin(body, success, error) {
      this.$http.post('https://' + window.location.host + '/api/login', body).then(success, error);
    },
    execLogout(success, error) {
      this.$http.post('https://' + window.location.host + '/api/logout', {}, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    execGetInfo(id, success, error) {
      this.$http.get('https://' + window.location.host + '/api/accounts/' + id, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    execChangePassword(password, id, success, error) {
      this.$http.put('https://' + window.location.host + '/api/accounts/' + id, {
        password: password
      }, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    }
  }
};
export default LoginService;
