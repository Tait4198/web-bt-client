import Vue from 'vue'
import {
    Button,
    Col,
    Row,
    Space,
    Card,
    List,
    Progress
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

Vue.config.productionTip = false

Vue.use(VueBus);
let app = new Vue({
    render: h => h(App),
})

ws(app.$bus)

app.$mount('#app')

