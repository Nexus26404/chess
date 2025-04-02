<script lang="ts" setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { http } from '../components/http/http';

const router = useRouter();
const username = ref('');
const password = ref('');
const errorMessage = ref('');

async function handleLogin(event: Event): Promise<void> {
    event.preventDefault();
    try {
        await http.post('/authenticate', {
            username: username.value,
            password: password.value
        });
        router.push('/home');
    } catch (error: unknown) {
        if (error instanceof Error) {
            errorMessage.value = error.message;
        } else {
            errorMessage.value = 'Login failed';
        }
    }
}
</script>

<template>
    <div class="login-container">
        <form @submit="handleLogin">
            <h1>Welcome Back</h1>
            <label for="username">Username:</label>
            <input type="text" id="username" v-model="username" placeholder="Enter your username" required>
            <label for="password">Password:</label>
            <input type="password" id="password" v-model="password" placeholder="Enter your password" required>
            <button type="submit" class="btn">Login</button>
            <p class="error-message">{{ errorMessage }}</p>
            <p>Don't have an account? <a href="/register">Sign up</a></p>
        </form>
    </div>
</template>

<style>
.error-message {
    color: #e74c3c;
    margin-top: 15px;
}
</style>
