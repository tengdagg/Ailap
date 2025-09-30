import request from './request'

export function analyzeLogs(payload) {
  return request.post('/ai/analyze-logs', payload, { timeout: 60000 })
}


