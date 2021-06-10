function getSocketUrl() {
    if (process.env.VUE_APP_WS_URL) {
        return `${process.env.VUE_APP_WS_URL}/ws/conn`
    } else {
        return `ws://${location.host}/ws/conn`
    }
}

let socket = new WebSocket(getSocketUrl());

export default function (bus, store) {
    socket.onopen = () => {
        store.commit('wsConnStatus', true)
    }
    socket.onclose = () => {
        store.commit('wsConnStatus', false)
    }
    socket.onmessage = function (message) {
        bus.emit('ws-message', message)
    }
    return socket
}
