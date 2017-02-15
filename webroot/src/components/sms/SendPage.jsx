import React, {Component} from 'react'
import {inject, observer} from 'mobx-react'
import {action, observable} from 'mobx'
import Paper from 'material-ui/Paper'
import TextField from 'material-ui/TextField'
import RaisedButton from 'material-ui/RaisedButton'
import SyntaxHighlighter from 'react-syntax-highlighter'
import {agate} from 'react-syntax-highlighter/dist/styles'

import {getAPI} from '../../utils'
import MessageStore from '../../stores/MessageStore'

@observer
export default class SendPage extends Component {

    static defaultProps = {
        store: new MessageStore()
    }

    @observable message = {
        to: '',
        from: '',
        body: ''
    }
    @observable response = 'null'

    constructor(props) {
        super(props)
        this.updateProperty = this.updateProperty.bind(this)
        this.post = this.post.bind(this)
        this.setResponse = this.setResponse.bind(this)
        this.reset = this.reset.bind(this)
    }

    componentDidMount() {
        this.reset()
    }

    @action updateProperty(event, value) {
        this.message[event.target.name] = value
    }

    @action setResponse(text) {
        this.response = text
    }

    @action reset() {
        this.message.to = ''
        this.message.from = ''
        this.message.body = ''
        this.response = 'null'
    }

    post() {
        fetch(getAPI('/api/messages'), {
                method: 'post',
                body: JSON.stringify({
                    'to': [this.message.to],
                    'from': this.message.from,
                    'body': this.message.body
                }),
                headers: new Headers({'Content-Type': 'application/json'})
            })
            .then(response => {
                return response.json()
            })
            .then(json => {
                this.setResponse(JSON.stringify(json, null, 4))
            })
    }

    render() {
        return (
            <div>
                <h2>Send an SMS</h2>

                <p>Request</p>

                <Paper style={{padding: 30}}>
                    <TextField
                        name="to"
                        hintText="To phone number (E.164 format): +886987654321"
                        value={this.message.to}
                        style={{width: '50%'}}
                        onChange={this.updateProperty}
                    />
                    <br />
                    <TextField
                        name="from"
                        hintText="Sender Id (phone number or alphanumeric)"
                        value={this.message.from}
                        style={{width: '50%'}}
                        onChange={this.updateProperty}
                    />
                    <br />
                    <TextField
                        name="body"
                        hintText="The text of the message"
                        value={this.message.body}
                        multiLine={true}
                        rows={2}
                        rowsMax={4}
                        style={{width: '50%'}}
                        onChange={this.updateProperty}
                    />
                    <br />
                    <div style={{textAlign: 'right'}}>
                        <RaisedButton
                            label="Send"
                            primary={true}
                            onTouchTap={this.post}
                        />
                    </div>
                </Paper>

                <p>Response</p>

                <div>
                    <SyntaxHighlighter
                        language='json'
                        wrapLines={true}
                        style={agate}
                    >{this.response}</SyntaxHighlighter>
                </div>
            </div>
        )
    }
}