<script setup>
import { ref, computed } from 'vue'

definePageMeta({
  layout: 'admin'
})

// Search query
const searchQuery = ref('')

// Fetch registered users (for admin users list)
const { data: usersResponse, refresh } = await useApiFetch('/users', { server: false })
const users = computed(() => {
  let list = usersResponse.value?.data || []
  if (searchQuery.value) {
    list = list.filter(item => 
      item.fullName?.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      item.email?.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }
  return list
})
</script>

<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-xl font-bold text-slate-800">Manajemen Pengguna</h2>
      <p class="text-xs text-gray-500">Kelola dan pantau seluruh akun pengguna terdaftar pada platform YPAI.</p>
    </div>

    <!-- Filter Card -->
    <div class="bg-white border border-gray-100 rounded-2xl p-6 shadow-sm">
      <div>
        <label class="block text-xs font-semibold text-gray-700 mb-2">Cari Nama / Email Pengguna</label>
        <div class="relative">
          <input v-model="searchQuery" type="text" placeholder="Masukkan nama atau email..." class="w-full pl-9 pr-3 py-2 text-xs border border-gray-200 rounded-lg focus:outline-none focus:border-green-500 bg-white text-slate-900">
          <i class="fas fa-search absolute left-3 top-3 text-gray-400 text-xs"></i>
        </div>
      </div>
    </div>

    <!-- Table List -->
    <div class="bg-white border border-gray-100 rounded-2xl shadow-sm overflow-hidden">
      <div class="p-6 border-b border-gray-100 flex items-center justify-between">
        <h3 class="font-bold text-slate-800 text-sm">Semua Pengguna Terdaftar</h3>
        <span class="text-xs text-gray-400 font-semibold">{{ users.length }} Akun ditemukan</span>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs border-collapse">
          <thead class="bg-slate-50 border-b border-gray-100 text-gray-500 font-semibold uppercase">
            <tr>
              <th class="px-6 py-3">Nama</th>
              <th class="px-6 py-3">Email</th>
              <th class="px-6 py-3">No. Telepon</th>
              <th class="px-6 py-3">Peran</th>
              <th class="px-6 py-3">Tanggal Bergabung</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50">
            <tr v-for="item in users" :key="item.id" class="hover:bg-slate-50 transition-colors">
              <td class="px-6 py-4 font-bold text-slate-700">{{ item.fullName }}</td>
              <td class="px-6 py-4 text-gray-600">{{ item.email }}</td>
              <td class="px-6 py-4 text-gray-500">{{ item.phoneNumber || '-' }}</td>
              <td class="px-6 py-4">
                <span :class="[
                  'px-2 py-1 rounded-full text-[10px] font-bold',
                  item.role === 'admin' ? 'bg-purple-50 text-purple-700' : 'bg-blue-50 text-blue-700'
                ]">
                  {{ item.role }}
                </span>
              </td>
              <td class="px-6 py-4 text-gray-400">
                {{ new Date(item.createdAt).toLocaleDateString('id-ID', { year: 'numeric', month: 'short', day: 'numeric' }) }}
              </td>
            </tr>
            <tr v-if="users.length === 0">
              <td colspan="5" class="text-center py-12 text-gray-400">Tidak ada data pengguna terdaftar.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
