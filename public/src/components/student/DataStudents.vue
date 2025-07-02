<template>
    <div class="card border-0 bg-body-tertiary" style="height: calc(100vh - 120px);">
        <div class="card-body">
            <div class="d-flex justify-content-between align-items-center mb-3">
               <RouterLink :to="{ name: 'create-student' }" class="btn btn-primary d-flex align-items-center">
                    <i class="bi bi-plus-lg me-2"></i> Tambah Siswa
                </RouterLink>
                <div class="bg-secondary-subtle rounded-pill px-3 py-1 text-primary">
                    <i class="bi bi-person me-2"></i>
                    <span class="small">{{ students.length }}</span>
                </div>
            </div>
            <div v-for="student in students" :key="student.id || student.full_name"
                class="d-flex align-items-center mb-3 shadow-sm p-3">
                <img :src="student.avatar || `https://i.pravatar.cc/150?u=${student.full_name}`"
                    :alt="`Avatar ${student.full_name}`" class="rounded-circle avatar" />
                <div class="ms-3 flex-grow-1">
                    <div class="row gx-3 w-100">
                        <div class="col-12 col-md-3">
                            <span class="fw-bold text-dark">{{ student.full_name }}</span>
                            <small class="d-block text-muted">{{ student.rfid }}</small>
                        </div>
                        <div class="col-12 col-md-3">
                            <small class="d-block text-muted">Kelas</small>
                            <span class="fw-bold text-dark">{{ student.class.name }}</span>
                        </div>
                        <div class="col-12 col-md-3">
                            <small class="d-block text-muted">Orang Tua</small>
                            <span class="fw-bold text-dark">{{ student.parent?.full_name || '-' }}</span>
                        </div>
                        <div class="col-12 col-md-3">
                            <small class="d-block text-muted">No. Whatsapp</small>
                            <span class="fw-bold text-dark">{{ student.parent?.whatsapp || '-' }}</span>
                        </div>
                    </div>
                </div>
                <div class="d-flex align-items-center">
                    <button @click="edit(student)" class="btn p-1 text-warning" title="Edit">
                        <i class="bi bi-pencil-square"></i>
                    </button>
                    <button @click="askDelete(student)" class="btn p-1 text-danger" title="Delete">
                        <i class="bi bi-trash"></i>
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { apiBaseUrl } from '@/config'

const students = ref([])

const router = useRouter()

function edit(student) {
    if (router && student?.id) {
        router.push(`/admin/students/${student.id}/edit`)
    }
}

function askDelete(student) {
    console.log('Ask delete for', student)
}

onMounted(() => {
    const token = localStorage.getItem('token')
    axios
        .get(`${apiBaseUrl}/students`, {
            headers: { Authorization: `Bearer ${token}` },
        })
        .then(res => { students.value = res.data.data })
        .catch(err => { console.error('Failed to fetch students', err) })
})
</script>
<style scoped>
.avatar {
    width: 40px;
    height: 40px;
    object-fit: cover;
}
</style>