import request from './request'

export function login(payload) {
  return request.post('/auth/login', payload)
}

export function logout() {
  return request.post('/auth/logout')
}

export function profile() {
  return request.get('/auth/profile')
}

export function changePassword(data) {
  return request.post('/auth/change-password', data)
}

export function updateProfile(data) {
  return request.post('/auth/update-profile', data)
}










