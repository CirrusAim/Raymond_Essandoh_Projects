const webpack = require("webpack");
const path = require("path");
const CopyWebpackPlugin = require("copy-webpack-plugin");

module.exports = {
  mode: 'development',
  entry: {
    app: './index.js'
  },
  output: {
    filename: "app.js",
    path: path.resolve(__dirname, "dist"),
  },
  node: {
    global: true
  },
  resolve: {
    fallback: {
      buffer: require.resolve('buffer/'),
      stream: require.resolve("stream-browserify"),
      process: require.resolve('process/browser'),
      os: require.resolve("os-browserify/browser"),
      http: require.resolve("stream-http"),
      https: require.resolve("https-browserify"),
      crypto: require.resolve("crypto-browserify"),
      assert: require.resolve("assert/")
    },
  },
  module: {
    rules: [{
      loader: 'babel-loader',
      test: /\.js$|jsx/,
      exclude: /node_modules/
    },
    {
      test: /\.(scss)$/,
      use: [
        // Adds CSS to the DOM by injecting a `<style>` tag
        'style-loader',
        // Resolves `@import` and `url()` like `import/require()`
        'css-loader',
        // Loader for webpack to process CSS with PostCSS
        'postcss-loader',
        // Loader for SASS/SCSS files and CSS compiler
        'sass-loader'
      ]
    }],
  },
  plugins: [
    new CopyWebpackPlugin({
      patterns: [
        { from: "./index.html", to: "index.html" },
      ],
    }),
    new webpack.ProvidePlugin({
      process: 'process/browser',
      Buffer: ['buffer', 'Buffer'],
    }),
  ],
  devServer: {
    static: {
      directory: path.join(__dirname, 'dist'),
    },
    compress: true,
  },
};