<template>
    <div class="card border-0 bg-body-tertiary">
        <div class="card-body">
            <div class="d-flex flex-wrap align-items-center justify-content-between mb-3 gap-2">
                <div>
                    <div class="d-flex align-items-center gap-2 mb-1">
                        <h5 class="mb-0 fw-bold">Laporan Grafis Kehadiran</h5>
                    </div>
                    <p class="text-muted small mb-0">
                        Visualisasi data kehadiran siswa dalam bentuk grafik interaktif.
                    </p>
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
                                Tampilkan Grafik
                            </button>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Loading State -->
            <div v-if="loading" class="text-center py-5">
                <div class="spinner-border text-primary mb-3" role="status"></div>
                <p class="text-muted">Memuat data grafik...</p>
            </div>

            <!-- No Data State -->
            <div v-else-if="!hasData" class="text-center py-5">
                <i class="bi bi-bar-chart-line display-4 text-secondary mb-3 d-block"></i>
                <p class="text-muted">Pilih periode dan klik "Tampilkan Grafik" untuk melihat visualisasi data.</p>
            </div>

            <!-- Charts Grid -->
            <div v-else class="row g-4">
                <!-- Row 1: Pie Chart and Bar Chart -->
                <div class="col-lg-5">
                    <div class="card border-0 shadow-sm h-100">
                        <div class="card-header bg-transparent border-0 py-3">
                            <h6 class="mb-0 fw-semibold">
                                <i class="bi bi-pie-chart me-2 text-primary"></i>
                                Status Kehadiran
                            </h6>
                        </div>
                        <div class="card-body">
                            <Doughnut :data="statusChartData" :options="doughnutOptions" />
                        </div>
                    </div>
                </div>

                <div class="col-lg-7">
                    <div class="card border-0 shadow-sm h-100">
                        <div class="card-header bg-transparent border-0 py-3">
                            <h6 class="mb-0 fw-semibold">
                                <i class="bi bi-bar-chart me-2 text-success"></i>
                                Kehadiran per Kelas
                            </h6>
                        </div>
                        <div class="card-body">
                            <Bar :data="classChartData" :options="barOptions" />
                        </div>
                    </div>
                </div>

                <!-- Row 2: Line Chart (Full Width) -->
                <div class="col-12">
                    <div class="card border-0 shadow-sm">
                        <div class="card-header bg-transparent border-0 py-3">
                            <h6 class="mb-0 fw-semibold">
                                <i class="bi bi-graph-up me-2 text-info"></i>
                                Trend Kehadiran Harian
                            </h6>
                        </div>
                        <div class="card-body">
                            <Line :data="dailyTrendData" :options="lineOptions" />
                        </div>
                    </div>
                </div>

                <!-- Row 3: Horizontal Bar and Area Chart -->
                <div class="col-lg-6">
                    <div class="card border-0 shadow-sm h-100">
                        <div class="card-header bg-transparent border-0 py-3">
                            <h6 class="mb-0 fw-semibold">
                                <i class="bi bi-clock-history me-2 text-warning"></i>
                                Top 10 Siswa Terlambat
                            </h6>
                        </div>
                        <div class="card-body">
                            <Bar :data="topLateData" :options="horizontalBarOptions" />
                        </div>
                    </div>
                </div>

                <div class="col-lg-6">
                    <div class="card border-0 shadow-sm h-100">
                        <div class="card-header bg-transparent border-0 py-3">
                            <h6 class="mb-0 fw-semibold">
                                <i class="bi bi-graph-up-arrow me-2 text-danger"></i>
                                Akumulasi Kehadiran
                            </h6>
                        </div>
                        <div class="card-body">
                            <Line :data="accumulationData" :options="areaOptions" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { apiBaseUrl } from '@/config'
import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    BarElement,
    ArcElement,
    Title,
    Tooltip,
    Legend,
    Filler
} from 'chart.js'
import { Bar, Line, Doughnut } from 'vue-chartjs'

// Register Chart.js components
ChartJS.register(
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    BarElement,
    ArcElement,
    Title,
    Tooltip,
    Legend,
    Filler
)

const loading = ref(false)
const classes = ref([])
const reportData = ref([])
const hasData = computed(() => reportData.value.length > 0)

const filters = reactive({
    startDate: new Date().toISOString().split('T')[0],
    endDate: new Date().toISOString().split('T')[0],
    classId: ''
})

