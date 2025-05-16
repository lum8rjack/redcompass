<template>
  <div class="mb-8">
    <h1 class="text-2xl font-bold text-white mb-2">{{ title }}</h1>
    <p class="text-gray-400">{{ description }}</p>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true
  },
  description: {
    type: String,
    required: true  
  }
})
</script>