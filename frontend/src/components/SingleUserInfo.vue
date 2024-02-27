<template>
  <div>
    <h1>User Information</h1>
    <p><strong>User:</strong> {{ data.user }}</p>
    <!-- <p><strong>Transaction:</strong> {{ data.transaction }}</p> -->
    <p><strong>Balance:</strong> {{ data.balance }}</p>
    <!-- <p><strong>Date:</strong> {{ new Date(data.date).toLocaleString() }}</p> -->
    <router-link :to="`/user/${data.user}/transaction`">Transact</router-link>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useRoute } from 'vue-router';

export default {
  name: 'SingleUserInfo',
  setup() {
    const data = ref({});

    const route = useRoute();
    const userName = route.params.userName;

    onMounted(async () => {
      try {
        const response = await axios.get(`/api/user/${userName}`);

        data.value = response.data;
      } catch (error) {
        console.error('API call failed:', error);
      }
    });

    return {
      data,
    };
  },
};
</script>
