<template>
    <div class="card border-0 bg-body-tertiary" style="height: calc(100vh - 120px);">
        <div class="card-body">
            <form @submit.prevent="submitForm">
                <div class="d-flex gap-2">
                    <div class="col">
                        <div class="form-floating mb-3">
                            <input type="text" id="rfid" v-model="student.rfid" class="form-control" placeholder="RFID"
                                required>
                            <label for="rfid">RFID</label>
                        </div>
                        <div class="form-floating mb-3">
                            <input type="text" id="nisn" v-model="student.nisn" class="form-control" placeholder="NISN">
                            <label for="nisn">NISN</label>
                        </div>
                        <div class="form-floating mb-3">
                            <input type="text" id="fullName" v-model="student.full_name" class="form-control"
                                placeholder="Nama Lengkap" required>
                            <label for="fullName">Nama Lengkap</label>
                        </div>
                        <div class="form-floating mb-3">
                            <input type="text" id="whatsapp" v-model="student.whatsapp" class="form-control"
                                placeholder="No. WhatsApp">
                            <label for="whatsapp">No. WhatsApp</label>
                        </div>
                        <div class="form-floating mb-3">
                            <input type="text" id="birthPlace" v-model="student.birth_place" class="form-control"
                                placeholder="Tempat Lahir">
                            <label for="birthPlace">Tempat Lahir</label>
                        </div>
                        <div class="form-floating mb-3">
                            <input type="date" id="bornDate" v-model="student.born" class="form-control"
                                placeholder="Tanggal Lahir">
                            <label for="bornDate">Tanggal Lahir</label>
                        </div>
                    </div>
                    <div class="col">
                        <div class="form-floating mb-3">
                            <select id="parentId" v-model="student.parent_id" class="form-select">
                                <option disabled value="">Pilih Orang Tua</option>
                                <option v-for="parent in parents" :key="parent.id" :value="parent.id">
                                    {{ parent.full_name }}
                                </option>
                            </select>
                            <label for="parentId">Orang Tua</label>
                        </div>
                    </div>
                </div>
                <div class="btn-group float-end">
                    <button type="button" class="btn btn-secondary" @click="() => router.push('/admin/students')">
                        Cancel
                    </button>
                    <button type="submit" class="btn btn-primary">Update</button>
                </div>
            </form>
        </div>
    </div>
</template>
<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { apiBaseUrl } from '@/config'

const route = useRoute()
const router = useRouter()
const studentId = route.params.id

const parents = ref([])

// Reactive object to bind to the form
const student = reactive({
    rfid: '',
    nisn: '',
    full_name: '',
    birth_place: '',
    born: '',
    whatsapp: '',
    parent_id: ''
})

/**
 * Fetch student data by ID and populate the form.
 */
const fetchStudent = async () => {
    try {
        const token = localStorage.getItem('token')
        if (!token) {
            alert('Token tidak ditemukan, silakan login kembali.')
            router.push('/login')
            return
        }

        const { data } = await axios.get(`${apiBaseUrl}/students/${studentId}`, {
            headers: { Authorization: `Bearer ${token}` }
        })

        // Map server response to the reactive object
        Object.assign(student, {
            rfid: data.rfid ?? '',
            nisn: data.nisn ?? '',
            full_name: data.full_name ?? '',
            birth_place: data.birth_place ?? '',
            born: data.born ? data.born.split('T')[0] : '',
            whatsapp: data.whatsapp ?? '',
            parent_id: data.parent?.id ?? ''
        })
    } catch (err) {
        console.error(err)
        alert('Gagal memuat data siswa')
    }
}

const fetchParents = async () => {
    try {
        const token = localStorage.getItem('token')
        if (!token) return

        const { data } = await axios.get(`${apiBaseUrl}/parents`, {
            headers: { Authorization: `Bearer ${token}` }
        })
        parents.value = data
    } catch (err) {
        console.error('Gagal memuat data orang tua:', err)
    }
}

onMounted(() => {
    fetchStudent()
    fetchParents()
})

/**
 * Update student data.
 */
const submitForm = async () => {
    try {
        const token = localStorage.getItem('token')
        if (!token) {
            alert('Token tidak ditemukan, silakan login kembali.')
            return
        }

        await axios.put(
            `${apiBaseUrl}/students/${studentId}`,
            student,
            { headers: { Authorization: `Bearer ${token}` } }
        )

        alert('Data siswa berhasil diperbarui!')
        router.push('/admin/students')
    } catch (err) {
        console.error(err)
        alert('Gagal memperbarui data siswa')
    }
}
</script>