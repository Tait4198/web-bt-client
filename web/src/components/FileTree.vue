<template>
  <a-tree :tree-data="treeData" :default-checked-keys="checkedKeys" checkable @check="onCheck">
    <template slot="custom" slot-scope="item">
      <div>
        <span>{{ item.title }}</span>
        <span style="margin-left: 24px">{{ fileSize(item.length) }}</span>
      </div>
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
    checkedKeys: {
      type: Array,
      default: () => {
        return []
      }
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
          scopedSlots: {
            title: 'custom'
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
    fileSize(length) {
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
          scopedSlots: {
            title: 'custom'
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

<style scoped>

</style>