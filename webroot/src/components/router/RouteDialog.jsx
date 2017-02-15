import React, {Component} from 'react'
import {inject, observer} from "mobx-react"
import {action, observable} from 'mobx'
import Dialog from 'material-ui/Dialog'
import TextField from 'material-ui/TextField'
import SelectField from 'material-ui/SelectField'
import MenuItem from 'material-ui/MenuItem'
import Toggle from 'material-ui/Toggle'
import FlatButton from 'material-ui/FlatButton'

@observer
export default class RouteDialog extends Component {

    constructor(props) {
        super(props)
        this.store = props.store
        this.route = props.route
        this.updateProperty = this.updateProperty.bind(this)
        this.updateProvider = this.updateProvider.bind(this)
        this.cancel = this.cancel.bind(this)
        this.submit = this.submit.bind(this)
    }

    updateProperty(event, value) {
        this.route[event.target.name] = value
    }

    updateProvider(event, index, value) {
        this.route.provider = value
    }

    cancel() {
        this.props.closeRouteDialog()
    }

    submit() {
        if (this.props.isNew) {
            this.store.create(this.route)
        } else {
            this.store.update(this.route)
        }
        this.props.closeRouteDialog()
    }

    render() {
        const actions = [
            <FlatButton
                label="Cancel"
                onTouchTap={this.cancel}
            />,
            <FlatButton
                label="Submit"
                primary={true}
                onTouchTap={this.submit}
            />,
        ]

        return (
            <Dialog
                title={(this.props.isNew ? 'Create a new' : 'Update a') + ' Route'}
                actions={actions}
                modal={true}
                open={this.props.isOpen}
            >
                <TextField
                    name="name"
                    hintText="Name"
                    value={this.route.name}
                    disabled={!this.props.isNew}
                    onChange={this.updateProperty}
                />
                <br />
                <TextField
                    name="pattern"
                    hintText="Pattern (Regular Expression)"
                    value={this.route.pattern}
                    onChange={this.updateProperty}
                />
                <br />
                <SelectField
                    floatingLabelText="Provider"
                    value={this.route.provider}
                    onChange={this.updateProvider}
                >
                    {this.store.providers.map((provider, i) => (
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
                    hintText="From (Sender ID)"
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
        )
    }
}