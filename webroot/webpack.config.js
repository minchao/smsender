const path = require('path')
const webpack = require('webpack')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')

module.exports = () => {
  const env = process.env.NODE_ENV
  const mode = (env === 'production') ? 'production' : 'development'
  const ifProd = plugin => (mode === 'production') ? plugin : undefined
  const ifDev = plugin => (mode === 'development') ? plugin : undefined
  const removeEmpty = array => array.filter(p => !!p)

  return {
    devtool: ifDev('source-map'),
    mode: mode,
    entry: {
      main: removeEmpty([
        ifDev('react-hot-loader/patch'),
        ifDev('webpack-dev-server/client?http://localhost:3000'),
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
          test: /\.(css)$/,
          use:
            [
              (mode === 'development') ? 'style-loader' : MiniCssExtractPlugin.loader,
              'css-loader?modules=true',
            ]
        }
      ]
    },
    optimization: {
      minimize: mode === 'production',
      splitChunks: {
        chunks: 'all'
      }
    },
    plugins: removeEmpty([
      new HtmlWebpackPlugin({
        template: path.join(__dirname, './src/index.html'),
        filename: 'index.html',
        inject: 'body'
      }),
      new webpack.DefinePlugin({
        __DEVELOPMENT__: mode === 'development',
        API_HOST: JSON.stringify(mode === 'development' ? 'http://localhost:8080' : '')
      }),
      ifDev(new webpack.HotModuleReplacementPlugin()),
      ifDev(new webpack.NamedModulesPlugin()),
      ifProd(new MiniCssExtractPlugin({
        filename: '[name].[hash].css'
      }))
    ])
  }
}
