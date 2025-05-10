<script setup>
import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'
import { ref, inject, onMounted, computed, watch } from 'vue'

const pocketbase = inject('$pocketbase')

// For file uploads
const htmlTemplateFile = ref(null)
const yamlConfigFile = ref(null)
const campaignName = ref('')
const targetGroup = ref('')
const uploadSuccess = ref(false)
const errorMessage = ref('')

// Campaign detail modal
const showCampaignModal = ref(false)
const selectedCampaign = ref(null)
const htmlPreview = ref('')
const yamlContent = ref('')

// Add a new ref for clipboard notification
const showClipboardNotification = ref(false)
const clipboardContent = ref('')  // To track what was copied

// Add these new refs for editing
const isEditingMetrics = ref(false)
const isEditingPhishlet = ref(false)
const isEditingCaddyfile = ref(false)
const editedCredentials = ref('')
const editedClickRate = ref('')
const editedPhishlet = ref('')
const editedCaddyfile = ref('')
const campaignNotes = ref('')

// Add a ref for tracking unsaved changes and modal to confirm closing
const hasUnsavedChanges = ref(false)
const showCloseConfirmation = ref(false)
const originalCampaign = ref(null)

// Add these new refs for editing campaign title and target
const isEditingCampaignInfo = ref(false)
const editedCampaignName = ref('')
const editedTargetGroup = ref('')

// Sample phishing tools and resources
const phishingTools = ref([
    {
        name: 'GoPhish',
        description: 'Open-source phishing toolkit designed for businesses to conduct phishing awareness campaigns',
        url: 'https://getgophish.com/',
        icon: 'GoPhish'
    },
    {
        name: 'EvilGinx2',
        description: 'Advanced phishing framework for capturing credentials and session cookies',
        url: 'https://github.com/kgretzky/evilginx2',
        icon: 'EvilGinx2'
    },
    {
        name: 'SocialFish',
        description: 'Educational tool for social media phishing',
        url: 'https://github.com/UndeadSec/SocialFish',
        icon: 'SocialFish'
    },
    {
        name: 'Lucy',
        description: 'Commercial phishing awareness solution',
        url: 'https://lucysecurity.com/',
        icon: 'Lucy'
    },
    {
        name: 'King Phisher',
        description: 'Tool for testing and promoting user awareness',
        url: 'https://github.com/rsmusllp/king-phisher',
        icon: 'King Phisher'
    }
])

// Replace with phishlets tracking
const phishlets = ref([
  {
    id: 1,
    name: 'Microsoft 365',
    uploadDate: '2024-05-15',
    uploadedBy: 'admin',
    notes: '',
    content: `domain: login.microsoftonline.com
  keys:
  - key.pem
  - cert.pem
  originURL: https://login.microsoftonline.com
  proxies:
    - type: template
      target: https://login.microsoftonline.com
      pathScope: /
      triggers:
        - trigger: login.microsoftonline.com
          params:
            subdomain: login
      headers:
        Set-Cookie: JSESSIONID={session_id}; Domain=.microsoftonline.com; Secure; HttpOnly
  credentials:
    username:
      key: login
      search: '(.*)'
      type: post
    password:
      key: passwd
      search: '(.*)'
      type: post`
  },
  {
    id: 2,
    name: 'Google Workspace',
    uploadDate: '2024-05-10',
    uploadedBy: 'analyst1',
    notes: '',
    content: `domain: accounts.google.com
  keys:
  - key.pem
  - cert.pem
  originURL: https://accounts.google.com
  proxies:
    - type: template
      target: https://accounts.google.com
      pathScope: /
      triggers:
        - trigger: accounts.google.com
          params:
            subdomain: accounts
      headers:
        Set-Cookie: SIDCC={session_id}; Domain=.google.com; Secure; HttpOnly
  credentials:
    username:
      key: identifier
      search: '(.*)'
      type: post
    password:
      key: password
      search: '(.*)'
      type: post`
  },
  {
    id: 3,
    name: 'Okta SSO',
    uploadDate: '2024-04-30',
    uploadedBy: 'redteam',
    notes: '',
    content: `domain: login.okta.com
  keys:
  - key.pem
  - cert.pem
  originURL: https://login.okta.com
  proxies:
    - type: template
      target: https://login.okta.com
      pathScope: /
      triggers:
        - trigger: login.okta.com
          params:
            subdomain: login
      headers:
        Set-Cookie: OKTA_SESSION={session_id}; Domain=.okta.com; Secure; HttpOnly
  credentials:
    username:
      key: username
      search: '(.*)'
      type: post
    password:
      key: password
      search: '(.*)'
      type: post`
  },
  {
    id: 4,
    name: 'Salesforce',
    uploadDate: '2024-04-22',
    uploadedBy: 'pentester',
    notes: '',
    content: `domain: login.salesforce.com
  keys:
  - key.pem
  - cert.pem
  originURL: https://login.salesforce.com
  proxies:
    - type: template
      target: https://login.salesforce.com
      pathScope: /
      triggers:
        - trigger: login.salesforce.com
          params:
            subdomain: login
      headers:
        Set-Cookie: sid={session_id}; Domain=.salesforce.com; Secure; HttpOnly
  credentials:
    username:
      key: username
      search: '(.*)'
      type: post
    password:
      key: password
      search: '(.*)'
      type: post`
  },
  {
    id: 5,
    name: 'AWS Console',
    uploadDate: '2024-04-15',
    uploadedBy: 'admin',
    notes: '',
    content: `domain: signin.aws.amazon.com
  keys:
  - key.pem
  - cert.pem
  originURL: https://signin.aws.amazon.com
  proxies:
    - type: template
      target: https://signin.aws.amazon.com
      pathScope: /
      triggers:
        - trigger: signin.aws.amazon.com
          params:
            subdomain: signin
      headers:
        Set-Cookie: aws-userInfo={session_id}; Domain=.amazon.com; Secure; HttpOnly
  credentials:
    username:
      key: email
      search: '(.*)'
      type: post
    password:
      key: password
      search: '(.*)'
      type: post`
  }
])

// Phishlet upload functionality
const phishletFile = ref(null)
const phishletName = ref('')
const showPhishletUploadModal = ref(false)
const phishletUploadSuccess = ref(false)
const phishletUploadError = ref('')

// Phishlet detail modal
const showPhishletDetailModal = ref(false)
const selectedPhishlet = ref(null)
const editedPhishletContent = ref('')
const isEditingPhishletContent = ref(false)
const phishletContentChanged = ref(false)

// Add these new refs for phishlet notes
const phishletNotes = ref('');
const phishletNotesChanged = ref(false);

// Sample HTML templates and YAML configs for the demo
const sampleHtmlTemplates = {
  1: `<!DOCTYPE html>
  <html>
  <head>
    <meta charset="UTF-8">
    <title>Important Security Update</title>
    <style>
      body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
      .header { background-color: #0066cc; color: white; padding: 10px; text-align: center; }
      .content { padding: 20px; }
      .button { background-color: #0066cc; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; }
    </style>
  </head>
  <body>
    <div class="header">
      <h1>IT Security Department</h1>
    </div>
    <div class="content">
      <h2>Action Required: Security Training Update</h2>
      <p>Dear Marketing Team Member,</p>
      <p>Our security system detected that your account requires an immediate security verification. Please complete your security awareness training by clicking the button below.</p>
      <p><a href="#" class="button">Complete Security Training</a></p>
      <p>This training must be completed within 24 hours to maintain access to company resources.</p>
      <p>Thank you for your cooperation,<br>IT Security Team</p>
    </div>
  </body>
  </html>`,
  2: `<!DOCTYPE html>
  <html>
  <head>
    <meta charset="UTF-8">
    <title>Executive Dashboard Access Update</title>
    <style>
      body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
      .header { background-color: #2c3e50; color: white; padding: 15px; text-align: center; }
      .content { padding: 25px; }
      .button { background-color: #2c3e50; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; }
    </style>
  </head>
  <body>
    <div class="header">
      <h1>Executive Portal</h1>
    </div>
    <div class="content">
      <h2>Urgent: Dashboard Access Update Required</h2>
      <p>Dear Executive Team Member,</p>
      <p>We've updated our executive dashboard with new security features. To maintain access, please verify your credentials by clicking below.</p>
      <p><a href="#" class="button">Verify Executive Access</a></p>
      <p>This verification is mandatory and will expire in 8 hours.</p>
      <p>Regards,<br>IT Department</p>
    </div>
  </body>
  </html>`,
  3: `<!DOCTYPE html>
  <html>
  <head>
    <meta charset="UTF-8">
    <title>New Employee Onboarding</title>
    <style>
      body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
      .header { background-color: #27ae60; color: white; padding: 10px; text-align: center; }
      .content { padding: 20px; }
      .button { background-color: #27ae60; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; }
    </style>
  </head>
  <body>
    <div class="header">
      <h1>Welcome to the Company!</h1>
    </div>
    <div class="content">
      <h2>Complete Your Onboarding</h2>
      <p>Dear New Team Member,</p>
      <p>Welcome to the team! To complete your onboarding process, please set up your account by clicking the button below.</p>
      <p><a href="#" class="button">Complete Onboarding</a></p>
      <p>Please complete this step within the next 48 hours.</p>
      <p>Best regards,<br>HR Department</p>
    </div>
  </body>
  </html>`
};