// Chart color palette
const colors = {
    success: 'rgba(25, 135, 84, 0.8)',
    successLight: 'rgba(25, 135, 84, 0.2)',
    warning: 'rgba(255, 193, 7, 0.8)',
    warningLight: 'rgba(255, 193, 7, 0.2)',
    danger: 'rgba(220, 53, 69, 0.8)',
    dangerLight: 'rgba(220, 53, 69, 0.2)',
    primary: 'rgba(13, 110, 253, 0.8)',
    primaryLight: 'rgba(13, 110, 253, 0.2)',
    info: 'rgba(13, 202, 240, 0.8)',
    infoLight: 'rgba(13, 202, 240, 0.2)'
}

// Helper: Get attendance status
const getStatus = (item) => {
    if (!item.check_in_at) return 'absent'
    const checkIn = new Date(item.check_in_at.replace('Z', ''))
    const limit = new Date(checkIn)
    limit.setHours(7, 30, 0, 0)
    return checkIn > limit ? 'late' : 'onTime'
}

// 1. Doughnut Chart - Status Kehadiran
const statusChartData = computed(() => {
    const counts = { onTime: 0, late: 0, absent: 0 }
    reportData.value.forEach(item => {
        const status = getStatus(item)
        counts[status]++
    })

    return {
        labels: ['Tepat Waktu', 'Terlambat', 'Tidak Hadir'],
        datasets: [{
            data: [counts.onTime, counts.late, counts.absent],
            backgroundColor: [colors.success, colors.warning, colors.danger],
            borderWidth: 0,
            hoverOffset: 10
        }]
    }
})

const doughnutOptions = {
    responsive: true,
    maintainAspectRatio: true,
    plugins: {
        legend: {
            position: 'bottom',
            labels: { padding: 20, usePointStyle: true }
        }
    },
    cutout: '60%'
}

// 2. Bar Chart - Kehadiran per Kelas
const classChartData = computed(() => {
    const classStats = {}

    reportData.value.forEach(item => {
        const className = item.student?.class?.name || 'Unknown'
        if (!classStats[className]) {
            classStats[className] = { onTime: 0, late: 0, absent: 0 }
        }
        const status = getStatus(item)
        classStats[className][status]++
    })

    const labels = Object.keys(classStats).sort()

    return {
        labels,
        datasets: [
            {
                label: 'Tepat Waktu',
                data: labels.map(c => classStats[c].onTime),
                backgroundColor: colors.success,
                borderRadius: 4
            },
            {
                label: 'Terlambat',
                data: labels.map(c => classStats[c].late),
                backgroundColor: colors.warning,
                borderRadius: 4
            },
            {
                label: 'Tidak Hadir',
                data: labels.map(c => classStats[c].absent),
                backgroundColor: colors.danger,
                borderRadius: 4
            }
        ]
    }
})

const barOptions = {
    responsive: true,
    maintainAspectRatio: true,
    plugins: {
        legend: {
            position: 'top',
            labels: { usePointStyle: true, padding: 15 }
        }
    },
    scales: {
        x: { grid: { display: false } },
        y: { beginAtZero: true, grid: { color: 'rgba(0,0,0,0.05)' } }
    }
}

// 3. Line Chart - Trend Harian
const dailyTrendData = computed(() => {
    const dailyStats = {}

    reportData.value.forEach(item => {
        const date = item.date?.split('T')[0] || ''
        if (!dailyStats[date]) {
            dailyStats[date] = { onTime: 0, late: 0, absent: 0 }
        }
        const status = getStatus(item)
        dailyStats[date][status]++
    })

    const labels = Object.keys(dailyStats).sort()
    const formattedLabels = labels.map(d => {
        const date = new Date(d)
        return new Intl.DateTimeFormat('id-ID', { day: 'numeric', month: 'short' }).format(date)
    })

    return {
        labels: formattedLabels,
        datasets: [
            {
                label: 'Tepat Waktu',
                data: labels.map(d => dailyStats[d].onTime),
                borderColor: colors.success,
                backgroundColor: colors.successLight,
                tension: 0.4,
                fill: false
            },
            {
                label: 'Terlambat',
                data: labels.map(d => dailyStats[d].late),
                borderColor: colors.warning,
                backgroundColor: colors.warningLight,
                tension: 0.4,
                fill: false
            },
            {
                label: 'Tidak Hadir',
                data: labels.map(d => dailyStats[d].absent),
                borderColor: colors.danger,
                backgroundColor: colors.dangerLight,
                tension: 0.4,
                fill: false
            }
        ]
    }
})

