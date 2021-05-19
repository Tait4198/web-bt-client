import Vue from 'vue'
import {
    Button,
    ButtonGroup,
    Divider,
} from 'element-ui';
import App from './App.vue';

Vue.config.productionTip = false

Vue.component(Button.name, Button);
Vue.component(ButtonGroup.name, ButtonGroup);
Vue.component(Divider.name, Divider);


new Vue({
    render: h => h(App),
}).$mount('#app')
