/**
	@file officer.js implements the feature to show an officer.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/app/officer */

/**
	module:private/components/app/officer is a component to show an officer.
	@name module:private/components/app/officer
	@type external:Mithril~Component
*/

import * as client from "../../client";
import * as container from "../container";
import * as progress from "../progress";
import ProgressSum from "../../../progress_sum";
import {officers} from "../table";

export class controller {
	constructor() {
		this.progress = new ProgressSum;
		this.load();
	}

	load() {
		this.progress.add(client.officerDetail(m.route.param("id")).then(officer => {
			this.officer = officer;
		}, xhr => {
			this.error = client.error(xhr) || "どうしようもないエラーが発生しました。";
		}));
	}
}

export function view(control) {
	return [
		m(progress, control.progress.html()),
		m(container, m("div", {className: "container"},
			control.error && m("div", {
				className: "alert alert-danger", role:      "alert",
			},
				m("span", {ariaHidden: "true"},
					m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
					" "
				), control.error
			),
			control.officer && m("div",
				m("h1", control.officer.name + "閣下の詳細情報"),
				m("div",
					m("h2", "権限"),
					m("ul", control.officer.scope.map(scope => m("li", {
						management: "メンバー情報を更新できる",
						privacy:    "メンバーの電話番号を閲覧できる",
					}[scope])))
				),
				m(officers, {
					members:     [control.officer.member],
					onloadstart: (function(promise) {
						this.progress.add(promise.done(
							submission => submission && this.load()));
					}).bind(control),
				})
			)
		)),
	];
}
