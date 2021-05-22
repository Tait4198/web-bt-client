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