const sampleYamlConfigs = {
  1: `# Phishing Campaign Configuration for 'Awareness Training Q2'
name: Awareness Training Q2
target:
  group: Marketing Team
  email_domain: company.com
  count: 75
schedule:
  start_date: 2024-04-10
  end_date: 2024-04-20
  send_time: '09:00'
template:
  subject: "ACTION REQUIRED: Complete Your Security Training"
  sender: "IT Security <security@company-training.com>"
  reply_to: "no-reply@company-training.com"
tracking:
  landing_page: "https://training.company.com/verify"
  click_tracking: true
  open_tracking: true
  form_capture: true
notifications:
  admin_email: "admin@company.com"
  success_alerts: true
  progress_updates: daily`,
  
  2: `# Phishing Campaign Configuration for 'Security Assessment'
name: Security Assessment
target:
  group: Executive Team
  email_domain: company.com
  count: 15
schedule:
  start_date: 2024-05-01
  end_date: 2024-05-15
  send_time: '10:30'
template:
  subject: "Urgent: Executive Dashboard Access Update"
  sender: "IT Department <it@company-portal.com>"
  reply_to: "support@company-portal.com"
tracking:
  landing_page: "https://portal.company.com/verify"
  click_tracking: true
  open_tracking: true
  form_capture: true
  credentials_capture: true
notifications:
  admin_email: "securityteam@company.com"
  success_alerts: true
  progress_updates: hourly`,
  
  3: `# Phishing Campaign Configuration for 'New Hire Training'
name: New Hire Training
target:
  group: New Employees
  email_domain: company.com
  count: 25
  filter: "hire_date > 2024-01-01"
schedule:
  start_date: 2024-06-01
  end_date: 2024-06-10
  send_time: '11:00'
template:
  subject: "Welcome - Complete Your Onboarding"
  sender: "HR Department <hr@company-onboarding.com>"
  reply_to: "onboarding@company-onboarding.com"
tracking:
  landing_page: "https://onboarding.company.com/setup"
  click_tracking: true
  open_tracking: true
  form_capture: true
notifications:
  admin_email: "training@company.com"
  success_alerts: true
  progress_updates: daily`
};

// Sample phishing campaigns
const phishingCampaigns = ref([])

// Sample Caddyfiles for demo
const sampleCaddyfiles = {
  1: `example.com {
  root * /var/www/example.com
  reverse_proxy /api/* localhost:8080
  tls internal
  log {
    output file /var/log/caddy/example.com.log
  }
}`,
  2: `executive-portal.company.com {
  tls /etc/certs/executive-portal.crt /etc/certs/executive-portal.key
  reverse_proxy localhost:3000
  handle_errors {
    respond "{http.error.status_code} {http.error.status_text}"
  }
  log {
    output file /var/log/caddy/executive-portal.log
  }
}`,
  3: `onboarding.company.com {
  tls {
    protocols tls1.2 tls1.3
    ciphers TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256 TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
  }
  root * /var/www/onboarding
  php_fastcgi unix//var/run/php-fpm.sock
  file_server
  log {
    output file /var/log/caddy/onboarding.log
  }
}`
};

onMounted(async () => {
  // You would typically load data from your backend here
  // For now, we'll use sample data
  phishingCampaigns.value = [
    {
      id: 1,
      name: 'Awareness Training Q2',
      target: 'Marketing Team',
      status: 'Completed',
      clickRate: '23%',
      submittedCredentials: '12%',
      htmlTemplate: sampleHtmlTemplates[1],
      yamlConfig: sampleYamlConfigs[1],
      caddyfile: sampleCaddyfiles[1],
      notes: 'Initial campaign targeting marketing team. Good response rate with 23% clicks. We should follow up with additional training for those who submitted credentials.'
    },
    {
      id: 2,
      name: 'Security Assessment',
      target: 'Executive Team',
      status: 'In Progress',
      clickRate: '15%',
      submittedCredentials: '7%',
      htmlTemplate: sampleHtmlTemplates[2],
      yamlConfig: sampleYamlConfigs[2],
      caddyfile: sampleCaddyfiles[2],
      notes: 'Executive team showing better awareness than previous quarter. Still concerned about the 7% credential submission rate.'
    },
    {
      id: 3,
      name: 'New Hire Training',
      target: 'New Employees',
      status: 'Scheduled',
      clickRate: 'N/A',
      submittedCredentials: 'N/A',
      htmlTemplate: sampleHtmlTemplates[3],
      yamlConfig: sampleYamlConfigs[3],
      caddyfile: sampleCaddyfiles[3],
      notes: 'Campaign scheduled to begin next week. Will target all employees hired in the last 6 months.'
    }
  ]
  
  // Calculate and update metrics based on existing project data
  phishingCampaigns.value.forEach(campaign => {
    const metrics = calculateAggregateMetrics(campaign.id);
    if (metrics.totalSent > 0) {
      campaign.clickRate = metrics.clickRate;
      campaign.submittedCredentials = metrics.submissionRate;
      campaign.totalSent = metrics.totalSent;
      campaign.totalClicked = metrics.totalClicked;
      campaign.totalSubmitted = metrics.totalSubmitted;
    }
  });
})

const handleHtmlTemplateChange = (event) => {
  const file = event.target.files[0]
  if (file && file.type === 'text/html') {
    htmlTemplateFile.value = file
    errorMessage.value = ''
    
    // Read the file for preview
    const reader = new FileReader()
    reader.onload = (e) => {
      htmlPreview.value = e.target.result
    }
    reader.readAsText(file)
  } else {
    errorMessage.value = 'Please select a valid HTML file'
    event.target.value = null
  }
}

const handleYamlConfigChange = (event) => {
  const file = event.target.files[0]
  if (file && (file.type === 'application/x-yaml' || file.name.endsWith('.yaml') || file.name.endsWith('.yml'))) {
    yamlConfigFile.value = file
    errorMessage.value = ''
    
    // Read the file content
    const reader = new FileReader()
    reader.onload = (e) => {
      yamlContent.value = e.target.result
    }
    reader.readAsText(file)
  } else {
    errorMessage.value = 'Please select a valid YAML file'
    event.target.value = null
  }
}

const createCampaign = async () => {
  if (!campaignName.value) {
    errorMessage.value = 'Please enter a campaign name'
    return
  }
  
  if (!htmlTemplateFile.value) {
    errorMessage.value = 'Please upload an HTML email template'
    return
  }
  
  if (!yamlConfigFile.value) {
    errorMessage.value = 'Please upload a YAML configuration file'
    return
  }
  
  // Here you would typically upload the files to the server
  // For now, we'll just simulate a successful upload
  
  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // Add to the list
    phishingCampaigns.value.unshift({
      id: Date.now(),
      name: campaignName.value,
      target: targetGroup.value || 'All Users',
      status: 'Scheduled',
      clickRate: 'N/A',
      submittedCredentials: 'N/A',
      htmlTemplate: htmlPreview.value || '<p>HTML template preview not available</p>',
      yamlConfig: yamlContent.value || '# YAML configuration not available',
      caddyfile: '',
      notes: ''
    })
    
    // Reset form
    campaignName.value = ''
    targetGroup.value = ''
    htmlTemplateFile.value = null
    yamlConfigFile.value = null
    htmlPreview.value = ''
    yamlContent.value = ''
    
    // Show success message
    uploadSuccess.value = true
    setTimeout(() => {
      uploadSuccess.value = false
    }, 3000)
    
  } catch (error) {
    console.error('Error creating campaign:', error)
    errorMessage.value = 'Error creating campaign. Please try again.'
  }
}

const viewCampaignDetails = (campaign) => {
  selectedCampaign.value = campaign
  // Store a deep copy of the original campaign to detect changes
  originalCampaign.value = JSON.parse(JSON.stringify(campaign))
  editedCredentials.value = campaign.submittedCredentials
  editedClickRate.value = campaign.clickRate
  editedPhishlet.value = campaign.yamlConfig
  editedCaddyfile.value = campaign.caddyfile
  campaignNotes.value = campaign.notes || ''
  editedCampaignName.value = campaign.name
  editedTargetGroup.value = campaign.target
  
  // Load campaign artifacts
  campaignArtifacts.value = sampleArtifacts[campaign.id] || []
  
  // Reset form fields
  artifactName.value = ''
  artifactDescription.value = ''
  artifactFile.value = null
  artifactUploadSuccess.value = false
  artifactUploadError.value = ''
  
  // Calculate and set aggregate metrics
  const metrics = calculateAggregateMetrics(campaign.id);
  selectedCampaign.value.totalSent = metrics.totalSent || 0;
  selectedCampaign.value.totalClicked = metrics.totalClicked || 0;
  selectedCampaign.value.totalSubmitted = metrics.totalSubmitted || 0;
  
  isEditingMetrics.value = false
  isEditingPhishlet.value = false
  isEditingCaddyfile.value = false
  isEditingCampaignInfo.value = false
  hasUnsavedChanges.value = false
  showCampaignModal.value = true
  
  // Reset the add project form
  newProjectName.value = '';
  newProjectSent.value = '';
  newProjectClicked.value = '';
  newProjectSubmitted.value = '';
  showAddProjectForm.value = false;
}

// Attempt to close the modal with confirmation if needed
const attemptCloseModal = () => {
  // Always check for unsaved changes directly instead of using the state variable
  const hasChanges = checkForUnsavedChanges();
  console.log('Attempting to close modal, hasUnsavedChanges:', hasChanges);
  console.log('Original notes:', originalCampaign.value?.notes);
  console.log('Current notes:', campaignNotes.value);
  
  if (hasChanges) {
    // Show confirmation dialog
    showCloseConfirmation.value = true;
  } else {
    // No changes, just close
    actuallyCloseModal();
  }
}

// The actual close modal function that doesn't check for changes
const actuallyCloseModal = () => {
  showCampaignModal.value = false;
  selectedCampaign.value = null;
  originalCampaign.value = null;
  showCloseConfirmation.value = false;
  hasUnsavedChanges.value = false;
}

// Close modal and discard changes
const discardChangesAndClose = () => {
  showCloseConfirmation.value = false;
  actuallyCloseModal();
}

// Save all changes and close modal
const saveAllChangesAndClose = () => {
  // Save campaign info if editing
  if (isEditingCampaignInfo.value) {
    updateCampaignInfo();
  }
  
  // Save metrics if editing
  if (isEditingMetrics.value) {
    updateCampaignMetrics();
  }
  
  // Save phishlet if editing
  if (isEditingPhishlet.value) {
    savePhishlet();
  }
  
  // Save caddyfile if editing
  if (isEditingCaddyfile.value) {
    saveCaddyfile();
  }
  
  // Save notes
  saveNotes();
  
  // Close confirmation and modal
  showCloseConfirmation.value = false;
  actuallyCloseModal();
}

