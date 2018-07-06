<template>
  <span>
    <button :disabled="status" v-if="isExpired(obj.subscription.valid_until) || plans.length > 0 && plans[plans.length-1].code != obj.subscription.subscription_plan.code"
      @click="showRenewModal(obj.id)" type="button" class="btn btn-primary">
      <span class="fa fa-shopping-cart"></span>
      {{isExpired(obj.subscription.valid_until) ? $t('payment.renew_button') : $t('payment.upgrade_button')}}
    </button>
    <div v-if="!(isExpired(obj.subscription.valid_until) || plans.length > 0 && plans[plans.length-1].code != obj.subscription.subscription_plan.code)">
      <span class="fa fa-star"></span>
      {{$t('payment.plans_reached')}}
    </div>
    <div class="modal fade" :id="'paymentModalRenew-'+obj.id" data-backdrop="static" tabindex="-1" role="dialog" aria-labelledby="myModalLabel"
      aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <button @click="hideRenewModal()" type="button" class="close" data-dismiss="modal" aria-hidden="true">
              <span class="pficon pficon-close"></span>
            </button>
            <h4 class="modal-title" id="myModalLabel">{{onUpgrade ? $t('payment.upgrade') : $t('payment.renew')}}</h4>
          </div>
          <div class="modal-body">
            <div class="card-pf-view">
              <div>
                <div class="card-pf-top-element">
                  <span :class="['fa', onUpgrade ? 'fa-arrow-up' : 'fa-refresh', 'card-pf-icon-circle', 'adjust-icon-size']"></span>
                </div>
                <h2 class="card-pf-title text-center">
                  {{onUpgrade ? $t('payment.upgrade_proceed') : $t('payment.renew_proceed')}}
                </h2>
                <div class="card-pf-items text-center">
                  <div class="card-pf-item details-pay-item">
                    <span class="card-pf-item-text">
                      <strong>{{$t('payment.plan')}}</strong>
                    </span>
                  </div>
                  <div class="card-pf-item details-pay-item">
                    <div class="dropdown">
                      <button class="btn btn-default dropdown-toggle" type="button" id="dropdownMenu1" data-toggle="dropdown">
                        {{currentPlan.name}}
                        <span class="caret"></span>
                      </button>
                      <ul class="dropdown-menu" role="menu" aria-labelledby="dropdownMenu1">
                        <li v-for="p in plans" v-if="p.code !== 'free'" v-bind:key="p.code" role="presentation" :class="[p.code == currentPlan.code ? 'selected' : '', p.code != obj.subscription.subscription_plan.code && (p.full_price || p.price) <= obj.subscription.subscription_plan.price ? 'disabled' : '']">
                          <a @click="p.code != obj.subscription.subscription_plan.code && (p.full_price || p.price) <= obj.subscription.subscription_plan.price ? null : changePlan(p)"
                            role="menuitem" tabindex="-1" href="#">{{p.name}}</a>
                        </li>
                      </ul>
                    </div>
                  </div>
                </div>
                <div class="card-pf-items text-center">
                  <div class="card-pf-item details-pay-item">
                    <span class="card-pf-item-text">
                      <strong>{{$t('payment.details')}}</strong>
                    </span>
                  </div>
                  <div class="card-pf-item details-pay-item">
                    <div v-if="!onUpgradePriceCalc" class="details-markdown text-left" v-html="markdownDescription"></div>
                    <div v-if="onUpgradePriceCalc" class="spinner spinner-sm"></div>
                  </div>
                </div>
                <div class="card-pf-items text-center">
                  <div class="card-pf-item details-pay-item">
                    <span class="card-pf-item-text">
                      <strong>{{$t('payment.price')}}</strong>
                    </span>
                  </div>
                  <div class="card-pf-item details-pay-item">
                    <span v-if="!onUpgradePriceCalc" class="card-pf-item-text">
                      <strong>{{currentPlan.full_price > 0 ? currentPlan.full_price : 0}}€</strong>
                    </span>
                    <div v-if="onUpgradePriceCalc" class="spinner spinner-sm"></div>
                  </div>
                </div>

                <div v-if="onUpgrade && currentPlan.price != currentPlan.full_price && discounts.annualDiscount > 0" class="card-pf-items text-center">
                  <div class="card-pf-item details-pay-item">
                    <span class="card-pf-item-text">
                      <strong>{{$t('payment.previous_license_discount')}}</strong>
                    </span>
                  </div>
                  <div class="card-pf-item details-pay-item">
                    <span v-if="!onUpgradePriceCalc" class="card-pf-item-text">
                      <strong>{{discounts.annualDiscount}}€</strong>
                    </span>
                    <div v-if="onUpgradePriceCalc" class="spinner spinner-sm"></div>
                  </div>
                </div>

                <div class="card-pf-items text-center">
                  <div class="card-pf-item details-pay-item">
                    <span class="card-pf-item-text">
                      <strong>{{$t('payment.period')}}</strong>
                    </span>
                  </div>
                  <div class="card-pf-item details-pay-item">
                    <span v-if="!onUpgrade && !onUpgradePriceCalc" class="card-pf-item-text">{{obj.subscription.valid_until | formatDate(false)}} - {{calculateSubscription(obj.subscription.valid_until,
                      obj.subscription.subscription_plan) | formatDate(false)}}</span>
                    <span v-if="onUpgrade && !onUpgradePriceCalc" class="card-pf-item-text">{{new Date().toISOString() | formatDate(false)}} - {{calculateSubscription(new Date().toISOString(),
                      currentPlan) | formatDate(false)}}</span>
                    <div v-if="onUpgradePriceCalc" class="spinner spinner-sm"></div>
                  </div>
                </div>

                <div class="card-pf-items text-center">
                  <div class="card-pf-item details-pay-item">
                    <span class="card-pf-item-text">
                      <strong>{{$t('payment.final_price')}}</strong>
                    </span>
                  </div>
                  <div class="card-pf-item details-pay-item">
                    <span v-if="!onUpgradePriceCalc" class="card-pf-item-text">
                      <strong>{{currentPlan.price}}€</strong>
                      <span>+ {{$t('payment.taxes')}}</span>
                    </span>
                    <div v-if="onUpgradePriceCalc" class="spinner spinner-sm"></div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="divider"></div>
          <div class="modal-footer text-center">
            <div v-if="payment.onProgress" class="spinner"></div>
            <div v-if="payment.done && !errors.state && !payment.onProgress" class="alert alert-success alert-dismissable">
              <span class="pficon pficon-ok"></span>
              <strong>{{$t('payment.confirmed')}}</strong>. {{$t('payment.payment_id_ref')}}
              <pre class="filters-container"><strong>{{payment.details.paymentID}}</strong></pre>
            </div>
            <button v-if="payment.done && !errors.state && !payment.onProgress" type="button" class="btn btn-primary done-payment-button"
              @click="hideRenewModal()">{{$t('payment.done')}}</button>
            <div v-if="!payment.done && !errors.state && !payment.onProgress" class="card-pf-item details-pay-item">
              <span class="card-pf-item-text">
                <strong>{{$t('payment.pay_with')}}</strong>
              </span>
            </div>
            <div v-show="!payment.done && !errors.state && !payment.onProgress" :id="'paypal-button-container-'+obj.id"></div>
            <div v-if="payment.done && errors.state && !payment.onProgress" class="alert alert-danger alert-dismissable">
              <span class="pficon pficon-error-circle-o"></span>
              <strong>{{$t('payment.payment_error')}}</strong>. {{errors.message || $t('payment.payment_error_details')}}
              <p>{{$t('payment.payment_id_ref')}}</p>
              <pre class="filters-container"><strong>{{payment.details.paymentID}}</strong></pre>
            </div>
          </div>
        </div>
      </div>
    </div>
  </span>
