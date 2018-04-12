var UserService = {
  methods: {
    userGetAll(hotspotId, accountType, success, error) {
      this.$http.get(this.$root.$options.api_scheme + this.$root.$options.api_host + '/api/users?' +
        (hotspotId && hotspotId != 0 ? '&hotspot=' + hotspotId : '') +
        (accountType && accountType.length > 0 ? '&type=' + accountType : ''), {
          headers: {
            'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
          }
        }).then(success, error);
    },
    userGet(id, success, error) {
      this.$http.get(this.$root.$options.api_scheme + this.$root.$options.api_host + '/api/users/' + id, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    userMACCreate(mac, digest, uuid, success, error) {
      this.$http.get(this.$root.$options.api_scheme + this.$root.$options.api_host + '/wax/register/mac/' +
        mac.username.toUpperCase().replace(/:/gi, '-') + '?' +
        '&digest=' + digest +
        '&uuid=' + uuid +
        '&name=' + mac.name + ' (' + mac.unit.name + ')' +
        '&kbps_down=' + mac.kbps_down +
        '&kbps_up=' + mac.kbps_up, {
          headers: {
            'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
          }
        }).then(success, error);
    },
    userModify(id, body, success, error) {
      this.$http.put(this.$root.$options.api_scheme + this.$root.$options.api_host + '/api/users/' + id, body, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
    userDelete(id, success, error) {
      this.$http.delete(this.$root.$options.api_scheme + this.$root.$options.api_host + '/api/users/' + id, {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
  }
};
export default UserService;
