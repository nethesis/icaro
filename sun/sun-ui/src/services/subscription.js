var SubscriptionService = {
  methods: {
    subscriptionPlansGetAll(success, error) {
      this.$http.get(this.$root.$options.api_scheme + this.$root.$options.api_host + '/api/subscription_plans/', {
        headers: {
          'Token': this.get('loggedUser') && this.get('loggedUser').token || ''
        }
      }).then(success, error);
    },
  }
};
export default SubscriptionService;
