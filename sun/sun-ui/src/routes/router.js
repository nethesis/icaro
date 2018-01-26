import Vue from 'vue'
import Router from 'vue-router'
import Dashboard from '../components/Dashboard.vue'
import Hotspots from '../components/Hotspots.vue'
import Users from '../components/Users.vue'
import Reports from '../components/Reports.vue'
import Accounts from '../components/Accounts.vue'
import Preferences from '../components/Preferences.vue'
import Profile from '../components/Profile.vue'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Dashboard',
      component: Dashboard
    },
    {
      path: '/hotspots',
      name: 'Hotspots',
      component: Hotspots
    },
    {
      path: '/users',
      name: 'Users',
      component: Users
    },
    {
      path: '/reports',
      name: 'Reports',
      component: Reports
    },
    {
      path: '/accounts',
      name: 'Accounts',
      component: Accounts
    },
    {
      path: '/preferences',
      name: 'Preferences',
      component: Preferences
    },
    {
      path: '/profile',
      name: 'Profile',
      component: Profile
    },
  ]
})
