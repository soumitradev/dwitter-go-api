import { createRouter, createWebHistory } from 'vue-router'
import EditUserVue from '../views/EditUser.vue'
import Home from '../views/Home.vue'
import LoginVue from '../views/Login.vue'
import UserVue from '../views/ViewUser.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/about',
    name: 'About',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/About.vue')
  },
  {
    path: '/user/:id',
    name: 'ViewUser',
    component: UserVue,
  },
  {
    path: '/login',
    name: 'Login',
    component: LoginVue,
  },
  {
    path: '/edit',
    name: 'Edit Profile',
    component: EditUserVue,
  },
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