// Current closeCampaignModal acts as a wrapper for actuallyCloseModal
const closeCampaignModal = attemptCloseModal;

const downloadYamlConfig = () => {
  if (!selectedCampaign.value) return
  
  const blob = new Blob([selectedCampaign.value.yamlConfig], { type: 'text/yaml' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${selectedCampaign.value.name.replace(/\s+/g, '_')}_config.yaml`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

const copyYamlToClipboard = async () => {
  if (!selectedCampaign.value) return
  
  try {
    await navigator.clipboard.writeText(selectedCampaign.value.yamlConfig)
    clipboardContent.value = 'Phishlet configuration'
    showClipboardNotification.value = true
    setTimeout(() => {
      showClipboardNotification.value = false
    }, 2000)
  } catch (err) {
    console.error('Failed to copy YAML to clipboard:', err)
    alert('Failed to copy to clipboard')
  }
}

const copyCaddyToClipboard = async () => {
  if (!selectedCampaign.value) return
  
  try {
    await navigator.clipboard.writeText(selectedCampaign.value.caddyfile)
    clipboardContent.value = 'Caddyfile'
    showClipboardNotification.value = true
    setTimeout(() => {
      showClipboardNotification.value = false
    }, 2000)
  } catch (err) {
    console.error('Failed to copy Caddyfile to clipboard:', err)
    alert('Failed to copy to clipboard')
  }
}

const downloadCaddyfile = () => {
  if (!selectedCampaign.value) return
  
  const blob = new Blob([selectedCampaign.value.caddyfile], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${selectedCampaign.value.name.replace(/\s+/g, '_')}_caddyfile`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

// Sanitized HTML for preview
const sanitizedHtml = computed(() => {
  if (!selectedCampaign.value) return ''
  return selectedCampaign.value.htmlTemplate
})

// Add this function to update campaign metrics
const updateCampaignMetrics = () => {
  if (!selectedCampaign.value) return
  
  // Update the campaign with edited values
  if (editedCredentials.value) {
    selectedCampaign.value.submittedCredentials = editedCredentials.value
  }
  
  if (editedClickRate.value) {
    selectedCampaign.value.clickRate = editedClickRate.value
  }
  
  // Update notes
  selectedCampaign.value.notes = campaignNotes.value
  
  // In a real app, you would save this to the server
  // For demo purposes, we'll just update the local state
  const campaignIndex = phishingCampaigns.value.findIndex(c => c.id === selectedCampaign.value.id)
  if (campaignIndex !== -1) {
    phishingCampaigns.value[campaignIndex].submittedCredentials = editedCredentials.value
    phishingCampaigns.value[campaignIndex].clickRate = editedClickRate.value
    phishingCampaigns.value[campaignIndex].notes = campaignNotes.value
  }
  
  // Exit edit mode
  isEditingMetrics.value = false
}

// Add this function to save notes
const saveNotes = () => {
  if (!selectedCampaign.value) return
  
  // Update the campaign with notes
  selectedCampaign.value.notes = campaignNotes.value
  
  // In a real app, you would save this to the server
  // For demo purposes, we'll just update the local state
  const campaignIndex = phishingCampaigns.value.findIndex(c => c.id === selectedCampaign.value.id)
  if (campaignIndex !== -1) {
    phishingCampaigns.value[campaignIndex].notes = campaignNotes.value
  }
  
  // Reset the unsaved changes flag and update the original campaign copy
  if (originalCampaign.value) {
    originalCampaign.value.notes = campaignNotes.value;
  }
  hasUnsavedChanges.value = false;
}

// Toggle editing phishlet
const toggleEditPhishlet = () => {
  isEditingPhishlet.value = !isEditingPhishlet.value
  if (isEditingPhishlet.value) {
    editedPhishlet.value = selectedCampaign.value.yamlConfig
  }
}

// Save phishlet changes
const savePhishlet = () => {
  if (!selectedCampaign.value) return
  
  // Update the campaign with edited phishlet
  selectedCampaign.value.yamlConfig = editedPhishlet.value
  
  // In a real app, you would save this to the server
  // For demo purposes, we'll just update the local state
  const campaignIndex = phishingCampaigns.value.findIndex(c => c.id === selectedCampaign.value.id)
  if (campaignIndex !== -1) {
    phishingCampaigns.value[campaignIndex].yamlConfig = editedPhishlet.value
  }
  
  // Exit edit mode
  isEditingPhishlet.value = false
}

// Toggle editing caddyfile
const toggleEditCaddyfile = () => {
  isEditingCaddyfile.value = !isEditingCaddyfile.value
  if (isEditingCaddyfile.value) {
    editedCaddyfile.value = selectedCampaign.value.caddyfile
  }
}

// Save caddyfile changes
const saveCaddyfile = () => {
  if (!selectedCampaign.value) return
  
  // Update the campaign with edited caddyfile
  selectedCampaign.value.caddyfile = editedCaddyfile.value
  
  // In a real app, you would save this to the server
  // For demo purposes, we'll just update the local state
  const campaignIndex = phishingCampaigns.value.findIndex(c => c.id === selectedCampaign.value.id)
  if (campaignIndex !== -1) {
    phishingCampaigns.value[campaignIndex].caddyfile = editedCaddyfile.value
  }
  
  // Exit edit mode
  isEditingCaddyfile.value = false
}

// Watch for changes in the notes to update hasUnsavedChanges flag
watch(campaignNotes, (newValue) => {
  if (selectedCampaign.value && originalCampaign.value) {
    const originalNotes = originalCampaign.value.notes || '';
    if (newValue !== originalNotes) {
      hasUnsavedChanges.value = true;
    }
  }
})

// Watch for changes in phishlet config
watch(editedPhishlet, (newValue) => {
  if (selectedCampaign.value && originalCampaign.value && isEditingPhishlet.value) {
    if (newValue !== originalCampaign.value.yamlConfig) {
      hasUnsavedChanges.value = true;
    }
  }
})

// Watch for changes in caddyfile
watch(editedCaddyfile, (newValue) => {
  if (selectedCampaign.value && originalCampaign.value && isEditingCaddyfile.value) {
    if (newValue !== originalCampaign.value.caddyfile) {
      hasUnsavedChanges.value = true;
    }
  }
})

// Watch for changes in metrics
watch([editedClickRate, editedCredentials], () => {
  if (selectedCampaign.value && originalCampaign.value && isEditingMetrics.value) {
    if (editedClickRate.value !== originalCampaign.value.clickRate || 
        editedCredentials.value !== originalCampaign.value.submittedCredentials) {
      hasUnsavedChanges.value = true;
    }
  }
})

// Check if there are any unsaved changes
const checkForUnsavedChanges = () => {
  if (!selectedCampaign.value || !originalCampaign.value) return false
  
  // Check for changes in campaign info
  if (isEditingCampaignInfo.value && 
      (editedCampaignName.value !== selectedCampaign.value.name || 
       editedTargetGroup.value !== selectedCampaign.value.target)) {
    console.log("Campaign info has changed");
    return true
  }
  
  // Check for changes in metrics
  if (isEditingMetrics.value && 
      (editedClickRate.value !== selectedCampaign.value.clickRate || 
       editedCredentials.value !== selectedCampaign.value.submittedCredentials)) {
    console.log("Metrics have changed");
    return true
  }
  
  // Check for changes in phishlet config
  if (isEditingPhishlet.value && editedPhishlet.value !== selectedCampaign.value.yamlConfig) {
    console.log("Phishlet config has changed");
    return true
  }
  
  // Check for changes in caddyfile
  if (isEditingCaddyfile.value && editedCaddyfile.value !== selectedCampaign.value.caddyfile) {
    console.log("Caddyfile has changed");
    return true
  }
  
  // Check for changes in notes - compare with original campaign notes to detect changes
  const originalNotes = originalCampaign.value.notes || '';
  const currentNotes = campaignNotes.value || '';
  
  if (originalNotes !== currentNotes) {
    console.log("Notes have changed");
    console.log("Original:", originalNotes);
    console.log("Current:", currentNotes);
    return true;
  }
  
  // No unsaved changes detected
  console.log("No changes detected");
  return false
}

// Add this function to handle notes changes
const handleNotesChange = () => {
  const originalNotes = originalCampaign.value?.notes || '';
  const currentNotes = campaignNotes.value || '';
  
  if (originalNotes !== currentNotes) {
    console.log('Notes changed from:', originalNotes, 'to:', currentNotes);
    hasUnsavedChanges.value = true;
  }
}

// Add watches for all editable fields to better track unsaved changes

// Handle phishlet file selection
const handlePhishletFileChange = (event) => {
  const file = event.target.files[0]
  if (file) {
    phishletFile.value = file
    
    // Read the file for preview
    const reader = new FileReader()
    reader.onload = (e) => {
      try {
        const content = e.target.result
        // Store content for later use
        phishletFileContent.value = content
      } catch (error) {
        console.error('Error reading phishlet file:', error)
        phishletUploadError.value = 'Error reading phishlet file'
      }
    }
    reader.readAsText(file)
  }
}

// Phishlet file content
const phishletFileContent = ref('')

// Upload phishlet
const uploadPhishlet = async () => {
  if (!phishletName.value) {
    phishletUploadError.value = 'Please enter a phishlet name'
    return
  }
  
  if (!phishletFile.value) {
    phishletUploadError.value = 'Please select a phishlet file'
    return
  }
  
  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 800))
    
    // Add to the list
    phishlets.value.unshift({
      id: Date.now(),
      name: phishletName.value,
      uploadDate: new Date().toISOString().split('T')[0],
      upvotes: 0,
      downvotes: 0,
      uploadedBy: 'current_user', // In a real app, use the logged-in user
      content: phishletFileContent.value || '# Empty phishlet content'
    })
    
    // Reset form
    phishletName.value = ''
    phishletFile.value = null
    phishletFileContent.value = ''
    
    // Show success and hide modal
    phishletUploadSuccess.value = true
    setTimeout(() => {
      phishletUploadSuccess.value = false
      showPhishletUploadModal.value = false
    }, 1500)
    
  } catch (error) {
    console.error('Error uploading phishlet:', error)
    phishletUploadError.value = 'Error uploading phishlet. Please try again.'
  }
}

