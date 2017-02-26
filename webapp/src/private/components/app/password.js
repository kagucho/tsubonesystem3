/**
	@file password.js implements password component.
	@author Akihiko Odaki <akihiko.odaki.4i@stu.hosei.ac.jp>
	@copyright 2017  {@link https://kagucho.net/|Kagucho}
	@license AGPL-3.0+
*/

/** @module private/components/app/password */

/**
	module:private/components/app/password is a component to update the
	password of the user.
	@name module:private/components/app/password
	@type !external:Mithril~Component
*/

import * as alert from "../alert";
import * as container from "../container";
import * as modal from "../../modal";
import * as primitive from "../password/primitive";
import * as progress from "../progress";

/**
	registerLoading registers a loading.
	@private
	@param {!external:jQuery~Promise} promise - A promise describing the
	loading.
	@returns {Undefined}
*/
function registerLoading(promise) {
	this.message = {
		attrs: {value: 0},
		type:  "inprogress",
	};

	const inprogress = modal.unshift(
		alert.inprogress(primitive.inprogress));

	promise.then(submission => {
		inprogress.remove();

		modal.unshift((submission ? alert.leavable : alert.closable)(
			m("span", {"aria-hidden": "true"},
				m("span", {className: "glyphicon glyphicon-ok"}),
				" "
			), primitive.success
		));
	}, xhr => {
		inprogress.remove();

		modal.unshift(alert.closable(
			m("span", {"aria-hidden": "true"},
				m("span", {className: "glyphicon glyphicon-exclamation-sign"}),
				" "
			), primitive.error(xhr)
		));
	}, event => this.progress = {
		max:   event.total,
		value: event.loaded,
	});
}

export function oninit() {
	this.primitive = primitive.newState();
	this.primitive.setOnloadstart(registerLoading.bind(this));
}

export function view() {
	return [
		this.progress && m(progress, this.progress),
		m(container,
			m("div", {className: "container", style: {textAlign: "center"}},
				m("form", {
					"aria-labelledby": "component-app-password-title",
					style:             {
						display:   "inline-block",
						textAlign: "left",
					},
				},
					m("h1", {id: "component-app-password-title"},
						primitive.title
					), m(this.primitive.body, {
						autofocus: true,
						oncreate:  () => setTimeout(
							this.primitive.focus.bind(this.primitive)),
					}), m("div", {style: {textAlign: "center"}},
						m(this.primitive.button)
					)
				)
			)
		),
	];
}
