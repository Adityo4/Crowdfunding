<script setup>
import { ref, computed } from 'vue'

definePageMeta({
  layout: 'admin'
})

// Search, status and categories filters
const searchQuery = ref('')
const selectedCategory = ref('')
const selectedStatus = ref('')

// Fetch campaigns from backend
const { data: charitiesResponse, refresh } = await useApiFetch('/charities', {
  server: false,
  params: {
    status: 'all' // admin should be able to fetch drafts and paused campaigns too
  }
})

const charities = computed(() => {
  let list = charitiesResponse.value?.data || []
  
  if (searchQuery.value) {
    list = list.filter(item => item.title.toLowerCase().includes(searchQuery.value.toLowerCase()))
  }
  if (selectedCategory.value) {
    list = list.filter(item => item.categoryId === selectedCategory.value)
  }
  if (selectedStatus.value) {
    list = list.filter(item => item.status === selectedStatus.value)
  }
  
  return list
})

const verifyCampaign = async (id) => {
  if (!confirm('Apakah Anda yakin ingin menyetujui dan mengaktifkan penggalangan dana ini?')) return

  try {
    const config = useRuntimeConfig()
    const token = localStorage.getItem('auth_token')
    
    await $fetch(`${config.public.apiBase}/charities/${id}/status`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`
      },
      body: {
        status: 'active'
      }
    })
    
    alert('Sukses menyetujui kampanye!')
    refresh()
  } catch (err) {
    alert('Gagal mengaktifkan kampanye: ' + (err.data?.error?.message || err.message))
  }
}

const toggleCampaignStatus = async (id, currentStatus) => {
  const nextStatus = currentStatus === 'active' ? 'paused' : 'active'
  const actionText = nextStatus === 'active' ? 'mengaktifkan kembali' : 'menonaktifkan (pause)'
  
  if (!confirm(`Apakah Anda yakin ingin ${actionText} penggalangan dana ini?`)) return

  try {
    const config = useRuntimeConfig()
    const token = localStorage.getItem('auth_token')
    
    await $fetch(`${config.public.apiBase}/charities/${id}/status`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`
      },
      body: {
        status: nextStatus
      }
    })
    
    alert(`Sukses mengubah status kampanye ke ${nextStatus}!`)
    refresh()
  } catch (err) {
    alert('Gagal memperbarui status kampanye: ' + (err.data?.error?.message || err.message))
  }
}

// Edit Modal State
const showEditModal = ref(false)
const editCampaignId = ref('')
const editTitle = ref('')
const editDescription = ref('')
const editTargetAmount = ref(0)
const editStatus = ref('')
const editCategoryId = ref('')
const editStartDate = ref('')
const editEndDate = ref('')
const editContactPerson = ref('')
const editContactEmail = ref('')
const editContactPhone = ref('')
const editOrganizationName = ref('')

// Categories
const { data: categoriesResponse } = await useApiFetch('/categories')
const categories = computed(() => categoriesResponse.value?.data || [])

const openEditModal = (item) => {
  editCampaignId.value = item.id
  editTitle.value = item.title
  editDescription.value = item.description
  editTargetAmount.value = item.targetAmount
  editStatus.value = item.status
  editCategoryId.value = item.categoryId
  editStartDate.value = item.startDate?.split('T')[0] || ''
  editEndDate.value = item.endDate?.split('T')[0] || ''
  editContactPerson.value = item.contactPerson || ''
  editContactEmail.value = item.contactEmail || ''
  editContactPhone.value = item.contactPhone || ''
  editOrganizationName.value = item.organizationName || ''
  showEditModal.value = true
}

