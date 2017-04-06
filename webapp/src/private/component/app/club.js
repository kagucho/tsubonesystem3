/**
	@file club.js implements club component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/component/app/club */

/**
	module:private/component/app/club is a component to show a club.
	@name module:private/component/app/club
	@type !external:Mithril~Component
*/

import * as container from "../container";
import * as progress from "../../progress";
import * as table from "../table";
import client from "../../client";

/**
	TODO
*/
function loadChief() {
	if (this.chiefStream) {
		this.chiefStream.end(true);
	}

	this.chiefStream = client.mapMember(this.club.chief, promise => {
		const loadingProgress = progress.add({
			"aria-describedby": "component-app-club-loading-chief",
			value:              0,
		});

		promise.then(detail => {
			this.club.chiefDetail = detail;
			loadingProgress.remove();
			m.redraw();
		}, error => {
			this.error = client.error(error);
			loadingProgress.remove();
			m.redraw();
		}, event => loadingProgress.updateValue(
			{max: event.total, value: event.loaded}));

		m.redraw();
	});
}

/**
	TODO
*/
function loadClub() {
	if (this.clubStream) {
		this.clubStream.end(true);
	}

	this.clubStream = client.mapClub(m.route.param("id"), promise => {
		const loadingProgress = progress.add({
			"aria-describedby": "component-app-club-loading-club",
			value:              0,
		});

		promise.then(club => {
			this.club = club;
			loadingProgress.remove();
			loadChief.call(this);
			m.redraw();
		}, error => {
			this.error = client.error(error);
			loadingProgress.remove();
			m.redraw();
		}, event => loadingProgress.updateValue(
			{max: event.total, value: event.loaded}));

		m.redraw();
		return promise;
	});
}

/**
	TODO
*/
function loadMembers() {
	if (this.membersStream) {
		this.membersStream.end(true);
	}

	this.membersStream = client.mapMembers(promise => {
		const loadingProgress = progress.add({
			"aria-describedby": "component-app-club-loading-members",
			value:              0,
		});

		promise.then(members => {
			loadingProgress.remove();
		}, error => {
			this.error = client.error(error);
			loadingProgress.remove();
			m.redraw();
		}, event => loadingProgress.updateValue(
			{max: event.total, value: event.loaded}));

		m.redraw();
		return promise.then(members => members());
	});
}

/**
	setMembersState sets the state of members table.
	@private
	@this module:private/component/app/club
	@param {!external:Mithril~Node} node - The node of created members
	table.
	@returns {Undefined}
*/
function setMembersState(node) {
	this.membersState = node.state;
}

/**
	setOfficersState sets the state of officers table.
	@private
	@this module:private/component/app/club
	@param {!external:Mithril~Node} node - The node of created officers
	table.
	@returns {Undefined}
*/
function setOfficersState(node) {
	this.officersState = node.state;
}

export function oninit() {
	loadClub.call(this);
	loadMembers.call(this);

	client.merge(this.clubStream, this.membersStream).map(
		promise => promise.done(
			(club, members) => {
				this.members = Array.from((function *() {
					for (const member of members) {
						if (club.members.has(member.id)) {
							yield member;
						}
					}
				})());

				m.redraw();
			}));
}

export function onbeforeremove() {
	if (this.chiefStream) {
		this.chiefStream.end(true);
	}

	if (this.clubStream) {
		this.clubStream.end(true);
	}

	if (this.membersStream) {
		this.membersStream.end(true);
	}

	if (this.membersState.modal) {
		this.membersState.modal.remove();
	}

	if (this.officersState.modal) {
		this.officersState.modal.remove();
	}
}

export function view() {
	return [
		m(container, m("div", {className: "container"},
			this.error && m("div", {
				className: "alert alert-danger",
				role:      "alert",
			},
				m("span", {"aria-hidden": "true"},
					m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
					" "
				), this.error),
			this.club && m("section",
				m("h1", {style: {fontSize: "x-large"}},
					this.club.name+"の詳細情報"),
				m("section",
					m("h2", {style: {fontSize: "large"}},
						"部長"),
					this.club.chiefDetail && m(table.officers, {
						members:     [this.club.chiefDetail],
						oncreate:    setOfficersState.bind(this),
					})),
				m("section",
					m("h2", {style: {fontSize: "large"}},
						this.club.name+"のいかれた仲間たち"),
					m("p", {style: {color: "gray"}},
						this.club.members.length+" 件"),
					this.members && m(table.members, {
						members:     this.members,
						oncreate:    setMembersState.bind(this),
					}))))),
		m("div", {
			"aria-hidden": (!this.chiefStream || this.chiefStream().state() != "pending").toString(),
			id:            "component-app-club-loading-chief",
			style:         {display: "none"},
		}, "部長の情報を読み込んでいます…"),
		m("div", {
			"aria-hidden": (!this.clubStream || this.clubStream().state() != "pending").toString(),
			id:            "component-app-club-loading-club",
			style:         {display: "none"},
		}, "部の情報を読み込んでいます…"),
		m("div", {
			"aria-hidden": (!this.membersStream || this.membersStream().state() != "pending").toString(),
			id:            "component-app-club-loading-members",
			style:         {display: "none"},
		}, "局員の情報を読み込んでいます…"),
	];
}
