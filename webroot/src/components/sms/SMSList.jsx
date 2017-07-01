import React, { Component } from 'react'
import { inject, observer } from 'mobx-react'
import { Link } from 'react-router'
import { Table, TableBody, TableHeader, TableHeaderColumn, TableRow, TableRowColumn } from 'material-ui/Table'
import RaisedButton from 'material-ui/RaisedButton'

@inject('routing')
@observer
export default class SMSList extends Component {
  constructor (props) {
    super(props)
    this.store = this.props.store
    this.push = this.props.routing.push
  }

  pagingPrev = () => {
    const since = this.store.since
    this.push('/console/sms' + since.substr(since.indexOf('?')))
  }

  pagingNext = () => {
    const until = this.store.until
    this.push('/console/sms' + until.substr(until.indexOf('?')))
  }

  render () {
    return (
      <div>
        <Table>
          <TableHeader adjustForCheckbox={false} displaySelectAll={false}>
            <TableRow>
              <TableHeaderColumn>MESSAGE ID</TableHeaderColumn>
              <TableHeaderColumn>TO</TableHeaderColumn>
              <TableHeaderColumn>ROUTE</TableHeaderColumn>
              <TableHeaderColumn>STATUS</TableHeaderColumn>
              <TableHeaderColumn>DATE</TableHeaderColumn>
            </TableRow>
          </TableHeader>
          <TableBody displayRowCheckbox={false}>
            {(this.store.messages.length === 0)
              ? (
                <TableRow>
                  <TableRowColumn>No data</TableRowColumn>
                </TableRow>
              )
              : this.store.messages.map((message) => (
                <TableRow key={message.id}>
                  <TableRowColumn>
                    <Link to={`/console/sms/${message.id}/details`}>{message.id}</Link>
                  </TableRowColumn>
                  <TableRowColumn>{message.to}</TableRowColumn>
                  <TableRowColumn>{message.route}</TableRowColumn>
                  <TableRowColumn>{message.status}</TableRowColumn>
                  <TableRowColumn>{message.created_time}</TableRowColumn>
                </TableRow>
              ))}
          </TableBody>
        </Table>

        <div style={{marginTop: 20, textAlign: 'center'}}>
          {this.store.since && <RaisedButton label="Prev" onTouchTap={this.pagingPrev} />}
          {this.store.until && <RaisedButton label="Next" onTouchTap={this.pagingNext} />}
        </div>
      </div>
    )
  }
}
