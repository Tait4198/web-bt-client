<template>
  <div>
    <a-list item-layout="vertical" :split="false" :data-source="taskList" :loading="taskListLoading">
      <a-list-item slot="renderItem" :key="index" slot-scope="item, index">
        <task-item :task-data="item"
                   @task-start="handleTaskStart"
                   @task-stop="handleTaskStop"
                   @task-part-download="handleTaskPartDownload"
                   @task-detail="handleShowDetail"
                   @task-delete="handleTaskDelete">
        </task-item>
      </a-list-item>
    </a-list>

    <a-modal v-model="partDownload.visible" title="部分下载" okText="开始下载"
             :maskClosable="false" :destroyOnClose="true" :width="780"
             :confirm-loading="partDownload.okLoading"
             @ok="handleFileCheckOk">
      <div class="tree-select">
        <div style="text-align:center; height: 36px;" v-if="partDownload.loading">
          <a-spin :spinning="partDownload.loading">
            <a-icon slot="indicator" type="loading" style="font-size: 24px;" spin/>
          </a-spin>
        </div>
        <file-tree v-else @on-file-check="handleFileCheck"
                   :default-checked-keys="partDownload.taskData.download_files"
                   :torrent-data="partDownload.taskData.torrentData"></file-tree>
      </div>
    </a-modal>

    <task-detail :visible="taskDetail.visible"
                 :task-data="taskDetail.taskData"
                 @on-close="handleCloseDetail"></task-detail>
  </div>
</template>

<script>
import {getTaskList, taskStart, taskStop, taskRestart, taskDelete, getTorrentInfo} from "@/http/api";
import TaskItem from "./TaskItem"
import FileTree from "./FileTree";
import TaskDetail from "./TaskDetail";
import {
  TaskInfo,
  TaskStats,
  TaskWait,
  TaskAdd,
  TaskPause,
  TaskComplete,
  TaskQueueStatus,
  TaskDelete,
} from "../constant/constant";

