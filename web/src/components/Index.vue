<template>
  <div style="height: 100%">
    <div class="header">
      <a-button icon="plus" type="primary">
        添加任务
      </a-button>
    </div>
    <div style="height: 64px"></div>
    <task-list></task-list>
  </div>
</template>

<script>
import TaskList from "./TaskList";

export default {
  name: "AppIndex",
  components: {
    TaskList
  },
  mounted() {

  },
  data() {
    return {}
  }
}
</script>

<style scoped>
.header {
  position: fixed;
  z-index: 999999;
  left: 0;
  top: 0;
  width: 100%;
  padding: 12px;
  background: #fafafa;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1)
}
</style>
