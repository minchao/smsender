class API {
  constructor () {
    this.baseURL = API_HOST
  }

  fetch (url, object, callback) {
    fetch(this.baseURL + url, object)
      .then(response => {
        if (response.ok) return response.json()
        return response.json()
      })
      .then(json => {
        return callback(json)
      })
      .catch(() => {
        return callback(null)
      })
  }

  getRoutes (callback) {
    this.fetch('/api/routes', {method: 'get'}, callback)
  }

  postRoute (route, callback) {
    this.fetch('/api/routes', {
      method: 'post',
      body: JSON.stringify(route),
      headers: new Headers({'Content-Type': 'application/json'})
    }, callback)
  }

  putRoute (route, callback) {
    this.fetch('/api/routes/' + route.name, {
      method: 'put',
      body: JSON.stringify(route),
      headers: new Headers({'Content-Type': 'application/json'})
    }, callback)
  }

  reorderRoutes (rangeStart, rangeLength, insertBefore, callback) {
    this.fetch('/api/routes', {
      method: 'put',
      body: JSON.stringify({
        'range_start': rangeStart,
        'range_length': rangeLength,
        'insert_before': insertBefore
      }),
      headers: new Headers({'Content-Type': 'application/json'})
    }, callback)
  }

  deleteRoute (name, callback) {
    this.fetch('/api/routes/' + name, {method: 'delete'}, callback)
  }

  postMessage (to, from, body, callback) {
    this.fetch('/api/messages', {
      method: 'post',
      body: JSON.stringify({
        'to': [to],
        'from': from,
        'body': body
      }),
      headers: new Headers({'Content-Type': 'application/json'})
    }, callback)
  }

  getMessages (query, callback) {
    this.fetch('/api/messages' + query, {method: 'get'}, callback)
  }

  getMessagesByIds (messageIds, callback) {
    this.fetch('/api/messages/byIds?ids=' + messageIds, {method: 'get'}, callback)
  }
}

const api = new API()

export default api
export { API }
