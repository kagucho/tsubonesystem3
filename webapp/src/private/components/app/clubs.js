/**
	@file clubs.js implements clubs component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/app/clubs */

/**
	module:private/components/app/clubs is a component to show clubs.
	@name module:private/components/app/clubs
	@type !external:Mithril~Component
*/

import * as container from "../container";
import * as progress from "../progress";
import ProgressSum from "../../../progress_sum";
import client from "../../client";
import {officers} from "../table";

/**
	load loads the remote content.
	@private
	@returns {Undefined}
*/
function load() {
	this.progress.add(client.clubList().then(clubs => {
		this.clubs = clubs;
	}, xhr => {
		this.error = client.error(xhr) || "どうしようもないエラーが発生しました。";
	}));
}

/**
	setOfficersState sets the state of officers table.
	@private
	@this module:private/components/app/clubs
	@param {!external:Mithril~Node} node - The node of officers table.
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
	if (this.officersState.modal) {
		this.officersState.modal.remove();
	}
}

export function view() {
	return [
		m(progress, this.progress.html()),
		m(container, m("div", {className: "container"},
			m("h1", {style: {fontSize: "x-large"}}, "Clubs"),
			this.error && m("div", {
				className: "alert alert-danger",
				role:      "alert",
			},
				m("span", {"aria-hidden": "true"},
					m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
					" "
				), this.error
			),
			this.clubs && m("div", this.clubs.map(club => m("div",
				m("h2", {style: {fontSize: "large"}},
					this.clubs.name
				), m("div",
					m("h3", "部長"),
					m(officers, {
						members:     [club.chief],
						oncreate:    setOfficersState.bind(this),
						onloadstart: promise => this.progress.add(promise.done(
							submission => submission && load.call(this))),
					})
				), m("h2", {style: {fontSize: "large"}},
					m("a", {href: "#!club?id="+club.id},
						club.name+"のいかれた仲間たちを見る"
					)
				)
			)))
		)),
	];
}