</template>
<script>
  import StorageService from "./../../services/storage";
  import SubscriptionService from "./../../services/subscription";

  import paypal from "paypal-checkout";
  import {
    setTimeout
  } from "timers";
  import marked from "marked";
  import _ from "lodash";

  export default {
    name: "RenewButton",
    props: ["obj", "status", "update"],
    mixins: [StorageService, SubscriptionService],
    mounted() {
      var context = this;
      paypal.Button.render({
          env: CONFIG.PAYPAL_PRODUCTION ? "production" : "sandbox",
          style: {
            layout: "vertical", // horizontal | vertical
            size: "medium", // medium | large | responsive
            shape: "rect", // pill | rect
            color: "gold" // gold | blue | silver | black
          },
          funding: {
            allowed: [paypal.FUNDING.CARD],
            disallowed: [paypal.FUNDING.CREDIT]
          },
          client: {
            sandbox: CONFIG.PAYPAL_SANDBOX,
            production: CONFIG.PAYPAL_PRODUCTION
          },
          payment: function (data, actions) {
            return actions.payment.create({
              payment: {
                transactions: [{
                  amount: {
                    total: Math.round(
                      (context.currentPlan.price +
                        context.currentPlan.price *
                        context.billingInfo.tax /
                        100) *
                      100
                    ) / 100,
                    currency: "EUR",
                    details: {
                      subtotal: Math.round(context.currentPlan.price * 100) / 100,
                      tax: Math.round(
                        context.currentPlan.price *
                        context.billingInfo.tax /
                        100 *
                        100
                      ) / 100
                    }
                  },
                  item_list: {
                    items: [{
                      name: context.currentPlan.code,
                      description: context.currentPlan.name,
                      sku: context.obj.uuid,
                      price: context.currentPlan.price,
                      currency: "EUR",
                      quantity: "1"
                    }]
                  }
                }]
              },
              experience: {
                input_fields: {
                  no_shipping: 1
                }
              }
            });
          },

          onAuthorize: function (data, actions) {
            return actions.payment.execute().then(function () {
              context.payment.onProgress = true;
              if (context.onUpgrade) {
                context.upgradeCheck(data);
              } else {
                context.renewCheck(data);
              }
            });
          }
        },
        "#paypal-button-container-" + this.obj.id
      );
    },
    data() {
      // get plans
      this.plansList();

      // read upgrade ref and show modal
      if (this.get("upgrade_ref", false)) {
        var context = this;
        setTimeout(function () {
          context.showRenewModal(context.get("upgrade_ref", false));
          context.delete("upgrade_ref");
        }, 0);
      }

      return {
        payment: {
          done: false,
          details: {},
          onProgress: false
        },
        errors: {
          message: "",
          state: false
        },
        plans: [],
        markdownDescription: "",
        currentPlan: this.obj.subscription.subscription_plan,
        billingInfo: {},
        discounts: {
          annualDiscount: 0
        },
        onUpgradePriceCalc: false,
        onUpgrade: false
      };
    },
    methods: {
      isExpired(date) {
        return new Date().toISOString() > date;
      },
      showRenewModal(id) {
        this.subscriptionBillingsGetAll(
          success => {
            this.billingInfo = success.body;
            this.payment.done = false;
            this.payment.onProgress = false;
            this.payment.details = {};
            this.errors.message = "";
            this.errors.state = false;
            this.currentPlan =
              this.obj.subscription.subscription_plan.code != "free" ?
              this.obj.subscription.subscription_plan :
              this.plans[1];
            if (this.obj.subscription.subscription_plan.code == "free") {
              this.currentPlan.full_price = this.plans[1].price;
            } else {
              this.currentPlan.full_price = this.currentPlan.price;
            }

            this.markdownDescription = marked(this.currentPlan.description, {
              sanitize: true
            });
            this.onUpgradePriceCalc = false;
            this.onUpgrade =
              this.obj.subscription.subscription_plan.code != "free" ?
              false :
              true;

            $("#paymentModalRenew-" + id).modal("toggle");
          },
          error => {
            this.$parent.$parent.action = "updateBilling";
            this.set("upgrade_ref", id);
            this.$router.push({
              path: "/profile"
            });
            console.error(error);
          }
        );
      },
      hideRenewModal() {
        $("#paymentModalRenew-" + this.obj.id).modal("hide");
      },
      plansList() {
        this.subscriptionPlansGetAll(
          success => {
            this.plans = _.orderBy(success.body, "price", "asc");
          },
          error => {
            console.error(error);
          }
        );
      },
      changePlan(plan) {
        if (plan.code !== this.obj.subscription.subscription_plan.code) {
          // handle upgrade
          var context = this;
          this.calculateUpgradePrice(plan.code, function (data) {
            context.onUpgrade = true;
            context.currentPlan = plan;
            context.currentPlan.price = data.price;
            context.currentPlan.full_price = data.full_price;
            context.discounts.annualDiscount =
              Math.round(data.discount * 100) / 100;
            context.markdownDescription = marked(plan.description, {
              sanitize: true
            });
          });
        } else {
          this.onUpgrade = false;
          this.currentPlan = plan;
          this.currentPlan.full_price = plan.price;
          this.markdownDescription = marked(plan.description, {
            sanitize: true
          });
        }
      },
      calculateUpgradePrice(plan, callback) {
        this.onUpgradePriceCalc = true;

        this.subscriptionUpgradePriceCheck(
          plan,
          success => {
            this.onUpgradePriceCalc = false;
            callback(success.body);
          },
          error => {
            console.error(error);
            this.onUpgradePriceCalc = false;
          }
        );
      },
      calculateSubscription(date, subscription) {
        var moment = require("patternfly/node_modules/moment/moment.js");
        return moment(date, "YYYY-MM-DDTHH:mm:ss").add(
          subscription.period,
          "days"
        );
      },
      renewCheck(payment) {
        this.subscriptionRenew(
          payment,
          success => {
            this.payment.onProgress = false;
            this.payment.details = payment;
            this.payment.done = true;
            this.update();
          },
          error => {
            console.error(error);
            this.payment.onProgress = false;
            this.payment.details = payment;
            this.payment.done = true;
            this.errors.state = true;
            this.errors.message = error.body.message;
          }
        );
      },
      upgradeCheck(payment) {
        this.subscriptionUpgrade(
          payment,
          this.currentPlan,
          success => {
            this.payment.onProgress = false;
            this.payment.details = payment;
            this.payment.done = true;
            this.update();
          },
          error => {
            console.error(error);
            this.payment.onProgress = false;
            this.payment.details = payment;
            this.payment.done = true;
            this.errors.state = true;
            this.errors.message = error.body.message;
          }
        );
      }
    }
  };
</script>
<style>
</style>