<template>
  <form @submit.prevent="add" class="row g-2">
    <div class="col-9">
      <input v-model="text" type="text"
             class="form-control"
             placeholder="Новый вариант (макс. 40 символов)"
             maxlength="40" required>
    </div>
    <div class="col-3">
      <button class="btn btn-success w-100"
              :disabled="loading">Добавить</button>
    </div>
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