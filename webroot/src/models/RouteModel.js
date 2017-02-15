import {observable} from 'mobx'

export default class RouteModel {
    @observable name
    @observable pattern
    @observable provider
    @observable is_active

    constructor(name, pattern, provider, is_active) {
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

    fromJS(object) {
        this.name = object.name
        this.pattern = object.pattern
        this.provider = object.provider
        this.is_active = object.is_active

        return this
    }

    static new(object) {
        return (new RouteModel().fromJS(object))
    }
}
