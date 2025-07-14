<template>
  <form @submit.prevent="add">
    <input v-model="text" placeholder="Новый вариант"
           class="border p-1 mr-2" maxlength="40">
    <button class="bg-green-600 text-white px-2 rounded"
            :disabled="loading">Добавить</button>
  </form>
</template>

<script setup>
import { ref } from 'vue'
const emit = defineEmits(['added'])
const text = ref('')
const loading = ref(false)

const add = async () => {
  loading.value = true
  const res = await fetch('/option', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ text: text.value })
  })
  loading.value = false
  if (res.ok) {
    text.value = ''
    emit('added')
  } else alert(await res.text())
}
</script>