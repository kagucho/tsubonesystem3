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

const SingleEntryPlugin = require("webpack/lib/SingleEntryPlugin");

const compilers = {};

module.exports = function() {
};

module.exports.pitch = function(request) {
  if (!compilers[request]) {
    // eslint-disable-next-line no-underscore-dangle
    compilers[request] = this._compilation.createChildCompiler(request);
    compilers[request].apply(
      new SingleEntryPlugin(this.context, "!!" + request, "entry"));
  }

  const sync = this.async();
  compilers[request].compile((error, compilation) => {
    if (error)
      return sync(error);

    if (compilation.errors.length > 0)
      return sync(compilation.errors[0]);

    for (const name in compilation.assets)
      if (name != "entry.js")
        // eslint-disable-next-line no-underscore-dangle
        this._compilation.assets[name] = compilation.assets[name];

    compilation.fileDependencies.forEach(this.addDependency);
    compilation.contextDependencies.forEach(this.addContextDependency);

    sync(null,
         "module.exports = " +
           JSON.stringify(compilation.assets["entry.js"].source()));
  });
};
