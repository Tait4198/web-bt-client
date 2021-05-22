<template>
  <div style="height: 100%">
    <div class="header">
      <a-button icon="plus" type="primary" @click="handleAddTask">
        添加任务
      </a-button>
    </div>
    <div style="height: 64px"></div>
    <task-list ref="taskList"></task-list>

    <a-modal v-model="addTaskVisible" title="添加任务"
             :maskClosable="false" :destroyOnClose="true"
             @cancel="addTaskCancel" @="addTaskOk">
      <a-form-model ref="ruleForm" :model="mgTask.form" :rules="mgTask.rules" v-bind="mgTask.modalLayout">
        <a-form-model-item has-feedback label="磁力链接" prop="uri">
          <a-input v-model="mgTask.form.uri" type="text"/>
        </a-form-model-item>
      </a-form-model>
    </a-modal>
  </div>
</template>

<script>
import TaskList from "./TaskList";

export default {
  name: "TaskIndex",
  components: {
    TaskList
  },
  mounted() {
  },
  methods: {
    handleAddTask() {
      this.addTaskVisible = true
    },
    addTaskCancel() {
      this.mgTask.form.uri = ''
    },
    addTaskOk() {

    }
  },
  data() {
    let validateUri = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入磁力链接'));
      }
      if (!value.match(/magnet:\?xt=urn:btih:[a-z0-9]{40}.*/)) {
        callback(new Error('无效磁力链接'));
      }
      callback();
    }

    return {
      addTaskVisible: false,
      mgTask: {
        modalLayout: {
          labelCol: {
            span: 4
          },
          wrapperCol: {
            span: 20
          }
        },
        form: {
          uri: ''
        },
        rules: {
          uri: [{validator: validateUri, trigger: 'change'}],
        }
      },
    }
  }
}
</script>

<style scoped>
.header {
  position: fixed;
  z-index: 999;
  left: 0;
  top: 0;
  width: 100%;
  padding: 12px;
  background: #fafafa;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1)
}
</style>
