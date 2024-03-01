<template>
  <div>
    <h1>Transaction Form</h1>
    <form @submit.prevent="submitTransaction">
      <p><strong>User:</strong> {{ userData.user }}</p>
      <p><strong>Balance:</strong> {{ userData.balance }}</p>

      <div class="transaction-container" v-if="userData.user">
        <div class="radio-group">
          <div>
            <input
              type="radio"
              id="deposit"
              value="deposit"
              v-model="transactionType"
            />
            <label for="deposit">Deposit</label>
          </div>

          <div>
            <input
              type="radio"
              id="withdraw"
              value="withdraw"
              v-model="transactionType"
            />
            <label for="withdraw">Withdraw</label>
          </div>
        </div>

        <div class="amount-container">
          <label for="amount">Amount:</label>
          <input
            id="amount"
            v-model.number="amount"
            type="number"
            placeholder="Enter amount"
          />
        </div>

        <button type="submit">Submit</button>
      </div>
    </form>
    <p v-if="error" class="error">{{ error }}</p>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useRoute } from 'vue-router';

export default {
  name: 'TransactionForm',
  setup() {
    const userData = ref({ user: '', balance: 0 });
    const transactionType = ref('deposit');
    const amount = ref(0);
    const error = ref('');
    const route = useRoute();

    function resetFields() {
      amount.value = 0;
      error.value = '';
    }

    onMounted(async () => {
      const userName = route.params.userName;
      try {
        const { data } = await axios.get(`/api/user/${userName}`);
        userData.value = data;
      } catch (error) {
        console.error('API call failed:', error);
        error.value = 'Failed to load user data.';
      }
    });

    const submitTransaction = async () => {
      try {
        const uppercaseTransactionType = transactionType.value.toUpperCase();

        const payload = {
          transaction:
            transactionType.value === 'deposit' ? amount.value : -amount.value,
          type: uppercaseTransactionType,
          user: userData.value.user,
        };

        const response = await axios.post('/api/transaction', payload);

        if (response.data.status === 'success') {
          const newBalance = response.data.newTransactionRecord.balance;

          userData.value.balance = newBalance;
          resetFields();
        }
      } catch (err) {
        console.error('Transaction failed:', err);
        error.value = err.response?.data || 'Transaction failed.';
      }
    };

    return {
      userData,
      transactionType,
      amount,
      submitTransaction,
      error,
    };
  },
};
</script>

<style>
.transaction-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  margin: 0 auto;
}

#amount,
button[type='submit'] {
  margin-top: 1rem;
  margin-bottom: 1rem;
  display: block;
}

.radio-group {
  margin-top: 2rem;
  display: flex;
  justify-content: space-around;
}

.amount-container {
  margin-top: 2rem;
  display: block;
}
</style>
