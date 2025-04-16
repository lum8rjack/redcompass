<script setup>
import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'
import { inject, ref, onMounted, computed } from 'vue'

const pocketbase = inject('$pocketbase')
const domainIdeas = ref([])
const loading = ref(true)
const error = ref(null)
const isViewer = ref(false)
const formMessage = ref('')

// Sorting state
const sortConfig = ref({
  key: 'Domain',
  direction: 'asc'
})

// Update sorting
const handleSort = (key) => {
  if (sortConfig.value.key === key) {
    // If clicking the same header, toggle direction
    sortConfig.value.direction = sortConfig.value.direction === 'asc' ? 'desc' : 'asc'
  } else {
    // If clicking a new header, set it with ascending direction
    sortConfig.value.key = key
    sortConfig.value.direction = 'asc'
  }
}

// Computed property for sorted domain ideas
const sortedDomainIdeas = computed(() => {
  return [...domainIdeas.value].sort((a, b) => {
    let aValue = a[sortConfig.value.key]
    let bValue = b[sortConfig.value.key]

    // Handle user name sorting
    if (sortConfig.value.key === 'User') {
      aValue = a.expand?.User?.name || ''
      bValue = b.expand?.User?.name || ''
    }

    if (sortConfig.value.direction === 'asc') {
      return aValue > bValue ? 1 : -1
    }
    return aValue < bValue ? 1 : -1
  })
})

// Form data
const newDomain = ref({
  name: '',
  price: '',
  reason: ''
})

// Fetch domain ideas
const fetchDomainIdeas = async () => {
  try {
    if(pocketbase.authStore.model.role === 'viewer') {
      isViewer.value = true;
    }

    const records = await pocketbase.collection('Domain_Ideas').getFullList({
      sort: '-Domain',
      fields: 'id,Domain,Price,Description,expand.User',
      expand: 'User'
    })
    domainIdeas.value = records
  } catch (err) {
    error.value = 'Failed to load domain ideas'
  } finally {
    loading.value = false
  }
}

// Submit new domain idea
const submitDomainIdea = async () => {
  try {
    formMessage.value = ''
    await pocketbase.collection('Domain_Ideas').create({
      Domain: newDomain.value.name,
      Price: newDomain.value.price,
      Description: newDomain.value.reason,
      User: pocketbase.authStore.model.id
    })
    
    // Reset form
    newDomain.value = {
      name: '',
      price: '',
      reason: ''
    }
    
    // Show success message
    formMessage.value = 'Domain idea added successfully'
    
    // Clear message after 2 seconds
    setTimeout(() => {
      formMessage.value = ''
    }, 2000)
    
    // Refresh list
    await fetchDomainIdeas()
  } catch (err) {
    formMessage.value = 'Failed to add domain idea'
  }
}

// Delete domain idea
const deleteDomainIdea = async (id) => {
  try {
    formMessage.value = ''
    await pocketbase.collection('Domain_Ideas').delete(id)
    formMessage.value = 'Domain idea deleted successfully'
    
    // Clear message after 2 seconds
    setTimeout(() => {
      formMessage.value = ''
    }, 2000)
    
    // Refresh list
    await fetchDomainIdeas()
  } catch (err) {
    formMessage.value = 'Failed to delete domain idea'
  }
}

onMounted(() => {
  fetchDomainIdeas()
})
</script>

