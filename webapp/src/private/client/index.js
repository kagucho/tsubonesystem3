/**
	@file index.js is the main file of client module.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/client */

import Client from "./client";

/**
	module:private/client is the client.
	@constant
	@type !module:private/client/client
*/
export default new Proxy(new Client, {
	get(target, key) {
		return target[key] ?
			new Proxy(target[key], {
				apply(propertyTarget, propertyThis, propertyArguments) {
					return propertyTarget.apply(target, propertyArguments);
				},
			}) : target.constructor[key];
	},

	has(target, key) {
		return key in target || key in target.constructor;
	},
});