const lineOptions = {
    responsive: true,
    maintainAspectRatio: true,
    plugins: {
        legend: {
            position: 'top',
            labels: { usePointStyle: true, padding: 15 }
        }
    },
    scales: {
        x: { grid: { display: false } },
        y: { beginAtZero: true, grid: { color: 'rgba(0,0,0,0.05)' } }
    },
    interaction: {
        intersect: false,
        mode: 'index'
    }
}

// 4. Horizontal Bar Chart - Top 10 Siswa Terlambat
const topLateData = computed(() => {
    const studentLate = {}

    reportData.value.forEach(item => {
        if (getStatus(item) === 'late') {
            const name = item.student?.full_name || 'Unknown'
            studentLate[name] = (studentLate[name] || 0) + 1
        }
    })

    const sorted = Object.entries(studentLate)
        .sort((a, b) => b[1] - a[1])
        .slice(0, 10)

    return {
        labels: sorted.map(s => s[0]),
        datasets: [{
            label: 'Jumlah Terlambat',
            data: sorted.map(s => s[1]),
            backgroundColor: colors.warning,
            borderRadius: 4
        }]
    }
})

const horizontalBarOptions = {
    responsive: true,
    maintainAspectRatio: true,
    indexAxis: 'y',
    plugins: {
        legend: { display: false }
    },
    scales: {
        x: { beginAtZero: true, grid: { color: 'rgba(0,0,0,0.05)' } },
        y: { grid: { display: false } }
    }
}

// 5. Area Chart - Akumulasi Kehadiran
const accumulationData = computed(() => {
    const dailyStats = {}

    reportData.value.forEach(item => {
        const date = item.date?.split('T')[0] || ''
        if (!dailyStats[date]) {
            dailyStats[date] = { present: 0, absent: 0 }
        }
        const status = getStatus(item)
        if (status === 'absent') {
            dailyStats[date].absent++
        } else {
            dailyStats[date].present++
        }
    })

    const labels = Object.keys(dailyStats).sort()
    const formattedLabels = labels.map(d => {
        const date = new Date(d)
        return new Intl.DateTimeFormat('id-ID', { day: 'numeric', month: 'short' }).format(date)
    })

    // Calculate accumulation
    let accumPresent = 0
    let accumAbsent = 0
    const presentAccum = labels.map(d => {
        accumPresent += dailyStats[d].present
        return accumPresent
    })
    const absentAccum = labels.map(d => {
        accumAbsent += dailyStats[d].absent
        return accumAbsent
    })

    return {
        labels: formattedLabels,
        datasets: [
            {
                label: 'Akumulasi Hadir',
                data: presentAccum,
                borderColor: colors.primary,
                backgroundColor: colors.primaryLight,
                tension: 0.4,
                fill: true
            },
            {
                label: 'Akumulasi Tidak Hadir',
                data: absentAccum,
                borderColor: colors.danger,
                backgroundColor: colors.dangerLight,
                tension: 0.4,
                fill: true
            }
        ]
    }
})

const areaOptions = {
    responsive: true,
    maintainAspectRatio: true,
    plugins: {
        legend: {
            position: 'top',
            labels: { usePointStyle: true, padding: 15 }
        }
    },
    scales: {
        x: { grid: { display: false } },
        y: { beginAtZero: true, grid: { color: 'rgba(0,0,0,0.05)' } }
    },
    interaction: {
        intersect: false,
        mode: 'index'
    }
}

// Fetch Report Data
const generateReport = async () => {
    loading.value = true
    try {
        const token = localStorage.getItem('token')
        const params = {
            start_date: filters.startDate,
            end_date: filters.endDate,
            class_id: filters.classId,
            page_size: 1000
        }

        const response = await axios.get(`${apiBaseUrl}/attendance`, {
            headers: { Authorization: `Bearer ${token}` },
            params
        })

        reportData.value = response.data.data || []
    } catch (error) {
        console.error('Failed to generate report:', error)
    } finally {
        loading.value = false
    }
}

// Fetch Classes on mount
const fetchClasses = async () => {
    try {
        const token = localStorage.getItem('token')
        const res = await axios.get(`${apiBaseUrl}/classes`, {
            headers: { Authorization: `Bearer ${token}` }
        })
        classes.value = res.data
    } catch (err) {
        console.error('Failed to fetch classes:', err)
    }
}

// Initialize
fetchClasses()
</script>
