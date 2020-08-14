const { CleanWebpackPlugin } = require("clean-webpack-plugin");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const Dotenv = require("dotenv-webpack");
const path = require("path");
const envpath = process.env.NODE_ENV ? `.env.${process.env.NODE_ENV}` : ".env";

module.exports = {
  entry: "./src/index.tsx",
  plugins: [
    new CleanWebpackPlugin({
      cleanAfterEveryBuildPatterns: ["public/build"],
    }),
    new HtmlWebpackPlugin({
      template: "src/index.html",
    }),
    new Dotenv({
      path: envpath,
      safe: true,
    }),
  ],
  output: {
    path: __dirname + "/public",
    filename: "build/[name].[contenthash].js",
    publicPath: "/",
  },
  resolve: {
    extensions: [".ts", ".tsx", ".js"],
    alias: {
      "~": path.resolve(__dirname, "src"),
    },
  },
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        loader: "ts-loader"
      },
      {
       test: /\.(jpg|png|svg|gif)$/,
       use: { loader: "url-loader" },
      }
    ],
  },
  devServer: {
    historyApiFallback: true,
    port: 3001,
  },
};
