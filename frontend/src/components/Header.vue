<template>
  <header class="bg-gray-800 shadow">
    <nav class="mx-auto flex max-w-7xl items-center justify-between p-6 lg:px-8">
      <div class="flex lg:flex-1">
        <a href="/domains" class="-m-1.5 p-1">
          <span class="sr-only">RedCompass</span>
          <img class="h-10 w-auto" src="../assets/redcompass.webp" alt="RedCompass" />
        </a>
      </div>
      <div class="flex lg:gap-x-12">
        <a href="/domains" class="text-md font-semibold leading-6 text-white">Domains</a>
        <a href="/projects" class="text-md font-semibold leading-6 text-white">Projects</a>
        <a href="/categorizations" class="text-md font-semibold leading-6 text-white">Categorizations</a>
        <a href="/phishing" class="text-md font-semibold leading-6 text-white hover:text-gray-200">Phishing</a>
        <a href="/domain-ideas" class="text-md font-semibold leading-6 text-white">Domain Ideas</a>
      </div>
      <div class="flex lg:flex-1 lg:justify-end">
        <div class="relative avatar-dropdown">
          <button
            @click.stop="toggleDropdown"
            class="flex items-center focus:outline-none"
          >
            <div v-if="avatarUrl">
              <img :src="avatarUrl" class="h-8 w-8 rounded-full" />
            </div>
            <div v-else>
              <div
                class="h-8 w-8 rounded-full bg-blue-600 flex items-center justify-center text-white font-medium"
              >
              {{ pocketbase.authStore.model.name.charAt(0) }}
              </div>
            </div>
          </button>
          
          <!-- Dropdown menu -->
          <div
            v-if="isDropdownOpen"
            class="absolute right-0 mt-2 w-48 rounded-md shadow-lg bg-gray-800 ring-1 ring-black ring-opacity-5"
          >
            <div class="py-1">
              <a
                href="/profile"
                class="block px-4 py-2 text-sm text-white hover:bg-gray-600"
              >
                Profile
              </a>
              <a
                v-if="isAdmin"
                href="/settings"
                class="block px-4 py-2 text-sm text-white hover:bg-gray-600"
              >
                Settings
              </a>
              <a
                href="#"
                @click.prevent="handleLogout"
                class="block px-4 py-2 text-sm text-white hover:bg-gray-600"
              >
                Log out
              </a>
            </div>
          </div>
        </div>
      </div>
    </nav>
  </header>
</template>

<script setup>
import { ref, inject, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const pocketbase = inject('$pocketbase')
const isDropdownOpen = ref(false)
const avatarUrl = ref('')
const isAdmin = ref(false)

onMounted(async () => {
  try {
    const url = pocketbase.files.getURL(pocketbase.authStore.model, pocketbase.authStore.model.avatar);
    avatarUrl.value = url;
    if (pocketbase.authStore.model.role === 'admin') {
      isAdmin.value = true;
    }
  } catch (error) {
    console.error('Error fetching data:', error);
  }
});

const toggleDropdown = () => {
  isDropdownOpen.value = !isDropdownOpen.value
}

const handleLogout = async () => {
  try {
    pocketbase.authStore.clear()
    router.push('/login')
  } catch (error) {
    console.error('Error during logout:', error)
  }
}

// Close dropdown when clicking outside
const closeDropdown = (e) => {
  if (!e.target.closest('.avatar-dropdown')) {
    isDropdownOpen.value = false
  }
}

// Add click outside listener
if (typeof window !== 'undefined') {
  window.addEventListener('click', closeDropdown)
}
</script>

<style scoped>
/* Add some hover effects for navigation links */
a {
  transition: all 0.2s ease;
}
</style> 