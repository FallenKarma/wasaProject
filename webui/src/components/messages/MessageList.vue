<!-- src/components/messages/MessageList.vue -->
<template>
  <div class="message-list" ref="messageListRef">
    <div v-if="isLoadingMessages && !messages.length" class="loading-state">
      <div class="loading-spinner"></div>
      <p>Loading messages...</p>
    </div>

    <div v-else-if="!messages.length" class="empty-state">
      <div class="empty-icon">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="40"
          height="40"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
        </svg>
      </div>
      <p>No messages yet</p>
      <p class="help-text">Start the conversation by sending a message below.</p>
    </div>

    <template v-else>
      <div class="messages-container">
        <template v-for="(message, index) in sortedMessages" :key="message.id">
          <!-- Date separator -->
          <div v-if="shouldShowDateSeparator(message, index)" class="date-separator">
            <span>{{ formatMessageDate(message.timestamp) }}</span>
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
      <button v-if="showScrollToBottom" @click="scrollToBottom" class="scroll-bottom-button">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <polyline points="6 9 12 15 18 9"></polyline>
        </svg>
      </button>
    </template>

    <!-- Reply box -->
    <div v-if="replyingTo" class="reply-box">
      <div class="reply-content">
        <div class="reply-indicator">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <polyline points="9 14 4 9 9 4"></polyline>
            <path d="M20 20v-7a4 4 0 0 0-4-4H4"></path>
          </svg>
        </div>
        <div class="reply-text">
          <div class="reply-author">{{ replyingTo.sender.name }}</div>
          <div class="reply-message">{{ truncateText(replyingTo.content, 50) }}</div>
        </div>
      </div>
      <button @click="cancelReply" class="cancel-reply-button">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUpdated, watch, nextTick } from 'vue'
import MessageItem from './MessageItem.vue'
import { useAuthStore } from '@/store/auth'

