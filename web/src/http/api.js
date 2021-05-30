import {GET, POST} from './http'

const baseURL = `http://127.0.0.1:8080`

export const getTaskList = (params) => GET('/task/list', params, {
    baseURL
})

export const taskStart = (params) => POST('/task/start', params, {
    baseURL
})

export const taskStop = (params) => GET('/task/stop', params, {
    baseURL
})

export const taskDelete = (params) => GET('/task/delete', params, {
    baseURL
})

export const taskRestart = (params) => POST('/task/restart', params, {
    baseURL
})

export const taskCreate = (params) => POST('/task/create', params, {
    baseURL
})

export const taskExists = (params) => GET('/task/exists', params, {
    baseURL
})

export const getPath = (params) => GET('/base/path', params, {
    baseURL
})

export const getSpace = (params) => GET('/base/space', params, {
    baseURL
})

export const getTorrentInfo = (params) => GET('/torrent/info', params, {
    baseURL
})

export const uploadTorrent = (params) => POST('/torrent/upload', params, {
    baseURL
})
