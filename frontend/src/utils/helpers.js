export function formatDate(date) {
  const d = date instanceof Date ? date : new Date(date)
  return d.toISOString()
}
















