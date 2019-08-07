import Vue from 'vue'
import App from './App.vue'
import BootstrapVue from 'bootstrap-vue'
import VueRouter from 'vue-router'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import { library } from '@fortawesome/fontawesome-svg-core'
import { faEyeSlash, faEdit, faTrash } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

library.add(faEyeSlash, faEdit, faTrash)
Vue.component('font-awesome-icon', FontAwesomeIcon)

import Home from './components/Home.vue'
import Settings from './components/Settings.vue'

Vue.use(BootstrapVue)
Vue.use(VueRouter)

Vue.config.productionTip = false

const routes = [
  { path: '/', component: Home },
  { path: '/settings', component: Settings },
]

const router = new VueRouter({
  routes // short for `routes: routes`
})

new Vue({
  render: h => h(App),
  router: router
}).$mount('#app')
