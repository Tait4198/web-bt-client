// let socket = new WebSocket(`ws://${location.host}/ws/conn`);
let socket = new WebSocket(`ws://127.0.0.1:8080/ws/conn`);

export default function (bus) {
    socket.onopen = () =>{
    }
    socket.onclose = () =>{
        alert('Websocket Close')
    }
    socket.onmessage = function (message) {
        bus.emit('ws-message', message)
    }
    return socket
}
