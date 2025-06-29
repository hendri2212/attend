<template>
    <div class="card border-0 bg-body-tertiary mb-3">
        <div class="card-body">
            <span class="fs-4 fw-bold mb-0 card-title">{{ currentTitle }}</span>
        </div>
    </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const currentTitle = computed(() => route.meta.title || 'Dashboard')
</script>