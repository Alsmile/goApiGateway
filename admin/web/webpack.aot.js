const webpack = require('webpack');
const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CompressionPlugin = require("compression-webpack-plugin");
const ngtools = require('@ngtools/webpack');

const ExtractTextPlugin = require('extract-text-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const LoaderOptionsPlugin = require('webpack/lib/LoaderOptionsPlugin');

module.exports = function (options) {
  if (!options) options = {};
  if (!options.env) options.env = 'production';

  return {
    entry: {
      polyfills: './src/polyfills.ts',
      vendor: './src/vendor.ts',
      app: './src/main.ts'
    },
    output: {
      publicPath: '/assets/',
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
          test: /\.css$/,
          loader: ExtractTextPlugin.extract({
            fallback: "style-loader",
            use: [{
              loader: "css-loader"
            }]
          })
        },
        {
          test: /\.html$/,
          use: 'raw-loader',
          exclude: ['./src/index.html', './src/index.inner.html']
        },
        {
          test: /\.(jpg|png|gif)$/,
          use: 'file-loader'
        }
      ]
    },
    plugins: [
      new webpack.DefinePlugin({
        'process.env': {
          NODE_ENV: JSON.stringify(options.env)
        }
      }),
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
        {
          from: 'node_modules/monaco-editor/min',
          to: 'monaco',
        }
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
                  'Safari >= 7.1'
                ]
              })
            ];
          }
        }
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
    }
  };
};
