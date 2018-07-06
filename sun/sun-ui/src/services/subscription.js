var SubscriptionService = {
  methods: {
    subscriptionPlansGetAll(success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/subscription_plans", {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    subscriptionBillingsGetAll(success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/billings", {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    subscriptionTaxesGetAll(success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/taxes", {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    subscriptionBillingUpdate(billingEmpty, success, error) {
      this.$http[billingEmpty ? "post" : "put"](
        this.$root.$options.api_scheme +
        this.$root.$options.api_host +
        "/api/billings",
        this.billingInfo, {
          headers: {
            Token:
              (this.get("loggedUser") && this.get("loggedUser").token) || ""
          }
        }
      ).then(success, error);
    },
    subscriptionRenew(payment, success, error) {
      this.$http
        .post(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/subscription_plans/renewal", {
            payment_id: payment.paymentID
          }, {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        ).then(success, error);
    },
    subscriptionUpgradePriceCheck(plan, success, error) {
      this.$http
        .get(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/subscription_plans/upgrade_price?plan=" +
          plan, {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        )
        .then(success, error);
    },
    subscriptionUpgrade(payment, plan, success, error) {
      this.$http
        .post(
          this.$root.$options.api_scheme +
          this.$root.$options.api_host +
          "/api/subscription_plans/upgrade", {
            payment_id: payment.paymentID,
            subscription_plan_id: plan.id
          }, {
            headers: {
              Token:
                (this.get("loggedUser") && this.get("loggedUser").token) || ""
            }
          }
        ).then(success, error);
    }
  }
};
export default SubscriptionService;
