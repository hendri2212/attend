<template>
    <div class="container-fluid py-3 bg-secondary-subtle d-flex">
        <div class="card border-0 bg-body-tertiary sidebar" style="width: 280px;">
            <div class="card-body">
                <a href="/" class="d-flex align-items-center mb-3 mb-md-0 me-md-auto text-decoration-none fw-bold">
                    <svg xmlns="http://www.w3.org/2000/svg" width="40" height="32" fill="currentColor"
                        class="bi bi-battery-charging me-2" viewBox="0 0 16 16">
                        <path
                            d="M9.585 2.568a.5.5 0 0 1 .226.58L8.677 6.832h1.99a.5.5 0 0 1 .364.843l-5.334 5.667a.5.5 0 0 1-.842-.49L5.99 9.167H4a.5.5 0 0 1-.364-.843l5.333-5.667a.5.5 0 0 1 .616-.09z" />
                        <path
                            d="M2 4h4.332l-.94 1H2a1 1 0 0 0-1 1v4a1 1 0 0 0 1 1h2.38l-.308 1H2a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2" />
                        <path
                            d="M2 6h2.45L2.908 7.639A1.5 1.5 0 0 0 3.313 10H2zm8.595-2-.308 1H12a1 1 0 0 1 1 1v4a1 1 0 0 1-1 1H9.276l-.942 1H12a2 2 0 0 0 2-2V6a2 2 0 0 0-2-2z" />
                        <path
                            d="M12 10h-1.783l1.542-1.639q.146-.156.241-.34zm0-3.354V6h-.646a1.5 1.5 0 0 1 .646.646M16 8a1.5 1.5 0 0 1-1.5 1.5v-3A1.5 1.5 0 0 1 16 8" />
                    </svg>
                    <span class="fs-4 text-nowrap">Absensi Baru</span> </a>
                <hr>
                <ul class="nav nav-pills flex-column mb-auto">
                    <li>
                        <RouterLink :to="{ name: 'dashboard' }"
                            :class="['nav-link', route.name === 'dashboard' ? 'active' : 'link-body-emphasis']">
                            <i class="bi bi-speedometer2 me-2" width="16" height="16"></i>
                            Dashboard
                        </RouterLink>
                    </li>
                    <!-- <li>
                      <RouterLink
                        :to="{ name: 'homeroom' }"
                        :class="['nav-link', route.name === 'homeroom' ? 'active' : 'link-body-emphasis']">
                        <i class="bi bi-pencil me-2" width="16" height="16"></i>
                        Wali Kelas
                      </RouterLink>
                    </li> -->
                    <li>
                        <RouterLink :to="{ name: 'classes' }"
                            :class="['nav-link', ['classes', 'create-class', 'edit-class'].includes(route.name) ? 'active' : 'link-body-emphasis']">
                            <i class="bi bi-door-open me-2" width="16" height="16"></i>
                            Kelas
                        </RouterLink>
                    </li>
                    <li>
                        <RouterLink :to="{ name: 'students' }"
                            :class="['nav-link', ['students', 'create-student', 'edit-student'].includes(route.name) ? 'active' : 'link-body-emphasis']">
                            <i class="bi bi-person-badge me-2" width="16" height="16"></i>
                            Siswa
                        </RouterLink>
                    </li>
                    <li>
                        <RouterLink :to="{ name: 'parents' }"
                            :class="['nav-link', ['parents', 'create-parent', 'edit-parent'].includes(route.name) ? 'active' : 'link-body-emphasis']">
                            <i class="bi bi-people-fill me-2" width="16" height="16"></i>
                            Orang Tua
                        </RouterLink>
                    </li>
                </ul>
            </div>
        </div>
        <div class="ms-3 w-100">
            <Navbar />
            <RouterView />
        </div>
    </div>
</template>
<script setup>
import Navbar from '../components/Navbar.vue'
import { RouterLink, useRoute } from 'vue-router'
const route = useRoute()
</script>
<style scoped>
.sidebar {
    position: sticky;
    top: 16px;
    align-self: flex-start;
    max-height: calc(100vh - 32px);
}
</style>
