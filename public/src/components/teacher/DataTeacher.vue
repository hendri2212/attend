<template>
    <div class="card border-0 bg-body-tertiary">
        <div class="card-body">
            <div class="d-flex flex-wrap align-items-center justify-content-between mb-3 gap-2">
                <div>
                    <div class="d-flex align-items-center gap-2 mb-1">
                        <h5 class="mb-0 fw-bold">Daftar Guru</h5>
                        <span class="badge rounded-pill bg-secondary-subtle text-primary d-flex align-items-center px-3 py-1">
                            <i class="bi bi-person-badge me-2"></i>
                            <span class="small">{{ teachers.length }} Guru</span>
                        </span>
                    </div>
                    <p class="text-muted small mb-0">
                        Kelola daftar guru, tambah, ubah, dan hapus sesuai kebutuhan.
                    </p>
                </div>
                <div class="d-flex gap-2 align-items-center">
                    <button type="button" class="btn btn-outline-primary d-flex align-items-center"
                        @click="openAddModal">
                        <i class="bi bi-plus-lg me-2"></i>
                        Tambah Guru
                    </button>
                </div>
            </div>
            <hr>

            <!-- Loading State -->
            <div v-if="loading" class="text-center py-5">
                <div class="spinner-border text-primary" role="status"></div>
                <p class="text-muted small mt-2">Memuat data guru...</p>
            </div>

            <!-- Empty State -->
            <div v-else-if="teachers.length === 0" class="text-center text-muted py-4 small">
                <i class="bi bi-person-x display-6 d-block mb-2"></i>
                Tidak ada data guru.
            </div>

            <!-- Teacher List -->
            <div v-else v-for="teacher in teachers" :key="teacher.user_id"
                class="d-flex align-items-center mb-3 shadow-sm p-3 rounded bg-white">
                <img :src="teacher.photo ? `${apiBaseUrl.replace('/api', '')}${teacher.photo}` : `https://i.pravatar.cc/150?u=${teacher.full_name}`"
                    :alt="`Avatar ${teacher.full_name}`" class="rounded-circle avatar" />
                <div class="ms-3 flex-grow-1">
                    <div class="row gx-3 w-100">
                        <div class="col-12 col-md-4">
                            <span class="fw-bold text-dark">{{ teacher.full_name }}</span>
                            <small class="d-block text-muted">NIP: {{ teacher.nip || '-' }}</small>
                        </div>
                        <div class="col-12 col-md-4">
                            <small class="d-block text-muted">Email</small>
                            <span class="fw-bold text-dark">{{ teacher.user?.email || '-' }}</span>
                        </div>
                        <div class="col-12 col-md-4">
                            <small class="d-block text-muted">Status</small>
                            <span class="badge bg-success-subtle text-success-emphasis">Aktif</span>
                        </div>
                    </div>
                </div>
                <div class="d-flex align-items-center">
                    <button @click="editTeacher(teacher)" class="btn p-1 text-warning" title="Edit">
                        <i class="bi bi-pencil-square"></i>
                    </button>
                    <button @click="askDelete(teacher)" class="btn p-1 text-danger" title="Delete">
                        <i class="bi bi-trash"></i>
                    </button>
                </div>
            </div>
        </div>

        <!-- Add/Edit Modal -->
        <div v-if="showModal" class="modal fade show d-block" style="background: rgba(0,0,0,0.5);" tabindex="-1">
            <div class="modal-dialog">
                <div class="modal-content">
                    <div class="modal-header">
                        <h5 class="modal-title">{{ isEditMode ? 'Edit Guru' : 'Tambah Guru' }}</h5>
                        <button type="button" class="btn-close" @click="closeModal"></button>
                    </div>
                    <div class="modal-body">
                        <div class="mb-3">
                            <label class="form-label">Sekolah <span class="text-danger">*</span></label>
                            <select class="form-select" v-model="form.school_id">
                                <option value="">Pilih Sekolah</option>
                                <option v-for="school in schools" :key="school.id" :value="school.id">
                                    {{ school.name }}
                                </option>
                            </select>
                        </div>
                        <div class="mb-3">
                            <label class="form-label">Nama Lengkap <span class="text-danger">*</span></label>
                            <input type="text" class="form-control" v-model="form.full_name"
                                placeholder="Masukkan nama lengkap">
                        </div>
                        <div class="mb-3">
                            <label class="form-label">NIP</label>
                            <input type="text" class="form-control" v-model="form.nip" placeholder="Masukkan NIP">
                        </div>
                        <div class="mb-3">
                            <label class="form-label">Email <span class="text-danger">*</span></label>
                            <input type="email" class="form-control" v-model="form.email" placeholder="email@example.com">
                        </div>
                        <div v-if="errorMessage" class="alert alert-danger small py-2">
                            {{ errorMessage }}
                        </div>
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-secondary" @click="closeModal">Batal</button>
                        <button type="button" class="btn btn-primary" @click="saveTeacher" :disabled="saving">
                            <span v-if="saving" class="spinner-border spinner-border-sm me-2"></span>
                            {{ isEditMode ? 'Simpan Perubahan' : 'Tambah Guru' }}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { apiBaseUrl } from '@/config'

