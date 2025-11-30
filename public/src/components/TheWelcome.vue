<template>
    <div class="d-flex flex-column vh-100 p-0">
        <header class="bg-primary text-white text-center py-3 d-flex align-items-center justify-content-center">
            <img src="/logo.png" alt="Logo" class="inline-block mr-6" style="width: 90px; height: auto; margin-right: 20px;" />
            <h1 class="display-1 fw-bold mb-0">SMPN 1 KOTABARU</h1>
        </header>
        <div class="d-flex flex-fill p-0">
            <!-- Left date panel -->
            <aside
                class="d-none d-md-flex flex-column justify-content-center align-items-center bg-light col-md-4 p-4 text-center">
                <h1 class="text-uppercase display-1 mb-3 fw-bold">{{ currentDay }}</h1>
                <h1 class="display-1 mb-3 fw-bold">{{ currentDate }}</h1>
                <h1 class="display-1 fw-bold text-muted">{{ currentMonth }} {{ currentYear }}</h1>
            </aside>

            <!-- Right attendance panel -->
            <main class="flex-fill d-flex flex-column bg-white">
                <!-- Content -->
                <div class="flex-fill d-flex flex-column justify-content-center align-items-center px-4">
                    <!-- Clock -->
                    <div class="display-1 fw-semibold mb-4">{{ currentTime }}</div>

                    <!-- Attendance info -->
                    <div v-if="currentAttendance" class="text-center">
                        <div class="h1 display-1 fw-bold text-capitalize text-muted">{{ currentAttendance.full_name }}</div>
                        <div class="h1 text-success mb-3">Berhasil presensi pukul {{ currentAttendance.time }} WITA
                        </div>
                        <div class="h1 display-1 badge bg-success py-2 px-4 text-uppercase">Status Hadir</div>
                    </div>

                    <!-- Error message (hidden, same styling as original) -->
                    <div v-if="errorMessage" class="text-center" style="display: none;">
                        <div class="h2 text-danger mb-3">{{ errorMessage }}</div>
                    </div>

                    <!-- RFID input -->
                    <input ref="rfidInputRef" type="text" id="rfidInput" class="form-control-plaintext text-center mb-4"
                        autofocus readonly />
                </div>

                <!-- Footer -->
                <footer class="text-center py-2">
                    <small class="text-muted">© 2025 SMPN 1 Kotabaru</small>
                </footer>
            </main>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, computed, nextTick } from 'vue'
import { apiBaseUrl } from '@/config';

// Refs
const rfidInputRef = ref(null)
const intervalId = ref(null)
const now = ref(new Date())
const currentAttendance = ref(null)
const errorMessage = ref('')

// Maps for tracking requests and cleanup
const pendingRequests = new Map()
const activeTimeouts = new Set()

// Computed properties for date/time
const currentDay = computed(() =>
    now.value.toLocaleDateString('id-ID', { weekday: 'long' })
)
const currentDate = computed(() =>
    now.value.getDate()
)
const currentMonth = computed(() =>
    now.value.toLocaleDateString('id-ID', { month: 'long' })
)
const currentYear = computed(() =>
    now.value.getFullYear()
)
const currentTime = computed(() =>
    now.value.toLocaleTimeString('id-ID')
)

// Utility functions
const throttle = (func, delay) => {
    let timeoutId
    let lastExecTime = 0
    return function (...args) {
        const currentTime = Date.now()

        if (currentTime - lastExecTime > delay) {
            func.apply(this, args)
            lastExecTime = currentTime
        } else {
            clearTimeout(timeoutId)
            timeoutId = setTimeout(() => {
                func.apply(this, args)
                lastExecTime = Date.now()
            }, delay - (currentTime - lastExecTime))
            activeTimeouts.add(timeoutId)
        }
    }
}

const debounce = (func, delay) => {
    let timeoutId
    return function (...args) {
        clearTimeout(timeoutId)
        timeoutId = setTimeout(() => {
            func.apply(this, args)
            activeTimeouts.delete(timeoutId)
        }, delay)
        activeTimeouts.add(timeoutId)
    }
}

// Core functions
const updateDateTime = throttle(() => {
    now.value = new Date()
    ensureFocus()
}, 1000)

