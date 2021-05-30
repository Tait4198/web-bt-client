<template>
  <div>
    <a-modal v-model="visible" title="添加任务"
             :maskClosable="false" :destroyOnClose="true"
             @cancel="cleanModal" @ok="handleTaskOk" @close="cleanModal"
             :confirm-loading="createTaskLoading">
      <a-form-model ref="ruleForm" :model="form" :rules="rules" v-bind="modalLayout">
        <a-form-item label="任务类型">
          <a-radio-group :default-value="taskType" v-model="taskType" button-style="solid">
            <a-radio-button value="uri">
              磁力链接
            </a-radio-button>
            <a-radio-button value="file">
              种子文件
            </a-radio-button>
          </a-radio-group>
        </a-form-item>

        <a-form-model-item v-if="taskType === 'uri'" label="磁力链接"
                           prop="magentUri" ref="magentUri">
          <a-input v-model="form.magentUri" type="text" placeholder="请输入有效的磁力链接"/>
        </a-form-model-item>

        <a-form-model-item v-else-if="taskType === 'file'" label="种子文件"
                           prop="torrentFile" ref="torrentFile">
          <a-upload
              :multiple="false"
              list-type="text"
              accept="application/x-bittorrent"
              :file-list="form.fileList"
              :before-upload="handleTorrentSelect"
              @change="handleTorrentChange">
            <a-button icon="upload" v-if="form.fileList.length < 1">
              上传种子
            </a-button>
          </a-upload>
        </a-form-model-item>

        <a-form-item label="下载类型">
          <a-radio-group :default-value="downloadType" v-model="downloadType" button-style="solid">
            <a-radio-button value="all">
              完整下载
            </a-radio-button>
            <a-radio-button value="part" :disabled="torrentData.files.length === 0">
              部分下载
            </a-radio-button>
          </a-radio-group>
        </a-form-item>

        <a-form-item label="文件选择" v-if="downloadType === 'part'">
          <a-button @click="fileCheck.visible = true">选择下载文件</a-button>
        </a-form-item>

        <a-form-item label="立即下载" v-if="downloadType === 'all'">
          <a-tooltip placement="topLeft">
            <span slot="title">是否在完成成功获取MetaInfo后立即开始下载</span>
            <a-switch default-checked @change="handleDownloadNowChange">
              <a-icon slot="checkedChildren" type="check"/>
              <a-icon slot="unCheckedChildren" type="close"/>
            </a-switch>
          </a-tooltip>
        </a-form-item>

        <a-form-model-item label="下载路径" prop="downloadPath" ref="downloadPath">
          <a-tooltip placement="topLeft">
            <span slot="title">{{ freeSpace }}</span>
            <a-input v-model="form.downloadPath" placeholder="请输入或选择有效的下载路径"
                     type="text" @blur="validateSpaceData">
              <div slot="addonAfter">
                <a-icon type="select" @click="pathSelect.visible = true"/>
              </div>
            </a-input>
          </a-tooltip>
        </a-form-model-item>
      </a-form-model>
    </a-modal>

    <a-modal v-model="pathSelect.visible" title="路径选择"
             :maskClosable="false" :destroyOnClose="true"
             @ok="handlePathSelectOk">
      <div class="tree-select">
        <path-tree @on-path-select="handlePathSelect"></path-tree>
      </div>
    </a-modal>

    <a-modal v-model="fileCheck.visible" title="文件选择"
             :maskClosable="false" :destroyOnClose="true" :width="780"
             @ok="handleFileCheckOk">
      <div class="tree-select">
        <file-tree @on-file-check="handleFileCheck"
                   :checked-keys="form.checkedFiles" :torrent-data="torrentData"></file-tree>
      </div>
    </a-modal>


  </div>
</template>

<script>
import PathTree from "./PathTree";
import FileTree from "./FileTree";
import byteSize from 'byte-size'

import {getSpace, uploadTorrent, getTorrentInfo, taskExists, taskCreate} from "../http/api";

