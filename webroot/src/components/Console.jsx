import React, {Component} from 'react';
import {inject, observer} from 'mobx-react';
import AppBar from 'material-ui/AppBar';
import Drawer from 'material-ui/Drawer';
import FlatButton from 'material-ui/FlatButton';
import Menu from 'material-ui/Menu';
import MenuItem from 'material-ui/MenuItem';
import {minBlack} from 'material-ui/styles/colors';
import SvgMessage from 'material-ui/svg-icons/communication/message';
import SvgList from 'material-ui/svg-icons/action/list';

@inject('routing')
@observer
class Routes extends Component {

    constructor(props) {
        super(props);
        this.menuItemStyle = this.menuItemStyle.bind(this)
    }

    menuItemStyle(targetPath) {
        if (this.props.routing.location.pathname == targetPath) {
            return {backgroundColor: minBlack};
        }
        return null;
    }

    render() {
        const {push} = this.props.routing;

        return (
            <div>
                <AppBar title="SMSender"
                        iconElementRight={<FlatButton label="Home" onTouchTap={() => push("/")} />}
                />

                <div style={{paddingLeft: 240, paddingRight: 40}}>
                    {this.props.children}
                </div>

                <Drawer
                    docked={true}
                    width={210}
                >
                    <AppBar title="SMSender" />
                    <Menu desktop={true}>
                        <MenuItem
                            onTouchTap={() => push("/console/sms")}
                            leftIcon={<SvgMessage />}
                        >SMS</MenuItem>
                        <MenuItem
                            onTouchTap={() => push("/console/sms/send")}
                            style={this.menuItemStyle('/console/sms/send')}
                            insetChildren={true}
                        >Send an SMS</MenuItem>
                        <MenuItem
                            onTouchTap={() => push("/console/sms")}
                            style={this.menuItemStyle('/console/sms')}
                            insetChildren={true}
                        >Delivery Logs</MenuItem>
                        <MenuItem
                            onTouchTap={() => push("/console/router")}
                            style={this.menuItemStyle('/console/router')}
                            leftIcon={<SvgList />}
                        >Router</MenuItem>
                    </Menu>
                </Drawer>
            </div>
        );
    }
}

export default Routes;