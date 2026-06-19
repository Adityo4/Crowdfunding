<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '~/stores/auth'

const authStore = useAuthStore()
const router = useRouter()

const profileName = ref('')
const profilePhone = ref('')
const profilePassword = ref('')
const isUpdatingProfile = ref(false)

onMounted(() => {
  authStore.init()
  if (!authStore.user) {
    alert('Silakan login terlebih dahulu untuk mengakses profil Anda.')
    router.push('/login')
    return
  }

  profileName.value = authStore.user.name || ''
  profilePhone.value = authStore.user.phoneNumber || ''
})

const handleUpdateProfile = async () => {
  if (!profileName.value) {
    alert('Nama lengkap wajib diisi!')
    return
  }

  isUpdatingProfile.value = true
  try {
    const config = useRuntimeConfig()
    const res = await $fetch(`${config.public.apiBase}/users/profile`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      },
      body: {
        fullName: profileName.value,
        phoneNumber: profilePhone.value,
        password: profilePassword.value
      }
    })

    if (res?.data) {
      // Update auth store & local storage
      const updatedUser = {
        id: res.data.id || authStore.user.id,
        name: res.data.fullName,
        email: res.data.email,
        role: res.data.role || authStore.user.role,
        phoneNumber: res.data.phoneNumber || ''
      }
      authStore.user = updatedUser
      localStorage.setItem('user_data', JSON.stringify(updatedUser))

      alert('Profil Anda berhasil diperbarui!')
      profilePassword.value = ''
    }
  } catch (err) {
    console.error('Gagal memperbarui profil:', err)
    alert('Gagal memperbarui profil: ' + (err.data?.error?.message || err.message))
  } finally {
    isUpdatingProfile.value = false
  }
}
</script>

<template>
  <div class="bg-slate-50 min-h-screen pt-24 pb-12">
    <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Profile Header Card -->
      <div class="bg-white border border-slate-100 rounded-2xl p-6 md:p-8 shadow-sm flex flex-col md:flex-row items-center gap-6 mb-8">
        <div class="w-20 h-20 bg-gradient-to-r from-green-500 to-blue-500 rounded-full flex items-center justify-center text-white font-extrabold text-3xl shadow-md">
          {{ authStore.user?.name?.charAt(0) || 'U' }}
        </div>
        <div class="text-center md:text-left flex-grow">
          <h2 class="text-2xl font-bold text-slate-900">{{ authStore.user?.name }}</h2>
          <p class="text-sm text-gray-500 mt-1">{{ authStore.user?.email }}</p>
          <span class="inline-block mt-3 px-3 py-1 bg-green-50 text-green-700 text-xs font-bold rounded-full capitalize">
            Role: {{ authStore.user?.role }}
          </span>
        </div>
      </div>

      <!-- Settings / Edit Profile Form Card -->
      <div class="bg-white border border-slate-100 rounded-2xl shadow-sm p-6 md:p-8 space-y-6">
        <div class="pb-4 border-b border-gray-100">
          <h3 class="font-bold text-slate-800 text-lg">Pengaturan Profil Pengguna</h3>
          <p class="text-xs text-gray-500 mt-1">Kelola data diri, nomor telepon, dan tingkatkan keamanan akun Anda.</p>
        </div>

        <form @submit.prevent="handleUpdateProfile" class="space-y-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label class="block text-xs font-semibold text-gray-700 mb-2">Nama Lengkap</label>
              <div class="relative">
                <input v-model="profileName" type="text" required class="w-full pl-10 pr-3 py-2.5 text-xs border border-gray-200 rounded-lg focus:outline-none focus:border-green-500 bg-white text-slate-900" placeholder="Masukkan nama lengkap...">
                <i class="fas fa-user absolute left-3.5 top-3.5 text-gray-400 text-xs"></i>
              </div>
            </div>

            <div>
              <label class="block text-xs font-semibold text-gray-700 mb-2">Nomor Telepon</label>
              <div class="relative">
                <input v-model="profilePhone" type="text" class="w-full pl-10 pr-3 py-2.5 text-xs border border-gray-200 rounded-lg focus:outline-none focus:border-green-500 bg-white text-slate-900" placeholder="Contoh: 08123456789">
                <i class="fas fa-phone absolute left-3.5 top-3.5 text-gray-400 text-xs"></i>
              </div>
            </div>

            <div class="md:col-span-2">
              <label class="block text-xs font-semibold text-gray-700 mb-2">Alamat Email <span class="text-[10px] text-gray-400 font-normal">(Tidak dapat diubah)</span></label>
              <div class="relative">
                <input :value="authStore.user?.email" type="email" disabled class="w-full pl-10 pr-3 py-2.5 text-xs border border-gray-100 bg-slate-50 rounded-lg text-gray-500 cursor-not-allowed">
                <i class="fas fa-envelope absolute left-3.5 top-3.5 text-gray-400 text-xs"></i>
              </div>
            </div>

            <div class="md:col-span-2">
              <label class="block text-xs font-semibold text-gray-700 mb-2">Password Baru <span class="text-[10px] text-gray-400 font-normal">(Kosongkan jika tidak ingin mengubah password)</span></label>
              <div class="relative">
                <input v-model="profilePassword" type="password" placeholder="Masukkan password baru..." class="w-full pl-10 pr-3 py-2.5 text-xs border border-gray-200 rounded-lg focus:outline-none focus:border-green-500 bg-white text-slate-900">
                <i class="fas fa-lock absolute left-3.5 top-3.5 text-gray-400 text-xs"></i>
              </div>
            </div>
          </div>

          <div class="pt-4 border-t border-gray-100 flex justify-end">
            <button type="submit" :disabled="isUpdatingProfile" class="w-full sm:w-auto bg-green-600 hover:bg-green-700 text-white font-bold text-xs px-6 py-3 rounded-lg transition-colors shadow-sm disabled:opacity-50 flex items-center justify-center">
              <i v-if="isUpdatingProfile" class="fas fa-spinner fa-spin mr-2"></i>
              {{ isUpdatingProfile ? 'Menyimpan Perubahan...' : 'Simpan Perubahan' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
