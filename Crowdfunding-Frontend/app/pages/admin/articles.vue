<script setup>
import { ref, computed } from 'vue'

definePageMeta({
  layout: 'admin'
})

// Search query
const searchQuery = ref('')

// Fetch articles dynamically (for admin articles overview)
const { data: articlesResponse, refresh } = await useApiFetch('/articles')
const articles = computed(() => {
  let list = articlesResponse.value?.data || []
  if (searchQuery.value) {
    list = list.filter(item => 
      item.title?.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }
  return list
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h2 class="text-xl font-bold text-slate-800">Manajemen Artikel</h2>
        <p class="text-xs text-gray-500">Kelola publikasi kabar, rilis berita, dan cerita sukses aksi sosial YPAI.</p>
      </div>
      <button @click="alert('Fitur tulis artikel akan segera hadir!')" class="inline-flex items-center justify-center bg-green-600 hover:bg-green-700 text-white font-bold text-xs px-4 py-2.5 rounded-lg transition-colors">
        <i class="fas fa-plus mr-2"></i>Tulis Artikel Baru
      </button>
    </div>

    <!-- Filter Card -->
    <div class="bg-white border border-gray-100 rounded-2xl p-6 shadow-sm">
      <div>
        <label class="block text-xs font-semibold text-gray-700 mb-2">Cari Judul Artikel</label>
        <div class="relative">
          <input v-model="searchQuery" type="text" placeholder="Masukkan judul..." class="w-full pl-9 pr-3 py-2 text-xs border border-gray-200 rounded-lg focus:outline-none focus:border-green-500 bg-white text-slate-900">
          <i class="fas fa-search absolute left-3 top-3 text-gray-400 text-xs"></i>
        </div>
      </div>
    </div>

    <!-- Table List -->
    <div class="bg-white border border-gray-100 rounded-2xl shadow-sm overflow-hidden">
      <div class="p-6 border-b border-gray-100 flex items-center justify-between">
        <h3 class="font-bold text-slate-800 text-sm">Semua Artikel</h3>
        <span class="text-xs text-gray-400 font-semibold">{{ articles.length }} Artikel ditemukan</span>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs border-collapse">
          <thead class="bg-slate-50 border-b border-gray-100 text-gray-500 font-semibold uppercase">
            <tr>
              <th class="px-6 py-3">Artikel</th>
              <th class="px-6 py-3">Tanggal Dibuat</th>
              <th class="px-6 py-3 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50">
            <tr v-for="item in articles" :key="item.id" class="hover:bg-slate-50 transition-colors">
              <td class="px-6 py-4">
                <div class="flex items-center space-x-3">
                  <img :src="item.coverImageUrl || 'https://images.unsplash.com/photo-1559027615-cd4628902d85?auto=format&fit=crop&w=100&q=80'" class="w-10 h-10 object-cover rounded-lg">
                  <div>
                    <h4 class="font-bold text-slate-800">{{ item.title }}</h4>
                    <span class="text-[10px] text-gray-400">Slug: {{ item.slug }}</span>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 text-gray-400">
                {{ new Date(item.createdAt).toLocaleDateString('id-ID', { year: 'numeric', month: 'short', day: 'numeric' }) }}
              </td>
              <td class="px-6 py-4 text-right">
                <NuxtLink :to="`/articles/${item.slug}`" class="p-1.5 text-gray-400 hover:text-green-600 transition-colors">
                  <i class="fas fa-eye text-sm"></i>
                </NuxtLink>
              </td>
            </tr>
            <tr v-if="articles.length === 0">
              <td colspan="3" class="text-center py-12 text-gray-400">Tidak ada artikel tersedia.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
