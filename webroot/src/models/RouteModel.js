import { action, observable } from 'mobx'

export default class RouteModel {
  @observable name
  @observable pattern
  @observable provider
  @observable from
  @observable is_active // eslint-disable-line

  @action set (key, value) {
    this[key] = value
  }

  @action fromJS (object) {
    this.name = object.name
    this.pattern = object.pattern
    this.provider = object.provider
    this.from = object.from
    this.is_active = object.is_active

    return this
  }

  toJS () {
    return {
      name: this.name,
      pattern: this.pattern,
      provider: this.provider,
      from: this.from,
      is_active: this.is_active
    }
  }

  static new (object) {
    return (new RouteModel().fromJS(object))
  }
}
