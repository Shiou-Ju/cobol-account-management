<template>
  <div>
    <button @click="showModal = true">Select User</button>
    <div v-if="showModal" class="modal">
      <div class="modal-content">
        <span class="close" @click="showModal = false">&times;</span>
        <h1>Select a user to join the chat!</h1>
        <ul>
          <li
            v-for="user in users"
            :key="user.user"
            @click="selectUserHandler(user)"
            :class="{ picked: user.status === 'picked' }"
          >
            <p>User: {{ user.user }}</p>
            <p>Balance: {{ user.balance }}</p>
            <p>Status: {{ user.status }}</p>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, ref } from 'vue';
import axios from 'axios';
import { userState } from '../states/userState';
import router from '@/router';

type UserPickedStatus = 'available' | 'picked' | null;

// TODO: new api?
interface UserData {
  date: string;
  transaction: number;
  balance: number;
  user: string;
  status: UserPickedStatus;
}

interface ReceivedUser {
  Date: string;
  Transaction: number;
  Balance: number;
  User: string;
}

const fetchUserStatus = async (user: string): Promise<UserPickedStatus> => {
  try {
    const response = await axios.get(
      `http://localhost:3001/go-api/user-state?username=${user}`,
    );
    return response.data.status;
  } catch (error) {
    console.error('Failed to fetch user status:', error);
    return null;
  }
};

const mapUserData = async (data: ReceivedUser[]): Promise<UserData[]> => {
  const usersWithStatusPromises = data.map(async (item) => {
    const status = await fetchUserStatus(item.User);
    return {
      date: item.Date,
      transaction: item.Transaction,
      balance: item.Balance,
      user: item.User,
      status,
    };
  });
  return Promise.all(usersWithStatusPromises);
};

export default defineComponent({
  name: 'SelectUser',
  setup() {
    const users = ref<UserData[]>([]);
    // TODO: default to show modal
    const showModal = ref(true);
    // TODO: shall be fetch from backend
    const userSelected = ref('');

    onMounted(async () => {
      try {
        const response = await axios.get('http://localhost:3001/go-api/users');

        const mapped = await mapUserData(response.data);
        users.value = mapped;
      } catch (error) {
        console.error('Failed to fetch users:', error);
      }
    });

    async function selectUserHandler(user: UserData) {
      if (user.status === 'picked') {
        // TODO: show banner here
        console.log('This user is already picked and cannot be selected.');
        return;
      }

      try {
        const lockResponse = await axios.post(
          'http://localhost:3001/go-api/try-lock-user',
          {
            username: user.user,
          },
        );

        if (lockResponse.status === 200) {
          console.log('User locked successfully:', lockResponse.data);

          userState.selectedUser = user.user;
          userState.isUserSelectedAndVerified = true;
          showModal.value = false;

          router.push('/chat');
        } else {
          console.error(
            'Failed to lock user, status code:',
            lockResponse.status,
          );
        }
      } catch (error) {
        if (axios.isAxiosError(error)) {
          console.error(
            'Error locking user:',
            error.response ? error.response.data : error.message,
          );

          console.log(error);

          const hasBeenLocked = error.response?.status;

          if (hasBeenLocked === 409) {
            console.log('User is already locked, refreshing page...');
            // TODO: refresh paged is too heavy maybe
            location.reload();
          }
        } else {
          console.error('unknown error locking user');
        }
      }
    }
    return {
      users,
      selectUserHandler,
      showModal,
      userSelected,
    };
  },
});
</script>

<style scoped>
.modal {
  position: fixed;
  z-index: 1;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  overflow: auto;
  background-color: rgba(0, 0, 0, 0.4);
}

.modal-content {
  background-color: #fefefe;
  margin: 15% auto;
  padding: 20px;
  border: 1px solid #888;
  width: 80%;
  max-width: 800px;
}

.close {
  color: #aaa;
  float: right;
  font-size: 28px;
  font-weight: bold;
}

.close:hover,
.close:focus {
  color: black;
  text-decoration: none;
  cursor: pointer;
}

li {
  cursor: pointer;
  background-color: #f0f0f0;
  padding: 10px;
  margin: 10px auto;
  border-radius: 5px;
  transition: background-color 0.3s ease;
  width: 80%;
}

li:hover:not(.disabled-user) {
  background-color: #e2e2e2;
}

.disabled-user {
  cursor: not-allowed;
  opacity: 0.5;
}

.picked {
  background-color: #ffcccc;
  cursor: not-allowed;
  opacity: 0.6;
}
</style>
