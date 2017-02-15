import {observable} from 'mobx'

export default class MessageModel {
    store
    id
    to
    route
    status
    created_time

    constructor(store, id, to, route, status, created_time) {
        this.store = store
        this.id = id
        this.to = to
        this.route = route
        this.status = status
        this.created_time = created_time
    }

    toJS() {
        return {
            id: this.id,
            to: this.to,
            route: this.route,
            status: this.status,
            created_time: this.created_time
        }
    }

    static fromJS(store, object) {
        return new MessageModel(store, object.id, object.to, object.route, object.status, object.created_time)
    }
}
