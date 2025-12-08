<template>
    <div class="card border-0 bg-body-tertiary">
        <div class="card-body">
            <div class="d-flex flex-wrap align-items-center justify-content-between mb-3 gap-2">
                <div>
                    <div class="d-flex align-items-center gap-2 mb-1">
                        <h5 class="mb-0 fw-bold">Manajemen User</h5>
                        <span class="badge rounded-pill bg-secondary-subtle text-primary">{{ users.length }} Users</span>
                    </div>
                    <p class="text-muted small mb-0">
                        Kelola data pengguna aplikasi (Superadmin Only).
                    </p>
                </div>
                <button class="btn btn-outline-primary d-flex align-items-center" @click="openCreateModal">
                    <i class="bi bi-plus-lg me-2"></i>User
                </button>
            </div>
            <hr>
            
            <div v-if="loading" class="text-center py-4">
                <div class="spinner-border text-primary" role="status"></div>
            </div>

            <div v-else-if="users.length === 0" class="text-center text-muted py-4 small">
                Tidak ada data user.
            </div>

            <div v-else class="table-responsive">
                <table class="table table-hover align-middle">
                    <thead>
                        <tr>
                            <th>Email</th>
                            <th>Role</th>
                            <th>School ID</th>
                            <th>Status/Aksi</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="user in users" :key="user.id">
                            <td>{{ user.email }}</td>
                            <td>
                                <span class="badge" :class="getRoleBadge(user.role)">{{ user.role }}</span>
                            </td>
                            <td>{{ user.school_id }}</td>
                            <td>
                                <div class="btn-group btn-group-sm">
                                    <button @click="editUser(user)" class="btn btn-outline-warning" title="Edit">
                                        <i class="bi bi-pencil"></i>
                                    </button>
                                    <button @click="deleteUser(user)" class="btn btn-outline-danger" title="Delete">
                                        <i class="bi bi-trash"></i>
                                    </button>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <!-- Modal Create/Edit -->
        <div v-if="showModal" class="modal fade show d-block" style="background: rgba(0,0,0,0.5);" tabindex="-1">
            <div class="modal-dialog">
                <div class="modal-content">
                    <div class="modal-header">
                        <h5 class="modal-title">{{ isEdit ? 'Edit User' : 'Tambah User' }}</h5>
                        <button type="button" class="btn-close" @click="closeModal"></button>
                    </div>
                    <form @submit.prevent="submitForm">
                        <div class="modal-body">
                            <div class="mb-3">
                                <label class="form-label">Email</label>
                                <input type="email" v-model="form.email" class="form-control" required>
                            </div>
                            <div class="mb-3">
                                <label class="form-label">Password</label>
                                <input type="password" v-model="form.password" class="form-control" :placeholder="isEdit ? '(Biarkan kosong jika tidak diganti)' : ''" :required="!isEdit">
                            </div>
                            <div class="mb-3">
                                <label class="form-label">Role</label>
                                <select v-model="form.role" class="form-select" required>
                                    <option value="student">Student</option>
                                    <option value="teacher">Teacher</option>
                                    <option value="admin">Admin</option>
                                    <option value="superadmin">Superadmin</option>
                                </select>
                            </div>
                            <div class="mb-3">
                                <label class="form-label">School</label>
                                <select v-model="form.school_id" class="form-select" required>
                                    <option value="" disabled>Pilih Sekolah</option>
                                    <option v-for="school in schools" :key="school.id" :value="school.id">
                                        {{ school.name }}
                                    </option>
                                </select>
                            </div>
                        </div>
                        <div class="modal-footer">
                            <button type="button" class="btn btn-secondary" @click="closeModal">Batal</button>
                            <button type="submit" class="btn btn-primary">Simpan</button>
                        </div>
                    </form>
                </div>
            </div>
        </div>

    </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { apiBaseUrl } from '@/config'

const users = ref([])
const loading = ref(false)
const showModal = ref(false)
const isEdit = ref(false)
const editId = ref(null)

const schools = ref([])
const form = reactive({
    email: '',
    password: '',
    role: 'student',
    school_id: ''
})

const getAuthHeaders = () => {
    const token = localStorage.getItem('token')
    return token ? { Authorization: `Bearer ${token}` } : {}
}

const fetchSchools = async () => {
    try {
        const { data } = await axios.get(`${apiBaseUrl}/schools`, {
            headers: getAuthHeaders()
        })
        schools.value = data || []
    } catch (err) {
        console.error('Failed to fetch schools', err)
    }
}

const fetchUsers = async () => {
    loading.value = true
    try {
        const { data } = await axios.get(`${apiBaseUrl}/users`, {
            headers: getAuthHeaders()
        })
        users.value = data || []
    } catch (err) {
        console.error(err)
        if (err.response?.status === 401) {
            // Token expired or invalid
        }
    } finally {
        loading.value = false
    }
}

const getRoleBadge = (role) => {
    switch (role) {
        case 'superadmin': return 'text-bg-danger'
        case 'admin': return 'text-bg-warning'
        case 'teacher': return 'text-bg-info'
        case 'student': return 'text-bg-secondary'
        default: return 'text-bg-light'
    }
}

const openCreateModal = () => {
    isEdit.value = false
    editId.value = null
    form.email = ''
    form.password = ''
    form.role = 'student'
    form.school_id = schools.value.length > 0 ? schools.value[0].id : ''
    showModal.value = true
}

const editUser = (user) => {
    isEdit.value = true
    editId.value = user.id
    form.email = user.email
    form.password = ''
    form.role = user.role
    form.school_id = user.school_id
    showModal.value = true
}

const closeModal = () => {
    showModal.value = false
}

const submitForm = async () => {
    try {
        const payload = { ...form }
        // Remove empty password if edit
        if (isEdit.value && !payload.password) delete payload.password
        
        // Ensure school_id is number
        payload.school_id = Number(payload.school_id)

        if (isEdit.value) {
            await axios.put(`${apiBaseUrl}/users/${editId.value}`, payload, {
                headers: getAuthHeaders()
            })
        } else {
            await axios.post(`${apiBaseUrl}/users`, payload, {
                headers: getAuthHeaders()
            })
        }
        closeModal()
        fetchUsers()
        alert('Data berhasil disimpan')
    } catch (err) {
        console.error(err)
        alert('Gagal menyimpan data: ' + (err.response?.data?.error || err.message))
    }
}

// ... existing code ...

onMounted(() => {
    fetchUsers()
    fetchSchools()
})
</script>
