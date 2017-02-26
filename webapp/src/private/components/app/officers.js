/**
	@file officers.js implements officers component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/app/officers */

/**
	module:private/components/app/officers is a component to show officers.
	@name module:private/components/app/officers
	@type !external:Mithril~Component
*/

import * as container from "../container";
import * as progress from "../progress";
import * as table from "../table";
import ProgressSum from "../../../progress_sum";
import client from "../../client";

/**
	reflectXHR reflects the progress and the result of the XHR.
	@private
	@this module:private/components/app/officers
	@param {!external:jQuery~Promise} promise - A promise describing the
	XHR.
	@returns {Undefined}
*/
function reflectXHR(promise) {
	this.progress.add(promise.catch(xhr => {
		this.error = client.error(xhr) || "どうしようもないエラーが発生しました。";
	}));
}

/**
	loadClubs loads clubs.
	@private
	@this module:private/components/app/officers
	@returns {Undefined}
*/
function loadClubs() {
	reflectXHR.call(this, client.clubList().then(clubs => {
		this.clubs = clubs;
	}));
}

/**
	loadOfficers loads officers.
	@private
	@this module:private/components/app/officers
	@returns {Undefined}
*/
function loadOfficers() {
	reflectXHR.call(this, client.officerList().then(officers => {
		this.officers = officers;
	}));
}

/**
	setClubsState sets the state of clubs table.
	@private
	@this module:private/components/app/officers
	@param {!external:Mithril~Node} node - The node of clubs table.
	@returns {Undefined}
*/
function setClubsState(node) {
	this.clubsState = node.state;
}

/**
	setOfficersState sets the state of officers table.
	@private
	@this module:private/components/app/officers
	@param {!external:Mithril~Node} node - The node of officers table.
	@returns {Undefined}
*/
function setOfficersState(node) {
	this.officersState = node.state;
}

export function oninit() {
	this.progress = new ProgressSum;
	loadOfficers.call(this);
	loadClubs.call(this);
}

export function onbeforeremove() {
	if (this.clubsState.modal) {
		this.clubsState.modal.remove();
	}

	if (this.officersState.modal) {
		this.officersState.modal.remove();
	}
}

export function view() {
	return [
		m(progress, this.progress.html()),
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
				this.officers && this.officers.map(officer => m("div",
					m("h2", {
						className: "text-center",
						style:     {fontSize: "x-large"},
					}, m("a", {href: "#!officer?id="+officer.id},
						officer.name)),
					m(table.officers, {
						members:     [officer.member],
						oncreate:    setOfficersState.bind(this),
						onloadstart: promise => this.progress.add(promise.done(
								submission => submission && loadOfficers.call(this))),
					})
				)),
				this.clubs && m("div",
					m("h2", {
						className: "text-center",
						style:     {fontSize: "x-large"},
					}, "各部長"),
					m("div", this.clubs.map(club => m("div",
						m("h3", {
							className: "lead",
							style:     {marginBottom: "0"},
						}, m("a", {href: "#!club?id="+club.id},
							club.name)),
						m(table.officers, {
							members:     [club.chief],
							oncreate:    setClubsState.bind(this),
							onloadstart: promise => this.progress.add(promise.done(
								submission => submission && loadClubs.call(this))),
						})
					)))
				)
			)
		)),
	];
}
