import { useRuntimeConfig } from '#app'

export const useApiFetch = (path, opts = {}) => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase || 'http://localhost:8080/v1'

  // Get token if it exists in localStorage
  let token = null
  if (typeof window !== 'undefined') {
    token = localStorage.getItem('auth_token')
  }

  return useFetch(path, {
    baseURL: apiBase,
    ...opts,
    async onRequest({ options }) {
      // Fetch token dynamically inside request hook to make it reactive on client side
      if (typeof window !== 'undefined') {
        const token = localStorage.getItem('auth_token')
        if (token) {
          options.headers = {
            ...options.headers,
            'Authorization': `Bearer ${token}`
          }
        }
      }
    },
    async onResponseError({ response }) {
      if (response.status === 401 && typeof window !== 'undefined') {
        localStorage.removeItem('auth_token')
        localStorage.removeItem('user_data')
        window.location.href = '/login'
      }
    }
  })
}
