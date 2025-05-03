import apiClient from '../client'

export const usersApi = {
  updateUsername(name) {
    return apiClient.put('/users/me/username', { name })
  },

  uploadPhoto(photoFile) {
    const formData = new FormData()
    formData.append('photo', photoFile)

    return apiClient.put('/users/me/photo', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
  },
}
