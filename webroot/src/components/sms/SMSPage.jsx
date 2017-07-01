import React, { Component } from 'react'
import { inject, observer } from 'mobx-react'
import { action, observable } from 'mobx'
import { Toolbar, ToolbarGroup } from 'material-ui/Toolbar'
import DropDownMenu from 'material-ui/DropDownMenu'
import MenuItem from 'material-ui/MenuItem'
import TextField from 'material-ui/TextField'
import RaisedButton from 'material-ui/RaisedButton'

import MessageStore from '../../stores/MessageStore'
import SMSList from './SMSList'

const status = [
  {text: 'All Status', value: ''},
  {text: 'Accepted', value: 'accepted'},
  {text: 'Queued', value: 'queued'},
  {text: 'Sending', value: 'sending'},
  {text: 'Failed', value: 'failed'},
  {text: 'Sent', value: 'sent'},
  {text: 'Unknown', value: 'unknown'},
  {text: 'Undelivered', value: 'undelivered'},
  {text: 'Delivered', value: 'delivered'}
]

@inject('routing')
@observer
export default class SMSPage extends Component {
  static defaultProps = {
    store: new MessageStore()
  }

  @observable form = {
    id: '',
    to: '',
    status: '',
    since: '',
    until: '',
    limit: 20
  }

  constructor (props) {
    super(props)
    this.queryString = null
    this.push = this.props.routing.push
    this.setForm = this.setForm.bind(this)
    this.resetForm = this.resetForm.bind(this)
    this.updateFormProperty = this.updateFormProperty.bind(this)
    this.updateFormStatus = this.updateFormStatus.bind(this)
  }

  componentDidMount () {
    this.setForm()
    this.fetch()
  }

  componentDidUpdate (prevProps) {
    const queryString = this.props.routing.location.pathname + this.props.routing.location.search

    if (this.queryString !== queryString) {
      this.setForm()
      this.fetch()
    }
  }

  @action setForm () {
    this.resetForm()
    const query = this.props.routing.location.query
    if (query.id) this.form.id = query.id
    if (query.to) this.form.to = query.to
    if (query.status) this.form.status = query.status
    if (query.since) this.form.since = query.since
    if (query.until) this.form.until = query.until
    if (query.limit) this.form.limit = query.limit

    this.queryString = this.props.routing.location.pathname + this.props.routing.location.search
  }

  @action resetForm () {
    this.form.id = ''
    this.form.to = ''
    this.form.status = ''
    this.form.since = ''
    this.form.until = ''
    this.form.limit = 20
  }

  @action updateFormProperty (event, value) {
    this.form[event.target.name] = value
  }

  @action updateFormStatus (event, index) {
    this.form.status = status[index].value
  }

  fetch = () => {
    if (this.form.id) {
      this.props.store.find(this.form.id)
    } else {
      this.props.store.search(this.form.to, this.form.status, this.form.since, this.form.until, this.form.limit)
    }
  }

  find = () => {
    this.push('/console/sms?id=' + this.form.id)
  }

  search = () => {
    const query = this.props.store.buildQueryString(this.form.to, this.form.status, '', '', this.form.limit)
    this.push('/console/sms' + query)
  }

  render () {
    return (
      <div>
        <h2>SMS Delivery Logs</h2>

        <p>Search by message ID</p>

        <Toolbar>
          <ToolbarGroup firstChild style={{width: '100%'}}>
            <TextField
              name="id"
              hintText="Message ID: b29f66182b317var78gg"
              value={this.form.id}
              fullWidth
              style={{marginLeft: 20, width: '100%'}}
              onChange={this.updateFormProperty}
            />
          </ToolbarGroup>
          <ToolbarGroup lastChild>
            <RaisedButton
              label="Find"
              primary
              onTouchTap={this.find}
            />
          </ToolbarGroup>
        </Toolbar>

        <p>Search by recipient phone number</p>

        <Toolbar>
          <ToolbarGroup firstChild>
            <TextField
              name="to"
              hintText="To Phone Number: +886987654321"
              value={this.form.to}
              style={{marginLeft: 20}}
              onChange={this.updateFormProperty}
            />
          </ToolbarGroup>
          <ToolbarGroup lastChild>
            <DropDownMenu
              name="status"
              value={this.form.status}
              onChange={this.updateFormStatus}
            >
              {status.map((s, i) => (
                <MenuItem key={i} value={s.value} primaryText={s.text} />
              ))}
            </DropDownMenu>
            <RaisedButton
              label="Search"
              primary
              onTouchTap={this.search}
            />
          </ToolbarGroup>
        </Toolbar>

        <SMSList store={this.props.store} />
      </div>
    )
  }
}
