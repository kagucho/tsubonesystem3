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

const LibraryTemplatePlugin = require("webpack/lib/LibraryTemplatePlugin");
const LimitChunkCountPlugin = require("webpack/lib/optimize/LimitChunkCountPlugin");
const NodeTemplatePlugin = require("webpack/lib/node/NodeTemplatePlugin");
const NodeTargetPlugin = require("webpack/lib/node/NodeTargetPlugin");
const SingleEntryPlugin = require("webpack/lib/SingleEntryPlugin");

const compilers = {};

module.exports = function() {
};

module.exports.pitch = function(request) {
	if (!compilers[request]) {
		const options = {
			filename:   "entry",
			publicPath: this._compilation.outputOptions.publicPath,
		};

		// eslint-disable-next-line no-underscore-dangle
		compilers[request] = this._compilation.createChildCompiler(request, options);

		compilers[request].apply(new NodeTemplatePlugin(options));
		compilers[request].apply(new LibraryTemplatePlugin(null, "commonjs2"));
		compilers[request].apply(new NodeTargetPlugin());
		compilers[request].apply(new SingleEntryPlugin(
			this.context, "!!" + request, "entry"));
	}

	const sync = this.async();
	compilers[request].compile((error, compilation) => {
		if (error) {
			return sync(error);
		}

		if (compilation.errors.length > 0) {
			return sync(compilation.errors[0]);
		}

		for (const name in compilation.assets) {
			if (name != "entry") {
				// eslint-disable-next-line no-underscore-dangle
				this._compilation.assets[name] = compilation.assets[name];
			}
		}

		compilation.fileDependencies.forEach(this.addDependency);
		compilation.contextDependencies.forEach(this.addContextDependency);

		sync(null, this.exec(compilation.assets.entry.source(), request));
	});
};
