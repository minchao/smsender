const path = require('path')
const webpack = require('webpack')

module.exports = () => {
  const env = process.env.NODE_ENV
  const ifDev = plugin => (env === 'development') ? plugin : undefined
  const removeEmpty = array => array.filter(p => !!p)

  return {
    devtool: ifDev('source-map'),
    entry: {
      main: removeEmpty([
        ifDev('react-hot-loader/patch'),
        ifDev(`webpack-dev-server/client?http://localhost:3000`),
        ifDev('webpack/hot/only-dev-server'),
        path.join(__dirname, './src/index')
      ])
    },
    resolve: {
      extensions: ['.js', '.jsx']
    },
    output: {
      filename: 'bundle.js',
      path: path.resolve(__dirname, 'public'),
      publicPath: '/static/'
    },
    module: {
      rules: [
        {
          test: /\.jsx?$/,
          exclude: /node_modules/,
          loader: ['babel-loader']
        }
      ]
    },
    plugins: removeEmpty([
      ifDev(new webpack.HotModuleReplacementPlugin())
    ])
  }
}
