import Vue from 'vue'
import axios from 'axios'
import { VueAxios } from './axios'
import message from 'ant-design-vue/es/message'
import notification from 'ant-design-vue/es/notification'
import { ACCESS_TOKEN } from '@/store/global/mutation-types'
import router from '@/router'
import store from '@/store'
import i18n from '@/lang'

const service = axios.create({
  baseURL: process.env.VUE_APP_API_BASE_URL, // api base_url
  timeout: 6000, // default request timeout
  withCredentials: true,
  crossDomain: true,
})

const err = (error) => {
  console.log(error)
  const reg = /5\d{2}/g
  if (error.response && reg.test(error.response.status)) {
    const errorMsg = ((error.response || {}).data || {}).message || i18n.t('requestServiceError')
    message.error(errorMsg)
  } else if (error.response.status === 404 && error.config.url.includes('ci_types')) {
    message.warning(i18n.t('requestContact'))
  } else if (error.response.status === 412) {
    let seconds = 5
    notification.warning({
      key: 'notification',
      message: 'WARNING',
      description: i18n.t('requestWait', {
        time: 5,
      }),
      duration: 5,
    })
    let interval = setInterval(() => {
      seconds -= 1
      if (seconds === 0) {
        clearInterval(interval)
        interval = null
        return
      }
      notification.warning({
        key: 'notification',
        message: 'WARNING',
        description: i18n.t('requestWait', {
          time: seconds,
        }),
        duration: seconds,
      })
    }, 1000)
  } else if (error.config.url === '/api/v0.1/ci_types/can_define_computed' || error.config.isShowMessage === false) {
  } else {
    const errorMsg = ((error.response || {}).data || {}).message || i18n.t('requestError')
    message.error(`${errorMsg}`)
  }
  if (error.response) {
    console.log(error.config.url)
    if (error.response.status === 401 && router.path === '/user/login') {
      window.location.href = '/user/logout'
    }
  }
  return Promise.reject(error)
}

// request interceptor
service.interceptors.request.use((config) => {
  const token = Vue.ls.get(ACCESS_TOKEN)
  if (token) {
    config.headers['Access-Token'] = token // each request carries a custom token
  }
  config.headers['Accept-Language'] = store?.state?.locale ?? 'zh'
  return config
}, err)

// response interceptor
service.interceptors.response.use((response) => {
  return response.data
}, err)

const installer = {
  vm: {},
  install(Vue) {
    Vue.use(VueAxios, service)
  },
}

export { installer as VueAxios, service as axios }
