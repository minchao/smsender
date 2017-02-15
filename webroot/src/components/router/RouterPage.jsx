import React, {Component} from 'react';
import {inject, observer} from 'mobx-react';
import {action, observable} from 'mobx';
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
};

@inject('routing')
@observer
export default class RouterPage extends Component {

    static defaultProps = {
        store: new RouteStore()
    };

    @observable isOpen = false;
    @observable route = {
        isNew: true,
        name: '',
        pattern: '',
        provider: '',
        is_active: false
    };
    @observable selected = [];

    constructor(props) {
        super(props);
        this.openRouteDialog = this.openRouteDialog.bind(this);
        this.closeRouteDialog = this.closeRouteDialog.bind(this);
        this.setRoute = this.setRoute.bind(this);
        this.createRoute = this.createRoute.bind(this);
        this.updateRoute = this.updateRoute.bind(this);
        this.deleteRoute = this.deleteRoute.bind(this);
        this.reorderUp = this.reorderUp.bind(this);
        this.reorderDown = this.reorderDown.bind(this);
    }

    componentDidMount() {
        this.props.store.sync()
    }

    @action openRouteDialog() {
        this.isOpen = true;
    }

    @action closeRouteDialog() {
        this.isOpen = false;
    }

    @action setRoute(route) {
        if (route) {
            this.route.isNew = false;
            this.route.name = route.name;
            this.route.pattern = route.pattern;
            this.route.provider = route.provider;
            this.route.is_active = route.is_active;
        } else {
            this.route.isNew = true;
            this.route.name = '';
            this.route.pattern = '';
            this.route.provider = '';
            this.route.is_active = false;
        }
    }

    createRoute() {
        this.setRoute(null);
        this.openRouteDialog();
    }

    updateRoute(e) {
        e.preventDefault();
        this.setRoute(this.props.store.getByName(e.target.name));
        this.openRouteDialog();
    }

    deleteRoute() {
        if (this.selected[0] != undefined) {
            const route = this.props.store.routes[this.selected[0]];
            if (route) {
                this.props.store.del(route.name);
            }
        }
    }

    reorderUp(index) {
        this.props.store.reorder(index, 1, index-1);
    }

    reorderDown(index) {
        this.props.store.reorder(index, 1, index+2);
    }

    render() {
        const hasRoutes = this.props.store.routes.length != 0;

        return (
            <div>
                <h2>Router</h2>

                <p>Manage routes</p>

                <Toolbar>
                    <ToolbarGroup firstChild={true}>
                    </ToolbarGroup>
                    <ToolbarGroup lastChild={true}>
                        <RaisedButton
                            label="Delete"
                            secondary={true}
                            style={{marginRight: 0}}
                            onTouchTap={this.deleteRoute}
                        />
                        <RaisedButton
                            label="Create"
                            primary={true}
                            onTouchTap={this.createRoute}
                        />
                    </ToolbarGroup>
                </Toolbar>

                <Table
                    multiSelectable={false}
                    onRowSelection={(rows) => {this.selected = rows}}
                >
                    <TableHeader displaySelectAll={false}>
                        <TableRow>
                            <TableHeaderColumn>NAME</TableHeaderColumn>
                            <TableHeaderColumn>PATTERN</TableHeaderColumn>
                            <TableHeaderColumn>PROVIDER</TableHeaderColumn>
                            <TableHeaderColumn>STATUS</TableHeaderColumn>
                            <TableHeaderColumn>REORDER</TableHeaderColumn>
                        </TableRow>
                    </TableHeader>
                    <TableBody
                        displayRowCheckbox={hasRoutes}
                        deselectOnClickaway={true}
                    >
                        {(!hasRoutes)
                            ?
                            (
                                <TableRow>
                                    <TableRowColumn>No data</TableRowColumn>
                                </TableRow>
                            )
                            :
                            this.props.store.routes.map((route, i) => (
                            <TableRow key={i}>
                                <TableRowColumn>
                                    <a
                                        href="#"
                                        name={route.name}
                                        onClick={this.updateRoute}
                                    >
                                    {route.name}
                                    </a>
                                </TableRowColumn>
                                <TableRowColumn>{route.pattern}</TableRowColumn>
                                <TableRowColumn>{route.provider}</TableRowColumn>
                                <TableRowColumn>{route.is_active ? "enable": "disable"}</TableRowColumn>
                                <TableRowColumn style={styles.reorder}>
                                    <IconButton>
                                        {i == 0
                                            ? null
                                            :
                                            <SvgIconKeyboardArrowUp
                                                color={blue500}
                                                onClick={() => this.reorderUp(i)}
                                            />}
                                    </IconButton>
                                    <IconButton>
                                        {i == this.props.store.routes.length - 1
                                            ? null
                                            :
                                            <SvgIconKeyboardArrowDown
                                                color={blue500}
                                                onClick={() => this.reorderDown(i)}
                                            />
                                        }
                                    </IconButton>
                                </TableRowColumn>

                            </TableRow>
                        ))}
                    </TableBody>
                </Table>

                <RouteDialog
                    isOpen={this.isOpen}
                    store={this.props.store}
                    route={this.route}
                    closeRouteDialog={this.closeRouteDialog}
                />
            </div>
        );
    }
}
