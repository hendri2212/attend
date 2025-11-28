<template>
    <div class="card border-0 bg-body-tertiary" style="height: calc(100vh - 120px);">
        <div class="card-body">
            <form @submit.prevent="submitForm">
                <div class="col-6">
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
                        <input type="text" id="address" v-model="student.address" class="form-control"
                            placeholder="Alamat">
                        <label for="address">Alamat</label>
                    </div>
                    <div class="btn-group float-end">
                        <button type="button" class="btn btn-secondary" @click="() => router.push('/admin/parents')">
                            Cancel
                        </button>
                        <button type="submit" class="btn btn-primary">Update</button>
                    </div>
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
    full_name: '',
    whatsapp: '',
    address: ''
})

/**
 * Fetch parent data by ID and populate the form.
 */
const fetchParent = async () => {
    try {
        const token = localStorage.getItem('token')
        if (!token) {
            alert('Token tidak ditemukan, silakan login kembali.')
            router.push('/login')
            return
        }

        const { data } = await axios.get(`${apiBaseUrl}/parents/${studentId}`, {
            headers: { Authorization: `Bearer ${token}` }
        })

        // Map server response to the reactive object
        Object.assign(student, {
            full_name: data.full_name ?? '',
            whatsapp: data.whatsapp ?? '',
            address: data.address ?? ''
        })
    } catch (err) {
        console.error(err)
        alert('Gagal memuat data orang tua')
    }
}

onMounted(() => {
    fetchParent()
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
            `${apiBaseUrl}/parents/${studentId}`,
            student,
            { headers: { Authorization: `Bearer ${token}` } }
        )

        alert('Data orang tua berhasil diperbarui!')
        router.push('/admin/parents')
    } catch (err) {
        console.error(err)
        alert('Gagal memperbarui data orang tua')
    }
}
</script>