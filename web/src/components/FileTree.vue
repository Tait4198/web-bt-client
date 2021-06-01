<template>
  <a-tree :tree-data="treeData"
          :selectable="false"
          :default-checked-keys="defaultCheckedKeys"
          checkable @check="onCheck">
    <template slot="custom" slot-scope="item">
      <div>
        <span>{{ item.title }}</span>
        <span style="margin-left: 24px">{{ byteSize(item.length) }}</span>
      </div>
    </template>
    <template slot="detail" slot-scope="item">
      <a-space>
        <p class="file-name">{{ item.title }}</p>
        <span>{{ byteSize(item.length) }}</span>
        <span v-if="item.stats">
          <a-button type="link" v-if="item.stats.completed_pieces === item.stats.pieces"
                    class="download-link" @click="onDownload(item.key)">
            下载
          </a-button>
          <span v-else-if="item.stats.percent > 0">
          {{ `${(item.stats.percent * 100).toFixed(2)}%` }}
          </span>
          <span v-else>
            0%
          </span>
        </span>
      </a-space>
    </template>
  </a-tree>
</template>

<script>
import byteSize from 'byte-size'

export default {
  name: "FreeTree",
  props: {
    torrentData: {
      type: Object,
      default: () => {
        return {}
      }
    },
    defaultCheckedKeys: {
      type: Array,
      default: () => {
        return []
      }
    },
    disableCheckbox: {
      type: Boolean,
      default: false
    },
    itemSlot: {
      type: String,
      default: 'custom'
    }
  },
  computed: {
    treeData() {
      let map = new Map()
      for (let i = 0; i < this.torrentData.files.length; i++) {
        let file = this.torrentData.files[i]
        this.convertMap(map, '', file, 0)
      }
      let array = this.mapToArray(map)
      if (this.torrentData.files.length > 1) {
        return [{
          key: this.torrentData.name,
          title: this.torrentData.name,
          children: array,
          length: this.torrentData.length,
          isLeaf: false,
          disableCheckbox: this.disableCheckbox,
          scopedSlots: {
            title: this.itemSlot
          }
        }]
      } else {
        return array
      }
    }
  },
  methods: {
    onCheck(checkedKeys) {
      this.$emit('on-file-check', checkedKeys)
    },
    onDownload(key) {
      this.$emit('on-download', key)
    },
    byteSize(length) {
      return byteSize(length)
    },
    convertMap(parentMap, parentKey, file, depth) {
      if (depth >= file.path.length) {
        return
      }
      let dPath = file.path[depth]
      if (!parentMap.has(dPath)) {
        let newNode = {
          title: dPath,
          length: 0
        }
        if (parentKey === '') {
          newNode.key = dPath
        } else {
          newNode.key = parentKey + '/' + dPath
        }
        if (depth + 1 < file.path.length) {
          newNode.map = new Map()
          newNode.isLeaf = false
        } else {
          newNode.isLeaf = true
          newNode.stats = {
            length: file.length,
            pieces: file.pieces,
            bytes_completed: file.bytes_completed,
            completed_pieces: file.completed_pieces,
            percent: file.bytes_completed / file.length
          }
        }
        parentMap.set(dPath, newNode)
      }
      let node = parentMap.get(dPath)
      node.length += file.length
      this.convertMap(node.map, node.key, file, depth + 1)
    },
    mapToArray(map) {
      let array = []
      for (let value of map.values()) {
        let obj = {
          key: value.key,
          title: value.title,
          length: value.length,
          isLeaf: value.isLeaf,
          stats: value.stats,
          disableCheckbox: this.disableCheckbox,
          scopedSlots: {
            title: this.itemSlot
          }
        }
        if (value.map) {
          obj.children = this.mapToArray(value.map)
        }
        array.push(obj)
      }
      return array
    }
  }
}
</script>

<style scoped lang="less">
.download-link {
  padding: 0 !important;
}

.file-name {
  max-width: 400px;
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
  margin: 0 !important;
}
</style>
