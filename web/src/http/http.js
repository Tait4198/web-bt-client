import axios from 'axios'

axios.interceptors.request.use(
    config => {
        return config
    },
    error => {
        return Promise.reject(error)
    }
)

axios.interceptors.response.use(function (response) {
    return response.data
}, function (error) {
    return Promise.reject(error)
})


export const GET = (url, params, options = {}) => {
    return axios.get(`${url}`, Object.assign({
        params
    }, options)).then(data => data)
}

export const POST = (url, params, options) => {
    return axios.post(`${url}`, params, options).then(data => data)
}
