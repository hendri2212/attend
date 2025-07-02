<template>
    <div class="card border-0 bg-body-tertiary" style="height: calc(100vh - 120px);">
        <div class="card-body">
            <div class="d-flex justify-content-between align-items-center mb-3">
                <RouterLink :to="{ name: 'create-parent' }" class="btn btn-primary d-flex align-items-center">
                    <i class="bi bi-plus-lg me-2"></i> Tambah Orang Tua
                </RouterLink>
                <div class="bg-secondary-subtle rounded-pill px-3 py-1 text-primary">
                    <i class="bi bi-person me-2"></i>
                    <span class="small">{{ parents.length }}</span>
                </div>
            </div>
            <div v-for="parent in parents" :key="parent.id || parent.full_name"
                class="d-flex align-items-center mb-3 shadow-sm p-3">
                <img :src="`https://i.pravatar.cc/150?u=${parent.full_name}`" :alt="`Avatar ${parent.full_name}`"
                    class="rounded-circle avatar" />
                <div class="ms-3 flex-grow-1">
                    <div class="row gx-3 w-100">
                        <div class="col-12 col-md-6">
                            <span class="fw-bold text-dark">{{ parent.full_name }}</span>
                        </div>
                        <div class="col-12 col-md-6">
                            <small class="d-block text-muted">No. Whatsapp</small>
                            <span class="fw-bold text-dark">{{ parent.whatsapp || '-' }}</span>
                        </div>
                    </div>
                </div>
                <div class="d-flex align-items-center">
                    <button @click="edit(parent)" class="btn p-1 text-warning" title="Edit">
                        <i class="bi bi-pencil-square"></i>
                    </button>
                    <button @click="askDelete(parent)" class="btn p-1 text-danger" title="Delete">
                        <i class="bi bi-trash"></i>
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { apiBaseUrl } from '@/config'

const parents = ref([])

onMounted(() => {
    const token = localStorage.getItem('token')
    axios.get(apiBaseUrl + '/parents', {
        headers: {
            Authorization: `Bearer ${token}`
        }
    })
        .then(response => {
            parents.value = response.data
        })
        .catch(error => {
            console.error('Gagal mengambil data orang tua:', error)
        })
})

function edit(parent) {
    // implementasi edit
}

function askDelete(parent) {
    // implementasi delete
}
</script>
<style scoped>
.avatar {
    width: 40px;
    height: 40px;
    object-fit: cover;
}
</style>