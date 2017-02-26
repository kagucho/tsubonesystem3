/**
	@file client.js provides the integrated feature of the client.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/client/client */

import * as api from "./api";
import Session from "./session";

/**
	sessions is a WeakMap to hold sessions.
	@private
	@type !WeakMap<!module:private/client/client, !module:private/client/session>
*/
const sessions = new WeakMap;

/**
	module:private/client/client is a class to provide the integrated
	feature of the client.
	@extends Object
*/
export default class {
	/**
		constructor constructs module:private/client/client.
		@returns {Undefined}
	*/
	constructor() {
		sessions.set(this, new Session);
		Object.freeze(this);
	}

	/**
		clubDetail returns the details of the club identified with the given ID.
		@param {!String} id - The ID.
		@returns {!module:private/promise} A promise resolved with the details.
	*/
	clubDetail(id) {
		return sessions.get(this).applyToken(
			token => api.clubDetail(token, id));
	}

	/**
		clubList returns the clubs.
		@returns {!module:private/promise} A promise resolved with the clubs.
	*/
	clubList() {
		return sessions.get(this).applyToken(api.clubList);
	}

	/**
		clubListnames returns the names of clubs associated with IDs.
		@returns {!module:private/promise} A promise resolved with the clubs.
	*/
	static clubListnames() {
		return api.clubListnames();
	}

	/**
		error returns the error message corresponding to jqXHR.
		@param   {!external:jQuery~jqXHR} - A failed jqXHR.
		@returns {!String} The error message.
	*/
	static error(xhr) {
		return api.error(xhr);
	}

	/**
		getFilling returns whether the user agent is prompting the user
		to fill his information.
		@returns {!Boolean} The boolean indicating whether the user
		agent is prompting the user to fill his information.
	*/
	getFilling() {
		return sessions.get(this).getFilling();
	}

	/**
		getID returns the ID of the member bound to the current session.
		@returns {!String} The ID.
	*/
	getID() {
		return sessions.get(this).getID();
	}

	/**
		getScope returns the scope of the current session.
		@returns {!String[]} The scope.
	*/
	getScope() {
		return sessions.get(this).getScope();
	}

	/**
		mail sends a email.
		@returns {!module:private/promise} A promise desribing the
		progress and the result.
	*/
	mail(properties) {
		return sessions.get(this).applyToken(
			token => api.mail(token, properties));
	}

	/**
		memberCreate creates a member.
		@param {!*} properties - The properties of the new member.
		@returns {!module:private/promise} A promis describing the
		progress and the result.
	*/
	memberCreate(properties) {
		return sessions.get(this).applyToken(
			token => api.memberCreate(token, properties));
	}

	/**
		memberDetail returns the details of the member identified with
		the given ID.
		@param {!String} id - The ID.
		@returns {!module:private/promise} A promise resolved with the
		details.
	*/
	memberDetail(id) {
		return sessions.get(this).applyToken(
			token => api.memberDetail(token, id));
	}

	/**
		memberDelete deletes a member identified with the given ID.
		@param {!String} id - The ID.
		@returns {!module:private/promise} A promise describing the
		result.
	*/
	memberDelete(id) {
		return sessions.get(this).applyToken(
			token => api.memberDelete(token, id));
	}

	/**
		memberList returns the members.
		@returns {!module:private/promise} A promise resolved with the members.
	*/
	memberList() {
		return sessions.get(this).applyToken(api.memberList);
	}

	/**
		memberListroles returns the roles of the members.
		@returns {!module:private/promise} A promise resolved with the
		roles of the members.
	*/
	memberListroles() {
		return sessions.get(this).applyToken(api.memberListroles);
	}

	/**
		officerDetail returns the details of the officer identified with
		the given ID.
		@param {!String} id - The ID.
		@returns {!module:private/promise} A promise resolved with the
		details.
	*/
	officerDetail(id) {
		return sessions.get(this).applyToken(
			token => api.officerDetail(token, id));
	}

	/**
		officerList returns the officers.
		@returns {!module:private/promise} A promise resolved with the officers.
	*/
	officerList() {
		return sessions.get(this).applyToken(api.officerList);
	}

	/**
		partyCreate creates a party.
		@param {!*} properties - Properties of the party.
		@returns {!module:private/promise} A promise describing the
		result.
	*/
	partyCreate(properties) {
		return sessions.get(this).applyToken(
			token => api.partyCreate(token, properties));
	}

	/**
		partyList lists the parties.
		@returns {!module:private/promise} A promise resolved with the
		parties.
	*/
	partyList() {
		return sessions.get(this).applyToken(api.partyList);
	}

	/**
		partyListnames lists the names of the parties.
		@returns {!module:provate/promise} A promise resolved with the
		names of the parties.
	*/
	partyListnames() {
		return sessions.get(this).applyToken(api.partyListnames);
	}

	/**
		TODO
	*/
	partyRespond(properties) {
		return sessions.get(this).applyToken(
			token => api.partyRespond(token, properties));
	}

	/**
		recoverSession recovers the session from sessionStorage.
		@returns {!module:private/promise} A promise resolved when
		recovered.
	*/
	recoverSession() {
		return sessions.get(this).recover();
	}

	/**
		setFillingToken sets the token to fill the information of the
		member.
		@param {!String} id - The ID of the user.
		@param {!String} token - The access token.
		@returns {Undefined}
	*/
	setFillingToken(id, token) {
		sessions.get(this).setFillingToken(id, token);
	}

	/**
		signin signs in.
		@param {!String} id - The ID.
		@param {!String} password - The password.
		@returns {!module:private/promise} A promise describing the
		progress and the result.
	*/
	signin(id, password) {
		return sessions.get(this).signin(id, password);
	}

	/**
		userConfirm confirms the email address of the user by submitting
		the token sent to the address.
		@param {!String} mailToken - The token sent to the address.
		@returns {!module:private/promise} A promise describing the
		progress and the result.
	*/
	userConfirm(mailToken) {
		return sessions.get(this).applyToken(
			token => api.userConfirm(token, mailToken));
	}

	/**
		userDeclareOB declares the user is OB.
		@returns {!module:private/promise} A promise describing the
		progress and the result.
	*/
	userDeclareOB() {
		return sessions.get(this).applyToken(api.userDeclareOB);
	}

	/**
		userDetail returns the details of the user.
		@returns {!module:private/promise} A promise resolved with the
		details.
	*/
	userDetail() {
		return sessions.get(this).applyToken(api.userDetail);
	}

	/**
		userUpdate updates properties of the user.
		@param {!*} properties - The properties to update.
		@returns {!module:private/promise} A promise describing the
		progress and the result.
	*/
	userUpdate(properties) {
		return sessions.get(this).applyToken(
			token => api.userUpdate(token, properties));
	}

	/**
		userUpdatePassword updates the password of the user.
		@param {!*} properties - The properties to update.
		@returns {!module:private/promise} A promise describing the
		progress and the result.
	*/
	userUpdatePassword(properties) {
		return sessions.get(this).applyToken(
			token => api.userUpdatePassword(token, properties));
	}
}
