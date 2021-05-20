<template>
  <div id="app">
    <app-index></app-index>
  </div>
</template>

<script>

import AppIndex from "@/components/AppIndex";

export default {
  name: 'App',
  components: {
    AppIndex
  },
  data() {
    return {}
  }
}
</script>

<style>

html, body {
  height: 100%;
  margin: 0;
}

#app {
  height: 100%;
}
</style>
