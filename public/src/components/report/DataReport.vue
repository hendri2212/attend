<template>
    <div class="card border-0 bg-body-tertiary">
        <div class="card-body">
            <div class="d-flex flex-wrap align-items-center justify-content-between mb-3 gap-2">
                <div>
                    <div class="d-flex align-items-center gap-2 mb-1">
                        <h5 class="mb-0 fw-bold">Laporan Absensi</h5>
                    </div>
                    <p class="text-muted small mb-0">
                        Buat dan unduh laporan absensi siswa berdasarkan periode dan kelas.
                    </p>
                </div>
                <div class="d-flex gap-2">
                    <button class="btn btn-outline-secondary d-flex align-items-center" disabled>
                        <i class="bi bi-printer me-2"></i>
                        Cetak
                    </button>
                    <button class="btn btn-primary d-flex align-items-center" disabled>
                        <i class="bi bi-download me-2"></i>
                        Export
                    </button>
                </div>
            </div>
            <hr>

            <!-- Filters -->
            <div class="card border-0 shadow-sm mb-4">
                <div class="card-body">
                    <div class="row g-3 align-items-end">
                        <div class="col-md-3">
                            <label class="form-label small fw-semibold">Tanggal Mulai</label>
                            <input type="date" class="form-control" v-model="filters.startDate">
                        </div>
                        <div class="col-md-3">
                            <label class="form-label small fw-semibold">Tanggal Selesai</label>
                            <input type="date" class="form-control" v-model="filters.endDate">
                        </div>
                        <div class="col-md-3">
                            <label class="form-label small fw-semibold">Kelas</label>
                            <select class="form-select" v-model="filters.classId">
                                <option value="">Semua Kelas</option>
                                <option v-for="cls in classes" :key="cls.id" :value="cls.id">
                                    {{ cls.name }}
                                </option>
                            </select>
                        </div>
                        <div class="col-md-3">
                            <button class="btn btn-primary w-100" @click="generateReport" :disabled="loading">
                                <span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status"
                                    aria-hidden="true"></span>
                                Tampilkan Laporan
                            </button>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Summary Cards -->
            <div class="row g-3 mb-4">
                <div class="col-md-3">
                    <div class="card border-0 shadow-sm h-100 border-start border-4 border-primary">
                        <div class="card-body">
                            <h6 class="text-muted small mb-1">Total Kehadiran</h6>
                            <h3 class="mb-0 fw-bold">{{ summary.total }}</h3>
                        </div>
                    </div>
                </div>
                <div class="col-md-3">
                    <div class="card border-0 shadow-sm h-100 border-start border-4 border-success">
                        <div class="card-body">
                            <h6 class="text-muted small mb-1">Tepat Waktu</h6>
                            <h3 class="mb-0 fw-bold text-success">{{ summary.onTime }}</h3>
                        </div>
                    </div>
                </div>
                <div class="col-md-3">
                    <div class="card border-0 shadow-sm h-100 border-start border-4 border-warning">
                        <div class="card-body">
                            <h6 class="text-muted small mb-1">Terlambat</h6>
                            <h3 class="mb-0 fw-bold text-warning">{{ summary.late }}</h3>
                        </div>
                    </div>
                </div>
                <div class="col-md-3">
                    <div class="card border-0 shadow-sm h-100 border-start border-4 border-danger">
                        <div class="card-body">
                            <h6 class="text-muted small mb-1">Tidak Hadir</h6>
                            <h3 class="mb-0 fw-bold text-danger">{{ summary.absent }}</h3>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Data Table -->
            <div class="table-responsive">
                <table class="table table-hover align-middle mb-0">
                    <thead class="table-light">
                        <tr>
                            <th scope="col">Tanggal</th>
                            <th scope="col">NISN</th>
                            <th scope="col">Siswa</th>
                            <th scope="col">Kelas</th>
                            <th scope="col">Masuk</th>
                            <th scope="col">Pulang</th>
                            <th scope="col">Status</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-if="loading">
                            <td colspan="7" class="text-center text-muted py-5">
                                <div class="spinner-border text-primary mb-2" role="status"></div>
                                <p class="mb-0 small">Memuat data laporan...</p>
                            </td>
                        </tr>
                        <tr v-else-if="reportData.length === 0">
                            <td colspan="7" class="text-center text-muted py-5">
                                <i class="bi bi-clipboard-data display-6 mb-3 d-block text-secondary"></i>
                                <p class="mb-0">Tidak ada data untuk periode yang dipilih.</p>
                            </td>
                        </tr>
                        <tr v-else v-for="(item, index) in reportData" :key="index">
                            <td class="text-nowrap">{{ formatDate(item.date) }}</td>
                            <td class="text-muted">{{ item.student?.nisn || '-' }}</td>
                            <td class="fw-semibold">{{ item.student?.full_name || '-' }}</td>
                            <td>{{ item.student?.class?.name || '-' }}</td>
                            <td>
                                <span v-if="item.check_in_at" class="badge bg-success-subtle text-success-emphasis">
                                    {{ formatTime(item.check_in_at) }}
                                </span>
                                <span v-else>-</span>
                            </td>
                            <td>
                                <span v-if="item.check_out_at" class="badge bg-primary-subtle text-primary-emphasis">
                                    {{ formatTime(item.check_out_at) }}
                                </span>
                                <span v-else>-</span>
                            </td>
                            <td>
                                <span :class="getStatusBadgeClass(item)">
                                    {{ getStatusLabel(item) }}
                                </span>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { apiBaseUrl } from '@/config'

