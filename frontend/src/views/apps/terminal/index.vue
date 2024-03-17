<script setup lang='ts'>
import { onMounted } from 'vue';
import { Terminal } from '@xterm/xterm';

onMounted(() => {
  const term = new Terminal();
        term.open(document.getElementById('terminal') as HTMLElement);
        term.write('Hello from \x1B[1;3;31mxterm.js\x1B[0m $ ')
})
</script>

<template>
  <div class="terminal">
    <h1 class="terminal-title"> Terminal App </h1>
    <div id="terminal" class="terminal-window bg-black"></div>
  </div>
</template>

<style lang='scss' scoped>
.terminal {
  display: grid;
  grid-template-columns: (10%, 90%);

  .terminal-title {
    @apply p-5;
    @apply text-blue-500;
  }

  .terminal-window {
    @apply p-5;
  }
}
</style>