// Vote on a phishlet
const voteOnPhishlet = (phishlet, isUpvote) => {
  const index = phishlets.value.findIndex(p => p.id === phishlet.id)
  if (index !== -1) {
    if (isUpvote) {
      phishlets.value[index].upvotes++
    } else {
      phishlets.value[index].downvotes++
    }
  }
}

// Add refs for tracking phishlet unsaved changes
const hasUnsavedPhishletChanges = ref(false)
const showPhishletCloseConfirmation = ref(false)
const originalPhishlet = ref(null)

// View phishlet details
const viewPhishletDetails = (phishlet) => {
  selectedPhishlet.value = JSON.parse(JSON.stringify(phishlet)); // Create a deep copy
  originalPhishlet.value = JSON.parse(JSON.stringify(phishlet)); // Store original for comparison
  editedPhishletContent.value = phishlet.content;
  phishletNotes.value = phishlet.notes || ''; // Set the notes
  isEditingPhishletContent.value = false;
  phishletContentChanged.value = false;
  isEditingRating.value = false;
  directEditUpvotes.value = phishlet.upvotes;
  directEditDownvotes.value = phishlet.downvotes;
  hasUnsavedPhishletChanges.value = false;
  phishletNotesChanged.value = false;
  showPhishletDetailModal.value = true;
}

// Close phishlet detail modal
const closePhishletDetailModal = () => {
  showPhishletDetailModal.value = false
  selectedPhishlet.value = null
  isEditingPhishletContent.value = false
}

// Toggle editing phishlet content
const toggleEditPhishletContent = () => {
  isEditingPhishletContent.value = !isEditingPhishletContent.value
}

// Save edited phishlet content
const savePhishletContent = () => {
  if (!selectedPhishlet.value) return

  // Update the selected phishlet
  selectedPhishlet.value.content = editedPhishletContent.value

  // Find and update the phishlet in the list
  const index = phishlets.value.findIndex(p => p.id === selectedPhishlet.value.id)
  if (index !== -1) {
    phishlets.value[index].content = editedPhishletContent.value
  }

  // Exit edit mode
  isEditingPhishletContent.value = false
  phishletContentChanged.value = false
}

// Copy phishlet content to clipboard
const copyPhishletContent = async () => {
  if (!selectedPhishlet.value) return

  try {
    await navigator.clipboard.writeText(selectedPhishlet.value.content)
    clipboardContent.value = 'Phishlet content'
    showClipboardNotification.value = true
    setTimeout(() => {
      showClipboardNotification.value = false
    }, 2000)
  } catch (err) {
    console.error('Failed to copy phishlet content to clipboard:', err)
    alert('Failed to copy to clipboard')
  }
}

// Delete a phishlet
const deletePhishlet = (phishlet) => {
  if (confirm(`Are you sure you want to delete the phishlet "${phishlet.name}"?`)) {
    const index = phishlets.value.findIndex(p => p.id === phishlet.id)
    if (index !== -1) {
      phishlets.value.splice(index, 1)
    }
  }
}

// Update votes directly from the detail modal
const updatePhishletVotes = (isUpvote, increment) => {
  if (!selectedPhishlet.value) return

  if (isUpvote) {
    // If incrementing, add one; if decrementing, subtract one
    selectedPhishlet.value.upvotes += increment ? 1 : -1
    // Ensure votes don't go below zero
    if (selectedPhishlet.value.upvotes < 0) selectedPhishlet.value.upvotes = 0
  } else {
    // If incrementing, add one; if decrementing, subtract one
    selectedPhishlet.value.downvotes += increment ? 1 : -1
    // Ensure votes don't go below zero
    if (selectedPhishlet.value.downvotes < 0) selectedPhishlet.value.downvotes = 0
  }

  // Find and update the phishlet in the list
  const index = phishlets.value.findIndex(p => p.id === selectedPhishlet.value.id)
  if (index !== -1) {
    phishlets.value[index].upvotes = selectedPhishlet.value.upvotes
    phishlets.value[index].downvotes = selectedPhishlet.value.downvotes
  }
}

// Add a new ref for campaign artifacts
const campaignArtifacts = ref([])
const artifactFile = ref(null)
const artifactName = ref('')
const artifactDescription = ref('')
const artifactUploadSuccess = ref(false)
const artifactUploadError = ref('')

// Sample campaign artifacts
const sampleArtifacts = {
  1: [
    {
      id: 1,
      name: 'Employee List',
      description: 'CSV file with target employee data',
      type: 'text/csv',
      size: '25 KB',
      uploadDate: '2024-05-10',
      uploadedBy: 'admin'
    },
    {
      id: 2,
      name: 'Results Analysis',
      description: 'Campaign results analysis document',
      type: 'application/pdf',
      size: '128 KB',
      uploadDate: '2024-05-18',
      uploadedBy: 'admin'
    }
  ],
  2: [
    {
      id: 3,
      name: 'Executive Screenshot',
      description: 'Screenshot of landing page for executives',
      type: 'image/png',
      size: '215 KB',
      uploadDate: '2024-05-11',
      uploadedBy: 'analyst1'
    }
  ],
  3: []
}

// Handle artifact file selection
const handleArtifactFileChange = (event) => {
  const file = event.target.files[0]
  if (file) {
    artifactFile.value = file
    
    // If no name provided, use the file name
    if (!artifactName.value) {
      artifactName.value = file.name.split('.')[0]
    }
  }
}

// Upload artifact
const uploadArtifact = async () => {
  if (!artifactName.value) {
    artifactUploadError.value = 'Please enter an artifact name'
    return
  }
  
  if (!artifactFile.value) {
    artifactUploadError.value = 'Please select a file'
    return
  }
  
  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 800))
    
    // Add to the list
    const newArtifact = {
      id: Date.now(),
      name: artifactName.value,
      description: artifactDescription.value || 'No description provided',
      type: artifactFile.value.type || 'application/octet-stream',
      size: formatFileSize(artifactFile.value.size),
      uploadDate: new Date().toISOString().split('T')[0],
      uploadedBy: 'current_user' // In a real app, use the logged-in user
    }
    
    campaignArtifacts.value.unshift(newArtifact)
    
    // Save to the sample data
    if (selectedCampaign.value && selectedCampaign.value.id) {
      if (!sampleArtifacts[selectedCampaign.value.id]) {
        sampleArtifacts[selectedCampaign.value.id] = []
      }
      sampleArtifacts[selectedCampaign.value.id].unshift(newArtifact)
    }
    
    // Reset form
    artifactName.value = ''
    artifactDescription.value = ''
    artifactFile.value = null
    
    // Show success message
    artifactUploadSuccess.value = true
    setTimeout(() => {
      artifactUploadSuccess.value = false
    }, 3000)
    
  } catch (error) {
    console.error('Error uploading artifact:', error)
    artifactUploadError.value = 'Error uploading artifact. Please try again.'
  }
}

// Delete artifact
const deleteArtifact = (artifact) => {
  if (confirm(`Are you sure you want to delete the artifact "${artifact.name}"?`)) {
    // Remove from the current list
    campaignArtifacts.value = campaignArtifacts.value.filter(a => a.id !== artifact.id)
    
    // Remove from the sample data
    if (selectedCampaign.value && selectedCampaign.value.id) {
      sampleArtifacts[selectedCampaign.value.id] = sampleArtifacts[selectedCampaign.value.id].filter(a => a.id !== artifact.id)
    }
  }
}

// Format file size
const formatFileSize = (bytes) => {
  if (bytes < 1024) {
    return bytes + ' B'
  } else if (bytes < 1024 * 1024) {
    return (bytes / 1024).toFixed(0) + ' KB'
  } else if (bytes < 1024 * 1024 * 1024) {
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
  } else {
    return (bytes / (1024 * 1024 * 1024)).toFixed(2) + ' GB'
  }
}

// Add state for direct editing of upvotes/downvotes
const isEditingRating = ref(false)
const directEditUpvotes = ref(0)
const directEditDownvotes = ref(0)

// Save direct rating edits
const saveRatingEdits = () => {
  if (!selectedPhishlet.value) return
  
  // Ensure non-negative values
  directEditUpvotes.value = Math.max(0, parseInt(directEditUpvotes.value) || 0)
  directEditDownvotes.value = Math.max(0, parseInt(directEditDownvotes.value) || 0)
  
  // Update the selected phishlet
  selectedPhishlet.value.upvotes = directEditUpvotes.value
  selectedPhishlet.value.downvotes = directEditDownvotes.value
  
  // Find and update the phishlet in the list
  const index = phishlets.value.findIndex(p => p.id === selectedPhishlet.value.id)
  if (index !== -1) {
    phishlets.value[index].upvotes = directEditUpvotes.value
    phishlets.value[index].downvotes = directEditDownvotes.value
  }
  
  // Exit edit mode
  isEditingRating.value = false
}

// Add this function to update campaign info
const updateCampaignInfo = () => {
  if (!selectedCampaign.value) return
  
  // Update the campaign with edited values
  if (editedCampaignName.value) {
    selectedCampaign.value.name = editedCampaignName.value
  }
  
  if (editedTargetGroup.value) {
    selectedCampaign.value.target = editedTargetGroup.value
  }
  
  // In a real app, you would save this to the server
  // For demo purposes, we'll just update the local state
  const campaignIndex = phishingCampaigns.value.findIndex(c => c.id === selectedCampaign.value.id)
  if (campaignIndex !== -1) {
    phishingCampaigns.value[campaignIndex].name = editedCampaignName.value
    phishingCampaigns.value[campaignIndex].target = editedTargetGroup.value
  }
  
  // Exit edit mode
  isEditingCampaignInfo.value = false
}

