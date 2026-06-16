<script setup>
import { ref, computed } from 'vue'

definePageMeta({
  layout: 'admin'
})

// Search query
const searchQuery = ref('')

// Fetch donations dynamically (for admin overview)
const { data: donationsResponse, refresh } = await useApiFetch('/donations?limit=50')
const donations = computed(() => {
  let list = donationsResponse.value?.data || []
  if (searchQuery.value) {
    list = list.filter(item => 
      item.guestName?.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      item.id?.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }
  return list
})
</script>

<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-xl font-bold text-slate-800">Riwayat Donasi</h2>
      <p class="text-xs text-gray-500">Pantau seluruh aliran dana dan transaksi dari donatur secara real-time.</p>
    </div>

    <!-- Filter Card -->
    <div class="bg-white border border-gray-100 rounded-2xl p-6 shadow-sm">
      <div>
        <label class="block text-xs font-semibold text-gray-700 mb-2">Cari ID Donasi / Donatur</label>
        <div class="relative">
          <input v-model="searchQuery" type="text" placeholder="Masukkan kata kunci..." class="w-full pl-9 pr-3 py-2 text-xs border border-gray-200 rounded-lg focus:outline-none focus:border-green-500 bg-white text-slate-900">
          <i class="fas fa-search absolute left-3 top-3 text-gray-400 text-xs"></i>
        </div>
      </div>
    </div>

    <!-- Table List -->
    <div class="bg-white border border-gray-100 rounded-2xl shadow-sm overflow-hidden">
      <div class="p-6 border-b border-gray-100 flex items-center justify-between">
        <h3 class="font-bold text-slate-800 text-sm">Semua Transaksi Masuk</h3>
        <span class="text-xs text-gray-400 font-semibold">{{ donations.length }} Transaksi ditemukan</span>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs border-collapse">
          <thead class="bg-slate-50 border-b border-gray-100 text-gray-500 font-semibold uppercase">
            <tr>
              <th class="px-6 py-3">ID Transaksi</th>
              <th class="px-6 py-3">Donatur</th>
              <th class="px-6 py-3">Nominal</th>
              <th class="px-6 py-3">Metode</th>
              <th class="px-6 py-3">Status</th>
              <th class="px-6 py-3">Tanggal</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50">
            <tr v-for="item in donations" :key="item.id" class="hover:bg-slate-50 transition-colors">
              <td class="px-6 py-4 font-bold text-slate-700">#{{ item.id?.slice(0, 8) }}</td>
              <td class="px-6 py-4 font-semibold text-slate-700">
                {{ item.anonymous ? 'Hamba Allah (Anonim)' : (item.guestName || 'User YPAI') }}
              </td>
              <td class="px-6 py-4 font-extrabold text-slate-800">Rp {{ item.amount?.toLocaleString('id-ID') }}</td>
              <td class="px-6 py-4 text-gray-500 uppercase">{{ item.paymentMethod }}</td>
              <td class="px-6 py-4">
                <span :class="[
                  'px-2 py-1 rounded-full text-[10px] font-bold',
                  item.status === 'paid' ? 'bg-green-50 text-green-700' : 'bg-yellow-50 text-yellow-700'
                ]">
                  {{ item.status === 'paid' ? 'Lunas' : 'Pending' }}
                </span>
              </td>
              <td class="px-6 py-4 text-gray-400">
                {{ new Date(item.createdAt).toLocaleDateString('id-ID', { year: 'numeric', month: 'short', day: 'numeric' }) }}
              </td>
            </tr>
            <tr v-if="donations.length === 0">
              <td colspan="6" class="text-center py-12 text-gray-400">Tidak ada riwayat donasi masuk.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
