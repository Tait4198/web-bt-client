import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
    state: {
        wsConnStatus: false
    },
    getters: {
        getWsConnStatus: (state) => {
            return state.wsConnStatus
        }
    },
    mutations: {
        wsConnStatus(state, val) {
            state.wsConnStatus = val
        }
    },
})