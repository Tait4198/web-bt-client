<template>
  <a-directory-tree :load-data="loadPathData" :tree-data="treeData" @select="nodeSelect">
  </a-directory-tree>
</template>

<script>
import {getPath} from "../http/api";

export default {
  name: "PathTree",
  mounted() {
    getPath().then(res => {
      let {data} = res
      this.treeData = data
    })
  },
  data() {
    return {
      treeData: []
    }
  },
  methods: {
    nodeSelect(selectedKeys) {
      this.$emit('on-path-select', selectedKeys[0])
    },
    loadPathData(treeNode) {
      return new Promise(resolve => {
        if (treeNode.dataRef.children) {
          resolve()
          return
        }
        getPath({parent: treeNode.dataRef.key}).then(res => {
          let {data} = res
          treeNode.dataRef.children = data
          this.treeData = [...this.treeData];
          resolve();
        })
      })
    }
  }
}
</script>

<style scoped>

</style>
