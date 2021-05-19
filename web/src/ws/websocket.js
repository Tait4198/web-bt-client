let socket = new WebSocket('ws://127.0.0.1:8080/ws/conn');

export default function (bus) {
    socket.onmessage = function (message) {
        bus.emit('ws-message', message)
    }
    return socket
}