const ensureFocus = () => {
    if (rfidInputRef.value && document.activeElement !== rfidInputRef.value) {
        rfidInputRef.value.focus()
    }
}

const clearMessages = debounce(() => {
    currentAttendance.value = null
    errorMessage.value = ''
}, 5000)

const handleRFIDInput = async (uid) => {
    // Validate UID length
    if (!uid || uid.length < 10) {
        return
    }

    // Prevent duplicate requests for same UID
    if (pendingRequests.has(uid)) {
        return
    }

    // Clear previous messages
    currentAttendance.value = null
    errorMessage.value = ''

    // Create abort controller for request cancellation
    const controller = new AbortController()
    pendingRequests.set(uid, controller)

    try {
        const response = await fetch(`${apiBaseUrl}/student/${uid}`, {
            signal: controller.signal,
            headers: {
                'Content-Type': 'application/json',
            }
        })

        if (!response.ok) {
            throw new Error(`Server responded with ${response.status}`)
        }

        const data = await response.json()

        currentAttendance.value = {
            full_name: data.full_name,
            className: data.class?.name || '',
            time: new Date().toLocaleTimeString('id-ID')
        }

        // Auto-clear success message after 5 seconds
        clearMessages()

    } catch (error) {
        if (error.name === 'AbortError') {
            return // Request was cancelled, ignore
        }

        console.error('RFID request failed:', error)
        // Keep original behavior - no error display to user

    } finally {
        pendingRequests.delete(uid)

        // Clear input
        if (rfidInputRef.value) {
            rfidInputRef.value.value = ''
        }

        // Ensure focus after processing
        nextTick(() => {
            ensureFocus()
        })
    }
}

// RFID Scanner handling with buffer
let rfidBuffer = ''
let bufferTimeout = null

const handleScannerKey = (e) => {
    // Clear existing buffer timeout
    if (bufferTimeout) {
        clearTimeout(bufferTimeout)
        activeTimeouts.delete(bufferTimeout)
    }

    if (e.key === 'Enter') {
        // Process complete UID
        const uid = rfidBuffer.trim()
        rfidBuffer = ''

        if (uid.length >= 10) {
            handleRFIDInput(uid)
        }
    } else if (e.key.length === 1) {
        // Accept any single character (maintain original behavior)
        rfidBuffer += e.key

        // Set buffer timeout to clear incomplete scans
        bufferTimeout = setTimeout(() => {
            rfidBuffer = ''
            activeTimeouts.delete(bufferTimeout)
        }, 1000)
        activeTimeouts.add(bufferTimeout)
    }
}

const handleBlur = () => {
    // Ensure input stays focused for scanner
    setTimeout(() => {
        ensureFocus()
    }, 10)
}

const handleVisibilityChange = () => {
    if (!document.hidden) {
        // Page became visible, ensure focus and update time
        ensureFocus()
        now.value = new Date()
    }
}

// Lifecycle hooks
onMounted(async () => {
    // Initial datetime update
    updateDateTime()

    // Start interval for time updates
    intervalId.value = setInterval(updateDateTime, 1000)

    // Wait for DOM to be ready
    await nextTick()

    // Setup RFID input focus management
    if (rfidInputRef.value) {
        ensureFocus()
        rfidInputRef.value.addEventListener('blur', handleBlur)
    }

    // Global event listeners
    document.addEventListener('keydown', handleScannerKey)
    document.addEventListener('visibilitychange', handleVisibilityChange)
    window.addEventListener('focus', ensureFocus)
})

onBeforeUnmount(() => {
    // Clear interval
    if (intervalId.value) {
        clearInterval(intervalId.value)
        intervalId.value = null
    }

    // Clear all active timeouts
    activeTimeouts.forEach(timeoutId => {
        clearTimeout(timeoutId)
    })
    activeTimeouts.clear()

    // Abort all pending requests
    pendingRequests.forEach(controller => {
        controller.abort()
    })
    pendingRequests.clear()

    // Remove event listeners
    if (rfidInputRef.value) {
        rfidInputRef.value.removeEventListener('blur', handleBlur)
    }

    document.removeEventListener('keydown', handleScannerKey)
    document.removeEventListener('visibilitychange', handleVisibilityChange)
    window.removeEventListener('focus', ensureFocus)
})
</script>