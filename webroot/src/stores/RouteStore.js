import {action, observable, computed, reaction} from 'mobx';

import {getAPI} from '../utils';
import RouteModel from '../models/RouteModel';

export default class RouteStore {
    @observable routes = [];
    @observable providers = [];

    @computed get count() {
        return this.routes.length;
    }

    @action sync() {
        fetch(getAPI('/api/routes'), {method: 'get'})
            .then(response => {
                if (!response.ok) throw new Error(response.statusText);
                return response.json();
            })
            .then(json => {
                this.clear();
                json.data.map(route => this.add(route));
                json.providers.map(provider => this.providers.push(provider.name));
            })
    }

    @action add(route) {
        this.routes.push(RouteModel.fromJS(this, route));
    }

    @action getByName(name) {
        for (let i = 0; i < this.routes.length; i++) {
            if (this.routes[i].name == name) {
                return this.routes[i];
            }
        }
        return null;
    }

    @action create(route) {

    }

    @action update(route) {

    }

    @action reorder(rangeStart, rangeLength, insertBefore) {

    }

    @action del(name) {

    }

    @action clear() {
        this.routes = [];
    }

    toJS() {
        return this.routes.map(route => route.toJS());
    }

    static fromJS(array) {
        const store = new RouteStore();
        store.routes = array.map(route => store.add(route));
        return store;
    }
}
