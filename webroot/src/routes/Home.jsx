import React, {Component} from 'react';
import AppBar from 'material-ui/AppBar';

class Home extends Component {

    render() {
        return (
            <div>
                <AppBar title="Hello, 世界" showMenuIconButton={false} />
            </div>
        )
    }
}

export default Home;
