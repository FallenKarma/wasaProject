<!-- src/components/messages/MessageReactions.vue -->
<template>
  <div class="message-reactions">
    <div class="reactions-list">
      <div
        v-for="reaction in groupedReactions"
        :key="reaction.emoji"
        class="reaction"
        :class="{ 'user-reacted': isUserReaction(reaction) }"
        @click="toggleReaction(reaction.emoji)"
      >
        <span class="reaction-emoji">{{ reaction.emoji }}</span>
        <span class="reaction-count">{{ reaction.count }}</span>
      </div>
    </div>
    <div v-if="showReactionSelector" class="reaction-selector">
      <div class="selector-content">
        <button
          v-for="emoji in commonEmojis"
          :key="emoji"
          class="emoji-button"
          @click="addReaction(emoji)"
        >
          {{ emoji }}
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import { useStore } from 'vuex'

export default {
  name: 'MessageReactions',
  props: {
    reactions: {
      type: Array,
      required: true,
      default: () => [],
    },
    messageId: {
      type: [String, Number],
      required: true,
    },
  },
  emits: ['add-reaction', 'remove-reaction'],
  setup(props, { emit }) {
    const store = useStore()
    const showReactionSelector = ref(false)

    // Common emoji reactions
    const commonEmojis = ['👍', '❤️', '😂', '😮', '😢', '👏', '🎉', '🤔']

    // Group reactions by emoji and count them
    const groupedReactions = computed(() => {
      const grouped = {}

      props.reactions.forEach((reaction) => {
        if (!grouped[reaction.emoji]) {
          grouped[reaction.emoji] = {
            emoji: reaction.emoji,
            count: 0,
            users: [],
          }
        }

        grouped[reaction.emoji].count++
        grouped[reaction.emoji].users.push({
          id: reaction.userId,
          username: reaction.username || 'User',
          reactionId: reaction.id,
        })
      })

      // Convert to array and sort by count (descending)
      return Object.values(grouped).sort((a, b) => b.count - a.count)
    })

    // Check if current user has reacted with this emoji
    const isUserReaction = (reaction) => {
      const currentUserId = store.getters['auth/user']?.id
      return reaction.users.some((user) => user.id === currentUserId)
    }

    // Get the current user's reaction ID for a specific emoji
    const getUserReactionId = (emoji) => {
      const currentUserId = store.getters['auth/user']?.id
      const userReaction = props.reactions.find(
        (r) => r.userId === currentUserId && r.emoji === emoji,
      )
      return userReaction?.id
    }

    // Toggle reaction selector
    const toggleReactionSelector = () => {
      showReactionSelector.value = !showReactionSelector.value
    }

    // Add a new reaction
    const addReaction = (emoji) => {
      showReactionSelector.value = false
      emit('add-reaction', emoji)
    }

    // Toggle a reaction (add if not present, remove if already there)
    const toggleReaction = (emoji) => {
      const hasReacted = isUserReaction(groupedReactions.value.find((r) => r.emoji === emoji))

      if (hasReacted) {
        const reactionId = getUserReactionId(emoji)
        if (reactionId) {
          store.dispatch('messages/removeReaction', {
            messageId: props.messageId,
            reactionId,
          })
        }
      } else {
        store.dispatch('messages/addReaction', {
          messageId: props.messageId,
          reaction: emoji,
        })
      }
    }

    return {
      showReactionSelector,
      commonEmojis,
      groupedReactions,
      isUserReaction,
      toggleReactionSelector,
      addReaction,
      toggleReaction,
    }
  },
}
</script>

<style scoped>
.message-reactions {
  position: relative;
  margin-top: 0.25rem;
}

.reactions-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
}

.reaction {
  display: flex;
  align-items: center;
  padding: 0.125rem 0.375rem;
  border-radius: 1rem;
  background-color: #f3f4f6;
  border: 1px solid #e5e7eb;
  cursor: pointer;
  user-select: none;
  transition: all 0.2s;
}

.reaction:hover {
  background-color: #e5e7eb;
}

.reaction.user-reacted {
  background-color: #dbeafe;
  border-color: #bfdbfe;
}

.reaction-emoji {
  font-size: 0.875rem;
  margin-right: 0.25rem;
}

.reaction-count {
  font-size: 0.75rem;
  color: #6b7280;
}

.user-reacted .reaction-count {
  color: #3b82f6;
}

.reaction-selector {
  position: absolute;
  bottom: calc(100% + 0.25rem);
  left: 0;
  background-color: white;
  border-radius: 0.5rem;
  box-shadow:
    0 4px 6px -1px rgba(0, 0, 0, 0.1),
    0 2px 4px -1px rgba(0, 0, 0, 0.06);
  border: 1px solid #e5e7eb;
  z-index: 10;
  padding: 0.5rem;
}

.selector-content {
  display: flex;
  gap: 0.25rem;
}

.emoji-button {
  font-size: 1.125rem;
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: none;
  cursor: pointer;
  border-radius: 0.5rem;
  transition: background-color 0.2s;
}

.emoji-button:hover {
  background-color: #f3f4f6;
}
</style>
