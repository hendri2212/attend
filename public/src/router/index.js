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
            path: '/register',
            name: 'register',
            component: () => import('../views/Login.vue'),
            meta: {
                title: 'Register',
            },
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
                    path: 'classes',
                    name: 'classes',
                    component: () => import('../components/school/DataClasses.vue'),
                    meta: {
                        title: 'Data Kelas',
                    },
                },
                {
                    path: 'classes/create',
                    name: 'create-class',
                    component: () => import('../components/school/AddClass.vue'),
                    meta: {
                        title: 'Data Kelas',
                    }
                },
                {
                    path: 'classes/:id/edit',
                    name: 'edit-class',
                    component: () => import('../components/school/EditClass.vue'),
                    meta: {
                        title: 'Data Kelas',
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
                        title: 'Data Siswa',
                    }
                },
                {
                    path: 'students/:id/edit',
                    name: 'edit-student',
                    component: () => import('../components/student/EditStudent.vue'),
                    meta: {
                        title: 'Data Siswa',
                    }
                },
                {
                    path: 'import-students',
                    name: 'import-students',
                    component: () => import('../components/student/ImportStudents.vue'),
                    meta: {
                        title: 'Data Siswa',
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
                        title: 'Data Orang Tua Siswa',
                    }
                },
                {
                    path: 'parents/:id/edit',
                    name: 'edit-parent',
                    component: () => import('../components/parent/EditParent.vue'),
                    meta: {
                        title: 'Data Orang Tua Siswa',
                    }
                },
                {
                    path: 'attendance',
                    name: 'attendance',
                    component: () => import('../components/attendance/DataAttendance.vue'),
                    meta: {
                        title: 'Data Kehadiran Siswa',
                    },
                },
                {
                    path: 'report',
                    name: 'report',
                    component: () => import('../components/report/DataReport.vue'),
                    meta: {
                        title: 'Laporan',
                    },
                },
                {
                    path: 'graphical-report',
                    name: 'graphical-report',
                    component: () => import('../components/report/GraphicalReport.vue'),
                    meta: {
                        title: 'Laporan Grafis',
                    },
                },
                {
                    path: 'teacher',
                    name: 'teacher',
                    component: () => import('../components/teacher/DataTeacher.vue'),
                    meta: {
                        title: 'Data Guru',
                    },
                },
                {
                    path: 'users',
                    name: 'users',
                    component: () => import('../components/user/DataUsers.vue'),
                    meta: {
                        title: 'Manajemen User',
                        restrictTo: ['superadmin']
                    },
                }
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

    // Check for role restriction
    if (to.meta.restrictTo) {
        const role = localStorage.getItem('role')
        if (!to.meta.restrictTo.includes(role)) {
            // Redirect to admin dashboard if not authorized
            // Or maybe show an alert? For now just redirect
            return { name: 'dashboard' }
        }
    }
})

router.afterEach((to) => {
    document.title = to.meta.title || 'Attend'
})
// ---------------------------------------------

export default router
