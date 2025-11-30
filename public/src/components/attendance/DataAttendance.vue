<template>
    <div class="card border-0 bg-body-tertiary">
        <div class="card-body">
            <div class="d-flex flex-wrap align-items-center justify-content-between mb-3 gap-2">
                <div>
                    <div class="d-flex align-items-center gap-2 mb-1">
                        <h5 class="mb-0 fw-bold">Data Absensi</h5>
                    </div>
                    <p class="text-muted small mb-0">
                        Kelola data absensi siswa, tinjau catatan kehadiran harian.
                    </p>
                </div>
                <button @click="fetchAttendance" class="btn btn-outline-primary d-flex align-items-center">
                    <i class="bi bi-arrow-clockwise me-2"></i>
                    Refresh
                </button>
            </div>
            <hr>
            <div class="table-responsive">
                <table class="table table-hover align-middle mb-0">
                    <thead class="table-light">
                        <tr>
                            <th scope="col">Tanggal</th>
                            <th scope="col">Siswa</th>
                            <th scope="col">NISN</th>
                            <th scope="col">Kelas</th>
                            <th scope="col">Masuk</th>
                            <th scope="col">Pulang</th>
                            <th scope="col">Metode</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-if="loading">
                            <td colspan="7" class="text-center text-muted py-4">
                                Loading...
                            </td>
                        </tr>
                        <tr v-else-if="attendances.length === 0">
                            <td colspan="7" class="text-center text-muted py-4">
                                Tidak ada data absensi
                            </td>
                        </tr>
                        <tr v-else v-for="item in attendances" :key="item.id">
                            <td class="text-nowrap">
                                {{ formatDate(item.date) }}
                            </td>
                            <td class="fw-semibold text-dark">
                                {{ item.student?.full_name || '-' }}
                            </td>
                            <td class="text-muted text-nowrap">
                                {{ item.student?.nisn || '-' }}
                            </td>
                            <td class="text-muted text-nowrap">
                                {{ item.student?.class?.name || '-' }}
                            </td>
                            <td class="text-nowrap">
                                <span v-if="item.check_in_at"
                                    class="badge bg-success-subtle text-success-emphasis px-2 py-1 small">
                                    {{ formatTime(item.check_in_at) }}
                                </span>
                                <span v-else class="text-muted">-</span>
                            </td>
                            <td class="text-nowrap">
                                <span v-if="item.check_out_at"
                                    class="badge bg-primary-subtle text-primary-emphasis px-2 py-1 small">
                                    {{ formatTime(item.check_out_at) }}
                                </span>
                                <span v-else class="text-muted">-</span>
                            </td>
                            <td class="text-nowrap">
                                <span class="badge bg-secondary-subtle text-secondary text-capitalize">
                                    {{ item.method || '-' }}
                                </span>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <div class="d-flex justify-content-between align-items-center mt-3">
                <p class="small text-muted mb-0">
                    Halaman <span class="fw-semibold">{{ page }}</span>
                </p>
                <div class="btn-group" role="group" aria-label="Pagination">
                    <button type="button" class="btn btn-outline-secondary btn-sm" @click="changePage(page - 1)"
                        :disabled="page === 1 || loading">
                        Previous
                    </button>
                    <button type="button" class="btn btn-outline-secondary btn-sm" @click="changePage(page + 1)"
                        :disabled="attendances.length < pageSize || loading">
                        Next
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { apiBaseUrl } from '@/config'

const attendances = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)

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
        minute: '2-digit',
        second: '2-digit'
    }).format(date)
}

const fetchAttendance = async () => {
    loading.value = true
    try {
        const token = localStorage.getItem('token')
        const response = await axios.get(
            apiBaseUrl + `/attendance?page=${page.value}&page_size=${pageSize.value}`,
            {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            }
        )
        attendances.value = response.data.data || []
    } catch (error) {
        console.error('Failed to fetch attendance:', error)
        // Bisa tambahkan toast/alert di sini kalau mau
    } finally {
        loading.value = false
    }
}

const changePage = (newPage) => {
    if (newPage < 1) return
    page.value = newPage
    fetchAttendance()
}

onMounted(() => {
    fetchAttendance()
})
</script>
