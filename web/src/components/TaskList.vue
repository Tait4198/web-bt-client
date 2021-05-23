<template>
  <div>
    <a-list item-layout="vertical" :split="false" :data-source="taskList" :loading="taskListLoading">
      <a-list-item slot="renderItem" :key="index" slot-scope="item, index">
        <task-item :task-data="item"
                   @task-start="handleTaskStart"
                   @task-stop="handleTaskStop">
        </task-item>
      </a-list-item>
    </a-list>
  </div>
</template>

<script>
import {getTaskList, taskStart, taskStop} from "@/http/task";
import TaskItem from "./TaskItem"
import {
  TaskInfo,
  TaskStats,
  TaskWait,
  TaskAdd,
  TaskPause,
  TaskComplete,
  TaskQueueStatus,
} from "../constant/constant";

export default {
  name: "TaskList",
  components: {
    TaskItem
  },
  mounted() {
    this.typFuncMap.set(TaskStats, this.taskStats)
    this.typFuncMap.set(TaskInfo, this.taskInfo)
    this.typFuncMap.set(TaskWait, this.taskWait)
    this.typFuncMap.set(TaskAdd, this.taskAdd)
    this.typFuncMap.set(TaskPause, this.taskPause)
    this.typFuncMap.set(TaskComplete, this.taskComplete)
    this.typFuncMap.set(TaskQueueStatus, this.taskQueueStatus)

    getTaskList().then(res => {
      let {data = []} = res
      for (let i = 0; i < data.length; i++) {
        let item = data[i]
        this.$set(this.tasks, item.info_hash, this.newTask(item))
      }
      this.taskListLoading = false
    })

    this.$bus.on("ws-message", (e) => {
      e.data.split('\n').forEach(item => {
        if (item) {
          let obj = JSON.parse(item)
          let infoHash = obj.info_hash
          let type = obj.type
          if (this.typFuncMap.has(type)) {
            if (this.tasks[infoHash]) {
              this.typFuncMap.get(type)(this.tasks[infoHash], obj)
            } else {
              this.typFuncMap.get(type)(obj)
            }
          }
        }
      })
    })
  },
  computed: {
    taskList() {
      return Object.values(this.tasks).sort((o1, o2) => o2.create_time - o1.create_time)
    }
  },
  methods: {
    taskStats(taskData, obj) {
      this.$set(taskData, 'stats', obj)
    },
    taskInfo(taskData, obj) {
      taskData.torrent_name = obj.name
      taskData.file_length = obj.length
      taskData.complete_file_length = 0
      taskData.meta_info = true
      this.$notification.success({
        message: 'MetaInfo 获取完成',
        description: `任务 ${taskData.torrent_name} 完成信息获取`,
      })
    },
    taskWait(taskData, obj) {
      taskData.wait = obj.status
    },
    taskAdd(obj) {
      this.$set(this.tasks, obj.info_hash, this.newTask(obj))
    },
    taskPause(taskData, obj) {
      taskData.pause = obj.status
    },
    taskComplete(taskData, obj) {
      taskData.complete = obj.status
      taskData.complete_file_length = obj.last_complete_length
      if (taskData.stats) {
        taskData.stats.bytes_completed = obj.last_complete_length
      }
    },
    taskQueueStatus(taskData, obj) {
      taskData.queue = obj.status
    },

    handleTaskStart(infoHash) {
      taskStart({
        info_hash: infoHash,
        download: true,
        update: true
      }).then(res => {
        if (res.status) {
          this.tasks[infoHash].pause = false
          this.tasks[infoHash].wait = true
          setTimeout(() => {
            this.tasks[infoHash].wait = false
          }, 3000)
          this.$message.success(`任务 ${this.tasks[infoHash].torrent_name} 启动成功`)
        } else {
          this.$message.error(res.message)
        }
      })
    },
    handleTaskStop(infoHash) {
      taskStop({
        hash: infoHash
      }).then(res => {
        if (res.status) {
          this.tasks[infoHash].pause = true
          this.tasks[infoHash].wait = true
          setTimeout(() => {
            this.tasks[infoHash].wait = false
          }, 3000)
          this.$message.success(`任务 ${this.tasks[infoHash].torrent_name} 暂停成功`)
        } else {
          this.$message.error(res.message)
        }
      })
    },

    newTask(item) {
      if (!item.wait) {
        item.wait = false
      }
      if (!item.queue) {
        item.queue = false
      }
      return item
    }
  },
  data() {
    return {
      tasks: [],
      typFuncMap: new Map(),
      taskListLoading: true
    }
  }
}
</script>

<style scoped>

</style>
