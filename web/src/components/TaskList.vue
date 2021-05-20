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
import {TorrentInfo, TorrentStats} from "../constant/constant";

export default {
  name: "TaskList",
  components: {
    TaskItem
  },
  mounted() {
    this.typFuncMap.set(TorrentStats, this.torrentStats)
    this.typFuncMap.set(TorrentInfo, this.torrentInfo)

    getTaskList().then(res => {
      let {data = []} = res
      for (let i = 0; i < data.length; i++) {
        let item = data[i]
        this.$set(this.tasks, item.info_hash, item)
      }
    })

    this.$bus.on("ws-message", (e) => {
      let obj = JSON.parse(e.data)
      let infoHash = obj.info_hash
      let type = obj.type
      if (this.typFuncMap.has(type) && this.tasks[infoHash]) {
        this.typFuncMap.get(type)(this.tasks[infoHash], obj)
      }
    })
  },
  computed: {
    taskList() {
      return Object.values(this.tasks).sort((o1, o2) => o1.id - o2.id)
    }
  },
  methods: {
    torrentStats(taskData, obj) {
      this.$set(taskData, 'stats', obj)
    },
    torrentInfo(taskData, obj) {
      taskData.torrent_name = obj.name
    }
  },
  data() {
    return {
      tasks: [],
      typFuncMap: new Map()
    }
  }
}
</script>

<style scoped>

</style>
