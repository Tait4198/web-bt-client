module.exports = {
    "presets": [["@babel/preset-env", {"modules": false}]],
    "plugins": [
        [
            "component",
            {
                "libraryName": "element-ui",
                "styleLibraryName": "theme-chalk"
            }
        ],
        [
            'import',
            {
                libraryName: '@byte-design/vue-icons',
                libraryDirectory: 'icons',
                style: () => '@byte-design/vue-icons/runtime/index.css',
                camel2DashComponentName: false,
                customName: name => '@byte-design/vue-icons/icons/' + name.slice(9)
            }
        ]
    ]
}
