<script setup>
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const slug = route.params.slug

// Fetch charity detail by slug
const { data: charityResponse, refresh } = await useApiFetch(`/charities/${slug}`)
const charity = computed(() => charityResponse.value?.data || {})

// Donation Form State
const donationAmount = ref(100000)
const paymentMethod = ref('gopay')
const anonymous = ref(false)
const message = ref('')

// Donation Result State
const isLoading = ref(false)
const isSuccess = ref(false)
const donationResult = ref(null)
const errorMessage = ref('')

const selectAmount = (amount) => {
  donationAmount.value = amount
}

const handleDonate = async () => {
  isLoading.value = true
  errorMessage.value = ''
  isSuccess.value = false
  
  try {
    const token = localStorage.getItem('auth_token')
    const headers = {}
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }
    
    const config = useRuntimeConfig()
    const response = await $fetch(`${config.public.apiBase}/donations`, {
      method: 'POST',
      headers,
      body: {
        charityId: charity.value.id,
        amount: Number(donationAmount.value),
        paymentMethod: paymentMethod.value,
        anonymous: anonymous.value,
        message: message.value
      }
    })
    
    if (response?.data) {
      donationResult.value = response.data
      isSuccess.value = true
    }
  } catch (error) {
    errorMessage.value = error.data?.error?.message || 'Gagal memproses donasi. Silakan coba lagi.'
  } finally {
    isLoading.value = false
  }
}

