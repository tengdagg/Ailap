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

export function history(type = 'recent') {
  return request.get('/logs/history', { params: { type } })
}

export function toggleFavorite(id) {
  return request.post(`/logs/history/${id}/favorite`)
}

export function updateNote(id, note) {
  return request.put(`/logs/history/${id}/note`, { note })
}

export function deleteHistory(id) {
  return request.delete(`/logs/history/${id}`)
}

export function inspect(params) {
  return request.get('/logs/inspect', { params })
}



