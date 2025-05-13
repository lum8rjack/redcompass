import pocketbase from '../pocketbase'

export const isAuthenticated = () => {
  return pocketbase.authStore.isValid
} 

export const canEdit = () => {
  return pocketbase.authStore.isValid && pocketbase.authStore.record.role !== 'viewer'
}