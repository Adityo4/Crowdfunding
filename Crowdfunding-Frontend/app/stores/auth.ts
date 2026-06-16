import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(null)
  const user = ref(null)

  // Initialize store state from localStorage
  const init = () => {
    if (typeof window !== 'undefined') {
      const storedToken = localStorage.getItem('auth_token')
      const storedUser = localStorage.getItem('user_data')
      
      if (storedToken) token.value = storedToken
      if (storedUser) {
        try {
          user.value = JSON.parse(storedUser)
        } catch (e) {
          user.value = null
        }
      }
    }
  }

  const login = (authToken, userData) => {
    if (typeof window !== 'undefined') {
      localStorage.setItem('auth_token', authToken)
      localStorage.setItem('user_data', JSON.stringify(userData))
    }
    token.value = authToken
    user.value = userData
  }

  const logout = () => {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('auth_token')
      localStorage.removeItem('user_data')
    }
    token.value = null
    user.value = null
  }

  const isAuthenticated = () => {
    return !!token.value
  }

  return {
    token,
    user,
    init,
    login,
    logout,
    isAuthenticated
  }
})
