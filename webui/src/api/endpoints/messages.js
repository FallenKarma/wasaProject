import apiClient from '../client'

export const messagesApi = {
  send(messageData) {
    return apiClient.post('/messages', messageData)
  },

  forward(messageId, targetConversationId) {
    return apiClient.post('/messages/forward', {
      messageId,
      targetConversationId,
    })
  },

  addComment(messageId, reactionData) {
    return apiClient.post(`/messages/${messageId}/comment`, reactionData)
  },

  removeComment(messageId) {
    return apiClient.delete(`/messages/${messageId}/comment`)
  },

  delete(messageId) {
    return apiClient.delete(`/messages/${messageId}`)
  },
}
