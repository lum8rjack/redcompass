<script setup>
import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'
import { ref, inject, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { formatDate } from '@/utils/dateUtils'

const route = useRoute()
const pocketbase = inject('$pocketbase')
const project = ref(null)
const domains = ref([])
const loading = ref(true)
const error = ref(null)
const showConfirmModal = ref(false)
const updateMessage = ref('')
const currentUser = pocketbase.authStore.model

// Computed property to check if current user is a project member
const isProjectMember = computed(() => {
  if (!project.value?.expand?.Project_Members) return false
  return project.value.expand.Project_Members.some(member => member.id === currentUser.id)
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
    await updateProjectStatus()
  }
}

const updateProjectStatus = async () => {
  try {
    const response = await pocketbase.collection('Projects').update(project.value.id, {
      "Completed": !project.value.Completed
    });
    
    if(response) {
      // If marking as complete, unassign all domains and set Last_Used
      if (!project.value.Completed) {
        // Get all domains assigned to this project
        const assignedDomains = await pocketbase.collection('Domains').getFullList({
          filter: `Assigned_Project = "${project.value.id}"`
        });

        // Unassign project from each domain and set Last_Used
        for (const domain of assignedDomains) {
          await pocketbase.collection('Domains').update(domain.id, {
            "Assigned_Project": null,
            "Last_Used": project.value.id
          });
        }
      }
      // Update local state
      project.value.Completed = response.Completed;
      // Refresh domains list if marking complete
      if (response.Completed) {
        domains.value = [];
      }
    }
  } catch (error) {
    console.error('Error updating project status:', error);
  } finally {
    showConfirmModal.value = false;
  }
}

const confirmComplete = () => {
  updateProjectStatus()
}

const cancelComplete = () => {
  showConfirmModal.value = false
}
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
            <div class="mb-6">
              <h1 class="text-2xl font-bold text-white">{{ project.Name }}</h1>
              <p class="text-gray-400 mt-2">
                {{ formatDate(project.Start_Date) }} - {{ formatDate(project.End_Date) }}
              </p>
            </div>

            <!-- Status -->
            <div class="mb-6">
              <h2 class="text-lg font-medium text-white mb-2">Status</h2>
              <span :class="{
                'px-3 py-1 rounded-full text-sm font-medium': true,
                'bg-green-100 text-green-800': project.Completed,
                'bg-gray-100 text-gray-800': !project.Completed
              }">
                {{ project.Completed ? 'Completed' : 'Active' }}
              </span>
            </div>

            <!-- Add this after the Status section -->
            <div v-if="isProjectMember" class="mb-6">
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

            <!-- Team Members -->
            <div class="mb-6">
              <h2 class="text-lg font-medium text-white mb-2">Team Members</h2>
              <div class="grid grid-cols-1 gap-2">
                <div 
                  v-for="member in project.expand?.Project_Members" 
                  :key="member.id"
                  class="bg-gray-700 rounded-lg p-3"
                >
                  <p class="text-white">{{ member.name }}</p>
                  <p class="text-gray-400 text-sm">{{ member.role }}</p>
                </div>
              </div>
              <p v-if="!project.expand?.Project_Members?.length" class="text-gray-400">
                No team members assigned
              </p>
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