import {observable, computed, reaction} from 'mobx';

import {getAPI} from '../utils';
import RouteModel from '../models/RouteModel';

export default class RouteStore {
    @observable routes = [];

    sync() {
        fetch(getAPI('/api/routes'), {method: 'get'})
            .then(response => {
                if (!response.ok) throw new Error(response.statusText)
                return response.json()
            })
            .then(json => {
                this.routes = [];
                json.data.map(route => this.add(route));
            })
    }

    add(route) {
        this.routes.push(RouteModel.fromJS(this, route))
    }

    toJS() {
        return this.routes.map(route => route.toJS())
    }

    static fromJS(array) {
        const store = new RouteStore();
        store.routes = array.map(route => store.add(route));
        return store;
    }
}