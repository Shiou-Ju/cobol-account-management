<template>
  <div class="navbar-container">
    <nav class="navbar">
      <router-link to="/" class="navbar-brand">User List</router-link>
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
  padding: 0.5rem 1rem;
  box-shadow: 0 0.1rem 0.2rem rgba(0, 0, 0, 0.1);
  text-align: left;
  width: 100%;
}

.navbar-brand {
  color: white;
  background-color: #007bff;
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: 0.25rem;
  transition: background-color 0.2s;
}

.navbar-brand:hover {
  background-color: #0056b3;
}
</style>
