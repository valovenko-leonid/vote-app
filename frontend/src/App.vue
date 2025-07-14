<template>
  <main class="p-4 max-w-xl mx-auto">
    <h1 class="text-2xl mb-4">Голосование</h1>
    <VoteList :options="options" :fp="fp" @voted="fetchOptions" />
    <AddOption @added="fetchOptions" />
  </main>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Fingerprint2 from 'fingerprintjs2'
import VoteList from './components/VoteList.vue'
import AddOption from './components/AddOption.vue'

const options = ref([])
const fp = ref('')

const fetchOptions = () => {
  fetch('/options').then(r => r.json()).then(data => options.value = data)
}

onMounted(() => {
  Fingerprint2.get(components => {
    fp.value = Fingerprint2.x64hash128(components.map(c => c.value).join(''), 31)
  })
  fetchOptions()

  const ws = new WebSocket(`ws://${location.host}/ws`)
  ws.onmessage = e => { options.value = JSON.parse(e.data) }
})
</script>