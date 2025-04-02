<script lang="ts" setup>
import { defineEmits, defineProps } from 'vue';

defineProps<{
  board: number[][];
  isGameStarted: boolean;
}>();

const emit = defineEmits<{
  (e: 'cellClick', row: number, col: number): void;
}>();

function handleCellClick(row: number, col: number) {
  emit('cellClick', row, col);
}
</script>

<template>
  <div v-if="isGameStarted" class="game-board">
    <div v-for="(row, rowIndex) in board" :key="rowIndex" class="board-row">
      <div 
        v-for="(cell, colIndex) in row" 
        :key="colIndex" 
        class="board-cell"
        @click="handleCellClick(rowIndex, colIndex)"
      >
        <div v-if="cell === 1" class="piece black"></div>
        <div v-if="cell === 2" class="piece white"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.game-board {
  justify-self: center;
  display: inline-grid;
  gap: 1px;
  padding: 2rem;
  background: url('data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEASABIAAD/4gHYSUNDX1BST0ZJTEUAAQEAAAHIAAAAAAQwAABtbnRyUkdCIFhZWiAH4AABAAEAAAAAAABhY3NwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAA9tYAAQAAAADTLQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAlkZXNjAAAA8AAAACRyWFlaAAABFAAAABRnWFlaAAABKAAAABRiWFlaAAABPAAAABR3dHB0AAABUAAAABRyVFJDAAABZAAAAChnVFJDAAABZAAAAChiVFJDAAABZAAAAChjcHJ0AAABjAAAADxtbHVjAAAAAAAAAAEAAAAMZW5VUwAAAAgAAAAcAHMAUgBHAEJYWVogAAAAAAAAb6IAADj1AAADkFhZWiAAAAAAAABimQAAt4UAABjaWFlaIAAAAAAAACSgAAAPhAAAts9YWVogAAAAAAAA9tYAAQAAAADTLXBhcmEAAAAAAAQAAAACZmYAAPKnAAANWQAAE9AAAApbAAAAAAAAAABtbHVjAAAAAAAAAAEAAAAMZW5VUwAAACAAAAAcAEcAbwBvAGcAbABlACAASQBuAGMALgAgADIAMAAxADb/2wBDABQODxIPDRQSEBIXFRQdHx4dHRsdHR4dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR3/2wBDAR0XFyAeIRshIRshHRsdIR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR0dHR3/wAARCAAIAAgDAREAAhEBAxEB/8QAFQABAQAAAAAAAAAAAAAAAAAAAAb/xAAUEAEAAAAAAAAAAAAAAAAAAAAA/8QAFQEBAQAAAAAAAAAAAAAAAAAAAAX/xAAUEQEAAAAAAAAAAAAAAAAAAAAA/9oADAMBAAIRAxEAPwCdABmX/9k=');
  padding: 1.5rem;
  border-radius: 12px;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15);
  display: inline-block;
  margin: 1.5rem 0;
}

.board-row {
  display: grid;
  grid-auto-flow: column;
  gap: 1px;
}

.board-cell {
  width: 42px;
  height: 42px;
  display: grid;
  place-items: center;
  border: 1px solid rgba(0, 0, 0, 0.2);
  cursor: pointer;
  position: relative;
  transition: all 0.2s;
}

.board-cell:hover {
  background: rgba(255, 255, 255, 0.1);
  box-shadow: inset 0 0 10px rgba(255, 255, 255, 0.5);
}

.piece {
  width: 38px;
  height: 38px;
  border-radius: 50%;
  transition: all 0.3s;
}

.piece.black {
  background: radial-gradient(circle at 30% 30%, #4a4a4a, #000);
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.3);
}

.piece.white {
  background: radial-gradient(circle at 30% 30%, #fff, #e0e0e0);
  border: 2px solid #2c3e50;
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.2);
}
</style>
