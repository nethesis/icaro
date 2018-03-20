import Vue from 'vue'
import Router from 'vue-router'

import Dashboard from '../components/Dashboard.vue'
import Hotspots from '../components/Hotspots.vue'
import HotspotsDetails from '../components/details-view/HotspotsDetails.vue'
import AccountsDetails from '../components/details-view/AccountsDetails.vue'
import SessionsDetails from '../components/details-view/SessionsDetails.vue'
import UnitsDetails from "../components/details-view/UnitsDetail.vue";
import Users from '../components/Users.vue'
import Reports from '../components/Reports.vue'
import Accounts from '../components/Accounts.vue'
import Profile from '../components/Profile.vue'
import Units from '../components/Units.vue'
import StorageService from "../services/storage"
import UtilService from "../services/util"

Vue.use(Router)

const router = new Router({
  routes: [{
      path: '/',
      name: 'Dashboard',
      component: Dashboard,
      meta: {
        roles: ['customer', 'reseller', 'desk', 'admin']
      },
    },
    {
      path: '/hotspots',
      name: 'Hotspots',
      component: Hotspots,
      meta: {
        roles: ['admin', 'reseller']
      },
    },
    {
      path: '/hotspots/:id',
      name: 'HotspotsDetails',
      component: HotspotsDetails,
      meta: {
        roles: ['admin', 'reseller', 'customer', 'desk']
      },
    },
    {
      path: '/users',
      name: 'Users',
      component: Users,
      meta: {
        roles: ['customer', 'reseller', 'desk', 'admin']
      },
    },
    {
      path: '/reports',
      name: 'Reports',
      component: Reports,
      meta: {
        roles: ['customer', 'reseller', 'admin']
      },
    },
    {
      path: '/units',
      name: 'Units',
      component: Units,
      meta: {
        roles: ['customer', 'reseller', 'admin']
      },
    },
    {
      path:'/units/:id',
      name:'UnitsDetails',
      component: UnitsDetails,
      meta: {
        roles: ['customer', 'reseller', 'admin']
      },
    },
    {
      path: '/sessions/:id',
      name: 'SessionsDetails',
      component: SessionsDetails,
      meta: {
        roles: ['customer', 'reseller', 'admin']
      },
    },
    {
      path: '/accounts',
      name: 'Accounts',
      component: Accounts,
      meta: {
        roles: ['admin', 'reseller']
      },
    },
    {
      path: '/accounts/:id',
      name: 'AccountsDetails',
      component: AccountsDetails,
      meta: {
        roles: ['admin', 'reseller']
      },
    },
    {
      path: '/profile',
      name: 'Profile',
      component: Profile,
      meta: {
        roles: ['customer', 'reseller', 'desk', 'admin']
      },
    },
  ]
})

router.beforeEach((to, from, next) => {

  var user = JSON.parse(localStorage.getItem("loggedUser"))
  if (user && to.meta.roles && to.meta.roles.indexOf(user.account_type) >= 0) {
    next()
  } else {
    next(false)
  }

})


export default router
