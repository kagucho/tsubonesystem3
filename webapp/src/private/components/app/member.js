/**
	@file member.js implements the feature to show a member.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/app/member */

/**
	module:private/components/app/member is a component to show a member.
	@name module:private/components/app/member
	@type external:Mithril~Component
*/

import * as container from "../container";
import * as primitive from "../member/primitive";
import * as progress from "../progress";
import ProgressSum from "../../../progress_sum";

export function controller() {
	const control = {
		primitive: new primitive.controller,
		progress:  new ProgressSum,
	};

	control.primitive.updateAttributes({
		id:             m.route.param("id"), leaveOnInvalid: true,

		onloadstart: (function(promise) {
			this.progress.add(promise.then(submitted => {
				if (submitted) {
					delete this.error;
				}
			}, message => {
				this.error = message;
			}));
		}).bind(control),
	});

	return control;
}

export function view(control) {
	return [
		m(progress, control.progress.html()),
		m(container,
			m("div", {style: {textAlign: "center"}},
				m(control.primitive.editable ? "form" : "div", {
					style: {
						display:   "inline-block",
						textAlign: "left",
					},
				},
					m("div", {
						ariaHidden: (control.error == null).toString(),
						style:      {minHeight: "8rem"},
					},
						control.error && m("div", {
							className: "alert alert-danger",
							role:      "alert",
						}, primitive.errorView(control.error))),
					m("div", {style: {float: "right"}},
						primitive.buttonView(control.primitive)),
					m("h1", {style: {fontSize: "x-large"}},
						primitive.headerView(control.primitive)),
					primitive.bodyView(control.primitive)
				)
			)
		),
		primitive.modalView(control.primitive),
	];
}
