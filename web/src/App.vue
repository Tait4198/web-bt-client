<template>
  <div id="app">
    <index></index>
  </div>
</template>

<script>

import Index from "@/components/Index";

export default {
  name: 'App',
  components: {
    Index
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
