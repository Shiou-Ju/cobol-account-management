import { reactive } from 'vue';

export const userState = reactive({
  isUserSelectedAndVerified: false,
  selectedUser: '',
});
