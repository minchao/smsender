const path = require('path')
const webpack = require('webpack')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const ExtractTextPlugin = require('extract-text-webpack-plugin')

module.exports = () => {
  const env = process.env.NODE_ENV
  const ifProd = plugin => (env === 'production') ? plugin : undefined
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
            use: 'css-loader'
          })
        }
      ]
    },
    plugins: removeEmpty([
      new webpack.DefinePlugin({
        'process.env.NODE_ENV': JSON.stringify(env),
        IS_DEV: Boolean(env === 'development'),
        API_HOST: JSON.stringify(env === 'development' ? 'http://localhost:8080' : '')
      }),
      new webpack.optimize.CommonsChunkPlugin({
        name: 'vendor',
        minChunks: function (module) {
          if (module.resource && (/^.*\.(css|scss)$/).test(module.resource)) {
            return false
          }
          return module.context && module.context.indexOf('node_modules') !== -1
        }
      }),
      new HtmlWebpackPlugin({
        template: path.resolve(__dirname, './src/index.html'),
        filename: 'index.html',
        inject: 'body'
      }),
      new ExtractTextPlugin({
        filename: '[name].[hash].css'
      }),
      ifProd(new webpack.optimize.UglifyJsPlugin({
        output: {
          comments: false
        }
      })),
      ifDev(new webpack.HotModuleReplacementPlugin()),
      ifDev(new webpack.NamedModulesPlugin())
    ])
  }
}
