/**
	@file index.js is the main file of client module.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0
*/

/** @module private/client */

import * as api from "./api";
import Session from "./session";

/**
	session is a variable to hold the current Session.
	@type !module:session
*/
const session = new Session;

/**
	clubDetail returns the details of the club identified with the given ID.
	@param {!external:ES.String} id - The ID.
	@returns {!external:jQuery.$.Deferred#promise} A promise resolved with
	the details.
*/
export function clubDetail(id) {
	return session.applyToken(token => api.clubDetail(token, id));
}

/**
	clubList returns the clubs.
	@returns {!external:jQuery.$.Deferred#promise} A promise resolved with
	the clubs.
*/
export function clubList() {
	return session.applyToken(token => api.clubList(token));
}

/**
	clubListName returns the names of clubs associated with IDs.
	@returns {!external:jQuery.$.Deferred#promise} A promise resolved with
	the clubs.
*/
export function clubListName() {
	return session.applyToken(token => api.clubListName(token));
}

/**
	error returns the error message corresponding to jqXHR.
	@param   {!external:jQuery~jqXHR} - A failed jqXHR.
	@returns {!external:ES.String} The error message.
*/
export const error = api.error;

/**
	getID returns the ID of the member bound to the current session.
	@returns {!external:ES.String} The ID.
*/
export const getID = session.getID.bind(session);

/**
	getScope returns the scope of the current session.
	@returns {!external:ES.String[]} The scope.
*/
export const getScope = session.getScope.bind(session);

/**
	TODO
*/
export function memberCreate(properties) {
	return session.applyToken(token => api.memberCreate(token, properties));
}

/**
	memberDeclareOB declares the user is OB.
	@returns {!external:jQuery.$.Deferred#promise} A promise describing
	the progress and the result.
*/
export function memberDeclareOB() {
	return session.applyToken(token => api.memberDeclareOB(token));
}

/**
	memberDetail returns the details of the member identified with the given
	ID.
	@param {!external:ES.String} id - The ID.
	@returns {!external:jQuery.$.Deferred#promise} A promise resolved with
	the details.
*/
export function memberDetail(id) {
	return session.applyToken(token => api.memberDetail(token, id));
}

/**
	memberDelete deletes a member identified with the given ID.
	@param {!external:ES.String} id - The ID.
	@returns {!external:jQuery.$.Deferred#promise} A promise describing the
	result.
*/
export function memberDelete(id) {
	return session.applyToken(token => api.memberDelete(token, id));
}

/**
	memberList returns the members.
	@returns {!external:jQuery.$.Deferred#promise} A promise resolved with
	the members.
*/
export function memberList() {
	return session.applyToken(token => api.memberList(token));
}

/**
	memberUpdate returns the result of updating properties of the user
	as a member.
	@param {!external:ES.Object} properties - The properties to update.
	@returns {!external:jQuery.$.Deferred#promise} A promise resolved with
	the result.
*/
export function memberUpdate(properties) {
	return session.applyToken(token => api.memberUpdate(token, properties));
}

/**
	officerDetail returns the details of the officer identified with the
	given ID.
	@param {!external:ES.String} id - The ID.
	@returns {!external:jQuery.$.Deferred#promise} A promise resolved with
	the details.
*/
export function officerDetail(id) {
	return session.applyToken(token => api.officerDetail(token, id));
}

/**
	officerList returns the officers.
	@returns {!external:jQuery.$.Deferred#promise} A promise resolved with
	the officers.
*/
export function officerList() {
	return session.applyToken(token => api.officerList(token));
}

/**
	recoverSession recovers the session from sessionStorage.
	@returns {!external:jQuery.$.Deferred#promise} A promise resolved when
	recovered.
*/
export function recoverSession() {
	return session.recover();
}

/**
	signin signs in.
	@param {!external:ES.String} id - The ID.
	@param {!external:ES.String} password - The password.
	@returns {!external:jQuery.$.Deferred#promise} A promise resolved when
	signed in.
*/
export function signin(id, password) {
	return session.signin(id, password);
}
