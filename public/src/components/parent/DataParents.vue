<template>
    <div class="card border-0 bg-body-tertiary">
        <div class="card-body">
            <div class="d-flex flex-wrap align-items-center justify-content-between mb-3 gap-2">
                <div>
                    <div class="d-flex align-items-center gap-2 mb-1">
                        <h5 class="mb-0 fw-bold">Daftar Orang Tua</h5>
                        <span class="badge rounded-pill bg-secondary-subtle text-primary d-flex align-items-center px-3 py-1">
                            <i class="bi bi-people me-2"></i>
                            <span class="small">{{ parents.length }} Orang Tua</span>
                        </span>
                    </div>
                    <p class="text-muted small mb-0">
                        Kelola daftar orang tua, tambah, ubah, dan hapus sesuai kebutuhan.
                    </p>
                </div>
                <RouterLink :to="{ name: 'create-parent' }" class="btn btn-outline-primary d-flex align-items-center">
                    <i class="bi bi-plus-lg me-2"></i>Orang Tua
                </RouterLink>
            </div>
            <hr>
            <div v-for="parent in parents" :key="parent.id || parent.full_name"
                class="d-flex align-items-center mb-3 shadow-sm p-3">
                <img :src="`https://i.pravatar.cc/150?u=${parent.full_name}`" :alt="`Avatar ${parent.full_name}`"
                    class="rounded-circle avatar" />
                <div class="ms-3 flex-grow-1">
                    <div class="row gx-3 w-100">
                        <div class="col-12 col-md-4">
                            <span class="fw-bold text-dark">{{ parent.full_name }}</span>
                        </div>
                        <div class="col-12 col-md-4">
                            <small class="d-block text-muted">No. Whatsapp</small>
                            <span class="fw-bold text-dark">{{ parent.whatsapp || '-' }}</span>
                        </div>
                        <div class="col-12 col-md-4">
                            <small class="d-block text-muted">Alamat</small>
                            <span class="fw-bold text-dark">{{ parent.address || '-' }}</span>
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
import { useRouter } from 'vue-router'
import { apiBaseUrl } from '@/config'

const parents = ref([])

const router = useRouter()

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
    if (router && parent?.id) {
        router.push({ name: 'edit-parent', params: { id: parent.id } })
    }
}

function askDelete(parent) {
    if (confirm(`Apakah Anda yakin ingin menghapus orang tua "${parent.full_name}"?`)) {
        const token = localStorage.getItem('token')
        axios.delete(`${apiBaseUrl}/parents/${parent.id}`, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        })
        .then(() => {
            parents.value = parents.value.filter(p => p.id !== parent.id)
            alert('Orang tua berhasil dihapus!')
        })
        .catch(error => {
            console.error('Gagal menghapus orang tua:', error)
            alert('Gagal menghapus orang tua')
        })
    }
}
</script>
<style scoped>
.avatar {
    width: 40px;
    height: 40px;
    object-fit: cover;
}
</style>