<template>
  <div class="card shadow-sm">
    <div class="card-header bg-primary text-white">
      <h3 class="mb-0">Голосование</h3>
    </div>
    <div class="card-body">
      <VoteList :options="options" :userId="userId" :admin="isAdmin" @voted="fetchOptions" @deleted="fetchOptions" />
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
import { v4 as uuidv4 } from 'uuid'

const options = ref([])
const fp = ref('')
const isAdmin = ref(false)
const fingerprint = ref();
const userId = ref();

const fetchOptions = () => {
  fetch('/options')
    .then(r => r.json())
    .then(data => {
      // сортировка по убыванию голосов
      options.value = data.sort((a, b) => b.votes - a.votes)
    })
}

onMounted(() => {
  let clientId = localStorage.getItem('client_id')
  if (!clientId) {
    clientId = uuidv4()
    localStorage.setItem('client_id', clientId)
  }

  Fingerprint2.get(components => {
    const fpHash = Fingerprint2.x64hash128(components.map(c => c.value).join(''), 31)
    fingerprint.value = fpHash

    // Пытаемся получить userId с сервера
    fetch(`/whoami?fp=${fpHash}`)
      .then(r => {
        if (r.status === 204) return null
        return r.json()
      })
      .then(data => {
        if (data && data.user_id) {
          userId.value = data.user_id
          localStorage.setItem('user_id', userId.value)
        } else {
          // Новый пользователь
          userId.value = uuidv4()
          localStorage.setItem('user_id', userId.value)
          fetch('/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ user_id: userId.value, fp: fpHash })
          })
        }
      })
  })

  fetchOptions()

  const urlParams = new URLSearchParams(window.location.search)
  isAdmin.value = urlParams.get('admin') === '1'

  const ws = new WebSocket(`${location.protocol === 'https:' ? 'wss' : 'ws'}://${location.host}/ws`)
  
  ws.onmessage = e => {
    const data = JSON.parse(e.data)
    options.value = data.sort((a, b) => b.votes - a.votes)
  }
})

</script>