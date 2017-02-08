import React from 'react';

export default class Title extends React.Component {
    render() {
        return <h2
            style={{
                marginTop: 40,
                marginBottom: 20,
                paddingBottom: 10,
                fontWeight: 400,
                fontSize: 28,
                color: "#000",
                borderBottom: "1px solid #eee"
            }}
        >{this.props.title}</h2>;
    }
}