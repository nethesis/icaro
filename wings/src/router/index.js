import Vue from 'vue'
import Router from 'vue-router'
import SplashPage from '@/components/SplashPage'
import LoginPage from '@/components/LoginPage'
import FacebookPage from '@/components/social/FacebookPage'
import LinkedInPage from '@/components/social/LinkedInPage'
import SMSPage from '@/components/others/SMSPage'
import EmailPage from '@/components/others/EmailPage'

Vue.use(Router)

export default new Router({
  mode: 'history',
  base: '/wings/',
  routes: [
    {
      path: '/',
      name: 'SplashPage',
      component: SplashPage
    },
    {
      path: '/login',
      name: 'LoginPage',
      component: LoginPage
    },
    {
      path: '/login/facebook',
      name: 'FacebookPage',
      component: FacebookPage
    },
    {
      path: '/login/linkedin',
      name: 'LinkedInPage',
      component: LinkedInPage
    },
    {
      path: '/login/sms',
      name: 'SMSPage',
      component: SMSPage
    },
    {
      path: '/login/email',
      name: 'EmailPage',
      component: EmailPage
    }
  ]
})
