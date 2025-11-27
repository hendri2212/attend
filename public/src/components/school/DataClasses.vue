<template>
    <div class="card border-0 bg-body-tertiary" style="height: calc(100vh - 120px);">
        <div class="card-body">
            <div class="d-flex justify-content-between align-items-center mb-3">
                <RouterLink :to="{ name: 'create-class' }" class="btn btn-primary d-flex align-items-center">
                    <i class="bi bi-plus-lg me-2"></i> Tambah Kelas
                </RouterLink>
                <div class="bg-secondary-subtle rounded-pill px-3 py-1 text-primary">
                    <i class="bi bi-person me-2"></i>
                    <span class="small">{{ classes.length }}</span>
                </div>
            </div>
            <div v-for="classItem in classes" :key="classItem.id || classItem.name"
                class="d-flex align-items-center mb-3 shadow-sm p-3">
                <div class="ms-3 flex-grow-1">
                    <div class="row gx-3 w-100">
                        <div class="col-12 col-md-6">
                            <span class="fw-bold text-dark">{{ classItem.name }}</span>
                        </div>
                    </div>
                </div>
                <div class="d-flex align-items-center">
                    <button @click="edit(classItem)" class="btn p-1 text-warning" title="Edit">
                        <i class="bi bi-pencil-square"></i>
                    </button>
                    <button @click="askDelete(classItem)" class="btn p-1 text-danger" title="Delete">
                        <i class="bi bi-trash"></i>
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>
<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { apiBaseUrl } from '@/config'

const classes = ref([])

const router = useRouter()

onMounted(() => {
    const token = localStorage.getItem('token')
    axios.get(apiBaseUrl + '/classes', {
        headers: {
            Authorization: `Bearer ${token}`
        }
    })
    .then(response => {
        classes.value = response.data
    })
    .catch(error => {
        console.error('Gagal mengambil data kelas:', error)
    })
})

function edit(classItem) {
    if (router && classItem?.id) {
        router.push(`/admin/classes/${classItem.id}/edit`)
    }
}

function askDelete(classItem) {
    if (confirm(`Apakah Anda yakin ingin menghapus kelas "${classItem.name}"?`)) {
        const token = localStorage.getItem('token')
        axios.delete(`${apiBaseUrl}/classes/${classItem.id}`, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        })
        .then(() => {
            alert('Kelas berhasil dihapus.')
            // Refresh the class list
            classes.value = classes.value.filter(c => c.id !== classItem.id)
        })
        .catch(error => {
            console.error('Gagal menghapus kelas:', error)
            alert('Gagal menghapus kelas.')
        })
    }
}
</script>
