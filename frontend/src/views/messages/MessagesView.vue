<script setup lang="ts">
import { ref } from 'vue';

// Mock data for messages
const messages = ref([
    { id: 1, sender: 'Prof. Bianchi', subject: 'Math Homework', date: '2023-10-25', preview: 'Please remember to complete...' },
    { id: 2, sender: 'School Office', subject: 'Holiday closing', date: '2023-10-20', preview: 'The school will be closed...' },
]);

const selectedMessage = ref<any>(null);

function selectMessage(msg: any) {
    selectedMessage.value = msg;
}
</script>

<template>
  <div class="p-4 flex gap-4 h-[calc(100vh-100px)]">
      <!-- Message List -->
      <div class="w-1/3 bg-base-100 rounded-box shadow overflow-y-auto">
          <ul class="menu p-2 bg-base-100 w-full">
              <li v-for="msg in messages" :key="msg.id" @click="selectMessage(msg)" class="border-b last:border-0">
                  <a :class="{ 'active': selectedMessage?.id === msg.id }">
                      <div class="flex flex-col gap-1 w-full">
                          <div class="flex justify-between font-bold">
                              <span>{{ msg.sender }}</span>
                              <span class="text-xs font-normal opacity-70">{{ msg.date }}</span>
                          </div>
                          <div class="text-sm truncate">{{ msg.subject }}</div>
                          <div class="text-xs opacity-70 truncate">{{ msg.preview }}</div>
                      </div>
                  </a>
              </li>
          </ul>
      </div>

      <!-- Message Detail -->
      <div class="w-2/3 bg-base-100 rounded-box shadow p-6 relative">
          <div v-if="selectedMessage">
              <div class="border-b pb-4 mb-4">
                  <h2 class="text-2xl font-bold mb-2">{{ selectedMessage.subject }}</h2>
                  <div class="flex justify-between items-center text-sm opacity-70">
                      <span>From: {{ selectedMessage.sender }}</span>
                      <span>{{ selectedMessage.date }}</span>
                  </div>
              </div>
              <div class="prose">
                  <p>
                      {{ selectedMessage.preview }}
                      <br><br>
                      [Full message content placeholder...]
                  </p>
              </div>
              <div class="absolute bottom-4 right-4">
                  <button class="btn btn-primary" @click="alert('Reply feature coming soon')">Reply</button>
              </div>
          </div>
          <div v-else class="flex items-center justify-center h-full text-lg opacity-50">
              Select a message to read
          </div>
      </div>
  </div>
</template>
