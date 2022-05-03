<template>
  <article>
    <div class="logo">
      <div class="image"></div>
    </div>
    <h1>Please select a user to interact with</h1>
    <ul>
      <li v-for="peer in state.peers" :key="peer.id">{{ peer.name }}</li>
    </ul>
  </article>
</template>

<script setup>
import { onMounted, reactive } from "vue";
import { ListPeers } from "../../wailsjs/go/main/App";

const state = reactive({
  peers: []
});

onMounted(async () => {
  await listPeers();
});

async function listPeers() {
  try {
    state.peers = await ListPeers();
  } catch (error) {
    LogError(error);
  }
}
</script>

<style scoped>
article {
  max-width: 90%;
  margin: auto;
  padding-top: 4rem;
}

.image {
  width: 10rem;
  height: 10rem;
  background-color: #aaa;
  margin: auto;
}

h1 {
  text-align: center;
}

ul {
  padding-inline-start: 0;
  list-style-type: none;
}

li {
  padding: 1rem;
  text-align: center;
  transition: ease background-color 250ms;
}

li:hover {
  cursor: pointer;
  background-color: #efee;
}
</style>