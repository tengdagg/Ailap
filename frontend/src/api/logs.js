import request from './request'

export function queryLogs(params) {
  return request.get('/logs/query', { params })
}

export function suggestions(params) {
  return request.get('/logs/suggestions', { params })
}

export function labelValues(params) {
  return request.get('/logs/label-values', { params })
}

export function history() {
  return request.get('/logs/history')
}

export function inspect(params) {
  return request.get('/logs/inspect', { params })
}



