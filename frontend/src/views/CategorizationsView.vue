<script setup>
import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'

const categorizationServices = [
    {
        name: 'Barracuda',
        url: 'https://www.barracudacentral.org/report',
        actions: 'Submit',
        icon: 'Barracuda'
    },
    {
        name: 'Check Point (requires login)',
        url: 'https://urlcat.checkpoint.com/urlcat/',
        actions: 'Submit',
        icon: 'Check Point'
    },
    {
        name: 'Cisco Talos',
        url: 'https://talosintelligence.com/',
        actions: 'Review',
        icon: 'Cisco Talos'
    },
    {
        name: 'Cloudflare',
        url: 'https://radar.cloudflare.com/domains/feedback',
        actions: 'Review and Submit',
        icon: 'Cloudflare'
    },
    {
        name: 'Google Transparency Report',
        url: 'https://transparencyreport.google.com/safe-browsing/search',
        actions: 'Review',
        icon: 'Google'
    },
    {
        name: 'McAfee',
        url: 'https://sitelookup.mcafee.com/',
        actions: 'Review and Submit',
        icon: 'McAfee'
    },
    {
        name: 'Norton',
        url: 'https://safeweb.norton.com/',
        actions: 'Review and Submit',
        icon: 'Norton'
    },
    {
        name: 'Palo Alto',
        url: 'https://urlfiltering.paloaltonetworks.com/',
        actions: 'Review and Submit',
        icon: 'Palo Alto'
    },
    {
        name: 'SonicWall',
        url: 'https://cfssupport.sonicwall.com/Support/web/eng/newui/viewRating.jsp',
        actions: 'Review and Submit',
        icon: 'SonicWall'
    },
    {
        name: 'Symantec',
        url: 'https://sitereview.bluecoat.com/#/',
        actions: 'Review and Submit',
        icon: 'Symantex'
    },
    {
        name: 'Trellix',
        url: 'https://trustedsource.org/',
        actions: 'Review and Submit',
        icon: 'Trellix'
    },
    {
        name: 'Trend Micro',
        url: 'https://global.sitesafety.trendmicro.com/',
        actions: 'Review and Submit',
        icon: 'Trend Micro'
    },
    {
        name: 'VirusTotal',
        url: 'https://www.virustotal.com/gui/home/url',
        actions: 'Review',
        icon: 'VirusTotal'
    }
]
</script>

<template>
  <div class="min-h-screen flex flex-col bg-gray-900">
    <Header />
    <main class="flex-grow">
      <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <!-- Page Description -->
        <div class="mb-8">
          <h1 class="text-2xl font-bold text-white mb-2">Domain Categorization</h1>
          <p class="text-gray-400">
            Domain categorization is a crucial aspect of web security and content filtering. It helps organizations ensure their domains are properly classified for business use, prevent legitimate domains from being blocked by security filters, maintain appropriate content categorization for compliance, and improve visibility and accessibility of business websites.
          </p>
        </div>

        <!-- Services Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div v-for="service in categorizationServices" :key="service.name" 
               class="bg-gray-800 shadow rounded-lg overflow-hidden">
            <div class="px-4 py-5 sm:p-6">
              <div class="flex items-center justify-between mb-4">
                <h2 class="text-xl font-semibold text-white">{{ service.name }}</h2>
              </div>
              <div class="flex flex-col space-y-3">
                <a :href="service.url" target="_blank" rel="noopener noreferrer"
                   class="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-gray-600 hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-400">
                  Actions: {{ service.actions }}
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
    <Footer />
  </div>
</template> 