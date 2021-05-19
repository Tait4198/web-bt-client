<template>
  <div style="height: 100%">
    <div class="header">
      <el-button icon="el-icon-plus" type="primary">
        添加任务
      </el-button>
    </div>

    <div>
    </div>

    <div class="footer">
      <p>Copyright © 2021 WEB_BT_CLIENT</p>
    </div>
  </div>
</template>

<script>

export default {
  name: "AppIndex",
  components: {},
  mounted() {

  },
  data() {
    return {
      tableData: []
    }
  }
}
</script>

<style scoped>
.header {
  position: fixed;
  z-index: 999999;
  left: 0;
  top: 0;
  width: 100%;
  padding: 12px;
  background: #fafafa;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1)
}

.footer {
  width: 100%;
  text-align: center;
  background: transparent;
  position: fixed;
  bottom: 0;
}
</style>