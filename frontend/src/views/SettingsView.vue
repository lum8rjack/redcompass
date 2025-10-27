<script setup>
  import Header from '@/components/Header.vue'
  import Footer from '@/components/Footer.vue'
  import { ref, computed, inject, onMounted } from 'vue'

  const pocketbase = inject('$pocketbase')
  const isAdmin = ref(false)
  const namecheapApiKey = ref('')
  const namecheapUsername = ref('')
  const namecheapIpAddress = ref('')
  const namecheapErrorMessage = ref('')
  const namecheapID = ref('')
  const namecheapCronJob = ref('')
  
  const porkbunApiKey = ref('')
  const porkbunSecretKey = ref('')
  const porkbunErrorMessage = ref('')
  const porkbunID = ref('')
  const porkbunCronJob = ref('')

  onMounted(async () => {
    try {
      if(pocketbase.authStore.model.role === 'admin') {
        isAdmin.value = true;
      }
      
      // Fetch Namecheap settings
      try {
        const namecheapRecord = await pocketbase.collection('Services').getFirstListItem('Provider="Namecheap"', {
          fields: 'id,Provider,Settings,Cron',
        });
        if (namecheapRecord) {
          namecheapApiKey.value = namecheapRecord.Settings.apiKey
          namecheapUsername.value = namecheapRecord.Settings.username
          namecheapIpAddress.value = namecheapRecord.Settings.ipAddress
          namecheapCronJob.value = namecheapRecord.Cron || ''
          namecheapID.value = namecheapRecord.id
        }
      } catch (error) {
        console.error('Error fetching Namecheap data:', error);
      }
      
      // Fetch Porkbun settings
      try {
        const porkbunRecord = await pocketbase.collection('Services').getFirstListItem('Provider="Porkbun"', {
          fields: 'id,Provider,Settings,Cron',
        });
        if (porkbunRecord) {
          porkbunApiKey.value = porkbunRecord.Settings.apiKey
          porkbunSecretKey.value = porkbunRecord.Settings.secretKey
          porkbunCronJob.value = porkbunRecord.Cron || ''
          porkbunID.value = porkbunRecord.id
        }
      } catch (error) {
        console.error('Error fetching Porkbun data:', error);
      }
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  });

  const saveNamecheapSettings = async () => {
    try {
      namecheapErrorMessage.value = '' // Clear any previous error
      if (namecheapApiKey.value.length < 8) {
        namecheapErrorMessage.value = 'API Key must be at least 8 characters'
        return
      }
      if (namecheapUsername.value.length < 5) {
        namecheapErrorMessage.value = 'Username must be at least 5 characters'
        return
      }
      if (namecheapIpAddress.value.length < 7) {
        namecheapErrorMessage.value = 'Please enter a valid IP address'
        return
      }

      const data = {
        "Provider": "Namecheap",
        "Settings": {
          "apiKey": namecheapApiKey.value,
          "username": namecheapUsername.value,
          "ipAddress": namecheapIpAddress.value,
        },
        "Cron": namecheapCronJob.value
      };

      if (namecheapID.value) {
        await pocketbase.collection('Services').update(namecheapID.value, data);
        namecheapErrorMessage.value = 'Settings updated'
      } else {
        const record =  await pocketbase.collection('Services').create(data);
        namecheapID.value = record.id
        namecheapErrorMessage.value = 'Settings saved'
      }
    } catch (error) {
      console.error('Error saving settings:', error)
      namecheapErrorMessage.value = 'An error occurred while saving settings'
    }
  }

  const deleteNamecheapSettings = async () => {
    try {
      namecheapErrorMessage.value = '' // Clear any previous error
      await pocketbase.collection('Services').delete(namecheapID.value)
      
      // Reset all form values after successful deletion
      namecheapApiKey.value = ''
      namecheapUsername.value = ''
      namecheapIpAddress.value = ''
      namecheapCronJob.value = ''
      namecheapID.value = ''
      
      namecheapErrorMessage.value = 'Settings deleted successfully'
    } catch (error) {
      namecheapErrorMessage.value = 'An error occurred while deleting settings'
    }
  }

  const savePorkbunSettings = async () => {
    try {
      porkbunErrorMessage.value = '' // Clear any previous error
      if (porkbunApiKey.value.length < 8) {
        porkbunErrorMessage.value = 'API Key must be at least 8 characters'
        return
      }
      if (porkbunSecretKey.value.length < 8) {
        porkbunErrorMessage.value = 'Secret Key must be at least 8 characters'
        return
      }

      if (!porkbunApiKey.value.startsWith('pk')) {
        porkbunErrorMessage.value = 'Invalid API Key. API Key must start with "pk"'
        return
      }
      if (!porkbunSecretKey.value.startsWith('sk')) {
        porkbunErrorMessage.value = 'Invalid Secret Key. Secret Key must start with "sk"'
        return
      }

      const data = {
        "Provider": "Porkbun",
        "Settings": {
          "apiKey": porkbunApiKey.value,
          "secretKey": porkbunSecretKey.value,
        },
        "Cron": porkbunCronJob.value
      };

      if (porkbunID.value) {
        await pocketbase.collection('Services').update(porkbunID.value, data);
        porkbunErrorMessage.value = 'Settings updated'
      } else {
        const record =  await pocketbase.collection('Services').create(data);
        porkbunID.value = record.id
        porkbunErrorMessage.value = 'Settings saved'
      }
    } catch (error) {
      console.error('Error saving settings:', error)
      porkbunErrorMessage.value = 'An error occurred while saving settings'
    }
  }

  const deletePorkbunSettings = async () => {
    try {
      porkbunErrorMessage.value = '' // Clear any previous error
      await pocketbase.collection('Services').delete(porkbunID.value)
      
      // Reset all form values after successful deletion
      porkbunApiKey.value = ''
      porkbunSecretKey.value = ''
      porkbunCronJob.value = ''
      porkbunID.value = ''
      
      porkbunErrorMessage.value = 'Settings deleted successfully'
    } catch (error) {
      porkbunErrorMessage.value = 'An error occurred while deleting settings'
    }
  }
</script>

<template>
  <div class="min-h-screen flex flex-col bg-gray-900">
    <Header />
    <main class="flex-grow">
      <div v-if="isAdmin" class="max-w-5xl mx-auto py-6 sm:px-6 lg:px-8">
        <div class="bg-gray-800 rounded-lg shadow p-6">
          <h2 class="text-white text-2xl font-bold mb-6">Namecheap API Settings</h2>
          <form @submit.prevent="saveNamecheapSettings" class="space-y-4">
            <div>
              <label for="apiKey" class="block text-sm font-medium text-gray-300">API Key</label>
              <input
                id="apiKey"
                v-model="namecheapApiKey"
                type="password"
                class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-indigo-500 focus:ring-indigo-500 py-1.5 pl-2"
                placeholder="Enter your Namecheap API key"
              />
            </div>
            <div>
              <label for="username" class="block text-sm font-medium text-gray-300">Username</label>
              <input
                id="username"
                v-model="namecheapUsername"
                type="text"
                class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-indigo-500 focus:ring-indigo-500 py-1.5 pl-2"
                placeholder="Enter your Namecheap username"
              />
            </div>
            <div>
              <label for="ipAddress" class="block text-sm font-medium text-gray-300">IP Address</label>
              <input
                id="ipAddress"
                v-model="namecheapIpAddress"
                type="text"
                class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-indigo-500 focus:ring-indigo-500 py-1.5 pl-2"
                placeholder="Enter your IP address"
              />
            </div>
            <div>
              <label for="cronJob" class="block text-sm font-medium text-gray-300">Cron Job Schedule</label>
              <input
                id="cronJob"
                v-model="namecheapCronJob"
                type="text"
                class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-indigo-500 focus:ring-indigo-500 py-1.5 pl-2"
                placeholder="Enter cron job schedule (e.g., */15 * * * *)"
              />
              <p class="mt-1 text-sm text-gray-400">
                Format: minute hour day month weekday (e.g., */15 * * * * runs every 15 minutes)
              </p>
            </div>
            <button
              type="submit"
              class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-gray-600 hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-400"
            >
              Save Settings
            </button>
            
            <!-- Delete button - only shown when settings exist -->
            <button
              v-if="namecheapID"
              type="button"
              @click="deleteNamecheapSettings"
              class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-400"
            >
              Delete Settings
            </button>

            <p v-if="namecheapErrorMessage" class="mt-2 text-sm text-white text-center">
              {{ namecheapErrorMessage }}
            </p>
          </form>
        </div>
        
        <div class="bg-gray-800 rounded-lg shadow p-6 mt-6">
          <h2 class="text-white text-2xl font-bold mb-6">Porkbun API Settings</h2>
          <form @submit.prevent="savePorkbunSettings" class="space-y-4">
            <div>
              <label for="porkbunApiKey" class="block text-sm font-medium text-gray-300">API Key</label>
              <input
                id="porkbunApiKey"
                v-model="porkbunApiKey"
                type="password"
                class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-indigo-500 focus:ring-indigo-500 py-1.5 pl-2"
                placeholder="Enter your Porkbun API key"
              />
            </div>
            <div>
              <label for="porkbunSecretKey" class="block text-sm font-medium text-gray-300">Secret Key</label>
              <input
                id="porkbunSecretKey"
                v-model="porkbunSecretKey"
                type="password"
                class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-indigo-500 focus:ring-indigo-500 py-1.5 pl-2"
                placeholder="Enter your Porkbun secret key"
              />
            </div>
            <div>
              <label for="porkbunCronJob" class="block text-sm font-medium text-gray-300">Cron Job Schedule</label>
              <input
                id="porkbunCronJob"
                v-model="porkbunCronJob"
                type="text"
                class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-indigo-500 focus:ring-indigo-500 py-1.5 pl-2"
                placeholder="Enter cron job schedule (e.g., */15 * * * *)"
              />
              <p class="mt-1 text-sm text-gray-400">
                Format: minute hour day month weekday (e.g., */15 * * * * runs every 15 minutes)
              </p>
            </div>
            <button
              type="submit"
              class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-gray-600 hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-400"
            >
              Save Settings
            </button>
            
            <!-- Delete button - only shown when settings exist -->
            <button
              v-if="porkbunID"
              type="button"
              @click="deletePorkbunSettings"
              class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-400"
            >
              Delete Settings
            </button>

            <p v-if="porkbunErrorMessage" class="mt-2 text-sm text-white text-center">
              {{ porkbunErrorMessage }}
            </p>
          </form>
        </div>
      </div>
      <div v-else class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <p class="text-white text-center">Unauthorized</p>
      </div>
    </main>
    <Footer />
  </div>
</template>

<style scoped>

</style> 