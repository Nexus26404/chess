<script setup lang="ts">
import { http } from '@/components/http/http'
import { onBeforeMount, ref } from 'vue'
import { useRouter } from 'vue-router'
import type { Room } from '../types/room'

const router = useRouter()
const rooms = ref<Room[]>([])
const newRoomName = ref('')
const errorMessage = ref('')
const userId = ref('')

const createRoom = async () => {
  if (!newRoomName.value.trim()) {
    errorMessage.value = 'Room name cannot be empty'
    return
  }
  errorMessage.value = ''
  try {
    const response = await http.get(`/room/create?roomName=${encodeURIComponent(newRoomName.value)}&userId=${userId.value}`)
    if (response.error) {
      errorMessage.value = response.error
    } else {
      newRoomName.value = ''
      fetchRooms()
      router.push({
        path: `/game/${response.roomId}`,
        query: { player: response.player }
      })
    }
  } catch (error) {
    console.error('Failed to create room:', error)
  }
}

const joinRoom = async (roomId: string) => {
  try {
    const response = await fetch(`/api/rooms/${roomId}/join`, {
      method: 'POST'
    })
    if (response.ok) {
      router.push(`/game/${roomId}`)
    }
  } catch (error) {
    console.error('Failed to join room:', error)
  }
}

const fetchRooms = async () => {
  try {
    const response = await fetch('/room/list')
    if (response.status === 401) {
      router.push('/login')
      return
    }
    if (response.ok) {
      const data = await response.json()
      rooms.value = data.rooms
    }
  } catch (error) {
    console.error('Failed to fetch rooms:', error)
  }
}

const fetchUserInfo = async () => {
  try {
    const response = await fetch('/user/info')
    if (response.ok) {
      const data = await response.json()
      userId.value = data.id
    }
  } catch (error) {
    console.error('Failed to fetch user info:', error)
  }
}

onBeforeMount(async () => {
  await Promise.all([fetchUserInfo(), fetchRooms()])
})
</script>

<template>
  <div class="container">
    <div class="lobby">
      <h1>Chess Game Lobby</h1>

      <!-- Create Room Form -->
      <div class="create-room">
        <div class="input-group">
          <input v-model="newRoomName" placeholder="Enter room name" />
          <button class="btn btn-success" @click="createRoom">Create Room</button>
        </div>
        <p class="error-message" v-if="errorMessage">{{ errorMessage }}</p>
      </div>

      <!-- Room List -->
      <div class="room-list">
        <h2>Available Rooms</h2>
        <div v-if="rooms.length === 0" class="no-rooms">
          No rooms available
        </div>
        <div v-else class="rooms">
          <div v-for="room in rooms" :key="room.id" class="room-item">
            <div class="room-info">
              <span class="room-name">{{ room.name }}</span>
              <span class="player-count">Players: {{ room.isFull ? 2 : 1 }}/2</span>
            </div>
            <button class="btn" @click="joinRoom(room.id)" :disabled="room.isFull">
              Join
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.lobby {
  text-align: center;
  padding: 2rem;
  max-width: 800px;
  margin: 0 auto;
}

.create-room {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  padding: 1rem;
  margin: 2rem auto;
  max-width: 600px;
  background-color: #f8f9fa;
  border-radius: 12px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.input-group {
  display: flex;
  justify-content: center;
  align-items: stretch;
  gap: 1rem;
}

.create-room input {
  margin: 0;
  height: 40px;
  width: 300px;
  padding: 0.5rem 1rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 1rem;
  box-sizing: border-box;
}

.create-room button {
  margin: 0;
  height: 40px;
  padding: 0.5rem 1.5rem;
  font-size: 1rem;
  border-radius: 6px;
  min-width: 120px;
  box-sizing: border-box;
}

.room-list {
  margin-top: 2rem;
  background-color: #f8f9fa;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.rooms {
  max-width: 600px;
  margin: 0 auto;
}

.room-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border: 1px solid #ddd;
  margin-bottom: 1rem;
  border-radius: 8px;
  background-color: white;
  transition: transform 0.2s, box-shadow 0.2s;
}

.room-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.room-info {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.5rem;
}

.room-name {
  font-size: 1.1rem;
  color: #2c3e50;
  font-weight: bold;
}

.player-count {
  color: #7f8c8d;
  font-size: 0.9rem;
}

.no-rooms {
  text-align: center;
  color: #7f8c8d;
  padding: 2rem;
  font-size: 1.1rem;
}

button:disabled {
  background-color: #bdc3c7;
  cursor: not-allowed;
}

h1 {
  color: #2c3e50;
  margin-bottom: 1.5rem;
}

h2 {
  color: #2c3e50;
  margin-bottom: 1.5rem;
  font-size: 1.5rem;
}

.error-message {
  color: #e74c3c;
  margin: 0;
  font-size: 0.9rem;
  text-align: center;
}
</style>
