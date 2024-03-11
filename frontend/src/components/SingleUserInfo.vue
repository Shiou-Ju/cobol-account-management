<template>
  <div class="content-container">
    <h1>User Information</h1>
    <p><strong>User:</strong> {{ userData.user }}</p>
    <p><strong>Balance:</strong> {{ userData.balance }}</p>
    <router-link :to="`/user/${userData.user}/transaction`"
      >Transact</router-link
    >

    <div v-if="transactions.length > 0" class="table-container">
      <h2>Nearest 10 Transactions</h2>
      <table>
        <thead>
          <tr>
            <th>Date</th>
            <th>Transaction</th>
            <th>Balance</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="transaction in transactions" :key="transaction.id">
            <td>{{ new Date(transaction.date).toLocaleString() }}</td>
            <td>{{ transaction.transaction }}</td>
            <td>{{ transaction.balance }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <p v-if="isLoading">Loading...</p>
    <p v-else-if="transactions.length === 0">No recent transactions found</p>
  </div>
</template>

<script lang="ts">
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useRoute } from 'vue-router';

interface UserData {
  user?: string;
  balance?: number;
}

interface Transaction {
  id: number;
  date: string;
  transaction: number;
  balance: number;
}

export default {
  name: 'SingleUserInfo',
  setup() {
    const userData = ref<UserData>({});
    const transactions = ref<Transaction[]>([]);
    const isLoading = ref(true);

    const route = useRoute();
    const userName = route.params.userName;

    onMounted(async () => {
      try {
        const userInfoResponse = await axios.get(`/api/user/${userName}`);
        const transactionResponse = await axios.get(
          `/api/user/${userName}/transactions`,
        );

        transactions.value = transactionResponse.data;
        userData.value = userInfoResponse.data;

        isLoading.value = false;
      } catch (err: unknown) {
        console.error('API call failed:', err);
        isLoading.value = false;
      }
    });

    return {
      userData,
      transactions,
      isLoading,
    };
  },
};
</script>

<style>
.content-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  margin: 0 auto;
}

.table-container {
  width: 70%;
  margin-top: 2rem;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th,
td {
  border: 0.1rem solid #ddd;
  padding: 0.2rem;
  text-align: left;
}

th {
  background-color: #f2f2f2;
}
</style>
