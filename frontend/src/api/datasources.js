import request from './request'

export function listDataSources(params) {
  return request.get('/datasources', { params })
}

export function createDataSource(data) {
  return request.post('/datasources', data)
}

export function updateDataSource(id, data) {
  return request.put(`/datasources/${id}`, data)
}

export function deleteDataSource(id) {
  return request.delete(`/datasources/${id}`)
}

export function testConnectionPayload(data) {
  return request.post('/datasources/test', data)
}

export function testConnection(id) {
  return request.post(`/datasources/${id}/test`)
}

