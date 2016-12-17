/*
  Copyright (C) 2016  Kagucho <kagucho.net@gmail.com>

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU Affero General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU Affero General Public License for more details.

  You should have received a copy of the GNU Affero General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

const Text = require("extract-text-webpack-plugin");
const css = new Text("[name].css");
const html = new Text("[name]");

function omitJS() {
  this.plugin("emit", (compilation, callback) => {
    for (const key in compilation.assets)
      if (key.endsWith(".js"))
        delete compilation.assets[key];

    callback();
  });
}

module.exports = {
  // You SHOULD add html entrypoints to vnu task of package.json.
  entry: {
    "error": ["./src/error/common.css"],
    "../error/301": ["./src/error/301.html"],
    "../error/400": ["./src/error/400.html"],
    "../error/404": ["./src/error/404.html"],
    "../error/405": ["./src/error/405.html"],
    "../error/406": ["./src/error/406.html"],
    "../error/408": ["./src/error/408.html"],
    "../error/412": ["./src/error/412.html"],
    "../error/413": ["./src/error/413.html"],
    "../error/414": ["./src/error/414.html"],
    "../error/416": ["./src/error/416.html"],
    "../error/417": ["./src/error/417.html"],
    "../error/421": ["./src/error/421.html"],
    "../error/500": ["./src/error/500.html"],
    "../error/501": ["./src/error/501.html"],
    "../error/503": ["./src/error/503.html"],
    "../graph": ["./src/error/graph.html"],
    "../error/unknown": ["./src/error/unknown.html"],
    "favicon": ["file?name=[name].[ext]!./src/favicon.ico"],
    "index": ["babel-polyfill", "./src/index.html"],
    "license": ["babel-polyfill", "./src/license.html"],
    "agpl-3.0": ["./src/agpl-3.0.html"],
    "private": ["babel-polyfill", "./src/private/index.html"],
  },
  module: {
    preLoaders: [
      {
        test: /\.tag$/,
        exclude: /node_modules/,
        loader: "riotjs",
      },
    ],
    loaders: [
      {
        test: /src\/error\/footer\.html$/,
        loader: "html",
      }, {
        test: /src\/error\/.+\.css$/,
        exclude: /index.css$/,
        loader: css.extract(["css", "postcss"]),
      }, {
        test: /\.css$/,
        exclude: /src\/error\/.+\.css$/,
        loaders: ["css", "postcss"],
      }, {
        test: /\.html$/,
        exclude: /src\/error\/footer\.html$/,
        loader: html.extract("html?interpolate=require"),
      }, {
        test: /\.(js|tag)$/,
        exclude: /node_modules/,
        loader: "babel",
        query: {presets: ["es2015", "es2016"]},
      }, {
        test: /\.(gif|jpg|png)?$/,
        loader: "file",
      },
    ],
  },
  output: {filename: "[name].js", publicPath: "/"},
  plugins: [omitJS, css, html],
  postcss: () => [
    require("stylelint")({
      config: {
        "extends": "stylelint-config-standard",
        "rules": {"selector-list-comma-newline-after": "never-multi-line"},
      },
    }),
    require("postcss-cssnext")({
      browsers: [
        "last 1 Chrome versions",
        "Edge >= 20",
        "Firefox ESR",
        "ie >= 11",
      ],
    }),
  ],
};
