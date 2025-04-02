<script lang="ts" setup>
import { defineEmits, defineProps } from 'vue';
import { RoomStatus } from './types';

defineProps<{
  roomId: string;
  isOwner: boolean;
  roomStatus: RoomStatus;
  player: string | undefined;
}>();

defineEmits<{
  (e: 'ready'): void;
  (e: 'start'): void;
  (e: 'cancelReady'): void;
}>();
</script>

<template>
  <div class="room-info">
    <div>Room ID: {{ roomId }}</div>
    <div class="players-list">
      <div>Player: {{ player }}</div>
    </div>
    <template v-if="isOwner">
      <div v-if="roomStatus === RoomStatus.Unready" class="status-message">
        Waiting for all players to be ready...
      </div>
      <button 
        v-if="roomStatus === RoomStatus.Ready"
        class="btn btn-primary"
        @click="$emit('start')"
      >
        Start Game
      </button>
    </template>
    <template v-else>
      <button 
        v-if="roomStatus === RoomStatus.Unready"
        class="btn btn-primary"
        @click="$emit('ready')"
      >
        Ready
      </button>
      <button 
        v-if="roomStatus === RoomStatus.Ready"
        class="btn btn-secondary"
        @click="$emit('cancelReady')"
      >
        Cancel Ready
      </button>
    </template>
  </div>
</template>

<style scoped>
.room-info {
  display: grid;
  grid-gap: 1rem;
  max-width: 600px;
  margin: 0 auto;
  width: 100%;
  background: linear-gradient(145deg, #f8f9fa, #e9ecef);
  padding: 1.5rem;
  border-radius: 12px;
  margin: 1.5rem 0;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.players-list {
  display: grid;
  grid-gap: 0.5rem;
  padding: 1rem;
  max-height: 200px;
  overflow-y: auto;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 8px;
}

.player {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 1rem;
  background: #fff;
  border-radius: 6px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  transition: transform 0.2s;
}

.player:hover {
  transform: translateX(5px);
}

.btn {
  min-width: 180px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 48px;
  padding: 0.8rem 2rem;
  font-size: 1.1rem;
  border-radius: 50px;
  cursor: pointer;
  transition: all 0.3s;
  border: none;
  text-transform: uppercase;
  font-weight: 600;
  letter-spacing: 1px;
  margin: 0 auto;
}

.btn-primary {
  background: linear-gradient(45deg, #2ecc71, #27ae60);
  color: white;
}

.btn-secondary {
  background: linear-gradient(45deg, #95a5a6, #7f8c8d);
  color: white;
}

.status-message {
  color: #666;
  margin: 1rem 0;
  font-style: italic;
  font-size: 1.1rem;
}
</style>
