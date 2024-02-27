<template>
  <div>
    <h1>Transaction Form</h1>
    <form @submit.prevent="submitTransaction">
      <p><strong>User:</strong> {{ userData.user }}</p>
      <p><strong>Balance:</strong> {{ userData.balance }}</p>

      <div class="transaction-container">
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

        <!-- TODO: make btn prettier -->
        <button type="submit">Submit</button>
      </div>
    </form>
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
    const route = useRoute();

    onMounted(async () => {
      const userName = route.params.userName;
      try {
        const { data } = await axios.get(`/api/user/${userName}`);
        userData.value = data;
      } catch (error) {
        console.error('API call failed:', error);
      }
    });

    const submitTransaction = async () => {
      try {
        const payload = {
          user: userData.value.user,
          transaction:
            transactionType.value === 'deposit' ? amount.value : -amount.value,
        };

        // TODO: api needde
        await axios.post('/api/transaction', payload);
        // TODO: upadate user info
      } catch (error) {
        console.error('Transaction failed:', error);
      }
    };

    return {
      userData,
      transactionType,
      amount,
      submitTransaction,
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
