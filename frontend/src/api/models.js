import request from './request'

export function listModels(params) {
  return request.get('/models', { params })
}

export function createModel(data) {
  return request.post('/models', data)
}

export function updateModel(id, data) {
  return request.put(`/models/${id}`, data)
}

export function deleteModel(id) {
  return request.delete(`/models/${id}`)
}



