<template>
  <div class="navbar-container">
    <nav class="navbar">
      <h1 class="navbar-title">FinChat</h1>
      <router-link to="/" class="navbar-brand">Transaction</router-link>
      <router-link to="/select-user" class="navbar-brand"
        >Chat Room</router-link
      >
    </nav>
  </div>
</template>

<script lang="ts">
export default {
  name: 'NavBar',
};
</script>

<style>
.navbar-container {
  max-width: 70%;
  margin: auto;
}

.navbar {
  background-color: beige;
  padding: 0 1rem;
  box-shadow: 0 0.1rem 0.2rem rgba(0, 0, 0, 0.1);
  text-align: left;
  width: 100%;
  display: flex;
  align-items: center;
}

.navbar-title {
  flex-grow: 1;
  color: #007bff;
  margin-right: auto;
  padding: 0.5rem 0;
}

.navbar-brand {
  color: white;
  background-color: #007bff;
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: 0.25rem;
  margin-right: 1rem;
  transition: background-color 0.2s;
  margin-left: 1rem;
}

.navbar-brand:hover {
  background-color: #0056b3;
}
</style>
