/**
	@file clubs.js implements the club component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0
*/

/** @module private/components/app/clubs */

/**
	module:private/components/app/clubs is a component to provide the
	feature to show clubs.
	@name module:private/components/app/clubs
	@type external:Mithril~Component
*/

import * as client from "../../client";
import * as container from "../container";
import * as progress from "../progress";
import {officers} from "../table";

export class controller {
	constructor() {
		client.clubList().then(clubs => {
			this.clubs = clubs;
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
	return [
		control.progress && m(progress, control.progress),
		m(container, m("div", {className: "container"},
			m("h1", {style: {fontSize: "x-large"}}, "Clubs"),
			control.error && m("div", {
				className: "alert alert-danger",
				role:      "alert",
			},
				m("span", {ariaHidden: "true"},
					m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
					" "
				), control.error
			),
			control.clubs && m("div", control.clubs.map(club => m("div",
				m("h2", {style: {fontSize: "large"}},
					control.clubs.name
				), m("div",
					m("h3", "部長"),
					m(officers, {
						members:     [club.chief],
						onloadstart: control.startProgress.bind(control),
						onloadend:   control.endProgress.bind(control),
						onprogress:  control.updateProgress.bind(control),
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
