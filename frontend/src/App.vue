<template>
  <div class="card shadow-sm">
    <div class="card-header bg-primary text-white">
      <h3 class="mb-0">Голосование</h3>
    </div>
    <div class="card-body">
      <VoteList :options="options" :fp="fp" :admin="isAdmin" @voted="fetchOptions" @deleted="fetchOptions" />
      <hr class="my-4" />
      <AddOption @added="fetchOptions" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Fingerprint2 from 'fingerprintjs2'
import VoteList from './components/VoteList.vue'
import AddOption from './components/AddOption.vue'

const options = ref([])
const fp = ref('')
const isAdmin = ref(false)

const fetchOptions = () => {
  fetch('/options')
    .then(r => r.json())
    .then(data => {
      // сортировка по убыванию голосов
      options.value = data.sort((a, b) => b.votes - a.votes)
    })
}

onMounted(() => {
  Fingerprint2.get(components => {
    fp.value = Fingerprint2.x64hash128(components.map(c => c.value).join(''), 31)
  })
  fetchOptions()

  const urlParams = new URLSearchParams(window.location.search)
  isAdmin.value = urlParams.get('admin') === '1'

  const ws = new WebSocket(`ws://${location.host}/ws`)
  ws.onmessage = e => {
    const data = JSON.parse(e.data)
    options.value = data.sort((a, b) => b.votes - a.votes)
  }
})
</script>