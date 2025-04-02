<script lang="ts" setup>
import { onBeforeMount, ref } from 'vue';
import { useRouter } from 'vue-router';
import { http } from '../components/http/http';
import RoomDialog from '../components/RoomDialog.vue';
import type { UserInfo } from './game/types';

interface RoomSelection {
  roomId: string;
  player: number;
}

const router = useRouter();
const user = ref<UserInfo | null>(null);
const error = ref('');
const showDialog = ref(false);

async function fetchUserInfo() {
  try {
    const response = await http.get('/user/info');
    user.value = response;
  } catch (err) {
    if (err instanceof Error) {
      error.value = err.message;
      if (error.value === 'Unauthorized' || error.value === 'Session expired') {
        router.push('/login');
      }
    } else {
      error.value = 'Failed to load user information';
    }
  }
}

const playNewGame = () => {
  showDialog.value = true;
}

const handleRoomSelection = (room: RoomSelection) => {
  showDialog.value = false;
  router.push({
    path: `/game/${room.roomId}`,
    query: { player: room.player.toString() }
  });
}

onBeforeMount(() => {
  fetchUserInfo();
});
</script>

<template>
  <div class="container">
    <div class="user-profile">
      <h1 class="title">Player Profile</h1>
      
      <div v-if="error" class="error-message">
        {{ error }}
      </div>

      <div v-if="user" class="profile-content">
        <div class="profile-section">
          <div class="profile-icon">♟</div>
          <h2 class="nickname">{{ user.nickname }}</h2>
        </div>
        
        <div class="stats-section">
          <h3>Game Statistics</h3>
          <ul class="feature-list">
            <li><span class="stat-label">Games Played:</span> {{ user.gamesPlayed }}</li>
            <li><span class="stat-label">Wins:</span> {{ user.wins }}</li>
            <li><span class="stat-label">Losses:</span> {{ user.losses }}</li>
            <li><span class="stat-label">Draws:</span> {{ user.draws }}</li>
          </ul>
        </div>

        <div class="actions">
          <button class="btn btn-primary" @click="playNewGame">Play New Game</button>
          <button class="btn btn-secondary">View Game History</button>
        </div>
      </div>

      <div v-else-if="!error" class="loading">
        Loading user information...
      </div>
    </div>
    <RoomDialog 
      :show="showDialog"
      :userId="user?.id.toString() ?? ''"
      @select="handleRoomSelection"
      @close="showDialog = false"
    />
  </div>
</template>

<style scoped>
.user-profile {
  text-align: center;
  padding: 2rem;
  animation: fadeIn 1s ease-in;
}

.title {
  color: #2ecc71;
  margin-bottom: 2rem;
}

.profile-content {
  background-color: #f8f9fa;
  padding: 2rem;
  border-radius: 12px;
  margin-top: 2rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.profile-section {
  margin: 2rem 0;
}

.profile-icon {
  font-size: 4rem;
  color: #2ecc71;
  margin-bottom: 1rem;
  transition: transform 0.3s ease;
}

.profile-icon:hover {
  transform: scale(1.1);
}

.nickname {
  color: #2c3e50;
  font-size: 2rem;
  margin: 0.5rem 0;
}

.stats-section {
  background-color: white;
  padding: 2rem;
  border-radius: 8px;
  margin: 2rem auto;
  max-width: 500px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.stats-section h3 {
  color: #2ecc71;
  margin-bottom: 1.5rem;
}

.feature-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.feature-list li {
  margin: 1rem 0;
  font-size: 1.1rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.feature-list li:hover {
  background-color: #f1f1f1;
}

.stat-label {
  color: #7f8c8d;
  font-weight: 500;
}

.loading {
  color: #7f8c8d;
  margin: 2rem 0;
  font-size: 1.1rem;
}

.actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
  margin-top: 2rem;
}

.btn {
  min-width: 150px;
  padding: 0.8rem 2rem;
  font-size: 1.1rem;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.2s, background-color 0.2s;
  border: none;
}

.btn:hover {
  transform: translateY(-2px);
}

.btn-primary {
  background-color: #2ecc71;
  color: white;
}

.btn-primary:hover {
  background-color: #27ae60;
}

.btn-secondary {
  background-color: #3498db;
  color: white;
}

.btn-secondary:hover {
  background-color: #2980b9;
}

.error-message {
  color: #e74c3c;
  margin: 1rem 0;
  padding: 1rem;
  background-color: #fdf0ed;
  border-radius: 8px;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