export default {
  name: 'MessageList',
  components: {
    MessageItem,
  },
  props: {
    messages: {
      type: Array,
      required: true,
    },
    isLoadingMessages: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['reply-message'],
  setup(props, { emit }) {
    const messageListRef = ref(null)
    const showScrollToBottom = ref(false)
    const replyingTo = ref(null)
    const autoScrollToBottom = ref(true)
    const authStore = useAuthStore()

    // Computed
    const currentUserId = computed(() => authStore.user?.id)

    // Sort messages to ensure newest messages are at the bottom
    const sortedMessages = computed(() => {
      if (!props.messages?.length) return []

      return [...props.messages].sort((a, b) => {
        const dateA = new Date(a.createdAt)
        const dateB = new Date(b.createdAt)
        return dateA - dateB // Ascending order (oldest first, newest at bottom)
      })
    })

    // Methods
    const isOwnMessage = (message) => {
      return message.sender.id === currentUserId.value
    }

    const shouldShowAvatar = (message, index) => {
      // Show avatar if it's the first message or if the previous message is from a different sender
      if (index === 0) return true

      const prevMessage = sortedMessages.value[index - 1]
      return prevMessage.sender.ID !== message.sender.ID
    }

    const shouldShowDateSeparator = (message, index) => {
      if (index === 0) return true

      const prevMessage = sortedMessages.value[index - 1]
      const prevDate = new Date(prevMessage.createdAt).toDateString()
      const currentDate = new Date(message.createdAt).toDateString()

      return prevDate !== currentDate
    }

    const formatMessageDate = (dateString) => {
      const date = new Date(dateString)
      const today = new Date()
      const yesterday = new Date(today)
      yesterday.setDate(yesterday.getDate() - 1)

      if (date.toDateString() === today.toDateString()) {
        return 'Today'
      } else if (date.toDateString() === yesterday.toDateString()) {
        return 'Yesterday'
      } else {
        return date.toLocaleDateString(undefined, {
          year: 'numeric',
          month: 'short',
          day: 'numeric',
        })
      }
    }

    const isScrolledToBottom = () => {
      if (!messageListRef.value) return false

      const { scrollTop, scrollHeight, clientHeight } = messageListRef.value
      return scrollHeight - scrollTop - clientHeight < 10 // Small threshold for precision
    }

    const handleScroll = () => {
      if (!messageListRef.value) return

      const scrolledToBottom = isScrolledToBottom()
      showScrollToBottom.value = !scrolledToBottom

      // Update auto-scroll behavior based on user's scroll position
      if (scrolledToBottom) {
        autoScrollToBottom.value = true
      } else {
        // Only disable auto-scroll if user intentionally scrolled up
        const { scrollTop, scrollHeight, clientHeight } = messageListRef.value
        if (scrollTop < scrollHeight - clientHeight - 100) {
          autoScrollToBottom.value = false
        }
      }
    }

    const scrollToBottom = (force = false) => {
      if (!messageListRef.value) return

      // Use nextTick to ensure DOM has been updated
      nextTick(() => {
        if (messageListRef.value && (autoScrollToBottom.value || force)) {
          messageListRef.value.scrollTop = messageListRef.value.scrollHeight
          autoScrollToBottom.value = true
          showScrollToBottom.value = false
        }
      })
    }

    const handleReaction = ({ messageId, reaction }) => {
      // TODO: Implement reaction handling
      console.log('Reaction:', messageId, reaction)
    }

    const handleReply = (message) => {
      replyingTo.value = message
      emit('reply-message', message)
    }

    const cancelReply = () => {
      replyingTo.value = null
      emit('reply-message', null)
    }

    const truncateText = (text, maxLength) => {
      if (!text) return ''
      return text.length > maxLength ? text.substring(0, maxLength) + '...' : text
    }

    // Lifecycle hooks
    onMounted(() => {
      if (messageListRef.value) {
        messageListRef.value.addEventListener('scroll', handleScroll)

        // Initial scroll to bottom after a short delay to ensure content is rendered
        setTimeout(() => {
          scrollToBottom(true)
        }, 100)
      }
    })

    onUpdated(() => {
      if (props.isLoadingMessages) return

      // Scroll to bottom for new messages if auto-scroll is enabled
      if (autoScrollToBottom.value) {
        scrollToBottom()
      }
    })

    // Watch for new messages
    let previousMessageCount = 0
    watch(
      () => sortedMessages.value?.length || 0,
      (newCount, oldCount) => {
        // Skip initial load
        if (oldCount === 0 && newCount > 0) {
          previousMessageCount = newCount
          nextTick(() => scrollToBottom(true))
          return
        }

        // Handle new messages
        if (newCount > previousMessageCount) {
          if (autoScrollToBottom.value) {
            scrollToBottom()
          } else {
            showScrollToBottom.value = true
          }
        }

        previousMessageCount = newCount
      },
      { immediate: true },
    )

    // Watch for changes in loading state
    watch(
      () => props.isLoadingMessages,
      (isLoading, wasLoading) => {
        // When loading finishes and we have messages, scroll to bottom if needed
        if (wasLoading && !isLoading && sortedMessages.value.length > 0) {
          if (autoScrollToBottom.value) {
            nextTick(() => scrollToBottom())
          }
        }
      },
    )

    return {
      messageListRef,
      showScrollToBottom,
      replyingTo,
      sortedMessages,
      isOwnMessage,
      shouldShowAvatar,
      shouldShowDateSeparator,
      formatMessageDate,
      scrollToBottom,
      handleReaction,
      handleReply,
      cancelReply,
      truncateText,
    }
  },
}
</script>

<style scoped>
.message-list {
  flex: 1;
  overflow-y: auto;
  padding: 1rem 1rem 0.5rem;
  display: flex;
  flex-direction: column;
  position: relative;
  scroll-behavior: smooth;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #6b7280;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #e5e7eb;
  border-top: 3px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #6b7280;
  text-align: center;
}

.empty-icon {
  color: #9ca3af;
  margin-bottom: 1rem;
}

.empty-state p {
  margin: 0;
  font-size: 0.875rem;
}

.empty-state .help-text {
  color: #9ca3af;
  margin-top: 0.5rem;
}

.messages-container {
  display: flex;
  flex-direction: column;
}

.date-separator {
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 1rem 0;
}

.date-separator span {
  background-color: #f3f4f6;
  color: #6b7280;
  font-size: 0.75rem;
  padding: 0.25rem 0.75rem;
  border-radius: 0.375rem;
}

.scroll-bottom-button {
  position: absolute;
  bottom: 1rem;
  right: 1rem;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background-color: #3b82f6;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border: none;
  transition: background-color 0.2s ease;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  z-index: 10;
}

.scroll-bottom-button:hover {
  background-color: #2563eb;
}

.reply-box {
  position: sticky;
  bottom: 0;
  background-color: #f9fafb;
  border-top: 1px solid #e5e7eb;
  padding: 0.75rem 1rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 0 -1rem -0.5rem;
}

.reply-content {
  display: flex;
  align-items: center;
  flex: 1;
}

.reply-indicator {
  color: #6b7280;
  margin-right: 0.5rem;
}

.reply-text {
  flex: 1;
}

.reply-author {
  font-size: 0.75rem;
  font-weight: 600;
  color: #374151;
}

.reply-message {
  font-size: 0.75rem;
  color: #6b7280;
  margin-top: 0.125rem;
}

.cancel-reply-button {
  background: none;
  border: none;
  color: #6b7280;
  cursor: pointer;
  padding: 0.25rem;
  border-radius: 0.25rem;
  transition: background-color 0.2s ease;
}

.cancel-reply-button:hover {
  background-color: #e5e7eb;
}
</style>
