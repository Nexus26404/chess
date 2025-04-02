<script lang="ts" setup>
import { http } from '@/components/http/http';
import { ref } from 'vue';

interface RoomSelectionData {
  roomId: string;
  player: number;
}

const props = defineProps<{
  userId: string;
  show: boolean;
}>();

const emit = defineEmits<{
  (event: 'close'): void;
  (event: 'select', data: RoomSelectionData): void;
}>();

const showJoinInput = ref(false);
const roomId = ref('');
const error = ref('');

interface RoomResponse {
  roomId: string;
  player: number;
  owner: boolean;
  error?: string;
}

async function createRoom() {
  try {
    const response = await http.get(`/room/create?userId=${props.userId}`) as RoomResponse;
    if (response.error) {
      error.value = response.error;
    } else {
      emit('select', { roomId: response.roomId, player: response.player });
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to create room';
  }
}

async function joinRoom() {
  if (!roomId.value) {
    error.value = 'Please enter a room ID';
    return;
  }
  try {
    const response = await http.get(`/room/join/${roomId.value}?userId=${props.userId}`) as RoomResponse;
    if (response.error) {
      error.value = response.error;
    } else {
      emit('select', { roomId: response.roomId, player: response.player });
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to join room';
  }
}
</script>

<template>
  <div v-if="show" class="dialog-overlay">
    <div class="dialog-content">
      <button class="close-button" @click="emit('close')">×</button>
      <h2 class="dialog-title">Play Game</h2>
      <div v-if="error" class="error-message">{{ error }}</div>
      
      <div v-if="!showJoinInput" class="dialog-actions">
        <button class="dialog-btn primary" @click="createRoom">Create Room</button>
        <button class="dialog-btn secondary" @click="showJoinInput = true">Join Room</button>
      </div>
      
      <div v-else class="join-form">
        <input v-model="roomId" placeholder="Enter Room ID" class="room-input" />
        <div class="dialog-actions">
          <button class="dialog-btn primary" @click="joinRoom">Join</button>
          <button class="dialog-btn secondary" @click="showJoinInput = false">Back</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(2px);
}

.dialog-content {
  position: relative;
  background: white;
  padding: 2.5rem 2rem 2rem;
  border-radius: 12px;
  min-width: 320px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.dialog-title {
  text-align: center;
  margin: 0 0 1.5rem;
  color: #2c3e50;
  font-size: 1.5rem;
}

.dialog-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1.5rem;
  justify-content: center;
}

.dialog-btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  cursor: pointer;
  min-width: 120px;
  transition: all 0.2s ease;
}

.dialog-btn.primary {
  background: #3498db;
  color: white;
}

.dialog-btn.primary:hover {
  background: #2980b9;
}

.dialog-btn.secondary {
  background: #ecf0f1;
  color: #2c3e50;
}

.dialog-btn.secondary:hover {
  background: #bdc3c7;
}

.join-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.room-input {
  padding: 0.75rem;
  border: 2px solid #ecf0f1;
  border-radius: 6px;
  font-size: 1rem;
  width: 100%;
  box-sizing: border-box;
  transition: border-color 0.2s ease;
}

.room-input:focus {
  outline: none;
  border-color: #3498db;
}

.error-message {
  color: #e74c3c;
  margin-bottom: 1rem;
  text-align: center;
  padding: 0.5rem;
  background: #fde8e8;
  border-radius: 4px;
}

.close-button {
  position: absolute;
  top: 12px;
  right: 12px;
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  padding: 0 8px;
  color: #95a5a6;
  transition: color 0.2s ease;
}

.close-button:hover {
  color: #2c3e50;
}
</style>
