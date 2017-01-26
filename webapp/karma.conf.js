/*
	Copyright (C) 2017  Kagucho <kagucho.net@gmail.com>

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published
	by the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

const webpackConfig = require("./webpack.config.js");
const webpack = require("webpack");

// Delete entries which use extract-text-webpack-plugin and fail.
delete webpackConfig.entry;

// Delete a plugin which prevents from emitting JS.
delete webpackConfig.plugins.shift();

// An alternative for CDN.
webpackConfig.plugins.push(new webpack.ProvidePlugin({$: "jquery"}));

module.exports = config => config.set({
	files: [
		"node_modules/mithril/mithril.js",
		"node_modules/babel-polyfill/dist/polyfill.js",
		{pattern: "src/**/*_test.js", watched: false},
	],
	frameworks: ["mocha"],
	plugins:    [
		"karma-chrome-launcher", "karma-coverage",
		"karma-edge-launcher", "karma-firefox-launcher",
		"karma-ie-launcher", "karma-mocha",
		"karma-phantomjs-launcher", "karma-sourcemap-loader",
		"karma-webpack",
	],
	preprocessors: {"src/**/*_test.js": ["webpack", "sourcemap", "coverage"]},
	proxies:       {"/api": process.env.TSUBONESYSTEM_URL + "api"},
	reporters:     ["coverage", "progress"],
	webpack:       webpackConfig,
});
