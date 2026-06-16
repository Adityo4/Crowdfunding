<script setup>
import { ref, watch } from 'vue'

const searchQuery = ref('')
const selectedCategory = ref('')
const selectedSort = ref('newest')
const currentPage = ref(1)

// Fetch categories from backend
const { data: categoriesResponse } = await useApiFetch('/categories')
const categories = computed(() => categoriesResponse.value?.data || [])

// Watch filters and trigger fetch dynamically
const { data: charitiesResponse, refresh } = await useAsyncData('charities', () => 
  $fetch('http://localhost:8080/v1/charities', {
    params: {
      q: searchQuery.value,
      category: selectedCategory.value,
      sort: selectedSort.value,
      page: currentPage.value,
      perPage: 9,
      status: 'active'
    }
  }), {
    watch: [selectedCategory, selectedSort, currentPage]
  }
)

const charities = computed(() => charitiesResponse.value?.data || [])
const pagination = computed(() => charitiesResponse.value?.pagination || {})

const handleSearch = () => {
  currentPage.value = 1
  refresh()
}

const changePage = (page) => {
  if (page >= 1 && page <= (pagination.value?.totalPages || 1)) {
    currentPage.value = page
  }
}

const filterCategory = (slug) => {
  selectedCategory.value = slug
  currentPage.value = 1
}
</script>

