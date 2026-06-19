<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '~/stores/auth'

const isSidebarOpen = ref(true)
const isMobileSidebarOpen = ref(false)
const isComponentsMenuOpen = ref(false)
const currentLang = ref('id')

const authStore = useAuthStore()
const router = useRouter()

onMounted(() => {
  authStore.init()
  // Guard admin routes
  if (authStore.user?.role !== 'admin') {
    alert('Akses ditolak. Halaman ini hanya untuk Administrator.')
    router.push('/login')
  }
})

const toggleSidebar = () => {
  isSidebarOpen.value = !isSidebarOpen.value
}

const toggleMobileSidebar = () => {
  isMobileSidebarOpen.value = !isMobileSidebarOpen.value
}

const toggleComponentsMenu = () => {
  isComponentsMenuOpen.value = !isComponentsMenuOpen.value
}

const switchLanguage = (lang) => {
  currentLang.value = lang
}

const logout = () => {
  authStore.logout()
  window.location.href = '/login'
}
</script>

<template>
  <div class="font-sans bg-gray-50 text-gray-900 antialiased min-h-screen flex">
    <!-- Desktop Sidebar -->
    <aside 
      :class="[
        'fixed inset-y-0 left-0 z-40 w-64 bg-white shadow-lg border-r border-gray-200 transition-transform duration-300 lg:translate-x-0 lg:static lg:flex lg:flex-col',
        isSidebarOpen ? 'translate-x-0' : '-translate-x-full lg:hidden'
      ]"
    >
      <!-- Sidebar Header -->
      <div class="flex items-center justify-between h-16 px-6 border-b border-gray-200">
        <div class="flex items-center space-x-3">
          <div class="w-8 h-8 bg-gradient-to-r from-green-500 to-blue-500 rounded-lg flex items-center justify-center">
            <i class="fas fa-crown text-white text-sm"></i>
          </div>
          <span class="text-lg font-bold text-gray-800">Admin Panel</span>
        </div>
        <button @click="toggleSidebar" class="hidden lg:block p-2 rounded-lg hover:bg-gray-100 text-gray-500">
          <i class="fas fa-bars"></i>
        </button>
      </div>

      <!-- Sidebar Navigation -->
      <nav class="flex-grow px-4 py-6 space-y-1 overflow-y-auto">
        <NuxtLink to="/" class="flex items-center space-x-3 px-4 py-3 rounded-lg font-medium text-gray-600 hover:text-green-600 hover:bg-green-50 transition-colors">
          <i class="fas fa-home w-5"></i>
          <span>Kembali ke Home</span>
        </NuxtLink>
        <NuxtLink to="/admin" class="flex items-center space-x-3 px-4 py-3 rounded-lg font-medium text-gray-600 hover:text-green-600 hover:bg-green-50 transition-colors" active-class="text-green-600 bg-green-50">
          <i class="fas fa-tachometer-alt w-5"></i>
          <span>Dashboard</span>
        </NuxtLink>
        <NuxtLink to="/admin/charities" class="flex items-center space-x-3 px-4 py-3 rounded-lg font-medium text-gray-600 hover:text-green-600 hover:bg-green-50 transition-colors" active-class="text-green-600 bg-green-50">
          <i class="fas fa-heart w-5"></i>
          <span>Charities</span>
        </NuxtLink>
        <NuxtLink to="/admin/donations" class="flex items-center space-x-3 px-4 py-3 rounded-lg font-medium text-gray-600 hover:text-green-600 hover:bg-green-50 transition-colors" active-class="text-green-600 bg-green-50">
          <i class="fas fa-donate w-5"></i>
          <span>Donations</span>
        </NuxtLink>
        <NuxtLink to="/admin/users" class="flex items-center space-x-3 px-4 py-3 rounded-lg font-medium text-gray-600 hover:text-green-600 hover:bg-green-50 transition-colors" active-class="text-green-600 bg-green-50">
          <i class="fas fa-users w-5"></i>
          <span>Users</span>
        </NuxtLink>
        <NuxtLink to="/admin/articles" class="flex items-center space-x-3 px-4 py-3 rounded-lg font-medium text-gray-600 hover:text-green-600 hover:bg-green-50 transition-colors" active-class="text-green-600 bg-green-50">
          <i class="fas fa-newspaper w-5"></i>
          <span>Articles</span>
        </NuxtLink>
        
        <!-- Components Section -->
        <div class="mt-6">
          <div class="flex items-center justify-between px-4 py-2 cursor-pointer" @click="toggleComponentsMenu">
            <span class="text-xs font-semibold text-gray-500 uppercase tracking-wide">Components</span>
            <i :class="['fas fa-chevron-down text-xs text-gray-400 transition-transform duration-200', isComponentsMenuOpen ? 'rotate-180' : '']"></i>
          </div>
          <div v-show="isComponentsMenuOpen" class="mt-2 space-y-1 pl-4">
            <NuxtLink to="/admin/components/form" class="flex items-center space-x-3 px-4 py-2 text-sm text-gray-600 hover:text-green-600 hover:bg-green-50 rounded-lg transition-colors">
              <i class="fas fa-wpforms w-4"></i>
              <span>Forms</span>
            </NuxtLink>
            <NuxtLink to="/admin/components/ui" class="flex items-center space-x-3 px-4 py-2 text-sm text-gray-600 hover:text-green-600 hover:bg-green-50 rounded-lg transition-colors">
              <i class="fas fa-puzzle-piece w-4"></i>
              <span>UI Components</span>
            </NuxtLink>
          </div>
        </div>
      </nav>

      <!-- Sidebar Footer -->
      <div class="p-4 border-t border-gray-200">
        <NuxtLink to="/dashboard" class="flex items-center space-x-3 mb-4 p-2 rounded-xl hover:bg-slate-50 border border-gray-100/50 hover:border-gray-200 transition-all text-left" v-if="authStore.user">
          <div class="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center text-green-700 font-bold flex-shrink-0">
            {{ authStore.user.name?.charAt(0) || 'A' }}
          </div>
          <div class="overflow-hidden flex-grow">
            <p class="text-xs font-semibold text-gray-800 truncate">{{ authStore.user.name }}</p>
            <p class="text-[10px] text-gray-500 truncate">{{ authStore.user.email }}</p>
          </div>
          <i class="fas fa-cog text-xs text-gray-400"></i>
        </NuxtLink>
        <button @click="logout" class="w-full flex items-center space-x-3 px-4 py-2.5 text-red-600 hover:bg-red-50 rounded-lg font-medium transition-colors">
          <i class="fas fa-sign-out-alt"></i>
          <span>Logout</span>
        </button>
      </div>
    </aside>

    <!-- Mobile Sidebar Backdrop -->
    <div 
      v-if="isMobileSidebarOpen" 
      @click="toggleMobileSidebar" 
      class="fixed inset-0 z-30 bg-black/40 lg:hidden"
    ></div>

    <!-- Mobile Sidebar Drawer -->
    <aside 
      :class="[
        'fixed inset-y-0 left-0 z-40 w-64 bg-white shadow-lg flex flex-col transition-transform duration-300 lg:hidden border-r border-gray-200',
        isMobileSidebarOpen ? 'translate-x-0' : '-translate-x-full'
      ]"
    >
      <div class="flex items-center justify-between h-16 px-6 border-b border-gray-200">
        <div class="flex items-center space-x-3">
          <div class="w-8 h-8 bg-gradient-to-r from-green-500 to-blue-500 rounded-lg flex items-center justify-center">
            <i class="fas fa-crown text-white text-sm"></i>
          </div>
          <span class="text-lg font-bold text-gray-800">Admin Panel</span>
        </div>
        <button @click="toggleMobileSidebar" class="p-2 rounded-lg hover:bg-gray-100 text-gray-500">
          <i class="fas fa-times"></i>
        </button>
      </div>

      <nav class="flex-grow px-4 py-6 space-y-1 overflow-y-auto">
        <NuxtLink to="/" @click="toggleMobileSidebar" class="flex items-center space-x-3 px-4 py-3 rounded-lg font-medium text-gray-600 hover:text-green-600 hover:bg-green-50 transition-colors">
          <i class="fas fa-home w-5"></i>
          <span>Kembali ke Home</span>
        </NuxtLink>
        <NuxtLink to="/admin" @click="toggleMobileSidebar" class="flex items-center space-x-3 px-4 py-3 rounded-lg font-medium text-gray-600 hover:text-green-600 hover:bg-green-50 transition-colors" active-class="text-green-600 bg-green-50">
          <i class="fas fa-tachometer-alt w-5"></i>
          <span>Dashboard</span>
        </NuxtLink>
        <NuxtLink to="/admin/charities" @click="toggleMobileSidebar" class="flex items-center space-x-3 px-4 py-3 rounded-lg font-medium text-gray-600 hover:text-green-600 hover:bg-green-50 transition-colors" active-class="text-green-600 bg-green-50">
          <i class="fas fa-heart w-5"></i>
          <span>Charities</span>
        </NuxtLink>
        <NuxtLink to="/admin/donations" @click="toggleMobileSidebar" class="flex items-center space-x-3 px-4 py-3 rounded-lg font-medium text-gray-600 hover:text-green-600 hover:bg-green-50 transition-colors" active-class="text-green-600 bg-green-50">
          <i class="fas fa-donate w-5"></i>
          <span>Donations</span>
        </NuxtLink>
        <NuxtLink to="/admin/users" @click="toggleMobileSidebar" class="flex items-center space-x-3 px-4 py-3 rounded-lg font-medium text-gray-600 hover:text-green-600 hover:bg-green-50 transition-colors" active-class="text-green-600 bg-green-50">
          <i class="fas fa-users w-5"></i>
          <span>Users</span>
        </NuxtLink>
        <NuxtLink to="/admin/articles" @click="toggleMobileSidebar" class="flex items-center space-x-3 px-4 py-3 rounded-lg font-medium text-gray-600 hover:text-green-600 hover:bg-green-50 transition-colors" active-class="text-green-600 bg-green-50">
          <i class="fas fa-newspaper w-5"></i>
          <span>Articles</span>
        </NuxtLink>
      </nav>

      <div class="p-4 border-t border-gray-200">
        <button @click="logout" class="w-full flex items-center space-x-3 px-4 py-2.5 text-red-600 hover:bg-red-50 rounded-lg font-medium transition-colors">
          <i class="fas fa-sign-out-alt"></i>
          <span>Logout</span>
        </button>
      </div>
    </aside>

    <!-- Content Area -->
    <div class="flex-grow flex flex-col min-w-0">
      <!-- Top Bar -->
      <header class="bg-white shadow-sm border-b border-gray-200 h-16 flex items-center justify-between px-6 z-10">
        <div class="flex items-center space-x-4">
          <button @click="toggleMobileSidebar" class="lg:hidden p-2 rounded-lg hover:bg-gray-100 text-gray-600">
            <i class="fas fa-bars"></i>
          </button>
          <button @click="toggleSidebar" class="hidden lg:block p-2 rounded-lg hover:bg-gray-100 text-gray-600">
            <i class="fas fa-bars"></i>
          </button>
          <h1 class="text-lg lg:text-xl font-bold text-gray-800">Admin Dashboard</h1>
        </div>

        <div class="flex items-center space-x-4">


          <!-- Notification Bell -->
          <button class="relative p-2 rounded-lg hover:bg-gray-100 text-gray-600">
            <i class="fas fa-bell"></i>
            <span class="absolute -top-1 -right-1 w-4 h-4 bg-red-500 text-white text-[10px] rounded-full flex items-center justify-center font-bold">3</span>
          </button>

          <!-- Avatar -->
          <div class="w-8 h-8 bg-green-600 text-white rounded-full flex items-center justify-center font-bold text-sm">
            A
          </div>
        </div>
      </header>

      <!-- Main Layout slot -->
      <main class="flex-grow p-6 overflow-y-auto">
        <slot />
      </main>
    </div>
  </div>
</template>
