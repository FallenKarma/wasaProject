<!-- src/components/messages/MessageItem.vue -->
<template>
  <div class="message-wrapper" :class="{ 'own-message': isOwn }">
    <div class="message-item" :class="{ 'with-avatar': showAvatar }">
      <!-- User avatar for non-own messages -->
      <div v-if="showAvatar && !isOwn" class="avatar-container">
        <img
          v-if="message.sender?.avatar"
          :src="message.sender.avatar"
          :alt="message.sender.username"
          class="avatar"
        />
        <div v-else class="default-avatar">
          {{ getInitials(message.sender?.username) }}
        </div>
      </div>
      <div v-else-if="!isOwn" class="avatar-spacer"></div>

      <div class="message-content-wrapper">
        <!-- Sender name for non-own messages -->
        <div v-if="showAvatar && !isOwn" class="sender-name">
          {{ message.sender?.username }}
        </div>

        <!-- Reply indicator if the message is a reply -->
        <div v-if="message.replyTo" class="reply-indicator">
          <div class="reply-line"></div>
          <div class="replied-content">
            <span class="replied-user">{{ message.replyTo.sender?.username }}</span>
            <span class="replied-text">{{ truncateText(message.replyTo.content, 40) }}</span>
          </div>
        </div>

        <!-- Message content -->
        <div class="message-content" :class="{ 'is-deleted': message.isDeleted }">
          <!-- Message text -->
          <div v-if="!message.isDeleted" class="message-text">
            {{ message.content }}
          </div>
          <div v-else class="deleted-message">
            <span class="deleted-icon">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="12"
                height="12"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <polyline points="3 6 5 6 21 6"></polyline>
                <path
                  d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
                ></path>
              </svg>
            </span>
            Message deleted
          </div>

          <!-- Attachments -->
          <div v-if="message.attachments?.length" class="attachments">
            <div
              v-for="(attachment, index) in message.attachments"
              :key="index"
              class="attachment"
              @click="handleAttachmentClick(attachment)"
            >
              <div v-if="isImage(attachment)" class="image-attachment">
                <img :src="attachment.url" :alt="attachment.filename" />
              </div>
              <div v-else class="file-attachment">
                <div class="file-icon">
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
                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                    <polyline points="14 2 14 8 20 8"></polyline>
                    <line x1="16" y1="13" x2="8" y2="13"></line>
                    <line x1="16" y1="17" x2="8" y2="17"></line>
                    <polyline points="10 9 9 9 8 9"></polyline>
                  </svg>
                </div>
                <div class="file-details">
                  <div class="file-name">{{ attachment.filename }}</div>
                  <div class="file-size">{{ formatFileSize(attachment.size) }}</div>
                </div>
              </div>
            </div>
          </div>

          <!-- Message reactions -->
          <MessageReactions
            v-if="message.reactions && message.reactions.length > 0"
            :reactions="message.reactions"
            :messageId="message.id"
            @add-reaction="$emit('reaction', { messageId: message.id, reaction: $event })"
            @remove-reaction="$emit('reaction', { messageId: message.id, reaction: null })"
          />

          <!-- Message timestamp -->
          <div class="message-time">
            {{ formatMessageTime(message.createdAt) }}
            <span v-if="message.edited" class="edited-indicator">(edited)</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Message actions -->
    <div class="message-actions" v-if="!message.isDeleted && showActions">
      <button class="action-button emoji-button" @click="toggleEmojiPicker">
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
          <circle cx="12" cy="12" r="10"></circle>
          <path d="M8 14s1.5 2 4 2 4-2 4-2"></path>
          <line x1="9" y1="9" x2="9.01" y2="9"></line>
          <line x1="15" y1="9" x2="15.01" y2="9"></line>
        </svg>
      </button>
      <button class="action-button reply-button" @click="$emit('reply', message)">
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
          <polyline points="9 17 4 12 9 7"></polyline>
          <path d="M20 18v-2a4 4 0 0 0-4-4H4"></path>
        </svg>
      </button>
      <button v-if="isOwn" class="action-button more-button" @click="toggleMoreActions">
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
          <circle cx="12" cy="12" r="1"></circle>
          <circle cx="19" cy="12" r="1"></circle>
          <circle cx="5" cy="12" r="1"></circle>
        </svg>
      </button>

      <!-- More actions dropdown -->
      <div v-if="showMoreActions" class="more-actions-dropdown">
        <button class="dropdown-item" @click="editMessage">
          <span class="dropdown-icon">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="14"
              height="14"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
              <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
            </svg>
          </span>
          Edit
        </button>
        <button class="dropdown-item delete" @click="confirmDelete">
          <span class="dropdown-icon">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="14"
              height="14"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <polyline points="3 6 5 6 21 6"></polyline>
              <path
                d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
              ></path>
            </svg>
          </span>
          Delete
        </button>
      </div>
    </div>

    <!-- Emoji picker -->
    <div v-if="showEmojiPicker" class="emoji-picker">
      <div class="emoji-list">
        <button
          v-for="emoji in commonEmojis"
          :key="emoji"
          class="emoji-item"
          @click="addReaction(emoji)"
        >
          {{ emoji }}
        </button>
      </div>
    </div>

    <!-- Delete confirmation -->
    <div v-if="showDeleteConfirmation" class="delete-confirmation">
      <div class="confirmation-dialog">
        <div class="confirmation-title">Delete message?</div>
        <div class="confirmation-text">This cannot be undone.</div>
        <div class="confirmation-actions">
          <button class="cancel-button" @click="showDeleteConfirmation = false">Cancel</button>
          <button class="delete-button" @click="deleteMessage">Delete</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useStore } from 'vuex'
