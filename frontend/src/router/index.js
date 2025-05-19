import { createRouter, createWebHistory } from 'vue-router'
import DomainsView from '../views/DomainsView.vue'
import LoginView from '../views/LoginView.vue'
import ProfileView from '@/views/ProfileView.vue'
import ProjectsView from '@/views/ProjectsView.vue'
import SettingsView from '@/views/SettingsView.vue'
import ProjectView from '@/views/ProjectView.vue'
import DomainView from '@/views/DomainView.vue'
import CategorizationsView from '../views/CategorizationsView.vue'
import DomainIdeasView from '../views/DomainIdeasView.vue'
import PhishingCampaignView from '@/views/PhishingCampaignView.vue'
import PhishingManagementView from '../views/PhishingManagementView.vue'
import PhishletsView from '../views/PhishletsView.vue'
import PhishingMetricsView from '../views/PhishingMetricsView.vue'
import { isAuthenticated } from '@/utils/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior: (to, from, savedPosition) => {
    return { top: 0 }
  },
  strict: true,
  routes: [
    {
      path: '/',
      redirect: to => {
        return isAuthenticated() ? '/domains' : '/login'
      }
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { requiresAuth: false, description: 'Login' },
      beforeEnter: (to, from, next) => {
        if (isAuthenticated()) {
          next('/domains')
        } else {
          next()
        }
      }
    },
    {
      path: '/categorizations',
      name: 'categorizations',
      component: CategorizationsView,
      meta: { requiresAuth: true, description: 'Categorization Management' }
    },
    {
      path: '/domains',
      name: 'domains',
      component: DomainsView,
      meta: { requiresAuth: true, description: 'Domain Management' }
    },
    {
      path: '/domain/:id',
      name: 'domain',
      component: DomainView,
      meta: { requiresAuth: true, description: 'Domain Details' }
    },
    {
      path: '/domain-ideas',
      name: 'domain-ideas',
      component: DomainIdeasView,
      meta: { requiresAuth: true, description: 'Domain Ideas' }
    },
    {
      path: '/phishing',
      name: 'phishing',
      component: PhishingManagementView,
      meta: { requiresAuth: true, description: 'Phishing Campaign Management' }
    },
    {
      path: '/phishing/:id',
      name: 'phishing-campaign',
      component: PhishingCampaignView,
      meta: { requiresAuth: true, description: 'Phishing Campaign Details' }
    },
    {
      path: '/phishing-metrics',
      name: 'phishing-metrics',
      component: PhishingMetricsView,
      meta: { requiresAuth: true, description: 'Phishing Metrics' }
    },
    {
      path: '/phishlets',
      name: 'phishlets',
      component: PhishletsView,
      meta: { requiresAuth: true, description: 'Phishlet Management' }
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView,
      meta: { requiresAuth: true, description: 'Profile' }
    },
    {
      path: '/projects',
      name: 'projects',
      component: ProjectsView,
      meta: { requiresAuth: true, description: 'Projects' }
    },
    {
      path: '/project/:id',
      name: 'project',
      component: ProjectView,
      meta: { requiresAuth: true, description: 'Project Details' }
    },
    {
      path: '/settings',
      name: 'settings',
      component: SettingsView,
      meta: { requiresAuth: true, description: 'Settings' }
    }
  ]
})

// Global navigation guard
router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!isAuthenticated()) {
      next('/login')
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router
