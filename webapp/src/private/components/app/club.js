/**
	@file club.js implements the club component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0
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

export class controller {
	constructor() {
		client.clubDetail(m.route.param("id")).then(club => {
			this.club = club;
			m.redraw();
		}, xhr => {
			this.error = client.error(xhr) || "どうしようもないエラーが発生しました。";
			this.endProgress();
			m.redraw();
		}, event => {
			this.updateProgress(event);
			m.redraw();
		});

		this.startProgress();
	}

	startProgress() {
		this.progress = {value: 0};
	}

	updateProgress(event) {
		this.progress = {max: event.total, value: event.loaded};
	}

	endProgress() {
		if (this.progress.value != this.progress.max) {
			delete this.progress;
		}
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
		const startProgress = control.startProgress.bind(control);
		const updateProgress = control.updateProgress.bind(control);
		const endProgress = control.endProgress.bind(control);

		content.push(m("div",
			m("h1", {style: {fontSize: "x-large"}},
				control.club.name+"の詳細情報"),
			m("div",
				m("h2", {style: {fontSize: "large"}},
					"部長"),
				m(table.officers, {
					members:     [control.club.chief],
					onloadstart: startProgress,
					onloadend:   endProgress,
					onprogress:  updateProgress,
				})
			), m("div",
				m("h2", {style: {fontSize: "large"}},
					control.club.name+"のいかれた仲間たち"),
				m("p", {style: {color: "gray"}},
					control.club.members.length+" 件"),
				m(table.members, {
					members:     control.club.members,
					onloadstart: startProgress,
					onloadend:   endProgress,
					onprogress:  updateProgress,
				})
			)
		));
	}

	return [
		control.progress && m(progress, control.progress),
		m(container, m("div", {className: "container"}, content)),
	];
}
