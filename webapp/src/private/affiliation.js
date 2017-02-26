/**
	@file affiliation.js implements an interface for affiliation datalist.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/affiliation */

/**
	id is the ID attribute of affiliation datalist.
	@type !String
*/
export const id = "affiliation";

/**
	update updates affiliation datalist with the given affiliations.
	@param {!Iterable.<String>} affiliations - The
	affiliations.
	@returns {Undefined}
*/
export function update(affiliations) {
	m.render(document.getElementById(id),
		Array.from(affiliations,
			affiliation => m("option", {value: affiliation})));
}
