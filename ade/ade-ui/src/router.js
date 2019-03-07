import Vue from 'vue'
import Router from 'vue-router'

import Home from './views/Home.vue'
import Feedbacks from './views/Feedbacks.vue'
import Reviews from './views/Reviews.vue'

import NotFound from './views/NotFound.vue'

Vue.use(Router)

export default new Router({
  routes: [{
      path: '/',
      redirect: '/not-found',
    },
    {
      path: '/feedbacks',
      redirect: '/not-found',
    },
    {
      path: '/feedbacks/:token',
      name: 'feedbacks',
      component: Feedbacks
    },
    {
      path: '/reviews',
      redirect: '/not-found',
    },
    {
      path: '/reviews/:token',
      name: 'reviews',
      component: Reviews
    },
    {
      path: '/not-found',
      name: 'not-found',
      component: NotFound
    },
    {
      path: '*',
      name: 'not-found',
      component: NotFound
    }
  ]
})