// global
import {
  NavigationGuardNext,
  RouteLocationNormalized,
  RouteRecordRaw,
  createRouter,
  createWebHistory,
} from 'vue-router';
import { userState } from '../states/userState';
// components
import AppHome from '@/components/AppHome.vue';
import SingleUserInfo from '@/components/SingleUserInfo.vue';
import TransactionForm from '@/components/TransactionForm.vue';
import ChatRoom from '@/components/ChatRoom.vue';
import SelectUser from '@/components/SelectUser.vue';
import axios from 'axios';

const routes: Array<RouteRecordRaw> = [
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
    path: '/select-user',
    name: 'SelectUser',
    component: SelectUser,
  },
  {
    path: '/chat',
    name: 'ChatRoom',
    component: ChatRoom,
    beforeEnter: (
      _to: RouteLocationNormalized,
      _from: RouteLocationNormalized,
      next: NavigationGuardNext,
    ) => {
      if (!userState.isUserSelectedAndVerified) {
        next('/select-user');
      } else {
        next();
      }
    },
  },
  {
    path: '/:catchAll(.*)',
    redirect: '/',
  },
];

const router = createRouter({
  // TODO:
  // eslint-disable-next-line no-undef
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

router.beforeEach(async (to, from, next) => {
  const isFromSelectUserToChatPage =
    to.path === '/chat' && from.path === '/select-user';

  if (isFromSelectUserToChatPage) {
    next();
  } else {
    try {
      await axios.post('http://localhost:3001/go-api/try-unlock-user', {
        username: userState.selectedUser,
      });

      next();
    } catch (error) {
      console.error('Error unlocking the user:', error);

      // TODO: to an error page
      next();
    }

    next();
  }
});

export default router;
