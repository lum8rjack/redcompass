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
import PhishingView from '../views/PhishingView.vue'
import PhishletsView from '../views/PhishletsView.vue'
import { isAuthenticated } from '@/utils/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
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
      beforeEnter: (to, from, next) => {
        if (isAuthenticated()) {
          next('/domains')
        } else {
          next()
        }
      }
    },
    {
      path: '/domains',
      name: 'domains',
      component: DomainsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/domain-ideas',
      name: 'domain-ideas',
      component: DomainIdeasView,
      meta: { requiresAuth: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView,
      meta: { requiresAuth: true }
    },
    {
      path: '/projects',
      name: 'projects',
      component: ProjectsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/settings',
      name: 'settings',
      component: SettingsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/project/:id',
      name: 'project',
      component: ProjectView,
      meta: { requiresAuth: true }
    },
    {
      path: '/domain/:id',
      name: 'domain',
      component: DomainView,
      meta: { requiresAuth: true }
    },
    {
      path: '/categorizations',
      name: 'categorizations',
      component: CategorizationsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/phishing',
      name: 'phishing',
      component: PhishingView,
      meta: { requiresAuth: true, description: 'Phishing Campaign Management' }
    },
    {
      path: '/phishlets',
      name: 'phishlets',
      component: PhishletsView,
      meta: { requiresAuth: true, description: 'Phishlet Management' }
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
