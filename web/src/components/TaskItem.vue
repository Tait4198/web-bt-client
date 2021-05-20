<template>
  <a-card class="item">
    <a-row>
      <a-col>
        <span>{{ taskData.torrent_name + ' ' + taskData.info_hash }}</span>
      </a-col>
      <a-col>
        <div class="progress">
          <a-progress :percent="percent" :strokeWidth="8"/>
        </div>
      </a-col>
      <a-col :lg="12" :sm="24" :xs="24">
        <div class="status">
          <a-space :size="20">
            <span>{{ read.speed }} / {{ write.speed }}</span>
            <span v-if="peers !== ''">Peers {{ peers }}</span>
          </a-space>
        </div>
      </a-col>

      <a-col :lg="12" ::sm="24" :xs="24">
        <div style="float: right">
          <a-space>
            <a-button icon="pause">
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
    }, 1000)
  },
  destroyed() {
    clearInterval(this.speedTimer)
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
      return ""
    }
  }
}
</script>

<style lang="less" scoped>
.item {
  margin: 0 12px;

  .progress {
    padding-top: 12px;
    padding-bottom: 12px
  }

  .status {
    line-height: 16px
  }
}
</style>
