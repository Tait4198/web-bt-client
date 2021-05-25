<template>
  <div id="app">
    <task-index></task-index>
  </div>
</template>

<script>

import TaskIndex from "@/components/TaskIndex";

export default {
  name: 'App',
  components: {
    TaskIndex,
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
