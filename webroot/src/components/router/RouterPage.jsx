import React, {Component} from 'react';
import {inject, observer} from 'mobx-react';
import {Toolbar, ToolbarGroup, ToolbarSeparator, ToolbarTitle} from 'material-ui/Toolbar';
import {Table, TableBody, TableHeader, TableHeaderColumn, TableRow, TableRowColumn} from 'material-ui/Table';
import RaisedButton from 'material-ui/RaisedButton';
import IconButton from 'material-ui/IconButton';
import {blue500} from 'material-ui/styles/colors';
import SvgIconKeyboardArrowUp from 'material-ui/svg-icons/hardware/keyboard-arrow-up';
import SvgIconKeyboardArrowDown from 'material-ui/svg-icons/hardware/keyboard-arrow-down';
import Title from './../Title';
import RouteModel from './RouteModel';

const styles = {
    reorder: {
        textAlign: "right"
    }
}

@inject('routing')
@observer
class RouterPage extends Component {

    constructor(props) {
        super(props);
        this.state = {
            open: false,
            value: "nexmo"
        };
    }

    handleOpenRouteModel = () => {
        this.setState({open: true});
    };

    handleCloseRouteModel = () => {
        this.setState({open: false});
    };

    render() {

        const { location, push, goBack } = this.props.routing;

        return (
            <div>
                <Title title="Router" />

                <p>Manage routes</p>

                <Toolbar>
                    <ToolbarGroup firstChild={true}>
                    </ToolbarGroup>
                    <ToolbarGroup lastChild={true}>
                        <RaisedButton
                            label="Delete"
                            secondary={true}
                            style={{marginRight: 0}}
                        />
                        <RaisedButton
                            label="Create"
                            primary={true}
                            onTouchTap={this.handleOpenRouteModel}
                        />
                    </ToolbarGroup>
                </Toolbar>

                <Table multiSelectable={true}>
                    <TableHeader>
                        <TableRow>
                            <TableHeaderColumn>NAME</TableHeaderColumn>
                            <TableHeaderColumn>PATTERN</TableHeaderColumn>
                            <TableHeaderColumn>PROVIDER</TableHeaderColumn>
                            <TableHeaderColumn>STATUS</TableHeaderColumn>
                            <TableHeaderColumn>REORDER</TableHeaderColumn>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        <TableRow>
                            <TableRowColumn>Taiwan</TableRowColumn>
                            <TableRowColumn>\+886</TableRowColumn>
                            <TableRowColumn>nexmo</TableRowColumn>
                            <TableRowColumn>active</TableRowColumn>
                            <TableRowColumn style={styles.reorder}>
                                <IconButton>
                                    <SvgIconKeyboardArrowDown color={blue500} />
                                </IconButton>
                            </TableRowColumn>
                        </TableRow>
                        <TableRow>
                            <TableRowColumn>USA</TableRowColumn>
                            <TableRowColumn>\+1</TableRowColumn>
                            <TableRowColumn>nexmo</TableRowColumn>
                            <TableRowColumn>active</TableRowColumn>
                            <TableRowColumn style={styles.reorder}>
                                <IconButton>
                                    <SvgIconKeyboardArrowUp color={blue500} />
                                </IconButton>
                                <IconButton>
                                    <SvgIconKeyboardArrowDown color={blue500} />
                                </IconButton>
                            </TableRowColumn>
                        </TableRow>
                        <TableRow>
                            <TableRowColumn>Japan</TableRowColumn>
                            <TableRowColumn>\+81</TableRowColumn>
                            <TableRowColumn>nexmo</TableRowColumn>
                            <TableRowColumn>active</TableRowColumn>
                            <TableRowColumn style={styles.reorder}>
                                <IconButton style={{zIndex: 9999}}>
                                    <SvgIconKeyboardArrowUp color={blue500} />
                                </IconButton>
                            </TableRowColumn>
                        </TableRow>
                    </TableBody>
                </Table>

                <RouteModel open={this.state.open} handleClose={this.handleCloseRouteModel} />

            </div>
        );
    }
}

export default RouterPage;