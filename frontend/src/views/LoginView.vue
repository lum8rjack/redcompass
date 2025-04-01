<template>
  <div class="min-h-screen flex flex-col bg-gray-900">
    <div class="flex-grow flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div class="max-w-md w-full space-y-8 bg-gray-800 p-8 rounded-lg shadow">
        <div>
          <img class="mx-auto h-32 w-auto" src="../assets/redcompass.webp" alt="RedCompass" />
          <h2 class="mt-6 text-center text-3xl font-extrabold text-white">
            {{ showForgotPassword ? 'Reset your password' : 'Sign in to your account' }}
          </h2>
        </div>

        <!-- Forgot Password Form -->
        <div v-if="showForgotPassword">
          <form class="mt-8 space-y-6" @submit.prevent="handleResetPassword" novalidate>
            <div class="space-y-4">
              <div>
                <label for="email-address" class="sr-only">Email address</label>
                <input
                  id="email-address"
                  name="email"
                  type="email"
                  required
                  class="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-600 placeholder-gray-400 text-white bg-gray-700 focus:outline-none focus:ring-gray-400 focus:border-gray-400 focus:z-10 sm:text-sm"
                  placeholder="Email address"
                  v-model="email"
                  @blur="validateEmail"
                />
              </div>
              <div v-if="emailError" class="text-sm text-white text-center">
                {{ emailError }}
              </div>
            </div>

            <div class="flex items-center justify-between space-x-3 mt-8">
              <button
                type="button"
                @click="handleBackToLogin"
                class="w-1/2 py-2 px-4 border border-gray-600 text-sm font-medium rounded-md text-white bg-transparent hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-400"
              >
                Back to login
              </button>
              <button
                type="submit"
                class="w-1/2 py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-gray-600 hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-400"
              >
                Reset Password
              </button>
            </div>
          </form>
        </div>

        <!-- Login Form -->
        <div v-else>
          <!-- Username/Password Form -->
          <div v-if="authMethods?.usernamePassword">
            <form class="mt-8 space-y-6" @submit.prevent novalidate>
              <div class="space-y-4">
                <div class="rounded-md shadow-sm -space-y-px">
                  <div>
                    <label for="email-address" class="sr-only">Email address</label>
                    <input
                      id="email-address"
                      name="email"
                      type="email"
                      required
                      class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-600 placeholder-gray-400 text-white bg-gray-700 rounded-t-md focus:outline-none focus:ring-gray-400 focus:border-gray-400 focus:z-10 sm:text-sm"
                      placeholder="Email address"
                      v-model="email"
                      @blur="validateEmail"
                    />
                  </div>
                  <div>
                    <label for="password" class="sr-only">Password</label>
                    <input
                      id="password"
                      name="password"
                      type="password"
                      required
                      class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-600 placeholder-gray-400 text-white bg-gray-700 rounded-b-md focus:outline-none focus:ring-gray-400 focus:border-gray-400 focus:z-10 sm:text-sm"
                      placeholder="Password"
                      v-model="password"
                    />
                  </div>
                </div>

                <div v-if="emailError" class="text-sm text-white text-center">
                  {{ emailError }}
                </div>
              </div>

              <div class="flex items-center justify-between">
                <div class="flex items-center">
                  <input
                    id="remember-me"
                    name="remember-me"
                    type="checkbox"
                    class="h-4 w-4 text-gray-600 focus:ring-gray-400 border-gray-600 bg-gray-700 rounded"
                    v-model="rememberMe"
                  />
                  <label for="remember-me" class="ml-2 block text-sm text-gray-300">
                    Remember me
                  </label>
                </div>

                <div class="text-sm">
                  <a 
                    href="#" 
                    class="font-medium text-gray-400 hover:text-gray-300"
                    @click.prevent="handleForgotPassword"
                  >
                    Forgot your password?
                  </a>
                </div>
              </div>

              <div>
                <button
                  type="submit"
                  @click.prevent="handleLogin"
                  class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-gray-600 hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-400"
                >
                  Sign in
                </button>
              </div>
            </form>

            <div class="relative py-4" v-if="authMethods?.authProviders?.length">
              <div class="absolute inset-0 flex items-center">
                <div class="w-full border-t border-gray-300"></div>
              </div>
              <div class="relative flex justify-center text-sm rounded-md">
                <span class="px-2 text-white bg-gray-700 rounded-md">Or continue with</span>
              </div>
            </div>
          </div>

          <!-- OAuth Providers -->
          <div v-if="authMethods?.authProviders?.length" class="space-y-3">
            <button
              v-for="provider in authMethods.authProviders"
              :key="provider.name"
              @click="handleOauthLogin(provider.name)"
              class="w-full flex justify-center items-center px-4 py-2 border border-gray-600 shadow-sm text-sm font-medium rounded-md text-white bg-gray-600 hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-400"
            >
              <!-- Add provider icon here if needed -->
              Sign in with {{ provider.name }}
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <div class="mt-auto">
      <Footer />
    </div>
  </div>