// Simulasi konfirmasi pembayaran
const handleConfirmPayment = async () => {
  if (!donationResult.value) return
  isLoading.value = true
  
  try {
    const config = useRuntimeConfig()
    const response = await $fetch(`${config.public.apiBase}/donations/callback`, {
      method: 'POST',
      body: {
        donationId: donationResult.value.donationId,
        paymentReference: 'PAY-SIM-' + donationResult.value.donationId.slice(0, 8),
        status: 'settlement' // simulate payment success
      }
    })
    
    if (response?.data?.status === 'paid') {
      alert('Pembayaran Terkonfirmasi! Terima kasih atas kebaikan Anda.')
      isSuccess.value = false
      donationResult.value = null
      message.value = ''
      refresh() // update collected amount
    }
  } catch (err) {
    alert('Konfirmasi pembayaran gagal.')
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div v-if="charity.id" class="bg-slate-50 min-h-screen py-10">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      
      <!-- Layout Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Main Info -->
        <div class="lg:col-span-2 space-y-6">
          <div class="bg-white rounded-2xl p-6 shadow-sm border border-gray-100">
            <div class="inline-flex items-center px-3 py-1 bg-green-50 rounded-full text-green-600 text-xs font-semibold mb-4">
              <i class="fas fa-heart mr-1.5"></i>
              <span>Charity Program</span>
            </div>
            <h1 class="text-2xl md:text-4xl font-extrabold text-slate-900 mb-2 leading-tight">
              {{ charity.title }}
            </h1>
            <p class="text-sm text-gray-500 mb-6">
              Diselenggarakan oleh <span class="font-semibold text-gray-700">{{ charity.organizationName || 'Pengguna Umum' }}</span> &bull; Kategori <span class="font-semibold text-green-600">{{ charity.category?.name }}</span>
            </p>
            
            <!-- Cover Image -->
            <div class="relative rounded-2xl overflow-hidden mb-6 h-64 md:h-96">
              <img :src="charity.coverImageUrl || 'https://images.unsplash.com/photo-1523050854058-8df90110c9f1?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80'" alt="Charity Cover" class="w-full h-full object-cover">
            </div>
            
            <h2 class="text-xl font-bold text-slate-900 mb-4">Tentang Program</h2>
            <div class="prose prose-sm max-w-none text-gray-600 leading-relaxed whitespace-pre-line">
              {{ charity.description }}
            </div>
          </div>
          
          <!-- Campaign Progress Summary -->
          <div class="bg-white rounded-2xl p-6 shadow-sm border border-gray-100">
            <h2 class="text-lg font-bold text-slate-900 mb-4">Perkembangan Penggalangan Dana</h2>
            <div class="flex justify-between items-end mb-2">
              <div>
                <span class="text-2xl font-black text-green-600">Rp {{ charity.currentAmount.toLocaleString('id-ID') }}</span>
                <span class="text-xs text-gray-500 block">Terkumpul dari target Rp {{ charity.targetAmount.toLocaleString('id-ID') }}</span>
              </div>
              <span class="text-lg font-bold text-slate-900">{{ ((charity.currentAmount / charity.targetAmount) * 100).toFixed(0) }}%</span>
            </div>
            <div class="w-full bg-gray-100 rounded-full h-2 overflow-hidden mb-6">
              <div class="bg-gradient-to-r from-green-500 to-emerald-500 h-2 rounded-full transition-all duration-500" 
                   :style="{ width: Math.min((charity.currentAmount / charity.targetAmount) * 100, 100) + '%' }"></div>
            </div>
            <div class="grid grid-cols-2 gap-4 text-center">
              <div class="p-3 bg-slate-50 rounded-xl">
                <span class="block text-lg font-bold text-slate-900">
                  {{ charity.status === 'active' ? 'Aktif' : charity.status }}
                </span>
                <span class="text-xs text-gray-500">Status</span>
              </div>
              <div class="p-3 bg-slate-50 rounded-xl">
                <span class="block text-lg font-bold text-slate-900">
                  {{ new Date(charity.endDate).toLocaleDateString('id-ID', { month: 'short', day: 'numeric', year: 'numeric' }) }}
                </span>
                <span class="text-xs text-gray-500">Batas Waktu</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Sidebar / Donation Widget -->
        <div class="lg:col-span-1 space-y-6">
          <div class="bg-white rounded-2xl shadow-sm border border-gray-100 p-6 sticky top-24">
            
            <!-- Success Modal / Invoice -->
            <div v-if="isSuccess && donationResult" class="space-y-4 text-center">
              <div class="w-12 h-12 bg-green-50 text-green-600 rounded-full flex items-center justify-center mx-auto text-xl">
                <i class="fas fa-file-invoice-dollar"></i>
              </div>
              <h3 class="text-lg font-bold text-slate-900">Invoice Pembayaran</h3>
              <p class="text-xs text-gray-500">Silakan scan QR Code atau bayar menggunakan tautan pembayaran di bawah ini.</p>
              
              <div class="p-4 bg-slate-50 rounded-xl border border-gray-100">
                <span class="text-xs text-gray-500 block">Total Pembayaran</span>
                <span class="text-xl font-black text-slate-900">Rp {{ donationResult.amount.toLocaleString('id-ID') }}</span>
              </div>
              
              <!-- QR Code Mock -->
              <div class="border border-dashed border-gray-200 p-3 bg-white rounded-xl inline-block">
                <i class="fas fa-qrcode text-8xl text-slate-900"></i>
              </div>
              
              <div class="space-y-2">
                <button @click="handleConfirmPayment" :disabled="isLoading" class="w-full bg-green-600 text-white py-2.5 rounded-xl text-sm font-semibold hover:bg-green-700 transition-all disabled:opacity-50">
                  <span v-if="isLoading">Menghubungkan...</span>
                  <span v-else>Simulasi Bayar Sukses</span>
                </button>
                <button @click="isSuccess = false" class="w-full text-xs text-gray-400 hover:text-gray-600 py-1">Kembali</button>
              </div>
            </div>

            <!-- Form Donasi -->
            <div v-else>
              <h3 class="text-lg font-bold text-slate-900 mb-4 border-b border-gray-100 pb-2">Kirim Donasi</h3>
              
              <div v-if="errorMessage" class="mb-4 p-3 bg-red-50 text-red-600 rounded-xl text-xs">
                {{ errorMessage }}
              </div>
              
              <form @submit.prevent="handleDonate" class="space-y-4">
                <!-- Preset Amounts -->
                <div class="grid grid-cols-2 gap-2">
                  <button type="button" @click="selectAmount(50000)" :class="{'bg-green-50 text-green-600 border-green-600': donationAmount === 50000, 'bg-slate-50 text-gray-700 border-transparent': donationAmount !== 50000}" class="py-2.5 border rounded-xl text-xs font-semibold hover:bg-green-50/50 transition-all">Rp 50K</button>
                  <button type="button" @click="selectAmount(100000)" :class="{'bg-green-50 text-green-600 border-green-600': donationAmount === 100000, 'bg-slate-50 text-gray-700 border-transparent': donationAmount !== 100000}" class="py-2.5 border rounded-xl text-xs font-semibold hover:bg-green-50/50 transition-all">Rp 100K</button>
                  <button type="button" @click="selectAmount(250000)" :class="{'bg-green-50 text-green-600 border-green-600': donationAmount === 250000, 'bg-slate-50 text-gray-700 border-transparent': donationAmount !== 250000}" class="py-2.5 border rounded-xl text-xs font-semibold hover:bg-green-50/50 transition-all">Rp 250K</button>
                  <button type="button" @click="selectAmount(500000)" :class="{'bg-green-50 text-green-600 border-green-600': donationAmount === 500000, 'bg-slate-50 text-gray-700 border-transparent': donationAmount !== 500000}" class="py-2.5 border rounded-xl text-xs font-semibold hover:bg-green-50/50 transition-all">Rp 500K</button>
                </div>
                
                <!-- Custom Amount -->
                <div class="space-y-1">
                  <label class="block text-xs font-semibold text-gray-600">Nominal Donasi Lainnya</label>
                  <div class="relative">
                    <input v-model="donationAmount" type="number" min="1000" class="w-full py-2.5 pl-10 pr-4 border border-gray-200 rounded-xl focus:outline-none focus:border-green-600" placeholder="0">
                    <span class="absolute left-3.5 top-1/2 transform -translate-y-1/2 text-sm text-gray-400 font-semibold">Rp</span>
                  </div>
                </div>

                <!-- Payment Method -->
                <div class="space-y-1">
                  <label class="block text-xs font-semibold text-gray-600">Metode Pembayaran</label>
                  <select v-model="paymentMethod" class="w-full px-3 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:border-green-600 text-sm">
                    <option value="gopay">GoPay</option>
                    <option value="ovo">OVO</option>
                    <option value="shopeepay">ShopeePay</option>
                    <option value="bank_transfer">Transfer Bank (Virtual Account)</option>
                  </select>
                </div>

                <!-- Message -->
                <div class="space-y-1">
                  <label class="block text-xs font-semibold text-gray-600">Pesan / Doa Donatur (Opsional)</label>
                  <textarea v-model="message" rows="3" class="w-full px-3 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:border-green-600 text-sm" placeholder="Tulis doa atau pesan dukungan..."></textarea>
                </div>

                <!-- Anonymous Checkbox -->
                <div class="flex items-center space-x-2">
                  <input v-model="anonymous" type="checkbox" id="anon" class="w-4 h-4 text-green-600 border-gray-300 rounded focus:ring-green-500">
                  <label for="anon" class="text-xs text-gray-600 cursor-pointer">Sembunyikan nama saya (Anonim)</label>
                </div>
                
                <button type="submit" :disabled="isLoading" class="w-full bg-green-600 hover:bg-green-700 text-white py-3 rounded-xl text-sm font-semibold transition-all shadow-sm">
                  <span v-if="isLoading"><i class="fas fa-spinner fa-spin mr-2"></i>Memproses...</span>
                  <span v-else>Kirim Donasi</span>
                </button>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div v-else class="text-center py-20 text-gray-500">
    <i class="fas fa-circle-notch fa-spin text-3xl mb-3 text-green-600"></i>
    <p>Memuat rincian penggalangan dana...</p>
  </div>
</template>
