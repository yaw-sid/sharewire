<template>
  <article>
    <div class="logo">
      <div class="image"></div>
    </div>
    <h1>Please select a user to interact with</h1>
    <transition-group name="list" tag="ul">
      <li v-for="peer in state.peers" :key="peer.id" @click="navigateToFiles(peer)">{{ peer.name }}</li>
    </transition-group>
  </article>
</template>

<script setup>
import { onMounted, reactive } from "vue";
import { useRouter } from "vue-router";
import { ListPeers } from "../../wailsjs/go/main/App";

const router = useRouter();

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

function navigateToFiles(peer) {
  router.push({
    name: "Files",
    query: {
      "peer": peer.name
    }
  });
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
  margin-bottom: 0.5rem;
  border: 1px solid #aaa;
  border-radius: 0.75rem;
}

li:hover {
  cursor: pointer;
  background-color: #efee;
}

.list-enter-active,
.list-leave-active {
  transition: all 0.5s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
</style>