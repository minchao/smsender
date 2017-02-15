import {observable} from 'mobx'

export default class RouteModel {
    store
    @observable name
    @observable pattern
    @observable provider
    @observable is_active

    constructor(store, name, pattern, provider, is_active) {
        this.store = store
        this.name = name
        this.pattern = pattern
        this.provider = provider
        this.is_active = is_active
    }

    toJS() {
        return {
            name: this.name,
            pattern: this.pattern,
            provider: this.provider,
            is_active: this.is_active
        }
    }

    static fromJS(store, object) {
        return new RouteModel(store, object.name, object.pattern, object.provider, object.is_active)
    }
}
