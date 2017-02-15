import React, {Component} from 'react';
import {inject, observer} from 'mobx-react';
import {action, observable} from 'mobx';
import {Toolbar, ToolbarGroup, ToolbarSeparator, ToolbarTitle} from 'material-ui/Toolbar';
import {Table, TableBody, TableHeader, TableHeaderColumn, TableRow, TableRowColumn} from 'material-ui/Table';
import DropDownMenu from 'material-ui/DropDownMenu';
import MenuItem from 'material-ui/MenuItem';
import TextField from 'material-ui/TextField';
import RaisedButton from 'material-ui/RaisedButton';

import MessageStore from '../../stores/MessageStore';

const status = [
    {text: 'All Status', value: ''},
    {text: 'Accepted', value: 'accepted'},
    {text: 'Queued', value: 'queued'},
    {text: 'Sending', value: 'sending'},
    {text: 'Failed', value: 'failed'},
    {text: 'Sent', value: 'sent'},
    {text: 'Unknown', value: 'unknown'},
    {text: 'Undelivered', value: 'undelivered'},
    {text: 'Delivered', value: 'delivered'},
];

@inject('routing')
@observer
export default class SMSPage extends Component {

    static defaultProps = {
        store: new MessageStore()
    };

    @observable form = {
        id: '',
        to: '',
        status: '',
        since: '',
        until: '',
        limit: 20
    };

    constructor(props) {
        super(props);
        this.queryString = null;
        this.push = this.props.routing.push;
        this.setForm = this.setForm.bind(this);
        this.resetForm = this.resetForm.bind(this);
        this.updateFormProperty = this.updateFormProperty.bind(this);
        this.updateFormStatus = this.updateFormStatus.bind(this);
    }

    componentDidMount() {
        this.setForm();
        this.fetch();
    };

    componentDidUpdate(prevProps) {
        const queryString = this.props.routing.location.pathname + this.props.routing.location.search;

        if (this.queryString != queryString) {
            this.setForm();
            this.fetch();
        }
    }

    @action setForm() {
        this.resetForm();
        const query = this.props.routing.location.query;
        if (query.id) this.form.id = query.id;
        if (query.to) this.form.to = query.to;
        if (query.status) this.form.status = query.status;
        if (query.since) this.form.since = query.since;
        if (query.until) this.form.until = query.until;
        if (query.limit) this.form.limit = query.limit;

        this.queryString = this.props.routing.location.pathname + this.props.routing.location.search;
    }

    @action resetForm() {
        this.form.id = '';
        this.form.to = '';
        this.form.status = '';
        this.form.since = '';
        this.form.until = '';
        this.form.limit = 20;
    }

    @action updateFormProperty(event, value) {
        this.form[event.target.name] = value;
    };

    @action updateFormStatus(event, index) {
        this.form.status = status[index].value;
    };

    fetch = () => {
        if (this.form.id) {
            this.props.store.find(this.form.id);
        } else {
            this.props.store.search(this.form.to, this.form.status, this.form.since, this.form.until, this.form.limit);
        }
    }

    find = () => {
        this.push('/console/sms?id=' + this.form.id);
    };

    search = () => {
        const query = this.props.store.buildQueryString(this.form.to, this.form.status, '', '', this.form.limit);
        this.push('/console/sms' + query);
    }

    pagingPrev = () => {
        const since = this.props.store.since;
        this.push('/console/sms' + since.substr(since.indexOf('?')));
    };

    pagingNext = () => {
        const until = this.props.store.until;
        this.push('/console/sms' + until.substr(until.indexOf('?')));
    };

    render() {
        return (
            <div>
                <h2>SMS Delivery Logs</h2>

                <p>Search by message ID</p>

                <Toolbar>
                    <ToolbarGroup firstChild={true} style={{width: "100%"}}>
                        <TextField
                            name="id"
                            hintText="Message ID: b29f66182b317var78gg"
                            value={this.form.id}
                            fullWidth={true}
                            style={{marginLeft: 20, width: "100%"}}
                            onChange={this.updateFormProperty}
                        />
                    </ToolbarGroup>
                    <ToolbarGroup lastChild={true}>
                        <RaisedButton
                            label="Find"
                            primary={true}
                            onTouchTap={this.find}
                        />
                    </ToolbarGroup>
                </Toolbar>

                <p>Search by recipient phone number</p>

                <Toolbar>
                    <ToolbarGroup firstChild={true}>
                        <TextField
                            name="to"
                            hintText="To Phone Number: +886987654321"
                            value={this.form.to}
                            style={{marginLeft: 20}}
                            onChange={this.updateFormProperty}
                        />
                    </ToolbarGroup>
                    <ToolbarGroup lastChild={true}>
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
                            primary={true}
                            onTouchTap={this.search}
                        />
                    </ToolbarGroup>
                </Toolbar>

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
                        {(this.props.store.messages.length == 0)
                            ?
                            (
                                <TableRow>
                                    <TableRowColumn>No data</TableRowColumn>
                                </TableRow>
                            )
                            :
                            this.props.store.messages.map((message, i) => (
                            <TableRow key={i}>
                                <TableRowColumn>{message.id}</TableRowColumn>
                                <TableRowColumn>{message.to}</TableRowColumn>
                                <TableRowColumn>{message.route}</TableRowColumn>
                                <TableRowColumn>{message.status}</TableRowColumn>
                                <TableRowColumn>{message.created_time}</TableRowColumn>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>

                <div style={{marginTop: 20, textAlign: "center"}}>
                    {this.props.store.since == null
                        ? null
                        : <RaisedButton label="Prev" onTouchTap={this.pagingPrev} />}
                    {this.props.store.until == null
                        ? null
                        : <RaisedButton label="Next" onTouchTap={this.pagingNext} />}
                </div>
            </div>
        );
    }
}
