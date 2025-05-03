// src/api/interceptors.js
export function setupInterceptors(apiClient) {
  // Request interceptor
  apiClient.interceptors.request.use(
    (config) => {
      // Get token from localStorage
      const token = localStorage.getItem('auth_token')

      // If token exists, add to headers
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }

      return config
    },
    (error) => Promise.reject(error),
  )

  // Response interceptor
  apiClient.interceptors.response.use(
    (response) => response,
    (error) => {
      // Handle 401 Unauthorized errors
      if (error.response && error.response.status === 401) {
        // Clear token and redirect to login
        localStorage.removeItem('auth_token')
        window.location.href = '/login'
      }
      return Promise.reject(error)
    },
  )
}
