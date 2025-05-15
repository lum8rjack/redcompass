<script setup>
import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'
import PageDescription from '@/components/PageDescription.vue'
import { ref, inject, onMounted, computed } from 'vue'
import { canEdit } from '@/utils/auth'
const pocketbase = inject('$pocketbase')

// For file uploads
const htmlTemplateFile = ref(null)
const htmlTemplateText = ref('')
const campaignName = ref('')
const targetGroup = ref('')
const emailSubject = ref('')
const emailFromAddress = ref('')
const phishlet = ref(null)
const updateMessage = ref('')

// Phishing campaigns
const phishingCampaigns = ref([])

// Phishlets
const phishlets = ref([])

const sortColumn = ref('Name');
const sortDirection = ref('asc');

function sortBy(column) {
  if (sortColumn.value === column) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc';
  } else {
    sortColumn.value = column;
    sortDirection.value = 'asc';
  }
}

const sortedPhishingCampaigns = computed(() => {
  const campaigns = [...phishingCampaigns.value];
  const col = sortColumn.value;
  const dir = sortDirection.value;
  return campaigns.sort((a, b) => {
    let aVal = a[col];
    let bVal = b[col];
    if (col === 'total_clicked' || col === 'total_submit') {
      aVal = Number(aVal) || 0;
      bVal = Number(bVal) || 0;
    } else {
      aVal = (aVal || '').toString().toLowerCase();
      bVal = (bVal || '').toString().toLowerCase();
    }
    if (aVal < bVal) return dir === 'asc' ? -1 : 1;
    if (aVal > bVal) return dir === 'asc' ? 1 : -1;
    return 0;
  });
});

async function getPhishingCampaigns() {
  try {
    const response = await pocketbase.collection('Phishing_Campaign_View').getFullList({
        fields: 'id, Name, Target_Group, total_clicked, total_submit',
        sort: 'Name',
    });
    phishingCampaigns.value = response;
  } catch (error) {
    console.error('Error fetching phishing campaigns:', error)
  }
}

onMounted(async () => {
  await getPhishingCampaigns()

  try {
    const response2 = await pocketbase.collection('Phishlets').getFullList({
        fields: 'id,Name',
    });
    phishlets.value = response2;
  } catch (error) {
    console.error('Error fetching phishlets:', error)
  }
})


const createCampaign = async () => {
  if (!campaignName.value) {
    updateMessageBox('Error: You must provide a campaign name')
    return
  }

  if (!targetGroup.value) {
    updateMessageBox('Error: You must provide a target group')
    return
  }
  
  if (!htmlTemplateText.value) {
    updateMessageBox('Error: You must provide an HTML email template')
    return
  }
  
  try {
    const data = {
      "Name": campaignName.value,
      "Example_Subject": emailSubject.value,
      "Example_From": emailFromAddress.value,
      "Target_Group": targetGroup.value,
      "HTML": htmlTemplateText.value,  
      "Phishlet": phishlet.value,
      "Updated_By": pocketbase.authStore.model.id,
    };
    
    await pocketbase.collection('Phishing_Campaign').create(data);
    
    updateMessageBox('Successfully created phishing campaign!')
    await getPhishingCampaigns()
  } catch (error) {
    updateMessageBox('Error creating phishing campaign. Please try again.')
  }
}

function updateMessageBox(message) {
  updateMessage.value = message
  setTimeout(() => {
    updateMessage.value = ''
  }, 2000)
}

const onHtmlFileChange = (event) => {
  const file = event.target.files[0];
  if (!file) {
    htmlTemplateText.value = '';
    return;
  }
  const reader = new FileReader();
  reader.onload = (e) => {
    htmlTemplateText.value = e.target.result;
  };
  reader.readAsText(file);
};
</script>

