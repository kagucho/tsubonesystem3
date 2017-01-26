/**
	@file api.js implements a WebAPI interface.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/client/api */

import {ensureRedraw} from "../mithril";

/**
	ajax returns fetched and decoded JSON.
	@param {!external:ES.String} uri - The URI.
	@param {!external:ES.String} method - The method.
	@param {?external:ES.String} token - The token.
	@param {?external:ES.Object} data - The data.
	@returns {!external:jQuery.$.Deferred#promise} A promise resolved with
	fetched JSON.
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

	xhr.onprogress = deferred.notify.bind(deferred);

	jqXHR.then(deferred.resolve, deferred.reject);

	return ensureRedraw(deferred);
}

/**
	error returns the error message corresponding to jqXHR.
	@param   {!external:jQuery~jqXHR} - A failed jqXHR.
	@returns {!external:ES.String} The error message.
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
	@param {!external:ES.String} username - The ID.
	@param {!external:ES.String} password - The password.
	@returns {!external:jQuery.$.Deferred#promise} A promise resolved with a
	token.
*/
export function getTokenWithPassword(username, password) {
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
	return ajax("/api/v0/token", "POST", null,
              {grant_type: "password", username, password});
	/* eslint-enable camelcase */
}

/**
	getTokenWithRefreshToken returns a token authorized with the given
	refresh token.
	@param {!external:ES.String} token - The refresh token.
	@returns {!external:jQuery.$.Deferred#promise} A promise resolved with a
	token.
*/
export function getTokenWithRefreshToken(token) {
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
	return ajax("/api/v0/token", "POST", null,
              {grant_type: "refresh_token", refresh_token: token});
	/* eslint-enable camelcase */
}

/**
	clubDetail returns the details of the club identified with the given ID.
	@param {!external:ES.String} token - The access token.
	@param {!external:ES.String} id - The ID.
	@returns {!external:jQuery.$.Deferred#promise} A promise resolved with
	the details.
*/
export function clubDetail(token, id) {
	return ajax("/api/v0/club/detail", "GET", token, {id});
}

/**
	clubList returns the clubs.
	@param {!external:ES.String} token - The access token.
	@returns {!external:jQuery.$.Deferred#promise} A promise resolved with the
	clubs.
*/
export function clubList(token) {
	return ajax("/api/v0/club/list", "GET", token);
}

/**
	clubListName returns the names of clubs associated with IDs.
	@param {!external:ES.String} token - The access token.
	@returns {!external:jQuery.$.Deferred#promise} A promise resolved with
	the club names.
*/
export function clubListName(token) {
	return ajax("/api/v0/club/listname", "GET", token);
}

/**
	TODO
*/
export function memberCreate(token, properties) {
	return ajax("/api/v0/member/create", "POST", token, properties);
}

/**
	memberDetail returns the details of the member identified with the given
	ID.
	@param {!external:ES.String} token - The access token.
	@param {!external:ES.String} id - The ID.
	@returns {!external:jQuery.$.Deferred#promise} A promise resolved with
	the details.
*/
export function memberDetail(token, id) {
	return ajax("/api/v0/member/detail", "GET", token, {id});
}

/**
	TODO
*/
export function memberDeclareOB(token) {
	return ajax("/api/v0/member/declareob", "POST", token);
}

/**
	memberDelete deletes a member identified with the given ID.
	@param {!external:ES.String} token - The access token.
	@param {!external:ES.String} id - The ID.
	@returns {!external:jQuery.$.Deferred#promise} A promise describing the
	result.
*/
export function memberDelete(token, id) {
	return ajax("/api/v0/member/delete", "POST", token, {id});
}

/**
	memberUpdate returns the result of updating properties of the user
	as a member.
	@param {!external:ES.String} token - The access token.
	@param {!external:ES.Object} properties - The properties to update.
	@returns {!external:jQuery.$.Deferred#promise} A promise resolved with
	the result.
*/
export function memberUpdate(token, properties) {
	return ajax("/api/v0/member/update", "POST", token, properties);
}

/**
	memberUpdatePassword returns the result of updating the user password.
	@param {!external:ES.String} token - The access token.
	@param {!external:ES.Object} properties - The properties to update.
	@returns {!external:jQuery.$.Deferred#promise} A promise resolved with
	the result.
*/
export function memberUpdatePassword(token, properties) {
	return ajax("/api/v0/member/updatepassword", "POST", token, properties);
}

/**
	memberList returns the members.
	@param {!external:ES.String} token - The access token.
	@returns {!external:jQuery.$.Deferred#promise} A promise resolved with
	the members.
*/
export function memberList(token) {
	return ajax("/api/v0/member/list", "GET", token);
}

/**
	officerDetail returns the details of the officer identified with the
	given ID.
	@param {!external:ES.String} token - The access token.
	@param {!external:ES.String} id - The ID.
	@returns {!external:jQuery.$.Deferred#promise} A promise resolved with
	the details.
*/
export function officerDetail(token, id) {
	return ajax("/api/v0/officer/detail", "GET", token, {id});
}

/**
	officerList returns the officers.
	@param {!external:ES.String} token - The access token.
	@returns {!external:jQuery.$.Deferred#promise} A promise resolved with
	the officers.
*/
export function officerList(token) {
	return ajax("/api/v0/officer/list", "GET", token);
}
