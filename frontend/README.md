# installation

```
$ npm install --save-dev glob
$ npm install --save-dev webpack webpack-cli webpack-dev-server
$ npm install --save-dev @babel/core @babel/preset-env @babel/preset-react
$ npm install --save-dev react react-dom
```

# structure

- frontend
    - public
        - js
        - index.html
    - src
        - App.js
    - webpack.config.js

## webpack.config.js

```
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
    port:3000
  },
  plugins: debug ? [] : [
    new webpack.optimize.OccurrenceOrderPlugin(),
    new webpack.optimize.UglifyJsPlugin({ mangle: false, sourcemap: false }),
  ]
};
```

## index.html

```
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>React Tutorials</title>
  </head>
  <body>
    <div id="app"></div>
    <script src="js/client.min.js"></script>
  </body>
</html>
```

## App.js

```
import React from "react";
import { createRoot } from 'react-dom/client';

class App extends React.Component {
  render() {
    return (
      <h1>It work!</h1>
    );
  }
}

const container = document.getElementById('app');
const root = createRoot(container);
root.render(<App />);
```

# client server

```
$ npm start
```

http://localhost:3000/

## SPA

devServer section

```
  historyApiFallback: {
     index: 'index.html'
  }
```

# Design

## CSS

```
$ npm install --save-dev css-loader
$ npm install --save-dev style-loader
```

added webpack.config.js(module.exports.module.rules)

```
{
  test: /\.css$/, 
  use: ['style-loader','css-loader']
}
```


## Material-Design(Lite)

```
$ npm install @mui/material @emotion/react @emotion/styled
```

```
<link
  rel="stylesheet"
  href="https://fonts.googleapis.com/icon?family=Material+Icons"
/>
```

# Router

```
$ npm install react-router-dom
```

# API Access

```
$ npm install --save axios
```

# Session(Cookie)


# Test


