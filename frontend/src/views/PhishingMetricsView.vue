<script setup>
import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'
import PageDescription from '@/components/PageDescription.vue'
import { ref, inject, onMounted, computed } from 'vue'
import { calculatePercentage } from '@/utils/calculateStats'

const pocketbase = inject('$pocketbase')
const loading = ref(true)
const error = ref(null)
const templates = ref([])
const campaigns = ref([])
const selectedTemplate = ref('all')
const sortColumn = ref('Date_Sent')
const sortDirection = ref('desc')

// Fetch all phishing templates for dropdown
async function getTemplates() {
  try {
    const response = await pocketbase.collection('Phishing_Templates').getFullList({
      fields: 'id,Name',
      sort: 'Name',
    })
    templates.value = response
  } catch (err) {
    error.value = 'Failed to load templates'
  }
}

// Fetch all phishing campaigns (metrics)
async function getCampaigns() {
  try {
    const response = await pocketbase.collection('Phishing_Metrics').getFullList({
      expand: 'Phishing_Template,Project',
      sort: '-Date_Sent',
      fields: 'id,Phishing_Template,expand.Phishing_Template.Name,expand.Project.Name,Date_Sent,Subject,From,Emails_Sent,Emails_Clicked,Creds_Submit',
    })
    campaigns.value = response
  } catch (err) {
    error.value = 'Failed to load campaigns'
  }
}

onMounted(async () => {
  loading.value = true
  await Promise.all([getTemplates(), getCampaigns()])
  loading.value = false
})

const filteredCampaigns = computed(() => {
  let data = selectedTemplate.value === 'all' ? campaigns.value : campaigns.value.filter(c => c.Phishing_Template === selectedTemplate.value)
  if (!sortColumn.value) return data
  return [...data].sort((a, b) => {
    let aVal, bVal
    switch (sortColumn.value) {
      case 'Template':
        aVal = a.expand?.Phishing_Template?.Name || ''
        bVal = b.expand?.Phishing_Template?.Name || ''
        break
      case 'Project':
        aVal = a.expand?.Project?.Name || ''
        bVal = b.expand?.Project?.Name || ''
        break
      case 'Date_Sent':
        aVal = a.Date_Sent || ''
        bVal = b.Date_Sent || ''
        break
      case 'Subject':
        aVal = a.Subject || ''
        bVal = b.Subject || ''
        break
      case 'From':
        aVal = a.From || ''
        bVal = b.From || ''
        break
      case 'Emails_Sent':
        aVal = a.Emails_Sent || 0
        bVal = b.Emails_Sent || 0
        break
      case 'Emails_Clicked':
        aVal = a.Emails_Clicked || 0
        bVal = b.Emails_Clicked || 0
        break
      case 'Creds_Submit':
        aVal = a.Creds_Submit || 0
        bVal = b.Creds_Submit || 0
        break
      default:
        aVal = ''
        bVal = ''
    }
    if (aVal < bVal) return sortDirection.value === 'asc' ? -1 : 1
    if (aVal > bVal) return sortDirection.value === 'asc' ? 1 : -1
    return 0
  })
})

const totalEmailsSent = computed(() => filteredCampaigns.value.reduce((sum, c) => sum + (c.Emails_Sent || 0), 0))
const totalEmailsClicked = computed(() => filteredCampaigns.value.reduce((sum, c) => sum + (c.Emails_Clicked || 0), 0))
const totalCredsSubmit = computed(() => filteredCampaigns.value.reduce((sum, c) => sum + (c.Creds_Submit || 0), 0))

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString()
}

function handleSort(col) {
  if (sortColumn.value === col) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortColumn.value = col
    sortDirection.value = 'asc'
  }
}
</script>

