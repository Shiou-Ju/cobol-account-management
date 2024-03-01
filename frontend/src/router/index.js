import { createRouter, createWebHistory } from 'vue-router';
import AppHome from '@/components/AppHome.vue';
import SingleUserInfo from '@/components/SingleUserInfo.vue';
import TransactionForm from '@/components/TransactionForm.vue';

const routes = [
  {
    path: '/',
    name: 'Home',
    component: AppHome,
  },
  {
    path: '/user/:userName',
    name: 'SingleUser',
    component: SingleUserInfo,
    props: true,
  },
  {
    path: '/user/:userName/transaction',
    name: 'UserTransaction',
    component: TransactionForm,
    props: true,
  },
  {
    path: '/:catchAll(.*)',
    redirect: '/',
  },
];

const router = createRouter({
  // eslint-disable-next-line no-undef
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

export default router;
