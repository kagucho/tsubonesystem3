/**
	@file officer.js implements the feature to show officers.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0
*/

/** @module officers */

/**
	module:private/components/app/officers is a component to show officers.
	@name module:private/components/app/officers
	@type external:Mithril~Component
*/

import * as client from "../../client";
import * as container from "../container";
import * as progress from "../progress";
import * as table from "../table";

export class controller {
	constructor() {
		const progressSum = {
			officers: {total: 1, loaded: 0},
			clubs:    {total: 1, loaded: 0},
			control:  this,

			reflect() {
				this.control.updateProgress({
					max:   this.officers.total + this.clubs.total,
					value: this.officers.loaded + this.clubs.loaded,
				});

				m.redraw();
			},
		};

		client.officerList().then(officers => {
			this.officers = officers;
			m.redraw();
		}, this.fail.bind(this), (function(ratio) {
			this.officers = ratio;
			this.reflect();
		}).bind(progressSum));

		client.clubList().then(clubs => {
			this.clubs = clubs;
			m.redraw();
		}, this.fail.bind(this), (function(ratio) {
			this.clubs = ratio;
			this.reflect();
		}).bind(progressSum));

		progressSum.reflect();
	}

	fail(xhr) {
		this.error = client.error(xhr) || "どうしようもないエラーが発生しました。";
		this.endProgress();
		m.redraw();
	}

	startProgress() {
		this.progress = {value: 0};
	}

	updateProgress(latest) {
		this.progress = latest;
	}

	updateProgressWithXHR(event) {
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
		m(progress, control.progress),
		m(container, m("div", {className: "container"},
			m("h1", "Officers"),
			control.error && m("div", {
				className: "alert alert-danger",
				role:      "alert",
			},
				m("span", {ariaHidden: "true"},
					m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
					" "
				), control.error
			),
			m("div",
				control.officers && control.officers.map(officer => m("div",
					m("h2", {
						className: "text-center",
						style:     {fontSize: "x-large"},
					}, m("a", {href: "#!officer?id="+officer.id},
						officer.name)),
					m(table.officers, {
						members:    [officer.member],
						onprogress: control.updateProgressWithXHR.bind(control),
					})
				)),
				control.clubs && m("div",
					m("h2", {
						className: "text-center",
						style:     {fontSize: "x-large"},
					}, "各部長"),
					m("div", control.clubs.map(club => m("div",
						m("h3", {
							className: "lead",
							style:     {marginBottom: "0"},
						}, m("a", {href: "#!club?id="+club.id},
							club.name)),
						m(table.officers, {
							members:     [club.chief],
							onloadstart: control.startProgress.bind(control),
							onloadend:   control.endProgress.bind(control),
							onprogress:  control.updateProgressWithXHR.bind(control),
						})
					)))
				)
			)
		)),
	];
}
