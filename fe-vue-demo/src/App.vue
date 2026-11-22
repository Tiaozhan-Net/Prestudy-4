<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

import Sidebar from './components/Sidebar.vue'
import BackHome from './components/BackHome.vue'

const route = useRoute()
const ifShow = computed(() => !['/login', '/register'].includes(route.path))
</script>

<template>
  <div class="app-container">
    <BackHome />
    <div class="sidebar">
      <Sidebar v-if="ifShow" />
    </div>
    <div class="page" v-if="ifShow">
      <router-view />
    </div>
    <div class="auth" v-else>
      <router-view />
    </div>
  </div>
</template>

<style scoped>
.app-container {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: row;
}

.sidebar {
  flex-shrink: 0;
}

.page {
  flex: 1;
  display: flex;
  flex-direction: column;
  width: calc(100% - 220px);
  min-height: auto;
  max-width: auto;
  margin-left: 220px;
  margin-right: none;
  padding: 25px 56px 30px;
}

.auth {
  flex: 1;
  display: flex;
  flex-direction: column;
  width: 100%;
  min-height: auto;
  max-width: auto;
  margin-left: none;
  margin-right: none;
}
</style>
