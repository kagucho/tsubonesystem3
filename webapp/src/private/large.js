/**
	@file large.js implements a detection of large display area.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/large */

/**
	module:private/large returns wether the display area is large.
	@returns {!Boolean} - The boolean showing whether the
	display area is large or not.
*/
export default function() {
	const container = $("#container");
	return container.width() / parseFloat(container.css("font-size")) > 64;
}
