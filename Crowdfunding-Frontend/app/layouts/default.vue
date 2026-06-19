<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '~/stores/auth'

const isMobileMenuOpen = ref(false)
const currentLang = ref('id')

const authStore = useAuthStore()

onMounted(() => {
  authStore.init()
})

const toggleMobileMenu = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value
}

const switchLanguage = (lang) => {
  currentLang.value = lang
}

const logout = () => {
  authStore.logout()
  window.location.href = '/'
}
</script>
<template>
  <div class="font-sans bg-white text-gray-900 antialiased min-h-screen flex flex-col pt-16">
    <!-- Header / Navbar -->
    <header class="fixed top-0 w-full z-50 bg-white shadow-sm border-b">
      <nav class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <!-- Logo -->
          <div class="flex items-center">
            <NuxtLink to="/" class="flex items-center">
              <div class="w-8 h-8 bg-gradient-to-r from-green-500 to-blue-500 rounded-lg flex items-center justify-center">
                <i class="fas fa-heart text-white text-sm"></i>
              </div>
              <span class="ml-3 text-xl font-bold text-gray-900 hidden sm:block">Yayasan Peduli Amal Indonesia</span>
              <span class="ml-3 text-xl font-bold text-gray-900 sm:hidden">YPAI</span>
            </NuxtLink>
          </div>
          
          <!-- Desktop Navigation -->
          <div class="hidden lg:flex items-center space-x-8">
            <NuxtLink to="/" class="text-gray-600 hover:text-green-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">Home</NuxtLink>
            <NuxtLink to="/#about" class="text-gray-600 hover:text-green-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">About</NuxtLink>
            <NuxtLink to="/charities" class="text-gray-600 hover:text-green-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">Programs</NuxtLink>
            <NuxtLink to="/articles" class="text-gray-600 hover:text-green-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">Articles</NuxtLink>
            <NuxtLink to="/#contact" class="text-gray-600 hover:text-green-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">Contact</NuxtLink>
            
            <template v-if="authStore.user">
              <NuxtLink :to="authStore.user?.role === 'admin' ? '/admin' : '/dashboard'" class="text-xs text-gray-700 hover:text-green-600 font-semibold flex items-center space-x-1 border border-gray-200 px-2.5 py-1 rounded-lg hover:bg-slate-50 transition-colors">
                <i class="fas fa-user-circle text-sm text-green-600"></i>
                <span>Hai, {{ authStore.user.name }}</span>
              </NuxtLink>
              <button @click="logout" class="text-gray-600 hover:text-red-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">Logout</button>
            </template>
            <template v-else>
              <NuxtLink to="/login" class="text-gray-600 hover:text-green-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">Login</NuxtLink>
            </template>
            
            <!-- Language Switch Dropdown -->
            <div class="relative group">
              <button class="flex items-center space-x-2 text-gray-600 hover:text-green-600 px-3 py-2 rounded-md text-sm font-medium transition-colors">
                <img :src="currentLang === 'id' ? 'https://flagcdn.com/w20/id.png' : 'https://flagcdn.com/w20/us.png'" alt="Language" class="w-5 h-4 rounded-sm">
                <i class="fas fa-chevron-down text-xs"></i>
              </button>
              <div class="absolute right-0 mt-2 w-36 bg-white rounded-lg shadow-lg border border-gray-100 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-300 z-50">
                <div class="py-2">
                  <button @click="switchLanguage('en')" class="flex items-center w-full px-3 py-2 text-sm text-gray-700 hover:bg-green-50 hover:text-green-600 transition-colors">
                    <img src="https://flagcdn.com/w20/us.png" alt="English" class="w-5 h-4 rounded-sm mr-3">
                    English
                  </button>
                  <button @click="switchLanguage('id')" class="flex items-center w-full px-3 py-2 text-sm text-gray-700 hover:bg-green-50 hover:text-green-600 transition-colors">
                    <img src="https://flagcdn.com/w20/id.png" alt="Indonesia" class="w-5 h-4 rounded-sm mr-3">
                    Indonesia
                  </button>
                </div>
              </div>
            </div>
          </div>
          
          <!-- Mobile Menu Button -->
          <div class="lg:hidden">
            <button @click="toggleMobileMenu" class="text-gray-600 hover:text-green-600 p-2 rounded-md transition-colors">
              <i class="fas fa-bars text-xl"></i>
            </button>
          </div>
        </div>
      </nav>
    </header>

    <!-- Mobile Sidebar -->
    <div :class="{'translate-x-0': isMobileMenuOpen, '-translate-x-full': !isMobileMenuOpen}" class="fixed inset-y-0 left-0 z-50 w-80 bg-white shadow-xl transform transition-all duration-300 ease-in-out lg:hidden border-r border-gray-100">
      <!-- Sidebar Header -->
      <div class="flex items-center justify-between h-16 px-6 border-b border-gray-100">
        <div class="flex items-center space-x-3">
          <div class="w-8 h-8 bg-green-500 rounded-lg flex items-center justify-center">
            <i class="fas fa-heart text-white text-sm"></i>
          </div>
          <span class="text-lg font-semibold text-gray-900">YPAI</span>
        </div>
        <button @click="toggleMobileMenu" class="p-2 rounded-lg hover:bg-gray-100 transition-colors">
          <i class="fas fa-times text-gray-500 hover:text-gray-700"></i>
        </button>
      </div>

      <!-- Sidebar Navigation -->
      <div class="p-6">
        <nav class="space-y-1">
          <NuxtLink to="/" @click="toggleMobileMenu" class="flex items-center px-3 py-2.5 text-sm font-medium text-gray-700 rounded-lg hover:bg-gray-50 hover:text-green-600 transition-colors">
            <i class="fas fa-home w-5 h-5 mr-3 text-gray-400"></i>
            Home
          </NuxtLink>
          <NuxtLink to="/#about" @click="toggleMobileMenu" class="flex items-center px-3 py-2.5 text-sm font-medium text-gray-700 rounded-lg hover:bg-gray-50 hover:text-green-600 transition-colors">
            <i class="fas fa-info-circle w-5 h-5 mr-3 text-gray-400"></i>
            About
          </NuxtLink>
          <NuxtLink to="/charities" @click="toggleMobileMenu" class="flex items-center px-3 py-2.5 text-sm font-medium text-gray-700 rounded-lg hover:bg-gray-50 hover:text-green-600 transition-colors">
            <i class="fas fa-heart w-5 h-5 mr-3 text-gray-400"></i>
            Programs
          </NuxtLink>
          <NuxtLink to="/articles" @click="toggleMobileMenu" class="flex items-center px-3 py-2.5 text-sm font-medium text-gray-700 rounded-lg hover:bg-gray-50 hover:text-green-600 transition-colors">
            <i class="fas fa-newspaper w-5 h-5 mr-3 text-gray-400"></i>
            Articles
          </NuxtLink>
          <NuxtLink to="/#contact" @click="toggleMobileMenu" class="flex items-center px-3 py-2.5 text-sm font-medium text-gray-700 rounded-lg hover:bg-gray-50 hover:text-green-600 transition-colors">
            <i class="fas fa-envelope w-5 h-5 mr-3 text-gray-400"></i>
            Contact
          </NuxtLink>
        </nav>
        
        <div class="mt-6 pt-6 border-t border-gray-200">
          <template v-if="authStore.user">
            <NuxtLink :to="authStore.user?.role === 'admin' ? '/admin' : '/dashboard'" @click="toggleMobileMenu" class="flex items-center px-3 py-2 text-xs font-semibold text-gray-700 hover:text-green-600 mb-2">
              <i class="fas fa-user-circle mr-2 text-green-600 text-sm"></i>
              <span>Hai, {{ authStore.user.name }}</span>
            </NuxtLink>
            <button @click="logout" class="flex items-center w-full px-3 py-2.5 text-sm font-medium text-red-600 rounded-lg hover:bg-red-50 transition-colors">
              <i class="fas fa-sign-out-alt w-5 h-5 mr-3 text-red-600"></i>
              Logout
            </button>
          </template>
          <template v-else>
            <NuxtLink to="/login" @click="toggleMobileMenu" class="flex items-center px-3 py-2.5 text-sm font-medium text-gray-700 rounded-lg hover:bg-gray-50 hover:text-green-600 transition-colors">
              <i class="fas fa-sign-in-alt w-5 h-5 mr-3 text-gray-400"></i>
              Login
            </NuxtLink>
          </template>
        </div>
        
        <div class="mt-6 pt-6 border-t border-gray-200">
          <div class="flex items-center justify-between mb-3">
            <span class="text-xs font-medium text-gray-500 uppercase tracking-wide">Language</span>
          </div>
          <div class="grid grid-cols-2 gap-2">
            <button @click="switchLanguage('id')" :class="{'bg-green-50 text-green-600': currentLang === 'id'}" class="flex items-center justify-center px-3 py-2 text-xs font-medium text-gray-700 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors">
              <img src="https://flagcdn.com/w20/id.png" alt="Indonesia" class="w-4 h-3 mr-2">
              ID
            </button>
            <button @click="switchLanguage('en')" :class="{'bg-green-50 text-green-600': currentLang === 'en'}" class="flex items-center justify-center px-3 py-2 text-xs font-medium text-gray-700 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors">
              <img src="https://flagcdn.com/w20/us.png" alt="English" class="w-4 h-3 mr-2">
              EN
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Mobile Sidebar Overlay -->
    <div v-if="isMobileMenuOpen" @click="toggleMobileMenu" class="fixed inset-0 bg-black bg-opacity-50 z-40 lg:hidden"></div>

    <!-- Main Content Area -->
    <main class="flex-grow">
      <slot />
    </main>

    <!-- Footer -->
    <footer class="bg-secondary text-white py-8 md:py-12">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div>
            <div class="flex items-center mb-4">
              <div class="w-8 h-8 bg-green-500 rounded-lg flex items-center justify-center">
                <i class="fas fa-heart text-white text-sm"></i>
              </div>
              <span class="ml-3 text-xl font-bold">YPAI</span>
            </div>
            <p class="text-gray-300 text-sm leading-relaxed">
              Yayasan Peduli Amal Indonesia berkomitmen untuk menghubungkan kebaikan dengan mereka yang membutuhkan secara transparan dan terpercaya.
            </p>
          </div>
          <div>
            <h3 class="text-base md:text-lg font-semibold mb-4">Contact Us</h3>
            <ul class="space-y-2 text-sm text-gray-300">
              <li class="flex items-start"><i class="fas fa-map-marker-alt mt-1 mr-3 text-green-500"></i>123 Hope Street, Charity City, CC 12345</li>
              <li class="flex items-center"><i class="fas fa-phone mr-3 text-green-500"></i>+1 (555) 123-4567</li>
              <li class="flex items-center"><i class="fas fa-envelope mr-3 text-green-500"></i>info@hopefoundation.org</li>
            </ul>
          </div>
          <div>
            <h3 class="text-base md:text-lg font-semibold mb-4">Quick Links</h3>
            <ul class="grid grid-cols-2 gap-2 text-sm text-gray-300">
              <li><NuxtLink to="/" class="hover:text-green-500 transition-colors">Home</NuxtLink></li>
              <li><NuxtLink to="/#about" class="hover:text-green-500 transition-colors">About</NuxtLink></li>
              <li><NuxtLink to="/charities" class="hover:text-green-500 transition-colors">Programs</NuxtLink></li>
              <li><NuxtLink to="/articles" class="hover:text-green-500 transition-colors">Articles</NuxtLink></li>
              <li><NuxtLink to="/login" class="hover:text-green-500 transition-colors">Login</NuxtLink></li>
            </ul>
          </div>
        </div>
        <div class="mt-8 pt-8 border-t border-gray-800 text-center">
          <p class="text-gray-400 text-xs md:text-sm">&copy; 2026 Yayasan Peduli Amal Indonesia. All rights reserved.</p>
        </div>
      </div>
    </footer>
  </div>
</template>