// Replace with these functions
// Attempt to close the phishlet modal with confirmation if needed
const attemptClosePhishletModal = () => {
  const hasChanges = checkForUnsavedPhishletChanges();
  
  if (hasChanges) {
    // Show confirmation dialog
    showPhishletCloseConfirmation.value = true;
  } else {
    // No changes, just close
    actuallyClosePhishletModal();
  }
}

// The actual close modal function that doesn't check for changes
const actuallyClosePhishletModal = () => {
  showPhishletDetailModal.value = false;
  selectedPhishlet.value = null;
  originalPhishlet.value = null;
  showPhishletCloseConfirmation.value = false;
  hasUnsavedPhishletChanges.value = false;
  isEditingPhishletContent.value = false;
  isEditingRating.value = false;
}

// Close phishlet modal and discard changes
const discardPhishletChangesAndClose = () => {
  showPhishletCloseConfirmation.value = false;
  actuallyClosePhishletModal();
}

// Save all phishlet changes and close modal
const saveAllPhishletChangesAndClose = () => {
  // Save content if editing
  if (isEditingPhishletContent.value) {
    savePhishletContent();
  }
  
  // Save rating if editing
  if (isEditingRating.value) {
    saveRatingEdits();
  }
  
  // Save notes if changed
  if (phishletNotesChanged.value) {
    savePhishletNotes();
  }
  
  // Close confirmation and modal
  showPhishletCloseConfirmation.value = false;
  actuallyClosePhishletModal();
}

// Check for unsaved phishlet changes
const checkForUnsavedPhishletChanges = () => {
  if (!selectedPhishlet.value || !originalPhishlet.value) return false;
  
  // Check for rating changes
  if (isEditingRating.value &&
      (directEditUpvotes.value !== selectedPhishlet.value.upvotes ||
       directEditDownvotes.value !== selectedPhishlet.value.downvotes)) {
    console.log("Phishlet rating has changed");
    return true;
  }
  
  // Check for content changes
  if (isEditingPhishletContent.value && 
      editedPhishletContent.value !== selectedPhishlet.value.content) {
    console.log("Phishlet content has changed");
    return true;
  }
  
  // Check if upvotes or downvotes have changed from original
  if (selectedPhishlet.value.upvotes !== originalPhishlet.value.upvotes ||
      selectedPhishlet.value.downvotes !== originalPhishlet.value.downvotes) {
    console.log("Phishlet votes have changed");
    return true;
  }
  
  // Check if content has changed from original
  if (selectedPhishlet.value.content !== originalPhishlet.value.content) {
    console.log("Phishlet content has changed from original");
    return true;
  }
  
  // Check if notes have changed
  if (phishletNotes.value !== (originalPhishlet.value.notes || '')) {
    console.log("Phishlet notes have changed");
    return true;
  }
  
  return false;
}

// Add to existing watch section
// Update watches for phishlet content changes
watch(editedPhishletContent, (newValue) => {
  if (selectedPhishlet.value && originalPhishlet.value && isEditingPhishletContent.value) {
    if (newValue !== originalPhishlet.value.content) {
      hasUnsavedPhishletChanges.value = true;
    }
  }
});

// Watch for changes in rating
watch([directEditUpvotes, directEditDownvotes], () => {
  if (selectedPhishlet.value && originalPhishlet.value && isEditingRating.value) {
    if (directEditUpvotes.value !== originalPhishlet.value.upvotes || 
        directEditDownvotes.value !== originalPhishlet.value.downvotes) {
      hasUnsavedPhishletChanges.value = true;
    }
  }
});

// Watch for changes in phishlet notes
watch(phishletNotes, (newValue) => {
  if (selectedPhishlet.value && originalPhishlet.value) {
    const originalNotes = originalPhishlet.value.notes || '';
    if (newValue !== originalNotes) {
      hasUnsavedPhishletChanges.value = true;
      phishletNotesChanged.value = true;
    }
  }
});

// Add a function to download artifact
const downloadArtifact = (artifact) => {
  // In a real app, this would call the server to get the file
  // For demo purposes, we'll simulate the download with a blob
  
  // Create a sample blob based on the artifact type
  let blobContent = `Sample content for ${artifact.name}`;
  let blobType = artifact.type || 'text/plain';
  
  // For demo purposes, generate different content based on file type
  if (artifact.type.includes('csv')) {
    blobContent = `id,name,email\n1,John Doe,john@example.com\n2,Jane Smith,jane@example.com`;
  } else if (artifact.type.includes('pdf')) {
    blobContent = `%PDF-1.4\nSample PDF content for ${artifact.name}`;
  } else if (artifact.type.includes('image')) {
    // For images, we can't really create demo content, so we'll use text
    blobContent = `[This would be ${artifact.name} image data in a real app]`;
  }
  
  const blob = new Blob([blobContent], { type: blobType });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = artifact.name;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
  
  // Show notification
  clipboardContent.value = `${artifact.name} downloaded`;
  showClipboardNotification.value = true;
  setTimeout(() => {
    showClipboardNotification.value = false;
  }, 2000);
}

// Add a computed property for sorted phishlets
const sortedPhishlets = computed(() => {
  return [...phishlets.value].sort((a, b) => {
    // Sort alphabetically by name
    return a.name.localeCompare(b.name);
  });
});

// Remove the JSZip mock and downloadCampaignAsZip function
const JSZip = {
  // This is a mock implementation for demo purposes
  // In a real app, you would import the JSZip library
  new: function() {
    return {
      file: function(name, content) {
        console.log(`Adding file to zip: ${name}`);
        return this;
      },
      generateAsync: function(options) {
        return new Promise(resolve => {
          setTimeout(() => {
            resolve(new Blob(['Zip content here'], { type: 'application/zip' }));
          }, 500);
        });
      }
    };
  }
};

