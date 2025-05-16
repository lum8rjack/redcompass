<script setup>
import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'
import PageDescription from '@/components/PageDescription.vue'
import { formatDate } from '@/utils/dateUtils'
import { ref, inject, computed, onMounted } from 'vue'
import { canEdit } from '@/utils/auth'
const pocketbase = inject('$pocketbase')

// Phishlets data
const phishlets = ref([])

// Sorting state
const sortConfig = ref({
  key: 'Name',
  direction: 'asc'
})

const handleSort = (key) => {
  if (sortConfig.value.key === key) {
    sortConfig.value.direction = sortConfig.value.direction === 'asc' ? 'desc' : 'asc'
  } else {
    sortConfig.value.key = key
    sortConfig.value.direction = 'asc'
  }
}

async function getPhishlets() {
  try {
      const response = await pocketbase.collection('Phishlets').getFullList({
          sort: 'Name',
          fields: 'id,Name,expand.Uploaded_By,Phishlet,Notes,updated',
          expand: 'Uploaded_By',
      });
      phishlets.value = response;
  } catch (error) {
      console.error('Error fetching phishlets:', error)
  }
}

onMounted(async () => {
  await getPhishlets()
})

const selectedPhishlet = ref(null)
const showPhishletDetailModal = ref(false)
const editedPhishletContent = ref('')
const editedPhishletNotes = ref('')
const newPhishletName = ref('')
const newPhishletContent = ref('')
const newPhishletNotes = ref('')
const showDeleteConfirmModal = ref(false)
const phishletToDelete = ref(null)
const updateMessage = ref('')
const sortedPhishlets = computed(() => {
  if (!phishlets.value) return []
  const sorted = [...phishlets.value]
  const direction = sortConfig.value.direction === 'asc' ? 1 : -1
  return sorted.sort((a, b) => {
    switch (sortConfig.value.key) {
      case 'Name':
        return direction * (a.Name || '').localeCompare(b.Name || '')
      case 'Last Updated':
        return direction * (new Date(a.updated).getTime() - new Date(b.updated).getTime())
      case 'Uploaded By':
        return direction * (a.Uploaded_By.Name || '').localeCompare(b.Uploaded_By.name || '')
      default:
        return 0
    }
  })
})

function viewPhishletDetails(phishlet) {
  selectedPhishlet.value = { ...phishlet }
  editedPhishletContent.value = phishlet.Phishlet || ''
  editedPhishletNotes.value = phishlet.Notes || ''
  showPhishletDetailModal.value = true
}

function closePhishletDetailModal() {
  showPhishletDetailModal.value = false
  selectedPhishlet.value = null
}

async function savePhishletDetails() {
  if (!selectedPhishlet.value) return
  
  // Create project in Pocketbase
  try {
    const data = {
        "Phishlet": editedPhishletContent.value,
        "Notes": editedPhishletNotes.value,
        "Uploaded_By": pocketbase.authStore.model.id,
    };

    await pocketbase.collection('Phishlets').update(selectedPhishlet.value.id, data)
    updateMessageBox('Successfully updated phishlet details')
    await getPhishlets()
    closePhishletDetailModal()
  } catch (error) {
    closePhishletDetailModal()
    updateMessageBox('Error updating phishlet details')
  }
}

function updateMessageBox(message) {
  updateMessage.value = message
  setTimeout(() => {
    updateMessage.value = ''
  }, 2000)
}

async function createPhishlet() {
  if (!newPhishletName.value || !newPhishletContent.value) {
    updateMessageBox('Error: You must provide a name and phishlet content')
    return
  }
  try {

    // Create project in Pocketbase
    const data = {
      "Name": newPhishletName.value,
      "Phishlet": newPhishletContent.value,
      "Notes": newPhishletNotes.value,
      "Uploaded_By": pocketbase.authStore.model.id,
    };
    
    await pocketbase.collection('Phishlets').create(data);
    // Fetch all phishlets again
    await getPhishlets()
    newPhishletName.value = ''
    newPhishletContent.value = ''
    newPhishletNotes.value = ''
    updateMessageBox('Successfully created phishlet')
  } catch (error) {
    updateMessageBox('Error creating phishlet')
  }
}