</template>

<script setup>
import Footer from '@/components/Footer.vue'
import { inject, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const pocketbase = inject('$pocketbase')
if (!pocketbase) {
  console.error('PocketBase not injected')
}

const router = useRouter()
const emailError = ref('')
const email = ref('')
const password = ref('')
const rememberMe = ref(false)
const showForgotPassword = ref(false)

const authMethods = ref({
  usernamePassword: false,
  mfa: false,
  otp: false,
  authProviders: []
})

const loadAuthMethods = async () => {
  try {
    const methods = await pocketbase.collection('users').listAuthMethods()
    
    // Check if password auth is available
    authMethods.value.usernamePassword = methods.password.enabled;

    // Check if MFA is available
    authMethods.value.mfa = methods.mfa.enabled;

    // Check if OTP is available
    authMethods.value.otp = methods.otp.enabled;
    
    // Get available OAuth providers
    authMethods.value.authProviders = []
    if (methods.oauth2.enabled) {
      methods.oauth2.providers.forEach(provider => {
        authMethods.value.authProviders.push(provider)
      })
    }
  } catch (error) {
    console.error('Failed to load auth methods:', error)
  }
}

onMounted(async () => {
  await loadAuthMethods()
  
  // Load saved form data
  const savedEmail = localStorage.getItem('rememberedEmail')
  const savedRememberMe = localStorage.getItem('rememberMe') === 'true'
  
  if (savedRememberMe) {
    email.value = savedEmail || ''
    rememberMe.value = true
  }
})

const validateEmail = () => {
  if (!email.value) {
    emailError.value = 'Please enter an email address'
  } else if (!/\S+@\S+\.\S+/.test(email.value)) {
    emailError.value = 'Please enter a valid email address'
  } else {
    emailError.value = ''
  }
}

const handleLogin = async () => {
  validateEmail()
  if (emailError.value) {
    console.log("email error")
    return
  }

  try {
    // Save or clear email based on remember me checkbox
    if (rememberMe.value) {
      localStorage.setItem('rememberedEmail', email.value)
      localStorage.setItem('rememberMe', 'true')
    } else {
      localStorage.removeItem('rememberedEmail')
      localStorage.removeItem('rememberMe')
    }

    const authData = await pocketbase.collection('users').authWithPassword(
      email.value,
      password.value
    )
    
    if (authData) {
      router.push('/domains')
    }
  } catch (error) {
    password.value = '' 
    emailError.value = 'Invalid email or password'
  }
}

const handleForgotPassword = () => {
  emailError.value = ''
  email.value = ''
  password.value = ''
  showForgotPassword.value = true
}

const handleResetPassword = async () => {
  validateEmail()
  if (emailError.value) return
  const res = await pocketbase.collection('users').requestPasswordReset(email.value)

  // if res is an error, show the error
  if (res.status !== 200) {
    emailError.value = "Error sending password reset email"
    email.value = ''
  } else {
    emailError.value = 'Password reset email sent to'
    email.value = ''
  }
}

const handleBackToLogin = () => {
  emailError.value = ''
  email.value = ''
  password.value = ''
  showForgotPassword.value = false
}

const handleOauthLogin = async (provider) => {
  await pocketbase.collection('users').authWithOAuth2({ provider: provider });

  // reload the page
  window.location.reload()
}
</script>

<style>
:root {
  --tooltip-background: rgb(31, 41, 55); /* bg-gray-800 */
  --tooltip-text: white;
}

input::-webkit-validation-bubble-message {
  background: var(--tooltip-background);
  color: var(--tooltip-text);
  border: none;
}

input::-moz-validation-bubble-message {
  background: var(--tooltip-background);
  color: var(--tooltip-text);
  border: none;
}
</style> 