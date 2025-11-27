<template>
    <div class="card border-0 bg-body-tertiary" style="height: calc(100vh - 120px);">
        <div class="card-body">
            <form @submit.prevent="submitForm">
                <div class="d-flex gap-2">
                    <div class="col">
                        <div class="form-floating mb-3">
                            <input type="text" id="fullName" v-model="classes.name" class="form-control"
                                placeholder="Nama Kelas" required>
                            <label for="fullName">Nama Kelas</label>
                        </div>
                    </div>
                </div>
                <div class="btn-group float-end">
                    <button type="button" class="btn btn-secondary" @click="() => router.push('/admin/classes')">
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
const classId = route.params.id

const parents = ref([])

// Reactive object to bind to the form
const classes = reactive({
    name: '',
})

/**
 * Fetch class data by ID and populate the form.
 */
const fetchClasses = async () => {
    try {
        const token = localStorage.getItem('token')
        if (!token) {
            alert('Token tidak ditemukan, silakan login kembali.')
            router.push('/login')
            return
        }

        const { data } = await axios.get(`${apiBaseUrl}/classes/${classId}`, {
            headers: { Authorization: `Bearer ${token}` }
        })

        // Map server response to the reactive object
        Object.assign(classes, {
            name: data.name ?? '',
        })
    } catch (err) {
        console.error(err)
        alert('Gagal memuat data kelas')
    }
}


onMounted(() => {
    fetchClasses()
})

/**
 * Update class data.
 */
const submitForm = async () => {
    try {
        const token = localStorage.getItem('token')
        if (!token) {
            alert('Token tidak ditemukan, silakan login kembali.')
            return
        }

        await axios.put(
            `${apiBaseUrl}/classes/${classId}`,
            classes,
            { headers: { Authorization: `Bearer ${token}` } }
        )

        alert('Data kelas berhasil diperbarui!')
        router.push('/admin/classes')
    } catch (err) {
        console.error(err)
        alert('Gagal memperbarui data kelas')
    }
}
</script>