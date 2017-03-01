import {action, observable, computed, reaction} from 'mobx'

import api from './API'
import MessageModel from '../models/MessageModel'

export default class MessageStore {
    @observable messages = []
    @observable since = null
    @observable until = null

    @action clear() {
        this.messages = []
        this.since = null
        this.until = null
    }

    @action initData(json) {
        this.clear()

        const rows = []
        json.data.map(message => rows.push(MessageModel.new(message)))
        this.messages = rows
        if (json.hasOwnProperty('paging')) {
            if (json.paging.hasOwnProperty('previous')) {
                this.since = json.paging.previous
            }
            if (json.paging.hasOwnProperty('next')) {
                this.until = json.paging.next
            }
        }
    }

    find(messageId = '') {
        api.getMessagesByIds(messageId, (json) => {
            this.initData(json)
        })
    }

    sync(query = '') {
        api.getMessages(query, (json) => {
            this.initData(json)
        })
    }

    search(to, status, since, until, limit) {
        this.sync(this.buildQueryString(to, status, since, until, limit))
    }

    buildQueryString(to = '', status = '', since = '', until = '', limit = 20) {
        let query = ''
        query += andWhere(query, 'to', to.replace('+', '%2B'))
        query += andWhere(query, 'status', status)
        query += andWhere(query, 'since', since)
        query += andWhere(query, 'until', until)
        query += andWhere(query, 'limit', limit)

        return '?' + query
    }
}

function andWhere(query, where, value) {
    if (!value) {
        return ''
    }

    where = where + '=' + value

    if (query) {
        where = '&' + where
    }
    return where
}