// Replace the downloadCampaignAsZip function
const downloadCampaignAsZip = async (campaign) => {
  // Show a loading notification
  clipboardContent.value = `Preparing zip package for "${campaign.name}"...`;
  showClipboardNotification.value = true;
  
  try {
    // In a real app, you would either:
    // 1. Import and use JSZip library client-side to create the zip
    // 2. Call an API endpoint to generate the zip server-side
    
    // Get the artifacts for this campaign
    const artifacts = sampleArtifacts[campaign.id] || [];
    
    // Create a new zip file
    const zip = JSZip.new();
    
    // Add campaign HTML template
    zip.file(`${campaign.name}_email_template.html`, campaign.htmlTemplate);
    
    // Add campaign YAML config
    zip.file(`${campaign.name}_config.yaml`, campaign.yamlConfig);
    
    // Add campaign Caddyfile
    zip.file(`${campaign.name}_caddyfile`, campaign.caddyfile);
    
    // Add campaign notes if available
    if (campaign.notes) {
      zip.file(`${campaign.name}_notes.txt`, campaign.notes);
    }
    
    // Add summary file with campaign details
    const summaryContent = `Campaign: ${campaign.name}
Target Group: ${campaign.target}
Status: ${campaign.status}
Click Rate: ${campaign.clickRate}
Submitted Credentials: ${campaign.submittedCredentials}
`;
    zip.file(`${campaign.name}_summary.txt`, summaryContent);
    
    // Add all artifacts
    artifacts.forEach(artifact => {
      // In a real app, you would fetch the actual file content from the server
      // For demo, we'll simulate different content based on file type
      let content = `Sample content for ${artifact.name}`;
      
      if (artifact.type.includes('csv')) {
        content = `id,name,email\n1,John Doe,john@example.com\n2,Jane Smith,jane@example.com`;
      } else if (artifact.type.includes('pdf')) {
        content = `%PDF-1.4\nSample PDF content for ${artifact.name}`;
      } else if (artifact.type.includes('image')) {
        content = `[This would be ${artifact.name} image data in a real app]`;
      }
      
      // Add to zip in an artifacts folder
      zip.file(`artifacts/${artifact.name}`, content);
    });
    
    // Generate the zip file asynchronously
    const zipBlob = await zip.generateAsync({ type: 'blob' });
    
    // Create download link and trigger download
    const url = URL.createObjectURL(zipBlob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${campaign.name.replace(/\s+/g, '_')}_package.zip`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    
    // Update notification
    clipboardContent.value = `Campaign "${campaign.name}" package downloaded`;
    setTimeout(() => {
      showClipboardNotification.value = false;
    }, 2000);
    
  } catch (error) {
    console.error('Error creating campaign zip package:', error);
    clipboardContent.value = `Error creating package for "${campaign.name}"`;
    setTimeout(() => {
      showClipboardNotification.value = false;
    }, 2000);
  }
};

// Add this new ref for project-specific metrics
const projectMetrics = ref({
  1: [ // Campaign ID 1
    { 
      projectName: "Marketing Q2 - HQ Team",
      sent: 30,
      clicked: 5,
      submitted: 3,
      clickRate: "16.7%",
      submissionRate: "10%",
      date: "2024-04-12"
    },
    { 
      projectName: "Marketing Q2 - Remote Team",
      sent: 45,
      clicked: 12,
      submitted: 6,
      clickRate: "26.7%",
      submissionRate: "13.3%",
      date: "2024-04-15"
    }
  ],
  2: [ // Campaign ID 2
    { 
      projectName: "Executive Assessment - US",
      sent: 8,
      clicked: 1,
      submitted: 0,
      clickRate: "12.5%",
      submissionRate: "0%",
      date: "2024-05-05"
    },
    { 
      projectName: "Executive Assessment - EMEA",
      sent: 7,
      clicked: 2,
      submitted: 1,
      clickRate: "28.6%",
      submissionRate: "14.3%",
      date: "2024-05-07"
    }
  ],
  3: [] // Campaign ID 3 - no projects yet
});

// Calculate aggregate metrics for a campaign
const calculateAggregateMetrics = (campaignId) => {
  const projects = projectMetrics.value[campaignId] || [];
  if (projects.length === 0) return { clickRate: 'N/A', submissionRate: 'N/A' };
  
  let totalSent = 0;
  let totalClicked = 0;
  let totalSubmitted = 0;
  
  projects.forEach(project => {
    totalSent += project.sent;
    totalClicked += project.clicked;
    totalSubmitted += project.submitted;
  });
  
  const clickRate = totalSent > 0 ? `${Math.round((totalClicked / totalSent) * 100)}%` : 'N/A';
  const submissionRate = totalSent > 0 ? `${Math.round((totalSubmitted / totalSent) * 100)}%` : 'N/A';
  
  return { clickRate, submissionRate, totalSent, totalClicked, totalSubmitted };
};

// Add a function to add a new project to a campaign
const newProjectName = ref('');
const newProjectSent = ref('');
const newProjectClicked = ref('');
const newProjectSubmitted = ref('');
const showAddProjectForm = ref(false);

const addProjectToCampaign = () => {
  if (!selectedCampaign.value) return;
  
  if (!newProjectName.value) {
    alert('Please enter a project name');
    return;
  }
  
  const sent = parseInt(newProjectSent.value) || 0;
  const clicked = parseInt(newProjectClicked.value) || 0;
  const submitted = parseInt(newProjectSubmitted.value) || 0;
  
  if (sent <= 0) {
    alert('Please enter a valid number of emails sent');
    return;
  }
  
  if (clicked > sent) {
    alert('Clicked count cannot be greater than sent count');
    return;
  }
  
  if (submitted > clicked) {
    alert('Submitted count cannot be greater than clicked count');
    return;
  }
  
  // Calculate rates
  const clickRate = `${Math.round((clicked / sent) * 100)}%`;
  const submissionRate = `${Math.round((submitted / sent) * 100)}%`;
  
  // Create new project
  const newProject = {
    projectName: newProjectName.value,
    sent,
    clicked,
    submitted,
    clickRate,
    submissionRate,
    date: new Date().toISOString().split('T')[0]
  };
  
  // Add to the appropriate campaign
  if (!projectMetrics.value[selectedCampaign.value.id]) {
    projectMetrics.value[selectedCampaign.value.id] = [];
  }
  
  projectMetrics.value[selectedCampaign.value.id].push(newProject);
  
  // Update the campaign's aggregate metrics
  const metrics = calculateAggregateMetrics(selectedCampaign.value.id);
  selectedCampaign.value.clickRate = metrics.clickRate;
  selectedCampaign.value.submittedCredentials = metrics.submissionRate;
  
  // Update the campaigns list
  const campaignIndex = phishingCampaigns.value.findIndex(c => c.id === selectedCampaign.value.id);
  if (campaignIndex !== -1) {
    phishingCampaigns.value[campaignIndex].clickRate = metrics.clickRate;
    phishingCampaigns.value[campaignIndex].submittedCredentials = metrics.submissionRate;
  }
  
  // Reset form
  newProjectName.value = '';
  newProjectSent.value = '';
  newProjectClicked.value = '';
  newProjectSubmitted.value = '';
  showAddProjectForm.value = false;
};

const copyHtmlToClipboard = async () => {
  if (!selectedCampaign.value) return
  
  try {
    await navigator.clipboard.writeText(selectedCampaign.value.htmlTemplate)
    clipboardContent.value = 'HTML template'
    showClipboardNotification.value = true
    setTimeout(() => {
      showClipboardNotification.value = false
    }, 2000)
  } catch (err) {
    console.error('Failed to copy HTML to clipboard:', err)
    alert('Failed to copy to clipboard')
  }
}

// Add function to save phishlet notes
const savePhishletNotes = () => {
  if (!selectedPhishlet.value) return;
  
  // Update the selected phishlet
  selectedPhishlet.value.notes = phishletNotes.value;
  
  // Find and update the phishlet in the list
  const index = phishlets.value.findIndex(p => p.id === selectedPhishlet.value.id);
  if (index !== -1) {
    phishlets.value[index].notes = phishletNotes.value;
  }
  
  // Reset change tracking
  phishletNotesChanged.value = false;
  
  // Update original phishlet reference
  if (originalPhishlet.value) {
    originalPhishlet.value.notes = phishletNotes.value;
  }

  // Show notification
  clipboardContent.value = 'Phishlet notes saved';
  showClipboardNotification.value = true;
  setTimeout(() => {
    showClipboardNotification.value = false;
  }, 2000);
};
</script>

<template>
  <div class="min-h-screen flex flex-col bg-gray-900">
    <Header />
    <main class="flex-grow">
      <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <!-- Page Description -->
        <div class="mb-8">
          <h1 class="text-2xl font-bold text-white mb-2">Phishing Campaign Management</h1>
          <p class="text-gray-400">
            Manage phishing campaigns for security awareness training. Plan, execute, and track phishing simulations to test and improve your organization's security awareness.
          </p>
        </div>

        <!-- Create New Campaign Section -->
        <div class="bg-gray-800 shadow rounded-lg mb-8 p-6">
          <h2 class="text-xl font-semibold text-white mb-4">Create New Campaign</h2>
          
          <div class="space-y-4">
            <!-- Campaign Name -->
            <div>
              <label for="campaignName" class="block text-sm font-medium text-white">Campaign Name</label>
              <input
                type="text"
                id="campaignName"
                v-model="campaignName"
                class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-500 focus:ring-opacity-50 text-white"
                placeholder="Enter campaign name"
              />
            </div>
            
            <!-- Target Group -->
            <div>
              <label for="targetGroup" class="block text-sm font-medium text-white">Target Group</label>
              <input
                type="text"
                id="targetGroup"
                v-model="targetGroup"
                class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-500 focus:ring-opacity-50 text-white"
                placeholder="e.g., Marketing Team, All Users"
              />
            </div>
            
            <!-- HTML Template Upload -->
            <div>
              <label for="htmlTemplate" class="block text-sm font-medium text-white">HTML Email Template</label>
              <input
                type="file"
                id="htmlTemplate"
                @change="handleHtmlTemplateChange"
                accept=".html"
                class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-500 focus:ring-opacity-50 text-white"
              />
              <p class="mt-1 text-xs text-gray-400">Upload your HTML email template file</p>
            </div>
            
            <!-- YAML Config Upload -->
            <div>
              <label for="yamlConfig" class="block text-sm font-medium text-white">YAML Configuration</label>
              <input
                type="file"
                id="yamlConfig"
                @change="handleYamlConfigChange"
                accept=".yaml,.yml"
                class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 shadow-sm focus:border-blue-500 focus:ring focus:ring-blue-500 focus:ring-opacity-50 text-white"
              />
              <p class="mt-1 text-xs text-gray-400">Upload your YAML configuration file</p>
            </div>
            
            <!-- Error Messages -->
            <div v-if="errorMessage" class="text-red-400 text-sm">{{ errorMessage }}</div>
            
            <!-- Success Message -->
            <div v-if="uploadSuccess" class="text-green-400 text-sm">Campaign created successfully!</div>
            
            <!-- Submit Button -->
            <div>
              <button
                @click="createCampaign"
                class="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
              >
                Create Campaign
              </button>
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
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider w-1/4">
                    Campaign Name
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider w-1/4">
                    Target Group
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider w-1/4">
                    Click Rate
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider w-1/4">
                    Submitted Credentials
                  </th>
                </tr>
              </thead>
              <tbody class="bg-gray-800 divide-y divide-gray-700">
                <tr 
                  v-for="campaign in phishingCampaigns" 
                  :key="campaign.id" 
                  class="hover:bg-gray-700 cursor-pointer"
                  @click="viewCampaignDetails(campaign)"
                >
                  <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-white">
                    {{ campaign.name }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300">
                    {{ campaign.target }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300">
                    {{ campaign.clickRate }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300">
                    {{ campaign.submittedCredentials }}
                  </td>
                </tr>
                <tr v-if="phishingCampaigns.length === 0">
                  <td colspan="4" class="px-6 py-4 text-center text-sm text-gray-400">
                    No campaigns found
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Phishing Tools Section -->
        <div>
          <h2 class="text-xl font-semibold text-white mb-4">Phishlet Tracking</h2>
          <div class="bg-gray-800 shadow overflow-hidden rounded-lg mb-4">
            <div class="p-4 flex justify-between items-center">
              <div>
                <p class="text-gray-400 text-sm">Track and manage standalone phishlets for use in your campaigns</p>
              </div>
              <button
                @click="showPhishletUploadModal = true"
                class="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
              >
                Upload Phishlet
              </button>
            </div>
            
            <table class="min-w-full divide-y divide-gray-700">
              <thead class="bg-gray-700">
                <tr>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider w-1/3">
                    Phishlet Name
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider w-1/3">
                    Upload Date
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider w-1/3">
                    Uploaded By
                  </th>
                  <th scope="col" class="px-6 py-3 text-right text-xs font-medium text-gray-300 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody class="bg-gray-800 divide-y divide-gray-700">
                <tr v-for="phishlet in sortedPhishlets" :key="phishlet.id" class="hover:bg-gray-700">
                  <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-white cursor-pointer" @click="viewPhishletDetails(phishlet)">
                    {{ phishlet.name }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300 cursor-pointer" @click="viewPhishletDetails(phishlet)">
                    {{ phishlet.uploadDate }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300 cursor-pointer" @click="viewPhishletDetails(phishlet)">
                    {{ phishlet.uploadedBy }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <button @click="viewPhishletDetails(phishlet)" class="text-blue-400 hover:text-blue-300 mr-3">
                      View
                    </button>
                    <button @click="deletePhishlet(phishlet)" class="text-red-400 hover:text-red-300">
                      Delete
                    </button>
                  </td>
                </tr>
                <tr v-if="phishlets.length === 0">
                  <td colspan="4" class="px-6 py-4 text-center text-sm text-gray-400">
                    No phishlets found
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </main>
    
    <!-- Campaign Details Modal -->
    <div v-if="showCampaignModal" class="fixed inset-0 z-50 overflow-y-auto">
      <!-- Modal Backdrop - Make it more translucent -->
      <div class="fixed inset-0 bg-black bg-opacity-60 backdrop-blur-sm transition-opacity" @click="attemptCloseModal"></div>
      
      <!-- Modal Content -->
      <div class="flex items-center justify-center min-h-screen p-4">
        <div class="bg-gray-800 rounded-lg overflow-hidden shadow-xl transform transition-all sm:max-w-6xl sm:w-full relative mx-auto my-8">
          <!-- Modal Header -->
          <div class="bg-gray-700 px-6 py-4 flex items-center justify-between">
            <h3 class="text-lg font-medium text-white">
              Campaign Details: {{ selectedCampaign?.name }}
            </h3>
            <button 
              @click="attemptCloseModal" 
              class="text-gray-400 hover:text-white focus:outline-none"
            >
              <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          
          <!-- Modal Body -->
          <div class="p-6 bg-gray-800">
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
                      {{ selectedCampaign?.name }}
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
                      {{ selectedCampaign?.target }}
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
                        'bg-green-100 text-green-800': selectedCampaign?.status === 'Completed',
                        'bg-yellow-100 text-yellow-800': selectedCampaign?.status === 'In Progress',
                        'bg-blue-100 text-blue-800': selectedCampaign?.status === 'Scheduled'
                      }">
                        {{ selectedCampaign?.status }}
                      </span>
                    </div>
                    
                    <div class="text-sm text-gray-400">Click Rate:</div>
                    <div v-if="!isEditingMetrics" class="text-sm text-white">
                      {{ selectedCampaign?.clickRate }}
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
                      {{ selectedCampaign?.submittedCredentials }}
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
                  
                  <!-- Save button for editing campaign info -->
                  <div v-if="isEditingCampaignInfo" class="mt-3 flex justify-end">
                    <button
                      @click="updateCampaignInfo"
                      class="bg-blue-600 hover:bg-blue-700 text-white text-xs font-medium py-1 px-3 rounded"
                    >
                      Save Campaign Changes
                    </button>
                  </div>
                  
                  <!-- Save button for editing metrics -->
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
                    <pre class="text-sm text-gray-300 whitespace-pre-wrap font-mono">{{ selectedCampaign?.yamlConfig }}</pre>
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
                    <pre class="text-sm text-gray-300 whitespace-pre-wrap font-mono">{{ selectedCampaign?.caddyfile }}</pre>
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
                    @change="handleNotesChange"
                    class="w-full bg-gray-700 border-gray-600 rounded-md text-white p-3"
                    style="height: 400px; resize: none;"
                    placeholder="Add campaign notes, observations, or action items here..."
                  ></textarea>
                </div>
              </div>
            </div>
            
            <!-- Campaign Metrics Section - Moved Above Artifacts and Made Wider -->
            <div class="mt-6 mb-6">
              <div class="bg-gray-700 rounded-lg p-4">
                <div class="flex justify-between items-center mb-2">
                  <h4 class="text-md font-medium text-white">Campaign Metrics</h4>
                  <div>
                    <button 
                      @click="showAddProjectForm = !showAddProjectForm" 
                      class="bg-blue-600 hover:bg-blue-700 text-white text-xs font-medium py-1 px-3 rounded"
                    >
                      {{ showAddProjectForm ? 'Cancel' : 'Add Project' }}
                    </button>
                  </div>
                </div>
                
                <!-- Aggregate Metrics -->
                <div class="bg-gray-800 p-3 rounded mb-3">
                  <h5 class="text-sm font-medium text-white mb-2">Aggregate Statistics</h5>
                  <div class="grid grid-cols-3 gap-4 text-center">
                    <div class="bg-gray-700 rounded p-2">
                      <div class="text-gray-400 text-xs">Total Sent</div>
                      <div class="text-white text-lg font-bold">{{ selectedCampaign.totalSent || 0 }}</div>
                    </div>
                    <div class="bg-gray-700 rounded p-2">
                      <div class="text-gray-400 text-xs">Total Clicked</div>
                      <div class="text-white text-lg font-bold">{{ selectedCampaign.totalClicked || 0 }}</div>
                      <div class="text-blue-400 text-xs">{{ selectedCampaign.clickRate }}</div>
                    </div>
                    <div class="bg-gray-700 rounded p-2">
                      <div class="text-gray-400 text-xs">Total Submitted</div>
                      <div class="text-white text-lg font-bold">{{ selectedCampaign.totalSubmitted || 0 }}</div>
                      <div class="text-blue-400 text-xs">{{ selectedCampaign.submittedCredentials }}</div>
                    </div>
                  </div>
                </div>
                
                <!-- Add Project Form -->
                <div v-if="showAddProjectForm" class="bg-gray-800 p-3 rounded mb-3">
                  <h5 class="text-sm font-medium text-white mb-2">Add New Project</h5>
                  <div class="space-y-2">
                    <div>
                      <label class="block text-xs text-gray-400 mb-1">Project Name</label>
                      <input 
                        v-model="newProjectName" 
                        type="text" 
                        class="w-full rounded-md bg-gray-700 border-gray-600 text-white text-sm py-1 px-2"
                        placeholder="e.g., Marketing Q2 - New York Office"
                      />
                    </div>
                    <div class="grid grid-cols-3 gap-2">
                      <div>
                        <label class="block text-xs text-gray-400 mb-1">Emails Sent</label>
                        <input 
                          v-model="newProjectSent" 
                          type="number" 
                          min="1"
                          class="w-full rounded-md bg-gray-700 border-gray-600 text-white text-sm py-1 px-2"
                        />
                      </div>
                      <div>
                        <label class="block text-xs text-gray-400 mb-1">Emails Clicked</label>
                        <input 
                          v-model="newProjectClicked" 
                          type="number" 
                          min="0"
                          :max="newProjectSent"
                          class="w-full rounded-md bg-gray-700 border-gray-600 text-white text-sm py-1 px-2"
                        />
                      </div>
                      <div>
                        <label class="block text-xs text-gray-400 mb-1">Credentials Submitted</label>
                        <input 
                          v-model="newProjectSubmitted" 
                          type="number" 
                          min="0"
                          :max="newProjectClicked"
                          class="w-full rounded-md bg-gray-700 border-gray-600 text-white text-sm py-1 px-2"
                        />
                      </div>
                    </div>
                    <div class="flex justify-end">
                      <button 
                        @click="addProjectToCampaign" 
                        class="bg-green-600 hover:bg-green-700 text-white text-xs font-medium py-1 px-3 rounded"
                      >
                        Add Project
                      </button>
                    </div>
                  </div>
                </div>
                
                <!-- Project-specific Metrics Table -->
                <div class="overflow-x-auto">
                  <table class="min-w-full divide-y divide-gray-600">
                    <thead class="bg-gray-800">
                      <tr>
                        <th scope="col" class="px-3 py-2 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Project</th>
                        <th scope="col" class="px-3 py-2 text-center text-xs font-medium text-gray-400 uppercase tracking-wider">Date</th>
                        <th scope="col" class="px-3 py-2 text-center text-xs font-medium text-gray-400 uppercase tracking-wider">Sent</th>
                        <th scope="col" class="px-3 py-2 text-center text-xs font-medium text-gray-400 uppercase tracking-wider">Clicked</th>
                        <th scope="col" class="px-3 py-2 text-center text-xs font-medium text-gray-400 uppercase tracking-wider">Submitted</th>
                        <th scope="col" class="px-3 py-2 text-center text-xs font-medium text-gray-400 uppercase tracking-wider">Click Rate</th>
                        <th scope="col" class="px-3 py-2 text-center text-xs font-medium text-gray-400 uppercase tracking-wider">Sub. Rate</th>
                      </tr>
                    </thead>
                    <tbody class="bg-gray-800 divide-y divide-gray-700">
                      <tr v-for="(project, index) in projectMetrics[selectedCampaign.id] || []" :key="index">
                        <td class="px-3 py-2 whitespace-nowrap text-xs text-white">{{ project.projectName }}</td>
                        <td class="px-3 py-2 whitespace-nowrap text-xs text-gray-300 text-center">{{ project.date }}</td>
                        <td class="px-3 py-2 whitespace-nowrap text-xs text-gray-300 text-center">{{ project.sent }}</td>
                        <td class="px-3 py-2 whitespace-nowrap text-xs text-gray-300 text-center">{{ project.clicked }}</td>
                        <td class="px-3 py-2 whitespace-nowrap text-xs text-gray-300 text-center">{{ project.submitted }}</td>
                        <td class="px-3 py-2 whitespace-nowrap text-xs text-green-400 text-center">{{ project.clickRate }}</td>
                        <td class="px-3 py-2 whitespace-nowrap text-xs text-red-400 text-center">{{ project.submissionRate }}</td>
                      </tr>
                      <tr v-if="!projectMetrics[selectedCampaign.id] || projectMetrics[selectedCampaign.id].length === 0">
                        <td colspan="7" class="px-3 py-2 text-center text-xs text-gray-400">
                          No project metrics available. Add a project to track metrics.
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
            
            <!-- Misc Artifact Uploads Section - Now Below Metrics -->
            <div>
              <div class="flex justify-between items-center mb-2">
                <h4 class="text-md font-medium text-white">Campaign Artifacts</h4>
                <!-- No button here -->
              </div>
              
              <!-- Artifact Upload Form -->
              <div class="bg-gray-700 rounded-lg p-4 mb-4">
                <div class="grid grid-cols-1 gap-3">
                  <div>
                    <label for="artifactName" class="block text-sm font-medium text-white mb-1">Artifact Name</label>
                    <input
                      type="text"
                      id="artifactName"
                      v-model="artifactName"
                      class="w-full rounded-md bg-gray-600 border-gray-500 text-white text-sm py-1 px-2"
                      placeholder="Enter artifact name"
                    />
                  </div>
                  
                  <div>
                    <label for="artifactDescription" class="block text-sm font-medium text-white mb-1">Description</label>
                    <input
                      type="text"
                      id="artifactDescription"
                      v-model="artifactDescription"
                      class="w-full rounded-md bg-gray-600 border-gray-500 text-white text-sm py-1 px-2"
                      placeholder="Enter artifact description"
                    />
                  </div>
                  
                  <div>
                    <label for="artifactFile" class="block text-sm font-medium text-white mb-1">File</label>
                    <input
                      type="file"
                      id="artifactFile"
                      @change="handleArtifactFileChange"
                      class="w-full rounded-md bg-gray-600 border-gray-500 text-white text-sm py-1 px-2"
                    />
                  </div>
                  
                  <div v-if="artifactUploadSuccess" class="text-green-400 text-xs">
                    Artifact uploaded successfully!
                  </div>
                  
                  <div v-if="artifactUploadError" class="text-red-400 text-xs">
                    {{ artifactUploadError }}
                  </div>
                  
                  <div>
                    <button
                      @click="uploadArtifact"
                      class="bg-blue-600 hover:bg-blue-700 text-white text-xs font-medium py-1 px-3 rounded"
                    >
                      Upload Artifact
                    </button>
                  </div>
                </div>
              </div>
              
              <!-- Artifacts List -->
              <div class="bg-gray-700 rounded-lg overflow-hidden">
                <table class="min-w-full divide-y divide-gray-600">
                  <thead class="bg-gray-800">
                    <tr>
                      <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">Name</th>
                      <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">Description</th>
                      <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">Type</th>
                      <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">Size</th>
                      <th scope="col" class="px-6 py-3 text-right text-xs font-medium text-gray-300 uppercase tracking-wider">Actions</th>
                    </tr>
                  </thead>
                  <tbody class="bg-gray-700 divide-y divide-gray-600">
                    <tr v-for="artifact in campaignArtifacts" :key="artifact.id">
                      <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-white">{{ artifact.name }}</td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{{ artifact.description }}</td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{{ artifact.type }}</td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{{ artifact.size }}</td>
                      <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                        <div class="flex justify-end space-x-2">
                          <button 
                            @click="downloadArtifact(artifact)" 
                            class="bg-blue-600 hover:bg-blue-700 text-white text-xs font-medium py-1 px-2 rounded"
                            title="Download this artifact"
                          >
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                              <path fill-rule="evenodd" d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z" clip-rule="evenodd" />
                            </svg>
                          </button>
                          <button 
                            @click="deleteArtifact(artifact)" 
                            class="bg-red-600 hover:bg-red-700 text-white text-xs font-medium py-1 px-2 rounded"
                            title="Delete this artifact"
                          >
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                              <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
                            </svg>
                          </button>
                        </div>
                      </td>
                    </tr>
                    <tr v-if="campaignArtifacts.length === 0">
                      <td colspan="5" class="px-6 py-4 text-center text-sm text-gray-400">
                        No artifacts found for this campaign
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Modal Footer -->
        <div class="bg-gray-700 px-6 py-4 flex justify-end">
          <button 
            @click="closeCampaignModal" 
            class="bg-gray-600 hover:bg-gray-500 text-white font-medium py-2 px-4 rounded"
          >
            Close
          </button>
        </div>
      </div>
    </div>
    
    <!-- Confirmation Modal for Closing with Unsaved Changes -->
    <div v-if="showCloseConfirmation" class="fixed inset-0 z-60 overflow-y-auto">
      <div class="fixed inset-0 bg-black bg-opacity-60 backdrop-blur-sm transition-opacity"></div>
      <div class="flex items-center justify-center min-h-screen p-4">
        <div class="bg-gray-800 rounded-lg overflow-hidden shadow-xl transform transition-all sm:max-w-lg sm:w-full relative mx-auto">
          <div class="bg-gray-700 px-6 py-4">
            <h3 class="text-lg font-medium text-white">Unsaved Changes</h3>
          </div>
          <div class="p-6 bg-gray-800">
            <p class="text-white">You have unsaved changes. What would you like to do?</p>
          </div>
          <div class="bg-gray-700 px-6 py-4 flex justify-end space-x-3">
            <button 
              @click="discardChangesAndClose" 
              class="bg-red-600 hover:bg-red-700 text-white font-medium py-2 px-4 rounded"
            >
              Discard Changes
            </button>
            <button 
              @click="saveAllChangesAndClose" 
              class="bg-green-600 hover:bg-green-700 text-white font-medium py-2 px-4 rounded"
            >
              Save All Changes
            </button>
            <button 
              @click="showCloseConfirmation = false" 
              class="bg-gray-600 hover:bg-gray-500 text-white font-medium py-2 px-4 rounded"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Clipboard Notification -->
    <div v-if="showClipboardNotification" class="fixed bottom-4 right-4 bg-gray-900 text-white px-4 py-2 rounded shadow-lg z-50">
      {{ clipboardContent }} copied to clipboard
    </div>
    
    <!-- Phishlet Details Modal -->
    <div v-if="showPhishletDetailModal" class="fixed inset-0 z-50 overflow-y-auto">
      <!-- Modal Backdrop -->
      <div class="fixed inset-0 bg-black bg-opacity-60 backdrop-blur-sm transition-opacity" @click="attemptClosePhishletModal"></div>
      
      <!-- Modal Content -->
      <div class="flex items-center justify-center min-h-screen p-4">
        <div class="bg-gray-800 rounded-lg overflow-hidden shadow-xl transform transition-all sm:max-w-5xl sm:w-full relative mx-auto my-8">
          <!-- Modal Header -->
          <div class="bg-gray-700 px-6 py-4 flex items-center justify-between">
            <h3 class="text-lg font-medium text-white">
              Phishlet Details: {{ selectedPhishlet?.name }}
            </h3>
            <button 
              @click="attemptClosePhishletModal" 
              class="text-gray-400 hover:text-white focus:outline-none"
            >
              <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          
          <!-- Modal Body -->
          <div class="p-6 bg-gray-800">
            <!-- Phishlet Information -->
            <div class="mb-6">
              <h4 class="text-md font-medium text-white mb-3">Phishlet Information</h4>
              <div class="grid grid-cols-2 gap-2 bg-gray-700 rounded-lg p-4">
                <div class="text-sm text-gray-400">Uploaded By:</div>
                <div class="text-sm text-white">{{ selectedPhishlet?.uploadedBy }}</div>
                
                <div class="text-sm text-gray-400">Upload Date:</div>
                <div class="text-sm text-white">{{ selectedPhishlet?.uploadDate }}</div>
                
            </div>
            </div>
            
            <!-- Phishlet Content -->
            <div>
              <div class="flex justify-between items-center mb-2">
                <h4 class="text-md font-medium text-white">Phishlet Content</h4>
                <div class="space-x-2">
                  <button 
                    v-if="!isEditingPhishletContent"
                    @click="copyPhishletContent" 
                    class="bg-blue-600 hover:bg-blue-700 text-white text-xs font-medium py-1 px-2 rounded"
                  >
                    Copy
                  </button>
                  <button 
                    v-if="!isEditingPhishletContent"
                    @click="toggleEditPhishletContent" 
                    class="bg-yellow-600 hover:bg-yellow-700 text-white text-xs font-medium py-1 px-2 rounded"
                  >
                    Edit
                  </button>
                  <button 
                    v-if="isEditingPhishletContent"
                    @click="savePhishletContent" 
                    class="bg-green-600 hover:bg-green-700 text-white text-xs font-medium py-1 px-2 rounded"
                  >
                    Save
                  </button>
                  <button 
                    v-if="isEditingPhishletContent"
                    @click="toggleEditPhishletContent" 
                    class="bg-red-600 hover:bg-red-700 text-white text-xs font-medium py-1 px-2 rounded"
                  >
                    Cancel
                  </button>
                </div>
              </div>
              <textarea
                v-model="editedPhishletContent"
                :readonly="!isEditingPhishletContent"
                @input="phishletContentChanged = true"
                class="w-full bg-gray-700 border-gray-600 rounded-md text-white p-3 font-mono"
                rows="15"
                style="resize: vertical;"
              ></textarea>
            </div>
            
            <!-- Phishlet Notes -->
            <div class="mt-4">
              <div class="flex justify-between items-center mb-2">
                <h4 class="text-md font-medium text-white">Phishlet Notes</h4>
                <button 
                  @click="savePhishletNotes" 
                  class="bg-blue-600 hover:bg-blue-700 text-white text-xs font-medium py-1 px-3 rounded"
                >
                  Save Notes
                </button>
              </div>
              <textarea
                v-model="phishletNotes"
                @input="phishletNotesChanged = true"
                class="w-full bg-gray-700 border-gray-600 rounded-md text-white p-3"
                style="height: 400px; resize: none;"
                placeholder="Add notes about this phishlet..."
              ></textarea>
            </div>
          </div>
          
        </div>
      </div>
    </div>
    
    <!-- Phishlet Close Confirmation Modal -->
    <div v-if="showPhishletCloseConfirmation" class="fixed inset-0 z-60 overflow-y-auto">
      <div class="fixed inset-0 bg-black bg-opacity-60 backdrop-blur-sm transition-opacity"></div>
      <div class="flex items-center justify-center min-h-screen p-4">
        <div class="bg-gray-800 rounded-lg overflow-hidden shadow-xl transform transition-all sm:max-w-lg sm:w-full relative mx-auto">
          <div class="bg-gray-700 px-6 py-4">
            <h3 class="text-lg font-medium text-white">Unsaved Phishlet Changes</h3>
          </div>
          <div class="p-6 bg-gray-800">
            <p class="text-white">You have unsaved changes to this phishlet. What would you like to do?</p>
          </div>
          <div class="bg-gray-700 px-6 py-4 flex justify-end space-x-3">
            <button 
              @click="discardPhishletChangesAndClose" 
              class="bg-red-600 hover:bg-red-700 text-white font-medium py-2 px-4 rounded"
            >
              Discard Changes
            </button>
            <button 
              @click="saveAllPhishletChangesAndClose" 
              class="bg-green-600 hover:bg-green-700 text-white font-medium py-2 px-4 rounded"
            >
              Save All Changes
            </button>
            <button 
              @click="showPhishletCloseConfirmation = false" 
              class="bg-gray-600 hover:bg-gray-500 text-white font-medium py-2 px-4 rounded"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Add any component-specific styles here */
th:nth-child(4):has(content="Duration"), 
td:nth-child(4):has(content*="to") {
  display: none !important;
}

pre {
  word-wrap: break-word;
}
</style> 