<script lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { RouterLink } from 'vue-router'

const currentTheme = ref('light')

async function changeTheme() {
  document.documentElement.className = currentTheme.value
  localStorage
}

</script>

<template>
  <header class="header">
    <div class="header-top">
      <h1 class="logo" v-if="$route.meta.isAdminPage">Amai: Admin</h1>
      <h1 class="logo">Amai</h1>
      <div class="settings">
        <div class="settings-item">
          <label for="theme-select">theme: </label>
          <select id="theme-select" v-model="currentTheme" @change="changeTheme">

          </select>
        </div>
      </div>
      <nav>
        <RouterLink to="/"><h3>main</h3></RouterLink>
        <RouterLink to="/about"><h3>about</h3></RouterLink>
        <RouterLink to="/articles"><h3>articles</h3></RouterLink>
      </nav>
    </div>
  </header>
</template>