<template>
  <div class="min-h-screen flex flex-col bg-gray-900">
    <Header />
    <main class="flex-grow">
      <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <PageDescription title="Phishing Metrics" description="View and filter all phishing campaign metrics by template." />
        <div v-if="loading" class="text-white text-center py-8">Loading metrics...</div>
        <div v-else-if="error" class="text-red-500 text-center py-8">{{ error }}</div>
        <div v-else>
          <div class="mb-6 flex flex-col sm:flex-row sm:items-center gap-4">
            <label for="template" class="block text-sm font-medium text-white">Phishing Template</label>
            <select id="template" v-model="selectedTemplate" class="rounded-md bg-gray-700 border-gray-600 text-white py-1.5 pl-2">
              <option value="all">All Templates</option>
              <option v-for="t in templates" :key="t.id" :value="t.id">{{ t.Name }}</option>
            </select>
          </div>
          <div class="bg-gray-800 shadow rounded-lg overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-700">
              <thead class="bg-gray-700">
                <tr>
                  <th @click="handleSort('Template')" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider cursor-pointer select-none">
                    Template
                    <span v-if="sortColumn === 'Template'">{{ sortDirection === 'asc' ? '▲' : '▼' }}</span>
                  </th>
                  <th @click="handleSort('Project')" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider cursor-pointer select-none">
                    Project
                    <span v-if="sortColumn === 'Project'">{{ sortDirection === 'asc' ? '▲' : '▼' }}</span>
                  </th>
                  <th @click="handleSort('Date_Sent')" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider cursor-pointer select-none">
                    Date Sent
                    <span v-if="sortColumn === 'Date_Sent'">{{ sortDirection === 'asc' ? '▲' : '▼' }}</span>
                  </th>
                  <th @click="handleSort('Subject')" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider cursor-pointer select-none">
                    Subject
                    <span v-if="sortColumn === 'Subject'">{{ sortDirection === 'asc' ? '▲' : '▼' }}</span>
                  </th>
                  <th @click="handleSort('From')" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider cursor-pointer select-none">
                    From
                    <span v-if="sortColumn === 'From'">{{ sortDirection === 'asc' ? '▲' : '▼' }}</span>
                  </th>
                  <th @click="handleSort('Emails_Sent')" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider cursor-pointer select-none">
                    Emails Sent
                    <span v-if="sortColumn === 'Emails_Sent'">{{ sortDirection === 'asc' ? '▲' : '▼' }}</span>
                  </th>
                  <th @click="handleSort('Emails_Clicked')" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider cursor-pointer select-none">
                    Clicks
                    <span v-if="sortColumn === 'Emails_Clicked'">{{ sortDirection === 'asc' ? '▲' : '▼' }}</span>
                  </th>
                  <th @click="handleSort('Creds_Submit')" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider cursor-pointer select-none">
                    Creds
                    <span v-if="sortColumn === 'Creds_Submit'">{{ sortDirection === 'asc' ? '▲' : '▼' }}</span>
                  </th>
                </tr>
              </thead>
              <tbody class="bg-gray-800 divide-y divide-gray-700">
                <tr v-for="c in filteredCampaigns" :key="c.id">
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-white">{{ c.expand?.Phishing_Template?.Name || 'Unknown' }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-white">{{ c.expand?.Project?.Name || 'Unknown' }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{{ formatDate(c.Date_Sent) }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{{ c.Subject }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{{ c.From }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{{ c.Emails_Sent }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{{ c.Emails_Clicked || 0 }} ({{ calculatePercentage(c.Emails_Sent, c.Emails_Clicked) }}%)</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{{ c.Creds_Submit || 0 }} ({{ calculatePercentage(c.Emails_Sent, c.Creds_Submit) }}%)</td>
                </tr>
                <tr v-if="filteredCampaigns.length === 0">
                  <td colspan="8" class="px-6 py-4 text-center text-sm text-gray-400">No phishing campaigns found</td>
                </tr>
                <tr v-else class="bg-gray-700 font-bold">
                  <td class="px-6 py-4 text-sm text-white" colspan="5">Totals</td>
                  <td class="px-6 py-4 text-sm text-white">{{ totalEmailsSent }}</td>
                  <td class="px-6 py-4 text-sm text-white">{{ totalEmailsClicked }} ({{ calculatePercentage(totalEmailsSent, totalEmailsClicked) }}%)</td>
                  <td class="px-6 py-4 text-sm text-white">{{ totalCredsSubmit }} ({{ calculatePercentage(totalEmailsSent, totalCredsSubmit) }}%)</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </main>
    <Footer />
  </div>
</template>

<style scoped>
</style> 