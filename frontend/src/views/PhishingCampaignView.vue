<script setup>
import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'
import { ref, inject, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { canEdit } from '@/utils/auth'
import { calculatePercentage } from '@/utils/calculateStats'

const pocketbase = inject('$pocketbase')
const route = useRoute()

const campaign = ref(null)
const loading = ref(true)
const error = ref(null)
const isEditingCampaignInfo = ref(false)
const isEditingCaddyfile = ref(false)
const editedCampaignName = ref('')
const editedTargetGroup = ref('')
const editedCaddyfile = ref('')
const campaignNotes = ref('')
const originalCampaign = ref(null)
const editedCampaignNotes = ref('')
const isEditingHtml = ref(false)
const editedHtml = ref('')
const editedSubject = ref('')
const editedFrom = ref('')
const campaignMessage = ref('')
const htmlMessage = ref('')
const phishletMessage = ref('')
const caddyfileMessage = ref('')
onMounted(async () => {
  try {
    const campaignId = route.params.id;
    
    // Fetch campaign from Pocketbase
    const record = await pocketbase.collection('Phishing_Campaign').getOne(campaignId, {
      fields: 'Name, Notes, Caddy, Example_From, Example_Subject, Target_Group, expand.Phishlet, HTML',
      expand: 'Phishlet'
    });
    console.log("record", record)

    const stats = await pocketbase.collection('Phishing_Campaign_View').getOne(campaignId, {
        fields: 'id, Name, Target_Group, total_sent, total_clicked, total_submit',
    });
    console.log("stats", stats)

    let tempPhishlet = record.expand.Phishlet || {}
    // Map backend fields to local campaign object
    campaign.value = {
      id: campaignId,
      name: record.Name,
      from: record.Example_From,
      subject: record.Example_Subject,
      target: record.Target_Group,
      htmlTemplate: record.HTML || '',
      yamlConfig: tempPhishlet.Phishlet || '',
      phishletName: tempPhishlet.Name || '',
      caddyfile: record.Caddy || '',
      notes: record.Notes || '',
      emailsSent: stats.total_sent || 0,
      emailsClicked: stats.total_clicked || 0,
      credentialsSubmitted: stats.total_submit || 0
    };

    editedCampaignName.value = campaign.value.name;
    editedTargetGroup.value = campaign.value.target;
    editedCaddyfile.value = campaign.value.caddyfile;
    campaignNotes.value = campaign.value.notes;
    editedCampaignNotes.value = campaign.value.notes;
    editedSubject.value = campaign.value.subject;
    editedFrom.value = campaign.value.from;
    originalCampaign.value = JSON.parse(JSON.stringify(campaign.value));
  } catch (err) {
    error.value = 'Failed to load campaign details';
    console.error(err);
  } finally {
    loading.value = false;
  }
});

async function updateCampaignInfo() {
  if (!campaign.value) return
  try {
    await pocketbase.collection('Phishing_Campaign').update(campaign.value.id, {
      Name: editedCampaignName.value,
      Target_Group: editedTargetGroup.value,
      Notes: editedCampaignNotes.value,
      Example_Subject: editedSubject.value,
      Example_From: editedFrom.value
    });

    campaign.value.name = editedCampaignName.value
    campaign.value.target = editedTargetGroup.value
    campaign.value.notes = editedCampaignNotes.value
    campaign.value.subject = editedSubject.value
    campaign.value.from = editedFrom.value
    campaignMessage.value = 'Successfully updated campaign info'
    setTimeout(() => {
      campaignMessage.value = ''
    }, 2000)
  } catch(err) {
    console.error('Failed to update campaign info:', err);
    campaignMessage.value = 'Error updating campaign info'
    setTimeout(() => {
      campaignMessage.value = ''
    }, 2000)
  } finally {
    isEditingCampaignInfo.value = false
  }
}

function toggleEditCaddyfile() {
  isEditingCaddyfile.value = !isEditingCaddyfile.value
  if (isEditingCaddyfile.value) {
    editedCaddyfile.value = campaign.value.caddyfile
  }
}

async function updateCaddyfile(newCaddyfile) {
  if (!campaign.value) return
  try {
    await pocketbase.collection('Phishing_Campaign').update(campaign.value.id, {
      Caddy: newCaddyfile
    });
  } catch(err) {
    throw err
  }
}

async function saveCaddyfile() {
  if (!campaign.value) return
  try {
    await updateCaddyfile(editedCaddyfile.value)
    campaign.value.caddyfile = editedCaddyfile.value
    caddyfileMessage.value = 'Successfully updated Caddyfile'
    setTimeout(() => {
      caddyfileMessage.value = ''
    }, 2000)
  } catch(err) {
    console.error('Failed to save Caddyfile:', err);
    caddyfileMessage.value = 'Error updating Caddyfile'
    setTimeout(() => {
      caddyfileMessage.value = ''
    }, 2000)
  } finally {
    isEditingCaddyfile.value = false
  }
}

function cancelEditCampaignInfo() {
  isEditingCampaignInfo.value = false;
  editedCampaignName.value = campaign.value.name;
  editedTargetGroup.value = campaign.value.target;
  editedCampaignNotes.value = campaign.value.notes;
  editedSubject.value = campaign.value.subject;
  editedFrom.value = campaign.value.from;
}

const sanitizedHtml = computed(() => {
  if (!campaign.value) return ''
  return campaign.value.htmlTemplate
})

async function saveHtmlTemplateToBackend(htmlTemplate) {
  if (!campaign.value) return;
  try {
    await pocketbase.collection('Phishing_Campaign').update(campaign.value.id, {
      HTML: htmlTemplate
    });
  } catch (err) {
    throw err
  }
}

async function toggleHtmlEdit() {
  if (!isEditingHtml.value) {
    editedHtml.value = campaign.value.htmlTemplate
    isEditingHtml.value = true
  } else {
    // Save and exit edit mode
    try {
      await saveHtmlTemplateToBackend(editedHtml.value)
      campaign.value.htmlTemplate = editedHtml.value
      editedHtml.value = ''
      htmlMessage.value = 'Successfully updated HTML template'  
      setTimeout(() => {
        htmlMessage.value = ''
      }, 2000)
    } catch(err) {
      console.error('Failed to save HTML template:', err);
      htmlMessage.value = 'Error updating HTML template'
      setTimeout(() => {
        htmlMessage.value = ''
      }, 2000)
    } finally {
      isEditingHtml.value = false
    }
  }
}

function cancelHtmlEdit() {
  isEditingHtml.value = false;
  editedHtml.value = '';
}
</script>

<template>
  <div class="min-h-screen flex flex-col bg-gray-900">
    <Header />
    <main class="flex-grow">
      <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div v-if="loading" class="text-white text-center py-8">Loading campaign details...</div>
        <div v-else-if="error" class="text-red-500 text-center py-8">{{ error }}</div>
        <div v-else-if="campaign">
          <!-- Campaign Details Box -->
          <div class="bg-gray-800 shadow rounded-lg px-4 py-5 sm:p-6 mb-8">
            <div class="flex justify-between items-start mb-4">
              <div>
                <h1 v-if="!isEditingCampaignInfo" class="text-2xl font-bold text-white mb-2">{{ campaign.name }}</h1>
                <input v-else v-model="editedCampaignName" type="text" class="text-2xl font-bold text-white mb-2 bg-gray-700 rounded px-2 py-1 w-full" />
                <div class="text-gray-400 mb-1">Target Group:
                  <span v-if="!isEditingCampaignInfo" class="text-white">{{ campaign.target }}</span>
                  <input v-else v-model="editedTargetGroup" type="text" class="text-white bg-gray-700 rounded px-2 py-1 ml-2" />
                </div>
                <div class="text-gray-400 mb-1">Subject:
                  <span v-if="!isEditingCampaignInfo" class="text-white">{{ campaign.subject }}</span>
                  <input v-else v-model="editedSubject" type="text" class="text-white bg-gray-700 rounded px-2 py-1 ml-2 w-full" />
                </div>
                <div class="text-gray-400 mb-1">From:
                  <span v-if="!isEditingCampaignInfo" class="text-white">{{ campaign.from }}</span>
                  <input v-else v-model="editedFrom" type="text" class="text-white bg-gray-700 rounded px-2 py-1 ml-2 w-full" />
                </div>
                <div class="text-gray-400 mb-1">Emails Sent: <span class="text-white">{{ campaign.emailsSent || 0 }}</span></div>
                <div class="text-gray-400 mb-1">Emails Clicked: <span class="text-white">{{ campaign.emailsClicked || 0 }} <span v-if="campaign.emailsSent > 0">({{ calculatePercentage(campaign.emailsSent, campaign.emailsClicked) }}%)</span></span></div>
                <div class="text-gray-400 mb-1">Credentials Submitted: <span class="text-white">{{ campaign.credentialsSubmitted || 0 }} <span v-if="campaign.emailsSent > 0">({{ calculatePercentage(campaign.emailsSent, campaign.credentialsSubmitted) }}%)</span></span></div>
              </div>
              <div>
                <button v-if="canEdit() && !isEditingCampaignInfo" @click="isEditingCampaignInfo = true" class="px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-400">Edit</button>
                <div v-if="isEditingCampaignInfo" class="flex gap-2">
                  <button @click="updateCampaignInfo" class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700">Save</button>
                  <button @click="cancelEditCampaignInfo" class="px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700">Cancel</button>
                </div>
              </div>
            </div>
            <div>
              <h2 class="text-lg font-medium text-white mb-2">Campaign Notes</h2>
              <div v-if="!isEditingCampaignInfo" class="bg-gray-700 rounded-lg p-3 min-h-[80px] text-white whitespace-pre-line">{{ campaign.notes }}</div>
              <textarea v-else v-model="editedCampaignNotes" class="w-full bg-gray-700 rounded-lg p-3 text-white min-h-[80px]" />
            </div>
            <p v-if="campaignMessage" class="mt-2 text-sm text-center" :class="{
                'text-green-400': campaignMessage.includes('Success'),
                'text-red-400': campaignMessage.includes('Error')
              }">
                {{ campaignMessage }}
              </p>
          </div>

          <!-- Campaign Setup Box -->
          <div class="bg-gray-800 shadow rounded-lg px-4 py-5 sm:p-6 mb-8">
            <h2 class="text-xl font-bold text-white mb-6">Campaign Setup</h2>
            <!-- HTML Email Preview -->
            <div class="mb-8">
              <div class="flex justify-between items-center mb-2">
                <h4 class="text-md font-medium text-white">HTML Email Preview</h4>
                <p v-if="htmlMessage" class="text-sm text-center" :class="{
                  'text-green-400': htmlMessage.includes('Success'),
                  'text-red-400': htmlMessage.includes('Error')
                }">
                  {{ htmlMessage }}
                </p>
                <div v-if="!isEditingHtml">
                  <button  
                    @click="toggleHtmlEdit"
                    class="bg-gray-600 hover:bg-gray-700 text-white text-xs font-medium py-1 px-3 rounded"
                  >
                    Edit
                  </button>
                </div>
                <div v-else class="flex gap-2">
                  <button
                    @click="toggleHtmlEdit"
                    class="bg-blue-600 hover:bg-blue-700 text-white text-xs font-medium py-1 px-3 rounded"
                  >
                    Save
                  </button>
                  <button
                    @click="cancelHtmlEdit"
                    class="bg-gray-600 hover:bg-gray-700 text-white text-xs font-medium py-1 px-3 rounded"
                  >
                    Cancel
                  </button>
                </div>
              </div>
              <div v-if="!isEditingHtml" class="bg-white rounded-lg overflow-hidden" style="height: 400px;">
                <iframe
                  sandbox="allow-same-origin"
                  class="w-full h-full"
                  title="Email Template Preview"
                  :srcdoc="sanitizedHtml"
                ></iframe>
              </div>
              <div v-else class="bg-gray-900 rounded-lg p-4 max-h-80 overflow-y-auto">
                <textarea
                  v-model="editedHtml"
                  class="w-full bg-gray-700 border-2 border-blue-500 shadow-lg rounded-md text-white p-3 font-mono text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                  rows="16"
                  style="resize: vertical;"
                ></textarea>
              </div>
            </div>

            <!-- Phishlet Configuration -->
            <div class="mb-8">
              <div class="flex justify-between items-center mb-2">
                <h4 class="text-md font-medium text-white">Phishlet Configuration: 
                  <a v-if="campaign.phishletName !== ''" href="/phishlets" target="_blank" class="text-blue-500 hover:text-blue-600">{{ campaign.phishletName }}</a>
                  <span v-else>No phishlet selected</span>
                </h4>
                <p v-if="phishletMessage" class="text-sm text-center" :class="{
                  'text-green-400': phishletMessage.includes('Success'),
                  'text-red-400': phishletMessage.includes('Error')
                }">
                  {{ phishletMessage }}
                </p>
              </div>
              <div class="bg-gray-900 rounded-lg p-4 max-h-80 overflow-y-auto">
                <pre class="text-sm text-gray-300 whitespace-pre-wrap font-mono">{{ campaign.yamlConfig }}</pre>
              </div>
            </div>

            <!-- Caddyfile Section -->
            <div>
              <div class="flex justify-between items-center mb-2">
                <h4 class="text-md font-medium text-white">Caddyfile</h4>
                <p v-if="caddyfileMessage" class="text-sm text-center" :class="{
                  'text-green-400': caddyfileMessage.includes('Success'),
                  'text-red-400': caddyfileMessage.includes('Error')
                }">
                  {{ caddyfileMessage }}
                </p>
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
                    class="bg-blue-600 hover:bg-blue-700 text-white text-xs font-medium py-1 px-2 rounded"
                  >
                    Save
                  </button>
                  <button 
                    v-if="isEditingCaddyfile"
                    @click="toggleEditCaddyfile" 
                    class="bg-gray-600 hover:bg-gray-700 text-white text-xs font-medium py-1 px-2 rounded"
                  >
                    Cancel
                  </button>
                </div>
              </div>
              <div v-if="!isEditingCaddyfile" class="bg-gray-900 rounded-lg p-4 max-h-80 overflow-y-auto">
                <pre class="text-sm text-gray-300 whitespace-pre-wrap font-mono">{{ campaign.caddyfile }}</pre>
              </div>
              <div v-else class="bg-gray-900 rounded-lg p-0 max-h-80 overflow-y-auto">
                <textarea
                  v-model="editedCaddyfile"
                  class="w-full bg-gray-700 border-2 border-blue-500 shadow-lg rounded-md text-white p-3 font-mono text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                  rows="12"
                  style="resize: vertical;"
                ></textarea>
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
pre {
  word-wrap: break-word;
}
</style> 