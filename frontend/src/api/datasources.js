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

export function testConnection(id) {
  // backend supports both /:id/test and /test with payload; we use id route
  return request.post(`/datasources/${id}/test`)
}

export function testConnectionPayload(data) {
  // use payload-based test before creating/saving
  return request.post('/datasources/test', data)
}

export async function getDataSourceById(id) {
  const { data } = await listDataSources()
  const items = data?.data?.items || []
  return items.find((it) => String(it.id) === String(id)) || null
}

