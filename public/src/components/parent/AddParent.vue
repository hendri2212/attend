<template>
    <div class="card border-0 bg-body-tertiary" style="height: calc(100vh - 120px);">
        <div class="card-body">
            <form @submit.prevent="submitForm">
                <div class="col-6">
                    <div class="form-floating mb-3">
                        <input type="text" id="fullName" v-model="parent.full_name" class="form-control"
                            placeholder="Nama Lengkap" required>
                        <label for="fullName">Nama Lengkap</label>
                    </div>
                    <div class="form-floating mb-3">
                        <input type="text" id="whatsapp" v-model="parent.whatsapp" class="form-control"
                            placeholder="WhatsApp" required>
                        <label for="whatsapp">WhatsApp</label>
                    </div>
                    <div class="form-floating mb-3">
                        <input type="text" id="address" v-model="parent.address" class="form-control"
                            placeholder="Alamat" required>
                        <label for="address">Alamat</label>
                    </div>
                    <div class="btn-group float-end">
                        <button type="button" class="btn btn-secondary" @click="() => $router.push('/admin/students')">
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
    full_name: '',
    whatsapp: '',
    address: ''
})

const submitForm = async () => {
    try {
        const token = localStorage.getItem('token')
        await axios.post(
            `${apiBaseUrl}/parents`,
            parent,
            {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            }
        )
        alert('Data orang tua berhasil disimpan!')
        router.push({ name: 'parents' })
    } catch (err) {
        console.error(err)
        alert('Gagal menyimpan data orang tua')
    }
}
</script>