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

"use strict";

console.log("You can test file.ServeHTTP serves a static file.");
console.log("Request /index.js and assert it writes the content of this file.");

console.log("You can test file.ServeHTTP doesn't set Content-Language field");
console.log("for files which have extensions (i.e. not HTML).");
console.log("Request /index.js and assert it doesn't.");
