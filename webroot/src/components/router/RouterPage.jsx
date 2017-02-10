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
import RouteDialog from './RouteDialog';
import RouteStore from '../../stores/RouteStore';

const styles = {
    reorder: {
        textAlign: "right"
    }
}

@inject('routing')
@observer
export default class RouterPage extends Component {

    static defaultProps = {
        store: new RouteStore()
    }

    constructor(props) {
        super(props);
        this.state = {
            open: false,
            value: "nexmo"
        };
    }

    componentDidMount() {
        this.props.store.sync()
    }

    handleOpenRouteDialog = () => {
        this.setState({open: true});
    };

    handleCloseRouteDialog = () => {
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
                            onTouchTap={this.handleOpenRouteDialog}
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
                        {this.props.store.routes.map((route, i) => (
                            <TableRow key={i}>
                                <TableRowColumn>{route.name}</TableRowColumn>
                                <TableRowColumn>{route.pattern}</TableRowColumn>
                                <TableRowColumn>{route.provider}</TableRowColumn>
                                <TableRowColumn>{route.is_active ? "enable": "disable"}</TableRowColumn>
                                <TableRowColumn style={styles.reorder}>
                                    <IconButton>
                                        {i == 0 ? null : <SvgIconKeyboardArrowUp color={blue500} />}
                                    </IconButton>
                                    <IconButton>
                                        {i == this.props.store.routes.length - 1 ? null : <SvgIconKeyboardArrowDown color={blue500} />}
                                    </IconButton>
                                </TableRowColumn>

                            </TableRow>
                        ))}
                    </TableBody>
                </Table>

                <RouteDialog open={this.state.open} handleClose={this.handleCloseRouteDialog} />

            </div>
        );
    }
}
