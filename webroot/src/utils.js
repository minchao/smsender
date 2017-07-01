let baseURL = ''

if (module.hot) {
  baseURL = 'http://localhost:8080'
}

export function getBaseURL () {
  return baseURL
}
