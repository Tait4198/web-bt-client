<template>
  <a-drawer width="700" placement="right" :closable="false" :visible="visible" @close="onClose">
    <a-descriptions title="任务信息" layout="vertical" :column="{ md: 1, sm: 1, xs: 1 }">
      <a-descriptions-item label="任务名称">
        {{ taskData.torrent_name }}
      </a-descriptions-item>
      <a-descriptions-item label="Hash">
        {{ taskData.info_hash }}
      </a-descriptions-item>
      <a-descriptions-item label="已下载">
        {{ fileSize }}
      </a-descriptions-item>
      <a-descriptions-item label="下载路径">
        {{ taskData.download_path }}
      </a-descriptions-item>
      <a-descriptions-item label="创建时间">
        {{ time(taskData.create_time / 1000000) }}
      </a-descriptions-item>
    </a-descriptions>

    <div v-if="taskData.torrentData">
      <a-descriptions title="Torrent" layout="vertical" :column="{ md: 2, sm: 1, xs: 1 }">
        <a-descriptions-item label="文件数量">
          {{ taskData.torrentData.files.length }}
        </a-descriptions-item>
        <a-descriptions-item label="Pieces">
          {{ pieces }}
        </a-descriptions-item>
      </a-descriptions>

      <file-tree v-if="visible"
                 :torrent-data="taskData.torrentData"
                 :default-checked-keys="taskData.download_files"
                 :disable-checkbox="true"></file-tree>
    </div>

  </a-drawer>
</template>

<script>
import byteSize from 'byte-size'
import FileTree from "./FileTree";

export default {
  name: "TaskDetail",
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    taskData: {
      type: Object,
      default: () => {
        return {}
      }
    }
  },
  components: {
    FileTree
  },
  computed: {
    fileSize() {
      if (this.taskData.stats) {
        return `${byteSize(this.taskData.stats.bytes_completed)} / ${byteSize(this.taskData.stats.length)}`
      } else {
        return `${byteSize(this.taskData.complete_file_length)} / ${byteSize(this.taskData.file_length)}`
      }
    },
    pieces() {
      if (this.taskData.stats) {
        return `${this.taskData.stats.completed_pieces} / ${this.taskData.stats.pieces}`
      } else {
        return `${this.taskData.torrentData.completed_pieces} / ${this.taskData.torrentData.pieces}`
      }
    }
  },
  methods: {
    onClose() {
      this.$emit('on-close')
    },
    byteSize(val) {
      return byteSize(val || 0)
    },
    time(timestamp) {
      let date = new Date(timestamp)
      return `${date.getFullYear()}-${this.padStart(date.getMonth() + 1)}-${this.padStart(date.getDate())}
      ${this.padStart(date.getHours())}:${this.padStart(date.getMinutes())}:${this.padStart(date.getSeconds())}`
    },
    padStart(val) {
      return val.toString().padStart(2, '0')
    }
  }
}
</script>

<style scoped>

</style>