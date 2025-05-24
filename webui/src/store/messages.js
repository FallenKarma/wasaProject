import { defineStore } from 'pinia'
import apiClient from '@/api/client'
import { messagesApi } from '@/api/endpoints/messages'

export const useMessageStore = defineStore('messages', {
  state: () => ({
    messages: [],
    isLoading: false,
    error: null,
    hasMoreMessages: true,
    currentPage: 1,
    lastMessageTimestamp: null,
  }),

  actions: {
    // Fetch messages for a conversation
    async fetchMessages(conversationId, reset = false) {
      if (reset) {
        this.resetPagination()
      }

      if (!this.hasMoreMessages && !reset) {
        return []
      }

      this.isLoading = true
      this.error = null

      try {
        const response = await apiClient.get(`/conversations/${conversationId}/messages`, {
          params: {
            page: this.currentPage,
            limit: 20,
            before: this.lastMessageTimestamp,
          },
        })

        const messages = response.data.messages || response.data

        if (reset) {
          this.messages = messages
        } else {
          // Add older messages to the end of the array
          this.messages = [...this.messages, ...messages]
        }

        // Update pagination state
        this.currentPage++
        this.hasMoreMessages = messages.length === 20

        if (messages.length > 0) {
          const lastMessage = messages[messages.length - 1]
          this.lastMessageTimestamp = lastMessage.created_at || lastMessage.timestamp
        }

        return messages
      } catch (error) {
        this.error = error.message || 'Failed to fetch messages'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Fetch more messages (used for pagination/infinite scroll)
    async fetchMoreMessages(conversationId) {
      return this.fetchMessages(conversationId, false)
    },

    // Send a new message
    async sendMessage(messageData) {
      this.isLoading = true
      this.error = null
      console.log('Sending message:', messageData)
      try {
        const response = await messagesApi.send(messageData)

        // Add new message to the beginning of the array
        this.messages.push(response.data)
        return response.data
      } catch (error) {
        this.error = error.message || 'Failed to send message'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Add a message locally (for optimistic updates)
    addMessage(message) {
      this.messages = [message, ...this.messages]
    },

    // Update a message (for editing)
    async updateMessage({ messageId, content }) {
      this.isLoading = true
      this.error = null

      try {
        const response = await apiClient.put(`/messages/${messageId}`, { content })

        const index = this.messages.findIndex((m) => m.id === messageId)
        if (index !== -1) {
          this.messages.splice(index, 1, response.data)
        }

        return response.data
      } catch (error) {
        this.error = error.message || 'Failed to update message'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Delete a message
    async deleteMessage(messageId) {
      this.isLoading = true
      this.error = null

      try {
        await apiClient.delete(`/messages/${messageId}`)
        this.messages = this.messages.filter((m) => m.id !== messageId)
        return true
      } catch (error) {
        this.error = error.message || 'Failed to delete message'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Add a reaction to a message
    async addReaction({ messageId, reaction }) {
      try {
        const response = await apiClient.post(`/messages/${messageId}/reactions`, { reaction })

        const index = this.messages.findIndex((m) => m.id === messageId)
        if (index !== -1) {
          this.messages.splice(index, 1, response.data)
        }

        return response.data
      } catch (error) {
        console.error('Failed to add reaction:', error)
        throw error
      }
    },

    // Remove a reaction from a message
    async removeReaction({ messageId, reactionId }) {
      try {
        const response = await apiClient.delete(`/messages/${messageId}/reactions/${reactionId}`)

        const index = this.messages.findIndex((m) => m.id === messageId)
        if (index !== -1) {
          this.messages.splice(index, 1, response.data)
        }

        return response.data
      } catch (error) {
        console.error('Failed to remove reaction:', error)
        throw error
      }
    },

    // Reset pagination
    resetPagination() {
      this.currentPage = 1
      this.hasMoreMessages = true
      this.lastMessageTimestamp = null
    },

    // Clear messages when leaving a conversation
    clearMessages() {
      this.messages = []
      this.resetPagination()
    },
  },
})
