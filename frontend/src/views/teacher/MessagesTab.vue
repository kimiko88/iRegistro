<template>
  <div class="h-full flex gap-4 bg-base-100 rounded-box shadow overflow-hidden">
      <!-- Conversation List -->
      <div class="w-1/3 border-r border-base-200 flex flex-col">
          <div class="p-4 border-b border-base-200 bg-base-50">
              <input type="text" placeholder="Search messages..." class="input input-bordered input-sm w-full" />
              <button class="btn btn-primary btn-sm w-full mt-2">+ New Message</button>
          </div>
          <div class="overflow-y-auto flex-1">
              <div v-for="conv in conversations" :key="conv.id" 
                   class="p-4 border-b border-base-100 hover:bg-base-200 cursor-pointer transition-colors"
                   :class="{ 'bg-primary/5 border-l-4 border-l-primary': selectedConvId === conv.id }"
                   @click="selectConv(conv.id)"
              >
                  <div class="flex justify-between items-start mb-1">
                      <span class="font-bold text-sm">{{ conv.with }}</span>
                      <span class="text-xs opacity-50">{{ conv.time }}</span>
                  </div>
                  <div class="text-xs text-base-content/70 truncate line-clamp-1">
                      {{ conv.lastMessage }}
                  </div>
                  <div class="flex gap-1 mt-2">
                       <span v-if="conv.unread > 0" class="badge badge-error badge-xs">{{ conv.unread }}</span>
                       <span class="badge badge-ghost badge-xs">{{ conv.class }}</span>
                  </div>
              </div>
          </div>
      </div>

      <!-- Chat Area -->
      <div class="flex-1 flex flex-col bg-base-50/50">
          <div v-if="selectedConv" class="flex-1 flex flex-col">
               <!-- Header -->
               <div class="p-4 border-b border-base-200 bg-base-100 flex justify-between items-center shadow-sm z-10">
                   <div>
                       <h3 class="font-bold">{{ selectedConv.with }}</h3>
                       <span class="text-xs opacity-60">Parents of {{ selectedConv.studentName }} ({{ selectedConv.class }})</span>
                   </div>
                   <button class="btn btn-ghost btn-circle btn-sm">⋮</button>
               </div>
               
               <!-- Messages -->
               <div class="flex-1 overflow-y-auto p-4 space-y-4">
                   <div v-for="msg in selectedMessages" :key="msg.id" class="chat" :class="msg.sender === 'me' ? 'chat-end' : 'chat-start'">
                       <div class="chat-header text-xs opacity-50 mb-1">
                           {{ msg.senderName }} <time class="text-[10px] ml-1">{{ msg.time }}</time>
                       </div>
                       <div class="chat-bubble shadow-sm" :class="msg.sender === 'me' ? 'chat-bubble-primary' : 'bg-white text-base-content'">
                           {{ msg.text }}
                       </div>
                       <div v-if="msg.status" class="chat-footer opacity-50 text-[10px]">
                           {{ msg.status }}
                       </div>
                   </div>
               </div>

               <!-- Input -->
               <div class="p-4 bg-base-100 border-t border-base-200">
                   <div class="join w-full">
                       <button class="btn join-item btn-square btn-ghost">📎</button>
                       <input type="text" placeholder="Type a message..." class="input input-bordered join-item w-full" />
                       <button class="btn join-item btn-primary">Send</button>
                   </div>
               </div>
          </div>

          <div v-else class="flex-1 flex flex-col items-center justify-center opacity-40">
              <div class="text-6xl mb-4">💬</div>
              <p>Select a conversation to start chatting</p>
          </div>
      </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

// Mock Data
const conversations = ref([
    { id: 1, with: 'Mario Rossi (Parent)', studentName: 'Luigi Rossi', class: '1A', lastMessage: 'Buongiorno, volevo chiedere informazioni sulla gita.', time: '10:30', unread: 2 },
    { id: 2, with: 'Giulia Bianchi (Parent)', studentName: 'Marco Bianchi', class: '1A', lastMessage: 'Grazie mille professore!', time: 'Yesterday', unread: 0 },
    { id: 3, with: 'Classe 2B (Broadcast)', studentName: 'Tutti', class: '2B', lastMessage: 'Ricordo a tutti la verifica di domani.', time: 'Yesterday', unread: 0 },
]);

const selectedConvId = ref<number | null>(null);

const selectConv = (id: number) => {
    selectedConvId.value = id;
};

const selectedConv = computed(() => conversations.value.find(c => c.id === selectedConvId.value));

const selectedMessages = ref([
    { id: 1, sender: 'them', senderName: 'Mario Rossi', text: 'Buongiorno Prof, volevo sapere se Luigi deve portare il pranzo al sacco per la gita di domani.', time: '10:25' },
    { id: 2, sender: 'me', senderName: 'Me', text: 'Buongiorno Signor Rossi. Sì, il pranzo è al sacco come indicato nella circolare.', time: '10:28', status: 'Read' },
    { id: 3, sender: 'them', senderName: 'Mario Rossi', text: 'Perfetto, la ringrazio molto.', time: '10:30' }
]);
</script>
