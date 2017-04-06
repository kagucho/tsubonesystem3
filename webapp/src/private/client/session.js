/**
	@file session.js implements session.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/client/session */

import * as api from "./api";

/**
	module:private/client/session is a class to implement session. TODO
	@extends Object
*/
export default function() {
	let filling;
	let id;
	let scope;
	let accessToken;
	let refreshToken;

	function storeID(newID) {
		sessionStorage.setItem("id", newID);
		id = newID;
	}

	function storeScope(newScope) {
		sessionStorage.setItem("scope", newScope);
		scope = newScope.split(" ");
	}

	function storeRefreshToken(newToken) {
		sessionStorage.setItem("refresh_token", newToken);
		refreshToken = newToken;
	}

	return {
		/**
			applyToken applies token to callback which returns a promise. TODO
			@param {module:session~Consumer} callback - The callback.
			@returns {!external:jQuery~Promise} A promise returned by the
			callback.
		*/
		applyToken(callback) {
			return callback(accessToken).catch(error => {
				if (error != "invalid_grant") {
					throw error;
				}

				return api.getTokenWithRefreshToken(refreshToken).then(data => {
					accessToken = data.access_token;

					if (data.refresh_token) {
						storeRefreshToken(data.refresh_token);
					}

					return applyToken(callback).then(
						null, null, progress => 0.5+progress);
				}, null, progress => progress/2);
			});
		},

		/**
			getFilling returns whether the session is limited for filling
			the information of the user. TODO
			@returns {?Boolean} true if the session is limited for filling
			the information of the user.
		*/
		getFilling() {
			return filling;
		},

		/**
			getID returns the ID of the member bound to the current session. TODO
			@returns {!String} The ID.
		*/
		getID() {
			return id;
		},

		/**
			getScope returns the scope of the current session. TODO
			@returns {!String[]} The scope.
		*/
		getScope() {
			return scope;
		},

		/**
			recover recovers session from sessionStorage. TODO
			@returns {!external:jQuery~Promise} A promise resolved when
			recovered session.
		*/
		recover() {
			const storedRefreshToken = sessionStorage.getItem("refresh_token");

			return api.getTokenWithRefreshToken(storedRefreshToken).then(data => {
				accessToken = data.access_token;
				id = sessionStorage.getItem("id");
				scope = sessionStorage.getItem("scope").split(" ");

				if (data.refresh_token) {
					storeRefreshToken(data.refresh_token);
				} else {
					refreshToken = storedRefreshToken;
				}
			});
		},

		/**
			setFillingToken sets the access token to fill the information of
			the user. TODO
			@param {!String} id - The ID.
			@param {!String} token - The access token.
			@returns {Undefined}
		*/
		setFillingToken(newID, token) {
			accessToken = token;
			filling = true;
			id = newID;
			scope = "update";
		},

		/**
			TODO
		*/
		updateToken(scope, newAccessToken, newRefreshToken) {
			accessToken = newAccessToken;
			sessionStorage.setItem("id", id);
			storeRefreshToken(newRefreshToken);
			storeScope(scope);
		},

		/**
			signin signs in. TODO
			@param {!String} id - The ID of the user.
			@param {!String} password - The password of the user.
			@returns {!external:jQuery~Promise} A promise resolved when
			signed in.
		*/
		signin(newID, password) {
			return api.getTokenWithPassword(newID, password).then(data => {
				accessToken = data.access_token;

				storeID(newID);
				storeScope(data.scope);
				storeRefreshToken(data.refresh_token);
			});
		},
	};
}

/**
	Consumer is a callback for applyToken.
	@param {!String} token - The access token.
	@returns {!external:jQuery~Promise} A promise which may reject with
	jqXHR.
	@callback module:session~Consumer
*/
