import React from 'react';
import {render} from 'react-dom';
import {AppContainer} from 'react-hot-loader';
import injectTapEventPlugin from 'react-tap-event-plugin';
import App from './App';

injectTapEventPlugin();

render(
    <AppContainer>
        <App />
    </AppContainer>,
    document.getElementById('root')
);

if (module.hot) {
    module.hot.accept('./App', () => {
        const NextApp = require('./App').default;

        render(
            <AppContainer>
                <NextApp />
            </AppContainer>,
            document.getElementById('root')
        );
    });
}
