/**
	@file member.js implements the feature to show a member.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0
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

export function controller() {
	const control = {primitive: new primitive.controller()};
	control.primitive.updateAttributes({
		id:             m.route.param("id"), leaveOnInvalid: true,

		onerror: (function(message) {
			this.message = {body: message, type: "error"};
		}).bind(control),

		onsuccess: (function(message) {
			this.message = {body: message, type: "success"};
		}).bind(control),

		onloadstart: (function() {
			this.progress = {value: 0};
		}).bind(control),

		onloadend: (function() {
			if (this.progress.value != this.progress.max) {
				delete this.progress;
			}
		}).bind(control),

		onprogress: (function(event) {
			this.progress = {max: event.total, value: event.loaded};
		}).bind(control),
	});

	return control;
}

export function view(control) {
	return [
		control.progress && m(progress, control.progress),
		m(container,
			m("div", {style: {textAlign: "center"}},
				m(control.primitive.editable ? "form" : "div", {
					style: {
						display:   "inline-block",
						textAlign: "left",
					},
				},
					m("div", {
						ariaHidden: (control.message == null).toString(),
						style:      {minHeight: "8rem"},
					},
						control.message && m("div", {
							className: control.message && {
								error:   "alert alert-danger",
								success: "alert alert-success",
							}[control.message.type],
							role: "alert",
						}, control.message && {
							error:   primitive.errorView,
							success: primitive.successView,
						}[control.message.type](control.message.body))),
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
