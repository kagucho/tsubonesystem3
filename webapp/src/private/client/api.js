/**
	@file api.js implements a WebAPI interface.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/client/api */

/**
	ajax returns fetched and decoded JSON.
	@private
	@param {!String} uri - The URI.
	@param {!String} method - The method.
	@param {?String} [token] - The token.
	@param {?*} [data] - The data.
	@returns {!module:private/promise} A promise resolved with fetched JSON. TODO
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
	let uploadProgress = {loaded: 0, total: 0};

	jqXHR.then(deferred.resolve,
		xhr => deferred.reject(
			xhr == 0 ?
				"network_error" :
				xhr.responseJSON && xhr.responseJSON.error));

	xhr.upload.onprogress = deferred.notify;

	return deferred;
}

/**
	error returns the error message corresponding to jqXHR.
	@function
	@param   {!external:jQuery~jqXHR} - A failed jqXHR.
	@returns {!String} The error message.
*/
export const error = id => ({
	/* eslint-disable camelcase */
	network_error:     "TsuboneSystemへの経路上に問題が発生しました。ネットワーク接続などを確認してください。",
	invalid_grant:     "あんた誰?って言われちゃいました。もう一度サインインしてください。",
	not_found:         "見つからないってよ",
	too_many_requests: "残念！！やりすぎです。ちょっと待ってください。",
	server_error:      "サーバー側のエラーです。がびーん。",
	/* eslint-enable camelcase */
}[id] || "どうしようもないエラーです");

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
	clubDetail returns the details of the club identified with the given ID. TODO
	@function
	@param {!String} token - The access token.
	@param {!String} id - The ID.
	@returns {!module:private/promise} A promise resolved with the details.
*/
export const getClub =
	(token, id) => ajax("/api/v0/club/" + id, "GET", token);

/**
	clubList returns the clubs.
	@function
	@param {!String} token - The access token.
	@returns {!module:private/promise} A promise resolved with the clubs.
*/
export const getClubs = token => ajax("/api/v0/clubs", "GET", token);

/**
	clubMapnames returns the names of clubs associated with IDs. TODO
	@function
	@returns {!module:private/promise} A promise resolved with the club
	names.
*/
export const getClubsNames = () => ajax("/api/v0/clubs/names", "GET");

/**
	mailCreate creates an email. TODO
	@function
	@param {!String} token - The access token.
	@param {!*} properties - The properties of the email.
	@returns {!module:private/promise} A promise describing the progress
	and the result.
*/
export const putMail =
	(token, subject, properties) => ajax("/api/v0/mail/" + subject, "PUT", token, properties);

/**
	TODO
*/
export const getMail =
	(token, subject) => ajax("/api/v0/mail" + subject, "GET", token);

/**
	TODO
*/
export const getMails = token => ajax("/api/v0/mails", "GET", token);

/**
	TODO
*/
export function patchMember(token, id, properties) {
	let url = "/api/v0/member";

	if (id) {
		url += "/" + id;
	}

	return ajax(url, "PATCH", token, properties);
};

/**
	memberCreate creates a new member. TODO
	@function
	@param {!String} token - The access token.
	@param {!*} properties - The properties of the new member.
	@returns {!module:private/promise} A promise describing the progress
	and the result.
*/
export function putMember(token, id, properties) {
	let url = "/api/v0/member";

	if (id) {
		url += "/" + id;
	}

	return ajax(url, "PUT", token, properties);
};

/**
	memberDetail returns the details of the member identified with the given
	ID. TODO
	@function
	@param {!String} token - The access token.
	@param {!String} id - The ID.
	@returns {!module:private/promise} A promise resolved with the details.
*/
export function getMember(token, id) {
	let url = "/api/v0/member";

	if (id) {
		url += "/" + id;
	}

	return ajax(url, "GET", token);
}

/**
	memberDelete deletes a member identified with the given ID.
	@function
	@param {!String} token - The access token.
	@param {!String} id - The ID.
	@returns {!module:private/promise} A promise describing the result.
*/
export const deleteMember =
	(token, id) => ajax("/api/v0/member/" + id, "DELETE", token);

/**
	TODO
*/
export const getMembersMails =
	token => ajax("/api/v0/members/mails", "GET", token);

/**
	members returns the members. TODO
	@function
	@param {!String} token - The access token.
	@returns {!module:private/promise} A promise resolved with the members.
*/
export const getMembers = token => ajax("/api/v0/members", "GET", token);

/**
	officerDetail returns the details of the officer identified with the
	given ID.
	@function
	@param {!String} token - The access token.
	@param {!String} id - The ID.
	@returns {!module:private/promise} A promise resolved with the details.
*/
export const getOfficer =
	(token, id) => ajax("/api/v0/officer/" + id, "GET", token);

/**
	officerList returns the officers. TODO
	@function
	@param {!String} token - The access token.
	@returns {!module:private/promise} A promise resolved with the officers.
*/
export const getOfficers = token => ajax("/api/v0/officers", "GET", token);

/**
	partyList lists the parties. TODO
	@function
	@param {!String} token - The access token.
	@returns {!module:private/promise} A promise resolved with the parties.
*/
export const getParties = token => ajax("/api/v0/parties", "GET", token);

/**
	TODO
*/
export const patchParty =
	(token, name, properties) => ajax("/api/v0/party/" + name,
		"PATCH", token, properties);

/**
	partyCreate creates a party.
	@function
	@param {!String} token - The access token.
	@param {!*} properties - Properties of the party.
	@returns {!module:private/promise} A promise describing the result.
*/
export const putParty =
	(token, name, properties) => ajax("/api/v0/party/" + name, "PUT", token, properties);