import MessageReactions from './MessageReactions.vue'

export default {
  name: 'MessageItem',
  components: {
    MessageReactions,
  },
  props: {
    message: {
      type: Object,
      required: true,
    },
    isOwn: {
      type: Boolean,
      default: false,
    },
    showAvatar: {
      type: Boolean,
      default: true,
    },
  },
  emits: ['reaction', 'reply'],
  setup(props) {
    const store = useStore()
    const showActions = ref(false)
    const showMoreActions = ref(false)
    const showEmojiPicker = ref(false)
    const showDeleteConfirmation = ref(false)

    // Common emojis
    const commonEmojis = ['👍', '❤️', '😂', '😮', '😢', '👏', '🎉', '🤔']

    // Methods
    const getInitials = (username) => {
      if (!username) return '?'
      return username
        .split(' ')
        .map((word) => word.charAt(0).toUpperCase())
        .slice(0, 2)
        .join('')
    }

    const formatMessageTime = (timestamp) => {
      if (!timestamp) return ''

      const date = new Date(timestamp)
      const now = new Date()
      const diffMs = now - date
      const diffMins = Math.floor(diffMs / 60000)
      const diffHours = Math.floor(diffMins / 60)
      const diffDays = Math.floor(diffHours / 24)

      if (diffMins < 1) {
        return 'Just now'
      } else if (diffMins < 60) {
        return `${diffMins}m ago`
      } else if (diffHours < 24) {
        return `${diffHours}h ago`
      } else if (diffDays === 1) {
        return 'Yesterday'
      } else if (diffDays < 7) {
        return `${diffDays}d ago`
      } else {
        return date.toLocaleDateString(undefined, {
          month: 'short',
          day: 'numeric',
          hour: '2-digit',
          minute: '2-digit',
        })
      }
    }

    const truncateText = (text, maxLength) => {
      if (!text) return ''
      return text.length > maxLength ? text.substring(0, maxLength) + '...' : text
    }

    const isImage = (attachment) => {
      if (!attachment || !attachment.mimeType) return false
      return attachment.mimeType.startsWith('image/')
    }

    const formatFileSize = (bytes) => {
      if (!bytes) return '0 B'
      const k = 1024
      const sizes = ['B', 'KB', 'MB', 'GB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
    }

    const handleAttachmentClick = (attachment) => {
      if (isImage(attachment)) {
        // Open image in a modal or lightbox
        // Implementation depends on your UI framework
      } else {
        // Download the file
        window.open(attachment.url, '_blank')
      }
    }

    const toggleMoreActions = () => {
      showMoreActions.value = !showMoreActions.value
      if (showMoreActions.value) {
        showEmojiPicker.value = false
      }
    }

    const toggleEmojiPicker = () => {
      showEmojiPicker.value = !showEmojiPicker.value
      if (showEmojiPicker.value) {
        showMoreActions.value = false
      }
    }

    const addReaction = (emoji) => {
      showEmojiPicker.value = false
      store.dispatch('messages/addReaction', { messageId: props.message.id, reaction: emoji })
    }

    const editMessage = () => {
      showMoreActions.value = false
      // Implement edit functionality
      // You might want to emit an event to parent component to handle this
    }

    const confirmDelete = () => {
      showMoreActions.value = false
      showDeleteConfirmation.value = true
    }

    const deleteMessage = () => {
      store.dispatch('messages/deleteMessage', props.message.id)
      showDeleteConfirmation.value = false
    }

    // Handle click outside to close dropdowns
    const handleClickOutside = (event) => {
      const isClickInsideMoreActions = event.target.closest('.more-button')
      const isClickInsideEmojiButton = event.target.closest('.emoji-button')

      if (!isClickInsideMoreActions && showMoreActions.value) {
        showMoreActions.value = false
      }

      if (!isClickInsideEmojiButton && showEmojiPicker.value) {
        showEmojiPicker.value = false
      }
    }

    // Mouse events
    const handleMouseEnter = () => {
      showActions.value = true
    }

    const handleMouseLeave = () => {
      showActions.value = false

      // Don't hide dropdowns when mouse leaves if they're active
      if (!showMoreActions.value && !showEmojiPicker.value) {
        showActions.value = false
      }
    }

    // Lifecycle hooks
    onMounted(() => {
      document.addEventListener('click', handleClickOutside)
    })

    onUnmounted(() => {
      document.removeEventListener('click', handleClickOutside)
    })

    return {
      showActions,
      showMoreActions,
      showEmojiPicker,
      showDeleteConfirmation,
      commonEmojis,
      getInitials,
      formatMessageTime,
      truncateText,
      isImage,
      formatFileSize,
      handleAttachmentClick,
      toggleMoreActions,
      toggleEmojiPicker,
      addReaction,
      editMessage,
      confirmDelete,
      deleteMessage,
      handleMouseEnter,
      handleMouseLeave,
    }
  },
}
</script>

<style scoped>
.message-wrapper {
  position: relative;
  margin-bottom: 0.5rem;
  padding: 0.25rem 0.5rem;
  border-radius: 0.5rem;
  transition: background-color 0.2s;
}

.message-wrapper:hover {
  background-color: #f9fafb;
}

.message-item {
  display: flex;
  position: relative;
}

.avatar-container {
  margin-right: 0.75rem;
  flex-shrink: 0;
}

.avatar-spacer {
  width: 2.25rem;
  margin-right: 0.75rem;
  flex-shrink: 0;
}

.avatar {
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 50%;
  object-fit: cover;
}

.default-avatar {
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 50%;
  background-color: #e5e7eb;
  color: #6b7280;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.875rem;
  font-weight: 500;
}

.message-content-wrapper {
  flex: 1;
  min-width: 0;
}

.sender-name {
  font-weight: 500;
  font-size: 0.875rem;
  margin-bottom: 0.25rem;
  color: #374151;
}

.message-content {
  position: relative;
  background-color: #f3f4f6;
  padding: 0.5rem 0.75rem;
  border-radius: 0.375rem;
  max-width: 85%;
  word-break: break-word;
}

.own-message .message-content {
  background-color: #dbeafe;
  margin-left: auto;
}

.message-text {
  white-space: pre-wrap;
  font-size: 0.9375rem;
  line-height: 1.5;
}

.message-time {
  font-size: 0.75rem;
  color: #6b7280;
  margin-top: 0.25rem;
  text-align: right;
}

.edited-indicator {
  font-size: 0.75rem;
  color: #9ca3af;
  margin-left: 0.25rem;
}

.deleted-message {
  color: #9ca3af;
  font-style: italic;
  font-size: 0.875rem;
  display: flex;
  align-items: center;
}

.deleted-icon {
  margin-right: 0.375rem;
  display: flex;
  align-items: center;
}

/* Reply styles */
.reply-indicator {
  display: flex;
  margin-bottom: 0.25rem;
  padding-left: 0.5rem;
  border-left: 2px solid #d1d5db;
}

.replied-content {
  font-size: 0.75rem;
  color: #6b7280;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.replied-user {
  font-weight: 500;
  color: #4b5563;
  margin-right: 0.25rem;
}

/* Attachments */
.attachments {
  margin-top: 0.5rem;
}

.attachment {
  margin-top: 0.5rem;
  cursor: pointer;
}

.image-attachment img {
  max-width: 100%;
  max-height: 200px;
  border-radius: 0.25rem;
  object-fit: contain;
}

.file-attachment {
  display: flex;
  align-items: center;
  padding: 0.5rem;
  background-color: rgba(255, 255, 255, 0.5);
  border-radius: 0.25rem;
  border: 1px solid #e5e7eb;
}

.file-icon {
  margin-right: 0.5rem;
  color: #6b7280;
}

.file-details {
  flex: 1;
  min-width: 0;
}

.file-name {
  font-size: 0.875rem;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-size {
  font-size: 0.75rem;
  color: #6b7280;
}

/* Message actions */
.message-actions {
  position: absolute;
  top: -0.75rem;
  right: 0.5rem;
  display: flex;
  background-color: white;
  border-radius: 0.375rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
  opacity: 0;
  transition: opacity 0.2s;
}

.message-wrapper:hover .message-actions {
  opacity: 1;
}

.action-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  border: none;
  background: none;
  color: #6b7280;
  cursor: pointer;
  border-radius: 0.25rem;
  transition: all 0.2s;
}

.action-button:hover {
  background-color: #f3f4f6;
  color: #374151;
}

/* Emoji picker */
.emoji-picker {
  position: absolute;
  top: -2.5rem;
  right: 0.5rem;
  background-color: white;
  border-radius: 0.375rem;
  box-shadow:
    0 4px 6px -1px rgba(0, 0, 0, 0.1),
    0 2px 4px -1px rgba(0, 0, 0, 0.06);
  border: 1px solid #e5e7eb;
  z-index: 10;
  padding: 0.5rem;
}

.emoji-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
}

