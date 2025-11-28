<template>
    <div class="card border-0 bg-body-tertiary">
        <div class="card-body">
            <div class="d-flex justify-content-between align-items-center mb-3">
                <div class="d-flex align-items-center gap-2">
                    <RouterLink :to="{ name: 'create-student' }" class="btn btn-primary d-flex align-items-center">
                        <i class="bi bi-plus-lg me-2"></i> Tambah Siswa
                    </RouterLink>
                    <button
                        type="button"
                        class="btn btn-outline-primary d-flex align-items-center"
                        @click="goToImport"
                    >
                        <i class="bi bi-upload me-2"></i>
                        Import Siswa
                    </button>
                </div>
                <div class="bg-secondary-subtle rounded-pill px-3 py-1 text-primary">
                    <i class="bi bi-person me-2"></i>
                    <span class="small">{{ students.length }}</span>
                </div>
            </div>

            <div v-if="students.length === 0" class="text-center text-muted py-4 small">
                Tidak ada data siswa.
            </div>

            <div
                v-else
                v-for="student in students"
                :key="student.id || student.full_name"
                class="d-flex align-items-center mb-3 shadow-sm p-3"
            >
                <img
                    :src="student.avatar || `https://i.pravatar.cc/150?u=${student.full_name}`"
                    :alt="`Avatar ${student.full_name}`"
                    class="rounded-circle avatar"
                />
                <div class="ms-3 flex-grow-1">
                    <div class="row gx-3 w-100">
                        <div class="col-12 col-md-3">
                            <span class="fw-bold text-dark">{{ student.full_name }}</span>
                            <small class="d-block text-muted">RFID: {{ student.rfid || '-' }}</small>
                            <small class="d-block text-muted">NISN: {{ student.nisn || '-' }}</small>
                        </div>
                        <div class="col-12 col-md-3">
                            <small class="d-block text-muted">Kelas</small>
                            <span class="fw-bold text-dark">{{ student.class?.name || '-' }}</span>
                        </div>
                        <div class="col-12 col-md-3">
                            <small class="d-block text-muted">Tempat / Tgl Lahir</small>
                            <span class="fw-bold text-dark">
                                {{ student.birth_place || '-' }}<template v-if="student.birth_place && student.born">, </template>{{ formatDateID(student.born) }}
                            </span>
                            <small class="d-block text-muted">No. Whatsapp</small>
                            <span class="fw-bold text-dark">{{ student.whatsapp || '-' }}</span>
                        </div>
                        <div class="col-12 col-md-3">
                            <small class="d-block text-muted">Orang Tua</small>
                            <span class="fw-bold text-dark">{{ student.parent?.full_name || '-' }}</span>
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
const loading = ref(false)

const router = useRouter()

const getAuthHeaders = () => {
    const token = localStorage.getItem('token')
    return token ? { Authorization: `Bearer ${token}` } : {}
}

const fetchStudents = async () => {
    loading.value = true
    try {
        const res = await axios.get(`${apiBaseUrl}/students`, {
            headers: getAuthHeaders(),
        })
        students.value = res.data?.data || []
    } catch (error) {
        console.error('Failed to fetch students', error)
    } finally {
        loading.value = false
    }
}

const deleteStudent = async (studentId) => {
    try {
        await axios.delete(`${apiBaseUrl}/students/${studentId}`, {
            headers: getAuthHeaders(),
        })
        students.value = students.value.filter(s => s.id !== studentId)
    } catch (error) {
        console.error('Failed to delete student', error)
    }
}

const edit = (student) => {
    if (router && student?.id) {
        router.push(`/admin/students/${student.id}/edit`)
    }
}

const askDelete = (student) => {
    if (!student?.id) return
    if (confirm(`Hapus data siswa "${student.full_name}"?`)) {
        deleteStudent(student.id)
    }
}

const goToImport = () => {
    // Sesuaikan dengan nama rute halaman import siswa di konfigurasi router-mu
    router.push({ name: 'import-students' })
}

onMounted(() => {
    fetchStudents()
})

const formatDateID = (value) => {
    if (!value) return '-'
    const date = new Date(value)
    if (Number.isNaN(date)) return '-'
    return new Intl.DateTimeFormat('id-ID', {
        day: '2-digit',
        month: 'long',
        year: 'numeric',
    }).format(date)
}
</script>

<style scoped>
.avatar {
    width: 40px;
    height: 40px;
    object-fit: cover;
}
</style>
