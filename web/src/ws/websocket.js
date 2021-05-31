let socket = new WebSocket(`ws://${location.host}/ws/conn`);

export default function (bus) {
    socket.onmessage = function (message) {
        bus.emit('ws-message', message)
    }
    return socket
}
