<script setup lang="ts">
import ListController from '@/components/ListController.vue'
import { ref } from 'vue'
import { useLists } from '../composables/useLists'

const { lists } = useLists()
const selectedListID = ref('')
</script>

<template>
  <!-- Side navigation contains header and links to all items -->
  <div>
    <div class="side-nav">
      <h3>rw youtube</h3>
      <div class="side-nav-flex border-top">
        <a
          v-for="list in lists"
          :key="list.id"
          :class="{ 'side-nav-selected': list.id == selectedListID }"
          @click="selectedListID = list.id"
        >
          {{ list.name }}
        </a>
      </div>
    </div>

    <!-- Main content located aside of navigation -->
    <div class="main-content">
      <ListController :list-id="selectedListID" />
    </div>
  </div>
</template>

<style scoped>
.side-nav {
  background: #f8f8f8;
  border-right: 1px solid #e7e7e7;
  height: 100vh;
  overflow: auto;
  position: fixed;
  top: 0;
  width: 200px;
  z-index: 1;
}

.side-nav > h3 {
  font-weight: 1200;
}

.side-nav-flex {
  display: flex;
  flex-direction: column;
  padding: 0;
  margin: 0;
}

.side-nav-flex > a {
  text-decoration: none;
  width: 100%;
  padding: 10px;
}

.side-nav-flex > a:hover,
.side-nav-selected {
  --tw-bg-opacity: 1;
  background-color: rgba(209, 213, 219, var(--tw-bg-opacity));
}

.border-top {
  border-top: 1px solid #e7e7e7;
}

.main-content {
  margin-left: 200px;
  padding: 20px;
}
</style>
