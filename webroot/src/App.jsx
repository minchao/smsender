import React, {Component} from 'react';
import DevTools from 'mobx-react-devtools';

class App extends Component {
    render() {
        return (
            <div>
                <p>Hello, 世界</p>
                <DevTools />
            </div>
        );
    }
};

export default App;
