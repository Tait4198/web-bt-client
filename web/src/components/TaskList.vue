<template>
  <div>
    <a-list item-layout="vertical" :split="false" :data-source="taskList">
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
import {TorrentInfo, TorrentStats, TorrentWait, TorrentAdd, TorrentPause, TorrentComplete} from "../constant/constant";

export default {
  name: "TaskList",
  components: {
    TaskItem
  },
  mounted() {
    this.typFuncMap.set(TorrentStats, this.torrentStats)
    this.typFuncMap.set(TorrentInfo, this.torrentInfo)
    this.typFuncMap.set(TorrentWait, this.torrentWait)
    this.typFuncMap.set(TorrentAdd, this.torrentAdd)
    this.typFuncMap.set(TorrentPause, this.torrentPause)
    this.typFuncMap.set(TorrentComplete, this.torrentComplete)

    getTaskList().then(res => {
      let {data = []} = res
      for (let i = 0; i < data.length; i++) {
        let item = data[i]
        this.$set(this.tasks, item.info_hash, this.newTask(item))
      }
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
    torrentStats(taskData, obj) {
      this.$set(taskData, 'stats', obj)
    },
    torrentInfo(taskData, obj) {
      taskData.torrent_name = obj.name
      taskData.file_length = obj.length
      taskData.complete_file_length = 0
      taskData.meta_info = true
      this.$notification.success({
        message: 'MetaInfo 获取完成',
        description: `任务 ${taskData.torrent_name} 完成信息获取`,
      })
    },
    torrentWait(taskData, obj) {
      taskData.wait = obj.status
    },
    torrentAdd(obj) {
      this.$set(this.tasks, obj.info_hash, this.newTask(obj))
    },
    torrentPause(taskData, obj) {
      taskData.pause = obj.status
    },
    torrentComplete(taskData, obj) {
      taskData.complete = obj.status
      taskData.complete_file_length = obj.last_complete_length
      if (taskData.stats) {
        taskData.stats.bytes_completed = obj.last_complete_length
      }
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
      item.wait = false
      return item
    }
  },
  data() {
    return {
      tasks: [],
      typFuncMap: new Map(),
    }
  }
}
</script>

<style scoped>

</style>
