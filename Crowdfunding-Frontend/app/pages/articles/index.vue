<script setup>
import { computed } from 'vue'

// Fetch articles from backend
const { data: articlesResponse } = await useFetch('http://localhost:8080/v1/articles')
const articles = computed(() => articlesResponse.value?.data || [])
</script>

<template>
  <div class="bg-slate-50 min-h-screen py-10">
    <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
      
      <!-- Page Header -->
      <div class="text-center mb-10">
        <div class="inline-flex items-center px-3 py-1 bg-green-50 rounded-full text-green-600 text-xs font-semibold mb-4">
          <i class="fas fa-newspaper mr-1.5"></i>
          <span>Kabar Terbaru</span>
        </div>
        <h1 class="text-3xl font-extrabold text-slate-900 tracking-tight sm:text-4xl">
          Artikel & Catatan Aksi Sosial
        </h1>
        <p class="text-base text-gray-600 max-w-2xl mx-auto mt-2">
          Ikuti terus kisah perjuangan relawan dan senyuman para penerima manfaat dari program donasi Anda.
        </p>
      </div>

      <!-- Grid Articles -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="article in articles" :key="article.id" class="group bg-white rounded-2xl border border-gray-100 overflow-hidden shadow-sm hover:shadow-xl transition-all duration-300 flex flex-col justify-between">
          <div>
            <div class="h-48 overflow-hidden relative bg-gray-100">
              <img :src="article.coverImageUrl || 'https://images.unsplash.com/photo-1559027615-cd4628902d85?auto=format&fit=crop&w=1000&q=80'" 
                   :alt="article.title" 
                   class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500">
            </div>
            <div class="p-6">
              <span class="text-xs text-gray-400 font-semibold block mb-2">
                {{ new Date(article.publishedAt || article.createdAt).toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' }) }}
              </span>
              <h3 class="text-base font-bold text-slate-900 line-clamp-2 group-hover:text-green-600 transition-colors">
                <NuxtLink :to="`/articles/${article.slug}`">{{ article.title }}</NuxtLink>
              </h3>
              <p class="text-gray-500 text-xs mt-2 line-clamp-3" v-html="article.content"></p>
            </div>
          </div>
          <div class="px-6 pb-6 pt-2">
            <NuxtLink :to="`/articles/${article.slug}`" class="text-sm font-semibold text-green-600 hover:text-green-700">
              Baca Selengkapnya &rarr;
            </NuxtLink>
          </div>
        </div>

        <!-- Empty state -->
        <div v-if="articles.length === 0" class="col-span-full text-center py-16 bg-white rounded-2xl border border-gray-100 text-gray-500">
          <i class="fas fa-newspaper text-4xl text-gray-300 mb-3"></i>
          <p>Belum ada artikel yang dipublikasikan saat ini.</p>
        </div>
      </div>
    </div>
  </div>
</template>
