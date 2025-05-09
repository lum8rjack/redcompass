<script setup>
  import Header from '@/components/Header.vue'
  import Footer from '@/components/Footer.vue'
  import { formatDate } from '@/utils/dateUtils'
  import { ref, computed, inject, onMounted } from 'vue'
  import { updateProjectStatus } from '@/utils/projectUtils'

  const pocketbase = inject('$pocketbase')
  const avatarUrl = ref('')
  const projects = ref('')
  const isViewer = ref(false)
  const showConfirmModal = ref(false)
  const projectToComplete = ref(null)

  onMounted(async () => {
    try {
      const url = pocketbase.files.getURL(pocketbase.authStore.model, pocketbase.authStore.model.avatar);
      avatarUrl.value = url;

      if(pocketbase.authStore.model.role === 'viewer') {
        isViewer.value = true;
      }

      const response = await pocketbase.collection('Projects').getFullList({
        sort: 'Name',
        fields: 'id,Name,Start_Date,End_Date,Completed,expand.Project_Members.id',
        expand: 'Project_Members',
        filter: `Project_Members.id ?= "${pocketbase.authStore.model.id}"`,
      });
      projects.value = response;
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  });

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

  // Computed property for sorted projects
  const sortedProjects = computed(() => {
    if (!projects.value) return []
    
    const sorted = [...projects.value]
    const direction = sortConfig.value.direction === 'asc' ? 1 : -1

    return sorted.sort((a, b) => {
      switch (sortConfig.value.key) {
        case 'name':
          return direction * (a.Name || '').localeCompare(b.Name || '')
        case 'timeline':
          const dateA = new Date(a.Start_Date || 0).getTime()
          const dateB = new Date(b.Start_Date || 0).getTime()
          return direction * (dateA - dateB)
        case 'status':
          const statusA = a.Completed ? 1 : 0
          const statusB = b.Completed ? 1 : 0
          return direction * (statusA - statusB)
        default:
          return 0
      }
    })
  })

  // Add function to handle status toggle
  const toggleStatus = async (project) => {
    if (!project.Completed) {
      // If marking as complete, show confirmation first
      projectToComplete.value = project
      showConfirmModal.value = true
    } else {
      // If reactivating, no confirmation needed
      await updateProjectStatus(pocketbase, project)
      // Update local projects array
      const projectIndex = projects.value.findIndex(p => p.id === project.id)
      if (projectIndex !== -1) {
        projects.value[projectIndex] = {
          ...projects.value[projectIndex],
          Completed: !project.Completed
        }
      }
    }
  }

  const confirmComplete = () => {
    if (projectToComplete.value) {
      updateProjectStatus(pocketbase, projectToComplete.value)
        .then(() => {
          const projectIndex = projects.value.findIndex(p => p.id === projectToComplete.value.id)
          if (projectIndex !== -1) {
            projects.value[projectIndex] = {
              ...projects.value[projectIndex],
              Completed: !projectToComplete.value.Completed
            }
          }
        })
        .finally(() => {
          showConfirmModal.value = false
          projectToComplete.value = null
        })
    }
  }

  const cancelComplete = () => {
    showConfirmModal.value = false
    projectToComplete.value = null
  }
</script>

<template>
  <div class="min-h-screen flex flex-col bg-gray-900">
    <Header />
    <main class="flex-grow">
      <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <!-- Profile Header -->
        <div class="bg-gray-800 shadow rounded-lg mb-6">
          <div class="px-4 py-5 sm:p-6">
            <div class="flex items-center">
              <template v-if="avatarUrl">
                <img
                  :src="avatarUrl"
                  :alt="pocketbase.authStore.model.name"
                  class="h-20 w-20 rounded-full object-cover"
                />
              </template>
              <template v-else>
                <div
                  class="h-20 w-20 rounded-full bg-blue-600 flex items-center justify-center text-white text-2xl font-medium"
                >
                  {{ pocketbase.authStore.model.name.charAt(0).toUpperCase() }}
                </div>
              </template>
              <div class="ml-6">
                <h2 class="text-2xl font-bold text-white">
                  {{ pocketbase.authStore.model.name }}
                </h2>
                <div class="mt-1 text-gray-400">
                  {{ pocketbase.authStore.model.email }}
                </div>
                <div class="mt-1 text-sm text-gray-400">
                  {{ pocketbase.authStore.model.role }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Projects Section -->
        <div class="bg-gray-800 shadow rounded-lg">
          <div class="px-4 py-5 sm:p-6">
            <h3 class="text-lg font-medium text-white mb-4">Projects</h3>
            <div class="overflow-hidden rounded-lg">
              <table class="min-w-full divide-y divide-gray-700">
                <thead class="bg-gray-700">
                  <tr>
                    <th 
                      v-for="(header, key) in {
                        name: 'Project Name',
                        timeline: 'Timeline',
                        status: 'Status',
                        actions: 'Actions'
                      }"
                      :key="key"
                      scope="col" 
                      class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider cursor-pointer hover:bg-gray-600 first:rounded-tl-lg last:rounded-tr-lg"
                      @click="key !== 'actions' ? handleSort(key) : null"
                    >
                      <div class="flex items-center space-x-1">
                        <span>{{ header }}</span>
                        <span v-if="sortConfig.key === key && key !== 'actions'" class="text-gray-400">
                          {{ sortConfig.direction === 'asc' ? '↑' : '↓' }}
                        </span>
                      </div>
                    </th>
                  </tr>
                </thead>
                <tbody class="bg-gray-800 divide-y divide-gray-700">
                  <tr 
                    v-for="project in sortedProjects" 
                    :key="project.id"
                    class="cursor-pointer hover:bg-gray-700 transition-colors duration-150"
                  >
                    <td 
                      @click="$router.push(`/project/${project.id}`)"
                      class="px-6 py-4 whitespace-nowrap text-sm font-medium text-white"
                    >
                      {{ project.Name }}
                    </td>
                    <td 
                      @click="$router.push(`/project/${project.id}`)"
                      class="px-6 py-4 whitespace-nowrap text-sm text-white"
                    >
                      {{ formatDate(project.Start_Date) }} - {{ formatDate(project.End_Date) }}
                    </td>
                    <td 
                      @click="$router.push(`/project/${project.id}`)"
                      class="px-6 py-4 whitespace-nowrap text-sm"
                    >
                      <span :class="{
                        'px-2 py-1 rounded-full text-xs font-medium': true,
                        'bg-green-100 text-green-800': project.Completed === true,
                        'bg-gray-100 text-gray-800': project.Completed === false
                      }">
                        {{ project.Completed === true ? 'Completed' : 'In Progress' }}
                      </span>
                    </td>
                    <td v-if="!isViewer" class="px-6 py-4 whitespace-nowrap text-sm">
                      <button
                        @click.stop="toggleStatus(project)"
                        class="px-3 py-1 rounded-md text-sm font-medium"
                        :class="{
                          'bg-green-600 hover:bg-green-700 text-white': project.Completed === true,
                          'bg-gray-600 hover:bg-gray-700 text-white': project.Completed === false
                        }"
                      >
                        {{ project.Completed === true ? 'Reactivate' : 'Mark Complete' }}
                      </button>
                    </td>
                  </tr>
                  <tr class="last:rounded-b-lg">
                    <td class="first:rounded-bl-lg"></td>
                    <td></td>
                    <td></td>
                    <td class="last:rounded-br-lg"></td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </main>
    <Footer />

    <!-- Confirmation Modal -->
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

<style scoped>
/* Add smooth transition for sort indicators */
th {
  transition: background-color 0.2s;
}
</style> 