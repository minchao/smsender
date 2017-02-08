import React, {Component} from 'react';
import {Provider, observer} from 'mobx-react';
import DevTools from 'mobx-react-devtools';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import lightBaseTheme from 'material-ui/styles/baseThemes/lightBaseTheme';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import {RouterStore, syncHistoryWithStore} from 'mobx-react-router';
import {Router, IndexRoute, Route, browserHistory} from 'react-router';
import Home from './routes/Home';

const routingStore = new RouterStore();

const stores = {
    routing: routingStore
};

const history = syncHistoryWithStore(browserHistory, routingStore);

class App extends Component {
    render() {
        return (
            <MuiThemeProvider muiTheme={getMuiTheme(lightBaseTheme)}>
                <div>
                    <Provider {...stores}>
                        <Router history={history}>
                            <Route path="/" component={Home} />
                        </Router>
                    </Provider>
                    <DevTools />
                </div>
            </MuiThemeProvider>
        )
    }
}

export default App;
