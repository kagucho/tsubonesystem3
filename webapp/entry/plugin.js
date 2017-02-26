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

const SingleEntryPlugin = require("webpack/lib/SingleEntryPlugin");
const name = require("fs").realpathSync(__dirname);

module.exports = function() {
	let children = {};

	this.plugin("compilation", compilation => {
		compilation.plugin("normal-module-loader", loader => {
			loader[name] = (context, request) => {
				if (!children[context]) {
					children[context] = {};
				}

				if (!children[context][request]) {
					children[context][request] = new Promise((resolve, reject) => {
						const compiler = compilation.createChildCompiler(request, {
							filename: "entry.js",
						});

						compiler.context = context;
						compiler.apply(new SingleEntryPlugin(
							context, "!!" + request, "entry"));
						compiler.apply.apply(compiler, compiler.options.plugins);

						compiler.compile((error, compilation) => {
							if (error) {
								return reject(error);
							}

							const assets = {};

							for (const key in compilation.assets) {
								assets[key] = {
									extracted: false,
									value: compilation.assets[key],
								};
							}

							resolve({compilation, assets});
						});
					});
				}

				return children[context][request];
			};
		});
	});

	this.plugin("emit", (compilation, sync) => Promise.all((function *() {
		for (const context in children) {
			for (const request in children[context]) {
				yield children[context][request].then(result => {
					for (const asset in result.assets) {
						if (!result.assets[asset].extracted) {
							compilation.assets[asset] = result.assets[asset].value;
						}
					}
				});
			}
		}

		children = {};
	}())).then(value => sync(null, value), error => sync(error)));
}
