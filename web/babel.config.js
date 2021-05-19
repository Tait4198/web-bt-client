module.exports = {
    "presets": [["@babel/preset-env", {"modules": false}]],
    "plugins": [
        ["import", {"libraryName": "ant-design-vue", "libraryDirectory": "es", "style": "css"}]
    ]
}
