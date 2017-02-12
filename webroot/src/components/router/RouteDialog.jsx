import React, {Component} from 'react';
import {inject, observer} from "mobx-react";
import Dialog from 'material-ui/Dialog';
import TextField from 'material-ui/TextField';
import SelectField from 'material-ui/SelectField';
import MenuItem from 'material-ui/MenuItem';
import Toggle from 'material-ui/Toggle';
import FlatButton from 'material-ui/FlatButton';

@observer
export default class RouteDialog extends Component {

    constructor(props) {
        super(props);
        this.isOpen = props.isOpen;
        this.providers = props.providers;
        this.route = props.route;
        this.closeRouteDialog = this.closeRouteDialog.bind(this);
        this.updateProperty = this.updateProperty.bind(this);
        this.updateProvider = this.updateProvider.bind(this);
    }

    closeRouteDialog() {
        this.props.closeRouteDialog();
    };

    updateProperty(event, value) {
        this.route[event.target.name] = value;
    }

    updateProvider(event, index, value) {
        this.route.provider = value;
    }

    render() {
        const actions = [
            <FlatButton
                label="Cancel"
                onTouchTap={this.closeRouteDialog}
            />,
            <FlatButton
                label="Submit"
                primary={true}
                onTouchTap={this.closeRouteDialog}
            />,
        ];

        return (
            <Dialog
                title={(this.route.isNew ? 'Create a new' : 'Update a') + ' Route'}
                actions={actions}
                modal={true}
                open={this.props.isOpen}
            >
                <TextField
                    name="name"
                    hintText="Name"
                    value={this.route.name}
                    onChange={this.updateProperty}
                />
                <br />
                <TextField
                    name="pattern"
                    hintText="Pattern"
                    value={this.route.pattern}
                    onChange={this.updateProperty}
                />
                <br />
                <SelectField
                    floatingLabelText="Provider"
                    value={this.route.provider}
                    onChange={this.updateProvider}
                >
                    {this.providers.map((provider, i) => (
                        <MenuItem
                            key={i}
                            value={provider}
                            primaryText={provider}
                        />
                    ))}
                </SelectField>
                <br />
                <TextField
                    name="from"
                    hintText="From"
                    value={this.route.from}
                    onChange={this.updateProperty}
                />
                <br />
                <br />
                <Toggle
                    name="is_active"
                    label="Is Active"
                    defaultToggled={this.route.is_active}
                    onToggle={this.updateProperty}
                />
            </Dialog>
        );
    }
}