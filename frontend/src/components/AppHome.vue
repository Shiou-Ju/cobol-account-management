<template>
  <div>
    <h1>User Information</h1>
    <p><strong>User:</strong> {{ data.user }}</p>
    <p><strong>Transaction:</strong> {{ data.transaction }}</p>
    <p><strong>Balance:</strong> {{ data.balance }}</p>
    <p><strong>Date:</strong> {{ new Date(data.date).toLocaleString() }}</p>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import axios from 'axios';

export default {
  name: 'AppHome',
  setup() {
    const data = ref({});

    onMounted(async () => {
      try {
        const response = await axios.get('/api/user/hunyo');
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
