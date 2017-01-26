/**
	@file officer.js implements the feature to show officers.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
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
import ProgressSum from "../../../progress_sum";

export class controller {
	constructor() {
		this.progress = new ProgressSum;
		this.loadOfficers();
		this.loadClubs();
	}

	loadClubs() {
		this.reflectXHR(client.clubList().then(clubs => {
			this.clubs = clubs;
		}));
	}

	loadOfficers() {
		this.reflectXHR(client.officerList().then(officers => {
			this.control.officers = officers;
		}));
	}

	reflectXHR(promise) {
		this.progress.add(promise.catch(xhr => {
			this.error = client.error(xhr) || "どうしようもないエラーが発生しました。";
		}));
	}
}

export function view(control) {
	return [
		m(progress, control.progress.html()),
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
						members:     [officer.member],
						onloadstart: (function(promise) {
							this.progress.add(promise.done(
								submission => submission && this.loadOfficers()));
						}).bind(control),
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
							onloadstart: (function(promise) {
								this.progress.add(promise.done(
									submission => submission && this.loadClubs()));
							}).bind(control),
						})
					)))
				)
			)
		)),
	];
}
