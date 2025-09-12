import request from './request'

export function queryLogs(params) {
  return request.get('/logs/query', { params })
}

export function suggestions(params) {
  return request.get('/logs/suggestions', { params })
}

