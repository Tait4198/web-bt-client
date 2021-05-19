<template>
  <div>
    <a-list item-layout="vertical" :split="false" :data-source="taskList">
      <a-list-item slot="renderItem" :key="index" slot-scope="item, index">
        <task-item :task-data="item"></task-item>
      </a-list-item>
    </a-list>
  </div>
</template>

<script>
import {getTaskList} from "@/http/task";
import TaskItem from "./TaskItem"

export default {
  name: "TaskList",
  components: {
    TaskItem
  },
  mounted() {
    getTaskList().then(res => {
      let {data = []} = res
      for (let i = 0; i < data.length; i++) {
        let item = data[i]
        this.$set(this.tasks, item.info_hash, item)
      }
    })

    this.$bus.on("ws-message", (e) => {
      let obj = JSON.parse(e.data)
      if (this.tasks[obj.info_hash]) {
        this.$set(this.tasks[obj.info_hash], 'stats', obj)
        console.log(this.tasks[obj.info_hash])
      }
    })
  },
  computed: {
    taskList() {
      return Object.values(this.tasks).sort((o1, o2) => o1.id - o2.id)
    }
  },
  data() {
    return {
      tasks: []
    }
  }
}
</script>

<style scoped>

</style>