const loading = ref(false)
const classes = ref([])
const filters = reactive({
    startDate: new Date().toISOString().split('T')[0],
    endDate: new Date().toISOString().split('T')[0],
    classId: ''
})

const summary = reactive({
    total: 0,
    onTime: 0,
    late: 0,
    absent: 0
})

const reportData = ref([])

const formatDate = (dateString) => {
    if (!dateString) return '-'
    const date = new Date(dateString.replace('Z', ''))
    return new Intl.DateTimeFormat('id-ID', {
        day: 'numeric',
        month: 'long',
        year: 'numeric'
    }).format(date)
}

const formatTime = (timeString) => {
    if (!timeString) return '-'
    const date = new Date(timeString.replace('Z', ''))
    return new Intl.DateTimeFormat('id-ID', {
        hour: '2-digit',
        minute: '2-digit'
    }).format(date)
}

const getStatusBadgeClass = (item) => {
    // Logic to determine status badge color
    // This is a placeholder logic
    if (!item.check_in_at) return 'badge bg-danger-subtle text-danger-emphasis'
    // Assuming 07:00 is the cut-off for late
    const checkIn = new Date(item.check_in_at.replace('Z', ''))
    const limit = new Date(checkIn)
    limit.setHours(7, 30, 0, 0)
    
    if (checkIn > limit) return 'badge bg-warning-subtle text-warning-emphasis'
    return 'badge bg-success-subtle text-success-emphasis'
}

const getStatusLabel = (item) => {
    if (!item.check_in_at) return 'Tidak Hadir'
    const checkIn = new Date(item.check_in_at.replace('Z', ''))
    const limit = new Date(checkIn)
    limit.setHours(7, 30, 0, 0)
    
    if (checkIn > limit) return 'Terlambat'
    return 'Tepat Waktu'
}

const generateReport = async () => {
    loading.value = true
    try {
        const token = localStorage.getItem('token')
        const params = {
            start_date: filters.startDate,
            end_date: filters.endDate,
            class_id: filters.classId,
            page_size: 1000 // Fetch all for report (or implement pagination if needed)
        }
        
        const response = await axios.get(`${apiBaseUrl}/attendance`, {
            headers: {
                Authorization: `Bearer ${token}`
            },
            params
        })
        
        reportData.value = response.data.data || []
        
        // Update Summary based on fetched data
        summary.total = reportData.value.length
        summary.onTime = reportData.value.filter(i => getStatusLabel(i) === 'Tepat Waktu').length
        summary.late = reportData.value.filter(i => getStatusLabel(i) === 'Terlambat').length
        summary.absent = reportData.value.filter(i => getStatusLabel(i) === 'Tidak Hadir').length
        
    } catch (error) {
        console.error('Failed to generate report:', error)
        // Handle error (e.g. show toast)
    } finally {
        loading.value = false
    }
}

onMounted(async () => {
    try {
        const token = localStorage.getItem('token')
        const res = await axios.get(`${apiBaseUrl}/classes`, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        })
        classes.value = res.data
    } catch (err) {
        console.error(err)
    }
})
</script>
