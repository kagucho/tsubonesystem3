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

const name = require("fs").realpathSync(__dirname);
const ModuleError = require("webpack/lib/ModuleError");

module.exports = function() {
};

module.exports.pitch = function(request) {
	const sync = this.async();

	this[name](this.context, request).then(result => {
		const query = this.query ? this.query.slice(1) : "entry.js";

		result.assets[query].extracted = true;
		result.compilation.contextDependencies.forEach(this.addContextDependency);
		result.compilation.fileDependencies.forEach(this.addDependency);

		this._module.errors.push(...result.compilation.errors);
		this._module.warnings.push(...result.compilation.warnings);

		if (result.compilation.errors.length) {
			sync("child compiler failed");
		} else {
			sync(null, "module.exports = " + JSON.stringify(result.assets[query].value.source()));
		}
	}).catch(sync);
};
