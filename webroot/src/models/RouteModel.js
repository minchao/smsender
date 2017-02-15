import {observable} from 'mobx'

export default class RouteModel {
    @observable name
    @observable pattern
    @observable provider
    @observable from
    @observable is_active

    constructor(name = '', pattern = '', provider = '', from = '', is_active = false) {
        this.name = name
        this.pattern = pattern
        this.provider = provider
        this.from = from
        this.is_active = is_active
    }

    toJS() {
        return {
            name: this.name,
            pattern: this.pattern,
            provider: this.provider,
            from: this.from,
            is_active: this.is_active
        }
    }

    fromJS(object) {
        this.name = object.name
        this.pattern = object.pattern
        this.provider = object.provider
        this.from = object.from
        this.is_active = object.is_active

        return this
    }

    static new(object) {
        return (new RouteModel().fromJS(object))
    }
}
