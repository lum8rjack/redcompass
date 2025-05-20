<script setup>
import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'
import { ref, computed, inject, onMounted } from 'vue'
import { onClickOutside } from '@vueuse/core'
import { formatDate } from '@/utils/dateUtils'
import PageDescription from '@/components/PageDescription.vue'
const pocketbase = inject('$pocketbase');

const availableUsers = ref([])
const projects = ref([])
const isViewer = ref(false)

onMounted(async () => {
    try {
      if(pocketbase.authStore.model.role === 'viewer') {
        isViewer.value = true;
      }

      const response = await pocketbase.collection('Projects').getFullList({
        sort: 'Start_Date',
        fields: 'id,Name,Start_Date,End_Date,Completed,expand.Project_Members.name',
        expand: 'Project_Members',
      });
      projects.value = response;

      const response2 = await pocketbase.collection('users').getFullList({
        sort: 'name',
        fields: 'id,name,role',
        filter: 'role != "viewer"',
      });
      availableUsers.value = response2;
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  });

// Form data
const newProject = ref({
  name: '',
  startDate: '',
  endDate: '',
  selectedUsers: []
})

// Filter states
const filters = ref({
  name: '',
  status: '',
  teamMember: ''
})

// Status options for filter
const statusOptions = ['All', 'Active', 'Completed']

// Sorting state
const sortConfig = ref({
  key: 'name',
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

// Update computed filtered projects to include sorting
const filteredProjects = computed(() => {
  let filtered = projects.value.filter(project => {
    const nameMatch = !filters.value.name || project.Name.toLowerCase().includes(filters.value.name.toLowerCase())
    const statusMatch = !filters.value.status || 
      filters.value.status === 'All' || 
      (filters.value.status === 'Active' && !project.Completed) ||
      (filters.value.status === 'Completed' && project.Completed)
    const teamMatch = !filters.value.teamMember || 
      (project.expand?.Project_Members || [])
        .some(member => member.name.toLowerCase()
        .includes(filters.value.teamMember.toLowerCase()))
    
    return nameMatch && statusMatch && teamMatch
  })

  // Apply sorting
  return filtered.sort((a, b) => {
    let aValue = a[sortConfig.value.key]
    let bValue = b[sortConfig.value.key]

    // Handle array values for Project_Members
    if (sortConfig.value.key === 'Project_Members') {
      aValue = (aValue || []).join(', ')
      bValue = (bValue || []).join(', ')
    }

    // Handle date fields
    if (['Start_Date', 'End_Date'].includes(sortConfig.value.key)) {
      aValue = new Date(aValue)
      bValue = new Date(bValue)
    }

    if (sortConfig.value.direction === 'asc') {
      return aValue > bValue ? 1 : -1
    }
    return aValue < bValue ? 1 : -1
  })
})

const handleCreateProject = async () => {
  console.log('Creating project:', newProject.value)

  try {
    // Create project in Pocketbase
    const data = {
      "Name": newProject.value.name,
      "Start_Date": newProject.value.startDate,
      "End_Date": newProject.value.endDate,
      "Project_Members": newProject.value.selectedUsers,
      "Completed": false
    };
    
    await pocketbase.collection('Projects').create(data);
    const response = await pocketbase.collection('Projects').getFullList({
      sort: 'Start_Date',
      fields: 'id,Name,Start_Date,End_Date,Completed,expand.Project_Members.name',
      expand: 'Project_Members',
    });
    projects.value = response;
  } catch (error) {
    console.log("error creating project", error)
  }

  // Reset the form
  newProject.value = {
    name: '',
    startDate: '',
    endDate: '',
    selectedUsers: []
  }
}

// Add new refs for search functionality
const searchQuery = ref('')
const showResults = ref(false)

// Computed property for filtered users based on search
const filteredUsers = computed(() => {
  if (searchQuery.value.length < 3) return []
  return availableUsers.value.filter(user => 
    user.name.toLowerCase().includes(searchQuery.value.toLowerCase()) &&
    !newProject.value.selectedUsers.includes(user.id)
  )
})

// Function to add user to selection
const addUser = (user) => {
  newProject.value.selectedUsers.push(user.id)
  searchQuery.value = ''
  showResults.value = false
}

// Function to remove user from selection
const removeUser = (userId) => {
  newProject.value.selectedUsers = newProject.value.selectedUsers.filter(id => id !== userId)
}

// Function to get user object by ID
const getUserById = (id) => availableUsers.value.find(user => user.id === id)

// Close search results when clicking outside
const searchContainer = ref(null)
onClickOutside(searchContainer, () => {
  showResults.value = false
})

// Add new refs for status dropdown
const showStatusDropdown = ref(false)
const statusContainer = ref(null)

// Close dropdown when clicking outside
onClickOutside(statusContainer, () => {
  showStatusDropdown.value = false
})

// Calculate progress based on start and end dates
const calculateProgress = (startDate, endDate) => {
  const start = new Date(startDate)
  const end = new Date(endDate)
  const today = new Date()

  // If project hasn't started yet
  if (today < start) return 0
  // If project is completed or past end date
  if (today > end) return 100

  const totalDuration = end - start
  const elapsed = today - start
  return Math.round((elapsed / totalDuration) * 100)
}
</script>

<template>
  <div class="min-h-screen flex flex-col bg-gray-900">
    <Header />
    <main class="flex-grow">
      <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <PageDescription title="Project Management" description="Manage your projects, track their progress, and assign team members. Create new projects and monitor their status from start to completion." />

        <!-- Project Creation Form -->
        <div v-if="!isViewer" class="bg-gray-800 shadow rounded-lg mb-6">
          <div class="px-4 py-5 sm:p-6">
            <h3 class="text-lg font-medium text-white mb-4">Create New Project</h3>
            <form @submit.prevent="handleCreateProject" class="space-y-4">
              <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <div>
                  <label for="projectName" class="block text-sm font-medium text-white">Project Name</label>
                  <input
                    id="projectName"
                    v-model="newProject.name"
                    type="text"
                    required
                    class="mt-1 block w-full rounded-md shadow-sm focus:border-gray-400 focus:ring-gray-400 sm:text-sm bg-gray-700 border-gray-600 text-white py-1.5 pl-2"
                  >
                </div>
                <div>
                  <label for="startDate" class="block text-sm font-medium text-white">Start Date</label>
                  <input
                    id="startDate"
                    v-model="newProject.startDate"
                    type="date"
                    required
                    class="mt-1 block w-full rounded-md shadow-sm focus:border-gray-400 focus:ring-gray-400 sm:text-sm bg-gray-700 border-gray-600 text-white py-1.5 pl-2"
                  >
                </div>
                <div>
                  <label for="endDate" class="block text-sm font-medium text-white">End Date</label>
                  <input
                    id="endDate"
                    v-model="newProject.endDate"
                    type="date"
                    required
                    class="mt-1 block w-full rounded-md shadow-sm focus:border-gray-400 focus:ring-gray-400 sm:text-sm bg-gray-700 border-gray-600 text-white py-1.5 pl-2"
                  >
                </div>
                <div>
                  <label for="teamSearch" class="block text-sm font-medium text-white">Assign Team Members</label>
                  <div class="mt-1 relative" ref="searchContainer">
                    <!-- Search input -->
                    <input
                      id="teamSearch"
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
                        v-for="userId in newProject.selectedUsers"
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
              </div>
              <div class="flex justify-center">
                <button
                  type="submit"
                  class="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-gray-600 hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-400"
                >
                  Create Project
                </button>
              </div>
            </form>
          </div>
        </div>

        <!-- Projects Table -->
        <div class="bg-gray-800 shadow rounded-lg">
          <div class="px-4 py-5 sm:p-6">
            <h3 class="text-lg font-medium text-white mb-4">All Projects</h3>
            
            <!-- Filter Section -->
            <div class="grid grid-cols-1 gap-4 sm:grid-cols-3 mb-6">
              <div>
                <label for="nameFilter" class="block text-sm font-medium text-white">Filter by Name</label>
                <input
                  type="text"
                  id="nameFilter"
                  v-model="filters.name"
                  placeholder="Search project name..."
                  class="mt-1 block w-full rounded-md shadow-sm focus:border-gray-400 focus:ring-gray-400 sm:text-sm bg-gray-700 border-gray-600 text-white py-1.5 pl-2"
                >
              </div>
              <div>
                <label for="statusFilter" class="block text-sm font-medium text-white">Filter by Status</label>
                <div class="mt-1 relative" ref="statusContainer" id="statusContainer">
                  <!-- Dropdown trigger button -->
                  <button
                    type="button"
                    id="statusFilter"
                    @click="showStatusDropdown = !showStatusDropdown"
                    class="w-full bg-gray-700 border border-gray-600 rounded-md py-1.5 pl-2 pr-10 text-left text-white focus:outline-none focus:ring-1 focus:ring-gray-400 focus:border-gray-400 sm:text-sm"
                  >
                    {{ filters.status || 'Select status' }}
                    <span class="absolute inset-y-0 right-0 flex items-center pr-2">
                      <svg class="h-5 w-5 text-gray-400" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
                      </svg>
                    </span>
                  </button>

                  <!-- Dropdown menu -->
                  <div 
                    v-if="showStatusDropdown"
                    class="absolute z-10 mt-1 w-full bg-gray-700 rounded-md shadow-lg"
                  >
                    <ul class="max-h-60 rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none sm:text-sm">
                      <li
                        v-for="status in statusOptions"
                        :key="status"
                        @click="() => { 
                          filters.status = status;
                          showStatusDropdown = false;
                        }"
                        class="text-white cursor-pointer select-none relative py-2 pl-3 pr-9 hover:bg-gray-600"
                        :class="{ 'bg-gray-600': filters.status === status }"
                      >
                        {{ status }}
                      </li>
                    </ul>
                  </div>
                </div>
              </div>
              <div>
                <label for="teamFilter" class="block text-sm font-medium text-white">Filter by Team Member</label>
                <input
                  type="text"
                  id="teamFilter"
                  v-model="filters.teamMember"
                  placeholder="Search team member..."
                  class="mt-1 block w-full rounded-md shadow-sm focus:border-gray-400 focus:ring-gray-400 sm:text-sm bg-gray-700 border-gray-600 text-white py-1.5 pl-2"
                >
              </div>
            </div>

            <!-- Results count -->
            <p class="text-sm text-gray-500 dark:text-gray-400 mb-4">
              Showing {{ filteredProjects.length }} of {{ projects.length }} projects
            </p>

            <div class="overflow-hidden rounded-lg">
              <table class="min-w-full divide-y divide-gray-700">
                <thead class="bg-gray-700">
                  <tr>
                    <th 
                      v-for="(header, key) in {
                        Name: 'Project Name',
                        Start_Date: 'Timeline',
                        Completed: 'Status',
                        Project_Members: 'Team Members',
                        Progress: 'Progress',
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
                  <tr 
                    v-for="project in filteredProjects" 
                    :key="project.id"
                    @click="$router.push(`/project/${project.id}`)"
                    class="cursor-pointer hover:bg-gray-700 transition-colors duration-150"
                  >
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">
                      {{ project.Name }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-white">
                      {{ formatDate(project.Start_Date) }} - {{ formatDate(project.End_Date) }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm">
                      <span :class="{
                        'px-2 py-1 rounded-full text-xs font-medium': true,
                        'bg-green-100 text-green-800': project.Completed === true,
                        'bg-gray-100 text-gray-800': project.Completed === false
                      }">
                        {{ project.Completed ? 'Completed' : 'Active' }}
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-white">
                      {{ project.expand?.Project_Members?.map(member => member.name)?.join(', ') || 'No team members' }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-white">
                      <div class="w-full bg-gray-200 rounded-full h-2.5">
                        <div
                          class="bg-gray-600 h-2.5 rounded-full"
                          :style="{ width: `${calculateProgress(project.Start_Date, project.End_Date)}%` }"
                        ></div>
                      </div>
                      <span class="text-xs mt-1 block">
                        {{ calculateProgress(project.Start_Date, project.End_Date) }}%
                      </span>
                    </td>
                  </tr>
                  <tr class="last:rounded-b-lg">
                    <td class="first:rounded-bl-lg"></td>
                    <td></td>
                    <td></td>
                    <td></td>
                    <td class="last:rounded-br-lg"></td>
                  </tr>
                </tbody>
              </table>

              <!-- No results message -->
              <div v-if="filteredProjects.length === 0" class="text-center py-4 text-gray-500 dark:text-gray-400">
                No projects match your filters
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
    <Footer />
  </div>
</template>

<style scoped>
select[multiple] {
  min-height: 120px;
}

select[multiple] option {
  padding: 0.5rem;
  cursor: pointer;
}

select[multiple] option:checked {
  background: linear-gradient(0deg, #4f46e5 0%, #4f46e5 100%);
  color: white;
}

/* Add smooth transition for sort indicators */
th {
  transition: background-color 0.2s;
}

/* Ensure dropdown appears above other elements */
.relative {
  z-index: 50;
}

/* Smooth transitions */
.absolute {
  transition: all 0.2s;
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

/* Ensure consistent input heights */
input #statusContainer {
  height: 38px;
}
</style> 