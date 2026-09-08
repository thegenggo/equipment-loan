import { createRouter, createWebHistory } from 'vue-router'

import { useAuthStore } from '@/stores/auth'

declare module 'vue-router' {
  interface RouteMeta {
    isPublic?: boolean
    adminOnly?: boolean
  }
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { isPublic: true }
    },
    {
      path: '/',
      name: 'equipments',
      component: () => import('@/views/staff/EquipmentListView.vue'),
    },
    {
      path: '/my-loans',
      name: 'my-loans',
      component: () => import('@/views/staff/MyLoansView.vue'),
    },
    {
      path: '/admin/loans',
      name: 'admin-loans',
      component: () => import('@/views/admin/LoanApprovalView.vue'),
      meta: { adminOnly: true }
    },
    {
      path: '/admin/equipments',
      name: 'admin-equipments',
      component: () => import('@/views/admin/EquipmentManageView.vue'),
      meta: { adminOnly: true }
    },
    { path: '/:pathMatch(.*)*', redirect: '/' }
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()

  if (to.meta.isPublic) {
    return auth.isAuthenticated ? { name: 'equipments' } : true
  }

  if (!auth.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  if (to.meta.adminOnly && !auth.isAdmin) {
    return { name: 'equipments' }
  }

  return true
})

export default router
