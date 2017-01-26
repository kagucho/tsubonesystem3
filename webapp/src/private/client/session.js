/**
	@file session.js implements session.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/client/session */

import * as api from "./api";

/**
	state holds private state of client module.
	@type !external:ES.WeakMap<module:session, external:ES.Object>
*/
const state = new WeakMap();

/**
	module:private/client/session is a class to implement session.
	@extends external:ES.Object
*/
export default class {
	/** constructor constructs a new instance. */
	constructor() {
		state.set(this, {
			storeID(id) {
				sessionStorage.setItem("id", id);
				this.id = id;
			},

			storeScope(scope) {
				sessionStorage.setItem("scope", scope);
				this.scope = scope.split(" ");
			},

			storeRefreshToken(token) {
				sessionStorage.setItem("refresh_token", token);
				this.refreshToken = token;
			},
		});

		Object.freeze(this);
	}

	/**
		applyToken applies token to callback which returns a promise.
		@param {module:session~tokenConsumer} callback - The callback.
		@returns {!external:jQuery.$.Deferred#promise} A promise
		returned by the callback.
	*/
	applyToken(callback) {
		const local = state.get(this);

		return callback(local.accessToken).catch(xhr => {
			if (xhr.status != 401) {
				throw xhr;
			}

			return api.getTokenWithRefreshToken(local.refreshToken).progress(
			progress => progress/2).then(data => {
				local.accessToken = data.access_token;

				if (data.refresh_token) {
					local.storeRefreshToken(data.refresh_token);
				}

				return this.applyToken(callback).progress(
					progress => 0.5+progress);
			});
		});
	}

	/**
		getID returns the ID of the member bound to the current session.
		@returns {!external:ES.String} The ID.
	*/
	getID() {
		return state.get(this).id;
	}

	/**
		getScope returns the scope of the current session.
		@returns {!external:ES.String[]} The scope.
	*/
	getScope() {
		return state.get(this).scope;
	}

	/**
		recover recovers session from sessionStorage.
		@returns {!external:jQuery.$.Deferred#promise} A promise
		resolved when recovered session.
	*/
	recover() {
		const refreshToken = sessionStorage.getItem("refresh_token");

		return api.getTokenWithRefreshToken(refreshToken).then(data => {
			const local = state.get(this);

			local.accessToken = data.access_token;
			local.id = sessionStorage.getItem("id");
			local.scope = sessionStorage.getItem("scope").split(" ");

			if (data.refresh_token) {
				local.storeRefreshToken(data.refresh_token);
			} else {
				local.refreshToken = refreshToken;
			}
		});
	}

	/**
		signin signs in.
		@returns {!external:jQuery.$.Deferred#promise} A promise
		resolved when signed in.
	*/
	signin(id, password) {
		return api.getTokenWithPassword(id, password).then(data => {
			const local = state.get(this);

			local.accessToken = data.access_token;

			local.storeID(id);
			local.storeScope(data.scope);
			local.storeRefreshToken(data.refresh_token);
		});
	}
}

/**
	tokenConsumer is a callback for applyToken.
	@param {!external:ES.String} token - The access token.
	@returns {!external:jQuery.$.Deferred#promise} A promise which may
	reject with jqXHR.
	@callback module:session~tokenConsumer
*/
