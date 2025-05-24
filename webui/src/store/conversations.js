import { defineStore } from 'pinia'
import conversationsApi from '@/api/endpoints/conversations'

export const useConversationStore = defineStore('conversations', {
  state: () => ({
    conversations: [],
    currentConversation: null,
    isLoading: false,
    error: null,
  }),

  getters: {
    allConversations: (state) => state.conversations,
  },

  actions: {
    // Fetch all conversations for the current user
    async fetchConversations() {
      this.isLoading = true
      this.error = null

      try {
        const response = await conversationsApi.getAll()
        this.conversations = response.data
        return response.data
      } catch (error) {
        this.error = error.message || 'Failed to fetch conversations'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Fetch a single conversation by ID
    async fetchConversation(conversationId) {
      this.isLoading = true
      this.error = null

      try {
        const response = await conversationsApi.getById(conversationId)
        this.currentConversation = response.data
        return response.data
      } catch (error) {
        this.error = error.message || 'Failed to fetch conversation'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Create a new conversation
    async createConversation(conversationData) {
      this.isLoading = true
      this.error = null

      try {
        const response = await conversationsApi.create(conversationData)
        this.conversations.unshift(response.data)
        return response.data
      } catch (error) {
        this.error = error.message || 'Failed to create conversation'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Update an existing conversation
    async updateConversation({ conversationId, data }) {
      this.isLoading = true
      this.error = null

      try {
        const response = await conversationsApi.update(conversationId, data)
        const updatedConversation = response.data

        // Update in the array
        const index = this.conversations.findIndex((c) => c.id === updatedConversation.id)
        if (index !== -1) {
          this.conversations.splice(index, 1, updatedConversation)
        }

        // Also update current conversation if it's the same one
        if (this.currentConversation && this.currentConversation.id === updatedConversation.id) {
          this.currentConversation = updatedConversation
        }

        return response.data
      } catch (error) {
        this.error = error.message || 'Failed to update conversation'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Delete a conversation
    async deleteConversation(conversationId) {
      this.isLoading = true
      this.error = null

      try {
        await conversationsApi.delete(conversationId)
        this.conversations = this.conversations.filter((c) => c.id !== conversationId)

        // Clear current conversation if it was the one removed
        if (this.currentConversation && this.currentConversation.id === conversationId) {
          this.currentConversation = null
        }

        return true
      } catch (error) {
        this.error = error.message || 'Failed to delete conversation'
        throw error
      } finally {
        this.isLoading = false
      }
    },

    async addMessageToConversation(conversationId, message) {
      const conversation = this.conversations.find((c) => c.id === conversationId)
      if (conversation) {
        conversation.messages = conversation.messages || []
        conversation.messages.unshift(message)
      }
    },
  },
})
