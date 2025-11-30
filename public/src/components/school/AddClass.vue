<template>
    <div class="card border-0 bg-body-tertiary" style="height: calc(100vh - 120px);">
        <div class="card-body">
            <div class="d-flex flex-wrap align-items-center justify-content-start mb-3 gap-2">
                <div>
                    <div class="d-flex align-items-center gap-2 mb-1">
                        <h5 class="mb-0 fw-bold">Tambah Kelas</h5>
                    </div>
                    <p class="text-muted small mb-0">
                        Kelola daftar kelas, tambah, ubah, dan hapus sesuai kebutuhan.
                    </p>
                </div>
            </div>
            <hr>
            <form @submit.prevent="submitForm">
                <div class="col-6">
                    <div class="form-floating mb-3">
                        <input type="text" id="name" v-model="parent.name" class="form-control"
                            placeholder="Nama Lengkap" required>
                        <label for="name">Nama Kelas</label>
                    </div>
                    <div class="btn-group float-end">
                        <button type="button" class="btn btn-secondary" @click="() => $router.push('/admin/classes')">
                            Cancel
                        </button>
                        <button type="submit" class="btn btn-primary">Save</button>
                    </div>
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
const parent = reactive({
    name: ''
})

const submitForm = async () => {
    try {
        const token = localStorage.getItem('token')
        await axios.post(
            `${apiBaseUrl}/classes`,
            parent,
            {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            }
        )
        alert('Data kelas berhasil disimpan!')
        router.push({ name: 'classes' })
    } catch (err) {
        console.error(err)
        alert('Gagal menyimpan data kelas')
    }
}
</script>