<template>
  <div class="bg-slate-50 min-h-screen py-10">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      
      <!-- Header Section -->
      <div class="text-center mb-10">
        <h1 class="text-3xl font-extrabold text-slate-900 tracking-tight sm:text-4xl mb-4">
          Program Penggalangan Dana
        </h1>
        <p class="text-base text-gray-600 max-w-2xl mx-auto">
          Temukan program aksi sosial kami dan mari wujudkan perubahan nyata.
        </p>
      </div>

      <!-- Search & Filters Container -->
      <div class="bg-white rounded-2xl shadow-sm border border-gray-100 p-6 mb-8 flex flex-col md:flex-row md:items-center justify-between gap-4">
        <!-- Search bar -->
        <form @submit.prevent="handleSearch" class="relative flex-grow max-w-md">
          <input v-model="searchQuery" type="text" placeholder="Cari program donasi..." class="w-full px-4 py-2.5 pl-10 border border-gray-200 rounded-xl focus:border-green-600 focus:ring-2 focus:ring-green-600/10 focus:outline-none transition-all duration-300">
          <i class="fas fa-search absolute left-3.5 top-1/2 transform -translate-y-1/2 text-gray-400"></i>
        </form>

        <!-- Sorting & Filters -->
        <div class="flex flex-wrap items-center gap-3">
          <select v-model="selectedSort" class="px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:border-green-600">
            <option value="newest">Terbaru</option>
            <option value="target_amount">Target Tertinggi</option>
            <option value="urgent">Mendesak</option>
          </select>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-4 gap-8">
        <!-- Sidebar Categories Filter -->
        <div class="lg:col-span-1">
          <div class="bg-white rounded-2xl shadow-sm border border-gray-100 p-5 sticky top-24">
            <h3 class="text-base font-bold text-slate-900 mb-4 pb-2 border-b border-gray-100">Kategori</h3>
            <ul class="space-y-1">
              <li>
                <button @click="filterCategory('')" :class="{'bg-green-50 text-green-600 font-semibold': selectedCategory === '', 'text-gray-600 hover:bg-gray-50 hover:text-green-600': selectedCategory !== ''}" class="w-full text-left px-3 py-2 rounded-lg text-sm transition-all">
                  Semua Kategori
                </button>
              </li>
              <li v-for="category in categories" :key="category.id">
                <button @click="filterCategory(category.slug)" :class="{'bg-green-50 text-green-600 font-semibold': selectedCategory === category.slug, 'text-gray-600 hover:bg-gray-50 hover:text-green-600': selectedCategory !== category.slug}" class="w-full text-left px-3 py-2 rounded-lg text-sm transition-all">
                  {{ category.name }}
                </button>
              </li>
            </ul>
          </div>
        </div>

        <!-- Grid Results -->
        <div class="lg:col-span-3">
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
            <div v-for="charity in charities" :key="charity.id" class="group relative bg-white rounded-2xl shadow-sm hover:shadow-xl border border-gray-100 transition-all duration-300 overflow-hidden flex flex-col justify-between">
              <!-- Image Container -->
              <div class="relative h-44 overflow-hidden">
                <img :src="charity.coverImageUrl || 'https://images.unsplash.com/photo-1523050854058-8df90110c9f1?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80'" 
                     :alt="charity.title" 
                     class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500">
                <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent"></div>
                
                <!-- Progress Badge -->
                <div class="absolute bottom-4 left-4 bg-white/90 backdrop-blur-sm rounded-full px-3 py-1">
                  <span class="text-xs font-semibold text-gray-800">
                    {{ ((charity.currentAmount / charity.targetAmount) * 100).toFixed(0) }}% Funded
                  </span>
                </div>
              </div>
              
              <!-- Content -->
              <div class="p-5 flex-grow flex flex-col justify-between">
                <div>
                  <span class="text-[10px] font-bold uppercase text-green-600 tracking-wider bg-green-50 px-2 py-0.5 rounded-full">{{ charity.category?.name || 'Sosial' }}</span>
                  <h3 class="text-base font-bold text-slate-900 mt-2 mb-2 line-clamp-2 hover:text-green-600 transition-colors">
                    <NuxtLink :to="`/charities/${charity.slug}`">{{ charity.title }}</NuxtLink>
                  </h3>
                  <p class="text-gray-500 mb-4 text-xs line-clamp-2">{{ charity.description }}</p>
                </div>
                
                <!-- Progress Bar -->
                <div>
                  <div class="flex justify-between text-[10px] text-gray-500 mb-1.5">
                    <span>Terkumpul: Rp {{ charity.currentAmount.toLocaleString('id-ID') }}</span>
                    <span>Target: Rp {{ charity.targetAmount.toLocaleString('id-ID') }}</span>
                  </div>
                  <div class="w-full bg-gray-100 rounded-full h-1 overflow-hidden mb-4">
                    <div class="bg-gradient-to-r from-green-500 to-emerald-500 h-1 rounded-full transition-all duration-500" 
                         :style="{ width: Math.min((charity.currentAmount / charity.targetAmount) * 100, 100) + '%' }"></div>
                  </div>
                  
                  <!-- Button -->
                  <NuxtLink :to="`/charities/${charity.slug}`" class="block w-full text-center bg-green-600 hover:bg-green-700 text-white py-2 rounded-xl text-sm font-semibold transition-all shadow-sm">
                    Detail Donasi
                  </NuxtLink>
                </div>
              </div>
            </div>

            <!-- Empty State -->
            <div v-if="charities.length === 0" class="col-span-full text-center py-16 bg-white rounded-2xl border border-gray-100 text-gray-500">
              <i class="fas fa-heart-broken text-4xl text-gray-300 mb-3"></i>
              <p>Tidak ada program penggalangan dana yang sesuai.</p>
            </div>
          </div>

          <!-- Pagination -->
          <div v-if="pagination.totalPages > 1" class="flex justify-center items-center gap-2">
            <button @click="changePage(currentPage - 1)" :disabled="!pagination.hasPrev" class="px-4 py-2 border border-gray-200 rounded-xl text-sm font-semibold bg-white text-gray-600 hover:bg-gray-50 disabled:opacity-50">
              Sebelumnya
            </button>
            <span class="text-sm font-semibold text-gray-700">Halaman {{ currentPage }} dari {{ pagination.totalPages }}</span>
            <button @click="changePage(currentPage + 1)" :disabled="!pagination.hasNext" class="px-4 py-2 border border-gray-200 rounded-xl text-sm font-semibold bg-white text-gray-600 hover:bg-gray-50 disabled:opacity-50">
              Selanjutnya
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
