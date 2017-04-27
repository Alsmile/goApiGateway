let webpack = require('webpack');
let path = require('path');
let webpackMerge = require('webpack-merge');
let HtmlWebpackPlugin = require('html-webpack-plugin');
let CompressionPlugin = require("compression-webpack-plugin");
let ngtools = require('@ngtools/webpack');

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
      publicPath: '/assets',
      path: path.resolve(__dirname, './dist/assets'),
      filename: '[name].[chunkhash].js',
      sourceMapFilename: '[name].[chunkhash].map',
      chunkFilename: '[id].[chunkhash].js'
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
            '@ngtools/webpack'
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
          exclude: ['./src/index.html']
        },
        {
          test: /\.(jpg|png|gif)$/,
          use: 'file-loader'
        }
      ]
    },
    plugins: [
      new ngtools.AotPlugin({
        tsConfigPath: './tsconfig.json',
        skipMetadataEmit: true,
        entryModule: 'src/app/app.module#AppModule'
      }),
      new ExtractTextPlugin('[name].[contenthash].css'),
      new webpack.optimize.CommonsChunkPlugin({
        name: ['app', 'vendor', 'polyfills']
      }),
      new CompressionPlugin({
        asset: "[path].gz[query]",
        algorithm: "gzip",
        test: /\.js$|\.html$/,
        threshold: 10240,
        minRatio: 0.3
      }),
      new webpack.optimize.UglifyJsPlugin(),
      new CopyWebpackPlugin([
        { from: 'src/assets/img', to: 'img' },
        { from: 'src/browser.html', to: '../browser.html' },
      ]),
      new HtmlWebpackPlugin({
        filename: '../index.html',
        template: 'src/index.html'
      }),
      new LoaderOptionsPlugin({
        minimize: true,
        debug: false,
        options: {
          postcss: function (bundler) {
            return [
              require('postcss-import')({ addDependencyTo: bundler }),
              require('precss')(),
              require('autoprefixer')({
                browsers: [
                  'Android >= 4',
                  'Chrome >= 35',
                  'Firefox >= 31',
                  'Explorer >= 10',
                  'iOS >= 7',
                  'Opera >= 12',
                  'Safari >= 7.1',
                ]
              })
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
      port: 8101,
      host: 'localhost',
      historyApiFallback: true,
      watchOptions: {
        aggregateTimeout: 300,
        poll: 1000
      },
      proxy: {
        "/api": {
          target: "http://www.i-dengta.com/api/proxy/58c965d26025d75788966da0",
          changeOrigin: true,
          // headers: {host: 'www.i-dengta.com'},
          secure: false
        },
        "/captcha": {
          target: "http://www.i-dengta.com",
          changeOrigin: true,
          // headers: {host: 'www.i-dengta.com'},
          secure: false
        }
      }
    }
  };
}