async function deletePhishlet(phishlet) {
  phishletToDelete.value = phishlet
  showDeleteConfirmModal.value = true
}

async function confirmDeletePhishlet() {
  if (!phishletToDelete.value || !phishletToDelete.value.id) return
  try {
    await pocketbase.collection('Phishlets').delete(phishletToDelete.value.id)
    await getPhishlets()
    updateMessageBox('Successfully deleted phishlet')
  } catch (error) {
    updateMessageBox('Error deleting phishlet')
  } finally {
    showDeleteConfirmModal.value = false
    phishletToDelete.value = null
  }
}

function cancelDeletePhishlet() {
  showDeleteConfirmModal.value = false
  phishletToDelete.value = null
}

</script>

<template>
  <div class="min-h-screen flex flex-col bg-gray-900">
    <Header />
    <main class="flex-grow">
      <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <PageDescription title="Phishlet Tracking" description="Track and manage standalone phishlets for use in your campaigns." />
        <!-- Create Phishlet Form -->
        <div v-if="canEdit()" class="bg-gray-800 shadow rounded-lg mb-8 p-6">
          <h2 class="text-xl font-semibold text-white mb-4">Create New Phishlet</h2>
          <form @submit.prevent="createPhishlet" class="space-y-4">
            <div>
              <label for="phishletName" class="block text-sm font-medium text-white">Phishlet Name</label>
              <input
                type="text"
                id="phishletName"
                v-model="newPhishletName"
                class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 shadow-sm focus:border-gray-400 focus:ring-gray-400 sm:text-sm text-white py-1.5 pl-2"
                placeholder="Enter phishlet name"
                required
              />
            </div>
            <div>
              <label for="phishletContent" class="block text-sm font-medium text-white">Phishlet Content</label>
              <textarea
                id="phishletContent"
                v-model="newPhishletContent"
                class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 shadow-sm focus:border-gray-400 focus:ring-gray-400 sm:text-sm text-white p-2 font-mono"
                rows="6"
                placeholder="Paste phishlet YAML or text here"
                required
              ></textarea>
            </div>
            <div>
              <label for="phishletNotes" class="block text-sm font-medium text-white">Notes</label>
              <textarea
                id="phishletNotes"
                v-model="newPhishletNotes"
                class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 shadow-sm focus:border-gray-400 focus:ring-gray-400 sm:text-sm text-white p-2"
                rows="3"
                placeholder="Add any notes about this phishlet (optional)"
              ></textarea>
            </div>
            <div class="flex flex-col items-center">
              <button
                type="submit"
                class="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-gray-600 hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-400"
              >
                Create Phishlet
              </button>
              <p v-if="updateMessage" class="mt-2 text-sm" :class="{
                'text-green-400': updateMessage.includes('Success'),
                'text-red-400': updateMessage.includes('Error')
              }">
                {{ updateMessage }}
              </p>
            </div>
          </form>
        </div>
        <!-- End Create Phishlet Form -->
        <div class="bg-gray-800 shadow overflow-hidden rounded-lg mb-4">
          <table class="min-w-full divide-y divide-gray-700">
            <thead class="bg-gray-700">
              <tr>
                <th 
                    v-for="(header, key) in {
                    Name: 'Phishlet Name',
                    Start_Date: 'Last Updated',
                    Completed: 'Uploaded By',
                    Project_Members: 'Actions',
                    }"
                    :key="key"
                    scope="col" 
                    class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider cursor-pointer hover:bg-gray-600 first:rounded-tl-lg last:rounded-tr-lg"
                    @click="handleSort(key)"
                >
                    <div class="flex items-center space-x-1">
                    <span>{{ header }}</span>
                    <span v-if="sortConfig.key === key" class="text-gray-400">
                        {{ sortConfig.direction === 'asc' ? '↑' : '↓' }}
                    </span>
                    </div>
                </th>
              </tr>
            </thead>
            <tbody class="bg-gray-800 divide-y divide-gray-700">
              <tr v-for="phishlet in sortedPhishlets" :key="phishlet.id" class="hover:bg-gray-700">
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-white cursor-pointer" @click="viewPhishletDetails(phishlet)">
                  {{ phishlet.Name }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300 cursor-pointer" @click="viewPhishletDetails(phishlet)">
                  {{ formatDate(phishlet.updated) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300 cursor-pointer" @click="viewPhishletDetails(phishlet)">
                  {{ phishlet.expand?.Uploaded_By?.name || 'Unknown' }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <button v-if="canEdit()" @click="deletePhishlet(phishlet)" class="bg-red-600 hover:bg-red-700 text-white text-sm font-medium py-1 px-3 rounded-md">
                    Delete
                  </button>
                </td>
              </tr>
              <tr v-if="phishlets.length === 0">
                <td colspan="4" class="px-6 py-4 text-center text-sm text-gray-400">
                  No phishlets found
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </main>
    <Footer />

    <!-- Phishlet Details Modal -->
    <div v-if="showPhishletDetailModal" class="fixed inset-0 z-50 overflow-y-auto">
      <div class="fixed inset-0 bg-black bg-opacity-60 backdrop-blur-sm transition-opacity" @click="closePhishletDetailModal"></div>
      <div class="flex items-center justify-center min-h-screen p-4">
        <div class="bg-gray-800 rounded-lg overflow-hidden shadow-xl transform transition-all sm:max-w-3xl w-full relative mx-auto my-8">
          <div class="bg-gray-700 px-6 py-4 flex items-center justify-between">
            <h3 class="text-lg font-medium text-white">
              Phishlet Details: {{ selectedPhishlet?.Name }}
            </h3>
            <button 
              @click="closePhishletDetailModal" 
              class="text-gray-400 hover:text-white focus:outline-none"
            >
              <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <div class="p-6 bg-gray-800">
            <div class="mb-4">
              <label class="block text-sm font-medium text-white mb-2">Phishlet Content</label>
              <textarea
                v-model="editedPhishletContent"
                class="w-full bg-gray-700 border-gray-600 rounded-md text-white p-3 font-mono"
                rows="10"
                style="resize: vertical;"
              ></textarea>
            </div>
            <div class="mb-4">
              <label class="block text-sm font-medium text-white mb-2">Notes</label>
              <textarea
                v-model="editedPhishletNotes"
                class="w-full bg-gray-700 border-gray-600 rounded-md text-white p-3"
                rows="5"
                style="resize: vertical;"
              ></textarea>
            </div>
            <div class="flex justify-center">
              <button
                v-if="canEdit()"
                @click="savePhishletDetails"
                class="bg-gray-600 hover:bg-gray-700 text-white text-sm font-medium py-2 px-4 rounded"
              >
                Save Details
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteConfirmModal" class="fixed inset-0 bg-gray-900 bg-opacity-75 flex items-center justify-center z-50">
      <div class="bg-gray-800 rounded-lg p-6 max-w-md w-full mx-4">
        <h3 class="text-lg font-medium text-white">Confirm Delete</h3>
        <p class="text-sm text-gray-300 mb-4">
          Are you sure you want to delete the phishlet '{{ phishletToDelete?.Name }}'? This action cannot be undone.
        </p>
        <div class="flex justify-end space-x-4">
          <button
            @click="cancelDeletePhishlet"
            class="px-4 py-2 bg-gray-700 text-white rounded-md hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-gray-500"
          >
            Cancel
          </button>
          <button
            @click="confirmDeletePhishlet"
            class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500"
          >
            Confirm
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Add any component-specific styles here */
</style> 