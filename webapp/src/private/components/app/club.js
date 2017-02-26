/**
	@file club.js implements club component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/app/club */

/**
	module:private/components/app/club is a component to show a club.
	@name module:private/components/app/club
	@type !external:Mithril~Component
*/

import * as container from "../container";
import * as progress from "../progress";
import * as table from "../table";
import ProgressSum from "../../../progress_sum";
import client from "../../client";

/**
	load loads the remote content.
	@private
	@this module:private/components/app/club
	@returns {Undefined}
*/
function load() {
	this.progress.add(client.clubDetail(m.route.param("id")).then(club => {
		this.club = club;
	}, xhr => {
		this.error = client.error(xhr) || "どうしようもないエラーが発生しました。";
	}));
}

/**
	registerLoading registers a promise describing a loading.
	@private
	@this module:private/components/app/club
	@param {!external:jQuery~Promise} - The promise describing a loding.
	@returns {Undefined}
*/
function registerLoading(promise) {
	this.progress.add(promise.done(submission => submission && load.call(this)));
}

/**
	setMembersState sets the state of members table.
	@private
	@this module:private/components/app/club
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
	@this module:private/components/app/club
	@param {!external:Mithril~Node} node - The node of created officers
	table.
	@returns {Undefined}
*/
function setOfficersState(node) {
	this.officersState = node.state;
}

export function oninit() {
	this.progress = new ProgressSum;
	load.call(this);
}

export function onbeforeremove() {
	if (this.membersState.modal) {
		this.membersState.modal.remove();
	}
	if (this.officersState.modal) {
		this.officersState.modal.remove();
	}
}

export function view() {
	const boundRegisterLoading = registerLoading.bind(this);

	return [
		m(progress, this.progress.html()),
		m(container, m("div", {className: "container"},
			this.error && m("div", {
				className: "alert alert-danger",
				role:      "alert",
			},
				m("span", {"aria-hidden": "true"},
					m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
					" "
				), this.error
			), this.club && m("div",
				m("h1", {style: {fontSize: "x-large"}},
					this.club.name+"の詳細情報"),
				m("div",
					m("h2", {style: {fontSize: "large"}},
						"部長"),
					m(table.officers, {
						members:     [this.club.chief],
						oncreate:    setOfficersState.bind(this),
						onloadstart: boundRegisterLoading,
					})
				), m("div",
					m("h2", {style: {fontSize: "large"}},
						this.club.name+"のいかれた仲間たち"),
					m("p", {style: {color: "gray"}},
						this.club.members.length+" 件"),
					m(table.members, {
						members:     this.club.members,
						oncreate:    setMembersState.bind(this),
						onloadstart: boundRegisterLoading,
					})
				)
			)
		)),
	];
}
