/**
	@file api.js implements a WebAPI interface.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/client/api */

import {Wrapper} from "../promise";

/**
	ajax returns fetched and decoded JSON.
	@private
	@param {!String} uri - The URI.
	@param {!String} method - The method.
	@param {?String} [token] - The token.
	@param {?*} [data] - The data.
	@returns {!module:private/promise} A promise resolved with fetched JSON.
*/
function ajax(uri, method, token, data) {
	const xhr = new XMLHttpRequest();

	const opts = {
		accepts:  {json: "application/json; charset=UTF-8"},
		data, dataType: "json", method,
		xhr() {
			return xhr;
		},
	};

	/*
		RFC 6750 - The OAuth 2.0 Authorization Framework: Bearer Token Usage
		1.  Introduction
		https://tools.ietf.org/html/rfc6750#section-3
		The Bearer authentication scheme is intended primarily for
		server authentication using the WWW-Authenticate and
		Authorization HTTP headers but does not preclude its use for
		proxy authentication.
	*/
	if (token) {
		opts.headers = {Authorization: "Bearer "+token};
	}

	const jqXHR = $.ajax(uri, opts);
	const deferred = $.Deferred();
	const wrapper = new Wrapper;
	let uploadProgress = {loaded: 0, total: 0};

	xhr.upload.onprogress = event => {
		deferred.notify(wrapper.progressProxy(event));
		uploadProgress = event;
	};

	xhr.onprogress = event => deferred.notify(wrapper.progressProxy(Object.defineProperties({}, {
		lengthComputable: {
			get() {
				return uploadProgress.lengthComputable || event.lengthComputable;
			},
		},

		loaded: {
			get() {
				return uploadProgress.loaded + event.loaded;
			},
		},

		total: {
			get() {
				return uploadProgress.total + event.total;
			},
		},
	})));

	jqXHR.then(
		(response, status, event) => deferred.resolve(response, status, wrapper.alwaysProxy(event)),
		(event, status, xhrError) => deferred.reject(wrapper.alwaysProxy(event), status, xhrError)
	);

	return wrapper.wrap(deferred);
}

/**
	error returns the error message corresponding to jqXHR.
	@param   {!external:jQuery~jqXHR} - A failed jqXHR.
	@returns {!String} The error message.
*/
export function error(xhr) {
	switch (xhr.status) {
	case 0:
		return "TsuboneSystemへの経路上に問題が発生しました。ネットワーク接続などを確認してください。";

	case 400:
		return {
			/* eslint-disable camelcase */
			invalid_grant: "あんた誰?って言われちゃいました。もう一度サインインしてください。",
			invalid_id:    "IDが違うってよ。",
			/* eslint-enable camelcase */
		}[xhr.responseJSON.error];

	case 429:
		return "残念！！やりすぎです。ちょっと待ってください。";

	default:
		return "サーバー側のエラーです。がびーん。";
	}
}

/**
	getTokenWithPassword returns a token authorized with the given password.
	@function
	@param {!String} username - The ID.
	@param {!String} password - The password.
	@returns {!module:private/promise} A promise resolved with a token.
*/
/*
	RFC 6749 - The OAuth 2.0 Authorization Framework
	4.3.2.  Access Token Request
	https://tools.ietf.org/html/rfc6749#section-4.3.2
	> grant_type
	>       REQUIRED.  Value MUST be set to "password".
	>
	> username
	>       REQUIRED.  The resource owner username.
	>
	> password
	>       REQUIRED.  The resource owner password.
*/
/* eslint-disable camelcase */
export const getTokenWithPassword =
	(username, password) => ajax("/api/v0/token", "POST", null,
		{grant_type: "password", username, password});
/* eslint-enable camelcase */

/**
	getTokenWithRefreshToken returns a token authorized with the given
	refresh token.
	@function
	@param {!String} token - The refresh token.
	@returns {!module:private/promise} A promise resolved with a token.
*/
/*
	6.  Refreshing an Access Token
	https://tools.ietf.org/html/rfc6749#section-6
	> grant_type
	>       REQUIRED.  Value MUST be set to "refresh_token".
	>
	> refresh_token
	>       REQUIRED.  The refresh token issued to the client.
*/
/* eslint-disable camelcase */
export const getTokenWithRefreshToken =
	token => ajax("/api/v0/token", "POST", null,
		{grant_type: "refresh_token", refresh_token: token});
/* eslint-enable camelcase */

/**
	clubDetail returns the details of the club identified with the given ID.
	@function
	@param {!String} token - The access token.
	@param {!String} id - The ID.
	@returns {!module:private/promise} A promise resolved with the details.
*/
export const clubDetail =
	(token, id) => ajax("/api/v0/club/detail", "GET", token, {id});

/**
	clubList returns the clubs.
	@function
	@param {!String} token - The access token.
	@returns {!module:private/promise} A promise resolved with the clubs.
*/
export const clubList = token => ajax("/api/v0/club/list", "GET", token);

