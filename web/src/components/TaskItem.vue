<template>
  <a-card style="margin: 0 12px">
    <a-row>
      <a-col>
        <span>{{ taskData.torrent_name }}</span>
      </a-col>
      <a-col>
        <div style="padding-top: 12px; padding-bottom: 12px">
          <a-progress :percent="30" :strokeWidth="8"/>
        </div>
      </a-col>
      <a-col :lg="12" :sm="24" :xs="24">
        <div style="line-height: 16px">
          <p>{{ taskData.status && taskData.status.bytes_read_data }}/300Kb</p>
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
export default {
  name: "TaskItem",
  props: {
    taskData: {
      type: Object,
      require: true
    }
  }
}
</script>

<style scoped>

</style>
