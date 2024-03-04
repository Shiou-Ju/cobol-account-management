<template>
  <div class="chat-room">
    <ul class="message-list">
      <li v-for="message in messages" :key="message.id" class="message">
        <div class="message-header">
          <span class="user">{{ message.user }}</span>
          <span class="time">{{ message.time }}</span>
        </div>
        <p class="text">{{ message.text }}</p>
      </li>
    </ul>
  </div>
  <input v-model="inputMessage" @keyup.enter="sendMessage" />
  <button v-if="isConnected" @click="sendMessage">Send</button>
  <button v-if="!isConnected" disabled @click="sendMessage">
    chat function not ready
  </button>
</template>

<!-- TODO: should not add `setup`, it return any type -->
<!-- <script lang="ts" setup> -->
<script lang="ts">
import { defineComponent, ref, onMounted } from 'vue';

const unknownUserConst = 'unknown-user' as const;

export interface ChatMessage {
  // TODO: shall be unique id
  id: number;
  text: string;
  user: string;
  time: string;
}

interface ReceivedMessage {
  username: string;
  message: string;
  time: number;
}

export default defineComponent({
  name: 'ChatRoom',
  setup() {
    const ws = ref<WebSocket | null>(null);
    const messages = ref<ChatMessage[]>([]);
    const inputMessage = ref<string>('');
    const isConnected = ref<boolean>(false);

    const connect = () => {
      // TODO: wss://
      ws.value = new WebSocket('ws://localhost:3001/ws');

      ws.value.onopen = () => {
        console.log('WebSocket Connected');
        isConnected.value = true;
      };

      ws.value.onmessage = (event: MessageEvent) => {
        const receivedMessage = JSON.parse(event.data) as ReceivedMessage;

        const hasTimeReceived =
          receivedMessage.time !== 0 && receivedMessage.time;

        const formattedTime = hasTimeReceived
          ? new Date(receivedMessage.time).toISOString()
          : '沒有顯示時間';

        const newMessage: ChatMessage = {
          // TODO: use UUID
          id: Date.now(),
          text: receivedMessage.message,
          user: receivedMessage.username || unknownUserConst,
          time: formattedTime,
        };

        console.log(`message received from go`, newMessage);

        messages.value.push(newMessage);

        // TODO: maybe longer?
        // const messageCountLimit = 15;
        const messageCountLimit = 5;

        const messagesLength = messages.value.length;
        const isOverLimit = messagesLength > messageCountLimit;

        if (isOverLimit) {
          messages.value.splice(0, messagesLength - messageCountLimit);
        }
      };

      ws.value.onclose = () => {
        console.log('WebSocket Disconnected');
        isConnected.value = false;

        const reconnectDelay = 1000;

        setTimeout(connect, reconnectDelay);
      };

      ws.value.onerror = (error: Event) => {
        console.error('WebSocket Error:', error);
      };
    };

    const sendMessage = () => {
      const canSend = inputMessage.value.trim() && isConnected.value;

      if (canSend) {
        // TODO: add Username
        const reqBody = { message: inputMessage.value };

        ws.value && ws.value.send(JSON.stringify(reqBody));

        const resetInput = '';

        inputMessage.value = resetInput;
      }
    };

    onMounted(() => {
      connect();
    });

    return {
      messages,
      inputMessage,
      sendMessage,
      isConnected,
    };
  },
});
</script>

<style>
.chat-room {
  display: flex;
  flex-direction: column;
  max-width: 600px;
  margin: auto;
}

.message-list {
  list-style: none;
  padding: 0;
  overflow-y: auto;
  max-height: 400px;
}

.message {
  margin-bottom: 10px;
  padding: 10px;
  background-color: #f4f4f4;
  border-radius: 8px;
}

.message-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 5px;
}

.user {
  font-weight: bold;
}

.time {
  font-size: 0.8rem;
  color: #666;
}

.text {
  margin: 0;
}

.message-input {
  display: flex;
  margin-top: 10px;
}

.message-input input {
  flex-grow: 1;
  margin-right: 10px;
}

.message-input button {
  padding: 5px 10px;
}
</style>