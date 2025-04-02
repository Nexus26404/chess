<script lang="ts" setup>
import { useRouter } from 'vue-router';

const router = useRouter();

const goHome = () => {
    router.push('/');
};
</script>

<template>
    <div class="container">
        <div class="not-found-content">
            <h1>404 - Page Not Found</h1>
            <div class="chess-piece">♟</div>
            <p>Oops! This move led to a dead end!</p>
            <p class="subtitle">The page you're looking for has been captured or never existed.</p>
            <div class="cta">
                <button class="btn btn-success" @click="goHome">Return to Home</button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.not-found-content {
    text-align: center;
    padding: 2rem;
}

.chess-piece {
    font-size: 6rem;
    color: #2ecc71;
    margin: 2rem 0;
    animation: bounce 2s infinite;
}

.subtitle {
    color: #7f8c8d;
    margin-bottom: 2rem;
}

@keyframes bounce {

    0%,
    100% {
        transform: translateY(0);
    }

    50% {
        transform: translateY(-20px);
    }
}
</style>
