// src/api/endpoints/conversations.js
import apiClient from '../client'

export const conversationsApi = {
  getAll() {
    return apiClient.get('/conversations')
  },

  getById(id) {
    return apiClient.get(`/conversations/${id}`)
  },
}
