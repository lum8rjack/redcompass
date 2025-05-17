<script setup>
import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'
import { ref, inject, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { formatDate } from '@/utils/dateUtils'
import { onClickOutside } from '@/composables/onClickOutside'
import { updateProjectStatus } from '@/utils/projectUtils'

const route = useRoute()
const pocketbase = inject('$pocketbase')
const project = ref(null)
const domains = ref([])
const loading = ref(true)
const error = ref(null)
const showConfirmModal = ref(false)
const updateMessage = ref('')
const currentUser = pocketbase.authStore.model

// Add new refs for edit modal
const showEditModal = ref(false)
const editForm = ref({
  startDate: '',
  endDate: '',
  selectedUsers: []
})
const availableUsers = ref([])

// Add new refs for search functionality
const searchQuery = ref('')
const showResults = ref(false)
const searchContainer = ref(null)

// Computed property to check if current user is a project member
const isProjectMember = computed(() => {
  if (!project.value?.expand?.Project_Members) return false
  return project.value.expand.Project_Members.some(member => member.id === currentUser.id)
})

// Add new computed property for filtered users
const filteredUsers = computed(() => {
  if (searchQuery.value.length < 3) return []
  return availableUsers.value.filter(user => 
    user.name.toLowerCase().includes(searchQuery.value.toLowerCase()) &&
    !editForm.value.selectedUsers.includes(user.id)
  )
})

onMounted(async () => {
  try {
    const projectId = route.params.id
    const [projectRecord, domainsRecord] = await Promise.all([
      pocketbase.collection('Projects').getOne(projectId, {
        expand: 'Project_Members'
      }),
      pocketbase.collection('Domains').getFullList({
        filter: `Assigned_Project = "${projectId}"`,
        sort: 'Name'
      })
    ])
    
    project.value = projectRecord
    domains.value = domainsRecord
  } catch (err) {
    error.value = 'Failed to load project details'
    console.error('Error fetching project:', err)
  } finally {
    loading.value = false
  }
})

const toggleStatus = async () => {
  if (!project.value.Completed) {
    // If marking as complete, show confirmation first
    showConfirmModal.value = true
  } else {
    // If reactivating, no confirmation needed
    await handleProjectStatusUpdate()
  }
}

const handleProjectStatusUpdate = async () => {
  try {
    const response = await updateProjectStatus(pocketbase, project.value)
    // Update local state
    project.value.Completed = response.Completed
    // Refresh domains list if marking complete
    if (response.Completed) {
      domains.value = []
    } else {
      // If reactivating, refetch domains
      const domainsRecord = await pocketbase.collection('Domains').getFullList({
        filter: `Assigned_Project = "${project.value.id}"`,
        sort: 'Name'
      })
      domains.value = domainsRecord
    }
  } catch (error) {
    console.error('Error updating project status:', error)
  } finally {
    showConfirmModal.value = false
  }
}

const confirmComplete = () => {
  handleProjectStatusUpdate()
}

const cancelComplete = () => {
  showConfirmModal.value = false
}

// Add new functions for edit modal
const openEditModal = async () => {
  try {
    // Fetch available users
    const response = await pocketbase.collection('users').getFullList({
      sort: 'name',
      fields: 'id,name,role',
    });
    availableUsers.value = response;

    // Format dates for date inputs (YYYY-MM-DD)
    const formatDateForInput = (dateString) => {
      const date = new Date(dateString);
      return date.toISOString().split('T')[0];
    };

    // Set current values in form
    editForm.value = {
      startDate: formatDateForInput(project.value.Start_Date),
      endDate: formatDateForInput(project.value.End_Date),
      selectedUsers: project.value.expand?.Project_Members?.map(member => member.id) || []
    }
    showEditModal.value = true
  } catch (error) {
    console.error('Error opening edit modal:', error)
  }
}

const closeEditModal = () => {
  showEditModal.value = false
  editForm.value = {
    startDate: '',
    endDate: '',
    selectedUsers: []
  }
}

const handleEdit = async () => {
  try {
    updateMessage.value = ''
    await pocketbase.collection('Projects').update(project.value.id, {
      Start_Date: editForm.value.startDate,
      End_Date: editForm.value.endDate,
      Project_Members: editForm.value.selectedUsers
    })
    
    // Refresh project data
    const projectRecord = await pocketbase.collection('Projects').getOne(project.value.id, {
      expand: 'Project_Members'
    })
    project.value = projectRecord
    
    updateMessage.value = 'Project updated successfully'
    setTimeout(() => {
      updateMessage.value = ''
    }, 2000)
    
    closeEditModal()
  } catch (error) {
    updateMessage.value = 'Failed to update project'
    console.error('Error updating project:', error)
  }
}

// Add function to add user to selection
const addUser = (user) => {
  editForm.value.selectedUsers.push(user.id)
  searchQuery.value = ''
  showResults.value = false
}

// Add function to remove user from selection
const removeUser = (userId) => {
  editForm.value.selectedUsers = editForm.value.selectedUsers.filter(id => id !== userId)
}

// Add function to get user object by ID
const getUserById = (id) => availableUsers.value.find(user => user.id === id)

// Close search results when clicking outside
onClickOutside(searchContainer, () => {
  showResults.value = false
})
</script>

<template>
  <div class="min-h-screen flex flex-col bg-gray-900">
    <Header />
    <main class="flex-grow">
      <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <!-- Loading State -->
        <div v-if="loading" class="text-white text-center py-8">
          Loading project details...
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="text-red-500 text-center py-8">
          {{ error }}
        </div>

        <!-- Project Details -->
        <div v-else-if="project" class="bg-gray-800 shadow rounded-lg">
          <div class="px-4 py-5 sm:p-6">
            <div class="flex justify-between items-start mb-6">
              <div>
                <h1 class="text-2xl font-bold text-white">{{ project.Name }}</h1>
                <p class="text-gray-400 mt-2">
                  {{ formatDate(project.Start_Date) }} - {{ formatDate(project.End_Date) }}
                </p>
              </div>
              <button
                v-if="isProjectMember && !project.Completed"
                @click="openEditModal"
                class="px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-400"
              >
                Edit Project
              </button>
            </div>

            <!-- Status and Team Members in a row -->
            <div class="flex flex-col sm:flex-row gap-6 mb-6">
              <!-- Status Section -->
              <div class="flex-1">
                <h2 class="text-lg font-medium text-white mb-2">Status</h2>
                <span :class="{
                  'px-3 py-1 rounded-full text-sm font-medium': true,
                  'bg-green-100 text-green-800': project.Completed,
                  'bg-gray-100 text-gray-800': !project.Completed
                }">
                  {{ project.Completed ? 'Completed' : 'Active' }}
                </span>

                <!-- Status Toggle Button -->
                <div v-if="isProjectMember" class="mt-4">
                  <button
                    @click="toggleStatus"
                    class="px-4 py-2 rounded-md text-sm font-medium"
                    :class="{
                      'bg-green-600 hover:bg-green-700 text-white': project.Completed,
                      'bg-gray-600 hover:bg-gray-700 text-white': !project.Completed
                    }"
                  >
                    {{ project.Completed ? 'Reactivate Project' : 'Mark Project Complete' }}
                  </button>
                </div>
              </div>

              <!-- Team Members Section -->
              <div class="flex-1">
                <h2 class="text-lg font-medium text-white mb-2">Team Members</h2>
                <div class="grid grid-cols-1 gap-2">
                  <div 
                    v-for="member in project.expand?.Project_Members" 
                    :key="member.id"
                    class="bg-gray-700 rounded-lg p-3"
                  >
                    <p class="text-white">{{ member.name }}</p>
                  </div>
                </div>
                <p v-if="!project.expand?.Project_Members?.length" class="text-gray-400">
                  No team members assigned
                </p>
              </div>
            </div>

            <!-- Progress -->
            <div>
              <h2 class="text-lg font-medium text-white mb-2">Progress</h2>
              <div class="w-full bg-gray-200 rounded-full h-2.5">
                <div
                  class="bg-gray-600 h-2.5 rounded-full"
                  :style="{ width: `${calculateProgress(project.Start_Date, project.End_Date)}%` }"
                ></div>
              </div>
              <span class="text-sm text-gray-400 mt-1 block">
                {{ calculateProgress(project.Start_Date, project.End_Date) }}% Complete
              </span>
            </div>
          </div>
        </div>

        <!-- Domains Section -->
        <div class="mt-6">
          <h2 class="text-lg font-medium text-white mb-2">Assigned Domains</h2>
          <div class="bg-gray-700 rounded-lg overflow-hidden">
            <div v-if="domains.length === 0" class="p-4 text-gray-400">
              No domains assigned to this project
            </div>
            <table v-else class="min-w-full divide-y divide-gray-600">
              <thead class="bg-gray-800">
                <tr>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">
                    Domain Name
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">
                    Health Status
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">
                    Purchased Date
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">
                    Expiry Date
                  </th>
                </tr>
              </thead>
              <tbody class="bg-gray-700 divide-y divide-gray-600">
                <tr 
                  v-for="domain in domains" 
                  :key="domain.id"
                  @click="$router.push(`/domain/${domain.id}`)"
                  class="cursor-pointer hover:bg-gray-600 transition-colors duration-150"
                >
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-white">
                    {{ domain.Name }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm">
                    <span :class="{
                      'px-2 py-1 rounded-full text-xs font-medium': true,
                      'bg-green-100 text-green-800': domain.Healthy,
                      'bg-red-100 text-red-800': !domain.Healthy
                    }">
                      {{ domain.Healthy ? 'Healthy' : 'Unhealthy' }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300">
                    {{ formatDate(domain.Purchased_Date) }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300">
                    {{ formatDate(domain.Expiration_Date) }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </main>
    <Footer />

    <!-- Edit Modal -->
    <div v-if="showEditModal && !project.Completed" class="fixed inset-0 bg-gray-900 bg-opacity-75 flex items-center justify-center z-50">
      <div class="bg-gray-800 rounded-lg p-6 max-w-md w-full mx-4">
        <h3 class="text-lg font-medium text-white mb-4">Edit Project</h3>
        <form @submit.prevent="handleEdit" class="space-y-4">
          <div>
            <label for="startDate" class="block text-sm font-medium text-white">Start Date</label>
            <input
              type="date"
              id="startDate"
              v-model="editForm.startDate"
              required
              class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-gray-400 focus:ring-gray-400 focus:ring-offset-gray-900 focus:ring-offset-0 sm:text-sm py-1.5 pl-2"
            />
          </div>

          <div>
            <label for="endDate" class="block text-sm font-medium text-white">End Date</label>
            <input
              type="date"
              id="endDate"
              v-model="editForm.endDate"
              required
              class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-gray-400 focus:ring-gray-400 focus:ring-offset-gray-900 focus:ring-offset-0 sm:text-sm py-1.5 pl-2"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-white mb-2">Team Members</label>
            <div class="mt-1 relative" ref="searchContainer">
              <!-- Search input -->
              <input
                type="text"
                v-model="searchQuery"
                @focus="showResults = true"
                placeholder="Search team members..."
                class="mt-1 block w-full rounded-md shadow-sm focus:border-gray-400 focus:ring-gray-400 sm:text-sm bg-gray-700 border-gray-600 text-white py-1.5 pl-2"
              >
              
              <!-- Search results dropdown -->
              <div 
                v-if="showResults && filteredUsers.length > 0"
                class="absolute z-10 mt-1 w-full bg-gray-700 rounded-md shadow-lg"
              >
                <ul class="max-h-60 overflow-auto rounded-md py-1 text-base">
                  <li
                    v-for="user in filteredUsers"
                    :key="user.id"
                    @click="addUser(user)"
                    class="cursor-pointer px-4 py-2 text-white hover:bg-gray-600"
                  >
                    {{ user.name }}
                  </li>
                </ul>
              </div>

              <!-- Selected users display -->
              <div class="mt-2 flex flex-wrap gap-2">
                <div
                  v-for="userId in editForm.selectedUsers"
                  :key="userId"
                  class="inline-flex items-center bg-gray-600 rounded-full px-3 py-1 text-sm text-white"
                >
                  <span>{{ getUserById(userId).name }}</span>
                  <button
                    @click="removeUser(userId)"
                    class="ml-2 text-gray-300 hover:text-white"
                  >
                    ×
                  </button>
                </div>
              </div>

              <!-- Helper text -->
              <p v-if="searchQuery.length > 0 && searchQuery.length < 3" class="mt-2 text-sm text-gray-400">
                Type at least 3 characters to search
              </p>
              <p v-else-if="searchQuery.length >= 3 && filteredUsers.length === 0" class="mt-2 text-sm text-gray-400">
                No matching team members found
              </p>
            </div>
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

    <!-- Add Confirmation Modal -->
    <div v-if="showConfirmModal" class="fixed inset-0 bg-gray-900 bg-opacity-75 flex items-center justify-center z-50">
      <div class="bg-gray-800 rounded-lg p-6 max-w-md w-full mx-4">
        <h3 class="text-lg font-medium text-white mb-4">Confirm Project Completion</h3>
        <p class="text-gray-300 mb-6">
          Marking this project as complete will unassign all associated domains. Are you sure you want to continue?
        </p>
        <div class="flex justify-end space-x-4">
          <button
            @click="cancelComplete"
            class="px-4 py-2 bg-gray-700 text-white rounded-md hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-gray-500"
          >
            Cancel
          </button>
          <button
            @click="confirmComplete"
            class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500"
          >
            Confirm
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
// Reuse the calculateProgress function from ProjectsView
function calculateProgress(startDate, endDate) {
  const start = new Date(startDate)
  const end = new Date(endDate)
  const today = new Date()

  if (today < start) return 0
  if (today > end) return 100

  const totalDuration = end - start
  const elapsed = today - start
  return Math.round((elapsed / totalDuration) * 100)
}
</script>

<style scoped>
/* Add smooth transition for sort indicators */
th {
  transition: background-color 0.2s;
}

/* Update the calendar icon and input styles */
input[type="date"]::-webkit-calendar-picker-indicator {
  filter: invert(1);
  cursor: pointer;
  margin-right: 4px;
}

input[type="date"]::-moz-calendar-picker-indicator {
  filter: invert(1);
  cursor: pointer;
  margin-right: 4px;
}

input[type="date"]::-ms-calendar-picker-indicator {
  filter: invert(1);
  cursor: pointer;
  margin-right: 4px;
}
</style> 