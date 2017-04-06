/**
	@file officer.js implements officer component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/component/app/officer */

/**
	module:private/component/app/officer is a component to show an officer.
	@name module:private/component/app/officer
	@type !external:Mithril~Component
*/

import * as container from "../container";
import * as progress from "../../progress";
import client from "../../client";
import {officers} from "../table";

/**
	setOfficersState sets the state of officers table.
	@param {!external:Mithril~Node} node - The node of officers table.
	@returns {Undefined}
*/
function setOfficersState(node) {
	this.officersState = node.state;
}

export function oninit() {
	this.stream = client.mapOfficer(m.route.param("id"), promise => {
		const loadingProgress = progress.add({
			"aria-describedby": "component-app-officer-loading",
			value:              0,
		});

		promise.then(officer => {
			this.officer = officer;
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

export function onbeforeremove() {
	if (this.officersState && this.officersState.modal) {
		this.officersState.modal.remove();
	}

	if (this.stream) {
		this.stream.end(true);
	}
}

export function view() {
	return [
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
				m("section",
					m("h2", "権限"),
					m("ul", this.officer.scope.map(scope => m("li", {
						management: "メンバー情報を更新できる",
						privacy:    "メンバーの電話番号を閲覧できる",
					}[scope])))),
				m(officers, {
					members:     [this.officer.member],
					oncreate:    setOfficersState.bind(this),
				})))),
		m("div", {
			"aria-hidden": (this.stream().state() != "pending").toString(),
			id:            "component-app-officer-loading",
			style:         {display: "none"},
		}, "読み込み中…"),
	];
}
