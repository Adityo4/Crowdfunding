<script setup>
import { ref, computed } from 'vue'
import { useRoute, useFetch } from '#app'

const route = useRoute()
const slug = route.params.slug

// Fetch article detail from backend
const { data: articleResponse, error } = await useFetch(`http://localhost:8080/v1/articles/${slug}`)

const article = computed(() => articleResponse.value?.data || null)

// Comments list (local mock/reactive state for comments since backend details may vary)
const comments = ref([
  {
    id: 1,
    name: 'John Smith',
    time: '2 jam yang lalu',
    content: 'Kisah yang sangat menginspirasi! Luar biasa sekali melihat bagaimana pendidikan mampu mengubah jalan hidup seseorang. Sukses selalu untuk YPAI!',
    avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?auto=format&fit=crop&w=100&q=80'
  },
  {
    id: 2,
    name: 'Emily Johnson',
    time: '4 jam yang lalu',
    content: 'Terima kasih YPAI atas kerja kerasnya menyalurkan kepedulian dari donatur secara transparan dan berdaya guna.',
    avatar: 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?auto=format&fit=crop&w=100&q=80'
  }
])

const newCommentName = ref('')
const newCommentEmail = ref('')
const newCommentBody = ref('')

const submitComment = (e) => {
  e.preventDefault()
  if (!newCommentName.value || !newCommentBody.value) return
  
  comments.value.unshift({
    id: Date.now(),
    name: newCommentName.value,
    time: 'Baru saja',
    content: newCommentBody.value,
    avatar: 'https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?auto=format&fit=crop&w=100&q=80'
  })
  
  newCommentName.value = ''
  newCommentEmail.value = ''
  newCommentBody.value = ''
}

const copyLink = () => {
  if (navigator.clipboard) {
    navigator.clipboard.writeText(window.location.href)
    alert('Link disalin ke clipboard!')
  }
}
</script>

