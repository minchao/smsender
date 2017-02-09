import React, {Component} from 'react';
import Dialog from 'material-ui/Dialog';
import TextField from 'material-ui/TextField';
import SelectField from 'material-ui/SelectField';
import MenuItem from 'material-ui/MenuItem';
import Toggle from 'material-ui/Toggle';
import FlatButton from 'material-ui/FlatButton';

export default class RouteModel extends Component {

    constructor(props) {
        super(props);
        this.state = {
            value: "nexmo"
        };
    }

    handleClose = () => {
        this.props.handleClose()
    };

    handleChange = (event, index, value) => this.setState({value});

    render() {

        const actions = [
            <FlatButton
                label="Cancel"
                onTouchTap={this.handleClose}
            />,
            <FlatButton
                label="Submit"
                primary={true}
                onTouchTap={this.handleClose}
            />,
        ];

        return (
            <Dialog
                title="Create a new Route"
                actions={actions}
                modal={true}
                open={this.props.open}
            >
                <TextField
                    hintText="Name"
                />
                <br />
                <TextField
                    hintText="Pattern"
                />
                <br />
                <SelectField
                    floatingLabelText="Provider"
                    value={this.state.value}
                    onChange={this.handleChange}
                >
                    <MenuItem value={"nexmo"} primaryText="nexmo" />
                    <MenuItem value={"twilio"} primaryText="twilio" />
                </SelectField>
                <br />
                <TextField
                    hintText="From"
                />
                <br />
                <br />
                <Toggle
                    label="Is Active"
                    defaultToggled={true}
                />
            </Dialog>
        );

    }
}