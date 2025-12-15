import request from '@/api/request'

export function listMonitors() {
    return request.get('/monitors')
}

export function createMonitor(data) {
    return request.post('/monitors', data)
}

export function updateMonitor(id, data) {
    return request.put(`/monitors/${id}`, data)
}

export function deleteMonitor(id) {
    return request.delete(`/monitors/${id}`)
}

export function getMonitor(id) {
    return request.get(`/monitors/${id}`)
}

export function listAlerts() {
    return request.get('/monitors/alerts')
}

export function getAlert(id) {
    return request.get(`/monitors/alerts/${id}`)
}

// Channels
export function listChannels() {
    return request.get('/channels')
}

export function createChannel(data) {
    return request.post('/channels', data)
}

export function updateChannel(id, data) {
    return request.put(`/channels/${id}`, data)
}

export function deleteChannel(id) {
    return request.delete(`/channels/${id}`)
}

export function getChannel(id) {
    return request.get(`/channels/${id}`)
}

export function testChannel(data) {
    return request.post('/channels/test', data)
}
