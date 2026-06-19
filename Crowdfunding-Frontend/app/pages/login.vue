<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

import { useAuthStore } from '~/stores/auth'

const router = useRouter()
const activeTab = ref('login') // 'login', 'register', 'forgot'
const authStore = useAuthStore()

// Form inputs
const email = ref('')
const password = ref('')
const fullName = ref('')
const phoneNumber = ref('')
const confirmPassword = ref('')
const agreeTerms = ref(false)

// UI Feedback
const isLoading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const showPassword = ref(false)
const showRegisterPassword = ref(false)
const showConfirmPassword = ref(false)

const handleLogin = async () => {
  isLoading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  
  try {
    const config = useRuntimeConfig()
    const response = await $fetch(`${config.public.apiBase}/auth/login`, {
      method: 'POST',
      body: {
        email: email.value,
        password: password.value
      }
    })
    
    if (response?.data?.token) {
      const userData = response.data.user || {}
      const parsedUser = {
        id: userData.id,
        name: userData.fullName || userData.email || 'User',
        email: userData.email,
        role: userData.role,
        phoneNumber: userData.phoneNumber || ''
      }
      
      // Save using Pinia store
      authStore.login(response.data.token, parsedUser)
      
      successMessage.value = 'Login sukses! Mengalihkan...'
      setTimeout(() => {
        router.push('/')
      }, 1000)
    }
  } catch (error) {
    errorMessage.value = error.data?.error?.message || 'Login gagal, periksa email dan password Anda.'
  } finally {
    isLoading.value = false
  }
}

const handleRegister = async () => {
  if (password.value !== confirmPassword.value) {
    errorMessage.value = 'Konfirmasi password tidak cocok'
    return
  }
  
  isLoading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  
  try {
    const config = useRuntimeConfig()
    const response = await $fetch(`${config.public.apiBase}/auth/register`, {
      method: 'POST',
      body: {
        fullName: fullName.value,
        email: email.value,
        password: password.value,
        phoneNumber: phoneNumber.value
      }
    })
    
    if (response?.data) {
      successMessage.value = 'Registrasi berhasil! Silakan login.'
      activeTab.value = 'login'
      password.value = ''
    }
  } catch (error) {
    errorMessage.value = error.data?.error?.message || 'Registrasi gagal, coba lagi.'
  } finally {
    isLoading.value = false
  }
}

const handleReset = () => {
  successMessage.value = 'Link reset password berhasil dikirim ke email Anda!'
}
</script>

