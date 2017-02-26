/**
	@file officer.js implements officer component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/app/officer */

/**
	module:private/components/app/officer is a component to show an officer.
	@name module:private/components/app/officer
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
	@this module:private/components/app/officer
	@returns {Undefined}
*/
function load() {
	this.progress.add(client.officerDetail(m.route.param("id")).then(officer => {
		this.officer = officer;
	}, xhr => {
		this.error = client.error(xhr) || "どうしようもないエラーが発生しました。";
	}));
}

/**
	setOfficersState sets the state of officers table.
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
			this.error && m("div", {
				className: "alert alert-danger", role:      "alert",
			},
				m("span", {"aria-hidden": "true"},
					m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
					" "
				), this.error
			),
			this.officer && m("div",
				m("h1", this.officer.name + "閣下の詳細情報"),
				m("div",
					m("h2", "権限"),
					m("ul", this.officer.scope.map(scope => m("li", {
						management: "メンバー情報を更新できる",
						privacy:    "メンバーの電話番号を閲覧できる",
					}[scope])))
				),
				m(officers, {
					members:     [this.officer.member],
					oncreate:    setOfficersState.bind(this),
					onloadstart: promise => this.progress.add(promise.done(
						submission => submission && load.call(this))),
				})
			)
		)),
	];
}
