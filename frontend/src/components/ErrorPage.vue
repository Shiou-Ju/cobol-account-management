<template>
  <div class="error-page">
    <h1>Configuration Error</h1>
    <p>
      The application configuration is incorrect, missing necessary environment
      variables.
    </p>
    <p>
      Please ensure that VUE_APP_API_BASE_URL and VUE_APP_WEBSOCKET_BASE_URL are
      properly set.
    </p>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
  name: 'ErrorPage',
});
</script>

<style>
.error-page {
  padding: 20px;
  text-align: center;
}
</style>
