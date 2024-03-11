<template>
  <div>
    <h1>User List</h1>
    <div class="users-container">
      <div class="users-list">
        <ul>
          <li v-for="user in users" :key="user.user">
            <router-link :to="`/user/${user.user}`">
              <p>
                <strong>User: </strong>
                <span class="user-name">{{ user.user }}</span>
              </p>
            </router-link>
            <p class="sub-desc">
              <strong>Last Transaction:</strong>
              {{ new Date(user.date).toLocaleString() }}
            </p>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { ref, onMounted } from 'vue';
import axios from 'axios';

interface User {
  user: string;
  date: string;
}

export default {
  name: 'UsersList',
  setup() {
    const users = ref<User[]>([]);

    onMounted(async () => {
      try {
        const response = await axios.get('/api/users');
        users.value = response.data;
      } catch (err: unknown) {
        console.error('API call failed:', err);
      }
    });

    return {
      users,
    };
  },
};
</script>

<style scoped>
.users-container {
  text-align: center;
}

.users-list {
  text-align: left;
  display: inline-block;
  text-align: left;
}

.sub-desc {
  font-size: 0.8em;
}
</style>
