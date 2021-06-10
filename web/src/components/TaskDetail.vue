<template>
  <a-drawer :title="taskData.torrent_name" :width="drawerWidth"
            placement="right" :closable="true" :visible="visible"
            @close="onClose">
    <a-descriptions title="任务信息" layout="vertical" :column="{ md: 1, sm: 1, xs: 1 }">
      <a-descriptions-item label="任务名称">
        {{ taskData.torrent_name }}
      </a-descriptions-item>
      <a-descriptions-item label="Hash">
        <a-space>
          <span> {{ taskData.info_hash }} </span>
          <a-button v-if="taskData.meta_info" type="link" class="download-link"
                    @click="handleDownloadTorrent(taskData.info_hash)">
            <a-icon type="arrow-down"/>
            种子下载
          </a-button>
        </a-space>
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

    <div>
      <div style="text-align: center">
        <a-spin :spinning="taskDataLoading"/>
      </div>
      <div v-if="!taskDataLoading && taskData.torrentData">
        <a-descriptions title="Torrent" layout="vertical" :column="{ md: 3, sm: 1, xs: 1 }">
          <a-descriptions-item label="文件数量">
            {{ taskData.torrentData.files.length }}
          </a-descriptions-item>
          <a-descriptions-item label="Pieces">
            {{ pieces }}
          </a-descriptions-item>
        </a-descriptions>

        <file-tree v-if="visible"
                   item-slot="detail"
                   @on-download="handleDownload"
                   @on-item-click="handleItemClick"
                   :torrent-data="taskData.torrentData"
                   :default-checked-keys="taskData.download_files"
                   :disable-checkbox="true"></file-tree>
      </div>
    </div>

  </a-drawer>
</template>

<script>
import byteSize from 'byte-size'
import FileTree from "./FileTree";
import {baseURL} from "../http/api";

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
    },
    taskDataLoading: {
      type: Boolean,
      default: true
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
    },
    drawerWidth() {
      let w = window.innerWidth;
      if (w < 680) {
        return w
      }
      return 680
    }
  },
  methods: {
    onClose() {
      this.$emit('on-close')
    },
    handleDownload(key) {
      window.open(this.downloadUrl(key), '_blank')
    },
    handleItemClick(key) {
      let textArea = document.createElement("textarea")
      textArea.style.position = 'fixed'
      textArea.style.top = '0'
      textArea.style.left = '0'
      textArea.style.width = '2em'
      textArea.style.height = '2em'
      textArea.style.padding = '0'
      textArea.style.border = 'none'
      textArea.style.outline = 'none'
      textArea.style.boxShadow = 'none'
      textArea.style.background = 'transparent'
      if (baseURL) {
        textArea.value = this.downloadUrl(key)
      } else {
        textArea.value = "http://" + location.host + this.downloadUrl(key)
      }
      document.body.appendChild(textArea)
      textArea.focus()
      textArea.select()
      document.execCommand('copy')
      let successful = false
      try {
        successful = document.execCommand('copy')
      } catch (err) {
        console.log(err)
      }
      document.body.removeChild(textArea)
      if (successful) {
        this.$message.success(`下载链接复制成功`)
      } else {
        this.$message.error(`下载链接复制失败`)
      }
    },
    handleDownloadTorrent(hash) {
      window.open(`${baseURL}/torrent/download/${hash}`, '_blank')
    },
    time(timestamp) {
      let date = new Date(timestamp)
      return `${date.getFullYear()}-${this.padStart(date.getMonth() + 1)}-${this.padStart(date.getDate())}
      ${this.padStart(date.getHours())}:${this.padStart(date.getMinutes())}:${this.padStart(date.getSeconds())}`
    },
    padStart(val) {
      return val.toString().padStart(2, '0')
    },
    downloadUrl(key) {
      return `${baseURL}/torrent/file/download?hash=${this.taskData.info_hash}&path=${key}`
    }
  }
}
</script>

<style scoped lang="less">
/deep/ .ant-tree li .ant-tree-node-content-wrapper:hover {
  background-color: transparent !important;
  cursor: default;
}

.download-link {
  padding: 0 !important;
}
</style>
