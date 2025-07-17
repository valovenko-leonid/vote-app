<template>
  <ul class="list-group">
    <li v-for="o in options" :key="o.id"
        class="list-group-item d-flex justify-content-between align-items-center"
        :class="{ 'active': hasVoted(o.id) }">
      
      <div class="d-flex align-items-center gap-3">
        <span class="badge bg-info text-dark fs-6">{{ o.votes }}</span>
        <span class="fs-5">{{ o.text }}</span>
      </div>

      <div class="btn-group">
        <button class="btn btn-sm"
                :class="hasVoted(o.id) ? 'btn-danger' : 'btn-outline-primary'"
                :disabled="loading"
                @click="toggleVote(o.id)">
          {{ hasVoted(o.id) ? 'Отменить' : 'Голос' }}
        </button>

        <button v-if="admin" class="btn btn-sm btn-outline-danger"
                @click="deleteOption(o.id)">
          🗑
        </button>
      </div>
    </li>
  </ul>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  options: Array,
  userId: String,
  admin: Boolean,
})
const emit = defineEmits(['voted', 'deleted'])

const voted = ref(new Set())
const loading = ref(false)

const fetchMyVotes = async () => {
  const res = await fetch(`/myvotes?user_id=${props.userId}`)
  const ids = await res.json()
  voted.value = new Set(ids)

  console.log(voted);
  
}

const hasVoted = (id) => voted.value.has(id)

const toggleVote = async (id) => {
  loading.value = true

  
  const res = await fetch('/vote', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ option_id: id, user_id: props.userId })
  })
  loading.value = false

  if (!res.ok) alert(await res.text())
  else {
    emit('voted')
    await fetchMyVotes()
  }
}

const deleteOption = async (id) => {
  if (!confirm("Удалить вариант?")) return
  await fetch(`/option?id=${id}`, { method: 'DELETE' })
  emit('deleted')
}
watch(() => props.userId, (new_userId) => {
  if (new_userId) fetchMyVotes()
}, { immediate: true })
</script>