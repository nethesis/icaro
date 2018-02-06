var UserService = {
  methods: {
    userGetAll(hotspotId, success, error) {
      this.$http.get('https://' + this.$root.$options.api_host + '/api/users' + (hotspotId && hotspotId != 0 ? '?hotspot=' + hotspotId : ''), {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    userGet(id, success, error) {
      this.$http.get('https://' + this.$root.$options.api_host + '/api/users/' + id, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    userModify(id, body, success, error) {
      this.$http.put('https://' + this.$root.$options.api_host + '/api/users/' + id, body, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    userDelete(id, success, error) {
      this.$http.delete('https://' + this.$root.$options.api_host + '/api/users/' + id, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
  }
};
export default UserService;
