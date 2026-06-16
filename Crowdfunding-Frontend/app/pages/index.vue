<script setup>
import { computed } from 'vue'

// Fetch campaign data from backend
const { data: charitiesResponse } = await useApiFetch('/charities?perPage=3&status=active')
const charities = computed(() => charitiesResponse.value?.data || [])

// Fetch articles from backend
const { data: articlesResponse } = await useApiFetch('/articles')
const articles = computed(() => articlesResponse.value?.data?.slice(0, 3) || [])
</script>

<template>
  <div>
    <!-- Hero Section - Split Layout -->
    <section id="home" class="py-8 md:py-16 bg-white overflow-hidden">
      <div class="max-w-7xl mx-auto px-4 md:px-6">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 md:gap-12 items-center min-h-[60vh] md:min-h-[80vh]">
          <!-- Content -->
          <div class="animate-fade-in">
            <div class="inline-flex items-center px-4 py-2 bg-green-50 rounded-full text-green-600 font-medium mb-6 hover:bg-green-100 transition-all duration-500 transform hover:scale-105">
              <i class="fas fa-heart mr-2 animate-pulse text-green-500"></i>
              Making Hope Happen
            </div>
            <h1 class="text-3xl sm:text-4xl md:text-5xl lg:text-6xl font-bold text-slate-900 mb-4 md:mb-6 leading-tight">
              Together We Can
              <span class="text-green-600 bg-gradient-to-r from-green-500 to-emerald-500 bg-clip-text text-transparent animate-gradient-x">Change Lives</span>
            </h1>
            <p class="text-base md:text-lg text-gray-600 mb-4 md:mb-6 leading-relaxed">
              Bergabunglah dengan ribuan donatur yang telah menciptakan perubahan positif untuk masyarakat yang membutuhkan. Kontribusi Anda sangat berarti bagi masa depan mereka.
            </p>
            <div class="flex flex-col sm:flex-row gap-3 md:gap-4">
              <NuxtLink to="/charities" class="bg-green-600 text-white px-6 md:px-8 py-3 md:py-4 rounded-full text-base md:text-lg font-semibold hover:bg-green-700 text-center transform hover:scale-105 hover:-translate-y-1 transition-all duration-500 shadow-lg hover:shadow-2xl group">
                <i class="fas fa-heart mr-2 group-hover:animate-bounce"></i>
                Donate Now
              </NuxtLink>
              <NuxtLink to="/#about" class="border-2 border-green-600 text-green-600 px-6 md:px-8 py-3 md:py-4 rounded-full text-base md:text-lg font-semibold hover:bg-green-600 hover:text-white text-center hover:scale-105 hover:-translate-y-1 transition-all duration-500 shadow-lg hover:shadow-xl group">
                Watch Our Story
              </NuxtLink>
            </div>
            <!-- Stats -->
            <div class="grid grid-cols-3 gap-4 md:gap-6 mt-6 md:mt-8">
              <div class="text-center group hover:scale-110 transition-all duration-500 cursor-pointer">
                <div class="text-2xl md:text-3xl font-bold text-green-600 mb-1 md:mb-2 group-hover:text-green-700 transition-colors duration-500">50K+</div>
                <div class="text-xs md:text-sm text-gray-600 group-hover:text-gray-800 transition-colors duration-500">Lives Impacted</div>
              </div>
              <div class="text-center group hover:scale-110 transition-all duration-500 cursor-pointer">
                <div class="text-2xl md:text-3xl font-bold text-green-600 mb-1 md:mb-2 group-hover:text-green-700 transition-colors duration-500">25</div>
                <div class="text-xs md:text-sm text-gray-600 group-hover:text-gray-800 transition-colors duration-500">Countries</div>
              </div>
              <div class="text-center group hover:scale-110 transition-all duration-500 cursor-pointer">
                <div class="text-2xl md:text-3xl font-bold text-green-600 mb-1 md:mb-2 group-hover:text-green-700 transition-colors duration-500">$2M+</div>
                <div class="text-xs md:text-sm text-gray-600 group-hover:text-gray-800 transition-colors duration-500">Raised</div>
              </div>
            </div>
          </div>
          
          <!-- Image -->
          <div class="relative animate-slide-up">
            <div class="relative">
              <img src="https://images.unsplash.com/photo-1559027615-cd4628902d85?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1000&q=80" 
                   alt="Volunteers working together" 
                   class="w-full h-64 md:h-96 object-cover rounded-2xl md:rounded-3xl shadow-2xl">
              <div class="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent rounded-3xl"></div>
            </div>
            <!-- Floating Cards -->
            <div class="absolute -top-3 -right-3 md:-top-6 md:-right-6 bg-white p-3 md:p-6 rounded-xl md:rounded-2xl shadow-xl hover:shadow-2xl transition-all duration-500 transform hover:scale-105 hover:-translate-y-2 animate-float">
              <div class="flex items-center space-x-2 md:space-x-3">
                <div class="w-8 h-8 md:w-12 md:h-12 bg-green-50 rounded-full flex items-center justify-center">
                  <i class="fas fa-users text-green-600 text-sm md:text-xl"></i>
                </div>
                <div>
                  <div class="text-lg md:text-2xl font-bold text-slate-900">1,200+</div>
                  <div class="text-gray-600 text-xs md:text-sm">Active Volunteers</div>
                </div>
              </div>
            </div>
            <div class="absolute -bottom-3 -left-3 md:-bottom-6 md:-left-6 bg-white p-3 md:p-6 rounded-xl md:rounded-2xl shadow-xl hover:shadow-2xl transition-all duration-500 transform hover:scale-105 hover:-translate-y-2 animate-float" style="animation-delay: 0.5s;">
              <div class="flex items-center space-x-2 md:space-x-3">
                <div class="w-8 h-8 md:w-12 md:h-12 bg-emerald-50 rounded-full flex items-center justify-center">
                  <i class="fas fa-check text-emerald-600 text-sm md:text-xl"></i>
                </div>
                <div>
                  <div class="text-lg md:text-2xl font-bold text-slate-900">98%</div>
                  <div class="text-gray-600 text-xs md:text-sm">Success Rate</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Donation Programs Section -->
    <section id="programs" class="py-12 md:py-16 bg-slate-50">
      <div class="max-w-6xl mx-auto px-4 md:px-6">
        <div class="text-center mb-8 md:mb-12 animate-fade-in">
          <h2 class="text-2xl md:text-3xl lg:text-4xl font-light text-slate-900 mb-3 md:mb-4">Program Penggalangan Dana</h2>
          <p class="text-sm md:text-base text-gray-600 max-w-2xl mx-auto leading-relaxed">
            Temukan berbagai inisiatif bantuan sosial kami dan bantu mereka dengan kontribusi nyata Anda hari ini.
          </p>
        </div>
        
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-6">
          <!-- Dinamis Campaign List dari Backend -->
          <div v-for="charity in charities" :key="charity.id" class="group relative bg-white rounded-2xl shadow-sm hover:shadow-2xl transition-all duration-700 overflow-hidden animate-slide-up transform hover:-translate-y-2 hover:scale-105">
            <!-- Image Container -->
            <div class="relative h-48 overflow-hidden">
              <img :src="charity.coverImageUrl || 'https://images.unsplash.com/photo-1523050854058-8df90110c9f1?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80'" 
                   :alt="charity.title" 
                   class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-700">
              <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent"></div>
              
              <!-- Progress Badge -->
              <div class="absolute bottom-4 left-4 bg-white/90 backdrop-blur-sm rounded-full px-3 py-1">
                <span class="text-sm font-semibold text-gray-800">
                  {{ ((charity.currentAmount / charity.targetAmount) * 100).toFixed(0) }}% Funded
                </span>
              </div>
            </div>
            
            <!-- Content -->
            <div class="p-4 md:p-5">
              <span class="text-xs font-semibold uppercase text-green-600 tracking-wider">{{ charity.category?.name || 'Sosial' }}</span>
              <h3 class="text-base md:text-lg font-bold text-slate-900 mt-1 mb-2 group-hover:text-green-600 transition-colors duration-300">
                <NuxtLink :to="`/charities/${charity.slug}`">{{ charity.title }}</NuxtLink>
              </h3>
              <p class="text-gray-600 mb-3 md:mb-4 leading-relaxed text-xs md:text-sm line-clamp-2">{{ charity.description }}</p>
              
              <!-- Progress Bar -->
              <div class="mb-3 md:mb-4">
                <div class="flex justify-between text-xs text-gray-500 mb-2">
                  <span>Terkumpul: Rp {{ charity.currentAmount.toLocaleString('id-ID') }}</span>
                  <span>Target: Rp {{ charity.targetAmount.toLocaleString('id-ID') }}</span>
                </div>
                <div class="w-full bg-gray-200 rounded-full h-1 md:h-1.5 overflow-hidden">
                  <div class="bg-gradient-to-r from-green-500 to-emerald-500 h-1 md:h-1.5 rounded-full transition-all duration-1000 ease-out" 
                       :style="{ width: Math.min((charity.currentAmount / charity.targetAmount) * 100, 100) + '%' }"></div>
                </div>
              </div>
              
              <!-- Button -->
              <NuxtLink :to="`/charities/${charity.slug}`" class="block w-full text-center bg-gradient-to-r from-green-500 to-emerald-600 text-white py-2 md:py-3 rounded-lg md:rounded-xl font-medium hover:from-emerald-600 hover:to-green-500 transform hover:scale-105 transition-all duration-300 shadow-lg hover:shadow-xl text-sm md:text-base">
                Detail Donasi
              </NuxtLink>
            </div>
          </div>

          <!-- Placeholder jika kosong -->
          <div v-if="charities.length === 0" class="col-span-full text-center py-12 text-gray-500">
            Belum ada program penggalangan dana aktif saat ini.
          </div>
        </div>
        
        <!-- Show More Button -->
        <div class="text-center mt-12 animate-fade-in">
          <NuxtLink to="/charities" class="inline-block bg-gradient-to-r from-green-600 to-emerald-600 text-white px-8 py-4 rounded-full font-semibold hover:from-emerald-600 hover:to-green-600 transform hover:scale-105 transition-all duration-300 shadow-lg hover:shadow-xl">
            <i class="fas fa-plus mr-2"></i>
            Lihat Semua Program
          </NuxtLink>
        </div>
      </div>
    </section>

    <!-- About Us Section -->
    <section id="about" class="py-12 md:py-16 bg-white">
      <div class="max-w-6xl mx-auto px-4 md:px-6">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 md:gap-12 items-center">
          <!-- Image -->
          <div class="order-2 lg:order-1 animate-fade-in">
            <div class="relative">
              <div class="h-64 md:h-96 bg-gradient-to-br from-green-500/10 to-slate-900/10 rounded-2xl md:rounded-3xl overflow-hidden relative">
                <img src="https://images.unsplash.com/photo-1559027615-cd4628902d85?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80" 
                     alt="Volunteers working together" 
                     class="w-full h-full object-cover">
                <div class="absolute inset-0 bg-gradient-to-br from-green-500/20 to-slate-900/20"></div>
                <div class="absolute inset-0 flex items-center justify-center">
                  <i class="fas fa-users text-white text-4xl md:text-6xl drop-shadow-lg animate-float"></i>
                </div>
              </div>
            </div>
          </div>
          
          <!-- Content -->
          <div class="order-1 lg:order-2 animate-slide-up">
            <h2 class="text-2xl md:text-3xl lg:text-4xl font-light text-slate-900 mb-4 md:mb-6">Tentang Yayasan Peduli Amal Indonesia</h2>
            <p class="text-sm md:text-base text-gray-600 mb-3 md:mb-4 leading-relaxed">
              Didirikan sejak 2015, Yayasan Peduli Amal Indonesia telah berkomitmen penuh untuk menyalurkan kebaikan kepada mereka yang membutuhkan dengan transparan dan akuntabel.
            </p>
            <p class="text-sm md:text-base text-gray-600 mb-6 md:mb-8 leading-relaxed">
              Misi utama kami adalah memberdayakan masyarakat terpencil melalui penyediaan akses pendidikan, kesehatan, rumah layak huni, dan pemenuhan kebutuhan pangan harian yang layak.
            </p>
            <div class="grid grid-cols-2 gap-4 md:gap-6">
              <div class="text-center p-3 md:p-4 bg-slate-50 rounded-xl md:rounded-2xl">
                <div class="text-xl md:text-2xl font-semibold text-green-600 mb-1">50K+</div>
                <div class="text-gray-600 font-medium text-xs md:text-sm">Lives Impacted</div>
              </div>
              <div class="text-center p-3 md:p-4 bg-slate-50 rounded-xl md:rounded-2xl">
                <div class="text-xl md:text-2xl font-semibold text-slate-900 mb-1">25</div>
                <div class="text-gray-600 font-medium text-xs md:text-sm">Countries</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Articles Section -->
    <section class="py-12 md:py-16 bg-white border-t border-gray-100">
      <div class="max-w-6xl mx-auto px-4 md:px-6">
        <div class="text-center mb-8 md:mb-12 animate-fade-in">
          <h2 class="text-2xl md:text-3xl lg:text-4xl font-light text-slate-900 mb-3 md:mb-4">Artikel & Kabar Terbaru</h2>
          <p class="text-sm md:text-base text-gray-600 max-w-2xl mx-auto">
            Dapatkan pembaruan kabar mengenai aksi sosial kami langsung di lapangan.
          </p>
        </div>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div v-for="article in articles" :key="article.id" class="group bg-white rounded-2xl border border-gray-100 overflow-hidden shadow-sm hover:shadow-xl transition-all duration-300">
            <div class="h-48 overflow-hidden relative">
              <img :src="article.coverImageUrl || 'https://images.unsplash.com/photo-1559027615-cd4628902d85?auto=format&fit=crop&w=1000&q=80'" 
                   :alt="article.title" 
                   class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500">
            </div>
            <div class="p-5">
              <span class="text-xs text-gray-400">{{ new Date(article.publishedAt || article.createdAt).toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' }) }}</span>
              <h3 class="text-base font-bold text-slate-900 mt-2 mb-3 line-clamp-2 group-hover:text-green-600 transition-colors">
                <NuxtLink :to="`/articles/${article.slug}`">{{ article.title }}</NuxtLink>
              </h3>
              <NuxtLink :to="`/articles/${article.slug}`" class="text-sm font-semibold text-green-600 hover:text-green-700">Baca Selengkapnya &rarr;</NuxtLink>
            </div>
          </div>

          <div v-if="articles.length === 0" class="col-span-full text-center py-12 text-gray-500">
            Belum ada artikel yang dipublikasikan saat ini.
          </div>
        </div>
      </div>
    </section>

    <!-- Call to Action Banner -->
    <section class="py-12 md:py-16 bg-gradient-to-br from-green-600 to-slate-900 relative overflow-hidden">
      <div class="absolute inset-0">
        <img src="https://images.unsplash.com/photo-1469571486292-0ba58a3f068b?ixlib=rb-4.0.3&auto=format&fit=crop&w=2070&q=80" 
             alt="People holding hands in unity" 
             class="w-full h-full object-cover opacity-20">
      </div>
      <div class="relative z-10 max-w-6xl mx-auto px-4 md:px-6 text-center">
        <div class="animate-fade-in text-white">
          <h2 class="text-2xl md:text-3xl lg:text-4xl font-bold mb-4 md:mb-6">
            Siap Untuk Membantu Sesama?
          </h2>
          <p class="text-sm md:text-base text-gray-200 mb-6 md:mb-8 max-w-3xl mx-auto leading-relaxed">
            Bergabunglah dengan gerakan amal global ini. Setiap rupiah kontribusi Anda akan sangat berarti bagi kehidupan jutaan saudara kita.
          </p>
          <NuxtLink to="/charities" class="inline-block bg-white text-green-600 px-8 md:px-12 py-3 md:py-4 rounded-full text-base md:text-lg font-semibold hover:bg-gray-100 transform hover:scale-105 transition-all duration-300 shadow-xl hover:shadow-2xl">
            Mulai Donasi Sekarang
          </NuxtLink>
        </div>
      </div>
    </section>
  </div>
</template>
