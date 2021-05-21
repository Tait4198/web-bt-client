<template>
  <a-card class="item">
    <a-row>
      <a-col>
        <span class="name">{{ taskData.torrent_name }}</span>
        <div class="file-size">
          <span>{{ fileSize }}</span>
        </div>
      </a-col>
      <a-col>
        <div class="progress">
          <a-progress :percent="percent" :strokeWidth="8"/>
        </div>
      </a-col>
      <a-col :lg="12" :sm="24" :xs="24">
        <div class="status">
          <a-space :size="16">
            <span>
               <a-icon type="arrow-down"/>
                {{ read.speed || '0 B' }}
            </span>

            <span>
               <a-icon type="arrow-up"/>
                {{ write.speed || '0 B' }}
            </span>

            <span v-if="!active">
              Peers
              {{ peers }}
            </span>
          </a-space>
        </div>
      </a-col>

      <a-col :lg="12" :sm="24" :xs="24">
        <div style="float: right">
          <a-space>
            <a-button icon="arrow-right"
                      type="primary"
                      :disabled="taskData.wait"
                      @click="handleTaskStart"
                      v-if="active">
              开始
            </a-button>
            <a-button icon="pause"
                      :disabled="taskData.wait"
                      @click="handleTaskStop"
                      v-else>
              暂停
            </a-button>

            <a-button icon="stock">
              详情
            </a-button>
            <a-button icon="delete" type="danger">
              删除
            </a-button>
          </a-space>
        </div>
      </a-col>
    </a-row>
  </a-card>
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
      this.$emit("task-start", this.taskData.info_hash)
    },
    handleTaskStop() {
      this.$emit("task-stop", this.taskData.info_hash)
    }
  },
  watch: {
    'taskData.stats.bytes_read_data'(val) {
      if (this.read.last === 0) {
        this.read.last = val
        this.read.speed = byteSize(0)
      } else {
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
    active() {
      return this.taskData && this.taskData.pause
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
    padding-top: 16px;
    padding-bottom: 16px
  }

  .status {
    line-height: 16px
  }

  .name {
    font-weight: bold;
    font-size: 16px;
  }

  .file-size {
    display: inline-block;
    float: right;
  }
}
</style>
