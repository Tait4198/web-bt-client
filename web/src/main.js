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
    FormModel,
    Input,
    Badge,
    Spin,
    message,
    notification
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
Vue.component(Input.name, Input)
Vue.component(Badge.name, Badge)
Vue.component(Spin.name, Spin)

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

