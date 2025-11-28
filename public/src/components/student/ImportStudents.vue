<template>
    <div class="card border-0 bg-body-tertiary shadow-sm">
        <div class="card-body">
            <div class="d-flex align-items-start justify-content-between">
                <div>
                    <h5 class="mb-1 fw-bold">Import Siswa Per Kelas</h5>
                    <p class="text-muted small mb-0">Pilih kelas terlebih dahulu, lalu unggah file template siswa.</p>
                </div>
                <a class="btn btn-outline-primary d-flex align-items-center" href="#" role="button">
                    <i class="bi bi-download me-2"></i>
                    Unduh Template
                </a>
            </div>

            <hr>

            <form class="row gy-3" @submit.prevent="handleImport">
                <div class="col-12 col-lg-6">
                    <label class="form-label fw-semibold">Pilih Kelas</label>
                    <div class="input-group">
                        <span class="input-group-text">
                            <i class="bi bi-door-open"></i>
                        </span>
                        <select v-model="selectedClass" class="form-select" required>
                            <option value="" disabled>Pilih kelas tujuan</option>
                            <option v-for="kelas in classes" :key="kelas.id" :value="kelas.id">
                                {{ kelas.name }}
                            </option>
                        </select>
                    </div>
                    <small class="text-muted">Import dilakukan per kelas, jadi pastikan memilih kelas yang benar.</small>
                </div>

                <div class="col-12 col-lg-6">
                    <label class="form-label fw-semibold">File Import</label>
                    <div
                        class="border rounded-3 p-3 d-flex align-items-center justify-content-between flex-wrap gap-3"
                        :class="selectedClass ? 'border-primary-subtle bg-white' : 'bg-secondary-subtle'"
                    >
                        <div class="d-flex align-items-center gap-3">
                            <div
                                class="rounded-circle d-flex align-items-center justify-content-center bg-primary-subtle text-primary"
                                style="width: 46px; height: 46px;"
                            >
                                <i class="bi bi-upload"></i>
                            </div>
                            <div>
                                <p class="mb-1 fw-semibold">
                                    {{ fileName || 'Belum ada file dipilih' }}
                                </p>
                                <small class="text-muted">
                                    Format yang didukung: .xlsx, .csv. Maksimal 5MB.
                                </small>
                            </div>
                        </div>
                        <div class="d-flex align-items-center gap-2">
                            <label :class="['btn', selectedClass ? 'btn-primary' : 'btn-secondary disabled']">
                                <input
                                    ref="fileInput"
                                    type="file"
                                    class="d-none"
                                    accept=".xlsx,.xls,.csv"
                                    :disabled="!selectedClass"
                                    @change="onFileChange"
                                >
                                Pilih File
                            </label>
                            <button
                                type="button"
                                class="btn btn-outline-secondary"
                                :disabled="!fileName"
                                @click="resetFile"
                            >
                                Reset
                            </button>
                        </div>
                    </div>
                    <small v-if="!selectedClass" class="text-danger d-block mt-1">
                        Pilih kelas dulu sebelum memilih file.
                    </small>
                </div>

                <div class="col-12 d-flex justify-content-end gap-2 mt-2">
                    <button type="button" class="btn btn-outline-secondary" @click="goBack">
                        Batal
                    </button>
                    <button type="submit" class="btn btn-success" :disabled="!canImport">
                        <i class="bi bi-check-circle me-2"></i>
                        Import
                    </button>
                </div>
            </form>

            <div class="mt-4 p-3 bg-light rounded-3 border">
                <div class="d-flex align-items-center mb-2">
                    <i class="bi bi-info-circle text-primary me-2"></i>
                    <span class="fw-semibold">Panduan singkat</span>
                </div>
                <ul class="mb-0 text-muted small">
                    <li>Unduh template, isi data siswa sesuai kolom yang tersedia.</li>
                    <li>Pastikan memilih kelas yang sesuai sebelum mengunggah file.</li>
                    <li>Periksa kembali data sebelum menekan tombol Import.</li>
                </ul>
            </div>
        </div>
    </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { apiBaseUrl } from '@/config'

const router = useRouter()
const classes = ref([])
const selectedClass = ref('')
const fileName = ref('')
const fileInput = ref(null)

const fetchClasses = async () => {
    try {
        const token = localStorage.getItem('token')
        const res = await axios.get(`${apiBaseUrl}/classes`, {
            headers: token ? { Authorization: `Bearer ${token}` } : {},
        })
        classes.value = res.data || []
    } catch (error) {
        console.error('Gagal memuat kelas', error)
    }
}

onMounted(() => {
    fetchClasses()
})

const onFileChange = (event) => {
    const file = event.target.files?.[0]
    fileName.value = file ? file.name : ''
}

const resetFile = () => {
    fileName.value = ''
    if (fileInput.value) {
        fileInput.value.value = ''
    }
}

const canImport = computed(() => !!selectedClass.value && !!fileName.value)

const handleImport = () => {
    // Hanya desain UI: aksi import akan dihubungkan ke backend nanti
    if (canImport.value) {
        alert(`Siap import siswa untuk kelas ID ${selectedClass.value} dengan file ${fileName.value}`)
    }
}

const goBack = () => {
    router.back()
}
</script>
