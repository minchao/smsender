import React, {Component} from 'react'
import {inject, observer} from 'mobx-react'
import {action, observable} from 'mobx'
import SyntaxHighlighter from 'react-syntax-highlighter'
import {agate} from 'react-syntax-highlighter/dist/styles'

import {getAPI} from '../../utils'
import MessageModel from '../../models/MessageModel'

@observer
export default class DetailsPage extends Component {
    message = new MessageModel()

    constructor(props) {
        super(props)
        this.fetch = this.fetch.bind(this)
    }

    componentDidMount() {
        this.fetch()
    }

    fetch() {
        fetch(getAPI('/api/messages/byIds?ids=' + this.props.params.messageId), {method: 'get'})
            .then(response => {
                if (!response.ok) throw new Error(response.statusText)
                return response.json()
            })
            .then(json => {
                if (json.data.length) {
                    this.message.fromJS(json.data[0])
                }
            })
    }

    render() {
        return (
            <div>
                <h2>Message details</h2>

                <h3>Properties</h3>

                <table className="sms-table">
                    <tbody>
                        <tr>
                            <th>Message ID</th>
                            <td>{this.message.id}</td>
                            <th>Route</th>
                            <td>{this.message.route}</td>
                        </tr>
                        <tr>
                            <th>From</th>
                            <td>{this.message.form}</td>
                            <th>Provider</th>
                            <td>{this.message.provider}</td>
                        </tr>
                        <tr>
                            <th>To</th>
                            <td>{this.message.to}</td>
                        </tr>
                        <tr>
                            <th>Body</th>
                            <td colSpan="3">
                                <pre style={{margin: 0, padding: 10, backgroundColor: '#e8e8e8'}}>
                                    {this.message.body}
                                </pre>
                            </td>
                        </tr>
                        <tr>
                            <th>Status</th>
                            <td>{this.message.status}</td>
                        </tr>
                        <tr>
                            <th>Created Time</th>
                            <td>{this.message.created_time}</td>
                        </tr>
                        <tr>
                            <th>Original Message ID</th>
                            <td>{this.message.original_message_id}</td>
                        </tr>
                    </tbody>
                </table>

                <h3>JSON</h3>
                <SyntaxHighlighter
                    language='json'
                    wrapLines={true}
                    style={agate}
                >{
                    this.message.json ? JSON.stringify(this.message.json, null, 4) : 'null'
                }</SyntaxHighlighter>
            </div>
        )
    }
}
