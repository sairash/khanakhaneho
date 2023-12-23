import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import {AuthStore} from '../stores/authStore'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: HomeView,
      meta: { requiresAuth: true },
    },
    {
      path: '/login',
      name: 'login',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/LoginView.vue'),
    }
  ]
})


router.beforeEach(async (to, from, next) => {
  const auth_store = AuthStore();
  const is_auth = await auth_store.IsAuth();
  console.log(is_auth);
  if (to.meta.requiresAuth && !is_auth) {
    next("/login")
  } else {
    next();
  }
});
export default router

