<!-- src/components/auth/LoginForm.vue -->
<template>
  <div class="login-container">
    <div class="login-card">
      <h2 class="login-title">Welcome Back</h2>

      <div v-if="errorMessage" class="error-message">
        {{ errorMessage }}
      </div>

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label for="username">Username</label>
          <input
            type="text"
            id="username"
            v-model="username"
            required
            class="form-input"
            placeholder="Enter your username"
          />
        </div>

        <div class="form-group">
          <label for="password">Password</label>
          <input
            type="password"
            id="password"
            v-model="password"
            required
            class="form-input"
            placeholder="Enter your password"
          />
          <div class="forgot-password">
            <a href="#" @click.prevent="forgotPassword">Forgot password?</a>
          </div>
        </div>

        <button type="submit" class="login-button" :disabled="isLoading">
          <span v-if="isLoading">Logging in...</span>
          <span v-else>Login</span>
        </button>

        <div class="register-link">
          Don't have an account?
          <router-link to="/register">Sign up</router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import authApi from '@/api/auth'
import { useAuthStore } from '@/store/auth'

export default {
  name: 'LoginForm',
  setup() {
    const username = ref('')
    const password = ref('')
    const isLoading = ref(false)
    const errorMessage = ref('')
    const router = useRouter()
    const authStore = useAuthStore()

    const handleLogin = async () => {
      try {
        isLoading.value = true
        errorMessage.value = ''

        const credentials = {
          username: username.value,
          password: password.value,
        }

        const response = await authApi.login(credentials)

        // Store authentication tokens
        authStore.setAuthToken(response.token)
        authStore.setUser(response.user)

        // Redirect to home page
        router.push('/')
      } catch (error) {
        console.error('Login failed:', error)
        errorMessage.value =
          error.response?.data?.message || 'Failed to login. Please check your credentials.'
      } finally {
        isLoading.value = false
      }
    }

    const forgotPassword = () => {
      // Handle forgot password functionality
      alert('Password reset functionality coming soon!')
    }

    return {
      username,
      password,
      isLoading,
      errorMessage,
      handleLogin,
      forgotPassword,
    }
  },
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f7fb;
}

.login-card {
  width: 100%;
  max-width: 400px;
  padding: 2rem;
  background-color: white;
  border-radius: 0.5rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.login-title {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 1.5rem;
  text-align: center;
  color: #333;
}

.form-group {
  margin-bottom: 1.25rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #4b5563;
}

.form-input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 0.375rem;
  font-size: 1rem;
  transition: border-color 0.15s ease-in-out;
}

.form-input:focus {
  border-color: #3b82f6;
  outline: none;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.login-button {
  width: 100%;
  padding: 0.75rem;
  background-color: #3b82f6;
  color: white;
  border: none;
  border-radius: 0.375rem;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.15s ease-in-out;
}

.login-button:hover {
  background-color: #2563eb;
}

.login-button:disabled {
  background-color: #93c5fd;
  cursor: not-allowed;
}

.error-message {
  padding: 0.75rem;
  margin-bottom: 1rem;
  background-color: #fee2e2;
  border: 1px solid #fecaca;
  border-radius: 0.375rem;
  color: #ef4444;
  font-size: 0.875rem;
}

.forgot-password {
  text-align: right;
  margin-top: 0.5rem;
}

.forgot-password a {
  color: #3b82f6;
  font-size: 0.875rem;
  text-decoration: none;
}

.forgot-password a:hover {
  text-decoration: underline;
}

.register-link {
  margin-top: 1.5rem;
  text-align: center;
  font-size: 0.875rem;
  color: #4b5563;
}

.register-link a {
  color: #3b82f6;
  font-weight: 500;
  text-decoration: none;
}

.register-link a:hover {
  text-decoration: underline;
}
</style>