<template>
  <div class="min-h-screen flex flex-col bg-gray-900">
    <Header />
    <main class="flex-grow">
      <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <PageDescription title="Phishing Campaign Management" description="Manage phishing campaigns for security awareness training. Plan, execute, and track phishing simulations to test and improve your organization's security awareness." />

        <!-- Create New Campaign Section -->
        <div v-if="canEdit()" class="bg-gray-800 shadow rounded-lg mb-8 p-6">
          <h2 class="text-xl font-semibold text-white mb-4">Create New Phishing Campaign</h2>
          
          <div class="space-y-4">
            <!-- Campaign Name -->
            <div>
              <label for="campaignName" class="block text-sm font-medium text-white">Phishing Campaign Name</label>
              <input
                type="text"
                id="campaignName"
                v-model="campaignName"
                class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-500 focus:ring-opacity-50 text-white py-1.5 pl-2"
                placeholder="Enter phishing campaign name"
                required
              />
            </div>

            <!-- Email Subject -->
            <div>
              <label for="emailSubject" class="block text-sm font-medium text-white">Email Subject</label>
              <input
                type="text"
                id="emailSubject"
                v-model="emailSubject"
                class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-500 focus:ring-opacity-50 text-white py-1.5 pl-2"
                placeholder="Enter email subject"
              />
            </div>

            <!-- Email From Address -->
            <div>
              <label for="emailFromAddress" class="block text-sm font-medium text-white">Email From Address</label>
              <input
                type="text"
                id="emailFromAddress"
                v-model="emailFromAddress"
                class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-500 focus:ring-opacity-50 text-white py-1.5 pl-2"
                placeholder="Enter email from address"
              />
            </div>
            
            <!-- Target Group -->
            <div>
              <label for="targetGroup" class="block text-sm font-medium text-white">Target Group</label>
              <input
                type="text"
                id="targetGroup"
                v-model="targetGroup"
                class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-500 focus:ring-opacity-50 text-white py-1.5 pl-2"
                placeholder="e.g., Marketing, Engineering, All Users"
                required
              />
            </div>
            
            <!-- HTML Template Upload -->
            <div>
              <label for="htmlTemplate" class="block text-sm font-medium text-white">HTML Email Template</label>
              <input
                type="file"
                id="htmlTemplateFile"
                accept=".html"
                @change="onHtmlFileChange"
                class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-500 focus:ring-opacity-50 text-white py-1.5 pl-2"
              />
              <p class="mt-1 text-xs text-gray-400">Upload your HTML email template file</p>
            </div>
            
            <!-- Phishlet -->
            <div>
              <label for="phishlet" class="block text-sm font-medium text-white">Phishlet</label>
              <select
                id="phishlet"
                v-model="phishlet"
                class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-500 focus:ring-opacity-50 text-white py-1.5 pl-2"
              >
                <option value="" disabled selected>Select a phishlet</option>
                <option 
                  v-for="phishlet in phishlets" 
                  :key="phishlet.id"
                  :value="phishlet.id"
                >
                  {{ phishlet.Name }}
                </option>
              </select>
              <p class="mt-1 text-xs text-gray-400">Select your phishlet</p>
            </div>
            
            <!-- Submit Button -->
            <div class="flex flex-col items-center">
              <button
                @click="createCampaign"
                class="bg-gray-600 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded"
              >
                Create Campaign
              </button>
              <p v-if="updateMessage" class="mt-2 text-sm" :class="{
                'text-green-400': updateMessage.includes('Success'),
                'text-red-400': updateMessage.includes('Error')
              }">
                {{ updateMessage }}
              </p>
            </div>
          </div>
        </div>

        <!-- Phishing Campaigns Section -->
        <div class="mb-8">
          <h2 class="text-xl font-semibold text-white mb-4">Phishing Campaigns</h2>
          
          <div class="bg-gray-800 shadow overflow-hidden rounded-lg">
            <table class="min-w-full divide-y divide-gray-700">
              <thead class="bg-gray-700">
                <tr>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider w-1/4 cursor-pointer" @click="sortBy('Name')">
                    Campaign Name
                    <span v-if="sortColumn === 'Name'">{{ sortDirection === 'asc' ? '▲' : '▼' }}</span>
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider w-1/4 cursor-pointer" @click="sortBy('Target_Group')">
                    Target Group
                    <span v-if="sortColumn === 'Target_Group'">{{ sortDirection === 'asc' ? '▲' : '▼' }}</span>
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider w-1/4 cursor-pointer" @click="sortBy('total_clicked')">
                    Click Rate
                    <span v-if="sortColumn === 'total_clicked'">{{ sortDirection === 'asc' ? '▲' : '▼' }}</span>
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider w-1/4 cursor-pointer" @click="sortBy('total_submit')">
                    Submitted Credentials
                    <span v-if="sortColumn === 'total_submit'">{{ sortDirection === 'asc' ? '▲' : '▼' }}</span>
                  </th>
                </tr>
              </thead>
              <tbody class="bg-gray-800 divide-y divide-gray-700">
                <tr 
                  v-for="campaign in sortedPhishingCampaigns" 
                  :key="campaign.id" 
                  class="hover:bg-gray-700 cursor-pointer"
                  @click="$router.push(`/phishing/${campaign.id}`)"
                >
                  <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-white">
                    {{ campaign.Name }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300">
                    {{ campaign.Target_Group }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300">
                    {{ campaign.total_clicked ? campaign.total_clicked : 0 }}%
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300">
                    {{ campaign.total_submit ? campaign.total_submit : 0 }}%
                  </td>
                </tr>
                <tr v-if="sortedPhishingCampaigns.length === 0">
                  <td colspan="4" class="px-6 py-4 text-center text-sm text-gray-400">
                    No campaigns found
                  </td>
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