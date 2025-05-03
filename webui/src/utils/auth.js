// src/utils/auth.js
export const TOKEN_KEY = 'wasa_auth_token'
export const USER_ID_KEY = 'wasa_user_id'

export const auth = {
  setToken(token) {
    localStorage.setItem(TOKEN_KEY, token)
  },

  getToken() {
    return localStorage.getItem(TOKEN_KEY)
  },

  removeToken() {
    localStorage.removeItem(TOKEN_KEY)
  },

  setUserId(id) {
    localStorage.setItem(USER_ID_KEY, id)
  },

  getUserId() {
    return localStorage.getItem(USER_ID_KEY)
  },

  removeUserId() {
    localStorage.removeItem(USER_ID_KEY)
  },

  isAuthenticated() {
    return !!this.getToken()
  },

  logout() {
    this.removeToken()
    this.removeUserId()
  },
}