const handleUpdateCampaign = async () => {
  try {
    const config = useRuntimeConfig()
    const token = localStorage.getItem('auth_token')
    
    // 1. Update status
    await $fetch(`${config.public.apiBase}/charities/${editCampaignId.value}/status`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`
      },
      body: {
        status: editStatus.value
      }
    })

    // 2. Update campaign payload details
    await $fetch(`${config.public.apiBase}/charities/${editCampaignId.value}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: {
        title: editTitle.value,
        description: editDescription.value,
        categoryId: editCategoryId.value,
        targetAmount: Number(editTargetAmount.value),
        startDate: editStartDate.value + 'T00:00:00Z',
        endDate: editEndDate.value + 'T00:00:00Z',
        contactPerson: editContactPerson.value,
        contactEmail: editContactEmail.value,
        contactPhone: editContactPhone.value,
        organization: editOrganizationName.value
      }
    })
    
    alert('Sukses mengupdate kampanye!')
    showEditModal.value = false
    refresh()
  } catch (err) {
    alert('Gagal memperbarui data: ' + (err.data?.error?.message || err.message))
  }
}

</script>

<template>
  <div class="space-y-6">
    <!-- Header Page Actions -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h2 class="text-xl font-bold text-slate-800">Manajemen Penggalangan Dana</h2>
        <p class="text-xs text-gray-500">Kelola, verifikasi, dan pantau status donasi kampanye sosial YPAI.</p>
      </div>
      <NuxtLink to="/charities/create" class="inline-flex items-center justify-center bg-green-600 hover:bg-green-700 text-white font-bold text-xs px-4 py-2.5 rounded-lg transition-colors">
        <i class="fas fa-plus mr-2"></i>Tambah Penggalangan Dana
      </NuxtLink>
    </div>

    <!-- Filters Section -->
    <div class="bg-white border border-gray-100 rounded-2xl p-6 shadow-sm">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
          <label class="block text-xs font-semibold text-gray-700 mb-2">Cari Kampanye</label>
          <div class="relative">
            <input v-model="searchQuery" type="text" placeholder="Masukkan judul..." class="w-full pl-9 pr-3 py-2 text-xs border border-gray-200 rounded-lg focus:outline-none focus:border-green-500 bg-white text-slate-900">
            <i class="fas fa-search absolute left-3 top-3 text-gray-400 text-xs"></i>
          </div>
        </div>
        <div>
          <label class="block text-xs font-semibold text-gray-700 mb-2">Status Kampanye</label>
          <select v-model="selectedStatus" class="w-full px-3 py-2 text-xs border border-gray-200 rounded-lg focus:outline-none focus:border-green-500 bg-white text-slate-900">
            <option value="">Semua Status</option>
            <option value="active">Aktif</option>
            <option value="completed">Selesai</option>
            <option value="paused">Ditangguhkan</option>
          </select>
        </div>
        <div>
          <label class="block text-xs font-semibold text-gray-700 mb-2">Urutan</label>
          <select class="w-full px-3 py-2 text-xs border border-gray-200 rounded-lg focus:outline-none focus:border-green-500 bg-white text-slate-900">
            <option value="newest">Terbaru</option>
            <option value="oldest">Terlama</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Table List -->
    <div class="bg-white border border-gray-100 rounded-2xl shadow-sm overflow-hidden">
      <div class="p-6 border-b border-gray-100 flex items-center justify-between">
        <h3 class="font-bold text-slate-800 text-sm">Daftar Kampanye Sosial</h3>
        <span class="text-xs text-gray-400 font-semibold">{{ charities.length }} Program ditemukan</span>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs border-collapse">
          <thead class="bg-slate-50 border-b border-gray-100 text-gray-500 font-semibold uppercase">
            <tr>
              <th class="px-6 py-3">Kampanye</th>
              <th class="px-6 py-3">Kategori</th>
              <th class="px-6 py-3">Progress Dana</th>
              <th class="px-6 py-3">Status</th>
              <th class="px-6 py-3 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50">
            <tr v-for="item in charities" :key="item.id" class="hover:bg-slate-50 transition-colors">
              <td class="px-6 py-4">
                <div class="flex items-center space-x-3">
                  <img :src="item.coverImageUrl || 'https://images.unsplash.com/photo-1559027615-cd4628902d85?auto=format&fit=crop&w=100&q=80'" class="w-10 h-10 object-cover rounded-lg">
                  <div>
                    <h4 class="font-bold text-slate-800">{{ item.title }}</h4>
                    <span class="text-[10px] text-gray-400">ID: {{ item.id }}</span>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4">
                <span class="px-2.5 py-1 bg-green-50 text-green-700 font-bold rounded-full text-[10px]">
                  {{ item.category?.name || 'Sosial' }}
                </span>
              </td>
              <td class="px-6 py-4">
                <div class="space-y-1 max-w-[150px]">
                  <div class="flex justify-between text-[10px] text-gray-500 font-medium">
                    <span>Rp {{ item.collectedAmount?.toLocaleString('id-ID') }}</span>
                    <span>Rp {{ item.targetAmount?.toLocaleString('id-ID') }}</span>
                  </div>
                  <div class="w-full bg-slate-100 h-1.5 rounded-full overflow-hidden">
                    <div class="bg-green-600 h-full rounded-full" :style="{ width: `${Math.min((item.collectedAmount / item.targetAmount) * 100, 100)}%` }"></div>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4">
                <span :class="[
                  'px-2 py-1 rounded-full text-[10px] font-bold',
                  item.status === 'active' ? 'bg-green-50 text-green-700' : (item.status === 'pending' ? 'bg-yellow-50 text-yellow-700' : 'bg-red-50 text-red-700')
                ]">
                  {{ item.status === 'active' ? 'Aktif' : (item.status === 'pending' ? 'Pending' : 'Nonaktif') }}
                </span>
              </td>
              <td class="px-6 py-4 text-right">
                <div class="inline-flex space-x-2 items-center">
                  <button v-if="item.status === 'pending'" @click="verifyCampaign(item.id)" class="px-2 py-1 bg-green-600 hover:bg-green-700 text-white font-bold text-[10px] rounded-lg transition-colors">
                    <i class="fas fa-check"></i> Aktifkan
                  </button>
                  <button v-else-if="item.status === 'active'" @click="toggleCampaignStatus(item.id, 'active')" class="px-2 py-1 bg-yellow-500 hover:bg-yellow-600 text-white font-bold text-[10px] rounded-lg transition-colors">
                    <i class="fas fa-pause"></i> Pause
                  </button>
                  <button v-else-if="item.status === 'paused'" @click="toggleCampaignStatus(item.id, 'paused')" class="px-2 py-1 bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-[10px] rounded-lg transition-colors">
                    <i class="fas fa-play"></i> Play
                  </button>
                  
                  <button @click="openEditModal(item)" class="p-1.5 text-blue-500 hover:text-blue-700 transition-colors">
                    <i class="fas fa-edit text-sm"></i>
                  </button>
                  <NuxtLink :to="`/charities/${item.slug}`" class="p-1.5 text-gray-400 hover:text-green-600 transition-colors">
                    <i class="fas fa-eye text-sm"></i>
                  </NuxtLink>
                </div>
              </td>
            </tr>
            <tr v-if="charities.length === 0">
              <td colspan="5" class="text-center py-12 text-gray-400">Tidak ada data kampanye sosial tersedia.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Edit Campaign Modal -->
    <div v-if="showEditModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 p-4">
      <div class="bg-white rounded-2xl max-w-lg w-full p-6 shadow-2xl border border-gray-100 space-y-4 max-h-[90vh] overflow-y-auto">
        <div class="flex justify-between items-center pb-2 border-b border-gray-100">
          <h3 class="font-bold text-slate-800 text-sm">Edit Data Kampanye</h3>
          <button @click="showEditModal = false" class="text-gray-400 hover:text-gray-600">
            <i class="fas fa-times"></i>
          </button>
        </div>
        <form @submit.prevent="handleUpdateCampaign" class="space-y-4 text-left">
          <div class="grid grid-cols-2 gap-4">
            <div class="col-span-2">
              <label class="block text-xs font-semibold text-gray-700 mb-1">Judul Kampanye *</label>
              <input v-model="editTitle" type="text" required class="w-full px-3 py-2 text-xs border border-gray-200 rounded-lg bg-white text-slate-900 focus:border-green-500 focus:outline-none">
            </div>
            
            <div class="col-span-2">
              <label class="block text-xs font-semibold text-gray-700 mb-1">Deskripsi Singkat *</label>
              <textarea v-model="editDescription" required rows="3" class="w-full px-3 py-2 text-xs border border-gray-200 rounded-lg bg-white text-slate-900 focus:border-green-500 focus:outline-none"></textarea>
            </div>

            <div>
              <label class="block text-xs font-semibold text-gray-700 mb-1">Target Dana (Rp) *</label>
              <input v-model="editTargetAmount" type="number" required class="w-full px-3 py-2 text-xs border border-gray-200 rounded-lg bg-white text-slate-900 focus:border-green-500 focus:outline-none">
            </div>

            <div>
              <label class="block text-xs font-semibold text-gray-700 mb-1">Kategori *</label>
              <select v-model="editCategoryId" required class="w-full px-3 py-2 text-xs border border-gray-200 rounded-lg bg-white text-slate-900 focus:border-green-500 focus:outline-none">
                <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-semibold text-gray-700 mb-1">Tanggal Mulai *</label>
              <input v-model="editStartDate" type="date" required class="w-full px-3 py-2 text-xs border border-gray-200 rounded-lg bg-white text-slate-900 focus:border-green-500 focus:outline-none">
            </div>

            <div>
              <label class="block text-xs font-semibold text-gray-700 mb-1">Tanggal Selesai *</label>
              <input v-model="editEndDate" type="date" required class="w-full px-3 py-2 text-xs border border-gray-200 rounded-lg bg-white text-slate-900 focus:border-green-500 focus:outline-none">
            </div>

            <div>
              <label class="block text-xs font-semibold text-gray-700 mb-1">Kontak Person *</label>
              <input v-model="editContactPerson" type="text" required class="w-full px-3 py-2 text-xs border border-gray-200 rounded-lg bg-white text-slate-900 focus:border-green-500 focus:outline-none">
            </div>

            <div>
              <label class="block text-xs font-semibold text-gray-700 mb-1">Email Kontak *</label>
              <input v-model="editContactEmail" type="email" required class="w-full px-3 py-2 text-xs border border-gray-200 rounded-lg bg-white text-slate-900 focus:border-green-500 focus:outline-none">
            </div>

            <div>
              <label class="block text-xs font-semibold text-gray-700 mb-1">No. HP Kontak</label>
              <input v-model="editContactPhone" type="text" class="w-full px-3 py-2 text-xs border border-gray-200 rounded-lg bg-white text-slate-900 focus:border-green-500 focus:outline-none">
            </div>

            <div>
              <label class="block text-xs font-semibold text-gray-700 mb-1">Nama Organisasi</label>
              <input v-model="editOrganizationName" type="text" class="w-full px-3 py-2 text-xs border border-gray-200 rounded-lg bg-white text-slate-900 focus:border-green-500 focus:outline-none">
            </div>

            <div class="col-span-2">
              <label class="block text-xs font-semibold text-gray-700 mb-1">Status Kampanye</label>
              <select v-model="editStatus" class="w-full px-3 py-2 text-xs border border-gray-200 rounded-lg focus:outline-none focus:border-green-500 bg-white text-slate-900">
                <option value="pending">Pending (Verifikasi)</option>
                <option value="active">Aktif (Menerima Donasi)</option>
                <option value="paused">Ditangguhkan (Pause)</option>
              </select>
            </div>
          </div>
          <div class="flex justify-end space-x-2 pt-2 border-t border-gray-100">
            <button type="button" @click="showEditModal = false" class="px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg font-bold text-xs">Batal</button>
            <button type="submit" class="px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg font-bold text-xs">Simpan Perubahan</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
