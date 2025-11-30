<template>
    <div class="card border-0 bg-body-tertiary" style="height: calc(100vh - 120px);">
        <div class="card-body">
            <div class="d-flex flex-wrap align-items-center justify-content-start mb-3 gap-2">
                <div>
                    <div class="d-flex align-items-center gap-2 mb-1">
                        <h5 class="mb-0 fw-bold">Tambah Siswa</h5>
                    </div>
                    <p class="text-muted small mb-0">
                        Kelola daftar siswa, tambah, ubah, dan hapus sesuai kebutuhan.
                    </p>
                </div>
            </div>
            <hr>
            <form @submit.prevent="submitForm">
                <div class="d-flex gap-2">
                    <div class="col">
                        <div class="form-floating mb-3">
                            <input type="text" id="rfid" v-model="student.rfid" class="form-control" placeholder="RFID"
                                required>
                            <label for="rfid">RFID</label>
                        </div>
                        <div class="form-floating mb-3">
                            <input type="email" id="email" v-model="student.email" class="form-control" placeholder="Email"
                                required>
                            <label for="email">Email</label>
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
                    </div>
                    <div class="col">
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
                        <div class="form-floating mb-3">
                            <select id="class" v-model="student.class_id" class="form-select" required>
                                <option disabled value="">Pilih Kelas</option>
                                <option v-for="cls in classes" :key="cls.id" :value="cls.id">
                                    {{ cls.name }}
                                </option>
                            </select>
                            <label for="class">Kelas</label>
                        </div>
                    </div>
                </div>
                <div class="btn-group float-end">
                    <button type="button" class="btn btn-secondary" @click="() => $router.push('/admin/students')">
                        Cancel
                    </button>
                    <button type="submit" class="btn btn-primary">Save</button>
                </div>
            </form>
        </div>
    </div>
</template>
<script setup>
import { useRouter } from 'vue-router'
import { reactive, onMounted } from 'vue'
import { apiBaseUrl } from '@/config'

const router = useRouter()
const student = reactive({
    rfid: '',
    email: '',
    nisn: '',
    full_name: '',
    birth_place: '',
    born: '',
    whatsapp: '',
    class_id: ''
})

const classes = reactive([])

onMounted(async () => {
    try {
        const token = localStorage.getItem('token')
        const res = await axios.get(`${apiBaseUrl}/classes`, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        })
        classes.push(...res.data)
    } catch (err) {
        console.error(err)
    }
})

const submitForm = async () => {
    try {
        const token = localStorage.getItem('token')
        await axios.post(
            `${apiBaseUrl}/students`,
            student,
            {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            }
        )
        alert('Data siswa berhasil disimpan!')
        router.push({ name: 'students' })
    } catch (err) {
        console.error(err)
        alert('Gagal menyimpan data siswa')
    }
}
</script>