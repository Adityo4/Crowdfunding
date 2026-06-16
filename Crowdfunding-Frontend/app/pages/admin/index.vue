<script setup>
import { ref } from 'vue'

definePageMeta({
  layout: 'admin'
})

// Mock or async stats info
const totalDonations = ref('Rp 2.5M')
const activeCharities = ref(24)
const totalUsers = ref(1247)
const successRate = ref('98%')

const recentDonations = ref([
  { id: 1, name: 'Budi Santoso', amount: 'Rp 500,000', campaign: 'Peduli Bencana Alam', status: 'Lunas', date: '2 jam yang lalu' },
  { id: 2, name: 'Siti Rahma', amount: 'Rp 1,000,000', campaign: 'Pendidikan untuk Anak Pesisir', status: 'Lunas', date: '4 jam yang lalu' },
  { id: 3, name: 'Hendra Wijaya', amount: 'Rp 250,000', campaign: 'Pangan Sehat Dhuafa', status: 'Pending', date: '5 jam yang lalu' }
])
</script>

<template>
  <div class="space-y-6">
    <!-- Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
      <div class="bg-white rounded-2xl p-6 border border-gray-100 shadow-sm flex items-center justify-between">
        <div>
          <p class="text-sm text-gray-500 font-medium">Total Donasi</p>
          <h3 class="text-2xl font-bold text-slate-800 mt-1">{{ totalDonations }}</h3>
          <span class="text-xs text-green-600 font-semibold block mt-1"><i class="fas fa-arrow-up mr-1"></i>+12.5% bulan ini</span>
        </div>
        <div class="w-12 h-12 bg-green-50 text-green-600 rounded-xl flex items-center justify-center text-xl">
          <i class="fas fa-donate"></i>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-6 border border-gray-100 shadow-sm flex items-center justify-between">
        <div>
          <p class="text-sm text-gray-500 font-medium">Program Aktif</p>
          <h3 class="text-2xl font-bold text-slate-800 mt-1">{{ activeCharities }}</h3>
          <span class="text-xs text-blue-600 font-semibold block mt-1"><i class="fas fa-plus mr-1"></i>3 baru minggu ini</span>
        </div>
        <div class="w-12 h-12 bg-blue-50 text-blue-600 rounded-xl flex items-center justify-center text-xl">
          <i class="fas fa-heart"></i>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-6 border border-gray-100 shadow-sm flex items-center justify-between">
        <div>
          <p class="text-sm text-gray-500 font-medium">Total Pengguna</p>
          <h3 class="text-2xl font-bold text-slate-800 mt-1">{{ totalUsers }}</h3>
          <span class="text-xs text-purple-600 font-semibold block mt-1"><i class="fas fa-arrow-up mr-1"></i>+8.2% bulan ini</span>
        </div>
        <div class="w-12 h-12 bg-purple-50 text-purple-600 rounded-xl flex items-center justify-center text-xl">
          <i class="fas fa-users"></i>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-6 border border-gray-100 shadow-sm flex items-center justify-between">
        <div>
          <p class="text-sm text-gray-500 font-medium">Tingkat Keberhasilan</p>
          <h3 class="text-2xl font-bold text-slate-800 mt-1">{{ successRate }}</h3>
          <span class="text-xs text-emerald-600 font-semibold block mt-1"><i class="fas fa-check mr-1"></i>Sangat Stabil</span>
        </div>
        <div class="w-12 h-12 bg-emerald-50 text-emerald-600 rounded-xl flex items-center justify-center text-xl">
          <i class="fas fa-chart-line"></i>
        </div>
      </div>
    </div>

    <!-- Main Section Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Recent Donations Table -->
      <div class="bg-white rounded-2xl p-6 border border-gray-100 shadow-sm lg:col-span-2 space-y-4">
        <div class="flex items-center justify-between">
          <h3 class="font-bold text-slate-800 text-base">Donasi Terbaru</h3>
          <NuxtLink to="/admin/donations" class="text-xs text-green-600 font-bold hover:underline">Lihat Semua</NuxtLink>
        </div>
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr class="border-b border-gray-100 text-gray-400 font-semibold uppercase">
                <th class="py-3">Donatur</th>
                <th class="py-3">Kampanye</th>
                <th class="py-3">Nominal</th>
                <th class="py-3">Status</th>
                <th class="py-3">Waktu</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-50">
              <tr v-for="don in recentDonations" :key="don.id" class="hover:bg-slate-50 transition-colors">
                <td class="py-3 font-semibold text-slate-700">{{ don.name }}</td>
                <td class="py-3 text-gray-500 truncate max-w-[150px]">{{ don.campaign }}</td>
                <td class="py-3 font-bold text-slate-800">{{ don.amount }}</td>
                <td class="py-3">
                  <span :class="[
                    'px-2 py-1 rounded-full text-[10px] font-bold',
                    don.status === 'Lunas' ? 'bg-green-50 text-green-700' : 'bg-yellow-50 text-yellow-700'
                  ]">
                    {{ don.status }}
                  </span>
                </td>
                <td class="py-3 text-gray-400">{{ don.date }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Quick Action Panel -->
      <div class="bg-white rounded-2xl p-6 border border-gray-100 shadow-sm space-y-4">
        <h3 class="font-bold text-slate-800 text-base">Aksi Cepat</h3>
        <div class="space-y-2">
          <NuxtLink to="/charities/create" class="flex items-center space-x-3 p-3 bg-green-50 rounded-xl hover:bg-green-100 transition-colors text-green-700">
            <i class="fas fa-plus"></i>
            <span class="text-xs font-bold">Buat Kampanye Baru</span>
          </NuxtLink>
          <NuxtLink to="/admin/articles" class="flex items-center space-x-3 p-3 bg-blue-50 rounded-xl hover:bg-blue-100 transition-colors text-blue-700">
            <i class="fas fa-pen"></i>
            <span class="text-xs font-bold">Tulis Artikel Baru</span>
          </NuxtLink>
          <NuxtLink to="/admin/users" class="flex items-center space-x-3 p-3 bg-purple-50 rounded-xl hover:bg-purple-100 transition-colors text-purple-700">
            <i class="fas fa-user-plus"></i>
            <span class="text-xs font-bold">Kelola Pengguna</span>
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>
