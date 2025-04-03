<script lang="ts" setup>
import { computed, onBeforeMount, onBeforeUnmount, ref } from 'vue';
import { useRouter } from 'vue-router';
import { http } from '../../components/http/http';
import Board from './Board.vue';
import Room from './Room.vue';
import { PlayerTurn, RoomStatus, type UserInfo } from './types';

const router = useRouter();
const route = useRouter().currentRoute.value;
const boardSize = 15;
const board = ref<number[][]>(Array(boardSize).fill(0).map(() => Array(boardSize).fill(0)));
const currentPlayer = ref<number>(1);
const winner = ref<number>(0);
const roomId = ref<string>(route.params.roomId as string);
const ws = ref<WebSocket | null>(null);
const roomStatus = ref<RoomStatus>(RoomStatus.Unready);
const playerTurn = ref<PlayerTurn>(PlayerTurn.None)
const owner = ref<string>('')
const userInfo = ref<UserInfo>()
const isGameStarted = computed<boolean>(() => {
  return roomStatus.value === RoomStatus.Started;
});

const gameStatus = computed(() => {
  if (!isGameStarted.value) {
    return 'Waiting for players';
  }
  if (winner.value > 0) {
    return winner.value === 1 ? 'Black wins!' : 'White wins!';
  }
  return playerTurn.value === PlayerTurn.Black ? 'Black\'s turn' : 'White\'s turn';
});

const isOwner = computed(() => {
  return owner.value === userInfo?.value?.id || false;
});

async function fetchUserInfo() {
  try {
    const response = await http.get('/user/info');
    userInfo.value = response;
  } catch (err) {
    if (err instanceof Error) {
      if (err.message === 'Unauthorized' || err.message === 'Session expired') {
        router.push('/login');
      }
    }
  }
}

onBeforeMount(async () => {
  await fetchUserInfo();
  if (!roomId.value) {
    router.push('/user');
    return;
  }
  initializeWebSocket();
});

function initializeWebSocket() {
  ws.value = new WebSocket(`wss://${window.location.host}/room/ws/${roomId.value}?userId=${userInfo.value?.id}`);

  ws.value.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log(data);
    
    if (data.type === 'update') {
      board.value = data.board;
      currentPlayer.value = data.player;
      winner.value = data.winner || 0;
      playerTurn.value = data.playerTurn;
    } else if (data.type === 'room_update') {
      owner.value = data.owner;
      roomStatus.value = data.roomStatus;
      playerTurn.value = data.playerTurn;

      if (isGameStarted.value) {
        currentPlayer.value = data.playerTurn;
      }
    }
  };

  ws.value.onopen = () => {
    ws.value?.send(JSON.stringify({ type: 'get_state' }));
  };
}

function setReady() {
  if (ws.value?.readyState === WebSocket.OPEN) {
    ws.value.send(JSON.stringify({
      type: 'ready',
      gameId: roomId.value
    }));
  }
}

function cancelReady() {
  if (ws.value?.readyState === WebSocket.OPEN) {
    ws.value.send(JSON.stringify({
      type: 'cancel_ready',
      gameId: roomId.value
    }));
  }
}

function startGame() {
  if (ws.value?.readyState === WebSocket.OPEN) {
    ws.value.send(JSON.stringify({
      type: 'start_game'
    }));
  }
}

onBeforeUnmount(() => {
  if (ws.value) {
    ws.value.close();
  }
});

async function handleCellClick(row: number, col: number) {
  if (roomStatus.value !== RoomStatus.Started) return;
  if (board.value[row][col] === 0 && ws.value?.readyState === WebSocket.OPEN) {
    ws.value.send(JSON.stringify({
      type: 'move',
      gameId: roomId.value,
      row: row,
      col: col
    }));
  }
}

async function resetGame() {
  if (ws.value?.readyState === WebSocket.OPEN) {
    ws.value.send(JSON.stringify({
      type: 'reset',
      gameId: roomId.value
    }));
  }
}
</script>

<template>
  <div class="container">
    <div class="game-content">
      <h1 class="title">Five in a Row</h1>

      <Room v-if="roomStatus !== RoomStatus.Started" :room-id="roomId" :player="userInfo?.nickname" :is-owner="isOwner"
        :room-status="roomStatus" @ready="setReady" @cancel-ready="cancelReady" @start="startGame" />

      <Board v-else :board="board" :is-game-started="isGameStarted" @cell-click="handleCellClick" />
      <div class="game-status">{{ gameStatus }}</div>

      <div v-if="isOwner" class="game-controls">
        <button class="btn btn-primary" @click="resetGame">New Game</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 2rem;
}

.game-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  max-width: 1000px;
  width: 95%;
  margin: 0 auto;
  gap: 2rem;
}

.title {
  color: #2c3e50;
  font-size: 2.5rem;
  margin-bottom: 1.5rem;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1);
}

.game-status {
  font-size: 1.4rem;
  color: #2c3e50;
  margin-bottom: 1.5rem;
  font-weight: 600;
}

.game-controls {
  display: flex;
  gap: 1.5rem;
  justify-content: center;
  align-items: center;
  width: 100%;
  padding: 1rem 0;
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

.room-info .btn {
  margin: 1rem auto;
  width: 100%;
  max-width: 300px;
}

.btn:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.2);
}

.btn:active {
  transform: translateY(-1px);
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

@media (max-width: 768px) {
  .game-content {
    grid-gap: 1rem;
  }

  .board-cell {
    width: 36px;
    height: 36px;
  }

  .piece {
    width: 32px;
    height: 32px;
  }

  .btn {
    min-width: 140px;
    height: 44px;
  }
}

@media (max-width: 480px) {
  .board-cell {
    width: 28px;
    height: 28px;
  }

  .piece {
    width: 24px;
    height: 24px;
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-30px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
