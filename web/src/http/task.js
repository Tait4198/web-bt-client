import {GET} from './http'

const baseURL = `http://localhost:8080`

export const getTaskList = (params) => GET('/task/list', params, {
    baseURL
})