const teachers = ref([])
const schools = ref([])
const loading = ref(false)
const saving = ref(false)
const showModal = ref(false)
const isEditMode = ref(false)
const errorMessage = ref('')
const editingTeacherId = ref(null)

const form = reactive({
    school_id: '',
    full_name: '',
    nip: '',
    email: ''
})

const getAuthHeaders = () => {
    const token = localStorage.getItem('token')
    return token ? { Authorization: `Bearer ${token}` } : {}
}

const fetchTeachers = async () => {
    loading.value = true
    try {
        const res = await axios.get(`${apiBaseUrl}/teachers`, {
            headers: getAuthHeaders()
        })
        teachers.value = res.data || []
    } catch (error) {
        console.error('Failed to fetch teachers', error)
    } finally {
        loading.value = false
    }
}

const fetchSchools = async () => {
    try {
        const res = await axios.get(`${apiBaseUrl}/schools`, {
            headers: getAuthHeaders()
        })
        schools.value = res.data || []
    } catch (error) {
        console.error('Failed to fetch schools', error)
    }
}

const openAddModal = () => {
    isEditMode.value = false
    editingTeacherId.value = null
    form.school_id = ''
    form.full_name = ''
    form.nip = ''
    form.email = ''
    errorMessage.value = ''
    showModal.value = true
}

const editTeacher = (teacher) => {
    isEditMode.value = true
    editingTeacherId.value = teacher.user_id
    form.school_id = teacher.user?.school_id || ''
    form.full_name = teacher.full_name
    form.nip = teacher.nip || ''
    form.email = teacher.user?.email || ''
    errorMessage.value = ''
    showModal.value = true
}

const closeModal = () => {
    showModal.value = false
    errorMessage.value = ''
}

const saveTeacher = async () => {
    if (!form.school_id) {
        errorMessage.value = 'Sekolah wajib dipilih'
        return
    }
    if (!form.full_name.trim()) {
        errorMessage.value = 'Nama lengkap wajib diisi'
        return
    }
    if (!form.email.trim()) {
        errorMessage.value = 'Email wajib diisi'
        return
    }

    saving.value = true
    errorMessage.value = ''

    try {
        if (isEditMode.value) {
            // Update existing teacher
            await axios.put(`${apiBaseUrl}/teachers/${editingTeacherId.value}`, {
                school_id: Number(form.school_id),
                full_name: form.full_name,
                nip: form.nip || null,
                email: form.email
            }, {
                headers: getAuthHeaders()
            })
        } else {
            // Create new teacher
            await axios.post(`${apiBaseUrl}/teachers`, {
                school_id: form.school_id,
                full_name: form.full_name,
                nip: form.nip || null,
                email: form.email
            }, {
                headers: getAuthHeaders()
            })
        }

        closeModal()
        fetchTeachers()
    } catch (error) {
        console.error('Failed to save teacher', error)
        errorMessage.value = error.response?.data?.error || 'Gagal menyimpan data guru'
    } finally {
        saving.value = false
    }
}

const deleteTeacher = async (teacherId) => {
    try {
        await axios.delete(`${apiBaseUrl}/users/${teacherId}`, {
            headers: getAuthHeaders()
        })
        teachers.value = teachers.value.filter(t => t.user_id !== teacherId)
    } catch (error) {
        console.error('Failed to delete teacher', error)
        alert(error.response?.data?.error || 'Gagal menghapus data guru')
    }
}

const askDelete = (teacher) => {
    if (!teacher?.user_id) return
    if (confirm(`Hapus data guru "${teacher.full_name}"?`)) {
        deleteTeacher(teacher.user_id)
    }
}

onMounted(() => {
    fetchTeachers()
    fetchSchools()
})
</script>

<style scoped>
.avatar {
    width: 50px;
    height: 50px;
    object-fit: cover;
}
</style>