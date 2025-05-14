// src/api/endpoints/conversations.js
import apiClient from '../client'

const conversationsApi = {
  getAll() {
    return apiClient.get('/conversations')
  },

  getById(id) {
    return apiClient.get(`/conversations/${id}`)
  },

  create(conversationData) {
    console.log('Creating conversation with data:', conversationData)
    return apiClient.post('/conversations', conversationData)
  },

  update(id, conversationData) {
    return apiClient.put(`/conversations/${id}`, conversationData)
  },

  delete(id) {
    return apiClient.delete(`/conversations/${id}`)
  },
}

export default conversationsApi