/**
	clubListnames returns the names of clubs associated with IDs.
	@function
	@returns {!module:private/promise} A promise resolved with the club
	names.
*/
export const clubListnames = () => ajax("/api/v0/club/listnames", "GET");

/**
	mail sends an email.
	@function
	@param {!String} token - The access token.
	@param {!*} properties - The properties of the email.
	@returns {!module:private/promise} A promise describing the progress
	and the result.
*/
export const mail =
	(token, properties) => ajax("/api/v0/mail", "POST", token, properties);

/**
	memberCreate creates a new member.
	@function
	@param {!String} token - The access token.
	@param {!*} properties - The properties of the new member.
	@returns {!module:private/promise} A promise describing the progress
	and the result.
*/
export const memberCreate =
	(token, properties) => ajax("/api/v0/member/create", "POST", token, properties);

/**
	memberDetail returns the details of the member identified with the given
	ID.
	@function
	@param {!String} token - The access token.
	@param {!String} id - The ID.
	@returns {!module:private/promise} A promise resolved with the details.
*/
export const memberDetail =
	(token, id) => ajax("/api/v0/member/detail", "GET", token, {id});

/**
	memberDelete deletes a member identified with the given ID.
	@function
	@param {!String} token - The access token.
	@param {!String} id - The ID.
	@returns {!module:private/promise} A promise describing the result.
*/
export const memberDelete =
	(token, id) => ajax("/api/v0/member/delete", "POST", token, {id});

/**
	memberList returns the members.
	@function
	@param {!String} token - The access token.
	@returns {!module:private/promise} A promise resolved with the members.
*/
export const memberList = token => ajax("/api/v0/member/list", "GET", token);

/**
	memberListroles lists the roles of the members.
	@function
	@param {!String} token - The access token.
	@returns {!module:private/promise} A promise resolved with the roles of
	the members.
*/
export const memberListroles =
	token => ajax("/api/v0/member/listroles", "GET", token);

/**
	officerDetail returns the details of the officer identified with the
	given ID.
	@function
	@param {!String} token - The access token.
	@param {!String} id - The ID.
	@returns {!module:private/promise} A promise resolved with the details.
*/
export const officerDetail =
	(token, id) => ajax("/api/v0/officer/detail", "GET", token, {id});

/**
	officerList returns the officers.
	@function
	@param {!String} token - The access token.
	@returns {!module:private/promise} A promise resolved with the officers.
*/
export const officerList = token => ajax("/api/v0/officer/list", "GET", token);

/**
	partyCreate creates a party.
	@function
	@param {!String} token - The access token.
	@param {!*} properties - Properties of the party.
	@returns {!module:private/promise} A promise describing the result.
*/
export const partyCreate =
	(token, properties) => ajax("/api/v0/party/create", "POST", token, properties);

/**
	partyList lists the parties.
	@function
	@param {!String} token - The access token.
	@returns {!module:private/promise} A promise resolved with the parties.
*/
export const partyList = token => ajax("/api/v0/party/list", "GET", token);

/**
	partyListnames lists the names of the parties.
	@function
	@param {!String} token - The access token.
	@returns {!module:provate/promise} A promise resolved with the names of
	the parties.
*/
export const partyListnames =
	token => ajax("/api/v0/party/listnames", "GET", token);

/**
	TODO
*/
export const partyRespond =
	(token, properties) =>
		ajax("/api/v0/party/respond", "GET", token, properties);

/**
	userConfirm confirms the email address of the user by submitting the
	token sent to the address.
	@function
	@param {!String} token - The access token.
	@param {!String} mailToken - The token sent to the address.
	@returns {!module:private/promise} A promise describing the progress
	and the result.
*/
export const userConfirm =
	(token, mailToken) =>
		ajax("/api/v0/user/confirm", "POST",
			token, {token: mailToken});

/**
	userDeclareOB declares the user is an OB.
	@function
	@param {!String} token - The access token.
	@returns {!module:private/promise} A promise describing the progress
	and the result.
*/
export const userDeclareOB =
	token => ajax("/api/v0/user/declareob", "POST", token);

/**
	userDetail returns the details of the user.
	@function
	@param {!String} token - The access token.
	@param {!String} id - The ID.
	@returns {!module:private/promise} A promise resolved with the details.
*/
export const userDetail =
	token => ajax("/api/v0/user/detail", "GET", token);

/**
	userUpdate returns the result of updating properties of the user.
	@function
	@param {!String} token - The access token.
	@param {!*} properties - The properties to update.
	@returns {!module:private/promise} A promise resolved with the result.
*/
export const userUpdate =
	(token, properties) =>
		ajax("/api/v0/user/update", "POST", token, properties);

/**
	userUpdatePassword returns the result of updating the user password.
	@function
	@param {!String} token - The access token.
	@param {!*} properties - The properties to update.
	@returns {!module:private/promise} A promise resolved with the result.
*/
export const userUpdatePassword =
	(token, properties) =>
		ajax("/api/v0/user/updatepassword", "POST", token, properties);
