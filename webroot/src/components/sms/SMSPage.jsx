import React, {Component} from 'react';
import {inject, observer} from 'mobx-react';
import {action, observable} from 'mobx';
import {Toolbar, ToolbarGroup, ToolbarSeparator, ToolbarTitle} from 'material-ui/Toolbar';
import {Table, TableBody, TableHeader, TableHeaderColumn, TableRow, TableRowColumn} from 'material-ui/Table';
import DropDownMenu from 'material-ui/DropDownMenu';
import MenuItem from 'material-ui/MenuItem';
import TextField from 'material-ui/TextField';
import RaisedButton from 'material-ui/RaisedButton';

import Title from './../Title';
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

@observer
export default class SMSPage extends Component {

    static defaultProps = {
        store: new MessageStore()
    };

    @observable form = {
        id: '',
        to: '',
        status: '',
        limit: 20
    };

    constructor(props) {
        super(props);
        this.updateProperty = this.updateProperty.bind(this);
        this.updateFilterStatus = this.updateFilterStatus.bind(this);
    }

    componentDidMount() {
        this.props.store.search(this.form.to, this.form.status, null, null, this.form.limit);
    };

    @action updateProperty(event, value) {
        this.form[event.target.name] = value;
    };

    @action updateFilterStatus(event, index) {
        this.form.status = status[index].value;
    };

    find = () => {
        this.props.store.find(this.form.id);
    };

    filter = () => {
        this.props.store.search(this.form.to, this.form.status, null, null, this.form.limit);
    };

    pagingPrev = () => {
        this.props.store.search(this.form.to, this.form.status, this.props.store.since, null, this.form.limit);
    };

    pagingNext = () => {
        this.props.store.search(this.form.to, this.form.status, null, this.props.store.until, this.form.limit);
    };

    render() {
        return (
            <div>
                <Title title="SMS" />

                <p>Messages delivery logs</p>

                <Toolbar>
                    <ToolbarGroup firstChild={true} style={{width: "100%"}}>
                        <TextField
                            name="id"
                            hintText="Search by Message ID: b29f66182b317var78gg"
                            value={this.form.id}
                            fullWidth={true}
                            style={{marginLeft: 20, width: "100%"}}
                            onChange={this.updateProperty}
                        />
                    </ToolbarGroup>
                    <ToolbarGroup lastChild={true}>
                        <RaisedButton
                            label="Search"
                            primary={true}
                            onTouchTap={this.find}
                        />
                    </ToolbarGroup>
                </Toolbar>

                <br />

                <Toolbar>
                    <ToolbarGroup firstChild={true}>
                        <TextField
                            name="to"
                            hintText="To Phone Number: +886987654321"
                            value={this.form.to}
                            style={{marginLeft: 20}}
                            onChange={this.updateProperty}
                        />
                    </ToolbarGroup>
                    <ToolbarGroup lastChild={true}>
                        <DropDownMenu
                            name="status"
                            value={this.form.status}
                            onChange={this.updateFilterStatus}
                        >
                            {status.map((s, i) => (
                                <MenuItem key={i} value={s.value} primaryText={s.text} />
                            ))}
                        </DropDownMenu>
                        <RaisedButton
                            label="Filter"
                            primary={true}
                            onTouchTap={this.filter}
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
