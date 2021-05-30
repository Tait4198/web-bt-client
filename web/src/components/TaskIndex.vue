<template>
  <div style="height: 100%">
    <div class="header">
      <div style="display: inline-block">
        <h3>WEB-BT-CLIENT</h3>
      </div>
      <div style="display: inline-block;float: right">
        <a-button icon="plus" type="primary" @click="handleShowTaskModal">
          添加任务
        </a-button>
      </div>
    </div>
    <div style="height: 72px"></div>
    <task-list ref="taskList"></task-list>
    <task-modal ref="taskModal"></task-modal>
  </div>
</template>

<script>
import TaskList from "./TaskList";
import TaskModal from "./TaskModal";

export default {
  name: "TaskIndex",
  components: {
    TaskList,
    TaskModal,
  },
  mounted() {
  },
  methods: {
    handleShowTaskModal() {
      this.$refs.taskModal.show()
    }
  },
  data() {
    return {}
  }
}
</script>

<style scoped>
.header {
  position: fixed;
  z-index: 999;
  left: 0;
  top: 0;
  width: 100%;
  padding: 12px;
  background: #fafafa;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}
</style>