<template>
  <div class="min-h-[80vh] flex items-center justify-center p-4 pt-10 bg-slate-50">
    <div class="w-full max-w-md">
      <!-- Logo and Header -->
      <div class="text-center mb-6">
        <div class="inline-flex items-center justify-center w-12 h-12 bg-green-600 rounded-xl mb-4 shadow-lg">
          <i class="fas fa-heart text-xl text-white"></i>
        </div>
        <h1 class="text-xl font-bold text-slate-900 mb-1">Yayasan Peduli Amal Indonesia</h1>
        <p class="text-gray-600 text-sm">Making Hope Happen Together</p>
      </div>

      <!-- Form Container -->
      <div class="bg-white rounded-2xl shadow-xl p-8 border border-gray-100">
        <!-- Form Tabs -->
        <div class="flex mb-6 border-b border-gray-200">
          <button @click="activeTab = 'login'" :class="{'text-green-600 border-b-2 border-green-600': activeTab === 'login', 'text-gray-500 hover:text-green-600': activeTab !== 'login'}" class="flex-1 py-3 px-4 text-center font-semibold transition-colors">
            <i class="fas fa-sign-in-alt mr-2"></i>Login
          </button>
          <button @click="activeTab = 'register'" :class="{'text-green-600 border-b-2 border-green-600': activeTab === 'register', 'text-gray-500 hover:text-green-600': activeTab !== 'register'}" class="flex-1 py-3 px-4 text-center font-semibold transition-colors">
            <i class="fas fa-user-plus mr-2"></i>Register
          </button>
          <button @click="activeTab = 'forgot'" :class="{'text-green-600 border-b-2 border-green-600': activeTab === 'forgot', 'text-gray-500 hover:text-green-600': activeTab !== 'forgot'}" class="flex-1 py-3 px-4 text-center font-semibold transition-colors">
            <i class="fas fa-key mr-2"></i>Reset
          </button>
        </div>

        <!-- Feedback Messages -->
        <div v-if="errorMessage" class="mb-4 p-3 bg-red-50 text-red-600 rounded-lg text-sm flex items-center">
          <i class="fas fa-exclamation-circle mr-2"></i>
          {{ errorMessage }}
        </div>
        <div v-if="successMessage" class="mb-4 p-3 bg-green-50 text-green-600 rounded-lg text-sm flex items-center">
          <i class="fas fa-check-circle mr-2"></i>
          {{ successMessage }}
        </div>
        
        <!-- Login Form -->
        <div v-if="activeTab === 'login'">
          <div class="text-center mb-6">
            <h2 class="text-xl font-bold text-slate-900 mb-1">Welcome Back</h2>
            <p class="text-gray-500 text-sm">Sign in to your account to continue making a difference</p>
          </div>
          
          <form @submit.prevent="handleLogin" class="space-y-4">
            <div class="space-y-1">
              <label class="block text-sm font-medium text-gray-700">Email Address</label>
              <div class="relative group">
                <input v-model="email" type="email" required class="w-full px-4 py-2.5 pl-12 border border-gray-200 rounded-lg focus:border-green-600 focus:ring-2 focus:ring-green-600/10 focus:outline-none transition-all duration-300" placeholder="Enter your email">
                <i class="fas fa-envelope absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400 group-hover:text-green-600 transition-colors"></i>
              </div>
            </div>
            
            <div class="space-y-1">
              <label class="block text-sm font-medium text-gray-700">Password</label>
              <div class="relative group">
                <input v-model="password" :type="showPassword ? 'text' : 'password'" required class="w-full px-4 py-2.5 pl-12 pr-12 border border-gray-200 rounded-lg focus:border-green-600 focus:ring-2 focus:ring-green-600/10 focus:outline-none transition-all duration-300" placeholder="Enter your password">
                <i class="fas fa-lock absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400 group-hover:text-green-600 transition-colors"></i>
                <button type="button" @click="showPassword = !showPassword" class="absolute right-4 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-green-600">
                  <i class="fas" :class="showPassword ? 'fa-eye-slash' : 'fa-eye'"></i>
                </button>
              </div>
            </div>
            
            <div class="flex items-center justify-between pt-1">
              <label class="flex items-center space-x-2 cursor-pointer">
                <input type="checkbox" class="w-4 h-4 text-green-600 border-gray-300 rounded focus:ring-green-500">
                <span class="text-sm text-gray-600">Remember me</span>
              </label>
              <button type="button" @click="activeTab = 'forgot'" class="text-sm text-green-600 hover:text-green-700">Forgot password?</button>
            </div>
            
            <button type="submit" :disabled="isLoading" class="w-full bg-green-600 text-white py-3 rounded-lg font-semibold hover:bg-green-700 transition-all duration-300 shadow-lg hover:shadow-xl disabled:opacity-50">
              <span v-if="isLoading"><i class="fas fa-spinner fa-spin mr-2"></i>Loading...</span>
              <span v-else>Sign In</span>
            </button>
          </form>
        </div>
        
        <!-- Register Form -->
        <div v-if="activeTab === 'register'">
          <div class="text-center mb-6">
            <h2 class="text-xl font-bold text-slate-900 mb-1">Create Account</h2>
            <p class="text-gray-500 text-sm">Join us in making a difference in the world</p>
          </div>
          
          <form @submit.prevent="handleRegister" class="space-y-4">
            <div class="space-y-1">
              <label class="block text-sm font-medium text-gray-700">Full Name</label>
              <div class="relative group">
                <input v-model="fullName" type="text" required class="w-full px-4 py-2.5 pl-12 border border-gray-200 rounded-lg focus:border-green-600 focus:ring-2 focus:ring-green-600/10 focus:outline-none" placeholder="Your full name">
                <i class="fas fa-user absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400 group-hover:text-green-600"></i>
              </div>
            </div>
            
            <div class="space-y-1">
              <label class="block text-sm font-medium text-gray-700">Email Address</label>
              <div class="relative group">
                <input v-model="email" type="email" required class="w-full px-4 py-2.5 pl-12 border border-gray-200 rounded-lg focus:border-green-600 focus:ring-2 focus:ring-green-600/10 focus:outline-none" placeholder="Your email address">
                <i class="fas fa-envelope absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400 group-hover:text-green-600"></i>
              </div>
            </div>

            <div class="space-y-1">
              <label class="block text-sm font-medium text-gray-700">Phone Number</label>
              <div class="relative group">
                <input v-model="phoneNumber" type="tel" class="w-full px-4 py-2.5 pl-12 border border-gray-200 rounded-lg focus:border-green-600 focus:ring-2 focus:ring-green-600/10 focus:outline-none" placeholder="Your phone number">
                <i class="fas fa-phone absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400 group-hover:text-green-600"></i>
              </div>
            </div>
            
            <div class="space-y-1">
              <label class="block text-sm font-medium text-gray-700">Password</label>
              <div class="relative group">
                <input v-model="password" :type="showRegisterPassword ? 'text' : 'password'" required class="w-full px-4 py-2.5 pl-12 pr-12 border border-gray-200 rounded-lg focus:border-green-600 focus:ring-2 focus:ring-green-600/10 focus:outline-none" placeholder="Create a password">
                <i class="fas fa-lock absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400 group-hover:text-green-600"></i>
                <button type="button" @click="showRegisterPassword = !showRegisterPassword" class="absolute right-4 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-green-600">
                  <i class="fas" :class="showRegisterPassword ? 'fa-eye-slash' : 'fa-eye'"></i>
                </button>
              </div>
            </div>

            <div class="space-y-1">
              <label class="block text-sm font-medium text-gray-700">Confirm Password</label>
              <div class="relative group">
                <input v-model="confirmPassword" :type="showConfirmPassword ? 'text' : 'password'" required class="w-full px-4 py-2.5 pl-12 pr-12 border border-gray-200 rounded-lg focus:border-green-600 focus:ring-2 focus:ring-green-600/10 focus:outline-none" placeholder="Confirm your password">
                <i class="fas fa-lock absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400 group-hover:text-green-600"></i>
                <button type="button" @click="showConfirmPassword = !showConfirmPassword" class="absolute right-4 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-green-600">
                  <i class="fas" :class="showConfirmPassword ? 'fa-eye-slash' : 'fa-eye'"></i>
                </button>
              </div>
            </div>
            
            <div class="flex items-start space-x-2 pt-1">
              <input v-model="agreeTerms" type="checkbox" id="terms" required class="w-4 h-4 text-green-600 border-gray-300 rounded focus:ring-green-500 mt-1">
              <label for="terms" class="text-sm text-gray-600">
                I agree to the <a href="#" class="text-green-600 hover:text-green-700">Terms of Service</a> and <a href="#" class="text-green-600 hover:text-green-700">Privacy Policy</a>
              </label>
            </div>
            
            <button type="submit" :disabled="isLoading" class="w-full bg-green-600 text-white py-3 rounded-lg font-semibold hover:bg-green-700 transition-all duration-300 shadow-lg hover:shadow-xl disabled:opacity-50">
              <span v-if="isLoading"><i class="fas fa-spinner fa-spin mr-2"></i>Loading...</span>
              <span v-else>Create Account</span>
            </button>
          </form>
        </div>

        <!-- Forgot Form -->
        <div v-if="activeTab === 'forgot'">
          <div class="text-center mb-6">
            <h2 class="text-xl font-bold text-slate-900 mb-1">Reset Password</h2>
            <p class="text-gray-500 text-sm">Enter your email address and we'll send you a link to reset your password</p>
          </div>
          
          <form @submit.prevent="handleReset" class="space-y-4">
            <div class="space-y-1">
              <label class="block text-sm font-medium text-gray-700">Email Address</label>
              <div class="relative group">
                <input v-model="email" type="email" required class="w-full px-4 py-2.5 pl-12 border border-gray-200 rounded-lg focus:border-green-600 focus:ring-2 focus:ring-green-600/10 focus:outline-none" placeholder="Enter your email">
                <i class="fas fa-envelope absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400 group-hover:text-green-600"></i>
              </div>
            </div>
            
            <button type="submit" class="w-full bg-green-600 text-white py-3 rounded-lg font-semibold hover:bg-green-700 transition-all duration-300 shadow-lg hover:shadow-xl">
              Send Reset Link
            </button>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>
