var debug   = process.env.NODE_ENV !== "production";
var webpack = require('webpack');
var glob    = require('glob');
var path    = require('path');

const entries = glob.sync("./src/**/*.js");

module.exports = {
  entry: entries,
  module: {
    rules: [{
      test: /\.jsx?$/,
      exclude: /(node_modules|bower_components)/,
      use: [{
        loader: 'babel-loader',
        options: {
          presets: ['@babel/preset-react', '@babel/preset-env']
        }
      }]
    },{ 
      test: /\.css$/, 
      use: ['style-loader','css-loader']
    }]
  },
  output: {
    path: path.resolve(__dirname,"public"),
    filename: "js/client.min.js"
  },
  devServer: {
    static: {
      directory: path.resolve(__dirname,"public")
    },
    historyApiFallback: {
      index: 'index.html'
    },
    proxy: {
        "/api":"http://localhost:5555"
    },
    port:3000
  },
  plugins: debug ? [] : [
    new webpack.optimize.OccurrenceOrderPlugin(),
    new webpack.optimize.UglifyJsPlugin({ mangle: false, sourcemap: false }),
  ]
};