export default {
  name: "TaskList",
  components: {
    TaskItem,
    FileTree,
    TaskDetail
  },
  mounted() {
    this.typFuncMap.set(TaskStats, this.wsTaskStats)
    this.typFuncMap.set(TaskInfo, this.wsTaskInfo)
    this.typFuncMap.set(TaskWait, this.wsTaskWait)
    this.typFuncMap.set(TaskAdd, this.wsTaskAdd)
    this.typFuncMap.set(TaskPause, this.wsTaskPause)
    this.typFuncMap.set(TaskComplete, this.wsTaskComplete)
    this.typFuncMap.set(TaskQueueStatus, this.wsTaskQueueStatus)
    this.typFuncMap.set(TaskDelete, this.wsTaskDelete)

    getTaskList().then(res => {
      let {data} = res
      if (data) {
        for (let i = 0; i < data.length; i++) {
          let item = data[i]
          this.$set(this.tasks, item.info_hash, this.newTask(item))
        }
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
    wsTaskStats(taskData, obj) {
      this.$set(taskData, 'stats', obj)
    },
    wsTaskInfo(taskData, obj) {
      taskData.torrent_name = obj.name
      taskData.file_length = obj.length
      taskData.complete_file_length = 0
      taskData.meta_info = true
      this.$notification.success({
        message: 'MetaInfo 获取完成',
        description: `任务 ${taskData.torrent_name} 完成信息获取`,
      })
    },
    wsTaskWait(taskData, obj) {
      taskData.wait = obj.status
    },
    wsTaskAdd(obj) {
      this.$set(this.tasks, obj.info_hash, this.newTask(obj))
      this.$notification.success({
        message: '任务创建成功',
        description: `任务 ${obj.torrent_name} 已创建成功`,
      })
    },
    wsTaskPause(taskData, obj) {
      taskData.pause = obj.status
    },
    wsTaskComplete(taskData, obj) {
      taskData.complete = obj.status
      taskData.complete_file_length = obj.last_complete_length
      if (taskData.stats) {
        taskData.stats.bytes_completed = obj.last_complete_length
      }
    },
    wsTaskQueueStatus(taskData, obj) {
      taskData.queue = obj.status
    },
    wsTaskDelete(taskData, obj) {
      if (obj.status) {
        this.$delete(this.tasks, taskData.info_hash)
      }
    },

    handleTaskStart(infoHash) {
      taskStart({
        info_hash: infoHash,
        download: true,
        update: true
      }).then(res => {
        if (res.status) {
          this.taskStart(this.tasks[infoHash])
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
    handleTaskPartDownload(taskData) {
      this.partDownload.taskData = taskData
      this.partDownload.checkedFilesUpdate = false
      this.partDownload.visible = true
      this.partDownload.tempCheckedKeys = []

      if (!taskData.torrentData) {
        this.partDownload.loading = true
        let that = this
        this.getTorrent(taskData, function () {
          that.partDownload.loading = false
        })
      } else {
        this.partDownload.taskData = taskData
        if (taskData.download_all) {
          this.partDownload.download_files = taskData.torrentData.files.map(f => f.path.join("/"))
        }
      }
    },
    handleFileCheckOk() {
      this.partDownload.okLoading = true
      let param = {
        info_hash: this.partDownload.taskData.info_hash,
        download: true,
        update: true
      }
      let taskData = this.partDownload.taskData
      if (this.partDownload.checkedFilesUpdate) {
        param.download_files = this.partDownload.tempCheckedKeys
        param.download_all = false
      }
      taskRestart(param).then(res => {
        let {data} = res
        if (data) {
          if (this.partDownload.checkedFilesUpdate) {
            taskData.download_files = this.partDownload.tempCheckedKeys
          }
          this.taskStart(taskData)
        }
      }).finally(() => {
        this.partDownload.okLoading = false
        this.partDownload.visible = false
      })
    },
    handleFileCheck(keys) {
      this.partDownload.tempCheckedKeys = keys
      this.partDownload.checkedFilesUpdate = true
    },
    handleTaskDelete(infoHash) {
      taskDelete({
        hash: infoHash
      }).then(res => {
        if (res.status) {
          this.$message.success(`任务 ${this.tasks[infoHash].torrent_name} 删除成功`)
        } else {
          this.$message.error(res.message)
        }
      })
    },
    handleShowDetail(taskData) {
      this.taskDetail.taskData = taskData
      this.taskDetail.visible = true

      // if (!taskData.torrentData) {
      this.taskDetail.loading = true
      let that = this
      this.getTorrent(taskData, function () {
        that.taskDetail.loading = false
      })
      // }
    },
    handleCloseDetail() {
      this.taskDetail.visible = false
    },

    getTorrent(taskData, finallyCb) {
      getTorrentInfo({hash: taskData.info_hash}).then(res => {
            let {data} = res
            if (data) {
              this.$set(taskData, 'torrentData', data)
              if (taskData.download_all) {
                taskData.download_files = data.files.map(f => f.path.join("/"))
              }
            }
          }
      ).finally(() => {
        if (finallyCb) {
          finallyCb()
        }
      })
    },
    taskStart(taskData) {
      taskData.pause = false
      taskData.wait = true
      taskData.complete = false
      setTimeout(() => {
        taskData.wait = false
      }, 3000)
      this.$message.success(`任务 ${taskData.torrent_name} 启动成功`)
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
      taskListLoading: true,
      partDownload: {
        visible: false,
        loading: false,
        okLoading: false,
        checkedFilesUpdate: false,
        tempCheckedKeys: [],
        checkedKeys: [],
        taskData: {
          info_hash: '',
          download_files: [],
          torrentData: {
            name: '',
            files: []
          }
        }
      },
      taskDetail: {
        visible: false,
        loading: false,
        taskData: {
          info_hash: '',
          download_files: [],
          torrentData: {
            name: '',
            files: []
          }
        }
      }
    }
  }
}
</script>

<style scoped>
.tree-select {
  max-height: 360px;
  overflow: auto;
}
</style>
