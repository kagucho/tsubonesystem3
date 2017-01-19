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

const getFileExtractUseEntries = name => [
	{
		loader:  "file-loader",
		options: {name},
	},
	"./loader/extract.js",
];

const htmlEntry = {
	loader: "html-loader",
	options: {interpolate: "require"},
};

module.exports = {
	// You SHOULD add html entrypoints to vnu task of package.json.
	entry: [
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
		"./src/graph.html",
		"./src/mail/html/creation.html",
		"./src/mail/text/creation.txt",
		"./src/footer.css",
		"./src/agpl-3.0.html",
		"./src/index.html",
		"./src/license.html",
		"./src/private/index.html",
		"file-loader?name=[name].[ext]!./src/favicon.ico",
	],
	module: {
		rules: [
			{
				include: path.resolve(__dirname, "src/footer.html"),
				loader: "html-loader",
			}, {
				include: path.resolve(__dirname, "src/footer.css"),
				use:     getFileExtractUseEntries("footer.css").concat(
						"css-loader", "postcss-loader"),
			}, {
				test:    /\.css$/,
				exclude: path.resolve(__dirname, "src/footer.css"),
				loaders: ["css-loader", "postcss-loader"],
			}, {
				test:    /\.html$/,
				include: path.resolve(__dirname, "src/error"),
				use:     getFileExtractUseEntries("../error/[name]").concat(htmlEntry),
			}, {
				include: path.resolve(__dirname, "src/graph.html"),
				use:     getFileExtractUseEntries("../graph").concat(htmlEntry),
			}, {
				include: path.resolve(__dirname, "src/mail/html"),
				use:     getFileExtractUseEntries("../mail/html/[name]").concat(
					"./loader/unix2dos.js", htmlEntry),
			}, {
				include: path.resolve(__dirname, "src/private/index.html"),
				use:     getFileExtractUseEntries("private").concat(htmlEntry),
			}, {
				include: [
					"src/agpl-3.0.html",
					"src/index.html",
					"src/license.html",
				].map(file => path.resolve(__dirname, file)),
				use: getFileExtractUseEntries("[name]").concat(htmlEntry),
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
					]
				},
			}, {
				test:   /\.(gif|jpg|png)?$/,
				loader: "file-loader",
			}, {
				include: path.resolve(__dirname, "src/mail/text"),
				use:     getFileExtractUseEntries("../mail/text/[name]").concat(
						"./loader/unix2dos.js"),
			},
		],
	},
	output:  {filename: "dummy", publicPath: "/"},
	plugins: [omitDummy],
};
