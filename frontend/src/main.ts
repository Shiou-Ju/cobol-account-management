import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import { apiBaseUrl, webSocketBaseUrl } from './config/config';

router.beforeEach((to, _from, next) => {
  if (!apiBaseUrl || !webSocketBaseUrl) {
    if (to.name !== 'ErrorPage') {
      next({ name: 'ErrorPage' });
    } else {
      next();
    }
  } else {
    next();
  }
});

createApp(App).use(router).mount('#app');
