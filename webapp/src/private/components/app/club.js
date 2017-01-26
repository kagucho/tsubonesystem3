/**
	@file club.js implements the club component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/app/club */

/**
	module:private/components/app/club is a component to provide the feature
	to show a club.
	@name module:private/components/app/club
	@type external:Mithril~Component
*/

import * as client from "../../client";
import * as container from "../container";
import * as progress from "../progress";
import * as table from "../table";
import ProgressSum from "../../../progress_sum";

export class controller {
	constructor() {
		this.progress = new ProgressSum;
		this.load();
	}

	load() {
		this.progress.add(client.clubDetail(m.route.param("id")).then(club => {
			this.club = club;
		}, xhr => {
			this.error = client.error(xhr) || "どうしようもないエラーが発生しました。";
		}));
	}
}

export function view(control) {
	const content = [];

	if (control.error) {
		content.push(m("div", {
			className: "alert alert-danger", role:      "alert",
		},
			m("span", {ariaHidden: "true"},
				m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
				" "
			), control.error
		));
	}

	if (control.club) {
		const onloadstart = (function(promise) {
			this.progress.add(promise.done(
				submission => submission && this.load()));
		}).bind(control);

		content.push(m("div",
			m("h1", {style: {fontSize: "x-large"}},
				control.club.name+"の詳細情報"),
			m("div",
				m("h2", {style: {fontSize: "large"}},
					"部長"),
				m(table.officers, {
					members: [control.club.chief],
					onloadstart,
				})
			), m("div",
				m("h2", {style: {fontSize: "large"}},
					control.club.name+"のいかれた仲間たち"),
				m("p", {style: {color: "gray"}},
					control.club.members.length+" 件"),
				m(table.members, {
					members: control.club.members,
					onloadstart,
				})
			)
		));
	}

	return [
		m(progress, control.progress.html()),
		m(container, m("div", {className: "container"}, content)),
	];
}