<template>
  <div class="min-h-screen flex flex-col bg-gray-900">
    <Header />
    <main class="flex-grow">
      <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <!-- Page Description -->
        <div class="mb-8">
          <h1 class="text-2xl font-bold text-white mb-2">Domain Ideas</h1>
          <p class="text-gray-400">
            This page is for collecting domain name ideas that could be valuable additions to our portfolio.
            Add domains you think would be worth purchasing, along with your reasoning and estimated price.
          </p>
        </div>

        <!-- Add New Domain Form -->
        <div v-if="!isViewer" class="bg-gray-800 shadow rounded-lg mb-8">
          <div class="px-4 py-5 sm:px-6">
            <h2 class="text-lg font-medium text-white mb-4">Add New Domain Idea</h2>
            <form @submit.prevent="submitDomainIdea" class="space-y-4">
              <div>
                <label for="domain-name" class="block text-sm font-medium text-white">Domain Name</label>
                <input
                  type="text"
                  id="domain-name"
                  name="domain-name"
                  v-model="newDomain.name"
                  required
                  class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-gray-400 focus:ring-gray-400 focus:ring-offset-gray-900 focus:ring-offset-0 sm:text-sm py-1.5 pl-2"
                  placeholder="example.com"
                />
              </div>

              <div>
                <label for="price" class="block text-sm font-medium text-white">Estimated Price</label>
                <input
                  type="text"
                  id="price"
                  name="price"
                  v-model="newDomain.price"
                  required
                  class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-gray-400 focus:ring-gray-400 focus:ring-offset-gray-900 focus:ring-offset-0 sm:text-sm py-1.5 pl-2"
                  placeholder="$1000"
                />
              </div>

              <div>
                <label for="reason" class="block text-sm font-medium text-white">Reason for Purchase</label>
                <textarea
                  id="reason"
                  name="reason"
                  v-model="newDomain.reason"
                  rows="3"
                  required
                  class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-gray-400 focus:ring-gray-400 focus:ring-offset-gray-900 focus:ring-offset-0 sm:text-sm py-1.5 pl-2"
                  placeholder="Why this domain would be valuable..."
                ></textarea>
              </div>

              <div class="flex flex-col items-center space-y-2">
                <button
                  type="submit"
                  class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-gray-600 hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-400"
                >
                  Add Domain Idea
                </button>
                <p v-if="formMessage" class="text-sm" :class="{
                  'text-green-400': formMessage.includes('success'),
                  'text-red-400': formMessage.includes('Failed')
                }">
                  {{ formMessage }}
                </p>
              </div>
            </form>
          </div>
        </div>

        <!-- Domain Ideas Table -->
        <div class="bg-gray-800 shadow rounded-lg">
          <div class="px-4 py-5 sm:p-6">
            <h2 class="text-lg font-medium text-white mb-4">Submitted Domain Ideas</h2>
            
            <!-- Loading State -->
            <div v-if="loading" class="text-white text-center py-8">
              Loading domain ideas...
            </div>

            <!-- Error State -->
            <div v-else-if="error" class="text-red-500 text-center py-8">
              {{ error }}
            </div>

            <!-- Table -->
            <div v-else class="overflow-x-auto ring-1 ring-gray-700 rounded-lg">
              <table class="min-w-full divide-y divide-gray-700">
                <thead class="bg-gray-700">
                  <tr>
                    <th 
                      scope="col" 
                      class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider first:rounded-tl-lg cursor-pointer hover:bg-gray-600"
                      @click="handleSort('Domain')"
                    >
                      <div class="flex items-center space-x-1">
                        <span>Domain</span>
                        <span v-if="sortConfig.key === 'Domain'" class="text-gray-400">
                          {{ sortConfig.direction === 'asc' ? '↑' : '↓' }}
                        </span>
                      </div>
                    </th>
                    <th 
                      scope="col" 
                      class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider cursor-pointer hover:bg-gray-600"
                      @click="handleSort('Price')"
                    >
                      <div class="flex items-center space-x-1">
                        <span>Price</span>
                        <span v-if="sortConfig.key === 'Price'" class="text-gray-400">
                          {{ sortConfig.direction === 'asc' ? '↑' : '↓' }}
                        </span>
                      </div>
                    </th>
                    <th 
                      scope="col" 
                      class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider cursor-pointer hover:bg-gray-600"
                      @click="handleSort('Description')"
                    >
                      <div class="flex items-center space-x-1">
                        <span>Reason</span>
                        <span v-if="sortConfig.key === 'Description'" class="text-gray-400">
                          {{ sortConfig.direction === 'asc' ? '↑' : '↓' }}
                        </span>
                      </div>
                    </th>
                    <th 
                      scope="col" 
                      class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider cursor-pointer hover:bg-gray-600"
                      @click="handleSort('User')"
                    >
                      <div class="flex items-center space-x-1">
                        <span>Submitted By</span>
                        <span v-if="sortConfig.key === 'User'" class="text-gray-400">
                          {{ sortConfig.direction === 'asc' ? '↑' : '↓' }}
                        </span>
                      </div>
                    </th>
                    <th 
                      v-if="!isViewer" 
                      scope="col" 
                      class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider last:rounded-tr-lg"
                    >
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody class="bg-gray-800 divide-y divide-gray-700">
                  <tr v-for="idea in sortedDomainIdeas" :key="idea.id" class="hover:bg-gray-700">
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-white">{{ idea.Domain }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-white">${{ idea.Price }}</td>
                    <td class="px-6 py-4 text-sm text-white">{{ idea.Description }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-white">{{ idea.expand?.User?.name }}</td>
                    <td v-if="!isViewer" class="px-6 py-4 whitespace-nowrap text-sm text-white">
                      <button
                        @click="deleteDomainIdea(idea.id)"
                        class="text-red-400 hover:text-red-300"
                      >
                        Delete
                      </button>
                    </td>
                  </tr>
                  <tr v-if="domainIdeas.length === 0">
                    <td :colspan="isViewer ? 4 : 5" class="px-6 py-4 text-center text-sm text-gray-400">
                      No domain ideas submitted yet
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </main>
    <Footer />
  </div>
</template> 