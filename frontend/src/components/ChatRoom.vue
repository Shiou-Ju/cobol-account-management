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
  <div class="message-input-area">
    <input
      v-if="isConnected"
      v-model="inputMessage"
      @keyup.enter="sendMessage"
      class="message-input"
    />
    <p v-if="!isConnected">Please refresh page</p>
    <button v-if="isConnected" @click="sendMessage" class="send-button">
      Send
    </button>
    <button v-if="!isConnected" disabled @click="sendMessage">
      chat function not ready
    </button>
  </div>
</template>

<!-- TODO: should not add `setup`, it return any type -->
<!-- <script lang="ts" setup> -->
<script lang="ts">
import { defineComponent, ref, onMounted } from 'vue';
import { v4 as uuidv4 } from 'uuid';
import { userState } from '../states/userState';
import axios from 'axios';
import router from '../router';
import { apiBaseUrl, webSocketBaseUrl } from '../config/config';

const defaultUnknownUserDisplayName = 'unknown-user' as const;

export interface ChatMessage {
  id: string;
  text: string;
  user: string;
  time: string;
}

interface ReceivedMessage {
  username: string;
  message: string;
  time: number;
}

type PublishMessageRequestBody = {
  username: string;
  message: string;
  time: number;
};

export default defineComponent({
  name: 'ChatRoom',
  setup() {
    const ws = ref<WebSocket | null>(null);
    const messages = ref<ChatMessage[]>([]);
    const inputMessage = ref<string>('');
    const isConnected = ref<boolean>(false);

    let isSendingTooFast = false;

    const connect = () => {
      // TODO: wss:// ?
      ws.value = new WebSocket(`${webSocketBaseUrl}/go-api/ws`);

      ws.value.onopen = () => {
        console.log('WebSocket Connected');
        isConnected.value = true;
      };

      ws.value.onmessage = (event: MessageEvent) => {
        const parsedDate = JSON.parse(event.data);

        const isSystemInfo = parsedDate.isMessage === 'false';

        console.log(event.data);

        if (isSystemInfo) {
          (async () => {
            try {
              const connection = parsedDate.connection;

              const response = await axios.post(
                `${apiBaseUrl}/go-api/add-connection-to-user`,
                {
                  username: userState.selectedUser,
                  connection: connection,
                },
              );
              console.log('Connection info saved:', response.data);
            } catch (error) {
              console.error('Error saving connection info:', error);

              console.log('return to user select');
              router.push('/chat');
            }
          })();
          return;
        }

        const receivedMessage = parsedDate as ReceivedMessage;

        const hasTimeReceived =
          receivedMessage.time !== 0 && receivedMessage.time;

        const formattedTime = hasTimeReceived
          ? new Date(receivedMessage.time).toISOString()
          : '沒有顯示時間';

        const newMessage: ChatMessage = {
          id: uuidv4(),
          text: receivedMessage.message,
          user: receivedMessage.username || defaultUnknownUserDisplayName,
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

        userState.isUserSelectedAndVerified = false;
        userState.selectedUser = '';

        // TODO: consider reconnect
        // const reconnectDelay = 1000;
        // setTimeout(connect, reconnectDelay);
      };

      ws.value.onerror = (error: Event) => {
        console.error('WebSocket Error:', error);
      };
    };

    const sendMessage = () => {
      const canSend = inputMessage.value.trim() && isConnected.value;

      const username = userState.selectedUser;

      if (canSend && !isSendingTooFast) {
        isSendingTooFast = true;

        const reqBody: PublishMessageRequestBody = {
          message: inputMessage.value,
          username,
          time: Date.now(),
        };

        ws.value && ws.value.send(JSON.stringify(reqBody));

        const resetInput = '';

        inputMessage.value = resetInput;

        const sendingGap = 100;
        setTimeout(() => {
          isSendingTooFast = false;
        }, sendingGap);
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

  async beforeRouteLeave(_to, _from, next) {
    try {
      await axios.post(`${apiBaseUrl}/go-api/try-unlock-user`, {
        username: userState.selectedUser,
      });
      next();
    } catch (error) {
      console.error('Error unlocking the user:', error);
      next();
    }
  },
});
</script>

<style>
.chat-room {
  display: flex;
  flex-direction: column;
  max-width: 600px;
  margin: auto;
  border: 1px solid #ccc;
  padding-bottom: 10px;
}

.message-list {
  list-style: none;
  padding: 0;
  overflow-y: auto;
  min-height: 400px;
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
  flex-grow: 1;
  border: 2px solid #007bff;
  border-radius: 20px;
  padding: 5px 15px;
  margin-right: 10px;
  margin-top: 15px;
}

.send-button:hover {
  background-color: #0056b3;
}

.message-input input {
  flex-grow: 1;
  margin-right: 10px;
}

.message-input button {
  padding: 5px 10px;
}

.send-button {
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 20px;
  padding: 5px 15px;
  cursor: pointer;
}
</style>
