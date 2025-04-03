<script lang="ts" setup>
import { useRouter } from 'vue-router';

const router = useRouter();

const navigateToLobby = () => {
  router.push('/lobby');
};
</script>

<template>
  <div class="container">
    <div class="home-content">
      <h1 class="title">Welcome to Chess Online</h1>
      <div class="welcome-content">
        <p class="subtitle">Experience the classic game of chess in a modern way</p>
        <div class="features">
          <h2>Features:</h2>
          <ul class="feature-list">
            <li><span class="chess-icon">♟</span> Play against other players online</li>
            <li><span class="chess-icon">♜</span> Practice with AI opponent</li>
            <li><span class="chess-icon">♝</span> Track your progress</li>
            <li><span class="chess-icon">♛</span> Join tournaments</li>
          </ul>
        </div>
        <div class="cta">
          <button class="btn btn-primary" @click="navigateToLobby">Start Playing</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.home-content {
  text-align: center;
  padding: 2rem;
}

.title {
  color: #2ecc71;
  margin-bottom: 2rem;
  animation: fadeIn 1s ease-in;
}

.subtitle {
  color: #7f8c8d;
  font-size: 1.2rem;
  margin-bottom: 2rem;
}

.welcome-content {
  background-color: #f8f9fa;
  padding: 2rem;
  border-radius: 12px;
  margin-top: 2rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.features {
  margin: 2rem 0;
}

.feature-list {
  list-style: none;
  padding: 0;
  max-width: 400px;
  margin: 0 auto;
}

.feature-list li {
  margin: 1rem 0;
  font-size: 1.1rem;
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 1rem;
  text-align: left;
}

.chess-icon {
  color: #2ecc71;
  font-size: 1.5rem;
}

.cta {
  margin-top: 2rem;
}

.btn-primary {
  background-color: #2ecc71;
  color: white;
  border: none;
  padding: 0.8rem 2rem;
  font-size: 1.2rem;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.2s, background-color 0.2s;
}

.btn-primary:hover {
  background-color: #27ae60;
  transform: translateY(-2px);
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
