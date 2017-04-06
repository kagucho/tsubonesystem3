/**
	@file clubs.js implements clubs component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/component/app/clubs */

/**
	module:private/component/app/clubs is a component to show clubs.
	@name module:private/component/app/clubs
	@type !external:Mithril~Component
*/

import * as container from "../container";
import * as progress from "../../progress";
import client from "../../client";
import {officers} from "../table";

function loadChiefs() {
	if (this.chiefStreams) {
		unloadChiefs();
	}

	this.chiefStreams = Array.from(this.clubs,
		club => client.mapMember(club.chief,
			promise => promise.done(
				detail => {
					club.chiefDetail = detail;
					m.redraw();
				})));

	client.merge(...this.chiefStreams).map(promise => {
		this.chiefLoading = "component-app-club-loading-chief";

		const loadingProgress = progress.add({
			"aria-describedby": this.chiefLoading,
			value:              0,
		});

		promise.then(() => {
			this.chiefLoading = null;
			loadingProgress.remove();
			m.redraw();
		}, error => {
			this.error = client.error(error);
			this.chiefLoading = null;
			loadingProgress.remove();
			m.redraw();
		}, event => loadingProgress.updateValue(
			{max: event.total, value: event.loaded}));

		m.redraw();
	});
}

function unloadChiefs() {
	for (const stream of this.chiefStreams) {
		stream.end(true);
	}
}

/**
	setOfficersState sets the state of officers table.
	@private
	@this module:private/component/app/clubs
	@param {!external:Mithril~Node} node - The node of officers table.
	@returns {Undefined}
*/
function setOfficersState(node) {
	this.officersState = node.state;
}

export function oninit() {
	this.clubsStream = client.mapClubs(promise => {
		this.clubLoading = "component-app-clubs-loading-club";

		const loadingProgress = progress.add(
			{"aria-describedby": this.clubLoading, value: 0});

		promise.then(clubs => {
			this.clubLoading = null;
			this.clubs = Array.from(clubs());
			loadChiefs.call(this);
			loadingProgress.remove();
			m.redraw();
		}, error => {
			this.clubLoading = null;
			this.error = client.error(error);
			loadingProgress.remove();
			m.redraw();
		}, event => loadingProgress.updateValue(
			{max: event.total, value: event.loaded}));

		m.redraw();
	});
}

export function onbeforeremove() {
	this.clubsStream.end(true);
	unloadChiefs.call(this);

	if (this.officersState.modal) {
		this.officersState.modal.remove();
	}
}

export function view() {
	return [
		m(container, m("div", {className: "container"},
			m("h1", "Clubs"),
			this.error && m("div", {
				className: "alert alert-danger",
				role:      "alert",
			},
				m("span", {"aria-hidden": "true"},
					m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
					" "),
				this.error),
			this.clubs && m("div", this.clubs.map(club => m("section",
				m("h2", club.name),
				m("section",
					m("h3", "部長"),
					club.chiefDetail && m(officers, {
						members:     [club.chiefDetail],
						oncreate:    setOfficersState.bind(this),
					})),
				m("div", {className: "h4"},
					m("a", {href: "#!club?id="+club},
						club.name + "のいかれた仲間たちを見る"))))))),
		m("div", {
			"aria-hidden": (!this.chiefLoading).toString(),
			id:            this.chiefLoading,
			style:         {display: "none"},
		}, "部長の情報を読み込んでいます…"),
		m("div", {
			"aria-hidden": (!this.clubLoading).toString(),
			id:            this.clubLoading,
			style:         {display: "none"},
		}, "部の情報を読み込んでいます…"),
	];
}
