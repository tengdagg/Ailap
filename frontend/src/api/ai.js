import request from './request'

export function analyzeLogs(payload, config = {}) {
  return request.post('/ai/analyze-logs', payload, { timeout: 60000, ...config })
}