export default {
  name: "TaskModal",
  components: {
    PathTree,
    FileTree
  },
  watch: {
    'taskType'(val) {
      if (val === 'uri') {
        this.form.fileList = []
        this.$refs.torrentFile.clearValidate()
      } else {
        this.form.magentUri = ''
        this.$refs.magentUri.clearValidate()
      }
      this.infoHash = ''
      this.downloadType = 'all'
      this.cleanTorrentData()
    },
    'infoHash'(val) {
      if (val !== this.torrentData.info_hash) {
        getTorrentInfo({hash: val}).then(res => {
          let {data} = res
          if (data) {
            this.torrentData = data
          } else {
            this.cleanTorrentData()
          }
        })
      }
    }
  },
  computed: {
    freeSpace() {
      return byteSize(this.pathFreeSpace)
    }
  },
  methods: {
    show() {
      this.visible = true
    },
    handleTaskOk() {
      this.$refs.ruleForm.validate(valid => {
        if (valid) {
          this.createTaskLoading = true
          let param = {
            task_type: this.taskType,
            download_files: this.form.checkedFiles,
            download_path: this.form.downloadPath,
            download_all: false
          }
          if (this.form.checkedFiles.length > 0) {
            param.download = true
          } else {
            param.download = this.download
            if (param.download) {
              param.download_all = true
            }
          }
          if (this.taskType === 'uri') {
            param.create_torrent_info = this.form.magentUri
          } else {
            param.create_torrent_info = this.infoHash
          }
          taskCreate(param).then(res => {
            if (!res.status) {
              this.$message.error(res.message)
            } else {
              this.cleanModal()
              this.visible = false
            }
          }).finally(() => {
            this.createTaskLoading = false
          })
        } else {
          return false;
        }
      });
    },
    handlePathSelectOk() {
      this.form.downloadPath = this.pathSelect.tempPath
      this.pathSelect.visible = false
      this.validateSpaceData()
    },
    handlePathSelect(path) {
      this.pathSelect.tempPath = path
    },
    handleTorrentSelect(file) {
      let form = new FormData()
      form.append('torrent', file)
      uploadTorrent(form).then(res => {
        let {data} = res
        if (data) {
          this.torrentData = data
          this.infoHash = data.info_hash
        } else {
          this.cleanTorrentData()
        }
      }).then(() => {
        if (this.infoHash) {
          taskExists({hash: this.infoHash}).then(res => {
            let {data} = res
            if (data) {
              this.form.fileList[0].status = 'error'
              this.form.fileList[0].response = `任务已存在`
              this.infoHash = ''
              this.cleanTorrentData()
            }
          })
        }
      })
      return false
    },
    handleFileCheck(keys) {
      this.fileCheck.tempCheckedKeys = keys
    },
    handleFileCheckOk() {
      this.form.checkedFiles = this.fileCheck.tempCheckedKeys
      this.fileCheck.visible = false
    },
    handleTorrentChange({fileList}) {
      this.form.fileList = fileList;
      this.$refs.torrentFile.onFieldChange()
    },
    handleDownloadNowChange(checked) {
      this.download = checked
    },

    cleanTorrentData() {
      this.torrentData.info_hash = ''
      this.torrentData.name = ''
      this.torrentData.files = []
    },
    cleanModal() {
      this.cleanTorrentData()
      this.$refs.ruleForm.resetFields()
      this.fileList = []
      this.taskType = 'uri'
      this.downloadType = 'all'
      this.pathFreeSpace = 0
    },

    validateSpaceData() {
      this.$refs.downloadPath.onFieldBlur()
    }
  },
  data() {
    let validateMagentUri = (rule, value, callback) => {
      if (!value) {
        callback(new Error('请输入磁力链接'));
      }
      let pattern = /magnet:\?xt=urn:btih:([a-zA-Z0-9]{40}).*/
      if (!value.match(pattern)) {
        callback(new Error('无效磁力链接'))
      } else {
        let hash = value.match(pattern)[1].toLowerCase()
        taskExists({hash}).then(res => {
          let {data} = res
          if (!data) {
            this.infoHash = hash
            callback()
          } else {
            callback(new Error(`任务 ${hash} 已经存在`))
          }
        })
      }
    }

    let validateDownloadPath = (rule, value, callback) => {
      if (value && value === this.pathSelect.lastPath) {
        callback()
      }
      this.pathSelect.lastPath = value
      this.pathFreeSpace = 0
      if (value === '') {
        callback(new Error('下载路径不能为空'));
      } else {
        getSpace({path: this.form.downloadPath}).then(res => {
          if (res.status) {
            let {data} = res
            if (data) {
              this.pathFreeSpace = data
              callback()
            } else {
              callback(new Error('所选下载路径剩余空间不足'));
            }
          } else {
            callback(new Error('所选下载路径无效'));
          }
        })
      }
    }

    let validateTorrentFile = (rule, value, callback) => {
      if (!this.form.fileList.length) {
        callback(new Error('请上传种子文件'))
      } else {
        if (this.form.fileList[0].status === 'error') {
          callback(new Error(this.form.fileList[0].response))
        } else {
          callback()
        }
      }
    }

    return {
      pathSelect: {
        visible: false,
        tempPath: '',
        lastPath: ''
      },
      fileCheck: {
        visible: false,
        tempCheckedKeys: [],
      },
      visible: false,
      infoHash: '',
      pathFreeSpace: 0,
      taskType: 'uri',
      downloadType: 'all',
      download: true,
      modalLayout: {
        labelCol: {
          span: 4
        },
        wrapperCol: {
          span: 20
        }
      },
      form: {
        magentUri: '',
        downloadPath: '',
        fileList: [],
        checkedFiles: []
      },
      rules: {
        magentUri: [{validator: validateMagentUri, trigger: 'change'}],
        downloadPath: [{validator: validateDownloadPath, trigger: 'blur'}],
        torrentFile: [{validator: validateTorrentFile, trigger: 'change'}],
      },
      torrentData: {
        info_hash: '',
        name: '',
        files: []
      },
      createTaskLoading: false
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