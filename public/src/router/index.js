// Setup Axios Interceptor (gunakan CDN axios)
if (window.axios) {
    window.axios.interceptors.response.use(
        response => response,
        error => {
            if (error.response && error.response.status === 401) {
                localStorage.removeItem('token')
                window.location.href = '/login'
            }
            return Promise.reject(error)
        }
    )
}
import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: HomeView,
        },
        {
            path: '/login',
            name: 'login',
            component: () => import('../views/Login.vue'),
            meta: {
                title: 'Login',
            },
        },
        {
            path: '/admin',
            name: 'admin',
            component: () => import('../views/DashboardView.vue'),
            meta: {
                requiresAuth: true,
                title: 'Dashboard Admin',
            },
            children: [
                {
                    path: '',
                    name: 'dashboard',
                    component: () => import('../components/Dashboard.vue'),
                    meta: {
                        title: 'Dashboard',
                    }
                },
                {
                    path: 'students',
                    name: 'students',
                    component: () => import('../components/student/DataStudents.vue'),
                    meta: {
                        title: 'Data Siswa',
                    },
                },
                {
                    path: 'students/create',
                    name: 'create-student',
                    component: () => import('../components/student/AddStudent.vue'),
                    meta: {
                        title: 'Tambah Siswa',
                    }
                },
                {
                    path: 'students/:id/edit',
                    name: 'edit-student',
                    component: () => import('../components/student/EditStudent.vue'),
                    meta: {
                        title: 'Edit Siswa',
                    }
                },
                {
                    path: 'parents',
                    name: 'parents',
                    component: () => import('../components/parent/DataParents.vue'),
                    meta: {
                        title: 'Data Orang Tua Siswa',
                    },
                },
                {
                    path: 'parents/create',
                    name: 'create-parent',
                    component: () => import('../components/parent/AddParent.vue'),
                    meta: {
                        title: 'Tambah Orang Tua',
                    }
                },
            ],
        },
    ],
})


// ---------------------------------------------
// Global guards: auth check & dynamic <title>
// ---------------------------------------------
router.beforeEach((to) => {
    // If route needs auth but no token, redirect to /login
    if (to.meta.requiresAuth && !localStorage.getItem('token')) {
        return { name: 'login', query: { redirect: to.fullPath } }
    }
})

router.afterEach((to) => {
    document.title = to.meta.title || 'Attend'
})
// ---------------------------------------------

export default router
