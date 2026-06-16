import { useRuntimeConfig } from '#app'

export const useApiFetch = (path, opts = {}) => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase || 'http://localhost:8080/v1'

  // Get token if it exists in localStorage
  let token = null
  if (typeof window !== 'undefined') {
    token = localStorage.getItem('auth_token')
  }

  const defaultHeaders = {}
  if (token) {
    defaultHeaders['Authorization'] = `Bearer ${token}`
  }

  return useFetch(path, {
    baseURL: apiBase,
    ...opts,
    headers: {
      ...defaultHeaders,
      ...opts.headers,
    },
    async onResponseError({ response }) {
      if (response.status === 401 && typeof window !== 'undefined') {
        // Handle token expiration or unauthorized requests
        localStorage.removeItem('auth_token')
        localStorage.removeItem('user_data')
        window.location.href = '/login'
      }
    }
  })
}
