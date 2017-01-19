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

module.exports = function(source) {
	this.cacheable();

	let destination;
	let destinationIndex;

	if (this.loaderIndex + 1 == this.loaders.length) {
		const prefix = "module.exports = \"";
		const suffix = "\";";

		destination = new Buffer(prefix.length + source.length * 4 + suffix.length);
		destinationIndex = destination.write(prefix);

		for (let sourceIndex = 0; sourceIndex < source.length; sourceIndex++) {
			switch (source[sourceIndex]) {
			case 10:
				destinationIndex += destination.write("\\r\\n", destinationIndex);
				break;

			case 34:
				destinationIndex += destination.write("\\\"", destinationIndex);
				break;

			case 92:
				destinationIndex += destination.write("\\\\", destinationIndex);
				break;

			default:
				if (source[sourceIndex] < 32) {
					this.emitError("invalid character "+source[sourceIndex]+" at "+sourceIndex);
					return;
				}

				destination[destinationIndex] = source[sourceIndex];
				destinationIndex++;
				break;
			}
		}

		destinationIndex += destination.write(suffix, destinationIndex);
	} else {
		destination = new Buffer(source.length * 2);
		destinationIndex = 0;

		for (let sourceIndex = 0; sourceIndex < source.length; sourceIndex++) {
			if ((source[sourceIndex] != 92 || source[sourceIndex] != 114) &&
				source[sourceIndex] == 92 && source[sourceIndex + 1] == 110) {
				destinationIndex += destination.write("\\r", destinationIndex);
			}

			destination[destinationIndex] = source[sourceIndex];
			destinationIndex++;
		}
	}

	return destination.slice(0, destinationIndex);
}

module.exports.raw = true;
