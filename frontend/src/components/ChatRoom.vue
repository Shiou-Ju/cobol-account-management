<template>
  <div>
    <ul>
      <li v-for="message in messages" :key="message.id">{{ message.text }}</li>
    </ul>
    <input v-model="inputMessage" @keyup.enter="sendMessage" />
    <button @click="sendMessage">Send</button>
  </div>
</template>

<!-- TODO: should not add `setup`, it return any type -->
<!-- <script lang="ts" setup> -->

<script lang="ts">
import { defineComponent, ref, onMounted } from 'vue';

export interface ChatMessage {
  id: number;
  text: string;
}

export default defineComponent({
  name: 'ChatRoom',
  setup() {
    const ws = ref<WebSocket | null>(null);
    const messages = ref<ChatMessage[]>([]);
    const inputMessage = ref<string>('');

    const connect = () => {
      ws.value = new WebSocket('ws://localhost:3001/ws');
      ws.value.onmessage = (event: MessageEvent) => {
        const message: ChatMessage = JSON.parse(event.data);
        messages.value.push(message);
      };
      ws.value.onerror = (error: Event) => {
        console.error('WebSocket Error:', error);
      };
    };

    const sendMessage = () => {
      if (inputMessage.value.trim() && ws.value) {
        ws.value.send(JSON.stringify({ message: inputMessage.value }));
        inputMessage.value = '';
      }
    };

    onMounted(() => {
      connect();
    });

    return {
      messages,
      inputMessage,
      sendMessage,
    };
  },
});
</script>
