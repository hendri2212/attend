<template>
    <div class="card border-0 bg-body-tertiary">
        <div class="card-body">
            <div class="d-flex flex-wrap align-items-center justify-content-between mb-3 gap-2">
                <div>
                    <div class="d-flex align-items-center gap-2 mb-1">
                        <h5 class="mb-0 fw-bold">Daftar Siswa</h5>
                        <span class="badge rounded-pill bg-secondary-subtle text-primary d-flex align-items-center px-3 py-1">
                            <i class="bi bi-person me-2"></i>
                            <span class="small">{{ students.length }} Siswa</span>
                        </span>
                    </div>
                    <p class="text-muted small mb-0">
                        Kelola daftar siswa, tambah, ubah, dan hapus sesuai kebutuhan.
                    </p>
                </div>
                <div class="d-flex gap-2">
                    <RouterLink :to="{ name: 'create-student' }" class="btn btn-outline-primary d-flex align-items-center">
                        <i class="bi bi-plus-lg me-2"></i>Siswa
                    </RouterLink>
                    <button
                        type="button"
                        class="btn btn-outline-primary d-flex align-items-center"
                        @click="openImportModal"
                    >
                        <i class="bi bi-upload me-2"></i>
                        Import Siswa
                    </button>
                </div>
            </div>
            <hr>
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

        <!-- Import Modal -->
        <div v-if="showImportModal" class="modal fade show d-block" style="background: rgba(0,0,0,0.5);" tabindex="-1">
            <div class="modal-dialog">
                <div class="modal-content">
                    <div class="modal-header">
                        <h5 class="modal-title">Import Data Siswa</h5>
                        <button type="button" class="btn-close" @click="showImportModal = false"></button>
                    </div>
                    <div class="modal-body">
                        <div class="alert alert-info small">
                            <i class="bi bi-info-circle me-1"></i>
                            Gunakan format Excel (.xlsx) yang sesuai.
                            <a href="/format_import_students.xlsx" download class="fw-bold">Download Template</a>
                        </div>
                        
                        <div class="mb-3">
                            <label class="form-label">Pilih File Excel</label>
                            <input type="file" class="form-control" accept=".xlsx" @change="handleFileChange">
                        </div>

                        <div v-if="importLoading" class="text-center py-2">
                            <div class="spinner-border text-primary" role="status"></div>
                            <p class="small text-muted mt-2">Sedang memproses...</p>
                        </div>

                        <div v-if="importResult.message" class="alert alert-success mt-3 small">
                            <i class="bi bi-check-circle me-1"></i>
                            {{ importResult.message }} ({{ importResult.successCount }} berhasil)
                        </div>

                        <div v-if="importResult.errors.length > 0" class="alert alert-danger mt-3 small" style="max-height: 200px; overflow-y: auto;">
                            <strong>Gagal:</strong>
                            <ul class="mb-0 ps-3">
                                <li v-for="(err, idx) in importResult.errors" :key="idx">{{ err }}</li>
                            </ul>
                        </div>
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-secondary" @click="showImportModal = false">Tutup</button>
                        <button type="button" class="btn btn-primary" :disabled="importLoading" @click="uploadFile">Upload</button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { apiBaseUrl } from '@/config'

const students = ref([])
const loading = ref(false)

// Import Logic
const showImportModal = ref(false)
const importFile = ref(null)
const importLoading = ref(false)
const importResult = reactive({
    message: '',
    successCount: 0,
    errors: []
})

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

// Open Import Modal
const openImportModal = () => {
    showImportModal.value = true
    importFile.value = null
    importResult.message = ''
    importResult.successCount = 0
    importResult.errors = []
}

const handleFileChange = (event) => {
    importFile.value = event.target.files[0]
}

const uploadFile = async () => {
    if (!importFile.value) {
        alert('Pilih file terlebih dahulu!')
        return
    }

    importLoading.value = true
    importResult.message = ''
    importResult.errors = []

    const formData = new FormData()
    formData.append('file', importFile.value)

    try {
        const { data } = await axios.post(`${apiBaseUrl}/students/import`, formData, {
            headers: {
                ...getAuthHeaders(),
                'Content-Type': 'multipart/form-data'
            }
        })
        
        importResult.message = data.message
        importResult.successCount = data.success_count
        importResult.errors = data.errors || []
        
        if (data.success_count > 0) {
            fetchStudents() // Refresh list
        }
    } catch (error) {
        console.error(error)
        importResult.errors = [error.response?.data?.error || 'Gagal mengupload file']
    } finally {
        importLoading.value = false
    }
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
