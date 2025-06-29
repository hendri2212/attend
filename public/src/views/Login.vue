<template>
    <div class="container d-flex justify-content-center align-items-center vh-100">
        <div class="card border-0 bg-body-tertiary rounded-4" style="width: 400px;">
            <div class="card-body">
                <div class="text-center mb-4">
                    <img src="/logo.png" alt="Logo" class="img-fluid" style="width: 120px; height: auto;">
                </div>
                <form @submit.prevent="handleLogin">
                    <div class="form-floating mb-3">
                        <input
                            type="text"
                            class="form-control"
                            id="floatingEmail"
                            placeholder="Email"
                            v-model="email"
                            required>
                        <label for="floatingEmail">Email</label>
                    </div>
                    <div class="form-floating mb-3">
                        <input
                            type="password"
                            class="form-control"
                            id="floatingPassword"
                            placeholder="Password"
                            v-model="password"
                            required>
                        <label for="floatingPassword">Password</label>
                    </div>
                    <button type="submit" class="btn btn-primary w-100">Login</button>
                </form>
                <div class="text-center mt-3">
                    <p class="mb-0">Belum punya akun? <RouterLink to="/register" class="text-decoration-none">Daftar
                        </RouterLink>
                    </p>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { apiBaseUrl } from '@/config'

// Reactive form fields
const email = ref('')
const password = ref('')

// Router instance
const router = useRouter()

// Handle login form submit
const handleLogin = async () => {
    try {
        // Axios automatically sets headers and parses JSON
        const { data } = await axios.post(`${apiBaseUrl}/login`, {
            email: email.value,
            password: password.value,
        })

        if (data.token) {
            localStorage.setItem('token', data.token)
            router.push('/admin')
        } else {
            throw new Error('Token not found in response')
        }
    } catch (err) {
        const message = err.response?.data?.message || err.message || 'Login failed'
        alert(message)
    }
}
</script>