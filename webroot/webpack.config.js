var path = require('path');
var webpack = require('webpack');

const NPM_TARGET = process.env.npm_lifecycle_event;

var DEV = false;

if (NPM_TARGET === 'dev') {
    DEV = true;
}

var config = {
    devtool: 'eval',
    entry: [
        './src/index'
    ],
    output: {
        path: path.join(__dirname, 'public'),
        filename: 'bundle.js',
        publicPath: '/static/'
    },
    plugins: [],
    resolve: {
        extensions: ['.js', '.jsx']
    },
    module: {
        loaders: [{
            test: /\.jsx?$/,
            loaders: ['babel-loader'],
            include: path.join(__dirname, 'src')
        }]
    }
};

// Development mode configuration
if (DEV) {
    config.entry = [
      'react-hot-loader/patch',
      'webpack-dev-server/client?http://localhost:3000',
      'webpack/hot/only-dev-server',
      './src/index'
    ]
    config.plugins.push(new webpack.HotModuleReplacementPlugin())
}

module.exports = config;
