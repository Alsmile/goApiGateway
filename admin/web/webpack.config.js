let webpack = require('webpack');
let path = require('path');
let webpackMerge = require('webpack-merge');
let HtmlWebpackPlugin = require('html-webpack-plugin');

let ExtractTextPlugin = require('extract-text-webpack-plugin');
let CopyWebpackPlugin = require('copy-webpack-plugin');
let LoaderOptionsPlugin = require('webpack/lib/LoaderOptionsPlugin');

module.exports = function (options) {
  return {
    entry: {
      polyfills: './src/polyfills.ts',
      vendor: './src/vendor.ts',
      app: './src/main.ts'
    },
    output: {
      path: path.resolve(__dirname, './dist'),
      filename: '[name].bundle.js',
      sourceMapFilename: '[file].map',
      chunkFilename: '[id].chunk.js',
    },
    resolve: {
      extensions: ['.ts', '.js'],
      modules: [path.resolve(__dirname, 'node_modules')]
    },

    module: {
      rules: [
        {
          test: /\.ts$/,
          use: [
            'awesome-typescript-loader',
            'angular-router-loader',
            'angular2-template-loader'
          ],
          exclude: [/\.(spec|e2e)\.ts$/]
        },
        {
          test: /\.pcss$/,
          loader: ExtractTextPlugin.extract({
            fallback: "style-loader",
            use: [{
              loader: "css-loader"
            }, {
              loader: "postcss-loader"
            }]
          })
        },
        {
          test: /\.html$/,
          use: 'raw-loader',
          exclude: [path.resolve(__dirname, './src/index.html')]
        },
        {
          test: /\.(jpg|png|gif)$/,
          use: 'file-loader'
        }
      ]
    },
    plugins: [
      new webpack.optimize.CommonsChunkPlugin({
        name: ['app', 'vendor', 'polyfills']
      }),
      new CopyWebpackPlugin([
        { from: 'src/assets/img', to: 'assets/img' },
        { from: 'src/browser.html', to: '../browser.html' },
      ]),
      new ExtractTextPlugin('[name].css'),
      new HtmlWebpackPlugin({
        filename: 'index.html',
        template: 'src/index.html'
      }),
      new LoaderOptionsPlugin({
        minimize: false,
        debug: true,
        options: {
          postcss: function (bundler) {
            return [
              require('postcss-import')({ addDependencyTo: bundler }),
              require('precss')()
            ];
          },
        },
      })
    ],
    node: {
      global: true,
      crypto: 'empty',
      __dirname: true,
      __filename: true,
      Buffer: false,
      clearImmediate: false,
      setImmediate: false
    },
    devServer: {
      port: 8501,
      host: 'localhost',
      historyApiFallback: true,
      watchOptions: {
        aggregateTimeout: 300,
        poll: 1000
      },
      proxy: {
        "/api": {
          target: "http://localhost:8500",
          changeOrigin: true,
          headers: {host: ''},
          secure: false
        },
        "/captcha": {
          target: "http://localhost:8500",
          changeOrigin: true,
          secure: false
        }
      }
    }
  };
}
