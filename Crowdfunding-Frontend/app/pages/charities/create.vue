<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// Form Data States
const title = ref('')
const description = ref('')
const categoryId = ref('')
const targetAmount = ref(10000000)
const startDate = ref('')
const endDate = ref('')
const contactPerson = ref('')
const contactEmail = ref('')
const contactPhone = ref('')
const organizationName = ref('')

// File inputs
const coverImageFile = ref(null)
const coverImagePreview = ref('')

// Categories
const categories = ref([])

// Feedback states
const isLoading = ref(false)
const errorMessage = ref('')
const isSuccess = ref(false)
const createdSlug = ref('')

onMounted(async () => {
  // Pastikan user sudah login
  const token = localStorage.getItem('auth_token')
  if (!token) {
    alert('Anda harus login terlebih dahulu untuk membuat penggalangan dana.')
    router.push('/login')
    return
  }

  // Ambil data kategori dari backend
  try {
    const response = await useApiFetch('/categories')
    if (response?.data?.value) {
      categories.value = response.data.value.data
    }
  } catch (err) {
    console.error('Gagal memuat kategori', err)
  }

  // Set default start date ke hari ini
  const today = new Date().toISOString().split('T')[0]
  startDate.value = today
})

const handleFileChange = (event) => {
  const file = event.target.files[0]
  if (file) {
    coverImageFile.value = file
    coverImagePreview.value = URL.createObjectURL(file)
  }
}

