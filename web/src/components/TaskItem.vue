<template>
  <div>
    <a-card class="item">
      <a-row>
        <a-col :lg="16" :sm="24" :xs="24">
          <a-badge :status="taskStatus"/>
          <span class="name">{{ taskData.torrent_name }}</span>
        </a-col>
        <a-col :lg="8" :sm="24" :xs="24">
          <div class="file-size" v-if="taskData.file_length">
            <span>{{ fileSize }}</span>
          </div>
        </a-col>
        <a-col class="progress">
          <div>
            <a-progress :percent="percent" :strokeWidth="8"/>
          </div>
        </a-col>
        <a-col :lg="12" :sm="24" :xs="24" class="item-bottom">
          <div class="status" v-if="active">
            <div class="icon-status" v-if="taskData.queue">
              <a-icon type="hourglass"/>
              <span>正在队列等待</span>
            </div>

            <div class="icon-status" v-else-if="!taskData.meta_info">
              <a-spin>
                <a-icon slot="indicator" type="loading" style="font-size: 16px" spin/>
              </a-spin>
              <span>正在获取信息</span>
            </div>

            <a-space :size="16" v-else>
            <span>
             <a-icon type="arrow-down"/> {{ read.speed || '0 B' }}
            </span>
              <span>
             <a-icon type="arrow-up"/> {{ write.speed || '0 B' }}
            </span>
              <span>
              Peers {{ peers }}
            </span>
              <span v-if="remainingTime">
              {{ remainingTime }}
            </span>
            </a-space>
          </div>
        </a-col>

        <a-col :lg="12" :sm="24" :xs="24" class="item-bottom">
          <div style="float: right">
            <a-space>
              <a-dropdown-button
                  :disabled="taskData.wait"
                  @click="handleTaskStart"
                  v-if="!active">
                <a-icon type="vertical-align-bottom"/>
                开始
                <a-menu slot="overlay">
                  <a-menu-item @click="handlePartDownload" :disabled="!taskData.meta_info">
                    <a-icon type="unordered-list"/>
                    部分下载
                  </a-menu-item>
                </a-menu>
              </a-dropdown-button>
              <a-button icon="pause"
                        :disabled="taskData.wait"
                        @click="handleTaskStop"
                        v-else>
                暂停
              </a-button>

              <a-button icon="stock" @click="handleShowDetail">
                详情
              </a-button>

              <a-popconfirm
                  :title="`确认删除?`"
                  ok-text="确认"
                  cancel-text="取消"
                  @confirm="handleDelete">
                <a-button icon="delete" type="danger">
                  删除
                </a-button>
              </a-popconfirm>
            </a-space>
          </div>
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script>
import byteSize from 'byte-size'

export default {
  name: "TaskItem",
  props: {
    taskData: {
      type: Object,
      require: true
    }
  },
  data() {
    return {
      read: {
        last: 0,
        auto: 0,
        speed: '0 B'
      },
      write: {
        last: 0,
        auto: 0,
        speed: '0 B'
      },
      remainingTime: '',
      speedTimer: null
    }
  },
  mounted() {
    this.speedTimer = setInterval(() => {
      if (this.read.last > 0) {
        if (this.read.last === this.read.auto) {
          this.read.last = 0
          this.read.auto = 0
          this.read.speed = '0 B'
          this.remainingTime = ''
        } else {
          this.read.auto = this.read.last
        }
      }
      if (this.write.last > 0) {
        if (this.write.last === this.write.auto) {
          this.write.last = 0
          this.write.auto = 0
          this.write.speed = '0 B'
        } else {
          this.write.auto = this.write.last
        }
      }
    }, 2000)
  },
  destroyed() {
    clearInterval(this.speedTimer)
  },
  methods: {
    handleTaskStart() {
      this.remainingTime = ''
      this.$emit("task-start", this.taskData.info_hash)
    },
    handleTaskStop() {
      this.$emit("task-stop", this.taskData.info_hash)
    },
    handlePartDownload() {
      this.$emit('task-part-download', this.taskData)
    },
    handleDelete() {
      this.$emit('task-delete', this.taskData.info_hash)
    },
    handleShowDetail() {
      this.$emit('task-detail', this.taskData)
    },
    calcTime(s) {
      let day = Math.floor(s / (24 * 3600))
      let hour = Math.floor((s - day * 24 * 3600) / 3600)
      let minute = Math.floor((s - day * 24 * 3600 - hour * 3600) / 60)
      let second = s - day * 24 * 3600 - hour * 3600 - minute * 60
      let time = []
      if (day) {
        time.push(`${day} 天`)
      }
      if (hour) {
        time.push(`${hour} 时`)
      }
      if (minute) {
        time.push(`${minute} 分`)
      }
      if (second) {
        time.push(`${second} 秒`)
      }
      if (time.length > 0) {
        return time.join(' ')
      } else {
        return ''
      }
    }
  },
  watch: {
    'taskData.stats.bytes_read_data'(val) {
      if (this.read.last === 0) {
        this.read.last = val
        this.read.speed = byteSize(0)
      } else {
        let sec = (this.taskData.stats.length - this.taskData.stats.bytes_completed) / (val - this.read.last)
        this.remainingTime = sec ? this.calcTime(Math.ceil(sec)) : ''

        this.read.speed = byteSize(val - this.read.last)
        this.read.last = val
      }
    },
    'taskData.stats.bytes_written_data'(val) {
      if (this.write.last === 0) {
        this.write.last = val
        this.write.speed = byteSize(0)
      } else {
        this.write.speed = byteSize(val - this.write.last)
        this.write.last = val
      }
    }
  },
  computed: {
    taskStatus() {
      if (this.taskData.complete) {
        return 'success'
      } else if (this.taskData.queue) {
        return 'warning'
      } else if (this.active) {
        return 'processing'
      }
      return 'default'
    },
    active() {
      return !this.taskData.complete && (this.taskData && !this.taskData.pause)
    },
    fileSize() {
      if (this.taskData.stats) {
        return `${byteSize(this.taskData.stats.bytes_completed)} / ${byteSize(this.taskData.stats.length)}`
      } else {
        return `${byteSize(this.taskData.complete_file_length)} / ${byteSize(this.taskData.file_length)}`
      }
    },
    percent() {
      let p = 0
      if (this.taskData.stats) {
        p = this.taskData.stats.bytes_completed / this.taskData.stats.length
      } else {
        p = this.taskData.complete_file_length / this.taskData.file_length
      }
      p = p * 100
      return p > 100 ? 100 : parseFloat(p.toFixed(2))
    },
    peers() {
      if (this.taskData.stats) {
        if (this.taskData.stats.total_peers > 0) {
          return `${this.taskData.stats.active_peers} / ${this.taskData.stats.total_peers}`
        }
      }
      return "0 / 0"
    }
  }
}
</script>

<style lang="less" scoped>
.item {
  margin: 0 16px;

  .progress {
    padding-top: 36px;
    padding-right: 12px;
  }

  .status {
    line-height: 24px;
  }

  .icon-status {
    display: inline-block;

    span {
      margin-left: 8px;
    }
  }

  .name {
    font-weight: bold;
    font-size: 16px;
  }

  .file-size {
    display: inline-block;
    float: right;
  }

  .item-bottom {
    margin-top: 18px;
  }
}
</style>
