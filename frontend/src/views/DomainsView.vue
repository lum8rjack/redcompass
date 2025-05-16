<script setup>
  import Header from '@/components/Header.vue'
  import Footer from '@/components/Footer.vue'
  import { inject, onMounted, ref, computed } from 'vue'
  import { onClickOutside } from '@vueuse/core'
  import { formatDate } from '@/utils/dateUtils'
  import PageDescription from '@/components/PageDescription.vue'

  const pocketbase = inject('$pocketbase');
  const domains = ref([]);

  onMounted(async () => {
    try {
      const response = await pocketbase.collection('Domains').getFullList({
        sort: '-Name',
        fields: 'id,Name,Purchased_Date,Expiration_Date,Is_Expired,Auto_Renew,Is_Locked,Assigned_Project,Tags,Healthy,Checked_Out_By,expand.Assigned_Project.Name',
        expand: 'Assigned_Project',
      });
      domains.value = response;
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  });


// Filter states
const filters = ref({
  name: '',
  healthStatus: '',
  assignedTo: '',
  tags: ''
})

// Filter options
const healthStatusOptions = ['All', 'Healthy', 'Unhealthy']

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

// Update computed filtered domains to include sorting
const filteredDomains = computed(() => {
  let filtered = domains.value.filter(domain => {
    const nameMatch = domain.Name.toLowerCase().includes(filters.value.name.toLowerCase())
    const healthMatch = !filters.value.healthStatus || 
      filters.value.healthStatus === 'All' || 
      (filters.value.healthStatus === 'Healthy' && domain.Healthy === true) ||
      (filters.value.healthStatus === 'Unhealthy' && domain.Healthy === false)
    const assignedMatch = !filters.value.assignedTo || 
      (domain.expand?.Assigned_Project?.Name?.toLowerCase() || '').includes(filters.value.assignedTo.toLowerCase())
    const tagsMatch = !filters.value.tags || 
      (domain.Tags && domain.Tags.some(tag => tag.toLowerCase().includes(filters.value.tags.toLowerCase())))

    return nameMatch && healthMatch && assignedMatch && tagsMatch
  })

  // Apply sorting
  return filtered.sort((a, b) => {
    let aValue = a[sortConfig.value.key]
    let bValue = b[sortConfig.value.key]

    // Handle null values for assignedTo
    if (sortConfig.value.key === 'assignedTo') {
      aValue = aValue || ''
      bValue = bValue || ''
    }

    // Handle boolean values for autoRenew
    if (sortConfig.value.key === 'autoRenew') {
      return sortConfig.value.direction === 'asc' ? 
        Number(aValue) - Number(bValue) : 
        Number(bValue) - Number(aValue)
    }

    // String comparison for other fields
    if (sortConfig.value.direction === 'asc') {
      return aValue > bValue ? 1 : -1
    }
    return aValue < bValue ? 1 : -1
  })
})

// Add new refs for dropdown
const showHealthDropdown = ref(false)
const healthContainer = ref(null)

// Close dropdown when clicking outside
onClickOutside(healthContainer, () => {
  showHealthDropdown.value = false
})
</script>

<template>
  <div class="min-h-screen flex flex-col bg-gray-900">
    <Header />
    <main class="flex-grow">
      <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <PageDescription title="Domain Management" description="Manage your domains, track their health status, and assign them to projects. Monitor expiration dates, auto-renewal settings, and domain health from a centralized dashboard." />

        <!-- Filter Section -->
        <div class="bg-gray-800 shadow rounded-lg mb-6">
          <div class="px-4 py-5 sm:p-6">
            <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
              <div>
                <label for="nameFilter" class="block text-sm font-medium text-white">Filter by Domain</label>
                <input
                  type="text"
                  id="nameFilter"
                  v-model="filters.name"
                  placeholder="Filter by domain name..."
                  class="filter-input"
                >
              </div>
              <div>
                <label for="assignedFilter" class="block text-sm font-medium text-white">Filter by Project</label>
                <input
                  type="text"
                  id="assignedFilter"
                  v-model="filters.assignedTo"
                  placeholder="Search assigned project..."
                  class="filter-input"
                >
              </div>
              <div>
                <label for="tagsFilter" class="block text-sm font-medium text-white">Filter by Tags</label>
                <input
                  type="text"
                  id="tagsFilter"
                  v-model="filters.tags"
                  placeholder="Search tags..."
                  class="filter-input"
                >
              </div>
              <div>
                <label for="healthFilter" class="block text-sm font-medium text-white">Filter by Health</label>
                <div class="mt-1 relative" ref="healthContainer" id="healthContainer">
                  <!-- Dropdown trigger button -->
                  <button
                    type="button"
                    id="healthFilter"
                    @click="showHealthDropdown = !showHealthDropdown"
                    class="filter-input text-left"
                  >
                    {{ filters.healthStatus || 'Select status' }}
                    <span class="absolute inset-y-0 right-0 flex items-center pr-2">
                      <svg class="h-5 w-5 text-gray-400" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
                      </svg>
                    </span>
                  </button>

                  <!-- Dropdown menu -->
                  <div 
                    v-if="showHealthDropdown"
                    class="absolute z-10 mt-1 w-full bg-gray-700 rounded-md shadow-lg"
                  >
                    <ul class="max-h-60 rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none sm:text-sm">
                      <li
                        v-for="status in healthStatusOptions"
                        :key="status"
                        @click="() => { 
                          filters.healthStatus = status;
                          showHealthDropdown = false;
                        }"
                        class="text-white cursor-pointer select-none relative py-2 pl-3 pr-9 hover:bg-gray-600"
                        :class="{ 'bg-gray-600': filters.healthStatus === status }"
                      >
                        {{ status }}
                      </li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Results count -->
        <p class="text-sm text-gray-500 mb-4">
          Showing {{ filteredDomains.length }} of {{ domains.length }} domains
        </p>

        <div class="bg-gray-800 shadow rounded-lg">
          <div class="px-4 py-5 sm:p-6">
            <div class="overflow-hidden rounded-lg">
              <table class="min-w-full divide-y divide-gray-700">
                <thead class="bg-gray-700">
                  <tr>
                    <th 
                      v-for="(header, key) in {
                        Name: 'Domain Name',
                        Purchased_Date: 'Purchase Date',
                        Expiration_Date: 'Expiration Date',
                        Auto_Renew: 'Auto Renew',
                        Is_Locked: 'Locked',
                        Assigned_Project: 'Assigned to',
                        Tags: 'Tags',
                        Healthy: 'Health Status'
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
                    v-for="domain in filteredDomains" 
                    :key="domain.id"
                    @click="$router.push(`/domain/${domain.id}`)"
                    class="cursor-pointer hover:bg-gray-700 transition-colors duration-150"
                  >
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-white">
                      {{ domain.Name }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-white">
                      {{ formatDate(domain.Purchased_Date) }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-white">
                      <span :class="domain.Is_Expired ? 'text-red-600' : 'text-white'">
                        {{ formatDate(domain.Expiration_Date) }}
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      <span :class="domain.Auto_Renew ? 'text-green-600' : 'text-red-600'">
                        {{ domain.Auto_Renew ? 'On' : 'Off' }}
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      <span :class="domain.Is_Locked ? 'text-red-600' : 'text-green-600'">
                        {{ domain.Is_Locked ? 'Yes' : 'No' }}
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      <span :class="domain.expand?.Assigned_Project?.Name ? 'text-white' : 'text-gray-500'">
                        {{ domain.expand?.Assigned_Project?.Name || 'Unassigned' }}
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-white">
                      {{ domain.Tags ? domain.Tags.join(', ') : '-' }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm">
                      <span :class="{
                        'px-2 py-1 rounded-full text-xs font-medium': true,
                        'bg-green-100 text-green-800': domain.Healthy === true,
                        'bg-red-100 text-red-800': domain.Healthy === false
                      }">
                        {{ domain.Healthy ? 'Healthy' : 'Unhealthy' }}
                      </span>
                    </td>
                  </tr>
                  <tr class="last:rounded-b-lg">
                    <td class="first:rounded-bl-lg"></td>
                    <td></td>
                    <td></td>
                    <td></td>
                    <td></td>
                    <td></td>
                    <td class="last:rounded-br-lg"></td>
                  </tr>
                </tbody>
              </table>

              <!-- No results message -->
              <div v-if="filteredDomains.length === 0" class="text-center py-4 text-gray-500">
                No domains match your filters
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
/* Common filter input styles */
.filter-input {
  margin-top: 0.25rem;
  display: block;
  width: 100%;
  border-radius: 0.375rem;
  border: 1px solid rgb(75, 85, 99);
  background-color: rgb(55, 65, 81);
  color: white;
  padding: 0.375rem 0.5rem;
  padding-right: 2.5rem;
  font-size: 0.875rem;
  line-height: 1.25rem;
}

.filter-input:focus {
  outline: none;
  border-color: white;
  box-shadow: 0 0 0 1px white;
}

/* Ensure dropdown appears above other elements */
.relative {
  z-index: 50;
}

/* Smooth transitions */
.absolute {
  transition: all 0.2s;
}

/* Add smooth transition for sort indicators */
th {
  transition: background-color 0.2s;
}

/* Style the dropdown options */
select option {
  padding: 0.5rem;
  margin: 0.25rem;
  border-radius: 0.375rem;
}

select option:checked {
  background: linear-gradient(0deg, #4f46e5 0%, #4f46e5 100%);
}

select option:hover {
  background-color: rgb(75, 85, 99);
}

/* Ensure consistent input heights */
input #healthContainer {
  height: 38px;
}
</style>
