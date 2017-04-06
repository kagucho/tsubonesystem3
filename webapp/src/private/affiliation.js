/**
	@file affiliation.js implements an interface for affiliation datalist.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/affiliation */

import client from "./client";

/**
	id is the ID attribute of affiliation datalist.
	@type !String
*/
export const id = "affiliation";

let stream;

/**
	update updates affiliation datalist with the given affiliations. TODO
	@param {!Iterable.<String>} affiliations - The
	affiliations.
	@returns {Undefined}
*/
export const listen = callback => {
	if (!stream) {
		stream = client.mapMembers().map(
			promise => promise.done(
			members => m.render(
				document.getElementById(id),
				Array.from((function *() {
					for (const id in members) {
						yield m("option", {value: members[id].affiliation});
					}
				})()))));
	}

	return stream.map(callback);
}