.emoji-item {
  font-size: 1.125rem;
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: none;
  cursor: pointer;
  border-radius: 0.25rem;
  transition: all 0.2s;
}

.emoji-item:hover {
  background-color: #f3f4f6;
}

/* More actions dropdown */
.more-actions-dropdown {
  position: absolute;
  top: 2rem;
  right: 0;
  background-color: white;
  border-radius: 0.375rem;
  box-shadow:
    0 4px 6px -1px rgba(0, 0, 0, 0.1),
    0 2px 4px -1px rgba(0, 0, 0, 0.06);
  border: 1px solid #e5e7eb;
  z-index: 10;
  width: 8rem;
}

.dropdown-item {
  display: flex;
  align-items: center;
  padding: 0.5rem 0.75rem;
  font-size: 0.875rem;
  color: #374151;
  cursor: pointer;
  border: none;
  background: none;
  width: 100%;
  text-align: left;
  transition: all 0.2s;
}

.dropdown-item:hover {
  background-color: #f3f4f6;
}

.dropdown-item.delete {
  color: #ef4444;
}

.dropdown-item.delete:hover {
  background-color: #fef2f2;
}

.dropdown-icon {
  margin-right: 0.5rem;
  display: flex;
  align-items: center;
}

/* Delete confirmation */
.delete-confirmation {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 50;
}

.confirmation-dialog {
  background-color: white;
  border-radius: 0.5rem;
  box-shadow:
    0 10px 15px -3px rgba(0, 0, 0, 0.1),
    0 4px 6px -2px rgba(0, 0, 0, 0.05);
  padding: 1.25rem;
  width: 20rem;
  max-width: 90%;
}

.confirmation-title {
  font-size: 1rem;
  font-weight: 600;
  color: #111827;
  margin-bottom: 0.5rem;
}

.confirmation-text {
  font-size: 0.875rem;
  color: #6b7280;
  margin-bottom: 1rem;
}

.confirmation-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
}

.cancel-button {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  border-radius: 0.375rem;
  border: 1px solid #d1d5db;
  background-color: white;
  color: #374151;
  cursor: pointer;
  transition: all 0.2s;
}

.cancel-button:hover {
  background-color: #f9fafb;
}

.delete-button {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  border-radius: 0.375rem;
  border: none;
  background-color: #ef4444;
  color: white;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
}

.delete-button:hover {
  background-color: #dc2626;
}
</style>
