<!-- src/components/messages/MessageList.vue -->
<template>
  <div class="message-list" ref="messageListRef">
    <div v-if="isLoadingMessages && !messages.length" class="loading-state">
      <div class="loading-spinner"></div>
      <p>Loading messages...</p>
    </div>
    
    <div v-else-if="!messages.length" class="empty-state">
      <div class="empty-icon">
        <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
        </svg>
      </div>
      <p>No messages yet</p>
      <p class="help-text">Start the conversation by sending a message below.</p>
    </div>
    
    <template v-else>
      <div 
        v-if="hasMoreMessages" 
        class="load-more-container"
        :class="{ loading: isLoadingMore }"
      >
        <button 
          v-if="!isLoadingMessages" 
          @click="loadMore" 
          class="load-more-button"
        >
          Load more messages
        </button>
        <div v-else class="loading-spinner small"></div>
      </div>
      
      <div class="messages-container">
        <template v-for="(message, index) in messages" :key="message.id">
          <!-- Date separator -->
          <div 
            v-if="shouldShowDateSeparator(message, index)" 
            class="date-separator"
          >
            <span>{{ formatMessageDate(message.createdAt) }}</span>
          </div>
          
          <!-- Message item -->
          <MessageItem 
            :message="message"
            :isOwn="isOwnMessage(message)"
            :showAvatar="shouldShowAvatar(message, index)"
            @reaction="handleReaction"
            @reply="handleReply"
          />
        </template>
      </div>
      
      <!-- Bottom scroll button -->
      <button 
        v-if="showScrollToBottom" 
        @click="scrollToBottom" 
        class="scroll-bottom-button"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="6 9 12 15 18 9"></polyline>
        </svg>
      </button>
    </template>
    
    <!-- Reply box -->
    <div v-if="replyingTo" class="reply-box">
      <div class="reply-content">
        <div class="reply-indicator">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="9 14 4 9 9 4"></polyline>
            <path d="M20 20v-7a4 4 0 0 0-4-4H4"></path>
          </svg>
        </div>
        <div class="reply-text">
          <div class="reply-author">{{ replyingTo.sender.username }}</div>
          <div class="reply-message">{{ truncateText(replyingTo.content, 50) }}</div>
        </div>
      </div>
      <button @click="cancelReply" class="cancel-reply-button">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUpdated, watch } from 'vue';
import MessageItem from './MessageItem.vue';
import { useAuthStore } from '@/store/auth';

export default {
  name: 'MessageList',
  components: {
    MessageItem
  },
  props: {
    messages: {
      type: Array,
      required: true
    },
    isLoadingMessages: {
      type: Boolean,
      default: false
    },
    hasMoreMessages: {
      type: Boolean,
      default: false
    }
  },
  emits: ['load-more', 'reply-message'],
  setup(props, { emit }) {
    const messageListRef = ref(null);
    const isLoadingMore = ref(false);
    const showScrollToBottom = ref(false);
    const replyingTo = ref(null);
    const autoScrollToBottom = ref(true);
    const authStore = useAuthStore();
    
    // Computed
    const currentUserId = computed(() => authStore.user?.id);
    
    // Methods
    const loadMore = () => {
      isLoadingMore.value = true;
      emit('load-more');
      
      // Reset loading state after a timeout (in case the API call fails)
      setTimeout(() => {
        isLoadingMore.value = false;
      }, 5000);
    };
    
    const isOwnMessage = (message) => {
      return message.sender.id === currentUserId.value;
    };
    
    const shouldShowAvatar = (message, index) => {
      // Show avatar if it's the first message or if the previous message is from a different sender
      if (index === 0) return true;
      
      const prevMessage = props.messages[index - 1];
      return prevMessage.sender.id !== message.sender.id;
    };
    
    const shouldShowDateSeparator = (message, index) => {
      if (index === 0) return true;
      
      const prevMessage = props.messages[index - 1];
      const prevDate = new Date(prevMessage.createdAt).toDateString();
      const currentDate = new Date(message.createdAt).toDateString();
      
      return prevDate !== currentDate;
    };
    
    const formatMessageDate = (dateString) => {
      const date = new Date(dateString);
      const today = new Date();
      const yesterday = new Date(today);
      yesterday.setDate(yesterday.getDate() - 1);
      
      if (date.toDateString() === today.toDateString()) {
        return 'Today';
      } else if (date.toDateString() === yesterday.toDateString()) {
        return 'Yesterday';
      } else {
        return date.toLocaleDateString(undefined, { 
          year: 'numeric', 
          month: 'short', 
          day: 'numeric' 
        });
      }
    };
    
    const handleScroll = () => {
      if (!messageListRef.value) return;
      
      const { scrollTop, scrollHeight, clientHeight } = messageListRef.value;
      const scrolledToBottom = scrollHeight - scrollTop - clientHeight < 50;
      
      showScrollToBottom.value = !scrolledToBottom;
      
      // Update auto-scroll behavior
      if (scrolledToBottom) {
        autoScrollToBottom.value = true;
      } else {
        // Only disable auto-scroll if user scrolls up and there are already messages
        if (props.messages.length > 0 && scrollTop < scrollHeight - clientHeight - 100) {
          autoScrollToBottom.value = false;
        }
      }
    };
    
    const scrollToBottom = () => {
      if (!messageListRef.value) return;
      
      messageListRef.value.scrollTop = messageListRef.value.scrollHeight;