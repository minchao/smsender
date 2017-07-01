import { action, observable } from 'mobx'

import api from './API'
import RouteModel from '../models/RouteModel'

export default class RouteStore {
  @observable routes = []
  @observable providers = []

  @action add (route) {
    this.routes.push(RouteModel.new(route))
  }

  @action clear () {
    this.routes = []
    this.providers = []
  }

  @action initData (json) {
    this.clear()
    json.data.map(route => this.add(route))
    json.providers.map(provider => this.providers.push(provider.name))
  }

  sync () {
    api.getRoutes((json) => {
      this.initData(json)
    })
  }

  getByName (name) {
    for (let i = 0; i < this.routes.length; i++) {
      if (this.routes[i].name === name) {
        return this.routes[i]
      }
    }
    return null
  }

  create (route) {
    api.postRoute(route, (json) => {
      this.sync()
    })
  }

  update (route) {
    api.putRoute(route, (json) => {
      this.sync()
    })
  }

  reorder (rangeStart, rangeLength, insertBefore) {
    api.reorderRoutes(rangeStart, rangeLength, insertBefore, (json) => {
      this.initData(json)
    })
  }

  del (name) {
    api.deleteRoute(name, (json) => {
      this.sync()
    })
  }
}
