<script setup>
import { computed } from 'vue'

definePageMeta({
  layout: 'admin'
})

// Fetch all donations, charities, and users from backend API (client-side only to ensure token is present)
const { data: donationsResponse } = await useApiFetch('/donations', { server: false })
const { data: charitiesResponse } = await useApiFetch('/charities', { server: false })
const { data: usersResponse } = await useApiFetch('/users', { server: false })

const donationsList = computed(() => donationsResponse.value?.data || [])
const charitiesList = computed(() => charitiesResponse.value?.data || [])
const usersList = computed(() => usersResponse.value?.data || [])

// Compute dynamic stats
const totalDonationsAmount = computed(() => {
  const total = donationsList.value
    .filter(item => item.status === 'paid')
    .reduce((sum, item) => sum + (item.amount || 0), 0)
  return `Rp ${total.toLocaleString('id-ID')}`
})

const activeCharitiesCount = computed(() => {
  return charitiesList.value.filter(item => item.status === 'active').length
})

const totalUsersCount = computed(() => {
  return usersList.value.length
})

// Dynamic success rate based on paid donations vs total donations
const successRateValue = computed(() => {
  const all = donationsList.value.length
  if (all === 0) return '100%'
  const paid = donationsList.value.filter(item => item.status === 'paid').length
  const percentage = Math.round((paid / all) * 100)
  return `${percentage}%`
})

// Retrieve 5 most recent paid donations
const recentDonations = computed(() => {
  return donationsList.value.filter(item => item.status === 'paid').slice(0, 5)
})

const formatTime = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = now - date
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMins / 60)
  const diffDays = Math.floor(diffHours / 24)

  if (diffMins < 1) return 'Baru saja'
  if (diffMins < 60) return `${diffMins} mnt yang lalu`
  if (diffHours < 24) return `${diffHours} jam yang lalu`
  return `${diffDays} hari yang lalu`
}
</script>

<template>
  <div class="space-y-6">
    <!-- Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
      <div class="bg-white rounded-2xl p-6 border border-gray-100 shadow-sm flex items-center justify-between">
        <div>
          <p class="text-sm text-gray-500 font-medium">Total Donasi</p>
          <h3 class="text-2xl font-bold text-slate-800 mt-1">{{ totalDonationsAmount }}</h3>
          <span class="text-xs text-green-600 font-semibold block mt-1"><i class="fas fa-arrow-up mr-1"></i>+12.5% bulan ini</span>
        </div>
        <div class="w-12 h-12 bg-green-50 text-green-600 rounded-xl flex items-center justify-center text-xl">
          <i class="fas fa-donate"></i>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-6 border border-gray-100 shadow-sm flex items-center justify-between">
        <div>
          <p class="text-sm text-gray-500 font-medium">Program Aktif</p>
          <h3 class="text-2xl font-bold text-slate-800 mt-1">{{ activeCharitiesCount }}</h3>
          <span class="text-xs text-blue-600 font-semibold block mt-1"><i class="fas fa-plus mr-1"></i>3 baru minggu ini</span>
        </div>
        <div class="w-12 h-12 bg-blue-50 text-blue-600 rounded-xl flex items-center justify-center text-xl">
          <i class="fas fa-heart"></i>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-6 border border-gray-100 shadow-sm flex items-center justify-between">
        <div>
          <p class="text-sm text-gray-500 font-medium">Total Pengguna</p>
          <h3 class="text-2xl font-bold text-slate-800 mt-1">{{ totalUsersCount }}</h3>
          <span class="text-xs text-purple-600 font-semibold block mt-1"><i class="fas fa-arrow-up mr-1"></i>+8.2% bulan ini</span>
        </div>
        <div class="w-12 h-12 bg-purple-50 text-purple-600 rounded-xl flex items-center justify-center text-xl">
          <i class="fas fa-users"></i>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-6 border border-gray-100 shadow-sm flex items-center justify-between">
        <div>
          <p class="text-sm text-gray-500 font-medium">Tingkat Keberhasilan</p>
          <h3 class="text-2xl font-bold text-slate-800 mt-1">{{ successRateValue }}</h3>
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
                <td class="py-3 font-semibold text-slate-700">
                  {{ don.anonymous ? 'Hamba Allah (Anonim)' : (don.guestName || (don.user ? don.user.fullName : 'User YPAI')) }}
                </td>
                <td class="py-3 text-gray-500 truncate max-w-[200px]" :title="don.charity?.title">
                  {{ don.charity ? don.charity.title : 'Kampanye' }}
                </td>
                <td class="py-3 font-bold text-slate-800">Rp {{ don.amount?.toLocaleString('id-ID') }}</td>
                <td class="py-3">
                  <span :class="[
                    'px-2 py-1 rounded-full text-[10px] font-bold',
                    don.status === 'paid' ? 'bg-green-50 text-green-700' : 'bg-yellow-50 text-yellow-700'
                  ]">
                    {{ don.status === 'paid' ? 'Lunas' : 'Pending' }}
                  </span>
                </td>
                <td class="py-3 text-gray-400">{{ formatTime(don.createdAt) }}</td>
              </tr>
              <tr v-if="recentDonations.length === 0">
                <td colspan="5" class="py-6 text-center text-gray-400">Belum ada donasi masuk.</td>
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
