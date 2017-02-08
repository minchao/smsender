import React, {Component} from 'react';
import {inject, observer} from 'mobx-react';
import {Toolbar, ToolbarGroup, ToolbarSeparator, ToolbarTitle} from 'material-ui/Toolbar';
import {Table, TableBody, TableHeader, TableHeaderColumn, TableRow, TableRowColumn} from 'material-ui/Table';
import DropDownMenu from 'material-ui/DropDownMenu';
import MenuItem from 'material-ui/MenuItem';
import TextField from 'material-ui/TextField';
import RaisedButton from 'material-ui/RaisedButton';
import Title from './../Title';

@inject('routing')
@observer
class SMSPage extends Component {

    constructor(props) {
        super(props);
        this.state = {
            value: "",
        };
    }

    handleChange = (event, index, value) => this.setState({value});

    render() {
        const { location, push, goBack } = this.props.routing;

        return (
            <div>
                <Title title="SMS" />

                <p>Message delivery logs</p>

                <Toolbar>
                    <ToolbarGroup firstChild={true} style={{width: "100%"}}>
                        <TextField
                            hintText="Search by Message ID: b29f66182b317var78gg"
                            fullWidth={true}
                            style={{marginLeft: 20, width: "100%"}}
                        />
                    </ToolbarGroup>
                    <ToolbarGroup lastChild={true}>
                        <RaisedButton
                            label="Search"
                            primary={true}
                        />
                    </ToolbarGroup>
                </Toolbar>

                <br />

                <Toolbar>
                    <ToolbarGroup firstChild={true}>
                        <TextField
                            hintText="To Phone Number: +886987654321"
                            style={{marginLeft: 20}}
                        />
                    </ToolbarGroup>
                    <ToolbarGroup lastChild={true}>
                        <DropDownMenu value={this.state.value} onChange={this.handleChange}>
                            <MenuItem value={""} primaryText="All Status" />
                            <MenuItem value={"accepted"} primaryText="Accepted" />
                            <MenuItem value={"queued"} primaryText="Queued" />
                            <MenuItem value={"sending"} primaryText="Sending" />
                            <MenuItem value={"failed"} primaryText="Failed" />
                            <MenuItem value={"sent"} primaryText="Sent" />
                            <MenuItem value={"unknown"} primaryText="Unknown" />
                            <MenuItem value={"undelivered"} primaryText="Undelivered" />
                            <MenuItem value={"delivered"} primaryText="Delivered" />
                        </DropDownMenu>
                        <RaisedButton
                            label="Filter"
                            primary={true}
                            style={{}}
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
                        <TableRow>
                            <TableRowColumn>b29f66182b317var78gg</TableRowColumn>
                            <TableRowColumn>+886987654321</TableRowColumn>
                            <TableRowColumn>nexmo</TableRowColumn>
                            <TableRowColumn>sent</TableRowColumn>
                            <TableRowColumn>2017-02-02T16:51:36</TableRowColumn>
                        </TableRow>
                    </TableBody>
                </Table>
            </div>
        );
    }
}

export default SMSPage;