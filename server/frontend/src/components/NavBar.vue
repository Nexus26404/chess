<script setup lang="ts">
</script>

<template>
  <nav>
    <div class="nav-container">
      <div class="logo">Chess Game</div>
      <div class="nav-links">
        <router-link to="/">Home</router-link>
        <router-link to="/lobby">Lobby</router-link>
        <router-link to="/user">User</router-link>
      </div>
    </div>
  </nav>
</template>

<style scoped>
nav {
  background: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 15px 0;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
}

.nav-container {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
}

.logo {
  font-size: 1.5rem;
  font-weight: bold;
  color: #2ecc71;
}

.nav-links {
  display: flex;
  gap: 20px;
}

nav a {
  font-weight: 600;
  color: #2c3e50;
  text-decoration: none;
  padding: 8px 16px;
  border-radius: 6px;
  transition: all 0.3s ease;
}

nav a:hover {
  background-color: #f0f2f5;
  color: #2ecc71;
}

nav a.router-link-exact-active {
  color: white;
  background-color: #2ecc71;
}
</style>
