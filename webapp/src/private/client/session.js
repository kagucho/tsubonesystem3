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
	@private
	@type !WeakMap.<module:private/client/session, module:private/client/session~State>
*/
const state = new WeakMap();

/**
	State is a class to represent the internal state of
	module:private/client/session.
	@private
	@extends Object
*/
class State {
	/**
		storeID stores the ID to the sessionStorage as well as the
		class.
		@param {!String} id - The ID.
		@returns {Undefined}
	*/
	storeID(id) {
		sessionStorage.setItem("id", id);
		this.id = id;
	}

	/**
		storeScope stores the scope to the sessionStorage as well as the
		class.
		@param {!String} scope - The scope.
		@returns {Undefined}
	*/
	storeScope(scope) {
		sessionStorage.setItem("scope", scope);
		this.scope = scope.split(" ");
	}

	/**
		storeRefreshToken stores the token to the sessionStorage as
		well as the class.
		@param {!String} token - The refresh token.
		@returns {Undefined}
	*/
	storeRefreshToken(token) {
		sessionStorage.setItem("refresh_token", token);
		this.refreshToken = token;
	}

	/**
		id is the ID of the user.
		@member {?String} module:private/client/session~State#id
	*/

	/**
		scope is the scope of the current session.
		@member {?String[]} module:private/client/session~State#scope
	*/

	/**
		accessToken is the access token of the current session
		@member {?String} module:private/client/session~State#accessToken
	*/

	/**
		refreshToken is the refresh token of the current session.
		@member {?String} module:private/client/session~State#refreshToken
	*/
}

/**
	module:private/client/session is a class to implement session.
	@extends Object
*/
export default class {
	/**
		constructor constructs module:private/client/session.
		@returns Undefined
	*/
	constructor() {
		state.set(this, new State);
		Object.freeze(this);
	}

	/**
		applyToken applies token to callback which returns a promise.
		@param {module:session~Consumer} callback - The callback.
		@returns {!external:jQuery~Promise} A promise returned by the
		callback.
	*/
	applyToken(callback) {
		const local = state.get(this);

		return callback(local.accessToken).catch(function(xhr) {
			if (xhr.status != 401) {
				return $.Deferred().reject(...arguments);
			}

			return api.getTokenWithRefreshToken(local.refreshToken).then(data => {
				local.accessToken = data.access_token;

				if (data.refresh_token) {
					local.storeRefreshToken(data.refresh_token);
				}

				return this.applyToken(callback).then(
					null, null, progress => 0.5+progress);
			}, null, progress => progress/2);
		}.bind(this));
	}

	/**
		getFilling returns whether the session is limited for filling
		the information of the user.
		@returns {?Boolean} true if the session is limited for filling
		the information of the user.
	*/
	getFilling() {
		return state.get(this).filling;
	}

	/**
		getID returns the ID of the member bound to the current session.
		@returns {!String} The ID.
	*/
	getID() {
		return state.get(this).id;
	}

	/**
		getScope returns the scope of the current session.
		@returns {!String[]} The scope.
	*/
	getScope() {
		return state.get(this).scope;
	}

	/**
		recover recovers session from sessionStorage.
		@returns {!external:jQuery~Promise} A promise resolved when
		recovered session.
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
		setFillingToken sets the access token to fill the information of
		the user.
		@param {!String} id - The ID.
		@param {!String} token - The access token.
		@returns {Undefined}
	*/
	setFillingToken(id, token) {
		const local = state.get(this);

		local.accessToken = token;
		local.filling = true;
		local.id = id;
		local.scope = "update";
	}

	/**
		signin signs in.
		@param {!String} id - The ID of the user.
		@param {!String} password - The password of the user.
		@returns {!external:jQuery~Promise} A promise resolved when
		signed in.
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
	Consumer is a callback for applyToken.
	@param {!String} token - The access token.
	@returns {!external:jQuery~Promise} A promise which may reject with
	jqXHR.
	@callback module:session~Consumer
*/
