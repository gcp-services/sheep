import Vue from 'vue'
import Router from 'vue-router'
import Status from '@/components/Status'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Status',
      component: Status
    }
  ]
})