const handleSubmit = async () => {
  isLoading.value = true
  errorMessage.value = ''
  
  try {
    const token = localStorage.getItem('auth_token')
    if (!token) {
      router.push('/login')
      return
    }

    // Menggunakan FormData untuk multipart file upload
    const formData = new FormData()
    formData.append('title', title.value)
    formData.append('description', description.value)
    formData.append('categoryId', categoryId.value)
    formData.append('targetAmount', targetAmount.value)
    formData.append('startDate', startDate.value)
    formData.append('endDate', endDate.value)
    formData.append('contactPerson', contactPerson.value)
    formData.append('contactEmail', contactEmail.value)
    formData.append('contactPhone', contactPhone.value)
    formData.append('organization', organizationName.value)
    
    if (coverImageFile.value) {
      formData.append('coverImage', coverImageFile.value)
    }

    const config = useRuntimeConfig()
    const response = await $fetch(`${config.public.apiBase}/charities`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`
      },
      body: formData
    })

    if (response?.data) {
      createdSlug.value = response.data.slug
      isSuccess.value = true
    }
  } catch (error) {
    errorMessage.value = error.data?.error?.message || 'Gagal membuat campaign. Harap periksa data form Anda.'
  } finally {
    isLoading.value = false
  }
}

const viewCampaign = () => {
  router.push(`/charities/${createdSlug.value}`)
}
</script>

<template>
  <div class="bg-slate-50 min-h-screen py-10">
    <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
      
      <!-- Header -->
      <div class="text-center mb-8">
        <h1 class="text-3xl font-extrabold text-slate-900 tracking-tight">Buat Penggalangan Dana</h1>
        <p class="text-gray-600 mt-2">Mulai campaign penggalangan dana sosial baru untuk membantu mereka yang membutuhkan</p>
      </div>

      <!-- Form Container -->
      <div class="bg-white rounded-2xl shadow-sm border border-gray-100 p-8">
        
        <div v-if="errorMessage" class="mb-6 p-4 bg-red-50 text-red-600 rounded-xl text-sm flex items-center">
          <i class="fas fa-exclamation-circle mr-2"></i>
          {{ errorMessage }}
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-6">
          
          <!-- Basic Information -->
          <div class="space-y-4">
            <h2 class="text-lg font-bold text-slate-900 border-b border-gray-100 pb-2">Informasi Dasar</h2>
            
            <div>
              <label class="block text-sm font-semibold text-gray-700 mb-1">Judul Campaign *</label>
              <input v-model="title" type="text" required class="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:border-green-600" placeholder="Tuliskan judul penggalangan dana...">
            </div>

            <div>
              <label class="block text-sm font-semibold text-gray-700 mb-1">Deskripsi Lengkap *</label>
              <textarea v-model="description" rows="5" required class="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:border-green-600" placeholder="Tuliskan cerita dan rincian mengenai campaign sosial Anda..."></textarea>
            </div>

            <div>
              <label class="block text-sm font-semibold text-gray-700 mb-1">Kategori Campaign *</label>
              <select v-model="categoryId" required class="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:border-green-600 text-sm">
                <option value="">Pilih Kategori</option>
                <option v-for="category in categories" :key="category.id" :value="category.id">{{ category.name }}</option>
              </select>
            </div>
          </div>

          <!-- Financial Target -->
          <div class="space-y-4 pt-4">
            <h2 class="text-lg font-bold text-slate-900 border-b border-gray-100 pb-2">Target Finansial</h2>
            
            <div>
              <label class="block text-sm font-semibold text-gray-700 mb-1">Target Dana (IDR) *</label>
              <div class="relative">
                <input v-model="targetAmount" type="number" min="100000" required class="w-full py-2.5 pl-10 pr-4 border border-gray-200 rounded-xl focus:outline-none focus:border-green-600" placeholder="10000000">
                <span class="absolute left-3.5 top-1/2 transform -translate-y-1/2 text-sm text-gray-400 font-semibold">Rp</span>
              </div>
            </div>
          </div>

          <!-- Time Settings -->
          <div class="space-y-4 pt-4">
            <h2 class="text-lg font-bold text-slate-900 border-b border-gray-100 pb-2">Batas Waktu</h2>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-semibold text-gray-700 mb-1">Tanggal Mulai *</label>
                <input v-model="startDate" type="date" required class="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:border-green-600">
              </div>
              <div>
                <label class="block text-sm font-semibold text-gray-700 mb-1">Tanggal Selesai *</label>
                <input v-model="endDate" type="date" required class="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:border-green-600">
              </div>
            </div>
          </div>

          <!-- Cover Image Upload -->
          <div class="space-y-4 pt-4">
            <h2 class="text-lg font-bold text-slate-900 border-b border-gray-100 pb-2">Media & Foto Sampul</h2>
            
            <div>
              <label class="block text-sm font-semibold text-gray-700 mb-2">Foto Utama (Cover) *</label>
              <div class="border-2 border-dashed border-gray-200 rounded-2xl p-6 text-center hover:border-green-500 transition-colors">
                <input type="file" id="cover" accept="image/*" required class="hidden" @change="handleFileChange">
                <label for="cover" class="cursor-pointer">
                  <div v-if="coverImagePreview" class="mb-4">
                    <img :src="coverImagePreview" class="max-h-48 mx-auto rounded-xl shadow-sm">
                  </div>
                  <div v-else>
                    <i class="fas fa-cloud-upload-alt text-4xl text-gray-300 mb-2"></i>
                    <p class="text-sm font-semibold text-gray-600">Klik untuk mengunggah foto sampul</p>
                    <p class="text-xs text-gray-400 mt-1">Format PNG, JPG, GIF hingga 10MB</p>
                  </div>
                </label>
              </div>
            </div>
          </div>

          <!-- Contact Information -->
          <div class="space-y-4 pt-4">
            <h2 class="text-lg font-bold text-slate-900 border-b border-gray-100 pb-2">Kontak Narahubung</h2>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-semibold text-gray-700 mb-1">Nama Narahubung *</label>
                <input v-model="contactPerson" type="text" required class="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:border-green-600" placeholder="Nama lengkap">
              </div>
              <div>
                <label class="block text-sm font-semibold text-gray-700 mb-1">Email Kontak *</label>
                <input v-model="contactEmail" type="email" required class="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:border-green-600" placeholder="email@domain.com">
              </div>
              <div>
                <label class="block text-sm font-semibold text-gray-700 mb-1">No. Handphone Kontak</label>
                <input v-model="contactPhone" type="tel" class="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:border-green-600" placeholder="08123456xxx">
              </div>
              <div>
                <label class="block text-sm font-semibold text-gray-700 mb-1">Nama Organisasi (Opsional)</label>
                <input v-model="organizationName" type="text" class="w-full px-4 py-2.5 border border-gray-200 rounded-xl focus:outline-none focus:border-green-600" placeholder="Nama Yayasan / Lembaga">
              </div>
            </div>
          </div>

          <!-- Consent -->
          <div class="space-y-2 pt-2 text-xs text-gray-600">
            <div class="flex items-start">
              <input type="checkbox" id="terms" required class="mt-0.5 h-4 w-4 text-green-600 border-gray-300 rounded focus:ring-green-500">
              <label for="terms" class="ml-2">Saya menyetujui Ketentuan Layanan dan Kebijakan Privasi Yayasan Peduli Amal Indonesia.</label>
            </div>
            <div class="flex items-start">
              <input type="checkbox" id="verify" required class="mt-0.5 h-4 w-4 text-green-600 border-gray-300 rounded focus:ring-green-500">
              <label for="verify" class="ml-2">Saya menjamin bahwa seluruh informasi yang disediakan adalah akurat, sah, dan saya memiliki wewenang penuh atas aksi penggalangan dana ini.</label>
            </div>
          </div>

          <!-- Submit Buttons -->
          <button type="submit" :disabled="isLoading" class="w-full bg-green-600 text-white py-3.5 rounded-xl font-bold hover:bg-green-700 transform hover:scale-[1.01] transition-all shadow-md disabled:opacity-50 text-sm">
            <span v-if="isLoading"><i class="fas fa-spinner fa-spin mr-2"></i>Memproses Pembuatan...</span>
            <span v-else><i class="fas fa-rocket mr-1.5"></i>Buat Penggalangan Dana</span>
          </button>
        </form>
      </div>
    </div>

    <!-- Success Modal -->
    <div v-if="isSuccess" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-2xl p-8 max-w-sm w-full text-center shadow-2xl">
        <div class="w-16 h-16 bg-green-50 text-green-600 rounded-full flex items-center justify-center mx-auto mb-4 text-2xl">
          <i class="fas fa-check"></i>
        </div>
        <h3 class="text-xl font-bold text-slate-900 mb-2">Campaign Dibuat!</h3>
        <p class="text-gray-500 text-xs mb-6">Program penggalangan dana Anda berhasil dibuat dan telah tersimpan dengan status pending untuk ditinjau oleh tim admin.</p>
        <button @click="viewCampaign" class="w-full bg-green-600 text-white py-2.5 rounded-xl text-sm font-semibold hover:bg-green-700 transition-all">
          Lihat Halaman Donasi
        </button>
      </div>
    </div>

  </div>
</template>
