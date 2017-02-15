import {action, observable, computed, reaction} from 'mobx'

import {getAPI} from '../utils'
import RouteModel from '../models/RouteModel'

export default class RouteStore {
    @observable routes = []
    @observable providers = []

    @action sync() {
        fetch(getAPI('/api/routes'), {method: 'get'})
            .then(response => {
                if (!response.ok) throw new Error(response.statusText)
                return response.json()
            })
            .then(json => {
                this.initData(json)
            })
    }

    @action add(route) {
        this.routes.push(RouteModel.new(route))
    }

    @action getByName(name) {
        for (let i = 0; i < this.routes.length; i++) {
            if (this.routes[i].name == name) {
                return this.routes[i]
            }
        }
        return null
    }

    @action create(route) {
        fetch(getAPI('/api/routes'), {
                method: 'post',
                body: JSON.stringify(route),
                headers: new Headers({'Content-Type': 'application/json'})
            })
            .then(response => {
                if (!response.ok) throw new Error(response.statusText)
                this.sync()
            })
    }

    @action update(route) {
        fetch(getAPI('/api/routes/' + route.name), {
                method: 'put',
                body: JSON.stringify(route),
                headers: new Headers({'Content-Type': 'application/json'})
            })
            .then(response => {
                if (!response.ok) throw new Error(response.statusText)
                this.sync()
            })
    }

    @action reorder(rangeStart, rangeLength, insertBefore) {
        fetch(getAPI('/api/routes'), {
                method: 'put',
                body: JSON.stringify({
                    'range_start': rangeStart,
                    'range_length': rangeLength,
                    'insert_before': insertBefore
                }),
                headers: new Headers({'Content-Type': 'application/json'})
            })
            .then(response => {
                if (!response.ok) throw new Error(response.statusText)
                return response.json()
            })
            .then(json => {
                this.initData(json)
            })
    }

    @action del(name) {
        fetch(getAPI('/api/routes/' + name), {method: 'delete'})
            .then(response => {
                if (!response.ok) throw new Error(response.statusText)
                this.sync()
            })
    }

    @action clear() {
        this.routes = []
        this.providers = []
    }

    @action initData(json) {
        this.clear()
        json.data.map(route => this.add(route))
        json.providers.map(provider => this.providers.push(provider.name))
    }

    toJS() {
        return this.routes.map(route => route.toJS())
    }
}
