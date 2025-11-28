<template>
    <div class="card border-0 bg-body-tertiary mb-3 navbar-sticky">
        <div class="card-body">
            <span class="fs-4 fw-bold mb-0 card-title">{{ currentTitle }}</span>
            <div class="dropdown float-end">
                <a class="d-flex align-items-center text-decoration-none" href="#" role="button" id="dropdownMenuLink"
                    data-bs-toggle="dropdown" aria-expanded="false">
                    <img src="https://i.pravatar.cc/150?u=hendri" alt="profile" class="rounded-circle me-2" width="40"
                        height="40">
                </a>
                <ul class="dropdown-menu dropdown-menu-end" aria-labelledby="dropdownMenuLink">
                    <li><a class="dropdown-item" href="#"><i class="bi bi-person-circle me-2"></i>Profile</a></li>
                    <li>
                        <hr class="dropdown-divider">
                    </li>
                    <li><a class="dropdown-item" href="#" @click.prevent="logout"><i
                                class="bi bi-box-arrow-right me-2"></i>Logout</a></li>
                </ul>
            </div>
        </div>
    </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { apiBaseUrl } from '@/config'

const route = useRoute()
const currentTitle = computed(() => route.meta.title || 'Dashboard')

function logout() {
    axios.post(`${apiBaseUrl}/logout`, {}, {
        headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`
        }
    }).then(() => {
        localStorage.removeItem('token')
        window.location.href = '/login'
    }).catch(error => {
        console.error('Logout failed:', error)
    })
}
</script>

<style scoped>
/* .navbar-sticky {
    position: sticky;
    top: 16px;
    z-index: 1020;
} */
</style>