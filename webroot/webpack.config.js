const path = require('path')
const webpack = require('webpack')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const ExtractTextPlugin = require('extract-text-webpack-plugin')

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
        path.join(__dirname, './src/index.jsx')
      ])
    },
    resolve: {
      extensions: ['.js', '.jsx']
    },
    output: {
      filename: '[name].[hash].js',
      sourceMapFilename: '[name].[hash].map.js',
      path: path.resolve(__dirname, 'dist'),
      publicPath: '/dist/'
    },
    module: {
      rules: [
        {
          test: /\.jsx?$/,
          exclude: /node_modules/,
          loader: ['babel-loader']
        },
        {
          test: /\.css$/,
          use: ExtractTextPlugin.extract({
            use: 'css-loader',
          })
        }
      ]
    },
    plugins: removeEmpty([
      new HtmlWebpackPlugin({
        template: path.resolve(__dirname, './src/index.html'),
        filename: 'index.html',
        inject: 'body'
      }),
      new ExtractTextPlugin({
        filename: '[name].[hash].css',
      }),
      ifDev(new webpack.HotModuleReplacementPlugin()),
      ifDev(new webpack.NamedModulesPlugin())
    ])
  }
}
