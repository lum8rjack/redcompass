import pocketbase from '../pocketbase'

export const isAuthenticated = () => {
  return pocketbase.authStore.isValid
} 