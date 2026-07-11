<script setup>
import { ref, computed } from 'vue'

definePageMeta({
  layout: 'admin'
})

// Search query
const searchQuery = ref('')

// Fetch articles dynamically (for admin articles overview)
const { data: articlesResponse, refresh } = await useApiFetch('/articles?all=true', { server: false })
const articles = computed(() => {
  let list = articlesResponse.value?.data || []
  if (searchQuery.value) {
    list = list.filter(item => 
      item.title?.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }
  return list
})

// Create article state
const isCreateModalOpen = ref(false)
const articleTitle = ref('')
const articleContent = ref('')
const articleIsPublished = ref(true)
const isSubmitting = ref(false)
const articleImages = ref([])
const imagePreviews = ref([])

const handleFileChange = (e) => {
  const files = e.target.files
  if (!files) return
  for (let i = 0; i < files.length; i++) {
    const file = files[i]
    articleImages.value.push(file)
    imagePreviews.value.push(URL.createObjectURL(file))
  }
}

const removeImage = (index) => {
  articleImages.value.splice(index, 1)
  URL.revokeObjectURL(imagePreviews.value[index])
  imagePreviews.value.splice(index, 1)
}

const openCreateModal = () => {
  articleTitle.value = ''
  articleContent.value = ''
  articleIsPublished.value = true
  articleImages.value = []
  imagePreviews.value.forEach(url => URL.revokeObjectURL(url))
  imagePreviews.value = []
  isCreateModalOpen.value = true
}

const closeCreateModal = () => {
  isCreateModalOpen.value = false
}

const submitArticle = async () => {
  if (!articleTitle.value || !articleContent.value) {
    alert('Judul dan Konten Artikel wajib diisi!')
    return
  }

  isSubmitting.value = true
  
  const formData = new FormData()
  formData.append('title', articleTitle.value)
  formData.append('content', articleContent.value)
  formData.append('isPublished', articleIsPublished.value ? 'true' : 'false')
  
  for (let i = 0; i < articleImages.value.length; i++) {
    formData.append('images', articleImages.value[i])
  }

  try {
    const { error } = await useApiFetch('/articles', {
      method: 'POST',
      body: formData
    })

    if (error.value) {
      alert('Gagal membuat artikel: ' + (error.value.data?.error?.message || error.value.message))
    } else {
      alert('Artikel berhasil dibuat!')
      closeCreateModal()
      refresh()
    }
  } catch (e) {
    console.error(e)
    alert('Terjadi kesalahan saat membuat artikel.')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h2 class="text-xl font-bold text-slate-800">Manajemen Artikel</h2>
        <p class="text-xs text-gray-500">Kelola publikasi kabar, rilis berita, dan cerita sukses aksi sosial YPAI.</p>
      </div>
      <button @click="openCreateModal" class="inline-flex items-center justify-center bg-green-600 hover:bg-green-700 text-white font-bold text-xs px-4 py-2.5 rounded-lg transition-colors">
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
              <th class="px-6 py-3">Status</th>
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
              <td class="px-6 py-4">
                <span :class="[
                  'px-2.5 py-1 rounded-full text-[10px] font-bold inline-block',
                  item.isPublished ? 'bg-green-50 text-green-700' : 'bg-yellow-50 text-yellow-700'
                ]">
                  {{ item.isPublished ? 'Diterbitkan' : 'Draft' }}
                </span>
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
              <td colspan="4" class="text-center py-12 text-gray-400">Tidak ada artikel tersedia.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create Article Modal -->
    <div v-if="isCreateModalOpen" class="fixed inset-0 z-50 overflow-y-auto bg-slate-900/60 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="bg-white rounded-2xl max-w-2xl w-full shadow-xl border border-gray-100 flex flex-col max-h-[90vh]">
        <!-- Modal Header -->
        <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
          <div>
            <h3 class="font-bold text-slate-800 text-base">Tulis Artikel Baru</h3>
            <p class="text-[11px] text-gray-500">Buat publikasi cerita sukses, kabar program, atau berita terbaru.</p>
          </div>
          <button @click="closeCreateModal" class="text-gray-400 hover:text-gray-600 p-1.5 rounded-lg hover:bg-gray-50 transition-colors">
            <i class="fas fa-times text-sm"></i>
          </button>
        </div>

        <!-- Modal Body -->
        <div class="p-6 overflow-y-auto space-y-4 flex-grow">
          <div>
            <label class="block text-xs font-semibold text-gray-700 mb-1">Judul Artikel <span class="text-red-500">*</span></label>
            <input v-model="articleTitle" type="text" placeholder="Contoh: Kisah Sukses Pembangunan Masjid di Pelosok" class="w-full px-3 py-2 text-xs border border-gray-200 rounded-lg focus:outline-none focus:border-green-500 bg-white text-slate-900">
          </div>

          <div>
            <label class="block text-xs font-semibold text-gray-700 mb-1">Unggah Gambar Artikel <span class="text-[10px] text-gray-400 font-normal">(Bisa memilih lebih dari 1 file)</span></label>
            <div class="mt-1 flex flex-col space-y-3">
              <div class="flex items-center justify-center w-full">
                <label class="flex flex-col items-center justify-center w-full h-32 border-2 border-gray-300 border-dashed rounded-xl cursor-pointer bg-slate-50 hover:bg-slate-100 transition-colors">
                  <div class="flex flex-col items-center justify-center pt-5 pb-6">
                    <i class="fas fa-cloud-upload-alt text-2xl text-gray-400 mb-2"></i>
                    <p class="text-xs text-gray-500 font-semibold">Klik untuk mengunggah gambar</p>
                    <p class="text-[10px] text-gray-400">PNG, JPG atau JPEG</p>
                  </div>
                  <input type="file" multiple accept="image/*" class="hidden" @change="handleFileChange">
                </label>
              </div>

              <!-- Image Previews -->
              <div v-if="imagePreviews.length > 0" class="grid grid-cols-3 sm:grid-cols-4 gap-4 p-2 bg-slate-50 rounded-xl">
                <div v-for="(preview, index) in imagePreviews" :key="index" class="relative group aspect-square rounded-lg overflow-hidden border border-gray-200 shadow-sm bg-white">
                  <img :src="preview" class="w-full h-full object-cover">
                  <span v-if="index === 0" class="absolute top-1 left-1 bg-green-600 text-white text-[9px] font-bold px-1.5 py-0.5 rounded shadow-sm">Cover</span>
                  <button type="button" @click="removeImage(index)" class="absolute top-1 right-1 w-6 h-6 bg-red-600 hover:bg-red-700 text-white rounded-full flex items-center justify-center shadow-md transition-colors">
                    <i class="fas fa-trash-alt text-[10px]"></i>
                  </button>
                </div>
              </div>
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-gray-700 mb-1">Konten Artikel <span class="text-red-500">*</span></label>
            <textarea v-model="articleContent" rows="10" placeholder="Tuliskan isi berita atau artikel lengkap di sini..." class="w-full px-3 py-2 text-xs border border-gray-200 rounded-lg focus:outline-none focus:border-green-500 bg-white text-slate-900 resize-y"></textarea>
          </div>

          <div class="flex items-center space-x-3 bg-slate-50 p-3 rounded-xl">
            <input v-model="articleIsPublished" type="checkbox" id="isPublished" class="w-4 h-4 text-green-600 border-gray-300 rounded focus:ring-green-500">
            <label for="isPublished" class="text-xs font-semibold text-gray-700 cursor-pointer">
              Publikasikan Langsung
              <span class="block text-[10px] text-gray-400 font-normal">Jika dicentang, artikel akan langsung muncul di halaman publik.</span>
            </label>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="px-6 py-4 border-t border-gray-100 flex items-center justify-end space-x-3 bg-slate-50 rounded-b-2xl">
          <button @click="closeCreateModal" class="px-4 py-2 rounded-lg text-xs font-bold text-gray-500 hover:text-gray-700 hover:bg-gray-100 transition-colors">
            Batal
          </button>
          <button @click="submitArticle" :disabled="isSubmitting" class="px-4 py-2 rounded-lg bg-green-600 hover:bg-green-700 text-white text-xs font-bold transition-colors disabled:opacity-50 flex items-center">
            <i v-if="isSubmitting" class="fas fa-spinner fa-spin mr-2"></i>
            {{ isSubmitting ? 'Menyimpan...' : 'Simpan Artikel' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
