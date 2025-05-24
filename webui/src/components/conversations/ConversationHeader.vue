<!-- src/components/conversations/ConversationHeader.vue -->
<template>
  <div class="conversation-header">
    <div class="conversation-info">
      <button v-if="showBackButton" class="back-button" @click="goBack">
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
          <line x1="19" y1="12" x2="5" y2="12"></line>
          <polyline points="12 19 5 12 12 5"></polyline>
        </svg>
      </button>

      <div class="avatar-container">
        <div v-if="conversation.type == 'group'" class="group-avatar">
          <span>{{ getGroupInitial() }}</span>
        </div>
        <div v-else class="user-avatar">
          <img
            v-if="otherUser && otherUser.avatarUrl"
            :src="otherUser.avatarUrl"
            alt="User Avatar"
          />
          <div v-else class="avatar-placeholder">{{ getInitials() }}</div>
        </div>
      </div>

      <div class="conversation-details">
        <h2 class="conversation-name">{{ conversationName }}</h2>
        <div class="conversation-status">
          <span v-if="conversation.type == 'group'">
            {{ conversation.participants.length }} members
          </span>
        </div>
      </div>
    </div>

    <div class="conversation-actions">
      <button
        v-if="conversation.isGroup"
        class="action-button"
        @click="openGroupInfo"
        title="Group info"
      >
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
          <circle cx="12" cy="12" r="10"></circle>
          <line x1="12" y1="16" x2="12" y2="12"></line>
          <line x1="12" y1="8" x2="12.01" y2="8"></line>
        </svg>
      </button>

      <button class="action-button" @click="openSearchMessages" title="Search messages">
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
          <circle cx="11" cy="11" r="8"></circle>
          <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
        </svg>
      </button>

      <button class="action-button" @click="toggleMenu" title="More options">
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
          <circle cx="12" cy="12" r="1"></circle>
          <circle cx="12" cy="5" r="1"></circle>
          <circle cx="12" cy="19" r="1"></circle>
        </svg>
      </button>
    </div>

    <!-- Dropdown menu -->
    <div v-if="showMenu" class="dropdown-menu" ref="menuRef">
      <div v-if="conversation.isGroup" class="menu-item" @click="leaveGroup">
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
          <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
          <polyline points="16 17 21 12 16 7"></polyline>
          <line x1="21" y1="12" x2="9" y2="12"></line>
        </svg>
        <span>Leave group</span>
      </div>

      <div class="menu-item" @click="clearChat">
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
          <polyline points="3 6 5 6 21 6"></polyline>
          <path
            d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
          ></path>
        </svg>
        <span>Clear chat</span>
      </div>

      <div class="menu-item danger" @click="deleteConversation">
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
          <polyline points="3 6 5 6 21 6"></polyline>
          <path
            d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"
          ></path>
          <line x1="10" y1="11" x2="10" y2="17"></line>
          <line x1="14" y1="11" x2="14" y2="17"></line>
        </svg>
        <span>Delete conversation</span>
      </div>
    </div>

    <!-- Search messages overlay -->
    <div v-if="showSearchMessages" class="search-overlay">
      <div class="search-header">
        <div class="search-input-container">
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
            class="search-icon"
          >
            <circle cx="11" cy="11" r="8"></circle>
            <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
          </svg>
          <input
            type="text"
            v-model="searchQuery"
            placeholder="Search in conversation..."
            class="search-input"
            ref="searchInputRef"
          />
        </div>
        <button class="close-search" @click="closeSearchMessages">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="18"
            height="18"
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
  </div>
</template>

