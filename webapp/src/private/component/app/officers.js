/**
	@file officers.js implements officers component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/component/app/officers */

/**
	module:private/component/app/officers is a component to show officers.
	@name module:private/component/app/officers
	@type !external:Mithril~Component
*/

import * as container from "../container";
import * as progress from "../../progress";
import * as table from "../table";
import client from "../../client";

/**
	TODO
*/
function showXHRError(error) {
	this.error = client.error(error);
}

/**
	TODO
*/
function loadChiefMembers() {
	if (this.chiefStreams) {
		for (const stream of this.chiefStreams) {
			stream.end(true);
		}
	}

	this.chiefStreams = this.clubs.map(
		club => client.mapMember(club.chief,
			promise => promise.done(
				detail => {
					club.chiefDetail = detail;
					m.redraw();
				})));

	client.merge(...this.chiefStreams).map(promise => {
		this.loadingChiefMembers = "component-app-club-loading-chief-members";

		const loadingProgress = progress.add({
			"aria-describedby": this.loadingChiefMembers,
			value:              0,
		});

		promise.then(() => {
			this.loadingChiefMembers = null;
			loadingProgress.remove();
			m.redraw();
		}, error => {
			this.loadingChiefMembers = null;
			showXHRError(error);
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
function loadOfficerMembers() {
	if (this.officerStreams) {
		for (const stream of this.officerStreams) {
			stream.end(true);
		}
	}

	this.officerStreams = this.officers.map(
		officer => client.mapMember(officer.member,
			promise => promise.done(
				detail => {
					officer.memberDetail = detail;
					m.redraw();
				})));

	client.merge(...officerStreams).map(promise => {
		this.loadingOfficerMembers = "component-app-club-loading-officer-members";

		const loadingProgress = progress.add({
			"aria-describedby": this.loadingOfficerMembers,
			value:              0,
		});

		promise.then(() => {
			this.loadingOfficerMembers = null;
			loadingProgress.remove();
			m.redraw();
		}, error => {
			showXHRError(error);

			this.loadingOfficerMembers = null;
			loadingProgress.remove();
			m.redraw();
		}, event => loadingProgress.updateValue(
			{max: event.total, value: event.loaded}));

		m.redraw();
	});
}

/**
	loadClubs loads clubs.
	@private
	@this module:private/component/app/officers
	@returns {Undefined}
*/
function loadClubs() {
	if (this.clubsStream) {
		this.clubsStream.end(true);
	}

	this.clubsStream = client.mapClubs(promise => {
		this.loadingClubs = "component-app-officers-loading";

		const loadingProgress = progress.add(
			{"aria-describedby": this.loadingClubs, value: 0});

		promise.then(clubs => {
			this.clubs = Array.from(clubs());
			this.loadingClubs = null;
			loadChiefMembers.call(this);
			loadingProgress.remove();
			m.redraw();
		}, error => {
			this.loadingClubs = null;
			loadingProgress.remove();
			showXHRError(error);
			m.redraw();
		}, event => loadingProgress.updateValue(
			{max: event.total, value: event.loaded}));

		m.redraw();
	});
}

/**
	loadOfficers loads officers.
	@private
	@this module:private/component/app/officers
	@returns {Undefined}
*/
function loadOfficers() {
	if (this.officersStream) {
		this.officersStream.end(true);
	}

	this.officersStream = client.mapOfficers(promise => {
		this.loadingOfficers = "component-app-officers-loading";

		const loadingProgress = progress.add(
			{"aria-describedby": this.loadingOfficers, value: 0});

		promise.then(officers => {
			this.loadingOfficers = null;
			this.officers = Array.from(officers());
			loadOfficerMembers.call(this);
			loadingProgress.remove();
			m.redraw();
		}, error => {
			this.loadingOfficers = null;
			loadingProgress.remove();
			showXHRError(error);
			m.redraw();
		}, event => loadingProgress.updateValue(
			{max: event.total, value: event.loaded}));

		m.redraw();
	});
}

/**
	setClubsState sets the state of clubs table.
	@private
	@this module:private/component/app/officers
	@param {!external:Mithril~Node} node - The node of clubs table.
	@returns {Undefined}
*/
function setClubsState(node) {
	this.clubsState = node.state;
}

/**
	setOfficersState sets the state of officers table.
	@private
	@this module:private/component/app/officers
	@param {!external:Mithril~Node} node - The node of officers table.
	@returns {Undefined}
*/
function setOfficersState(node) {
	this.officersState = node.state;
}

export function oninit() {
	loadOfficers.call(this);
	loadClubs.call(this);
}

export function onbeforeremove() {
	for (const streams of [
		this.chiefStreams, this.officerStreams,
		[this.clubsStream, this.officersStream],
	]) {
		if (streams) {
			for (const stream of streams) {
				if (stream) {
					stream.end(true);
				}
			}
		}
	}

	if (this.clubsState.modal) {
		this.clubsState.modal.remove();
	}

	if (this.officersState.modal) {
		this.officersState.modal.remove();
	}
}

export function view() {
	return [
		m(container, m("div", {className: "container"},
			m("h1", "Officers"),
			this.error && m("div", {
				className: "alert alert-danger",
				role:      "alert",
			},
				m("span", {"aria-hidden": "true"},
					m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
					" "
				), this.error
			),
			m("div",
				this.officers && this.officers.map(officer => m("section",
					m("h2", {
						className: "text-center",
						style:     {fontSize: "x-large"},
					}, m("a", {href: "#!officer?id="+officer.id},
						officer.name)),
					officer.memberDetail && m(table.officers, {
						members:     [officer.memberDetail],
						oncreate:    setOfficersState.bind(this),
					})
				)),
				this.clubs && m("section",
					m("h2", {
						className: "text-center",
						style:     {fontSize: "x-large"},
					}, "各部長"),
					this.clubs.map(club => m("section",
						m("h3", {
							className: "lead",
							style:     {marginBottom: "0"},
						}, m("a", {href: "#!club?id="+club.id},
							club.name)),
						club.chiefDetail && m(table.officers, {
							members:     [club.chiefDetail],
							oncreate:    setClubsState.bind(this),
						}))))))),
		m("div", {
			"aria-hidden": (!this.loadingChiefMembers).toString(),
			id:            this.loadingChiefMembers,
			style:         {display: "none"},
		}, "部長の情報を読み込んでいます"),
		m("div", {
			"aria-hidden": (!this.loadingClubs).toString(),
			id:            this.loadingClubs,
			style:         {display: "none"},
		}, "部の情報を読み込んでいます…"),
		m("div", {
			"aria-hidden": (!this.loadingOfficerMembers).toString(),
			id:            this.loadingOfficerMembers,
			style:         {display: "none"},
		}, "えらい人たちの詳細情報を読み込んでいます"),
		m("div", {
			"aria-hidden": (!this.loadingOfficers).toString(),
			id:            this.loadingOfficers,
			style:         {display: "none"},
		}, "えらい人たちの情報を読み込んでいます…"),
	];
}
