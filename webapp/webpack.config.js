/*
	Copyright (C) 2017  Kagucho <kagucho.net@gmail.com>

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published
	by the Free Software Foundation, either version 3 of the License, or (at
	your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

const path = require("path");

function omitDummy() {
	this.plugin("emit", (compilation, callback) => {
		delete compilation.assets.dummy;
		callback();
	});
}

const getFileEntry = name => ({
	loader:  "file-loader",
	options: {name},
});

const htmlEntry = {
	loader:  "html-loader",
	options: {
		ignoreCustomComments: [/Copyright/],
		interpolate: true,
		minifyCSS:   false,
	},
};

module.exports = {
	// You SHOULD add html entrypoints to vnu task of package.json.
	entry: [
		"!!file-loader?name=../graph!./loader/extract!html-loader?interpolate=require?removeComments=false!./src/graph.html",
		"!!file-loader?name=error.css!./loader/extract!css-loader!postcss-loader!./src/error/common.css",
		"!!file-loader?name=favicon.ico!./src/favicon.ico",
		"!!file-loader?name=footer.css!./loader/extract!css-loader!postcss-loader!./src/footer.css",
		"./src/error/301.html",
		"./src/error/400.html",
		"./src/error/404.html",
		"./src/error/405.html",
		"./src/error/406.html",
		"./src/error/408.html",
		"./src/error/412.html",
		"./src/error/413.html",
		"./src/error/414.html",
		"./src/error/416.html",
		"./src/error/417.html",
		"./src/error/421.html",
		"./src/error/500.html",
		"./src/error/501.html",
		"./src/error/503.html",
		"./src/error/unknown.html",
		"./src/mail/html/message.html",
		"./src/mail/html/confirmation.html",
		"./src/mail/html/creation.html",
		"./src/mail/html/invitation.html",
		"./src/mail/text/message.txt",
		"./src/mail/text/confirmation.txt",
		"./src/mail/text/creation.txt",
		"./src/mail/text/invitation.txt",
		"./src/agpl-3.0.html",
		"./src/index.html",
		"./src/license.html",
		"./src/private/index.html",
	],
	module: {
		rules: [
			{
				include: path.resolve(__dirname, "src/footer.html"),
				loader:  "html-loader",
			}, {
				test:    /\.css$/,
				loaders: ["css-loader", "postcss-loader"],
			}, {
				test:    /\.html$/,
				include: path.resolve(__dirname, "src/error"),
				use:     [
					getFileEntry("../error/[name]"),
					"./loader/extract",
					htmlEntry,
				],
			}, {
				include: path.resolve(__dirname, "src/mail/html"),
				use:     [
					getFileEntry("../mail/html/[name]"),
					"./loader/unix2dos",
					"./loader/extract",
					htmlEntry,
				],
			}, {
				include: path.resolve(__dirname, "src/private/index.html"),
				use:     [
					getFileEntry("private"),
					"./loader/extract",
					htmlEntry,
				],
			}, {
				include: [
					"src/agpl-3.0.html",
					"src/index.html",
					"src/license.html",
				].map(file => path.resolve(__dirname, file)),
				use: [
					getFileEntry("[name]"),
					"./loader/extract",
					htmlEntry,
				],
			}, {
				test:    /\.(js|tag)$/,
				exclude: /node_modules/,
				loader:  "babel-loader",
				options: {
					presets: [
						[
							"es2015",
							{modules: false},
						], [
							"es2016",
						],
					],
				},
			}, {
				test:   /\.(gif|jpg|png)?$/,
				loader: "file-loader",
			}, {
				include: path.resolve(__dirname, "src/mail/text"),
				use:     [
					getFileEntry("../mail/text/[name]"),
					"./loader/unix2dos",
				],
			},
		],
	},
	output:  {filename: "dummy", publicPath: "/"},
	plugins: [omitDummy, require("./entry/plugin")],
};

if (process.env.NODE_ENV == "production") {
	module.exports.plugins.push(
		new (require("webpack/lib/optimize/UglifyJsPlugin"))({
			collapse_vars: true,
			comments:      null,
			pure_getters:  true,
			reduce_vars:   true,
			unsafe:        true,
			unsafe_comps:  true,
			unsafe_proto:  true,
			warnings:      true,
		})
	);
}
