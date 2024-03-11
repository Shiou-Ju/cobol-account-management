<template>
  <div class="navbar-container">
    <nav class="navbar">
      <h1 class="navbar-title">FinChat</h1>
      <button @click="toggleNav" class="navbar-toggler" v-show="!isDesktop">
        &#9776;
      </button>
      <div v-show="isNavOpen || isDesktop" class="navbar-collapse">
        <router-link to="/" class="navbar-brand">Transaction</router-link>
        <router-link to="/select-user" class="navbar-brand"
          >Chat Room</router-link
        >
      </div>
    </nav>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, computed } from 'vue';

export default defineComponent({
  name: 'NavBar',
  setup() {
    const isNavOpen = ref(false);
    const screenWidth = ref(window.innerWidth);

    window.addEventListener('resize', () => {
      screenWidth.value = window.innerWidth;
    });

    const isDesktop = computed(() => screenWidth.value >= 768);

    function toggleNav() {
      isNavOpen.value = !isNavOpen.value;
    }

    return { isNavOpen, toggleNav, isDesktop };
  },
});
</script>

<style>
.navbar-container {
  max-width: 70%;
  margin: auto;
  margin-bottom: 0.5rem;
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

.navbar-toggler {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
}

.navbar-collapse {
  display: flex;
}

.navbar-brand {
  color: white;
  background-color: #007bff;
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: 0.25rem;
  margin-right: 1rem;
  transition: background-color 0.2s;
}

.navbar-brand:hover {
  background-color: #0056b3;
}

@media (max-width: 768px) {
  .navbar-collapse {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .navbar-collapse {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .navbar-brand {
    margin: 0 10px;
  }
}
</style>
