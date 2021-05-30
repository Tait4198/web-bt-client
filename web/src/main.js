import Vue from 'vue'
import {
    Button,
    Col,
    Row,
    Space,
    Card,
    List,
    Progress,
    Icon,
    Modal,
    Form,
    FormModel,
    Input,
    Badge,
    Spin,
    Tree,
    Radio,
    Upload,
    Tooltip,
    Switch,
    Dropdown,
    Menu,
    message,
    notification,
} from 'ant-design-vue'
import App from './App.vue';
import VueBus from 'vue-bus';
import ws from "./ws/websocket";

Vue.component(Button.name, Button)
Vue.component(Col.name, Col)
Vue.component(Row.name, Row)
Vue.component(Space.name, Space)
Vue.component(List.name, List)
Vue.component(List.Item.name, List.Item)
Vue.component(Progress.name, Progress)
Vue.component(Card.name, Card)
Vue.component(Icon.name, Icon)
Vue.component(FormModel.name, FormModel)
Vue.component(FormModel.Item.name, FormModel.Item)
Vue.component(Form.name, Form)
Vue.component(Form.Item.name, Form.Item)
Vue.component(Input.name, Input)
Vue.component(Badge.name, Badge)
Vue.component(Spin.name, Spin)
Vue.component(Tree.name, Tree)
Vue.component(Tree.DirectoryTree.name, Tree.DirectoryTree)
Vue.component(Radio.name, Radio)
Vue.component(Radio.Group.name, Radio.Group)
Vue.component(Radio.Button.name, Radio.Button)
Vue.component(Upload.name, Upload)
Vue.component(Tooltip.name, Tooltip)
Vue.component(Switch.name, Switch)
Vue.component(Dropdown.name, Dropdown)
Vue.component(Dropdown.Button.name, Dropdown.Button)
Vue.component(Menu.name, Menu)
Vue.component(Menu.Item.name, Menu.Item)

Vue.config.productionTip = false

Vue.prototype.$message = message;
Vue.prototype.$notification = notification;

Vue.use(Modal);
Vue.use(VueBus);
let app = new Vue({
    render: h => h(App),
})

ws(app.$bus)

app.$mount('#app')

