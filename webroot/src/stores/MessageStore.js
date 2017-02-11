import {action, observable, computed, reaction} from 'mobx';

import {getAPI} from '../utils';
import MessageModel from '../models/MessageModel';

export default class MessageStore {
    @observable messages = [];
    @observable since = null;
    @observable until = null;

    @action find(messageId = '') {
        fetch(getAPI('/api/messages/byIds?ids=' + messageId), {method: 'get'})
            .then(response => {
                if (!response.ok) throw new Error(response.statusText);
                return response.json();
            })
            .then(json => {
                this.clear();
                json.data.map(message => this.add(message));
            })
    }

    @action sync(query = '') {
        fetch(getAPI('/api/messages' + query), {method: 'get'})
            .then(response => {
                if (!response.ok) throw new Error(response.statusText);
                return response.json();
            })
            .then(json => {
                this.clear();
                json.data.map(message => this.add(message));
                if (json.paging.hasOwnProperty('previous')) {
                    this.since = getParameterByName(json.paging.previous, 'since');
                }
                if (json.paging.hasOwnProperty('next')) {
                    this.until = getParameterByName(json.paging.next, 'until');
                }
            })
    }

    @action search(to = '', status = '', since = null, until = null, limit = 10) {
        let query = '';
        query += andWhere(query, 'to', to.replace('+', '%2B'));
        query += andWhere(query, 'status', status);
        query += andWhere(query, 'since', since);
        query += andWhere(query, 'until', until);
        query += andWhere(query, 'limit', limit);

        this.sync('?' + query);
    }

    @action add(message) {
        this.messages.push(MessageModel.fromJS(this, message));
    }

    @action clear() {
        this.messages = [];
        this.since = null;
        this.until = null;
    }
}

function getParameterByName(url, name) {
    name = name.replace(/[\[\]]/g, '\$&');
    let regex = new RegExp('[?&]' + name + '(=([^&#]*)|&|#|$)'),
        results = regex.exec(url);
    if (!results) return null;
    if (!results[2]) return '';
    return decodeURIComponent(results[2].replace(/\+/g, ' '));
}

function andWhere(query, where, value) {
    if (!value) {
        return '';
    }

    where = where + '=' + value;

    if (query) {
        where = '&' + where;
    }
    return where;
}
