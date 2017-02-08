import React, {Component} from 'react';
import {inject, observer} from 'mobx-react';
import RaisedButton from 'material-ui/RaisedButton';

const styles = {
    home: {
        paddingTop: 200,
        paddingBottom: 60,
        textAlign: "center",
        backgroundColor: "#00bcd4",
        height: "100%",
        overflow: "hidden"
    },
    h1: {
        margin: 0,
        paddingBottom: 20,
        fontWeight: 300,
        fontSize: "50px",
        color: "#fff"
    },
    h2: {
        paddingBottom: 20,
        fontWeight: 300,
        fontSize: "24px",
        lineHeight: "32px",
        color: "#fff",
        webkitFontSmoothing: "antialiased"
    }
};

@inject('routing')
@observer
class Home extends Component {

    render() {
        const { location, push, goBack } = this.props.routing;

        return (
            <div style={styles.home}>
                <h1 style={styles.h1}>SMSender Console</h1>
                <h2 style={styles.h2}>A SMS server written in Go</h2>
                <RaisedButton
                    label="Console"
                    labelColor="#00bcd4"
                    onTouchTap={() => push("console")}
                />
            </div>
        )
    }
}

export default Home;