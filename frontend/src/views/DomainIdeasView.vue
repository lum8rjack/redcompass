<script setup>
import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'
import PageDescription from '@/components/PageDescription.vue'
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

// Add new refs for confirmation modal
const showConfirmModal = ref(false)
const ideaToDelete = ref(null)

// Add new refs for edit modal
const showEditModal = ref(false)
const editingIdea = ref(null)
const editForm = ref({
  domain: '',
  price: '',
  reason: ''
})

// Add new functions for confirmation modal
const confirmDelete = () => {
  if (ideaToDelete.value) {
    deleteDomainIdea(ideaToDelete.value)
    showConfirmModal.value = false
    ideaToDelete.value = null
  }
}

const cancelDelete = () => {
  showConfirmModal.value = false
  ideaToDelete.value = null
}

// Add new functions for edit modal
const openEditModal = (idea) => {
  editingIdea.value = idea
  editForm.value = {
    domain: idea.Domain,
    price: idea.Price,
    reason: idea.Description
  }
  showEditModal.value = true
}

const closeEditModal = () => {
  showEditModal.value = false
  editingIdea.value = null
  editForm.value = {
    domain: '',
    price: '',
    reason: ''
  }
}

const handleEdit = async () => {
  try {
    formMessage.value = ''
    await pocketbase.collection('Domain_Ideas').update(editingIdea.value.id, {
      Domain: editForm.value.domain,
      Price: editForm.value.price,
      Description: editForm.value.reason
    })
    
    formMessage.value = 'Domain idea updated successfully'
    
    // Clear message after 2 seconds
    setTimeout(() => {
      formMessage.value = ''
    }, 2000)
    
    // Refresh list
    await fetchDomainIdeas()
    closeEditModal()
  } catch (err) {
    formMessage.value = 'Failed to update domain idea'
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
        <PageDescription title="Domain Ideas" description="This page is for collecting domain name ideas that could be valuable additions to our portfolio. Add domains you think would be worth purchasing, along with your reasoning and estimated price." />

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
                <label for="price" class="block text-sm font-medium text-white">Estimated Price (do not include $)</label>
                <input
                  type="text"
                  id="price"
                  name="price"
                  v-model="newDomain.price"
                  required
                  class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-gray-400 focus:ring-gray-400 focus:ring-offset-gray-900 focus:ring-offset-0 sm:text-sm py-1.5 pl-2"
                  placeholder="19.98"
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
                  placeholder="What would this domain be used for?"
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
                    <td 
                      @click="openEditModal(idea)"
                      class="px-6 py-4 whitespace-nowrap text-sm text-white cursor-pointer"
                    >{{ idea.Domain }}</td>
                    <td 
                      @click="openEditModal(idea)"
                      class="px-6 py-4 whitespace-nowrap text-sm text-white cursor-pointer"
                    >${{ idea.Price }}</td>
                    <td 
                      @click="openEditModal(idea)"
                      class="px-6 py-4 text-sm text-white whitespace-pre-line cursor-pointer"
                    >{{ idea.Description }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-white">{{ idea.expand?.User?.name }}</td>
                    <td v-if="!isViewer" class="px-6 py-4 whitespace-nowrap text-sm text-white">
                      <button
                        @click.stop="() => { ideaToDelete = idea.id; showConfirmModal = true; }"
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

    <!-- Edit Modal -->
    <div v-if="showEditModal" class="fixed inset-0 bg-gray-900 bg-opacity-75 flex items-center justify-center z-50">
      <div class="bg-gray-800 rounded-lg p-6 max-w-md w-full mx-4">
        <h3 class="text-lg font-medium text-white mb-4">Edit Domain Idea</h3>
        <form @submit.prevent="handleEdit" class="space-y-4">
          <div>
            <label for="edit-domain" class="block text-sm font-medium text-white">Domain Name</label>
            <input
              type="text"
              id="edit-domain"
              v-model="editForm.domain"
              required
              class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-gray-400 focus:ring-gray-400 focus:ring-offset-gray-900 focus:ring-offset-0 sm:text-sm py-1.5 pl-2"
              placeholder="example.com"
            />
          </div>
          
          <div>
            <label for="edit-price" class="block text-sm font-medium text-white">Estimated Price (do not include $)</label>
            <input
              type="text"
              id="edit-price"
              v-model="editForm.price"
              required
              class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-gray-400 focus:ring-gray-400 focus:ring-offset-gray-900 focus:ring-offset-0 sm:text-sm py-1.5 pl-2"
              placeholder="19.98"
            />
          </div>

          <div>
            <label for="edit-reason" class="block text-sm font-medium text-white">Reason for Purchase</label>
            <textarea
              id="edit-reason"
              v-model="editForm.reason"
              rows="3"
              required
              class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-gray-400 focus:ring-gray-400 focus:ring-offset-gray-900 focus:ring-offset-0 sm:text-sm py-1.5 pl-2"
              placeholder="What would this domain be used for?"
            ></textarea>
          </div>

          <div class="flex justify-end space-x-4">
            <button
              type="button"
              @click="closeEditModal"
              class="px-4 py-2 bg-gray-700 text-white rounded-md hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-gray-500"
            >
              Cancel
            </button>
            <button
              type="submit"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              Save Changes
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Confirmation Modal -->
    <div v-if="showConfirmModal" class="fixed inset-0 bg-gray-900 bg-opacity-75 flex items-center justify-center z-50">
      <div class="bg-gray-800 rounded-lg p-6 max-w-md w-full mx-4">
        <h3 class="text-lg font-medium text-white mb-4">Confirm Domain Idea Deletion</h3>
        <p class="text-gray-300 mb-6">
          Are you sure you want to delete this domain idea? This action cannot be undone.
        </p>
        <div class="flex justify-end space-x-4">
          <button
            @click="cancelDelete"
            class="px-4 py-2 bg-gray-700 text-white rounded-md hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-gray-500"
          >
            Cancel
          </button>
          <button
            @click="confirmDelete"
            class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500"
          >
            Delete
          </button>
        </div>
      </div>
    </div>
  </div>
</template> 