/**
	@file index.js is the main file of client module.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/client */

import * as client from "./client";
import * as api from "./api";

const instance = client.default();
instance.error = api.error;
instance.merge = client.merge;

/**
	module:private/client is the client.
	@constant
	@type !module:private/client/client
*/
export default instance;
