<template>
  <div style="height: 100%">
    <div class="header">
      <a-button icon="plus" type="primary" @click="handleAddTask">
        添加任务
      </a-button>
    </div>
    <div style="height: 64px"></div>
    <task-list ref="taskList"></task-list>

    <a-modal v-model="mgTask.visible" title="添加任务"
             :maskClosable="false" :destroyOnClose="true"
             @cancel="addTaskCancel" @ok="addTaskOk">
      <a-form-model ref="ruleForm" :model="mgTask.form" :rules="mgTask.rules" v-bind="mgTask.modalLayout">
        <a-form-model-item has-feedback label="磁力链接" prop="uri">
          <a-input v-model="mgTask.form.uri" type="text"/>
        </a-form-model-item>
        <a-form-model-item has-feedback label="文件路径" prop="path" ref="mgTaskPath">
          <a-input v-model="mgTask.form.path" type="text" @blur="getSpaceData"/>
          <div>
            <span>剩余空间 {{ pathFreeSpace }}</span>
            <div style="display: inline-block; float: right;">
              <a-button @click="pathSelect.visible = true" icon="select" size="small">
                选择下载路径
              </a-button>
            </div>

          </div>
        </a-form-model-item>
      </a-form-model>
    </a-modal>

    <a-modal v-model="pathSelect.visible" title="路径选择"
             :maskClosable="false" :destroyOnClose="true"
             @ok="selectPathOk">
      <div class="file-select">
        <path-tree @on-path-select="handlePathSelect"></path-tree>
      </div>
    </a-modal>
  </div>
</template>

<script>
import TaskList from "./TaskList";
import PathTree from "./PathTree";
import byteSize from 'byte-size'
import {getSpace} from "../http/api";

export default {
  name: "TaskIndex",
  components: {
    PathTree,
    TaskList
  },
  mounted() {
  },
  computed: {
    pathFreeSpace() {
      return byteSize(this.mgTask.pathFreeSpace)
    }
  },
  methods: {
    handleAddTask() {
      this.mgTask.visible = true
    },
    addTaskCancel() {
      this.mgTask.form.uri = ''
    },
    addTaskOk() {

    },
    selectPathOk() {
      this.mgTask.form.path = this.pathSelect.tempPath
      this.pathSelect.visible = false
      console.log(this.$refs.mgTaskPath)
      this.$refs.mgTaskPath.onFieldChange()
    },
    handlePathSelect(path) {
      this.pathSelect.tempPath = path
    },
    getSpaceData() {

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
      callback()
    }

    let validatePath = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('路径不能为空'));
      }
      getSpace({path: this.mgTask.form.path}).then(res => {
        let {data} = res
        if (data) {
          this.mgTask.pathFreeSpace = data
          callback()
        } else {
          callback(new Error('当前路径无效'));
        }
      })
    }

    return {
      pathSelect: {
        visible: false,
        tempPath: '',
      },
      mgTask: {
        visible: false,
        pathFreeSpace: 0,
        modalLayout: {
          labelCol: {
            span: 4
          },
          wrapperCol: {
            span: 20
          }
        },
        form: {
          uri: '',
          path: ''
        },
        rules: {
          uri: [{validator: validateUri, trigger: 'change'}],
          path: [{validator: validatePath, trigger: 'change'}],
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

.file-select {
  max-height: 360px;
  overflow: auto;
}
</style>
