import Vue from "vue";
import Router from "vue-router";
import VueAnalytics from "vue-analytics";

import Dashboard from "../components/Dashboard.vue";
import Hotspots from "../components/Hotspots.vue";
import HotspotsDetails from "../components/details-view/HotspotsDetails.vue";
import AccountsDetails from "../components/details-view/AccountsDetails.vue";
import SessionsDetails from "../components/details-view/SessionsDetails.vue";
import UnitsDetails from "../components/details-view/UnitsDetail.vue";
import Users from "../components/Users.vue";
import Sessions from "../components/Sessions.vue";
import Accounts from "../components/Accounts.vue";
import Profile from "../components/Profile.vue";
import Units from "../components/Units.vue";
import Reports from "../components/Reports.vue";
import Integration from "../components/Integration.vue";
import Marketing from "../components/Marketing.vue";
import Devices from "../components/Devices.vue";

Vue.use(Router);

const router = new Router({
  routes: [{
      path: "/",
      name: "Dashboard",
      component: Dashboard,
      meta: {
        roles: ["admin", "reseller", "customer", "desk"]
      }
    },
    {
      path: '/login/callback',
      name: 'LoginCallback',
      component: {
        template: '<div>Processing login...</div>',
        async created() {
          // Handle OIDC callback with secure code exchange
          const code = this.$route.query.code;

          if (code) {
            try {
              // Exchange the temporary code for token
              const response = await this.$http.post(
                this.$root.$options.api_scheme +
                this.$root.$options.api_host +
                "/api/auth/oidc/exchange",
                { code: code }
              );

              if (response.status === 200) {
                const data = response.data;

                // Create user object similar to regular login
                const loggedUser = {
                  token: data.token,
                  expires: new Date(parseInt(data.expires) * 1000).toString(),
                  id: parseInt(data.id),
                  account_type: data.account_type
                };

                // Save to localStorage
                localStorage.setItem('loggedUser', JSON.stringify(loggedUser));

                // Force page reload to trigger App.vue initialization with new login status
                window.location.href = '/';
              } else {
                // Exchange failed, redirect to login with error
                window.location.href = '/?error=exchange_failed';
              }
            } catch (error) {
              console.error('OIDC code exchange failed:', error);
              // Exchange failed, redirect to login with error
              window.location.href = '/?error=exchange_failed';
            }
          } else {
            // No code provided, something went wrong
            window.location.href = '/?error=missing_code';
          }
        }
      }
    },
    {
      path: "/hotspots",
      name: "Hotspots",
      component: Hotspots,
      meta: {
        roles: ["admin", "reseller"]
      }
    },
    {
      path: "/hotspots/:id",
      name: "HotspotsDetails",
      component: HotspotsDetails,
      meta: {
        roles: ["admin", "reseller", "customer", "desk"]
      }
    },
    {
      path: "/users",
      name: "Users",
      component: Users,
      meta: {
        roles: ["admin", "reseller", "customer", "desk"]
      }
    },
    {
      path: "/sessions",
      name: "Sessions",
      component: Sessions,
      meta: {
        roles: ["admin", "reseller", "customer"]
      }
    },
    {
      path: "/units",
      name: "Units",
      component: Units,
      meta: {
        roles: ["admin", "reseller", "customer"]
      }
    },
    {
      path: "/devices",
      name: "Devices",
      component: Devices,
      meta: {
        roles: ["admin", "reseller", "customer"]
      }
    },
    {
      path: "/units/:id",
      name: "UnitsDetails",
      component: UnitsDetails,
      meta: {
        roles: ["admin", "reseller", "customer"]
      }
    },
    {
      path: "/units",
      name: "Units",
      component: Units,
      meta: {
        roles: ["admin", "reseller", "customer"]
      }
    },
    {
      path: "/units/:id",
      name: "UnitsDetails",
      component: UnitsDetails,
      meta: {
        roles: ["admin", "reseller", "customer"]
      }
    },
    {
      path: "/sessions/:id",
      name: "SessionsDetails",
      component: SessionsDetails,
      meta: {
        roles: ["admin", "reseller", "customer"]
      }
    },
    {
      path: "/accounts",
      name: "Accounts",
      component: Accounts,
      meta: {
        roles: ["admin", "reseller"]
      }
    },
    {
      path: "/accounts/:id",
      name: "AccountsDetails",
      component: AccountsDetails,
      meta: {
        roles: ["admin", "reseller"]
      }
    },
    {
      path: "/profile",
      name: "Profile",
      component: Profile,
      meta: {
        roles: ["admin", "reseller", "customer", "desk"]
      }
    },
    {
      path: "/reports",
      name: "Reports",
      component: Reports,
      meta: {
        roles: ["admin", "reseller", "customer", "desk"]
      }
    },
    {
      path: "/marketing",
      name: "Marketing",
      component: Marketing,
      meta: {
        roles: ["admin", "reseller"]
      }
    },
    {
      path: "/integrations",
      name: "Integration",
      component: Integration,
      meta: {
        roles: ["admin", "reseller"]
      }
    }
  ]
});

router.beforeEach((to, from, next) => {
  // Always allow OIDC callback route
  if (to.name === 'LoginCallback') {
    next();
    return;
  }

  var user = JSON.parse(localStorage.getItem("loggedUser"));
  if (user && to.meta.roles && to.meta.roles.indexOf(user.account_type) >= 0) {
    next();
  } else {
    next(false);
  }
});

Vue.use(VueAnalytics, {
  id: CONFIG.UA_ANALYTICS,
  router
});

export default router;