<template>
  <div class="bg-slate-50 min-h-screen">
    <!-- Error/Not Found State -->
    <div v-if="error || !article" class="max-w-4xl mx-auto px-6 py-20 text-center">
      <i class="fas fa-exclamation-triangle text-5xl text-yellow-500 mb-4 animate-bounce"></i>
      <h1 class="text-2xl font-bold text-slate-800">Artikel Tidak Ditemukan</h1>
      <p class="text-gray-500 mt-2">Maaf, artikel yang Anda cari tidak dapat ditemukan atau telah dihapus.</p>
      <NuxtLink to="/articles" class="mt-6 inline-block bg-green-600 hover:bg-green-700 text-white font-semibold px-6 py-2.5 rounded-lg transition-colors">
        Kembali ke Daftar Artikel
      </NuxtLink>
    </div>

    <div v-else>
      <!-- Article Hero Section -->
      <section class="pt-24 pb-12 bg-gradient-to-br from-green-50/50 to-white">
        <div class="max-w-4xl mx-auto px-6">
          <!-- Breadcrumb -->
          <nav class="mb-8">
            <ol class="flex items-center space-x-2 text-sm text-gray-500">
              <li><NuxtLink to="/" class="hover:text-green-600 transition-colors">Home</NuxtLink></li>
              <li><i class="fas fa-chevron-right text-[10px]"></i></li>
              <li><NuxtLink to="/articles" class="hover:text-green-600 transition-colors">Artikel</NuxtLink></li>
              <li><i class="fas fa-chevron-right text-[10px]"></i></li>
              <li class="text-slate-800 font-medium line-clamp-1">{{ article.title }}</li>
            </ol>
          </nav>

          <!-- Article Meta -->
          <div class="mb-8">
            <div class="flex items-center space-x-4 mb-4">
              <span class="px-3 py-1 bg-green-50 text-green-700 text-xs font-semibold rounded-full">Kisah Sukses</span>
              <span class="text-xs text-gray-400">•</span>
              <span class="text-xs text-gray-500">
                {{ new Date(article.publishedAt || article.createdAt).toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' }) }}
              </span>
            </div>
            <h1 class="text-3xl md:text-4xl font-extrabold text-slate-900 mb-6 leading-tight">
              {{ article.title }}
            </h1>
            
            <!-- Author Info -->
            <div class="flex items-center space-x-3">
              <div class="w-10 h-10 bg-green-100 rounded-full flex items-center justify-center text-green-700 font-bold">
                A
              </div>
              <div>
                <div class="font-semibold text-slate-900 text-sm">Administrator</div>
                <div class="text-xs text-gray-500">Yayasan Peduli Amal Indonesia</div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Article Content -->
      <section class="py-10 bg-white border-t border-gray-100">
        <div class="max-w-4xl mx-auto px-6">
          <div class="grid grid-cols-1 lg:grid-cols-4 gap-12">
            <!-- Main Content -->
            <div class="lg:col-span-3">
              <!-- Featured Image -->
              <div class="mb-10 rounded-2xl overflow-hidden shadow-md bg-slate-100">
                <img :src="article.coverImageUrl || 'https://images.unsplash.com/photo-1523050854058-8df90110c9f1?auto=format&fit=crop&w=1200&q=80'" 
                     :alt="article.title" 
                     class="w-full h-auto object-cover max-h-[400px]">
              </div>

              <!-- Article Body -->
              <article class="prose max-w-none text-slate-700 leading-relaxed space-y-6" v-html="article.content">
              </article>

              <!-- Share Buttons -->
              <div class="mt-12 pt-8 border-t border-gray-100">
                <h3 class="text-sm font-bold text-slate-900 mb-4 uppercase tracking-wider">Bagikan Artikel</h3>
                <div class="flex flex-wrap gap-3">
                  <a href="#" class="flex items-center space-x-2 px-4 py-2 bg-blue-600 text-white rounded-lg text-xs font-semibold hover:bg-blue-700 transition-colors">
                    <i class="fab fa-facebook-f"></i>
                    <span>Facebook</span>
                  </a>
                  <a href="#" class="flex items-center space-x-2 px-4 py-2 bg-sky-500 text-white rounded-lg text-xs font-semibold hover:bg-sky-600 transition-colors">
                    <i class="fab fa-twitter"></i>
                    <span>Twitter</span>
                  </a>
                  <button @click="copyLink" class="flex items-center space-x-2 px-4 py-2 bg-slate-600 text-white rounded-lg text-xs font-semibold hover:bg-slate-700 transition-colors">
                    <i class="fas fa-link"></i>
                    <span>Salin Link</span>
                  </button>
                </div>
              </div>
            </div>

            <!-- Sidebar -->
            <div class="lg:col-span-1">
              <!-- Author Card -->
              <div class="bg-slate-50 border border-slate-100 rounded-2xl p-6 mb-6">
                <h3 class="text-sm font-bold text-slate-900 mb-4">Mengenai Penulis</h3>
                <div class="flex items-center space-x-3 mb-4">
                  <div class="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center text-green-700 font-bold text-lg">
                    A
                  </div>
                  <div>
                    <div class="font-semibold text-slate-900 text-sm">Administrator</div>
                    <div class="text-xs text-gray-500">Tim YPAI</div>
                  </div>
                </div>
                <p class="text-xs text-gray-500 leading-relaxed">
                  Administrator resmi Yayasan Peduli Amal Indonesia. Memberikan informasi terpercaya mengenai aksi kemanusiaan dan penyaluran bantuan sosial.
                </p>
              </div>

              <!-- Newsletter Banner -->
              <div class="bg-gradient-to-br from-green-600 to-emerald-700 rounded-2xl p-6 text-white">
                <h3 class="font-bold text-base mb-2">Pantau Aksi Kami</h3>
                <p class="text-xs text-green-100 mb-4">Dapatkan update berkala mengenai program donasi dan kisah inspiratif lainnya.</p>
                <form class="space-y-2" @submit.prevent="alert('Terima kasih sudah berlangganan!')">
                  <input type="email" required placeholder="Email Anda" class="w-full px-3 py-2 text-xs rounded-lg text-slate-950 placeholder-gray-400 focus:outline-none">
                  <button type="submit" class="w-full bg-white text-green-700 py-2 rounded-lg font-bold text-xs hover:bg-green-50 transition-colors">
                    Langganan
                  </button>
                </form>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Comments Section -->
      <section class="py-12 bg-slate-50 border-t border-gray-100">
        <div class="max-w-4xl mx-auto px-6">
          <div class="bg-white border border-slate-100 rounded-2xl p-6 md:p-8">
            <h3 class="text-lg font-bold text-slate-900 mb-6">Diskusi & Komentar ({{ comments.length }})</h3>
            
            <!-- Comment Form -->
            <form @submit="submitComment" class="mb-8 p-4 bg-slate-50 rounded-xl space-y-3">
              <h4 class="text-xs font-bold text-slate-700 mb-2">Tulis komentar Anda</h4>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                <input v-model="newCommentName" type="text" required placeholder="Nama lengkap *" class="px-3 py-2 text-xs border border-gray-200 rounded-lg focus:border-green-500 focus:outline-none bg-white text-slate-900">
                <input v-model="newCommentEmail" type="email" placeholder="Alamat email" class="px-3 py-2 text-xs border border-gray-200 rounded-lg focus:border-green-500 focus:outline-none bg-white text-slate-900">
              </div>
              <textarea v-model="newCommentBody" required placeholder="Komentar Anda *" rows="3" class="w-full px-3 py-2 text-xs border border-gray-200 rounded-lg focus:border-green-500 focus:outline-none bg-white text-slate-900"></textarea>
              <button type="submit" class="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-lg text-xs font-bold transition-colors">
                Kirim Komentar
              </button>
            </form>

            <!-- Comments List -->
            <div class="space-y-6">
              <div v-for="comment in comments" :key="comment.id" class="flex space-x-4">
                <img :src="comment.avatar" :alt="comment.name" class="w-9 h-9 rounded-full object-cover">
                <div class="flex-1">
                  <div class="flex items-center space-x-2 mb-1">
                    <span class="font-bold text-slate-800 text-sm">{{ comment.name }}</span>
                    <span class="text-xs text-gray-400">{{ comment.time }}</span>
                  </div>
                  <p class="text-gray-600 text-xs leading-relaxed">{{ comment.content }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.prose :deep(h1), .prose :deep(h2), .prose :deep(h3) {
  color: #0f172a;
  font-weight: 700;
  margin-top: 1.5rem;
  margin-bottom: 0.75rem;
}
.prose :deep(h2) {
  font-size: 1.25rem;
}
.prose :deep(p) {
  margin-bottom: 1rem;
}
.prose :deep(blockquote) {
  border-left: 4px solid #059669;
  padding-left: 1rem;
  margin: 1.5rem 0;
  font-style: italic;
  color: #4b5563;
}
</style>
