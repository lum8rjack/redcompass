<script setup>
import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'
import { ref, inject, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { canEdit } from '@/utils/auth'

const pocketbase = inject('$pocketbase')
const route = useRoute()
const router = useRouter()

const campaign = ref(null)
const loading = ref(true)
const error = ref(null)
const isEditingCampaignInfo = ref(false)
const isEditingMetrics = ref(false)
const isEditingPhishlet = ref(false)
const isEditingCaddyfile = ref(false)
const editedCampaignName = ref('')
const editedTargetGroup = ref('')
const editedClickRate = ref('')
const editedCredentials = ref('')
const editedPhishlet = ref('')
const editedCaddyfile = ref('')
const campaignNotes = ref('')
const hasUnsavedChanges = ref(false)
const showCloseConfirmation = ref(false)
const originalCampaign = ref(null)
const clipboardContent = ref('')
const showClipboardNotification = ref(false)

onMounted(async () => {
  try {
    const campaignId = route.params.id
    // Fetch campaign from backend (replace with real API call)
    // For now, simulate fetch
    // TODO: Replace with real fetch logic
    // Example:
    // campaign.value = await pocketbase.collection('PhishingCampaigns').getOne(campaignId)
    // For now, just simulate:
    campaign.value = {
      id: campaignId,
      name: 'Sample Campaign',
      target: 'Marketing Team',
      status: 'Completed',
      clickRate: '23%',
      submittedCredentials: '12%',
      htmlTemplate: '',
      yamlConfig: '',
      caddyfile: '',
      notes: 'Sample notes.'
    }
    editedCampaignName.value = campaign.value.name
    editedTargetGroup.value = campaign.value.target
    editedClickRate.value = campaign.value.clickRate
    editedCredentials.value = campaign.value.submittedCredentials
    editedPhishlet.value = campaign.value.yamlConfig
    editedCaddyfile.value = campaign.value.caddyfile
    campaignNotes.value = campaign.value.notes
    originalCampaign.value = JSON.parse(JSON.stringify(campaign.value))
  } catch (err) {
    error.value = 'Failed to load campaign details'
  } finally {
    loading.value = false
  }
})

function updateCampaignInfo() {
  if (!campaign.value) return
  campaign.value.name = editedCampaignName.value
  campaign.value.target = editedTargetGroup.value
  isEditingCampaignInfo.value = false
}

function updateCampaignMetrics() {
  if (!campaign.value) return
  campaign.value.clickRate = editedClickRate.value
  campaign.value.submittedCredentials = editedCredentials.value
  isEditingMetrics.value = false
}

function toggleEditPhishlet() {
  isEditingPhishlet.value = !isEditingPhishlet.value
  if (isEditingPhishlet.value) {
    editedPhishlet.value = campaign.value.yamlConfig
  }
}

function savePhishlet() {
  if (!campaign.value) return
  campaign.value.yamlConfig = editedPhishlet.value
  isEditingPhishlet.value = false
}

function toggleEditCaddyfile() {
  isEditingCaddyfile.value = !isEditingCaddyfile.value
  if (isEditingCaddyfile.value) {
    editedCaddyfile.value = campaign.value.caddyfile
  }
}

function saveCaddyfile() {
  if (!campaign.value) return
  campaign.value.caddyfile = editedCaddyfile.value
  isEditingCaddyfile.value = false
}

function saveNotes() {
  if (!campaign.value) return
  campaign.value.notes = campaignNotes.value
}

function copyYamlToClipboard() {
  if (!campaign.value) return
  navigator.clipboard.writeText(campaign.value.yamlConfig)
  clipboardContent.value = 'Phishlet configuration'
  showClipboardNotification.value = true
  setTimeout(() => { showClipboardNotification.value = false }, 2000)
}

function copyCaddyToClipboard() {
  if (!campaign.value) return
  navigator.clipboard.writeText(campaign.value.caddyfile)
  clipboardContent.value = 'Caddyfile'
  showClipboardNotification.value = true
  setTimeout(() => { showClipboardNotification.value = false }, 2000)
}

function downloadYamlConfig() {
  if (!campaign.value) return
  const blob = new Blob([campaign.value.yamlConfig], { type: 'text/yaml' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${campaign.value.name.replace(/\s+/g, '_')}_config.yaml`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

function downloadCaddyfile() {
  if (!campaign.value) return
  const blob = new Blob([campaign.value.caddyfile], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${campaign.value.name.replace(/\s+/g, '_')}_caddyfile`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

function copyHtmlToClipboard() {
  if (!campaign.value) return
  navigator.clipboard.writeText(campaign.value.htmlTemplate)
  clipboardContent.value = 'HTML template'
  showClipboardNotification.value = true
  setTimeout(() => { showClipboardNotification.value = false }, 2000)
}

const sanitizedHtml = computed(() => {
  if (!campaign.value) return ''
  return campaign.value.htmlTemplate
})
</script>

<template>
  <div class="min-h-screen flex flex-col bg-gray-900">
    <Header />
    <main class="flex-grow">
      <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div v-if="loading" class="text-white text-center py-8">Loading campaign details...</div>
        <div v-else-if="error" class="text-red-500 text-center py-8">{{ error }}</div>
        <div v-else-if="campaign" class="bg-gray-800 shadow rounded-lg">
          <div class="px-4 py-5 sm:p-6">
            <div class="flex justify-between items-start mb-6">
              <div>
                <h1 class="text-2xl font-bold text-white">{{ campaign.name }}</h1>
                <p class="text-gray-400 mt-2">Target: {{ campaign.target }}</p>
                <p class="text-gray-400 mt-1">Status: {{ campaign.status }}</p>
              </div>
              <button @click="router.back()" class="px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-400">Back</button>
            </div>
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <!-- Campaign Info -->
              <div>
                <div class="bg-gray-700 rounded-lg p-4 mb-4">
                  <div class="flex justify-between items-center mb-2">
                    <h4 class="text-md font-medium text-white">Campaign Information</h4>
                    <div>
                      <button 
                        @click="isEditingCampaignInfo = !isEditingCampaignInfo" 
                        class="text-xs font-medium text-blue-400 hover:text-blue-300 mr-2"
                      >
                        {{ isEditingCampaignInfo ? 'Cancel' : 'Edit Campaign' }}
                      </button>
                      <button 
                        @click="isEditingMetrics = !isEditingMetrics" 
                        class="text-xs font-medium text-blue-400 hover:text-blue-300"
                      >
                        {{ isEditingMetrics ? 'Cancel' : 'Edit Metrics' }}
                      </button>
                    </div>
                  </div>
                  <div class="grid grid-cols-2 gap-2">
                    <div class="text-sm text-gray-400">Campaign Name:</div>
                    <div v-if="!isEditingCampaignInfo" class="text-sm text-white">
                      {{ campaign.name }}
                    </div>
                    <div v-else class="text-sm">
                      <input 
                        v-model="editedCampaignName" 
                        type="text" 
                        class="w-full rounded-md bg-gray-600 border-gray-500 text-white text-sm py-1 px-2"
                        placeholder="Campaign Name"
                      />
                    </div>
                    <div class="text-sm text-gray-400">Target Group:</div>
                    <div v-if="!isEditingCampaignInfo" class="text-sm text-white">
                      {{ campaign.target }}
                    </div>
                    <div v-else class="text-sm">
                      <input 
                        v-model="editedTargetGroup" 
                        type="text" 
                        class="w-full rounded-md bg-gray-600 border-gray-500 text-white text-sm py-1 px-2"
                        placeholder="Target Group"
                      />
                    </div>
                    <div class="text-sm text-gray-400">Status:</div>
                    <div class="text-sm text-white">
                      <span :class="{
                        'px-2 py-1 rounded-full text-xs font-medium': true,
                        'bg-green-100 text-green-800': campaign.status === 'Completed',
                        'bg-yellow-100 text-yellow-800': campaign.status === 'In Progress',
                        'bg-blue-100 text-blue-800': campaign.status === 'Scheduled'
                      }">
                        {{ campaign.status }}
                      </span>
                    </div>
                    <div class="text-sm text-gray-400">Click Rate:</div>
                    <div v-if="!isEditingMetrics" class="text-sm text-white">
                      {{ campaign.clickRate }}
                    </div>
                    <div v-else class="text-sm">
                      <input 
                        v-model="editedClickRate" 
                        type="text" 
                        class="w-full rounded-md bg-gray-600 border-gray-500 text-white text-sm py-1 px-2"
                        placeholder="e.g., 15%"
                      />
                    </div>
                    <div class="text-sm text-gray-400">Submitted Credentials:</div>
                    <div v-if="!isEditingMetrics" class="text-sm text-white">
                      {{ campaign.submittedCredentials }}
                    </div>
                    <div v-else class="text-sm">
                      <input 
                        v-model="editedCredentials" 
                        type="text" 
                        class="w-full rounded-md bg-gray-600 border-gray-500 text-white text-sm py-1 px-2"
                        placeholder="e.g., 7%"
                      />
                    </div>
                  </div>
                  <div v-if="isEditingCampaignInfo" class="mt-3 flex justify-end">
                    <button
                      @click="updateCampaignInfo"
                      class="bg-blue-600 hover:bg-blue-700 text-white text-xs font-medium py-1 px-3 rounded"
                    >
                      Save Campaign Changes
                    </button>
                  </div>
                  <div v-if="isEditingMetrics" class="mt-3 flex justify-end">
                    <button
                      @click="updateCampaignMetrics"
                      class="bg-blue-600 hover:bg-blue-700 text-white text-xs font-medium py-1 px-3 rounded"
                    >
                      Save Metrics Changes
                    </button>
                  </div>
                </div>
                <!-- YAML Configuration -->
                <div>
                  <div class="flex justify-between items-center mb-2">
                    <h4 class="text-md font-medium text-white">Phishlet Configuration</h4>
                    <div class="flex space-x-2">
                      <button 
                        v-if="!isEditingPhishlet"
                        @click="toggleEditPhishlet" 
                        class="bg-gray-600 hover:bg-gray-500 text-white text-xs font-medium py-1 px-2 rounded"
                      >
                        Edit
                      </button>
                      <button 
                        v-if="isEditingPhishlet"
                        @click="savePhishlet" 
                        class="bg-green-600 hover:bg-green-700 text-white text-xs font-medium py-1 px-2 rounded"
                      >
                        Save
                      </button>
                      <button 
                        v-if="isEditingPhishlet"
                        @click="toggleEditPhishlet" 
                        class="bg-red-600 hover:bg-red-700 text-white text-xs font-medium py-1 px-2 rounded"
                      >
                        Cancel
                      </button>
                      <button 
                        v-if="!isEditingPhishlet"
                        @click="copyYamlToClipboard" 
                        class="bg-gray-600 hover:bg-gray-500 text-white text-xs font-medium py-1 px-2 rounded"
                      >
                        Copy to Clipboard
                      </button>
                      <button 
                        v-if="!isEditingPhishlet"
                        @click="downloadYamlConfig" 
                        class="bg-blue-600 hover:bg-blue-700 text-white text-xs font-medium py-1 px-2 rounded"
                      >
                        Download YAML
                      </button>
                    </div>
                  </div>
                  <div v-if="!isEditingPhishlet" class="bg-gray-900 rounded-lg p-4 max-h-80 overflow-y-auto">
                    <pre class="text-sm text-gray-300 whitespace-pre-wrap font-mono">{{ campaign.yamlConfig }}</pre>
                  </div>
                  <div v-else class="bg-gray-900 rounded-lg p-0 max-h-80 overflow-y-auto">
                    <textarea
                      v-model="editedPhishlet"
                      class="w-full bg-gray-800 border-gray-700 rounded-md text-white p-3 font-mono text-sm"
                      rows="12"
                      style="resize: vertical;"
                    ></textarea>
                  </div>
                </div>
                <!-- Caddyfile Section -->
                <div class="mt-4">
                  <div class="flex justify-between items-center mb-2">
                    <h4 class="text-md font-medium text-white">Caddyfile</h4>
                    <div class="flex space-x-2">
                      <button 
                        v-if="!isEditingCaddyfile"
                        @click="toggleEditCaddyfile" 
                        class="bg-gray-600 hover:bg-gray-500 text-white text-xs font-medium py-1 px-2 rounded"
                      >
                        Edit
                      </button>
                      <button 
                        v-if="isEditingCaddyfile"
                        @click="saveCaddyfile" 
                        class="bg-green-600 hover:bg-green-700 text-white text-xs font-medium py-1 px-2 rounded"
                      >
                        Save
                      </button>
                      <button 
                        v-if="isEditingCaddyfile"
                        @click="toggleEditCaddyfile" 
                        class="bg-red-600 hover:bg-red-700 text-white text-xs font-medium py-1 px-2 rounded"
                      >
                        Cancel
                      </button>
                      <button 
                        v-if="!isEditingCaddyfile"
                        @click="copyCaddyToClipboard" 
                        class="bg-gray-600 hover:bg-gray-500 text-white text-xs font-medium py-1 px-2 rounded"
                      >
                        Copy to Clipboard
                      </button>
                      <button 
                        v-if="!isEditingCaddyfile"
                        @click="downloadCaddyfile" 
                        class="bg-blue-600 hover:bg-blue-700 text-white text-xs font-medium py-1 px-2 rounded"
                      >
                        Download Caddyfile
                      </button>
                    </div>
                  </div>
                  <div v-if="!isEditingCaddyfile" class="bg-gray-900 rounded-lg p-4 max-h-80 overflow-y-auto">
                    <pre class="text-sm text-gray-300 whitespace-pre-wrap font-mono">{{ campaign.caddyfile }}</pre>
                  </div>
                  <div v-else class="bg-gray-900 rounded-lg p-0 max-h-80 overflow-y-auto">
                    <textarea
                      v-model="editedCaddyfile"
                      class="w-full bg-gray-800 border-gray-700 rounded-md text-white p-3 font-mono text-sm"
                      rows="12"
                      style="resize: vertical;"
                    ></textarea>
                  </div>
                </div>
              </div>
              <!-- HTML Email Preview and Notes -->
              <div>
                <div class="flex justify-between items-center mb-2">
                  <h4 class="text-md font-medium text-white">HTML Email Preview</h4>
                  <button 
                    @click="copyHtmlToClipboard" 
                    class="bg-blue-600 hover:bg-blue-700 text-white text-xs font-medium py-1 px-3 rounded"
                  >
                    Copy HTML
                  </button>
                </div>
                <div class="bg-white rounded-lg overflow-hidden" style="height: 400px;">
                  <iframe
                    sandbox="allow-same-origin"
                    class="w-full h-full"
                    title="Email Template Preview"
                    :srcdoc="sanitizedHtml"
                  ></iframe>
                </div>
                <!-- Campaign Notes -->
                <div class="mt-4">
                  <div class="flex justify-between items-center mb-2">
                    <h4 class="text-md font-medium text-white">Campaign Notes</h4>
                    <button 
                      @click="saveNotes" 
                      class="bg-blue-600 hover:bg-blue-700 text-white text-xs font-medium py-1 px-3 rounded"
                    >
                      Save Notes
                    </button>
                  </div>
                  <textarea
                    v-model="campaignNotes"
                    @input="hasUnsavedChanges = true"
                    class="w-full bg-gray-700 border-gray-600 rounded-md text-white p-3"
                    style="height: 400px; resize: none;"
                    placeholder="Add campaign notes, observations, or action items here..."
                  ></textarea>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
    <Footer />
    <div v-if="showClipboardNotification" class="fixed bottom-4 right-4 bg-gray-900 text-white px-4 py-2 rounded shadow-lg z-50">
      {{ clipboardContent }} copied to clipboard
    </div>
  </div>
</template>

<style scoped>
pre {
  word-wrap: break-word;
}
</style> 