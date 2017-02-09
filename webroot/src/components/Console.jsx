import React, {Component} from 'react';
import {inject, observer} from 'mobx-react';
import AppBar from 'material-ui/AppBar';
import Drawer from 'material-ui/Drawer';
import FlatButton from 'material-ui/FlatButton';
import MenuItem from 'material-ui/MenuItem';
import ClearFix from 'material-ui/internal/ClearFix';

@inject('routing')
@observer
class Routes extends Component {

    render() {
        const { location, push, goBack } = this.props.routing;

        return (
            <div>
                <AppBar title="SMSender"
                        iconElementRight={<FlatButton label="Home" onTouchTap={() => push("/")} />}
                />

                <ClearFix style={{paddingLeft: 240, paddingRight: 40}}>
                    {this.props.children}
                </ClearFix>

                <Drawer
                    docked={true}
                    width={200}
                >
                    <AppBar title="SMSender" />
                    <MenuItem onTouchTap={() => push("/console/sms")}>SMS</MenuItem>
                    <MenuItem onTouchTap={() => push("/console/router")}>Router</MenuItem>
                </Drawer>
            </div>
        );
    }
}

export default Routes;