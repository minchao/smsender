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
                const rows = [];
                json.data.map(message => rows.push(MessageModel.fromJS(this, message)));
                this.messages = rows;
                this.since = null;
                this.until = null;
                if (json.paging.hasOwnProperty('previous')) {
                    this.since = json.paging.previous;
                }
                if (json.paging.hasOwnProperty('next')) {
                    this.until = json.paging.next;
                }
            })
    }

    @action add(message) {
        this.messages.push(MessageModel.fromJS(this, message));
    }

    @action clear() {
        this.messages = [];
        this.since = null;
        this.until = null;
    }

    search(to, status, since, until, limit) {
        this.sync(this.buildQueryString(to, status, since, until, limit));
    }

    buildQueryString(to = '', status = '', since = '', until = '', limit = 20) {
        let query = '';
        query += andWhere(query, 'to', to.replace('+', '%2B'));
        query += andWhere(query, 'status', status);
        query += andWhere(query, 'since', since);
        query += andWhere(query, 'until', until);
        query += andWhere(query, 'limit', limit);

        return '?' + query;
    }
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
