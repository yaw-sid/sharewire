<template>
  <article>
    <header>
      <button @click="goBack">&#10092;</button>
      <h1>{{ peer }}</h1>
    </header>

    <section>
      <ul>
        <li v-for="(file, index) in files" :key="index">
          <p>{{ file.name }}</p>
          <button>&darr;</button>
        </li>
      </ul>
    </section>
  </article>
</template>

<script setup>
import { computed, reactive } from "vue";
import { useRoute, useRouter } from "vue-router";

const route = useRoute();
const router = useRouter();

const peer = computed(() => route.query.peer);

const files = reactive([
  {
    name: "File 1"
  }, {
    name: "File 2"
  }, {
    name: "File 3"
  }
]);

function goBack() {
  router.back();
}
</script>

<style scoped>
article {
  max-width: 90%;
  margin: auto;
}

header {
  display: flex;
  align-items: center;
  gap: 2rem;
}

ul {
  padding-inline-start: 0;
  list-style-type: none;
}

li {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  transition: ease background-color 250ms;
  margin-bottom: 0.5rem;
  border: 1px solid #aaa;
  border-radius: 0.75rem;
}

li:hover {
  background-color: #efee;
}

li p {
  margin: 0;
}

button:hover {
  cursor: pointer;
}
</style>