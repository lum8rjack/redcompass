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
        <div class="relative domains-dropdown">
          <button 
            @click.stop="toggleDomainsDropdown"
            class="text-md font-semibold leading-6 text-white hover:text-gray-200 flex items-center"
          >
            Domains
            <svg class="h-4 w-4 ml-1" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
            </svg>
          </button>
          <div
            v-if="isDomainsDropdownOpen" 
            class="absolute left-0 mt-2 w-56 rounded-md shadow-lg bg-gray-800 ring-1 ring-black ring-opacity-5"
          >
            <div class="py-1">
              <a href="/domains" class="block px-4 py-2 text-sm text-white hover:bg-gray-600">Domain Management</a>
              <a href="/categorizations" class="block px-4 py-2 text-sm text-white hover:bg-gray-600">Domain Categorization</a>
              <a href="/domain-ideas" class="block px-4 py-2 text-sm text-white hover:bg-gray-600">Domain Ideas</a>
            </div>
          </div>
        </div>
        <a href="/projects" class="text-md font-semibold leading-6 text-white hover:text-gray-200">Projects</a>
        <div class="relative phishing-dropdown">
          <button 
            @click.stop="togglePhishingDropdown"
            class="text-md font-semibold leading-6 text-white hover:text-gray-200 flex items-center"
          >
            Phishing
            <svg class="h-4 w-4 ml-1" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
            </svg>
          </button>
          
          <div
            v-if="isPhishingDropdownOpen" 
            class="absolute left-0 mt-2 w-48 rounded-md shadow-lg bg-gray-800 ring-1 ring-black ring-opacity-5"
          >
            <div class="py-1">
              <a
                href="/phishing"
                class="block px-4 py-2 text-sm text-white hover:bg-gray-600"
              >
                Campaign Templates
              </a>
              <a
                href="/phishing-metrics"
                class="block px-4 py-2 text-sm text-white hover:bg-gray-600"
              >
                Metrics
              </a>
              <a
                href="/phishlets"
                class="block px-4 py-2 text-sm text-white hover:bg-gray-600"
              >
                Phishlets
              </a>
            </div>
          </div>
        </div>
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
const isPhishingDropdownOpen = ref(false)
const isDomainsDropdownOpen = ref(false)
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

const togglePhishingDropdown = () => {
  isPhishingDropdownOpen.value = !isPhishingDropdownOpen.value
}

const toggleDomainsDropdown = () => {
  isDomainsDropdownOpen.value = !isDomainsDropdownOpen.value
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
  if (!e.target.closest('.phishing-dropdown')) {
    isPhishingDropdownOpen.value = false
  }
  if (!e.target.closest('.domains-dropdown')) {
    isDomainsDropdownOpen.value = false
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