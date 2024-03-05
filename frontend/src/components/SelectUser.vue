<template>
  <div>
    <div class="button-container">
      <button @click="showModal = true" class="select-user-button">
        Select User
      </button>
    </div>
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
import { Ref, defineComponent, onMounted, onUnmounted, ref } from 'vue';
import axios from 'axios';
import { userState } from '../states/userState';
import router from '@/router';

type UserPickedStatus = 'available' | 'picked' | null;

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

const fetchSingleUserStatus = async (
  user: string,
): Promise<UserPickedStatus> => {
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

const mapUsersData = async (data: ReceivedUser[]): Promise<UserData[]> => {
  const usersWithStatusPromises = data.map(async (item) => {
    const status = await fetchSingleUserStatus(item.User);
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

const fetchAndMapAllUsers = async (users: Ref<UserData[]>) => {
  try {
    const response = await axios.get('http://localhost:3001/go-api/users');
    const mapped = await mapUsersData(response.data);
    users.value = mapped;
  } catch (error) {
    console.error('Failed to fetch users:', error);
  }
};

export default defineComponent({
  name: 'SelectUser',
  setup() {
    const users = ref<UserData[]>([]);
    const showModal = ref(true);
    const intervalId = ref<number | null>(null);

    onMounted(async () => {
      await fetchAndMapAllUsers(users);

      intervalId.value = window.setInterval(async () => {
        await fetchAndMapAllUsers(users);
      }, 5000);
    });

    onUnmounted(() => {
      if (intervalId.value !== null) {
        clearInterval(intervalId.value);
        intervalId.value = null;
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

          if (intervalId.value !== null) {
            clearInterval(intervalId.value);
            intervalId.value = null;
          }
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
  overflow: hidden;
  background-color: rgba(0, 0, 0, 0.4);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-content {
  background-color: #fefefe;
  margin: 15% auto;
  padding: 20px;
  border: 1px solid #888;
  width: 80%;
  max-width: 800px;
  min-width: 300px;
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
  list-style-type: none;
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

.button-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.select-user-button {
  padding: 10px 20px;
  background-color: brown;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.select-user-button:hover {
  background-color: purple;
}
</style>