<script>
import { ref, computed, onMounted, onBeforeUnmount, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'

export default {
  name: 'ConversationHeader',
  props: {
    conversation: {
      type: Object,
      required: true,
    },
    showBackButton: {
      type: Boolean,
      default: true,
    },
  },
  emits: ['search-messages', 'clear-search'],
  setup(props, { emit }) {
    const authStore = useAuthStore()
    const router = useRouter()

    // State
    const showMenu = ref(false)
    const showSearchMessages = ref(false)
    const searchQuery = ref('')
    const menuRef = ref(null)
    const searchInputRef = ref(null)

    // Computed
    const currentUserId = authStore.user?.id

    const conversationName = computed(() => {
      if (props.conversation.type == 'group') {
        return props.conversation.name
      } else {
        // For direct conversations, show the other user's name
        return otherUser.value?.name || 'Unknown User'
      }
    })

    const otherUser = computed(() => {
      if (props.conversation.isGroup) {
        return null
      }

      return props.conversation.participants.find((member) => member.id !== currentUserId)
    })

    // Methods
    const goBack = () => {
      router.push('/conversations')
    }

    const getInitials = () => {
      if (!otherUser.value || !otherUser.value.name) return '?'

      return otherUser.value.name
        .split(' ')
        .map((word) => word.charAt(0).toUpperCase())
        .join('')
        .substring(0, 2)
    }

    const getGroupInitial = () => {
      if (!props.conversation.name) return '#'
      return props.conversation.name.charAt(0).toUpperCase()
    }

    const toggleMenu = () => {
      showMenu.value = !showMenu.value
    }

    const openGroupInfo = () => {
      // Implement navigation to group info page
      router.push(`/conversations/${props.conversation.id}/info`)
    }

    const openSearchMessages = () => {
      showSearchMessages.value = true
      // Focus on search input after DOM update
      nextTick(() => {
        searchInputRef.value?.focus()
      })
    }

    const closeSearchMessages = () => {
      showSearchMessages.value = false
      searchQuery.value = ''
      emit('clear-search')
    }

    const leaveGroup = async () => {
      if (confirm(`Are you sure you want to leave ${props.conversation.name}?`)) {
        try {
          await store.dispatch('conversations/leaveGroup', props.conversation.id)
          showMenu.value = false
          router.push('/conversations')
        } catch (error) {
          console.error('Failed to leave group:', error)
        }
      }
    }

    const clearChat = async () => {
      if (confirm('Are you sure you want to clear all messages? This cannot be undone.')) {
        try {
          await store.dispatch('conversations/clearMessages', props.conversation.id)
          showMenu.value = false
        } catch (error) {
          console.error('Failed to clear chat:', error)
        }
      }
    }

    const deleteConversation = async () => {
      if (confirm('Are you sure you want to delete this conversation? This cannot be undone.')) {
        try {
          await store.dispatch('conversations/deleteConversation', props.conversation.id)
          showMenu.value = false
          router.push('/conversations')
        } catch (error) {
          console.error('Failed to delete conversation:', error)
        }
      }
    }

    const handleClickOutside = (event) => {
      if (menuRef.value && !menuRef.value.contains(event.target) && showMenu.value) {
        showMenu.value = false
      }
    }

    // Watch for search query changes
    watch(searchQuery, (newValue) => {
      if (showSearchMessages.value) {
        emit('search-messages', newValue)
      }
    })

    // Lifecycle hooks
    onMounted(() => {
      document.addEventListener('click', handleClickOutside)
    })

    onBeforeUnmount(() => {
      document.removeEventListener('click', handleClickOutside)
    })

    return {
      showMenu,
      showSearchMessages,
      searchQuery,
      menuRef,
      searchInputRef,
      conversationName,
      otherUser,
      getInitials,
      getGroupInitial,
      goBack,
      toggleMenu,
      openGroupInfo,
      openSearchMessages,
      closeSearchMessages,
      leaveGroup,
      clearChat,
      deleteConversation,
    }
  },
}
</script>

<style scoped>
.conversation-header {
  position: relative;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid #e5e7eb;
  background-color: white;
  height: 64px;
}

.conversation-info {
  display: flex;
  align-items: center;
  overflow: hidden;
}

.back-button {
  background: none;
  border: none;
  padding: 0.5rem;
  margin-right: 0.5rem;
  cursor: pointer;
  color: #6b7280;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 0.375rem;
  transition: background-color 0.2s ease;
}

.back-button:hover {
  background-color: #f3f4f6;
}

.avatar-container {
  margin-right: 0.75rem;
}

.user-avatar,
.group-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  background-color: #3b82f6;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 1rem;
}

.group-avatar {
  background-color: #60a5fa;
  color: white;
  font-weight: 600;
  font-size: 1.125rem;
}

.conversation-details {
  overflow: hidden;
}

.conversation-name {
  font-size: 1rem;
  font-weight: 600;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: #111827;
}

.conversation-status {
  font-size: 0.75rem;
  color: #6b7280;
}

.online-status {
  color: #10b981;
}

.offline-status {
  color: #6b7280;
}

.member-count {
  color: #6b7280;
}

.conversation-actions {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.action-button {
  background: none;
  border: none;
  padding: 0.5rem;
  cursor: pointer;
  color: #6b7280;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 0.375rem;
  transition: background-color 0.2s ease;
}

.action-button:hover {
  background-color: #f3f4f6;
  color: #4b5563;
}

/* Dropdown menu */
.dropdown-menu {
  position: absolute;
  top: 60px;
  right: 12px;
  width: 220px;
  background-color: white;
  border: 1px solid #e5e7eb;
  border-radius: 0.5rem;
  box-shadow:
    0 4px 6px -1px rgba(0, 0, 0, 0.1),
    0 2px 4px -1px rgba(0, 0, 0, 0.06);
  z-index: 10;
  overflow: hidden;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 0.75rem 1rem;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.menu-item:hover {
  background-color: #f3f4f6;
}

.menu-item svg {
  margin-right: 0.75rem;
  color: #6b7280;
}

.menu-item.danger {
  color: #ef4444;
}

.menu-item.danger svg {
  color: #ef4444;
}

/* Search overlay */
.search-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  background-color: white;
  z-index: 20;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.search-header {
  display: flex;
  align-items: center;
  padding: 0.75rem;
  border-bottom: 1px solid #e5e7eb;
}

.search-input-container {
  position: relative;
  flex: 1;
}

.search-icon {
  position: absolute;
  left: 0.75rem;
  top: 50%;
  transform: translateY(-50%);
  color: #9ca3af;
}

.search-input {
  width: 100%;
  padding: 0.5rem 0.5rem 0.5rem 2.25rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 0.875rem;
}

.search-input:focus {
  outline: none;
  border-color: #4a6cf7;
  box-shadow: 0 0 0 2px rgba(74, 108, 247, 0.2);
}
</style>
