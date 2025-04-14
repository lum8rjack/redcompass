<script setup>
import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'
import { ref, inject, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { formatDate } from '@/utils/dateUtils'
import { onClickOutside } from '@vueuse/core'

const route = useRoute()
const pocketbase = inject('$pocketbase')
const domain = ref(null)
const dnsRecords = ref([])
const loading = ref(true)
const error = ref(null)
const notes = ref('')
const isHealthy = ref(false)
const updateMessage = ref('')
const availableProjects = ref([])
const searchQuery = ref('')
const showResults = ref(false)
const searchRef = ref(null)
const newTag = ref('')
const showTagInput = ref(false)

// Predefined list of approved tags
const approvedTags = [
  'Generic',
  'Admin',
  'C2',
  'Email',
  'Hosting'
]

// Computed property for available tags (tags that haven't been used yet)
const availableTags = computed(() => {
  return approvedTags.filter(tag => !domain.value?.Tags?.includes(tag))
})

onMounted(async () => {
  try {
    const domainId = route.params.id
    const record = await pocketbase.collection('Domains').getOne(domainId, {
      expand: 'Assigned_Project,Last_Used'
    })
    domain.value = record
    console.log('domain.value', domain.value)
    notes.value = record.Notes || ''
    isHealthy.value = record.Healthy || false

    // Fetch DNS records for the domain
    const records = await pocketbase.collection('Domain_Records').getFullList({
      filter: `Domain = "${domainId}"`,
      sort: '+Record_Type,+Record_Name'
    })
    dnsRecords.value = records

    // Fetch available projects
    const projectsResponse = await pocketbase.collection('Projects').getFullList({
      fields: 'id,Name,Completed,expand.Project_Members.name',
      sort: 'Name',
      filter: 'Completed = false'  // Only show active projects
    })
    availableProjects.value = projectsResponse
  } catch (err) {
    error.value = 'Failed to load domain details'
  } finally {
    loading.value = false
  }
})

// Computed property for filtered projects based on search
const filteredProjects = computed(() => {
  if (!searchQuery.value) return []
  return availableProjects.value.filter(project => 
    project.Name.toLowerCase().includes(searchQuery.value.toLowerCase()) &&
    project.id !== domain.value?.Assigned_Project
  )
})

// Add this computed property to check if user is not a viewer
const canEdit = computed(() => {
  const user = pocketbase.authStore.model
  return user && user.role !== 'viewer'
})

const updateDomain = async () => {
  try {
    updateMessage.value = ''
    const data = {
      "Notes": notes.value,
      "Healthy": isHealthy.value
    }
    await pocketbase.collection('Domains').update(domain.value.id, data)
    updateMessage.value = 'Domain updated successfully'
  } catch (err) {
    updateMessage.value = 'Failed to update domain'
  }
}

// Function to assign project
const assignProject = async (project) => {
  try {
    updateMessage.value = ''
    await pocketbase.collection('Domains').update(domain.value.id, {
      "Assigned_Project": project.id
    })
    // Update local state
    domain.value.Assigned_Project = project.id
    domain.value.expand = {
      ...domain.value.expand,
      Assigned_Project: project
    }
    updateMessage.value = 'Project assigned successfully'
    // Clear search
    searchQuery.value = ''
    showResults.value = false
  } catch (err) {
    updateMessage.value = 'Failed to assign project'
  }
}

// Function to remove project assignment
const removeProject = async () => {
  try {
    updateMessage.value = ''
    await pocketbase.collection('Domains').update(domain.value.id, {
      "Assigned_Project": null
    })
    // Update local state
    domain.value.Assigned_Project = null
    domain.value.expand.Assigned_Project = null
    updateMessage.value = 'Project removed successfully'
  } catch (err) {
    updateMessage.value = 'Failed to remove project'
  }
}

// Function to add a new tag
const addTag = async (tagToAdd) => {
  if (!tagToAdd) return
  try {
    updateMessage.value = ''
    const updatedTags = [...(domain.value.Tags || []), tagToAdd]
    await pocketbase.collection('Domains').update(domain.value.id, {
      "Tags": updatedTags
    })
    domain.value.Tags = updatedTags
    showTagInput.value = false
    updateMessage.value = 'Tag added successfully'
    // Clear message after 2 seconds
    setTimeout(() => {
      updateMessage.value = ''
    }, 2000)
  } catch (err) {
    updateMessage.value = 'Failed to add tag'
  }
}

// Function to remove a tag
const removeTag = async (tagToRemove) => {
  try {
    updateMessage.value = ''
    const updatedTags = (domain.value.Tags || []).filter(tag => tag !== tagToRemove)
    await pocketbase.collection('Domains').update(domain.value.id, {
      "Tags": updatedTags
    })
    domain.value.Tags = updatedTags
    updateMessage.value = 'Tag removed successfully'
    // Clear message after 2 seconds
    setTimeout(() => {
      updateMessage.value = ''
    }, 2000)
  } catch (err) {
    updateMessage.value = 'Failed to remove tag'
  }
}

// Function to update health status
const updateHealth = async (status) => {
  try {
    updateMessage.value = ''
    await pocketbase.collection('Domains').update(domain.value.id, {
      "Healthy": status
    })
    isHealthy.value = status
    updateMessage.value = 'Health status updated successfully'
    // Clear message after 2 seconds
    setTimeout(() => {
      updateMessage.value = ''
    }, 2000)
  } catch (err) {
    updateMessage.value = 'Failed to update health status'
    // Revert the UI state on error
    isHealthy.value = !status
  }
}

// Close search results when clicking outside
onClickOutside(searchRef, () => {
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
          Loading domain details...
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="text-red-500 text-center py-8">
          {{ error }}
        </div>

        <!-- Domain Details -->
        <div v-else-if="domain" class="bg-gray-800 shadow rounded-lg">
          <div class="px-4 py-5 sm:p-6">
            <!-- Domain Name and Basic Info -->
            <div class="mb-6">
              <h1 class="text-2xl font-bold text-white">{{ domain.Name }}</h1>
              <p class="text-gray-400 mt-2">
                Registered: {{ formatDate(domain.Purchased_Date) }} | 
                Expires: {{ formatDate(domain.Expiration_Date) }}
              </p>
            </div>

            <!-- Domain Status Grid -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
              <!-- Provider and Health Status (Left Side) -->
              <div class="space-y-4 flex flex-col">
                <div class="bg-gray-700 rounded-lg p-3 flex-grow">
                  <h3 class="text-sm font-medium text-gray-400 mb-1">Domain Provider</h3>
                  <p class="text-white">{{ domain.Domain_Provider || 'Not specified' }}</p>
                </div>

                <div class="bg-gray-700 rounded-lg p-3 flex-grow">
                  <h3 class="text-sm font-medium text-gray-400 mb-2">Health Status</h3>
                  <div class="flex items-center space-x-2">
                    <button
                      v-if="canEdit"
                      @click="updateHealth(true)"
                      :class="{
                        'px-4 py-2 rounded-md text-sm font-medium transition-colors duration-150': true,
                        'bg-green-600 text-white': isHealthy,
                        'bg-gray-600 text-gray-300 hover:bg-gray-500': !isHealthy
                      }"
                    >
                      Healthy
                    </button>
                    <span
                      v-else
                      :class="{
                        'px-4 py-2 rounded-md text-sm font-medium': true,
                        'bg-green-600 text-white': isHealthy,
                        'bg-gray-600 text-gray-300': !isHealthy
                      }"
                    >
                      Healthy
                    </span>
                    <button
                      v-if="canEdit"
                      @click="updateHealth(false)"
                      :class="{
                        'px-4 py-2 rounded-md text-sm font-medium transition-colors duration-150': true,
                        'bg-red-600 text-white': !isHealthy,
                        'bg-gray-600 text-gray-300 hover:bg-gray-500': isHealthy
                      }"
                    >
                      Unhealthy
                    </button>
                    <span
                      v-else
                      :class="{
                        'px-4 py-2 rounded-md text-sm font-medium': true,
                        'bg-red-600 text-white': !isHealthy,
                        'bg-gray-600 text-gray-300': isHealthy
                      }"
                    >
                      Unhealthy
                    </span>
                  </div>
                </div>

                <div class="bg-gray-700 rounded-lg p-3">
                  <h3 class="text-sm font-medium text-gray-400 mb-2">Assigned Project</h3>
                  
                  <!-- Show current project if assigned -->
                  <div v-if="domain.expand?.Assigned_Project" class="mb-2 flex items-center justify-between">
                    <span class="text-white">{{ domain.expand.Assigned_Project.Name }}</span>
                    <button
                      v-if="canEdit"
                      @click="removeProject"
                      class="text-sm text-red-400 hover:text-red-300"
                    >
                      Remove
                    </button>
                  </div>

                  <!-- Project search -->
                  <div v-else-if="canEdit" ref="searchRef" class="relative">
                    <input
                      type="text"
                      v-model="searchQuery"
                      @focus="showResults = true"
                      name="project-search"
                      id="project-search"
                      placeholder="Search for a project..."
                      class="w-full rounded-md bg-gray-600 border-gray-500 text-white placeholder-gray-400 text-sm focus:border-gray-400 focus:ring-gray-400 focus:ring-offset-gray-900 focus:ring-offset-0 py-1.5 pl-2"
                    >
                    
                    <!-- Search results dropdown -->
                    <div
                      v-if="showResults && filteredProjects.length > 0"
                      class="absolute z-10 mt-1 w-full bg-gray-700 rounded-md shadow-lg max-h-60 overflow-auto"
                    >
                      <ul class="py-1">
                        <li
                          v-for="project in filteredProjects"
                          :key="project.id"
                          @click="assignProject(project)"
                          class="px-3 py-2 text-white hover:bg-gray-600 cursor-pointer text-sm"
                        >
                          {{ project.Name }}
                        </li>
                      </ul>
                    </div>
                    
                    <!-- No results message -->
                    <div
                      v-if="showResults && searchQuery && filteredProjects.length === 0"
                      class="absolute z-10 mt-1 w-full bg-gray-700 rounded-md shadow-lg"
                    >
                      <p class="px-3 py-2 text-sm text-gray-400">No matching projects found</p>
                    </div>
                  </div>

                  <!-- Show message for viewers when no project is assigned -->
                  <div v-else-if="!domain.expand?.Assigned_Project" class="text-gray-400 text-sm">
                    No project assigned
                  </div>
                </div>

                <div class="bg-gray-700 rounded-lg p-3">
                  <h3 class="text-sm font-medium text-gray-400 mb-2">Tags</h3>
                  
                  <!-- Display existing tags -->
                  <div class="flex flex-wrap gap-2 mb-2">
                    <div
                      v-for="tag in domain.Tags"
                      :key="tag"
                      class="flex items-center bg-gray-600 text-white px-3 py-1 rounded-full text-sm"
                    >
                      <span>{{ tag }}</span>
                      <button
                        v-if="canEdit"
                        @click="removeTag(tag)"
                        class="ml-2 text-gray-400 hover:text-gray-300"
                      >
                        ×
                      </button>
                    </div>
                    <div
                      v-if="domain.Tags?.length === 0"
                      class="text-gray-400 text-sm"
                    >
                      No tags
                    </div>
                  </div>

                  <!-- Add tag button and input -->
                  <div v-if="canEdit" class="mt-2">
                    <div v-if="showTagInput" class="flex flex-col gap-2">
                      <div class="flex flex-wrap gap-2">
                        <button
                          v-for="tag in availableTags"
                          :key="tag"
                          @click="addTag(tag)"
                          class="px-3 py-1 bg-gray-600 text-white rounded-full hover:bg-gray-500 text-sm"
                        >
                          + {{ tag }}
                        </button>
                      </div>
                      <button
                        @click="showTagInput = false"
                        class="px-3 py-1 bg-gray-600 text-white rounded-md hover:bg-gray-500 text-sm self-start"
                      >
                        Cancel
                      </button>
                    </div>
                    <button
                      v-else
                      @click="showTagInput = true"
                      class="text-sm text-gray-400 hover:text-gray-300"
                    >
                      + Add Tag
                    </button>
                  </div>
                </div>
              </div>

              <!-- Status Indicators (Right Side) -->
              <div class="space-y-4 flex flex-col">
                <div class="flex items-center justify-between bg-gray-700 rounded-lg p-4 flex-grow">
                  <span class="text-white">Expiration Status</span>
                  <span :class="{
                    'px-2 py-1 rounded-full text-xs font-medium': true,
                    'bg-red-100 text-red-800': domain.Is_Expired,
                    'bg-green-100 text-green-800': !domain.Is_Expired
                  }">
                    {{ domain.Is_Expired ? 'Expired' : 'Active' }}
                  </span>
                </div>

                <div class="flex items-center justify-between bg-gray-700 rounded-lg p-4 flex-grow">
                  <span class="text-white">Auto Renew</span>
                  <span :class="{
                    'px-2 py-1 rounded-full text-xs font-medium': true,
                    'bg-green-100 text-green-800': domain.Auto_Renew,
                    'bg-red-100 text-red-800': !domain.Auto_Renew
                  }">
                    {{ domain.Auto_Renew ? 'Enabled' : 'Disabled' }}
                  </span>
                </div>

                <div class="flex items-center justify-between bg-gray-700 rounded-lg p-4 flex-grow">
                  <span class="text-white">Custom DNS</span>
                  <span :class="{
                    'px-2 py-1 rounded-full text-xs font-medium': true,
                    'bg-red-100 text-red-800': domain.Custom_DNS,
                    'bg-green-100 text-green-800': !domain.Custom_DNS
                  }">
                    {{ domain.Custom_DNS ? 'Yes' : 'No' }}
                  </span>
                </div>

                <div class="flex items-center justify-between bg-gray-700 rounded-lg p-4 flex-grow">
                  <span class="text-white">Domain Lock</span>
                  <span :class="{
                    'px-2 py-1 rounded-full text-xs font-medium': true,
                    'bg-red-100 text-red-800': domain.Is_Locked,
                    'bg-green-100 text-green-800': !domain.Is_Locked
                  }">
                    {{ domain.Is_Locked ? 'Yes' : 'No' }}
                  </span>
                </div>

                <div class="flex items-center justify-between bg-gray-700 rounded-lg p-4 flex-grow">
                  <span class="text-white">Previous Project</span>
                  <span class="px-2 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                    {{ domain.expand?.Last_Used?.Name || 'None' }}
                  </span>
                </div>
              </div>
            </div>

            <!-- Notes Section -->
            <div class="mb-6">
              <h2 class="text-lg font-medium text-white mb-2">Notes</h2>
              <!-- Editable textarea for non-viewers -->
              <textarea
                v-if="canEdit"
                v-model="notes"
                rows="4"
                name="domain-notes"
                id="domain-notes"
                class="block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-gray-400 focus:ring-gray-400 focus:ring-offset-gray-900 focus:ring-offset-0 sm:text-sm p-2"
                placeholder="Add notes about this domain..."
              ></textarea>
              <!-- Read-only view for viewers -->
              <div
                v-else
                class="block w-full rounded-md bg-gray-700 border-gray-600 text-white p-2 min-h-[96px]"
              >
                {{ notes || 'No notes available' }}
              </div>
            </div>

            <!-- Update Button Section - Only show if user can edit -->
            <div v-if="canEdit" class="flex flex-col items-center">
              <button
                @click="updateDomain"
                class="w-full md:w-auto px-6 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-gray-600 hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-400"
              >
                Update Domain
              </button>
              <p v-if="updateMessage" class="mt-2 text-sm" :class="{
                'text-green-400': updateMessage.includes('success'),
                'text-red-400': updateMessage.includes('Failed')
              }">
                {{ updateMessage }}
              </p>
            </div>
          </div>
        </div>

        <!-- DNS Records Section -->
        <div class="bg-gray-800 shadow rounded-lg mt-6">
          <div class="px-4 py-5 sm:p-6">
            <h2 class="text-lg font-medium text-white mb-4">DNS Records</h2>
            <div class="overflow-x-auto ring-1 ring-gray-700 rounded-lg">
              <table class="min-w-full divide-y divide-gray-700">
                <thead class="bg-gray-700">
                  <tr>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider first:rounded-tl-lg">
                      Type
                    </th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">
                      Host
                    </th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider last:rounded-tr-lg">
                      Value
                    </th>
                  </tr>
                </thead>
                <tbody class="bg-gray-800 divide-y divide-gray-700">
                  <tr v-for="record in dnsRecords" :key="record.id" class="hover:bg-gray-700">
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-white">{{ record.Record_Type }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-white">{{ record.Record_Name }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-white">{{ record.Address }}</td>
                  </tr>
                  <tr v-if="dnsRecords.length === 0">
                    <td colspan="3" class="px-6 py-4 text-center text-sm text-gray-400">
                      No DNS records found
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