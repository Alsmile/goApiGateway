const webpack = require('webpack');
const webpackMerge = require('webpack-merge');
const prodConfig = require('./webpack.aot');

module.exports = function (options) {
  return prodConfig({env: 'dev'});
  // return webpackMerge(prodConfig(options), {
  //   plugins: [
  //     new webpack.DefinePlugin({
  //       'process.env': {
  //         NODE_ENV: JSON.stringify('dev')
  //       }
  //     })
  //   ]
  // });
};
